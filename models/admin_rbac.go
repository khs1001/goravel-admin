package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminRole 对齐 owl-admin admin_roles 表结构。
type AdminRole struct {
	orm.Model
	Name string `gorm:"column:name;size:50;unique" json:"name"`
	Slug string `gorm:"column:slug;size:50;unique" json:"slug"`
}

func (r *AdminRole) TableName() string { return "admin_roles" }

// AdminRoleUsers 角色-用户关联（owl-admin admin_role_users）。
type AdminRoleUsers struct {
	RoleID int `gorm:"column:role_id" json:"role_id"`
	UserID int `gorm:"column:user_id" json:"user_id"`
}

func (r *AdminRoleUsers) TableName() string { return "admin_role_users" }

// AdminPermissions 权限表（owl-admin admin_permissions）。
type AdminPermission struct {
	orm.Model
	Name        string  `gorm:"column:name;size:50;unique" json:"name"`
	Slug        string  `gorm:"column:slug;size:50;unique" json:"slug"`
	HTTPMethod  *string `gorm:"column:http_method" json:"http_method"`
	HTTPPath    *string `gorm:"column:http_path" json:"http_path"`
	CustomOrder int     `gorm:"column:custom_order;default:0" json:"custom_order"`
	ParentID    int     `gorm:"column:parent_id;default:0" json:"parent_id"`
}

func (r *AdminPermission) TableName() string { return "admin_permissions" }

// AdminRolePermissions 角色-权限关联（owl-admin admin_role_permissions）。
type AdminRolePermissions struct {
	RoleID       int `gorm:"column:role_id" json:"role_id"`
	PermissionID int `gorm:"column:permission_id" json:"permission_id"`
}

func (r *AdminRolePermissions) TableName() string { return "admin_role_permissions" }

// AdminPermissionMenu 权限-菜单关联（owl-admin admin_permission_menu）。
type AdminPermissionMenu struct {
	PermissionID int `gorm:"column:permission_id" json:"permission_id"`
	MenuID       int `gorm:"column:menu_id" json:"menu_id"`
}

func (r *AdminPermissionMenu) TableName() string { return "admin_permission_menu" }
