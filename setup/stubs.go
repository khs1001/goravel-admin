package main

import "strings"

func config(configPackage string, facadesImport, facadesPackage string) string {
	content := `package DummyPackage

import (
	"DummyFacadesImport"
)

func init() {
	config := DummyFacadesPackage.Config()
	config.Add("admin", map[string]any{
		"name":           config.Env("ADMIN_APP_NAME", "Goravel Admin"),
		"logo":           config.Env("ADMIN_LOGO", "/admin-assets/logo.png"),
		"default_avatar": config.Env("ADMIN_DEFAULT_AVATAR", "/admin-assets/default-avatar.png"),

		"route": map[string]any{
			"prefix":               config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
			"without_extra_routes": []string{"/dashboard"},
		},

		"auth": map[string]any{
			"login_captcha": false,
			"enable":        true,
			"permission":    false,
		},

		"layout": map[string]any{
			"title": "%title%",
			"locale_options": map[string]string{
				"en":    "English",
				"zh_CN": "简体中文",
			},
		},

		"show_development_tools": false,
	})
}
`
	content = strings.ReplaceAll(content, "DummyPackage", configPackage)
	content = strings.ReplaceAll(content, "DummyFacadesImport", facadesImport)
	content = strings.ReplaceAll(content, "DummyFacadesPackage", facadesPackage)

	return content
}
