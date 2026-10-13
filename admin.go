package admin

import (
	"github.com/goravel/framework/contracts/foundation"
)

// Admin is the main service of the package. It exposes the resolved admin
// configuration to the host through the facades.Admin facade.
type Admin struct {
	app foundation.Application
}

// NewAdmin builds the Admin service bound to the given application.
func NewAdmin(app foundation.Application) *Admin {
	return &Admin{app: app}
}

// Name returns the configured admin application name.
func (r *Admin) Name() string {
	return r.config("admin.name", "Goravel Admin")
}

// RoutePrefix returns the configured API route prefix, e.g. "admin-api".
func (r *Admin) RoutePrefix() string {
	return r.config("admin.route.prefix", "admin-api")
}

// Path returns the configured SPA entry path, e.g. "admin".
func (r *Admin) Path() string {
	return r.config("admin.path", "admin")
}

// Assets returns the static asset base path, e.g. "admin-assets".
func (r *Admin) Assets() string {
	return r.config("admin.assets", "admin-assets")
}

func (r *Admin) config(key string, def string) string {
	if r.app == nil {
		return def
	}

	config := r.app.MakeConfig()
	if config == nil {
		return def
	}

	return config.GetString(key, def)
}
