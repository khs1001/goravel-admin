package services

import (
	"strings"
	"testing"
)

func TestIsLoginLocked(t *testing.T) {
	tests := []struct {
		name        string
		failures    int
		maxAttempts int
		expected    bool
	}{
		{"zero failures", 0, 5, false},
		{"below threshold", 4, 5, false},
		{"at threshold", 5, 5, true},
		{"above threshold", 9, 5, true},
		{"threshold one", 1, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoginLocked(tt.failures, tt.maxAttempts); got != tt.expected {
				t.Errorf("IsLoginLocked(%d, %d) = %v, want %v", tt.failures, tt.maxAttempts, got, tt.expected)
			}
		})
	}
}

func TestLoginThrottleKey(t *testing.T) {
	tests := []struct {
		ip       string
		expected string
	}{
		{"127.0.0.1", "admin:login_throttle:127.0.0.1"},
		{"", "admin:login_throttle:"},
	}

	for _, tt := range tests {
		if got := LoginThrottleKey(tt.ip); got != tt.expected {
			t.Errorf("LoginThrottleKey(%q) = %q, want %q", tt.ip, got, tt.expected)
		}
	}
}

func TestLoginThrottleMessage(t *testing.T) {
	tests := []struct {
		minutes  int
		contains string
	}{
		{10, "10"},
		{1, "1"},
	}

	for _, tt := range tests {
		got := LoginThrottleMessage(tt.minutes)
		if got == "" {
			t.Errorf("LoginThrottleMessage(%d) is empty", tt.minutes)
		}

		if !strings.Contains(got, tt.contains) {
			t.Errorf("LoginThrottleMessage(%d) = %q, want contains %q", tt.minutes, got, tt.contains)
		}
	}
}
