package middleware

import (
	"encoding/json"
	"strings"
	"testing"
)

// sanitizeInput：敏感字段脱敏（键名大小写不敏感）。
func TestSanitizeInputMasksSensitive(t *testing.T) {
	out := sanitizeInput(map[string]any{
		"username": "admin",
		"Password": "secret",
		"old_password": "a",
		"New_Password": "b",
		"confirm_password": "b",
		"token": "jwt-value",
	})

	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output should be valid json: %v", err)
	}

	if parsed["username"] != "admin" {
		t.Fatalf("non-sensitive field should pass through")
	}
	for _, key := range []string{"Password", "old_password", "New_Password", "confirm_password", "token"} {
		if parsed[key] != "***" {
			t.Errorf("field %q should be masked, got %v", key, parsed[key])
		}
	}
}

// sanitizeInput：空输入返回空串。
func TestSanitizeInputEmpty(t *testing.T) {
	if got := sanitizeInput(map[string]any{}); got != "" {
		t.Fatalf("empty map should yield empty string, got %q", got)
	}
}

// sanitizeInput：超长内容截断到 2000 字节。
func TestSanitizeInputTruncates(t *testing.T) {
	long := strings.Repeat("a", 5000)
	out := sanitizeInput(map[string]any{"data": long})

	if len(out) != 2000 {
		t.Fatalf("output should be truncated to 2000 bytes, got %d", len(out))
	}
}
