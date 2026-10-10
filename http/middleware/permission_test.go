package middleware

import (
	"testing"
)

func strPtr(s string) *string { return &s }

// matchPattern：精确、尾部通配、全通配、斜杠归一。
func TestMatchPattern(t *testing.T) {
	cases := []struct {
		pattern, path string
		want          bool
	}{
		{"system/admin_users", "system/admin_users", true},
		{"system/admin_users", "system/admin_roles", false},
		{"system/admin_users", "system/admin_users/1", false}, // 精确匹配不含子路径
		{"system/admin_users*", "system/admin_users/1", true},
		{"system/admin_users*", "system/admin_users", true},
		{"system/admin_users*", "system/admin_usersxxx", true}, // 尾部通配为前缀匹配
		{"*", "anything/here", true},
		{"/system/admin_users/", "system/admin_users", true}, // 首尾斜杠归一
		{"", "", true},
		{"", "system", false},
	}

	for _, c := range cases {
		if got := matchPattern(c.pattern, c.path); got != c.want {
			t.Errorf("matchPattern(%q, %q) = %v, want %v", c.pattern, c.path, got, c.want)
		}
	}
}

// methodAllowed：nil/空/非法 JSON/空数组 = 不限方法；数组内大小写不敏感匹配。
func TestMethodAllowed(t *testing.T) {
	cases := []struct {
		httpMethod *string
		method     string
		want       bool
	}{
		{nil, "GET", true},
		{strPtr(""), "POST", true},
		{strPtr("not-json"), "DELETE", true},  // 非法 JSON 视为不限（fail-open）
		{strPtr("[]"), "PUT", true},           // 空数组不限
		{strPtr(`["GET"]`), "GET", true},
		{strPtr(`["get"]`), "GET", true},      // 大小写不敏感
		{strPtr(`["GET","POST"]`), "POST", true},
		{strPtr(`["GET"]`), "POST", false},
		{strPtr(`["GET"]`), "DELETE", false},
	}

	for _, c := range cases {
		if got := methodAllowed(c.httpMethod, c.method); got != c.want {
			t.Errorf("methodAllowed(%v, %q) = %v, want %v", c.httpMethod, c.method, got, c.want)
		}
	}
}

// inExcept：免检路径（菜单/登录/当前用户等）。
func TestInExcept(t *testing.T) {
	excepted := []string{"menus", "current-user", "user_setting", "login", "logout", "no-content", "_settings", "captcha"}
	for _, p := range excepted {
		if !(&Permission{}).inExcept(p) {
			t.Errorf("path %q should be excepted", p)
		}
	}

	for _, p := range []string{"system/admin_users", "system/admin_roles/1", "dashboard"} {
		if (&Permission{}).inExcept(p) {
			t.Errorf("path %q should NOT be excepted", p)
		}
	}
}

// pathAllowed：http_path 多行 + 可选方法前缀 + /admin-api 前缀剥离（对齐 owl-admin 种子格式）。
func TestPathAllowed(t *testing.T) {
	cases := []struct {
		name              string
		httpPath          string
		method, path      string
		want              bool
	}{
		{"不带方法前缀=全部方法放行", "/admin-api/system/admin_users*", "DELETE", "system/admin_users/3", true},
		{"方法前缀匹配", "GET /admin-api/system/admin_users", "GET", "system/admin_users", true},
		{"方法前缀不匹配", "GET /admin-api/system/admin_users", "POST", "system/admin_users", false},
		{"多行中任一命中", "GET /admin-api/system/admin_positions\nPOST /admin-api/system/admin_positions", "POST", "system/admin_positions", true},
		{"多行全不命中", "GET /admin-api/system/admin_positions", "DELETE", "system/admin_positions/1", false},
		{"空行跳过", "\n  /admin-api/menus  \n", "GET", "menus", true},
		{"路径前缀已剥离时仍匹配", "system/admin_menus/save_order", "POST", "system/admin_menus/save_order", true},
		{"尾部通配放行子资源", "GET /admin-api/system/admin_menus*", "GET", "system/admin_menus/5", true},
	}

	for _, c := range cases {
		if got := pathAllowed(c.httpPath, c.method, c.path); got != c.want {
			t.Errorf("%s: pathAllowed(%q, %q, %q) = %v, want %v", c.name, c.httpPath, c.method, c.path, got, c.want)
		}
	}
}
