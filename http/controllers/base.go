package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
)

// BaseController holds helpers shared by the admin controllers.
type BaseController struct{}

// user returns the authenticated admin user (set by the auth middleware).
func (b *BaseController) user(ctx http.Context) *models.SysUser {
	return current.User(ctx)
}

// denied enforces a RuoYi permission identifier. It returns a response when the
// request must be rejected, or nil when it is allowed.
func (b *BaseController) denied(ctx http.Context, perm string) http.Response {
	if !facades.Config().GetBool("admin.auth.permission", true) {
		return nil
	}

	if services.NewMenuService().HasPermission(b.user(ctx), perm) {
		return nil
	}

	response := ctx.Response().Json(http.StatusOK, owl.FailSilent("无操作权限"))
	_ = response.Abort()

	return response
}

// page reads the amis pagination parameters.
func (b *BaseController) page(ctx http.Context) (int, int) {
	page := ctx.Request().QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage := ctx.Request().QueryInt("perPage", 20)
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	return page, perPage
}

// defaultValue returns value when non-empty, otherwise fallback.
func defaultValue(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

// parseUint parses a decimal unsigned integer, returning 0 on error.
func parseUint(value string) uint {
	parsed, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0
	}

	return uint(parsed)
}
