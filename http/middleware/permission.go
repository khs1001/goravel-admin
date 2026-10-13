package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
)

// permissionMiddleware enforces a RuoYi permission identifier on a route.
type permissionMiddleware struct {
	perm string
}

func (m *permissionMiddleware) Signature() string {
	return "admin.permission:" + m.perm
}

func (m *permissionMiddleware) Handle(ctx http.Context) {
	if !facades.Config().GetBool("admin.auth.permission", true) {
		ctx.Request().Next()
		return
	}

	if services.NewMenuService().HasPermission(current.User(ctx), m.perm) {
		ctx.Request().Next()
		return
	}

	ctx.Response().Json(http.StatusOK, owl.FailSilent("无操作权限")).Abort()
}

// Permission returns a route-level middleware enforcing a permission id.
func Permission(perm string) http.Middleware {
	return &permissionMiddleware{perm: perm}
}
