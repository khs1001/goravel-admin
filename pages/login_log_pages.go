package pages

// 登录日志页面 schema（只读列表 + 详情抽屉）。

// LoginLogsPage 登录日志列表页。
func LoginLogsPage() any {
	filter := []any{
		textControl("username", "用户名", false),
		textControl("ip", "IP", false),
		map[string]any{
			"type": "select", "name": "status", "label": "状态", "clearable": true,
			"options": []any{
				map[string]any{"label": "登录成功", "value": "success"},
				map[string]any{"label": "密码错误", "value": "failed"},
				map[string]any{"label": "限流锁定", "value": "throttled"},
				map[string]any{"label": "账号禁用", "value": "disabled"},
			},
		},
	}

	columns := []any{
		idColumn(),
		textColumn("username", "用户名"),
		mappingTagColumn("status", "状态", map[string]any{
			"success":   map[string]any{"label": "登录成功", "color": "success"},
			"failed":    map[string]any{"label": "密码错误", "color": "danger"},
			"throttled": map[string]any{"label": "限流锁定", "color": "warning"},
			"disabled":  map[string]any{"label": "账号禁用", "color": "danger"},
		}),
		textColumn("message", "说明"),
		textColumn("ip", "IP"),
		textColumn("user_agent", "UA"),
		datetimeColumn("created_at", "时间"),
	}

	crud := crudBase("/system/admin_login_logs?_action=getData", columns, filter, nil)
	crud["rowActions"] = []any{
		map[string]any{
			"type":       "button",
			"label":      "详情",
			"level":      "link",
			"actionType": "drawer",
			"drawer": drawerForm("登录日志详情", "", "/system/admin_login_logs/${id}/edit", []any{
				map[string]any{"type": "static", "name": "username", "label": "用户名"},
				map[string]any{"type": "static", "name": "status", "label": "状态"},
				map[string]any{"type": "static", "name": "message", "label": "说明"},
				map[string]any{"type": "static", "name": "ip", "label": "IP"},
				map[string]any{"type": "static", "name": "user_agent", "label": "User-Agent"},
			}),
		},
	}

	return listPage("登录日志", crud)
}
