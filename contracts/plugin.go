package contracts

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
)

// Plugin is the base contract every admin plugin must satisfy.
//
// A plugin becomes visible to goravel-admin by binding an implementation under
// the key `admin.plugin.<name>` (see support/plugin.Binding). Because Go
// interfaces are structural, a plugin does NOT need to import this package: it
// only needs matching method signatures.
//
// Every signature below intentionally uses framework types or plain Go types
// (never goravel-admin types), so a zero-dependency plugin can satisfy the
// capability contracts without importing this module.
type Plugin interface {
	// Name is the unique plugin name. It is the binding suffix and the value
	// stored in `sys_extension`.
	Name() string
}

// MetadataProvider is optional and enriches the extension-management page.
type MetadataProvider interface {
	Version() string
	Description() string
	Author() string
}

// MenuProvider contributes menu / permission nodes. A node is a plain map with
// the keys consumed by the menu merge:
//
//	title, path, parent (parent path/name; empty = top level), order (int),
//	icon, perms, component (page sign), menuType (M/C/F).
type MenuProvider interface {
	Menus() []map[string]any
}

// RouteProvider registers routes on a sub-router that is already wrapped with
// the admin auth, operation-log and plugin-guard middleware.
type RouteProvider interface {
	Routes(router route.Router)
}

// PageProvider contributes amis page schemas keyed by a short sign; the sign is
// namespaced as "<plugin>.<short>".
type PageProvider interface {
	Pages() map[string]map[string]any
}

// SettingProvider contributes default settings (lowest precedence above the
// host config, below the database).
type SettingProvider interface {
	Settings() map[string]any
}

// AssetProvider contributes frontend assets: {"js": []string, "css": []string}.
type AssetProvider interface {
	Assets() map[string]any
}

// NavProvider contributes navbar amis schemas: {"prependNav": any, "appendNav": any}.
type NavProvider interface {
	Nav() map[string]any
}

// CommandProvider contributes console commands.
type CommandProvider interface {
	Commands() []console.Command
}

// MigrationProvider contributes migrations (registered at boot, run by the host
// `migrate`; not affected by the enable/disable gate).
type MigrationProvider interface {
	Migrations() []schema.Migration
}

// SeederProvider contributes seeders.
type SeederProvider interface {
	Seeders() []seeder.Seeder
}

// Hooks are optional and independent: implement only the ones you need. They
// are invoked (for enabled plugins, in discovery order) at the matching point.

type UserAuthenticatedHook interface {
	OnUserAuthenticated(ctx http.Context, user any)
}

type LoginHook interface {
	OnLoginSucceeded(ctx http.Context, username string)
	OnLoginFailed(ctx http.Context, username string, msg string)
}

type LogoutHook interface {
	OnLogout(ctx http.Context, username string)
}

type OperationLoggedHook interface {
	OnOperationLogged(ctx http.Context, entry map[string]any)
}

type MenusBuiltHook interface {
	OnMenusBuilt(ctx http.Context, routes []map[string]any) []map[string]any
}

type SettingsBuiltHook interface {
	OnSettingsBuilt(ctx http.Context, payload map[string]any)
}
