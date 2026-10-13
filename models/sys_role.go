package models

// SysRole is the RuoYi `sys_role` table.
type SysRole struct {
	RoleID            uint   `gorm:"column:role_id;primaryKey;autoIncrement" json:"roleId"`
	RoleName          string `gorm:"column:role_name" json:"roleName"`
	RoleKey           string `gorm:"column:role_key" json:"roleKey"`
	RoleSort          int    `gorm:"column:role_sort" json:"roleSort"`
	DataScope         string `gorm:"column:data_scope" json:"dataScope"`
	MenuCheckStrictly bool   `gorm:"column:menu_check_strictly" json:"menuCheckStrictly"`
	DeptCheckStrictly bool   `gorm:"column:dept_check_strictly" json:"deptCheckStrictly"`
	Status            string `gorm:"column:status" json:"status"`
	DelFlag           string `gorm:"column:del_flag" json:"delFlag"`
	BaseModel

	// MenuIDs is transient, only used when reading/writing role-menu bindings.
	MenuIDs []uint `gorm:"-" json:"menuIds,omitempty"`
}

// TableName sets the table name of the model.
func (r *SysRole) TableName() string {
	return "sys_role"
}
