package models

import (
	"github.com/goravel/framework/support/carbon"
)

// AuditModel carries the common RuoYi audit columns (create_by/create_time/
// update_by/update_time) shared by most tables.
type AuditModel struct {
	CreateBy   string           `gorm:"column:create_by" json:"createBy"`
	CreateTime *carbon.DateTime `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateBy   string           `gorm:"column:update_by" json:"updateBy"`
	UpdateTime *carbon.DateTime `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

// BaseModel is AuditModel plus the RuoYi `remark` column.
type BaseModel struct {
	AuditModel
	Remark string `gorm:"column:remark" json:"remark"`
}

// SysUser is the RuoYi `sys_user` table.
type SysUser struct {
	UserID      uint             `gorm:"column:user_id;primaryKey;autoIncrement" json:"userId"`
	DeptID      uint             `gorm:"column:dept_id" json:"deptId"`
	UserName    string           `gorm:"column:user_name" json:"userName"`
	NickName    string           `gorm:"column:nick_name" json:"nickName"`
	UserType    string           `gorm:"column:user_type" json:"userType"`
	Email       string           `gorm:"column:email" json:"email"`
	Phonenumber string           `gorm:"column:phonenumber" json:"phonenumber"`
	Sex         string           `gorm:"column:sex" json:"sex"`
	Avatar      string           `gorm:"column:avatar" json:"avatar"`
	Password    string           `gorm:"column:password" json:"-"`
	Status      string           `gorm:"column:status" json:"status"`
	DelFlag     string           `gorm:"column:del_flag" json:"delFlag"`
	LoginIP     string           `gorm:"column:login_ip" json:"loginIp"`
	LoginDate   *carbon.DateTime `gorm:"column:login_date" json:"loginDate"`
	BaseModel

	// RoleIDs is transient, only used when reading/writing user-role bindings.
	RoleIDs []uint `gorm:"-" json:"roleIds,omitempty"`
	// PostIDs is transient, only used when reading/writing user-post bindings.
	PostIDs []uint `gorm:"-" json:"postIds,omitempty"`
}

// TableName sets the table name of the model.
func (r *SysUser) TableName() string {
	return "sys_user"
}

// IsAdmin reports whether the user is the built-in super administrator (id 1).
func (r *SysUser) IsAdmin() bool {
	return r.UserID == 1
}
