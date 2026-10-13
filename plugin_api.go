package admin

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/http/middleware"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// Plugin contract re-exports.
//
// A plugin never has to import this package: Go interfaces are structural, so
// implementing the matching method signatures is enough. These aliases exist for
// plugins that DO choose to import the admin package and want the named types.
type (
	Plugin                = contracts.Plugin
	MetadataProvider      = contracts.MetadataProvider
	MenuProvider          = contracts.MenuProvider
	RouteProvider         = contracts.RouteProvider
	PageProvider          = contracts.PageProvider
	SettingProvider       = contracts.SettingProvider
	AssetProvider         = contracts.AssetProvider
	NavProvider           = contracts.NavProvider
	CommandProvider       = contracts.CommandProvider
	MigrationProvider     = contracts.MigrationProvider
	SeederProvider        = contracts.SeederProvider
	UserAuthenticatedHook = contracts.UserAuthenticatedHook
	LoginHook             = contracts.LoginHook
	LogoutHook            = contracts.LogoutHook
	OperationLoggedHook   = contracts.OperationLoggedHook
	MenusBuiltHook        = contracts.MenusBuiltHook
	SettingsBuiltHook     = contracts.SettingsBuiltHook
)

// PluginBinding returns the container key a plugin must bind under to be
// discovered (e.g. PluginBinding("blog") == "admin.plugin.blog").
func PluginBinding(name string) string {
	return plugin.Binding(name)
}

// Permission returns a route-level middleware enforcing a permission id, for
// use in a plugin's RouteProvider.
func Permission(perm string) http.Middleware {
	return middleware.Permission(perm)
}

// Authorize reports whether the current user holds the given permission id,
// for handler-level checks inside a plugin (mirrors the built-in controllers).
func Authorize(ctx http.Context, perm string) bool {
	if !facades.Config().GetBool("admin.auth.permission", true) {
		return true
	}

	return services.NewMenuService().HasPermission(current.User(ctx), perm)
}
