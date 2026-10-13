package config

import (
	"github.com/goravel/framework/facades"
)

// init registers the package configuration under the "admin" key.
//
// The file is published to the host `config/admin.go` by `setup/setup.go`, so
// the host can override any value. It also mirrors Owl Admin's config shape so
// the shipped `admin-views` frontend can be driven without changes.
func init() {
	config := facades.Config()
	config.Add("admin", map[string]any{
		// Application name shown on the login page and the browser title.
		"name": config.Env("ADMIN_APP_NAME", "Goravel Admin"),
		// Logo and default avatar served from the published frontend assets.
		"logo":           config.Env("ADMIN_LOGO", "/admin-assets/logo.png"),
		"default_avatar": config.Env("ADMIN_DEFAULT_AVATAR", "/admin-assets/default-avatar.png"),

		// The SPA entry path and the base path of the frontend assets.
		"path":   config.Env("ADMIN_PATH", "admin"),
		"assets": config.Env("ADMIN_ASSETS", "admin-assets"),

		// API route prefix. The built frontend hardcodes `window.$adminApiPrefix`
		// to `/admin-api`; change both if you change this.
		"route": map[string]any{
			"prefix": config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
		},

		"auth": map[string]any{
			// Whether authentication is enforced at all.
			"enable": config.Env("ADMIN_ENABLE_AUTH", true),
			// Whether the per-request permission check is enforced.
			"permission": config.Env("ADMIN_ENABLE_PERMISSION", true),
			// Whether the login captcha is required.
			"login_captcha": config.Env("ADMIN_LOGIN_CAPTCHA", true),
			// JWT time to live (minutes). 0 falls back to the host `jwt.ttl`.
			"token_expiration": config.Env("ADMIN_TOKEN_EXPIRATION", 0),
		},

		"layout": map[string]any{
			"title":  config.Env("ADMIN_SITE_TITLE", "%title% | Goravel Admin"),
			"footer": "<a href=\"https://github.com/khs1001/goravel-admin\" target=\"_blank\">Goravel Admin</a>",
			"header": map[string]any{
				"refresh":       true,
				"dark":          true,
				"full_screen":   true,
				"locale_toggle": true,
				"theme_config":  true,
			},
			"locale_options": []any{
				map[string]any{"label": "简体中文", "value": "zh_CN"},
				map[string]any{"label": "English", "value": "en"},
			},
			"keep_alive_exclude": []any{},
		},
	})
}
