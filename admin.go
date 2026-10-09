package admin

import (
	"sync"

	"github.com/goravel/framework/contracts/foundation"

	"github.com/khs1001/goravel-admin/contracts"
)

// Admin admin 包运行时实现，经容器绑定 "admin" 解析。
type Admin struct {
	app foundation.Application

	mu      sync.RWMutex
	plugins []contracts.Plugin
}

func NewAdmin(app foundation.Application) *Admin {
	return &Admin{app: app}
}

// Name 后台应用名称。
func (r *Admin) Name() string {
	return r.app.MakeConfig().GetString("admin.name", "Goravel Admin")
}

// Logo 后台 logo 地址。
func (r *Admin) Logo() string {
	return r.app.MakeConfig().GetString("admin.logo", "")
}

// RoutePrefix 后台 API 路由前缀。
func (r *Admin) RoutePrefix() string {
	return r.app.MakeConfig().GetString("admin.route.prefix", "admin-api")
}

// RegisterPlugin 手动注册插件（除容器绑定外的补充方式）。
func (r *Admin) RegisterPlugin(plugin contracts.Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plugins = append(r.plugins, plugin)
}

// Plugins 已注册的全部插件（含手动注册 + 容器 PluginKey 绑定）。
func (r *Admin) Plugins() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]Plugin{}, r.plugins...)
}

// Plugin 是 contracts.Plugin 的包内别名，避免使用方双重导入。
type Plugin = contracts.Plugin
