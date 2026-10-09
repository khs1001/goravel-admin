package middleware

import (
	"encoding/json"
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/jobs"
	"github.com/khs1001/goravel-admin/models"
)

// OperationLog 操作日志中间件：记录所有写请求（POST/PUT/PATCH/DELETE），
// 经 Queue 异步落库。需挂在 Authenticate 之后（依赖已 Parse 的登录态）。
type OperationLog struct{}

func NewOperationLog() *OperationLog { return &OperationLog{} }

// Signature 中间件签名。
func (r *OperationLog) Signature() string { return "admin.operation_log" }

// Handle 记录写请求并放行。
func (r *OperationLog) Handle(ctx contractshttp.Context) {
	method := ctx.Request().Method()
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		ctx.Request().Next()

		return
	}

	payload := jobs.OperationLogPayload{
		Action: "request",
		Method: method,
		Path:   ctx.Request().Path(),
		IP:     ctx.Request().Ip(),
	}

	// 登录用户（Authenticate 已 Parse）
	user := &models.AdminUser{}
	if err := facades.Auth(ctx).Guard("admin").User(user); err == nil && user.ID != 0 {
		payload.UserID = int(user.ID)
		payload.Username = user.Username
	}

	payload.Input = sanitizeInput(ctx.Request().All())

	dispatch(payload)

	ctx.Request().Next()
}

// sanitizeInput 请求参数序列化：脱敏密码字段并截断超长内容。
func sanitizeInput(all map[string]any) string {
	if len(all) == 0 {
		return ""
	}

	sensitive := map[string]bool{
		"password": true, "confirm_password": true, "old_password": true,
		"new_password": true, "token": true,
	}
	for key := range all {
		if sensitive[strings.ToLower(key)] {
			all[key] = "***"
		}
	}

	raw, err := json.Marshal(all)
	if err != nil {
		return ""
	}

	input := string(raw)
	if len(input) > 2000 {
		input = input[:2000]
	}

	return input
}

// dispatch 推送队列（失败仅告警，不影响主流程）。
func dispatch(payload jobs.OperationLogPayload) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		facades.Log().Warningf("[admin] marshal operation log failed: %v", err)

		return
	}

	if err := facades.Queue().Job(&jobs.OperationLogJob{}, []contractsqueue.Arg{
		{Value: string(encoded), Type: "string"},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("[admin] dispatch operation log failed: %v", err)
	}
}
