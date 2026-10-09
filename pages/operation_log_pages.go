package pages

// 操作日志页面 schema（只读列表 + 详情抽屉）。

// OperationLogsPage 操作日志列表页。
func OperationLogsPage() any {
	filter := []any{
		textControl("username", "操作人", false),
		textControl("path", "路径", false),
		map[string]any{
			"type": "select", "name": "action", "label": "类型", "clearable": true,
			"options": []any{
				map[string]any{"label": "请求", "value": "request"},
				map[string]any{"label": "创建", "value": "created"},
				map[string]any{"label": "更新", "value": "updated"},
				map[string]any{"label": "删除", "value": "deleted"},
			},
		},
	}

	columns := []any{
		map[string]any{"type": "text", "name": "id", "label": "ID", "sortable": true},
		map[string]any{"type": "text", "name": "username", "label": "操作人"},
		map[string]any{
			"type": "tag", "name": "action", "label": "类型",
			"map": map[string]any{
				"request": map[string]any{"label": "请求", "color": "active"},
				"created": map[string]any{"label": "创建", "color": "success"},
				"updated": map[string]any{"label": "更新", "color": "warning"},
				"deleted": map[string]any{"label": "删除", "color": "danger"},
			},
		},
		map[string]any{"type": "text", "name": "resource", "label": "资源"},
		map[string]any{"type": "text", "name": "method", "label": "方法"},
		map[string]any{"type": "text", "name": "path", "label": "路径"},
		map[string]any{"type": "text", "name": "ip", "label": "IP"},
		map[string]any{"type": "datetime", "name": "created_at", "label": "时间", "sortable": true},
	}

	crud := crudBase("/system/admin_operation_logs?_action=getData", columns, filter, nil)
	crud["rowActions"] = []any{
		map[string]any{
			"type":       "button",
			"label":      "详情",
			"level":      "link",
			"actionType": "drawer",
			"drawer": drawerForm("日志详情", "", "/system/admin_operation_logs/${id}/edit", []any{
				map[string]any{"type": "static", "name": "username", "label": "操作人"},
				map[string]any{"type": "static", "name": "action", "label": "类型"},
				map[string]any{"type": "static", "name": "resource", "label": "资源"},
				map[string]any{"type": "static", "name": "method", "label": "方法"},
				map[string]any{"type": "static", "name": "path", "label": "路径"},
				map[string]any{"type": "static", "name": "ip", "label": "IP"},
				map[string]any{"type": "static", "name": "input", "label": "参数"},
			}),
		},
	}

	return listPage("操作日志", crud)
}
