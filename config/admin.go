package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("admin", map[string]any{
		// 应用名称
		"name":           config.Env("ADMIN_APP_NAME", "Goravel Admin"),
		// 应用 logo
		"logo":           config.Env("ADMIN_LOGO", "/admin-assets/logo.png"),
		"default_avatar": config.Env("ADMIN_DEFAULT_AVATAR", "/admin-assets/default-avatar.png"),

		// 路由配置
		"route": map[string]any{
			"prefix": config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
			// 不追加 新增/详情/编辑 额外路由的路径
			"without_extra_routes": []string{"/dashboard"},
		},

		// 认证与权限
		"auth": map[string]any{
			// 登录验证码（图形验证码，一次性消费，5 分钟过期）
			"login_captcha": false,
			// 是否开启认证
			"enable": true,
			// 是否开启权限鉴权（RBAC 接口级拦截，超管不受限）
			"permission": false,
		},

		// 布局配置（下发给前端 /_settings.layout）
		"layout": map[string]any{
			"title": "%title%",
			"locale_options": map[string]string{
				"en":    "English",
				"zh_CN": "简体中文",
			},
		},

		// 是否展示开发工具
		"show_development_tools": false,
	})
}
