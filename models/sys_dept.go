package models

// SysDept is the RuoYi `sys_dept` table (organisation tree).
type SysDept struct {
	DeptID    uint   `gorm:"column:dept_id;primaryKey;autoIncrement" json:"deptId"`
	ParentID  uint   `gorm:"column:parent_id" json:"parentId"`
	Ancestors string `gorm:"column:ancestors" json:"ancestors"`
	DeptName  string `gorm:"column:dept_name" json:"deptName"`
	OrderNum  int    `gorm:"column:order_num" json:"orderNum"`
	Leader    string `gorm:"column:leader" json:"leader"`
	Phone     string `gorm:"column:phone" json:"phone"`
	Email     string `gorm:"column:email" json:"email"`
	Status    string `gorm:"column:status" json:"status"`
	DelFlag   string `gorm:"column:del_flag" json:"delFlag"`
	AuditModel

	// Children is transient, used to render the department tree.
	Children []SysDept `gorm:"-" json:"children,omitempty"`
}

// TableName sets the table name of the model.
func (r *SysDept) TableName() string {
	return "sys_dept"
}

// SysUserPost is the RuoYi `sys_user_post` join table.
type SysUserPost struct {
	UserID uint `gorm:"column:user_id;primaryKey" json:"userId"`
	PostID uint `gorm:"column:post_id;primaryKey" json:"postId"`
}

// TableName sets the table name of the model.
func (r *SysUserPost) TableName() string {
	return "sys_user_post"
}

// SysRoleDept is the RuoYi `sys_role_dept` join table, used for the "custom"
// data scope of a role.
type SysRoleDept struct {
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
	DeptID uint `gorm:"column:dept_id;primaryKey" json:"deptId"`
}

// TableName sets the table name of the model.
func (r *SysRoleDept) TableName() string {
	return "sys_role_dept"
}
