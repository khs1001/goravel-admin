package pages

// 系统管理页面 schema（用户/角色/权限/菜单），对齐 owl-admin 各 Admin*Controller 的 list/form。
// amis api 均为相对路径（前端 axios baseURL 已含 /admin-api），写操作带 method: 前缀。
// 树形选项/下拉选项通过 source 端点动态拉取（/system/admin_menus/tree 等）。

// ---- 通用构件 ----

// listPage CRUD 列表页骨架（对齐 ElementTrait::baseList/baseCRUD 的简化版）。
func listPage(title string, crud map[string]any) map[string]any {
	return map[string]any{
		"type":  "page",
		"title": title,
		"body":  crud,
	}
}

// crudBase CRUD 表格骨架：api、分页、工具栏。
func crudBase(api string, columns []any, filter []any, headerActions []any) map[string]any {
	crud := map[string]any{
		"type":                 "crud",
		"api":                  api,
		"perPage":              20,
		"alwaysShowPagination": true,
		"affixHeader":          false,
		"syncLocation":         false,
		"primaryField":         "id",
		"columns":              columns,
		"headerToolbar":        append(headerActions,
			map[string]any{"type": "reload", "align": "right"},
			map[string]any{"type": "filter-toggler", "align": "right"},
		),
		"footerToolbar": []any{"statistics", "pagination"},
	}
	if filter != nil {
		crud["filter"] = map[string]any{
			"type":  "form",
			"title": "",
			"body":  filter,
			"actions": []any{
				map[string]any{"type": "button", "label": "重置", "actionType": "clear-and-submit"},
				map[string]any{"type": "submit", "label": "搜索", "level": "primary"},
			},
		}
		crud["filterTogglable"] = true
		crud["filterDefaultVisible"] = false
	}

	return crud
}

// bulkDelete 批量删除按钮。
func bulkDelete(api string) []any {
	return []any{
		map[string]any{
			"type":        "button",
			"label":       "删除",
			"level":       "danger",
			"actionType":  "ajax",
			"api":         api,
			"confirmText": "确认删除选中的记录？",
		},
	}
}

// rowActions 操作列：编辑抽屉 + 删除确认（hiddenOnDelete 控制隐藏）。
func rowActions(editDrawer map[string]any, deleteAPI, hiddenOnDelete string) []any {
	return []any{
		map[string]any{
			"type":       "button",
			"label":      "编辑",
			"level":      "link",
			"actionType": "drawer",
			"drawer":     editDrawer,
		},
		map[string]any{
			"type":       "button",
			"label":      "删除",
			"level":      "link",
			"hiddenOn":   hiddenOnDelete,
			"actionType": "drawer",
			"drawer": map[string]any{
				"title":     "",
				"className": "py-2",
				"actions": []any{
					map[string]any{"type": "button", "actionType": "cancel", "label": "取消"},
					map[string]any{"type": "button", "actionType": "submit", "label": "删除", "level": "danger"},
				},
				"body": map[string]any{
					"type":          "form",
					"wrapWithPanel": false,
					"api":           deleteAPI,
					"body": []any{
						map[string]any{"type": "tpl", "className": "py-2", "tpl": "确认删除该记录？"},
					},
				},
			},
		},
	}
}

// drawerForm 抽屉表单（新增/编辑共用；initApi 为空表示新增）。
func drawerForm(title, api, initApi string, body []any) map[string]any {
	form := map[string]any{
		"type":  "form",
		"title": "",
		"mode":  "normal",
		"api":   api,
		"body":  body,
	}
	if initApi != "" {
		form["initApi"] = initApi
	}

	return map[string]any{
		"title":          title,
		"closeOnEsc":     true,
		"closeOnOutside": true,
		"size":           "md",
		"body":           form,
	}
}

// createDrawerAction 头部"新增"抽屉按钮。
func createDrawerAction(title string, drawer map[string]any) map[string]any {
	return map[string]any{
		"type":       "button",
		"label":      title,
		"icon":       "fa fa-add",
		"level":      "primary",
		"actionType": "drawer",
		"drawer":     drawer,
	}
}

// textControl 文本输入。
func textControl(name, label string, required bool) map[string]any {
	c := map[string]any{"type": "input-text", "name": name, "label": label}
	if required {
		c["required"] = true
	}

	return c
}

// switchControl 开关。
func switchControl(name, label, disabledOn string) map[string]any {
	c := map[string]any{
		"type":    "switch",
		"name":    name,
		"label":   label,
		"onText":  "启用",
		"offText": "禁用",
		"value":   1,
	}
	if disabledOn != "" {
		c["disabledOn"] = disabledOn
	}

	return c
}

// ---- 系统管理入口页（GET /system） ----

// SystemIndex 系统管理入口：导航卡片。
func SystemIndex() any {
	card := func(title, desc, link string) map[string]any {
		return map[string]any{
			"type": "card",
			"header": map[string]any{"title": title, "subTitle": desc},
			"body": map[string]any{
				"type":       "button",
				"label":      "前往管理",
				"level":      "primary",
				"actionType": "link",
				"link":       link,
			},
		}
	}

	return map[string]any{
		"type":  "page",
		"title": "系统管理",
		"body": map[string]any{
			"type": "grid",
			"columns": []any{
				map[string]any{"md": 3, "body": card("管理员", "后台账号与角色分配", "/system/admin_users")},
				map[string]any{"md": 3, "body": card("角色", "角色与权限设置", "/system/admin_roles")},
				map[string]any{"md": 3, "body": card("权限", "接口权限定义", "/system/admin_permissions")},
				map[string]any{"md": 3, "body": card("菜单", "后台菜单配置", "/system/admin_menus")},
				map[string]any{"md": 3, "body": card("部门", "树形组织架构", "/system/admin_departments")},
				map[string]any{"md": 3, "body": card("岗位", "岗位定义", "/system/admin_positions")},
				map[string]any{"md": 3, "body": card("字典", "业务字典与选项", "/system/admin_dictionaries")},
				map[string]any{"md": 3, "body": card("操作日志", "后台操作审计", "/system/admin_operation_logs")},
			},
		},
	}
}

// ---- 管理员管理 ----

// UsersPage 管理员列表页。
func UsersPage() any {
	filter := []any{
		textControl("username", "用户名", false),
		textControl("name", "昵称", false),
		map[string]any{
			"type": "select", "name": "roles", "label": "角色",
			"clearable": true, "searchable": true, "multiple": true,
			"labelField": "name", "valueField": "id",
			"source": "/system/admin_roles/options",
		},
		map[string]any{
			"type": "select", "name": "enabled", "label": "状态", "clearable": true,
			"options": []any{
				map[string]any{"label": "启用", "value": 1},
				map[string]any{"label": "禁用", "value": 0},
			},
		},
	}

	columns := []any{
		idColumn(),
		textColumn("username", "用户名"),
		textColumn("name", "昵称"),
		eachTagsColumn("roles", "角色"),
		switchColumn("enabled", "状态", "${id == 1}"),
		datetimeColumn("created_at", "创建时间"),
	}

	crud := crudBase("/system/admin_users?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_users", "", userForm())),
	})
	crud["itemCheckableOn"] = "${id != 1}"
	crud["bulkActions"] = bulkDelete("delete:/system/admin_users?ids=${ids|raw}")
	crud["rowActions"] = rowActions(
		drawerForm("编辑", "put:/system/admin_users/${id}", "/system/admin_users/${id}/edit", userForm()),
		"delete:/system/admin_users/${id}", "${id == 1}",
	)

	return listPage("管理员", crud)
}

func userForm() []any {
	return []any{
		textControl("username", "用户名", true),
		textControl("name", "昵称", true),
		map[string]any{"type": "input-password", "name": "password", "label": "密码"},
		map[string]any{"type": "input-password", "name": "confirm_password", "label": "确认密码"},
		map[string]any{
			"type": "select", "name": "roles", "label": "角色",
			"searchable": true, "multiple": true,
			"labelField": "name", "valueField": "id",
			"joinValues": false, "extractValue": true,
			"hiddenOn": "${id == 1}",
			"source":   "/system/admin_roles/options",
		},
		switchControl("enabled", "状态", "${id == 1}"),
	}
}

// ---- 角色管理 ----

// RolesPage 角色列表页（含设置权限抽屉）。
func RolesPage() any {
	filter := []any{
		textControl("name", "角色名", false),
		textControl("slug", "标识", false),
	}

	columns := []any{
		idColumn(),
		textColumn("name", "角色名"),
		tagColumn("slug", "标识"),
		datetimeColumn("created_at", "创建时间"),
		datetimeColumn("updated_at", "更新时间"),
	}

	actions := []any{
		map[string]any{
			"type":       "button",
			"label":      "设置权限",
			"level":      "link",
			"actionType": "drawer",
			"drawer": drawerForm("设置权限", "post:/system/admin_roles/save_permissions", "/system/admin_roles/${id}/edit", []any{
				map[string]any{
					"type": "tree-select", "name": "permissions", "label": "权限",
					"source":       "/system/admin_permissions/tree",
					"labelField":   "name", "valueField": "id",
					"multiple":     true, "joinValues": false, "extractValue": true,
					"cascade":      true,
					"heightAuto":   true,
				},
			}),
		},
	}
	actions = append(actions, rowActions(
		drawerForm("编辑", "put:/system/admin_roles/${id}", "/system/admin_roles/${id}/edit", roleForm()),
		"delete:/system/admin_roles/${id}", "${slug == 'administrator'}",
	)...)

	columns = append(columns, map[string]any{"type": "operation", "label": "操作", "buttons": actions})

	crud := crudBase("/system/admin_roles?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_roles", "", roleForm())),
	})
	crud["bulkActions"] = bulkDelete("delete:/system/admin_roles?ids=${ids|raw}")

	return listPage("角色", crud)
}

func roleForm() []any {
	return []any{
		textControl("name", "角色名", true),
		map[string]any{
			"type": "input-text", "name": "slug", "label": "标识", "required": true,
			"description": "唯一英文标识，如 editor",
		},
	}
}

// ---- 权限管理 ----

// PermissionsPage 权限列表页。
func PermissionsPage() any {
	filter := []any{
		textControl("name", "权限名", false),
		textControl("slug", "标识", false),
	}

	columns := []any{
		idColumn(),
		textColumn("name", "权限名"),
		tagColumn("slug", "标识"),
		textColumn("http_path", "HTTP 路径"),
		datetimeColumn("created_at", "创建时间"),
	}

	crud := crudBase("/system/admin_permissions?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_permissions", "", permissionForm())),
	})
	crud["bulkActions"] = bulkDelete("delete:/system/admin_permissions?ids=${ids|raw}")
	crud["rowActions"] = rowActions(
		drawerForm("编辑", "put:/system/admin_permissions/${id}", "/system/admin_permissions/${id}/edit", permissionForm()),
		"delete:/system/admin_permissions/${id}", "",
	)

	return listPage("权限", crud)
}

func permissionForm() []any {
	return []any{
		textControl("name", "权限名", true),
		textControl("slug", "标识", true),
		map[string]any{
			"type": "tree-select", "name": "parent_id", "label": "父级",
			"source": "/system/admin_permissions/tree",
			"labelField": "name", "valueField": "id",
			"clearable": true,
		},
		map[string]any{
			"type": "checkboxes", "name": "http_method", "label": "HTTP 方法",
			"joinValues": false, "extractValue": true,
			"options": []any{"GET", "POST", "PUT", "DELETE"},
		},
		map[string]any{
			"type": "textarea", "name": "http_path", "label": "HTTP 路径",
			"description": "每行一个，支持通配符 *，如 /admin-api/system/admin_users*",
		},
		map[string]any{"type": "input-number", "name": "custom_order", "label": "排序", "value": 0},
		map[string]any{
			"type": "tree-select", "name": "menus", "label": "挂载菜单",
			"source": "/system/admin_menus/tree",
			"labelField": "title", "valueField": "id",
			"multiple": true, "joinValues": false, "extractValue": true,
		},
	}
}

// ---- 菜单管理 ----

// MenusPage 菜单列表页（树形数据）。
func MenusPage() any {
	filter := []any{
		textControl("title", "标题", false),
		textControl("url", "路由", false),
	}

	columns := []any{
		idColumn(),
		textColumn("title", "标题"),
		textColumn("icon", "图标"),
		textColumn("url", "路由"),
		sortableTextColumn("custom_order", "排序"),
		datetimeColumn("created_at", "创建时间"),
	}

	crud := crudBase("/system/admin_menus?_action=getData", columns, filter, []any{
		createDrawerAction("新增", drawerForm("新增", "post:/system/admin_menus", "", menuForm())),
	})
	crud["perPage"] = 999
	crud["loadDataOnce"] = true
	// 拖拽排序：拖完按新顺序提交 save_order（ids= 顶层顺序）
	crud["draggable"] = true
	crud["saveOrderApi"] = "post:/system/admin_menus/save_order?ids=${ids|raw}"
	crud["rowActions"] = rowActions(
		drawerForm("编辑", "put:/system/admin_menus/${id}", "/system/admin_menus/${id}/edit", menuForm()),
		"delete:/system/admin_menus/${id}", "",
	)

	return listPage("菜单", crud)
}

func menuForm() []any {
	return []any{
		map[string]any{
			"type": "tree-select", "name": "parent_id", "label": "父级",
			"source": "/system/admin_menus/tree",
			"labelField": "title", "valueField": "id",
			"clearable": true,
		},
		textControl("title", "标题", true),
		textControl("icon", "图标", false),
		map[string]any{
			"type": "select", "name": "url_type", "label": "类型", "value": 1,
			"options": []any{
				map[string]any{"label": "路由", "value": 1},
				map[string]any{"label": "外链", "value": 2},
				map[string]any{"label": "Iframe", "value": 3},
				map[string]any{"label": "页面", "value": 4},
			},
		},
		textControl("url", "路由", false),
		map[string]any{
			"type": "input-text", "name": "component", "label": "页面标识",
			"description": "url_type=页面 时的 page_schema sign；路由类型留空即用 amis 渲染",
		},
		map[string]any{"type": "switch", "name": "visible", "label": "显示", "value": 1,
			"onText": "显示", "offText": "隐藏"},
		map[string]any{"type": "switch", "name": "is_home", "label": "首页", "value": 0},
		map[string]any{"type": "input-number", "name": "custom_order", "label": "排序", "value": 0},
	}
}
