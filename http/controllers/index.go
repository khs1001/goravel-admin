package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/resources/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
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

// Menus returns the permission-filtered route tree consumed by the frontend.
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

	return ctx.Response().Json(http.StatusOK, owl.Success(routes))
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

// Settings assembles the `/_settings` payload.
func (c *IndexController) Settings(ctx http.Context) http.Response {
	config := facades.Config()
	configService := services.NewConfigService()

	theme := configService.GetJSON("system_theme_setting")
	if theme == nil {
		theme = owl.DefaultThemeSetting()
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{
		"nav":                    http.Json{},
		"assets":                 http.Json{},
		"app_name":               config.GetString("admin.name", "Goravel Admin"),
		"locale":                 configService.Get("admin_locale", "zh_CN"),
		"layout":                 config.Get("admin.layout"),
		"logo":                   config.GetString("admin.logo", ""),
		"login_captcha":          config.GetBool("admin.auth.login_captcha", true),
		"locale_options":         owl.LocaleOptions(config.Get("admin.layout.locale_options")),
		"show_development_tools": false,
		"system_theme_setting":   theme,
	}))
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

// PageSchema returns the embedded amis schema for the given page sign.
func (c *IndexController) PageSchema(ctx http.Context) http.Response {
	sign := ctx.Request().Query("sign")
	schema, err := pages.Get(sign)
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("页面不存在: "+sign))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(schema))
}

func stringify(value any) string {
	if s, ok := value.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", value)
}
