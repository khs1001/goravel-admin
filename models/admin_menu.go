package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminMenu 对齐 owl-admin admin_menus 表结构。
type AdminMenu struct {
	orm.Model
	ParentID    int     `gorm:"column:parent_id;default:0" json:"parent_id"`
	CustomOrder int     `gorm:"column:custom_order;default:0" json:"custom_order"`
	Title       string  `gorm:"column:title;size:100" json:"title"`
	Icon        *string `gorm:"column:icon" json:"icon"`
	URL         *string `gorm:"column:url" json:"url"`
	UrlType     int     `gorm:"column:url_type;default:1" json:"url_type"`
	Visible     int     `gorm:"column:visible;default:1" json:"visible"`
	IsHome      int     `gorm:"column:is_home;default:0" json:"is_home"`
	KeepAlive   *int    `gorm:"column:keep_alive" json:"keep_alive"`
	IframeURL   *string `gorm:"column:iframe_url" json:"iframe_url"`
	Component   *string `gorm:"column:component" json:"component"`
	IsFull      int     `gorm:"column:is_full;default:0" json:"is_full"`
	Extension   *string `gorm:"column:extension" json:"extension"`
}

func (r *AdminMenu) TableName() string { return "admin_menus" }

// GetURL 菜单路由（null 视为空字符串）。
func (r *AdminMenu) GetURL() string {
	if r.URL == nil {
		return ""
	}

	return *r.URL
}

// GetComponent 菜单组件（null 视为空字符串）。
func (r *AdminMenu) GetComponent() string {
	if r.Component == nil {
		return ""
	}

	return *r.Component
}

// GetIcon 菜单图标（null 视为空字符串）。
func (r *AdminMenu) GetIcon() string {
	if r.Icon == nil {
		return ""
	}

	return *r.Icon
}
