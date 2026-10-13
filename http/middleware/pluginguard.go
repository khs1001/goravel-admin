package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/khs1001/goravel-admin/support/owl"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// pluginGuard rejects requests to a disabled plugin's routes with 404, so a
// disabled plugin is invisible without needing a process restart.
type pluginGuard struct {
	name string
}

func (m *pluginGuard) Signature() string {
	return "admin.plugin.guard:" + m.name
}

func (m *pluginGuard) Handle(ctx http.Context) {
	if plugin.R().IsEnabled(m.name) {
		ctx.Request().Next()
		return
	}

	ctx.Response().Json(http.StatusNotFound, owl.FailSilent("插件未启用")).Abort()
}

// PluginGuard returns the middleware that gates a plugin's routes on its
// enabled state.
func PluginGuard(name string) http.Middleware {
	return &pluginGuard{name: name}
}
