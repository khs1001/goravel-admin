package middleware

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
)

// OperLog records write operations (POST/PUT/DELETE) into `sys_oper_log`.
//
// It must run after the auth middleware so the operator is known.
type OperLog struct{}

// Signature returns the unique signature of the middleware.
func (m *OperLog) Signature() string {
	return "admin.operlog"
}

// Handle processes the incoming request.
func (m *OperLog) Handle(ctx http.Context) {
	start := time.Now()

	ctx.Request().Next()

	method := strings.ToUpper(ctx.Request().Method())
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		return
	}

	name := ""
	if user := current.User(ctx); user != nil {
		name = user.UserName
	}

	status, msg := responseOutcome(ctx)

	services.NewLogService().Operation(models.SysOperLog{
		Title:         operTitle(ctx.Request().Path()),
		BusinessType:  businessType(method),
		RequestMethod: method,
		OperatorType:  1,
		OperName:      name,
		OperURL:       ctx.Request().Path(),
		OperIP:        ctx.Request().Ip(),
		OperParam:     paramString(ctx),
		Status:        status,
		ErrorMsg:      msg,
		OperTime:      carbon.NewDateTime(carbon.Now()),
		CostTime:      int(time.Since(start).Milliseconds()),
	})
}

// responseOutcome reads the Owl envelope from the rendered response (best
// effort) to record success/failure and the message.
func responseOutcome(ctx http.Context) (int, string) {
	origin := ctx.Response().Origin()
	if origin == nil || origin.Body() == nil {
		return 0, ""
	}

	text := strings.TrimSpace(origin.Body().String())
	if text == "" {
		return 0, ""
	}

	var envelope struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}
	if err := json.Unmarshal([]byte(text), &envelope); err != nil {
		return 0, ""
	}

	return envelope.Status, envelope.Msg
}

func paramString(ctx http.Context) string {
	params := ctx.Request().All()
	for _, key := range []string{"password", "old_password", "confirm_password"} {
		delete(params, key)
	}

	raw, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	result := string(raw)
	if len(result) > 2000 {
		result = result[:2000]
	}

	return result
}

func businessType(method string) int {
	switch method {
	case "POST":
		return models.BusinessTypeInsert
	case "PUT", "PATCH":
		return models.BusinessTypeUpdate
	case "DELETE":
		return models.BusinessTypeDelete
	default:
		return models.BusinessTypeOther
	}
}

func operTitle(path string) string {
	for prefix, title := range map[string]string{
		"/system/user":       "用户管理",
		"/system/role":       "角色管理",
		"/system/menu":       "菜单管理",
		"/system/dept":       "部门管理",
		"/system/post":       "岗位管理",
		"/system/dict":       "字典管理",
		"/system/config":     "参数管理",
		"/system/operlog":    "操作日志",
		"/system/logininfor": "登录日志",
		"/upload":            "文件上传",
		"/user_setting":      "个人设置",
	} {
		if strings.Contains(path, prefix) {
			return title
		}
	}

	return path
}
