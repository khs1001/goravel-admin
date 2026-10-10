package support

import (
	"encoding/json"
	"testing"
)

// decode 反序列化 envelope 为 map 便于断言 JSON 形状。
func decode(t *testing.T, e Envelope) map[string]any {
	t.Helper()

	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return out
}

// 契约：成功响应 status=0，前端取 data。
func TestSuccess(t *testing.T) {
	env := Success(map[string]any{"k": "v"}, "ok")
	if env.Status != 0 || env.Msg != "ok" {
		t.Fatalf("unexpected envelope: %+v", env)
	}

	out := decode(t, env)
	if out["status"].(float64) != 0 {
		t.Fatalf("status should be 0, got %v", out["status"])
	}
	if out["data"].(map[string]any)["k"] != "v" {
		t.Fatalf("data lost: %v", out["data"])
	}
	// code 在成功时不输出
	if _, exists := out["code"]; exists {
		t.Fatalf("code should be omitted on success")
	}
}

// 契约：data 为 nil 时输出 {}（对齐 PHP success(null) 的对象转换，前端 amis 不报错）。
func TestSuccessNilData(t *testing.T) {
	out := decode(t, Success(nil, ""))

	data, ok := out["data"].(map[string]any)
	if !ok || len(data) != 0 {
		t.Fatalf("nil data should serialize to {}, got %v", out["data"])
	}
}

// 契约：SuccessMessage 的 data 为空数组（对齐 PHP success([])）。
func TestSuccessMessage(t *testing.T) {
	env := SuccessMessage("saved")
	if env.Status != 0 {
		t.Fatalf("status should be 0")
	}

	raw, _ := json.Marshal(env.Data)
	if string(raw) != "[]" {
		t.Fatalf("data should be [], got %s", raw)
	}
}

// 契约：失败响应 status=1，data 为 null。
func TestFail(t *testing.T) {
	out := decode(t, Fail("没有权限"))
	if out["status"].(float64) != 1 {
		t.Fatalf("status should be 1")
	}
	if out["msg"] != "没有权限" {
		t.Fatalf("msg mismatch")
	}
	if v, exists := out["data"]; !exists || v != nil {
		t.Fatalf("data should be null, got %v", out["data"])
	}
}

// 契约：未登录响应 code=401 且 doNotDisplayToast=1，前端据此跳登录页并不弹 toast。
func TestUnauthorized(t *testing.T) {
	out := decode(t, Unauthorized("请先登录"))
	if out["code"].(float64) != 401 {
		t.Fatalf("code should be 401")
	}
	if out["status"].(float64) != 1 {
		t.Fatalf("status should be 1")
	}
	if out["doNotDisplayToast"].(float64) != 1 {
		t.Fatalf("doNotDisplayToast should be 1")
	}
}
