package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminSetting 对齐 owl-admin admin_settings 表结构（key 为主键，values 为 JSON 文本）。
// 注意：Key 不能加 default 标签——GORM 会对带 default 的零值字段在 INSERT 时使用列默认值。
type AdminSetting struct {
	Key    string  `gorm:"column:key;primaryKey" json:"key"`
	Values *string `gorm:"column:values" json:"values"`
	orm.Timestamps
}

func (r *AdminSetting) TableName() string { return "admin_settings" }

// AdminPage 对齐 owl-admin admin_pages 表结构（可视化编辑器存储的页面 schema）。
type AdminPage struct {
	orm.Model
	Title  string `gorm:"column:title" json:"title"`
	Sign   string `gorm:"column:sign" json:"sign"`
	Schema string `gorm:"column:schema" json:"schema"`
}

func (r *AdminPage) TableName() string { return "admin_pages" }
