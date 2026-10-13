// Package demo is an example admin plugin.
//
// It is intentionally ZERO-DEPENDENCY: it never imports goravel-admin. It only
// uses the Goravel framework contracts plus plain Go types (maps), which proves
// that a plugin can be discovered and contribute functionality without being
// coupled to the admin package.
//
// Discovery works because the plugin is bound under the agreed key
// "admin.plugin.demo" and its methods structurally satisfy the admin capability
// contracts.
package demo

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
)

// BindingKey is the container key this plugin binds under. The prefix
// "admin.plugin." is the discovery convention.
const BindingKey = "admin.plugin.demo"

// Plugin implements several admin capability contracts.
type Plugin struct{}

// Name implements the base contract.
func (p *Plugin) Name() string { return "demo" }

// Version / Description / Author implement the optional metadata contract.
func (p *Plugin) Version() string     { return "0.1.0" }
func (p *Plugin) Description() string { return "示例插件（零依赖）" }
func (p *Plugin) Author() string      { return "khs1001" }

// Menus contributes menu nodes: a directory under 系统管理, and two pages.
func (p *Plugin) Menus() []map[string]any {
	return []map[string]any{
		{"title": "示例插件", "path": "/demo", "parent": "/system", "order": 9, "icon": "mdi:emoticon-happy-outline", "menuType": "M"},
		{"title": "示例首页", "path": "/demo/index", "parent": "/demo", "order": 1, "component": "index", "perms": "demo:index:list", "menuType": "C"},
		{"title": "示例详情", "path": "/demo/detail", "parent": "/demo", "order": 2, "component": "detail", "perms": "demo:index:detail", "menuType": "C"},
	}
}

// Pages contributes amis page schemas (signs are namespaced as "demo.<short>").
func (p *Plugin) Pages() map[string]map[string]any {
	return map[string]map[string]any{
		"index": {
			"type":  "page",
			"title": "示例首页",
			"body":  map[string]any{"type": "tpl", "tpl": "<div style='padding:24px'>来自 demo 插件（零依赖）的首页</div>"},
		},
		"detail": {
			"type":  "page",
			"title": "示例详情",
			"body":  map[string]any{"type": "tpl", "tpl": "<div style='padding:24px'>demo 插件详情页</div>"},
		},
	}
}

// Settings contributes a default setting (overridable in sys_config).
func (p *Plugin) Settings() map[string]any {
	return map[string]any{"demo_greeting": "hello from demo"}
}

// Assets contributes extra frontend assets (none here).
func (p *Plugin) Assets() map[string]any {
	return map[string]any{"js": []string{}, "css": []string{}}
}

// Routes registers a demo route on the guarded sub-router.
func (p *Plugin) Routes(router route.Router) {
	router.Get("/ping", func(ctx http.Context) http.Response {
		return ctx.Response().Json(http.StatusOK, http.Json{
			"status": 0, "msg": "", "doNotDisplayToast": 0,
			"data": http.Json{"pong": true},
		})
	})
}

// Commands contributes a console command.
func (p *Plugin) Commands() []console.Command {
	return []console.Command{&DemoCommand{}}
}

// DemoCommand is the plugin's demo command.
type DemoCommand struct{}

func (c *DemoCommand) Signature() string { return "demo:hello" }
func (c *DemoCommand) Description() string {
	return "示例插件命令"
}
func (c *DemoCommand) Extend() command.Extend { return command.Extend{Category: "demo"} }
func (c *DemoCommand) Handle(ctx console.Context) error {
	ctx.Success("hello from demo plugin")

	return nil
}

// ServiceProvider binds the plugin so goravel-admin can discover it.
type ServiceProvider struct{}

// Relationship returns the relationship of the service provider.
func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{},
		Dependencies: []string{binding.Config},
		ProvideFor:   []string{},
	}
}

// Register binds the plugin instance under the discovery key.
func (r *ServiceProvider) Register(app foundation.Application) {
	app.Bind(BindingKey, func(app foundation.Application) (any, error) {
		return &Plugin{}, nil
	})
}

// Boot performs post-registration booting of services.
func (r *ServiceProvider) Boot(app foundation.Application) {}
