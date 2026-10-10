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
		idColumn(),
		textColumn("username", "操作人"),
		mappingTagColumn("action", "类型", map[string]any{
			"request": map[string]any{"label": "请求", "color": "active"},
			"created": map[string]any{"label": "创建", "color": "success"},
			"updated": map[string]any{"label": "更新", "color": "warning"},
			"deleted": map[string]any{"label": "删除", "color": "danger"},
		}),
		textColumn("resource", "资源"),
		textColumn("method", "方法"),
		textColumn("path", "路径"),
		textColumn("ip", "IP"),
		datetimeColumn("created_at", "时间"),
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
				map[string]any{"type": "static", "name": "resource_id", "label": "资源 ID"},
				map[string]any{"type": "static", "name": "method", "label": "方法"},
				map[string]any{"type": "static", "name": "path", "label": "路径"},
				map[string]any{"type": "static", "name": "ip", "label": "IP"},
				map[string]any{"type": "static", "name": "input", "label": "参数"},
				// 字段级变更对照（仅 updated 动作有数据）
				map[string]any{
					"type": "static-table", "name": "changes", "label": "变更对照",
					"columns": []any{
						map[string]any{"name": "field", "label": "字段"},
						map[string]any{"name": "old", "label": "变更前"},
						map[string]any{"name": "new", "label": "变更后"},
					},
				},
			}),
		},
	}

	return listPage("操作日志", crud)
}
