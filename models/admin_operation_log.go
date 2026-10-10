package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminOperationLog 操作日志表。
//
// 三类来源共用一张表，以 action 区分：
//   - request：HTTP 中间件记录的写请求（POST/PUT/DELETE）
//   - created / updated / deleted：ORM 模型事件记录的数据变更
//
// changes 仅 updated 动作记录：字段级变更 JSON {col: {old, new}}。
type AdminOperationLog struct {
	orm.Model
	UserID     int     `gorm:"column:user_id;default:0" json:"user_id"`
	Username   string  `gorm:"column:username;size:64;default:''" json:"username"`
	Action     string  `gorm:"column:action;size:16" json:"action"`
	Resource   string  `gorm:"column:resource;size:64;default:''" json:"resource"`
	ResourceID string  `gorm:"column:resource_id;size:32;default:''" json:"resource_id"`
	Method     string  `gorm:"column:method;size:10;default:''" json:"method"`
	Path       string  `gorm:"column:path;size:255;default:''" json:"path"`
	IP         string  `gorm:"column:ip;size:64;default:''" json:"ip"`
	Input      *string `gorm:"column:input" json:"input"`
	Changes    *string `gorm:"column:changes" json:"changes"`
}

func (r *AdminOperationLog) TableName() string { return "admin_operation_logs" }
