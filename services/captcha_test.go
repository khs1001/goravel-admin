package services

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
)

// TestNormalizeCaptchaInput 输入归一化：去首尾空白、转小写。
func TestNormalizeCaptchaInput(t *testing.T) {
	cases := map[string]string{
		" 1234 ": "1234",
		"AbCd":   "abcd",
		"\t1a2b": "1a2b",
		"":       "",
		"   ":    "",
	}
	for input, want := range cases {
		if got := NormalizeCaptchaInput(input); got != want {
			t.Errorf("NormalizeCaptchaInput(%q) = %q, want %q", input, got, want)
		}
	}
}

// TestCaptchaCacheKey 缓存 key 由归一化答案派生，同答案必同 key，不同答案不同 key。
func TestCaptchaCacheKey(t *testing.T) {
	want := func(answer string) string {
		sum := md5.Sum([]byte(answer))
		return "admin:captcha:" + hex.EncodeToString(sum[:])
	}

	if got := CaptchaCacheKey("1234"); got != want("1234") {
		t.Errorf("CaptchaCacheKey(1234) = %s", got)
	}
	if CaptchaCacheKey(NormalizeCaptchaInput(" ABCD ")) != CaptchaCacheKey("abcd") {
		t.Errorf("normalized inputs should derive identical keys")
	}
	if CaptchaCacheKey("1234") == CaptchaCacheKey("4321") {
		t.Errorf("different answers must derive different keys")
	}
}
