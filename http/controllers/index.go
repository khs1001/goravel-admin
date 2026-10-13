package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/resources/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// IndexController implements the Owl index endpoints: /menus, /_settings and
// /page_schema.
type IndexController struct {
	BaseController
}

// NewIndexController creates an IndexController.
func NewIndexController() *IndexController {
	return &IndexController{}
}

// Menus returns the permission-filtered route tree consumed by the frontend,
// merged with the (enabled) plugin contributions.
func (c *IndexController) Menus(ctx http.Context) http.Response {
	user := current.User(ctx)

	menuService := services.NewMenuService()
	menus, err := menuService.Routable(user)
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	homeID := menuService.HomeMenuID(menus)
	tree := owl.BuildMenuTree(menus, 0)

	routes := owl.ToRoutes(tree, homeID)
	routes = append(routes, c.extraRoutes()...)
	routes = owl.MergePluginRoutes(routes, c.pluginMenuNodes(user))
	routes = plugin.R().RunMenusBuilt(ctx, routes)

	return ctx.Response().Json(http.StatusOK, owl.Success(routes))
}

// pluginMenuNodes returns the plugin menu nodes the user may see.
func (c *IndexController) pluginMenuNodes(user *models.SysUser) []map[string]any {
	nodes := plugin.R().MenuNodes()

	menuService := services.NewMenuService()
	if menuService.IsSuperAdmin(user) {
		return nodes
	}

	result := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		perm := fmt.Sprintf("%v", node["perms"])
		if perm == "" || menuService.HasPermission(user, perm) {
			result = append(result, node)
		}
	}

	return result
}

// extraRoutes returns the non-database routes (hidden pages) added on top of the
// menu tree, e.g. the personal settings page.
func (c *IndexController) extraRoutes() []map[string]any {
	if !facades.Config().GetBool("admin.auth.enable", true) {
		return nil
	}

	return []map[string]any{
		{
			"name":       "user_setting",
			"path":       "/user_setting",
			"component":  "amis",
			"page_sign":  "user_setting",
			"is_home":    0,
			"url_type":   owl.UrlTypePage,
			"keep_alive": 0,
			"is_full":    0,
			"is_link":    false,
			"meta": map[string]any{
				"hide":  true,
				"title": "个人设置",
				"icon":  "material-symbols:manage-accounts",
			},
		},
	}
}

// Settings assembles the `/_settings` payload, including plugin assets, nav and
// default settings.
func (c *IndexController) Settings(ctx http.Context) http.Response {
	config := facades.Config()
	configService := services.NewConfigService()

	theme := configService.GetJSON("system_theme_setting")
	if theme == nil {
		theme = owl.DefaultThemeSetting()
	}

	payload := map[string]any{
		"nav":                    plugin.R().Nav(),
		"assets":                 plugin.R().Assets(),
		"app_name":               config.GetString("admin.name", "Goravel Admin"),
		"locale":                 configService.Get("admin_locale", "zh_CN"),
		"layout":                 config.Get("admin.layout"),
		"logo":                   config.GetString("admin.logo", ""),
		"login_captcha":          config.GetBool("admin.auth.login_captcha", true),
		"locale_options":         owl.LocaleOptions(config.Get("admin.layout.locale_options")),
		"show_development_tools": false,
		"system_theme_setting":   theme,
	}

	// Plugin defaults: host config < plugin default < DB sys_config.
	for key, value := range plugin.R().Settings() {
		if raw := configService.Get(key, ""); raw != "" {
			payload[key] = raw
		} else {
			payload[key] = value
		}
	}

	plugin.R().RunSettingsBuilt(ctx, payload)

	return ctx.Response().Json(http.StatusOK, owl.Success(payload))
}

// SaveSettings persists the settings sent by the frontend.
func (c *IndexController) SaveSettings(ctx http.Context) http.Response {
	configService := services.NewConfigService()

	for key, value := range ctx.Request().All() {
		if key == "system_theme_setting" {
			_ = configService.SetJSON(key, value)
			continue
		}

		_ = configService.Set(key, stringify(value))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "保存成功"))
}

// PageSchema returns the amis schema for the given sign: the embedded pages
// first, then the plugin pages (namespaced as "<plugin>.<short>").
func (c *IndexController) PageSchema(ctx http.Context) http.Response {
	sign := ctx.Request().Query("sign")

	if schema, err := pages.Get(sign); err == nil {
		return ctx.Response().Json(http.StatusOK, owl.Success(schema))
	}

	if schema, ok := plugin.R().Pages()[sign]; ok {
		return ctx.Response().Json(http.StatusOK, owl.Success(schema))
	}

	return ctx.Response().Json(http.StatusOK, owl.Fail("页面不存在: "+sign))
}

func stringify(value any) string {
	if s, ok := value.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", value)
}
