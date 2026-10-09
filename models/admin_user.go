package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminUser 对齐 owl-admin admin_users 表结构。
type AdminUser struct {
	orm.Model
	Username      string  `gorm:"column:username;size:120;unique" json:"username"`
	Password      string  `gorm:"column:password;size:80" json:"-"`
	Enabled       int     `gorm:"column:enabled;default:1" json:"enabled"`
	Name          string  `gorm:"column:name;default:''" json:"name"`
	Avatar        *string `gorm:"column:avatar" json:"avatar"`
	RememberToken *string `gorm:"column:remember_token;size:100" json:"-"`
}

func (r *AdminUser) TableName() string { return "admin_users" }

// IsEnabled 是否启用。
func (r *AdminUser) IsEnabled() bool { return r.Enabled == 1 }
