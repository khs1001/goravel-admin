package models

import (
	"github.com/goravel/framework/database/orm"
)

// 登录日志状态常量。
const (
	LoginStatusSuccess   = "success"   // 登录成功
	LoginStatusFailed    = "failed"    // 用户名或密码错误
	LoginStatusThrottled = "throttled" // 触发限流锁定
	LoginStatusDisabled  = "disabled"  // 账号被禁用
)

// AdminLoginLog 登录日志表：登录成功/失败/限流/禁用全量留痕（安全审计）。
type AdminLoginLog struct {
	orm.Model
	UserID    int    `gorm:"column:user_id;default:0" json:"user_id"`
	Username  string `gorm:"column:username;size:64;default:''" json:"username"`
	Status    string `gorm:"column:status;size:16" json:"status"`
	Message   string `gorm:"column:message;size:255;default:''" json:"message"`
	IP        string `gorm:"column:ip;size:64;default:''" json:"ip"`
	UserAgent string `gorm:"column:user_agent;size:255;default:''" json:"user_agent"`
}

func (r *AdminLoginLog) TableName() string { return "admin_login_logs" }
