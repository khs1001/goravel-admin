package pages

// 内置 amis 页面注册表：sign → 页面 schema。
// 业务插件通过 contracts.Plugin.Pages() 追加，DB（admin_pages）可覆盖。

// Builtin 全部内置页面。
func Builtin() map[string]any {
	return map[string]any{
		"dashboard":    Dashboard(),
		"user_setting": UserSetting(),
	}
}

// Dashboard 控制台页面（GET /admin-api/dashboard 与 sign=dashboard 均可命中）。
func Dashboard() any {
	statCard := func(title, name, icon string) map[string]any {
		return map[string]any{
			"md": 3,
			"body": map[string]any{
				"type": "card",
				"header": map[string]any{
					"title":   title,
					"subTitle": icon,
				},
				"body": map[string]any{
					"type":       "tpl",
					"tpl":        "<span class='text-2xl font-bold'>${" + name + "|default:0}</span>",
					"className":  "py-2",
				},
			},
		}
	}

	return map[string]any{
		"type":    "page",
		"title":   "仪表盘",
		"initApi": "/dashboard_stats",
		"body": []any{
			map[string]any{
				"type": "grid",
				"columns": []any{
					statCard("管理员", "user_count", "fa fa-users"),
					statCard("角色", "role_count", "fa fa-id-badge"),
					statCard("菜单", "menu_count", "fa fa-bars"),
					statCard("今日操作日志", "today_log_count", "fa fa-history"),
				},
			},
			map[string]any{
				"type": "grid",
				"columns": []any{
					map[string]any{
						"md": 6,
						"body": map[string]any{
							"type": "card",
							"header": map[string]any{
								"title": "欢迎使用 Goravel Admin",
								"subTitle": "基于 goravel 扩展包机制的后台平台",
							},
							"body": map[string]any{
								"type": "tpl",
								"tpl":  "后端已由 Goravel 驱动，页面全部由 amis JSON 下发。<br/>登录、鉴权、菜单、页面 schema、RBAC、操作日志全链路已打通。",
							},
						},
					},
					map[string]any{
						"md": 6,
						"body": map[string]any{
							"type": "card",
							"header": map[string]any{"title": "最近操作日志"},
							"body": map[string]any{
								"type": "table",
								"source": "${recent_logs}",
								"columns": []any{
									map[string]any{"name": "username", "label": "操作人"},
									map[string]any{"name": "action", "label": "类型"},
									map[string]any{"name": "method", "label": "方法"},
									map[string]any{"name": "path", "label": "路径"},
									map[string]any{"name": "created_at", "label": "时间"},
								},
							},
						},
					},
				},
			},
		},
	}
}

// UserSetting 个人设置页面（GET /admin-api/user_setting）。
func UserSetting() any {
	return map[string]any{
		"type":  "page",
		"title": "个人设置",
		"body": []any{
			map[string]any{
				"type":            "form",
				"title":           "",
				"panelClassName":  "px-48 m:px-0",
				"mode":            "horizontal",
				"initApi":         "/current-user",
				"api":             "put:/user_setting",
				"body": []any{
					map[string]any{
						"type": "input-image", "label": "头像", "name": "avatar",
						"receiver": "/upload",
					},
					map[string]any{
						"type": "input-text", "label": "名称", "name": "name", "required": true,
					},
					map[string]any{
						"type": "input-password", "label": "旧密码", "name": "old_password",
					},
					map[string]any{
						"type": "input-password", "label": "新密码", "name": "password",
					},
					map[string]any{
						"type": "input-password", "label": "确认密码", "name": "confirm_password",
					},
				},
			},
		},
	}
}

// UserDropdown 登录用户顶栏下拉菜单（current-user 返回的 menus 字段）。
// 对齐 owl-admin AuthController::currentUser。
func UserDropdown(name string, avatar any) any {
	return map[string]any{
		"type":          "dropdown-button",
		"hideCaret":     true,
		"trigger":       "hover",
		"label":         name,
		"className":     "h-full w-full",
		"btnClassName":  "navbar-user w-full",
		"menuClassName": "min-w-0",
		"icon":          avatar,
		"buttons": []any{
			map[string]any{
				"type":          "button",
				"iconClassName": "pr-2",
				"icon":          "fa fa-user-gear",
				"label":         "个人设置",
				"onClick":       "window.location.hash = '#/user_setting'",
			},
			map[string]any{
				"type":          "button",
				"iconClassName": "pr-2",
				"label":         "退出登录",
				"icon":          "fa-solid fa-right-from-bracket",
				"onClick":       "window.$owl.logout()",
			},
		},
	}
}
