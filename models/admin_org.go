package models

import (
	"github.com/goravel/framework/database/orm"
)

// AdminDepartment 部门表（组织架构，树形 parent_id）。
type AdminDepartment struct {
	orm.Model
	ParentID    int     `gorm:"column:parent_id;default:0" json:"parent_id"`
	Name        string  `gorm:"column:name;size:50" json:"name"`
	CustomOrder int     `gorm:"column:custom_order;default:0" json:"custom_order"`
	Remark      *string `gorm:"column:remark" json:"remark"`
}

func (r *AdminDepartment) TableName() string { return "admin_departments" }

// AdminPosition 岗位表。
type AdminPosition struct {
	orm.Model
	Name        string  `gorm:"column:name;size:50;unique" json:"name"`
	CustomOrder int     `gorm:"column:custom_order;default:0" json:"custom_order"`
	Remark      *string `gorm:"column:remark" json:"remark"`
}

func (r *AdminPosition) TableName() string { return "admin_positions" }

// AdminDictionary 字典表（key 唯一，业务侧按 key 取选项）。
type AdminDictionary struct {
	orm.Model
	Name   string  `gorm:"column:name;size:50" json:"name"`
	Key    string  `gorm:"column:key;size:50;unique" json:"key"`
	Remark *string `gorm:"column:remark" json:"remark"`
}

func (r *AdminDictionary) TableName() string { return "admin_dictionaries" }

// AdminDictionaryItem 字典项（挂在字典下，custom_order 排序）。
type AdminDictionaryItem struct {
	orm.Model
	DictionaryID int    `gorm:"column:dictionary_id;index" json:"dictionary_id"`
	Label        string `gorm:"column:label;size:50" json:"label"`
	Value        string `gorm:"column:value;size:50" json:"value"`
	CustomOrder  int    `gorm:"column:custom_order;default:0" json:"custom_order"`
}

func (r *AdminDictionaryItem) TableName() string { return "admin_dictionary_items" }
