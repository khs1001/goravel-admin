package admin

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/jobs"
	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
)

// Binding admin 服务的容器绑定 key。
const Binding = "admin"

// App 应用实例（Register 阶段赋值，供 facades.Admin() 解析使用）。
var App foundation.Application

type ServiceProvider struct {
}

// Relationship 返回服务提供者的绑定关系。
func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{Binding},
		Dependencies: []string{},
		ProvideFor:   []string{},
	}
}

// Register 注册服务提供者：绑定 Admin 实现，并收集容器中所有
// 以 contracts.PluginKey 为 key 的插件绑定。
func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return NewAdmin(app), nil
	})
}

// Boot 启动：收集插件（容器 PluginKey 绑定 → contracts.Plugin），
// 注册数据迁移（扩展包自身 + 实现 contracts.Migrator 的插件）、
// 操作日志队列任务与 ORM 观察者，并注册后台全部路由。
func (r *ServiceProvider) Boot(app foundation.Application) {
	plugins := r.collectPlugins(app)

	// 手动注册（Admin.RegisterPlugin）的插件同样生效
	if instance, err := app.Make(Binding); err == nil {
		if admin, ok := instance.(*Admin); ok {
			plugins = append(plugins, admin.Plugins()...)
		}
	}

	registerMigrations(plugins)
	registerOperationLog()

	RegisterRoutes(plugins)
}

// registerMigrations 注册数据迁移。
//
// 迁移注册规范：只在扩展包内部注册，宿主应用不罗列——
// admin 包自身的迁移来自 Migrations()（单一注册源），
// 二级插件迁移来自其 contracts.Migrator 可选实现。
func registerMigrations(plugins []contracts.Plugin) {
	migrations := Migrations()
	for _, plugin := range plugins {
		migrator, ok := plugin.(contracts.Migrator)
		if !ok {
			continue
		}

		pluginMigrations := migrator.Migrations()
		facades.Log().Debugf("[admin] plugin %s registered %d migrations", plugin.Name(), len(pluginMigrations))
		migrations = append(migrations, pluginMigrations...)
	}

	facades.Schema().Register(migrations)
}

// registerOperationLog 注册操作日志 Job（Queue）与 ORM 模型观察者。
// 各表 TraceColumns 声明参与字段级 diff 的业务列（仅 updated 动作记录变更前后值）。
func registerOperationLog() {
	facades.Queue().Register([]contractsqueue.Job{&jobs.OperationLogJob{}})

	observer := func(resource string) services.OperationLogObserver {
		return services.OperationLogObserver{Resource: resource}
	}
	tracingObserver := func(resource string, columns ...string) services.OperationLogObserver {
		return services.OperationLogObserver{Resource: resource, TraceColumns: columns}
	}

	facades.Orm().Observe(&models.AdminUser{}, tracingObserver("admin_users", "username", "name", "avatar", "enabled", "password"))
	facades.Orm().Observe(&models.AdminRole{}, tracingObserver("admin_roles", "name", "slug"))
	facades.Orm().Observe(&models.AdminPermission{}, tracingObserver("admin_permissions", "name", "slug", "http_method", "http_path", "custom_order", "parent_id"))
	facades.Orm().Observe(&models.AdminMenu{}, tracingObserver("admin_menus", "title", "url", "url_type", "visible", "custom_order", "parent_id", "icon"))
	facades.Orm().Observe(&models.AdminDepartment{}, tracingObserver("admin_departments", "name", "parent_id", "custom_order", "remark"))
	facades.Orm().Observe(&models.AdminPosition{}, tracingObserver("admin_positions", "name", "custom_order", "remark"))
	facades.Orm().Observe(&models.AdminDictionary{}, tracingObserver("admin_dictionaries", "name", "key", "remark"))
	facades.Orm().Observe(&models.AdminDictionaryItem{}, tracingObserver("admin_dictionary_items", "label", "value", "custom_order", "dictionary_id"))
	facades.Orm().Observe(&models.AdminSetting{}, observer("admin_settings"))
}

// collectPlugins 枚举容器绑定：所有 PluginKey 类型的 key 视为插件注册，
// 经容器解析后断言为 contracts.Plugin。
func (r *ServiceProvider) collectPlugins(app foundation.Application) []contracts.Plugin {
	var plugins []contracts.Plugin
	for _, key := range app.Bindings() {
		pluginKey, ok := key.(contracts.PluginKey)
		if !ok {
			continue
		}

		instance, err := app.Make(pluginKey)
		if err != nil {
			facades.Log().Warningf("[admin] resolve plugin %s failed: %v", pluginKey, err)

			continue
		}

		if plugin, ok := instance.(contracts.Plugin); ok {
			plugins = append(plugins, plugin)
			facades.Log().Debugf("[admin] plugin registered: %s", plugin.Name())
		}
	}

	return plugins
}
