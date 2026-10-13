package services

import (
	"strings"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// LogService persists operation and login logs.
type LogService struct{}

// NewLogService creates a LogService.
func NewLogService() *LogService {
	return &LogService{}
}

// Operation stores an operation log entry.
func (s *LogService) Operation(entry models.SysOperLog) {
	_ = facades.Orm().Query().Create(&entry)
}

// Login stores a login log entry.
func (s *LogService) Login(entry models.SysLogininfor) {
	_ = facades.Orm().Query().Create(&entry)
}

// ClientInfo derives a coarse browser and operating system name from a
// User-Agent header.
func ClientInfo(userAgent string) (browser, os string) {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "edg"):
		browser = "Edge"
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	default:
		browser = "Unknown"
	}

	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	default:
		os = "Unknown"
	}

	return browser, os
}
