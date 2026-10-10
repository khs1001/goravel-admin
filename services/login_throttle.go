package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
)

// 登录限流默认参数（可通过环境变量覆盖）。
const (
	DefaultLoginMaxAttempts = 5  // 最大连续失败次数
	DefaultLoginLockMinutes = 10 // 锁定时长（分钟）
)

// envInt 从环境读取整型值（Config 契约未提供 EnvInt，经字符串转换）。
func envInt(name string, fallback int) int {
	raw := facades.Config().EnvString(name, "")
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}

	return value
}

// loginMaxAttempts 从环境读取最大失败次数。
func loginMaxAttempts() int {
	return envInt("ADMIN_LOGIN_MAX_ATTEMPTS", DefaultLoginMaxAttempts)
}

// loginLockDuration 从环境读取锁定时长。
func loginLockDuration() time.Duration {
	minutes := envInt("ADMIN_LOGIN_LOCK_MINUTES", DefaultLoginLockMinutes)
	if minutes <= 0 {
		minutes = DefaultLoginLockMinutes
	}

	return time.Duration(minutes) * time.Minute
}

// LoginThrottleKey 生成基于 IP 的限流缓存 key。
func LoginThrottleKey(ip string) string {
	return "admin:login_throttle:" + ip
}

// IsLoginLocked 判定是否达到锁定阈值（纯函数）。
func IsLoginLocked(failures, maxAttempts int) bool {
	return failures >= maxAttempts
}

// LoginThrottleMessage 锁定提示文案（纯函数）。
func LoginThrottleMessage(lockMinutes int) string {
	return fmt.Sprintf("尝试次数过多，账号已临时锁定，请 %d 分钟后重试", lockMinutes)
}

// LoginThrottle 基于 Cache 的登录失败限流（按 IP 维度）。
// 连续失败达到 ADMIN_LOGIN_MAX_ATTEMPTS 次（默认 5）后锁定
// ADMIN_LOGIN_LOCK_MINUTES 分钟（默认 10），登录成功后清零。
type LoginThrottle struct{}

func NewLoginThrottle() *LoginThrottle { return &LoginThrottle{} }

// TooManyAttempts 是否已锁定。
func (r *LoginThrottle) TooManyAttempts(ip string) bool {
	return IsLoginLocked(facades.Cache().GetInt(LoginThrottleKey(ip)), loginMaxAttempts())
}

// Hit 记录一次失败。滑动窗口：每次失败刷新 TTL，持续尝试则持续锁定。
// （memory 驱动的 Increment 仅支持原子指针类型，为跨驱动一致采用读-写计数。）
func (r *LoginThrottle) Hit(ip string) {
	key := LoginThrottleKey(ip)
	count := facades.Cache().GetInt(key) + 1
	if err := facades.Cache().Put(key, count, loginLockDuration()); err != nil {
		facades.Log().Warningf("[admin] throttle put failed for %s: %v", ip, err)
	}
}

// Clear 登录成功（或手动解锁）后清零。
func (r *LoginThrottle) Clear(ip string) {
	facades.Cache().Forget(LoginThrottleKey(ip))
}
