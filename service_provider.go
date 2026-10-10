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
// observer 记录 created/deleted；updated 的字段级 diff 由各控制器
// 调用 services.DispatchChangeLog 记录（Save 模式下 ORM 事件拿不到 Original）。
func registerOperationLog() {
	facades.Queue().Register([]contractsqueue.Job{
		&jobs.OperationLogJob{},
		&jobs.LoginLogJob{},
	})

	observer := func(resource string) services.OperationLogObserver {
		return services.OperationLogObserver{Resource: resource}
	}

	facades.Orm().Observe(&models.AdminUser{}, observer("admin_users"))
	facades.Orm().Observe(&models.AdminRole{}, observer("admin_roles"))
	facades.Orm().Observe(&models.AdminPermission{}, observer("admin_permissions"))
	facades.Orm().Observe(&models.AdminMenu{}, observer("admin_menus"))
	facades.Orm().Observe(&models.AdminDepartment{}, observer("admin_departments"))
	facades.Orm().Observe(&models.AdminPosition{}, observer("admin_positions"))
	facades.Orm().Observe(&models.AdminDictionary{}, observer("admin_dictionaries"))
	facades.Orm().Observe(&models.AdminDictionaryItem{}, observer("admin_dictionary_items"))
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
