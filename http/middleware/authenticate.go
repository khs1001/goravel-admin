package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// Authenticate parses the JWT bearer token, resolves the admin user and stores
// it on the request context. It aborts with the Owl 401 envelope on failure,
// which makes the frontend redirect to the login page.
type Authenticate struct{}

// Signature returns the unique signature of the middleware.
func (m *Authenticate) Signature() string {
	return "admin.auth"
}

// Handle processes the incoming request.
func (m *Authenticate) Handle(ctx http.Context) {
	if !facades.Config().GetBool("admin.auth.enable", true) {
		ctx.Request().Next()
		return
	}

	token := ctx.Request().Header("Authorization")
	if token == "" {
		ctx.Response().Json(http.StatusOK, owl.Unauthorized("请先登录")).Abort()
		return
	}

	guard := facades.Auth(ctx).Guard("admin")
	if _, err := guard.Parse(token); err != nil {
		ctx.Response().Json(http.StatusOK, owl.Unauthorized("登录状态已失效，请重新登录")).Abort()
		return
	}

	var user models.SysUser
	if err := guard.User(&user); err != nil || user.UserID == 0 {
		ctx.Response().Json(http.StatusOK, owl.Unauthorized("请先登录")).Abort()
		return
	}

	current.Set(ctx, &user)
	plugin.R().RunUserAuthenticated(ctx, &user)

	ctx.Request().Next()
}
