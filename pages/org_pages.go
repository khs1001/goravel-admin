package pages

// 组织架构与字典页面 schema（M3）：部门（树形）/ 岗位 / 字典（含字典项抽屉 CRUD）。

// ---- 部门管理 ----

// DepartmentsPage 部门列表页（树形数据，同菜单页 loadDataOnce 模式）。
func DepartmentsPage() any {
	filter := []any{
		textControl("name", "部门名称", false),
	}

	columns := []any{
		map[string]any{"type": "text", "name": "name", "label": "部门名称"},
		map[string]any{"type": "text", "name": "custom_order", "label": "排序", "sortable": true},
		map[string]any{"type": "text", "name": "remark", "label": "备注"},
		map[string]any{"type": "datetime", "name": "created_at", "label": "创建时间", "sortable": true},
	}

	crud := crudBase("/system/admin_departments?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_departments", "", departmentForm())),
	})
	crud["perPage"] = 999
	crud["loadDataOnce"] = true
	crud["rowActions"] = rowActions(
		drawerForm("编辑", "put:/system/admin_departments/${id}", "/system/admin_departments/${id}/edit", departmentForm()),
		"delete:/system/admin_departments/${id}", "",
	)

	return listPage("部门", crud)
}

func departmentForm() []any {
	return []any{
		map[string]any{
			"type": "tree-select", "name": "parent_id", "label": "上级部门",
			"source": "/system/admin_departments/tree",
			"labelField": "name", "valueField": "id",
			"clearable": true,
		},
		textControl("name", "部门名称", true),
		map[string]any{"type": "input-number", "name": "custom_order", "label": "排序", "value": 0},
		map[string]any{"type": "textarea", "name": "remark", "label": "备注"},
	}
}

// ---- 岗位管理 ----

// PositionsPage 岗位列表页。
func PositionsPage() any {
	filter := []any{
		textControl("name", "岗位名称", false),
	}

	columns := []any{
		map[string]any{"type": "text", "name": "id", "label": "ID", "sortable": true},
		map[string]any{"type": "text", "name": "name", "label": "岗位名称"},
		map[string]any{"type": "text", "name": "custom_order", "label": "排序", "sortable": true},
		map[string]any{"type": "text", "name": "remark", "label": "备注"},
		map[string]any{"type": "datetime", "name": "created_at", "label": "创建时间", "sortable": true},
	}

	crud := crudBase("/system/admin_positions?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_positions", "", positionForm())),
	})
	crud["bulkActions"] = bulkDelete("delete:/system/admin_positions?ids=${ids|raw}")
	crud["rowActions"] = rowActions(
		drawerForm("编辑", "put:/system/admin_positions/${id}", "/system/admin_positions/${id}/edit", positionForm()),
		"delete:/system/admin_positions/${id}", "",
	)

	return listPage("岗位", crud)
}

func positionForm() []any {
	return []any{
		textControl("name", "岗位名称", true),
		map[string]any{"type": "input-number", "name": "custom_order", "label": "排序", "value": 0},
		map[string]any{"type": "textarea", "name": "remark", "label": "备注"},
	}
}

// ---- 字典管理 ----

// DictionariesPage 字典列表页（行操作含字典项抽屉 CRUD）。
func DictionariesPage() any {
	filter := []any{
		textControl("name", "字典名称", false),
		textControl("key", "key", false),
	}

	columns := []any{
		map[string]any{"type": "text", "name": "id", "label": "ID", "sortable": true},
		map[string]any{"type": "text", "name": "name", "label": "字典名称"},
		map[string]any{"type": "tag", "name": "key", "label": "key"},
		map[string]any{"type": "text", "name": "item_count", "label": "字典项数"},
		map[string]any{"type": "text", "name": "remark", "label": "备注"},
		map[string]any{"type": "datetime", "name": "created_at", "label": "创建时间", "sortable": true},
	}

	actions := []any{
		map[string]any{
			"type":       "button",
			"label":      "字典项",
			"level":      "link",
			"actionType": "drawer",
			"drawer":     dictionaryItemsDrawer(),
		},
	}
	actions = append(actions, rowActions(
		drawerForm("编辑", "put:/system/admin_dictionaries/${id}", "/system/admin_dictionaries/${id}/edit", dictionaryForm()),
		"delete:/system/admin_dictionaries/${id}", "",
	)...)

	columns = append(columns, map[string]any{"type": "operation", "label": "操作", "buttons": actions})

	crud := crudBase("/system/admin_dictionaries?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_dictionaries", "", dictionaryForm())),
	})
	crud["bulkActions"] = bulkDelete("delete:/system/admin_dictionaries?ids=${ids|raw}")

	return listPage("字典", crud)
}

func dictionaryForm() []any {
	return []any{
		textControl("name", "字典名称", true),
		map[string]any{
			"type": "input-text", "name": "key", "label": "key", "required": true,
			"description": "唯一英文标识，业务侧按 key 读取选项，如 gender",
		},
		map[string]any{"type": "textarea", "name": "remark", "label": "备注"},
	}
}

// dictionaryItemsDrawer 字典项管理抽屉（嵌套 CRUD）。
func dictionaryItemsDrawer() map[string]any {
	itemsCrud := map[string]any{
		"type":                 "crud",
		"api":                  "/system/admin_dictionary_items?dictionary_id=${id}&_action=getData",
		"perPage":              999,
		"syncLocation":         false,
		"primaryField":         "id",
		"columns": []any{
			map[string]any{"type": "text", "name": "label", "label": "标签"},
			map[string]any{"type": "text", "name": "value", "label": "值"},
			map[string]any{"type": "text", "name": "custom_order", "label": "排序", "sortable": true},
			map[string]any{
				"type": "operation", "label": "操作",
				"buttons": rowActions(
					drawerForm("编辑", "put:/system/admin_dictionary_items/${id}", "/system/admin_dictionary_items/${id}/edit", dictionaryItemForm()),
					"delete:/system/admin_dictionary_items/${id}", "",
				),
			},
		},
		"headerToolbar": []any{
			map[string]any{
				"type":       "button",
				"label":      "新增",
				"icon":       "fa fa-add",
				"level":      "primary",
				"actionType": "drawer",
				"drawer":     drawerForm("新增字典项", "post:/system/admin_dictionary_items?dictionary_id=${id}", "", dictionaryItemForm()),
			},
			map[string]any{"type": "reload", "align": "right"},
		},
		"footerToolbar": []any{"statistics"},
	}

	return map[string]any{
		"title":          "字典项 - ${name}",
		"size":           "lg",
		"closeOnEsc":     true,
		"closeOnOutside": true,
		"body":           itemsCrud,
	}
}

func dictionaryItemForm() []any {
	return []any{
		textControl("label", "标签", true),
		textControl("value", "值", true),
		map[string]any{"type": "input-number", "name": "custom_order", "label": "排序", "value": 0},
	}
}
