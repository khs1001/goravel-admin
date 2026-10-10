package controllers

import (
	"encoding/json"
	"testing"
)

// intSlice：amis 表单数组的各种形态（M2/M4 曾踩 Input() 数组返回空串的坑，此处固化正确解析）。
func TestIntSlice(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  []int
	}{
		{"[]any float64", []any{float64(1), float64(2)}, []int{1, 2}},
		{"[]any string", []any{"3", "4"}, []int{3, 4}},
		{"[]any mixed", []any{float64(1), "2"}, []int{1, 2}},
		{"[]int", []int{5, 6}, []int{5, 6}},
		{"json string", "[7,8]", []int{7, 8}},
		{"single number string", "9", []int{9}},
		{"float64", float64(10), []int{10}},
		{"invalid json string", "not-json", nil},
		{"invalid items skipped", []any{float64(1), "x", nil}, []int{1}},
		{"nil", nil, nil},
		{"empty array", []any{}, []int{}},
	}

	for _, c := range cases {
		got := intSlice(c.value)
		if !intSliceEqual(got, c.want) {
			t.Errorf("%s: intSlice(%v) = %v, want %v", c.name, c.value, got, c.want)
		}
	}
}

// strSlice：http_method 多选形态。
func TestStrSlice(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  []string
	}{
		{"[]any", []any{"GET", "POST"}, []string{"GET", "POST"}},
		{"[]any empty strings skipped", []any{"GET", "", "POST"}, []string{"GET", "POST"}},
		{"[]string", []string{"PUT"}, []string{"PUT"}},
		{"single string", "DELETE", []string{"DELETE"}},
		{"empty string", "", nil},
		{"nil", nil, nil},
		{"non-string items skipped", []any{"GET", float64(1)}, []string{"GET"}},
	}

	for _, c := range cases {
		got := strSlice(c.value)
		if !strSliceEqual(got, c.want) {
			t.Errorf("%s: strSlice(%v) = %v, want %v", c.name, c.value, got, c.want)
		}
	}
}

// parseIDList：批量删除 ids 串。
func TestParseIDList(t *testing.T) {
	cases := []struct {
		raw  string
		want []int
	}{
		{"1,2,3", []int{1, 2, 3}},
		{" 1 , 2 ", []int{1, 2}},
		{"", nil},
		{"a,b", []int{}},
		{"0,-1,2", []int{2}}, // 非正数跳过
		{"5", []int{5}},
	}

	for _, c := range cases {
		got := parseIDList(c.raw)
		if !intSliceEqual(got, c.want) {
			t.Errorf("parseIDList(%q) = %v, want %v", c.raw, got, c.want)
		}
	}
}

// listEnvelope：amis CRUD 消费 {items, total} 契约。
func TestListEnvelope(t *testing.T) {
	env := listEnvelope([]any{map[string]any{"id": 1}}, 42)

	if env.Status != 0 {
		t.Fatalf("status should be 0")
	}

	raw, _ := json.Marshal(env.Data)
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatalf("data should be object: %v", err)
	}

	if data["total"].(float64) != 42 {
		t.Fatalf("total mismatch: %v", data["total"])
	}
	if items, ok := data["items"].([]any); !ok || len(items) != 1 {
		t.Fatalf("items mismatch: %v", data["items"])
	}
}

func intSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func strSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
