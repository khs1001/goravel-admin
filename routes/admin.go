// Package routes registers the admin HTTP routes.
//
// The routes are registered at runtime from the service provider's Boot method,
// never through `modify.RegisterRoute`, so installing/uninstalling the package
// never edits host source files.
package routes

import (
	"path/filepath"
	"runtime"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/http/controllers"
	"github.com/khs1001/goravel-admin/http/middleware"
)

// Admin registers the API, the SPA entry and the static assets.
func Admin(router route.Route) {
	config := facades.Config()
	prefix := config.GetString("admin.route.prefix", "admin-api")
	assets := config.GetString("admin.assets", "admin-assets")
	uiPath := config.GetString("admin.path", "admin")

	registerAssets(router, assets, uiPath)

	authController := controllers.NewAuthController()
	indexController := controllers.NewIndexController()
	userController := controllers.NewSysUserController()
	roleController := controllers.NewSysRoleController()
	menuController := controllers.NewSysMenuController()
	deptController := controllers.NewSysDeptController()
	postController := controllers.NewSysPostController()
	dictController := controllers.NewSysDictController()
	configController := controllers.NewSysConfigController()

	router.Prefix(prefix).Group(func(router route.Router) {
		// Public endpoints.
		router.Post("/login", authController.Login)
		router.Get("/logout", authController.Logout)
		router.Get("/captcha", authController.Captcha)
		router.Get("/_settings", indexController.Settings)
		router.Post("/_settings", indexController.SaveSettings)

		// Authenticated endpoints.
		router.Middleware(&middleware.Authenticate{}).Group(func(router route.Router) {
			router.Get("/current-user", authController.CurrentUser)
			router.Get("/menus", indexController.Menus)
			router.Get("/page_schema", indexController.PageSchema)

			router.Prefix("/system").Group(func(router route.Router) {
				// User (sys_user).
				router.Get("/user/list", userController.List)
				router.Get("/user/roleOptions", userController.RoleOptions)
				router.Get("/user/{id}", userController.Get)
				router.Post("/user", userController.Add)
				router.Put("/user", userController.Edit)
				router.Delete("/user/{id}", userController.Remove)

				// Role (sys_role).
				router.Get("/role/list", roleController.List)
				router.Get("/role/menuTree/{id}", roleController.MenuTree)
				router.Get("/role/deptTree/{id}", roleController.DeptTree)
				router.Get("/role/{id}", roleController.Get)
				router.Post("/role", roleController.Add)
				router.Put("/role", roleController.Edit)
				router.Delete("/role/{id}", roleController.Remove)

				// Menu (sys_menu).
				router.Get("/menu/list", menuController.List)
				router.Get("/menu/tree", menuController.Tree)
				router.Get("/menu/{id}", menuController.Get)
				router.Post("/menu", menuController.Add)
				router.Put("/menu", menuController.Edit)
				router.Delete("/menu/{id}", menuController.Remove)

				// Dept (sys_dept).
				router.Get("/dept/list", deptController.List)
				router.Get("/dept/tree", deptController.Tree)
				router.Get("/dept/{id}", deptController.Get)
				router.Post("/dept", deptController.Add)
				router.Put("/dept", deptController.Edit)
				router.Delete("/dept/{id}", deptController.Remove)

				// Post (sys_post).
				router.Get("/post/list", postController.List)
				router.Get("/post/options", postController.Options)
				router.Get("/post/{id}", postController.Get)
				router.Post("/post", postController.Add)
				router.Put("/post", postController.Edit)
				router.Delete("/post/{id}", postController.Remove)

				// Dict type (sys_dict_type).
				router.Get("/dict/type/list", dictController.TypeList)
				router.Get("/dict/type/{id}", dictController.TypeGet)
				router.Post("/dict/type", dictController.TypeAdd)
				router.Put("/dict/type", dictController.TypeEdit)
				router.Delete("/dict/type/{id}", dictController.TypeRemove)

				// Dict data (sys_dict_data).
				router.Get("/dict/data/list", dictController.DataList)
				router.Get("/dict/data/type/{dictType}", dictController.DataByType)
				router.Get("/dict/data/{id}", dictController.DataGet)
				router.Post("/dict/data", dictController.DataAdd)
				router.Put("/dict/data", dictController.DataEdit)
				router.Delete("/dict/data/{id}", dictController.DataRemove)

				// Config / parameters (sys_config).
				router.Get("/config/list", configController.List)
				router.Get("/config/{id}", configController.Get)
				router.Post("/config", configController.Add)
				router.Put("/config", configController.Edit)
				router.Delete("/config/{id}", configController.Remove)
			})
		})
	})
}

// registerAssets serves the bundled SPA and its static assets from the package.
//
// The frontend uses hash-based routing, so only the exact entry path needs to
// return index.html; `/admin-assets/*` serves the rest.
func registerAssets(router route.Route, assets, uiPath string) {
	dist := distPath()
	if dist == "" {
		return
	}

	router.Static(assets, dist)

	indexFile := filepath.Join(dist, "index.html")
	router.Get("/"+uiPath, func(ctx http.Context) http.Response {
		return ctx.Response().File(indexFile)
	})
}

// distPath resolves the packaged `admin-views/dist` directory at runtime, so it
// works even when the package runs from the module cache.
func distPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Join(filepath.Dir(file), "..", "admin-views", "dist")
}
