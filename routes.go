package admin

import (
	contractsroute "github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/http/controllers"
	"github.com/khs1001/goravel-admin/http/middleware"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/services"
)

// RegisterRoutes 注册后台全部路由（对齐 owl-admin src/Support/Cores/Route.php 第一阶段子集）。
//
// 路由前缀 config admin.route.prefix（默认 admin-api）；
// 免鉴权基座路由 + 鉴权路由两层；插件 Routes() 挂在鉴权层。
func RegisterRoutes(plugins []contracts.Plugin) {
	prefix := facades.Config().GetString("admin.route.prefix", "admin-api")

	authController := controllers.NewAuthController()
	settingService := services.NewSettingService()
	pageService := services.NewPageService(pages.Builtin(), plugins)

	// 插件代码级菜单 → 菜单树（对齐 owl-admin Plugin::menus）
	pluginMenus := make([]contracts.MenuItem, 0, len(plugins))
	for _, plugin := range plugins {
		pluginMenus = append(pluginMenus, plugin.Menus()...)
	}
	indexController := controllers.NewIndexController(services.NewMenuService(), settingService, pageService)
	indexController.SetPluginMenus(pluginMenus)

	// 前端静态资源：public/admin-assets → /admin-assets/*（admin-views 构建产物）
	facades.Route().Static("admin-assets", "public/admin-assets")
	// 上传文件访问：storage/app/public → /storage/*（public 磁盘）
	facades.Route().Static("storage", "storage/app/public")

	facades.Route().Prefix("/" + prefix).Group(func(router contractsroute.Router) {		// 基座路由（免鉴权）
		router.Post("/login", authController.Login)
		router.Get("/logout", authController.Logout)
		router.Get("/captcha", authController.Captcha)
		router.Get("/no-content", indexController.NoContent)
		router.Get("/_settings", indexController.Settings)

		// 鉴权路由
		router.Middleware(middleware.NewAuthenticate()).Group(func(authed contractsroute.Router) {
			authed.Post("/_settings", indexController.SaveSettings)
			authed.Get("/menus", indexController.Menus)
			authed.Get("/current-user", authController.CurrentUser)
			authed.Get("/user_setting", authController.UserSetting)
			authed.Put("/user_setting", authController.SaveUserSetting)
			authed.Get("/page_schema", indexController.PageSchema)

			// 统一上传（登录即可，不挂权限）
			uploadController := controllers.NewUploadController()
			authed.Post("/upload", uploadController.Store)

			// 权限路由（RBAC 接口级拦截 + 写操作日志）
			authed.Middleware(middleware.NewPermission(), middleware.NewOperationLog()).Group(func(permitted contractsroute.Router) {
				// 内置页面路由（菜单 url_type=路由 时，前端按路径拉取 schema）
				permitted.Get("/dashboard", indexController.Dashboard)
				permitted.Get("/dashboard_stats", indexController.DashboardStats)

				// 系统管理（RBAC）：对齐 owl-admin /system 资源路由
				userController := controllers.NewAdminUserController()
				roleController := controllers.NewAdminRoleController()
				permissionController := controllers.NewAdminPermissionController()
				menuController := controllers.NewAdminMenuController()

				permitted.Get("/system", indexController.SystemIndex)

				permitted.Get("/system/admin_users", userController.Index)
				permitted.Post("/system/admin_users", userController.Store)
				permitted.Get("/system/admin_users/{id}/edit", userController.Edit)
				permitted.Put("/system/admin_users/{id}", userController.Update)
				permitted.Delete("/system/admin_users/{id}", userController.Destroy)
				permitted.Delete("/system/admin_users", userController.Destroy)

				permitted.Get("/system/admin_roles", roleController.Index)
				permitted.Get("/system/admin_roles/options", roleController.Options)
				permitted.Post("/system/admin_roles", roleController.Store)
				permitted.Post("/system/admin_roles/save_permissions", roleController.SavePermissions)
				permitted.Get("/system/admin_roles/{id}/edit", roleController.Edit)
				permitted.Put("/system/admin_roles/{id}", roleController.Update)
				permitted.Delete("/system/admin_roles/{id}", roleController.Destroy)
				permitted.Delete("/system/admin_roles", roleController.Destroy)

				permitted.Get("/system/admin_permissions", permissionController.Index)
				permitted.Get("/system/admin_permissions/tree", permissionController.Tree)
				permitted.Post("/system/admin_permissions", permissionController.Store)
				permitted.Get("/system/admin_permissions/{id}/edit", permissionController.Edit)
				permitted.Put("/system/admin_permissions/{id}", permissionController.Update)
				permitted.Delete("/system/admin_permissions/{id}", permissionController.Destroy)
				permitted.Delete("/system/admin_permissions", permissionController.Destroy)

				permitted.Get("/system/admin_menus", menuController.Index)
				permitted.Get("/system/admin_menus/tree", menuController.TreeOptions)
				permitted.Post("/system/admin_menus", menuController.Store)
				permitted.Post("/system/admin_menus/save_order", menuController.SaveOrder)
				permitted.Get("/system/admin_menus/{id}/edit", menuController.Edit)
				permitted.Put("/system/admin_menus/{id}", menuController.Update)
				permitted.Delete("/system/admin_menus/{id}", menuController.Destroy)
				permitted.Delete("/system/admin_menus", menuController.Destroy)

				// 组织架构与字典（M3）
				departmentController := controllers.NewAdminDepartmentController()
				positionController := controllers.NewAdminPositionController()
				dictionaryController := controllers.NewAdminDictionaryController()

				permitted.Get("/system/admin_departments", departmentController.Index)
				permitted.Get("/system/admin_departments/tree", departmentController.TreeOptions)
				permitted.Post("/system/admin_departments", departmentController.Store)
				permitted.Get("/system/admin_departments/{id}/edit", departmentController.Edit)
				permitted.Put("/system/admin_departments/{id}", departmentController.Update)
				permitted.Delete("/system/admin_departments/{id}", departmentController.Destroy)
				permitted.Delete("/system/admin_departments", departmentController.Destroy)

				permitted.Get("/system/admin_positions", positionController.Index)
				permitted.Get("/system/admin_positions/options", positionController.Options)
				permitted.Post("/system/admin_positions", positionController.Store)
				permitted.Get("/system/admin_positions/{id}/edit", positionController.Edit)
				permitted.Put("/system/admin_positions/{id}", positionController.Update)
				permitted.Delete("/system/admin_positions/{id}", positionController.Destroy)
				permitted.Delete("/system/admin_positions", positionController.Destroy)

				permitted.Get("/system/admin_dictionaries", dictionaryController.Index)
				permitted.Post("/system/admin_dictionaries", dictionaryController.Store)
				permitted.Get("/system/admin_dictionaries/{id}/edit", dictionaryController.Edit)
				permitted.Put("/system/admin_dictionaries/{id}", dictionaryController.Update)
				permitted.Delete("/system/admin_dictionaries/{id}", dictionaryController.Destroy)
				permitted.Delete("/system/admin_dictionaries", dictionaryController.Destroy)

				permitted.Get("/system/admin_dictionary_items", dictionaryController.ItemsIndex)
				permitted.Post("/system/admin_dictionary_items", dictionaryController.ItemsStore)
				permitted.Get("/system/admin_dictionary_items/{id}/edit", dictionaryController.ItemsEdit)
				permitted.Put("/system/admin_dictionary_items/{id}", dictionaryController.ItemsUpdate)
				permitted.Delete("/system/admin_dictionary_items/{id}", dictionaryController.ItemsDestroy)
				permitted.Delete("/system/admin_dictionary_items", dictionaryController.ItemsDestroy)

				// 操作日志（只读）
				logController := controllers.NewAdminOperationLogController()
				permitted.Get("/system/admin_operation_logs", logController.Index)
				permitted.Get("/system/admin_operation_logs/{id}/edit", logController.Edit)

				// 插件业务路由
				for _, plugin := range plugins {
					plugin.Routes(permitted)
				}
			})
		})
	})
}
