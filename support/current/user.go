// Package current stores and retrieves the authenticated admin user on the
// request context, so controllers do not depend on the auth middleware package.
package current

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/khs1001/goravel-admin/models"
)

// key is the context key under which the authenticated user is stored.
const key = "goravel_admin:user"

// Set stores the authenticated user on the request context.
func Set(ctx http.Context, user *models.SysUser) {
	ctx.WithValue(key, user)
}

// User returns the authenticated user, or nil when the request is not
// authenticated.
func User(ctx http.Context) *models.SysUser {
	value := ctx.Value(key)
	if value == nil {
		return nil
	}

	user, ok := value.(*models.SysUser)
	if !ok {
		return nil
	}

	return user
}
