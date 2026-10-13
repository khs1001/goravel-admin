package owl

// DefaultThemeSetting mirrors the frontend's default `system_theme_setting`.
func DefaultThemeSetting() map[string]any {
	return map[string]any{
		"darkTheme":          false,
		"footer":             false,
		"breadcrumb":         true,
		"themeColor":         "#1677ff",
		"layoutMode":         "default",
		"siderTheme":         "light",
		"topTheme":           "light",
		"animateInType":      "alpha",
		"animateInDuration":  600,
		"animateOutType":     "alpha",
		"animateOutDuration": 600,
		"loginTemplate":      "default",
		"keepAlive":          false,
		"enableTab":          false,
		"tabIcon":            true,
		"accordionMenu":      false,
		"timeZone":           "Asia/Shanghai",
	}
}

// LocaleOptions returns the ordered `[{label, value}]` list the frontend locale
// dropdown expects.
//
// The configuration value must be a list (not a map): viper lowercases map keys,
// which would corrupt the `zh_CN` value.
func LocaleOptions(options any) []map[string]any {
	if list, ok := options.([]any); ok {
		result := make([]map[string]any, 0, len(list))
		for _, item := range list {
			if mapping, ok := item.(map[string]any); ok {
				result = append(result, mapping)
			}
		}
		if len(result) > 0 {
			return result
		}
	}

	return []map[string]any{
		{"label": "简体中文", "value": "zh_CN"},
		{"label": "English", "value": "en"},
	}
}
