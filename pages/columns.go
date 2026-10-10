package pages

// 通用 amis 列构建器：收敛列表页重复的列定义模式，
// 各页面 schema 以函数调用替代手写 map，保持前端契约集中可维护。

// idColumn 主键列（可排序）。
func idColumn() map[string]any {
	return map[string]any{"type": "text", "name": "id", "label": "ID", "sortable": true}
}

// textColumn 普通文本列。
func textColumn(name, label string) map[string]any {
	return map[string]any{"type": "text", "name": name, "label": label}
}

// sortableTextColumn 可排序文本列。
func sortableTextColumn(name, label string) map[string]any {
	return map[string]any{"type": "text", "name": name, "label": label, "sortable": true}
}

// tagColumn 标签列。
func tagColumn(name, label string) map[string]any {
	return map[string]any{"type": "tag", "name": name, "label": label}
}

// mappingTagColumn 带取值映射的标签列（map: value → {label, color}）。
func mappingTagColumn(name, label string, mapping map[string]any) map[string]any {
	return map[string]any{"type": "tag", "name": name, "label": label, "map": mapping}
}

// datetimeColumn 时间列（可排序）。
func datetimeColumn(name, label string) map[string]any {
	return map[string]any{"type": "datetime", "name": name, "label": label, "sortable": true}
}

// switchColumn 开关列（disabledOn 为空则不限制）。
func switchColumn(name, label, disabledOn string) map[string]any {
	col := map[string]any{"type": "switch", "name": name, "label": label}
	if disabledOn != "" {
		col["disabledOn"] = disabledOn
	}

	return col
}

// eachTagsColumn 数组字段的多标签列（items 内取 ${name} 展示）。
func eachTagsColumn(name, label string) map[string]any {
	return map[string]any{
		"type": "each", "name": name, "label": label,
		"items": map[string]any{"type": "tag", "label": "${name}", "className": "my-1"},
	}
}

// operationColumn 操作列。
func operationColumn(buttons []any) map[string]any {
	return map[string]any{"type": "operation", "label": "操作", "buttons": buttons}
}
