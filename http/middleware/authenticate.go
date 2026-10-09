package middleware

import (
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/support"
)

// Authenticate 后台鉴权中间件。
//
// 免鉴权路由对齐 owl-admin Permission::$authExcept：
// login / logout / no-content / _settings(GET) / captcha。
// 未登录返回 HTTP 200 + {status:1, code:401, doNotDisplayToast:1}，
// 前端拦截器据响应体 code==401 跳转登录页。
type Authenticate struct{}

func NewAuthenticate() *Authenticate { return &Authenticate{} }

// Signature 中间件签名。
func (r *Authenticate) Signature() string { return "admin.auth" }

// Handle 执行鉴权检查：先 Parse token（goravel 要求），
// 失败/未携带时返回 code=401 响应体。
func (r *Authenticate) Handle(ctx contractshttp.Context) {
	prefix := facades.Config().GetString("admin.route.prefix", "admin-api")
	path := strings.Trim(strings.TrimPrefix(ctx.Request().Path(), "/"+prefix), "/")

	if r.inExcept(path, ctx.Request().Method()) {
		ctx.Request().Next()

		return
	}

	token := ctx.Request().Header("Authorization", "")
	if token == "" {
		ctx.Response().Json(contractshttp.StatusOK, support.Unauthorized("请先登录")).Abort()

		return
	}

	// Parse 校验签名/有效期并把 claims 写入请求上下文，后续 User() 才可用
	if _, err := facades.Auth(ctx).Guard("admin").Parse(token); err != nil {
		facades.Log().Debugf("[admin] token parse failed: %v", err)
		ctx.Response().Json(contractshttp.StatusOK, support.Unauthorized("请先登录")).Abort()

		return
	}

	ctx.Request().Next()
}

var authExcept = map[string]bool{
	"login":            true,
	"logout":           true,
	"no-content":       true,
	"_settings":        true,
	"captcha":          true,
	"_download_export": true,
}

// inExcept 判断是否免鉴权；_settings 的写操作（POST/PUT/PATCH）仍需登录
// （对齐 authIntercept 中对 _settings 写方法的特判）。
func (r *Authenticate) inExcept(path, method string) bool {
	if !authExcept[path] {
		return false
	}

	if path == "_settings" {
		switch method {
		case "POST", "PUT", "PATCH":
			return false
		}
	}

	return true
}
