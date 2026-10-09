package middleware

import (
	"encoding/json"
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support"
)

// Permission 接口权限拦截（对齐 owl-admin Support/Cores/Permission::permissionIntercept）：
// 超管放行；permissionExcept 内的路径放行；其余按用户权限的 http_path/http_method 匹配。
type Permission struct {
	prefix string
}

func NewPermission() *Permission {
	return &Permission{
		prefix: "/" + facades.Config().GetString("admin.route.prefix", "admin-api"),
	}
}

// Signature 中间件签名。
func (r *Permission) Signature() string { return "admin.permission" }

// permissionExcept 免权限检查路径（相对前缀，对齐 owl-admin Permission::$permissionExcept）。
var permissionExcept = []string{
	"menus", "current-user", "user_setting", "login", "logout",
	"no-content", "_settings", "captcha", "_download_export",
}

// Handle 执行权限检查（挂在 Authenticate 之后）。
func (r *Permission) Handle(ctx contractshttp.Context) {
	if !facades.Config().GetBool("admin.auth.permission", true) {
		ctx.Request().Next()

		return
	}

	path := strings.TrimPrefix(ctx.Request().Path(), r.prefix)
	path = strings.Trim(path, "/")

	if r.inExcept(path) {
		ctx.Request().Next()

		return
	}

	user := &models.AdminUser{}
	if err := facades.Auth(ctx).Guard("admin").User(user); err != nil || user.ID == 0 {
		ctx.Request().Next()

		return
	}

	permissionService := services.NewPermissionService()
	if permissionService.IsAdministrator(int(user.ID)) {
		ctx.Request().Next()

		return
	}

	if r.allowed(ctx, permissionService.UserPermissionIDs(int(user.ID)), path) {
		ctx.Request().Next()

		return
	}

	ctx.Response().Json(contractshttp.StatusOK, support.Fail("没有权限访问该接口")).Abort()
}

// inExcept 免检路径判断。
func (r *Permission) inExcept(path string) bool {
	for _, except := range permissionExcept {
		if matchPattern(except, path) {
			return true
		}
	}

	return false
}

// allowed 权限 http_path 匹配。
func (r *Permission) allowed(ctx contractshttp.Context, permissionIDs []int, path string) bool {
	if len(permissionIDs) == 0 {
		return false
	}

	var permissions []models.AdminPermission
	ids := make([]any, 0, len(permissionIDs))
	for _, id := range permissionIDs {
		ids = append(ids, id)
	}
	if err := facades.Orm().Query().WhereIn("id", ids).Find(&permissions); err != nil {
		facades.Log().Errorf("[admin] load permissions for check failed: %v", err)

		return false
	}

	method := strings.ToUpper(ctx.Request().Method())
	for _, permission := range permissions {
		if !methodAllowed(permission.HTTPMethod, method) {
			continue
		}
		if permission.HTTPPath == nil {
			continue
		}

		for _, line := range strings.Split(*permission.HTTPPath, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 行格式："GET /admin-api/xxx" 或 "/admin-api/xxx"（不带方法 = 全部方法）
			pattern := line
			if fields := strings.Fields(line); len(fields) == 2 {
				if !strings.EqualFold(fields[0], method) {
					continue
				}
				pattern = fields[1]
			}
			pattern = strings.Trim(strings.TrimPrefix(pattern, "/admin-api"), "/")

			if matchPattern(pattern, path) {
				return true
			}
		}
	}

	return false
}

// methodAllowed http_method JSON 数组为空表示不限方法。
func methodAllowed(httpMethod *string, method string) bool {
	if httpMethod == nil || *httpMethod == "" {
		return true
	}

	var methods []string
	if err := json.Unmarshal([]byte(*httpMethod), &methods); err != nil {
		return true
	}
	if len(methods) == 0 {
		return true
	}

	for _, m := range methods {
		if strings.EqualFold(m, method) {
			return true
		}
	}

	return false
}

// matchPattern 支持尾部 * 通配。
func matchPattern(pattern, path string) bool {
	pattern = strings.Trim(pattern, "/")
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "*"))
	}

	return pattern == path
}
