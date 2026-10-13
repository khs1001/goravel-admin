package models

// Menu types, following RuoYi.
const (
	MenuTypeDir    = "M" // 目录
	MenuTypeMenu   = "C" // 菜单
	MenuTypeButton = "F" // 按钮
)

// sys_menu menu_type values are stored as a single character (M/C/F); the
// helpers below keep call sites readable.

// SysMenu is the RuoYi `sys_menu` table. It carries both navigation nodes
// (M/C) and button-level permissions (F, via the `perms` field).
type SysMenu struct {
	MenuID    uint   `gorm:"column:menu_id;primaryKey;autoIncrement" json:"menuId"`
	MenuName  string `gorm:"column:menu_name" json:"menuName"`
	ParentID  uint   `gorm:"column:parent_id" json:"parentId"`
	OrderNum  int    `gorm:"column:order_num" json:"orderNum"`
	Path      string `gorm:"column:path" json:"path"`
	Component string `gorm:"column:component" json:"component"`
	Query     string `gorm:"column:query" json:"query"`
	IsFrame   int    `gorm:"column:is_frame" json:"isFrame"`
	IsCache   int    `gorm:"column:is_cache" json:"isCache"`
	MenuType  string `gorm:"column:menu_type" json:"menuType"`
	Visible   string `gorm:"column:visible" json:"visible"`
	Status    string `gorm:"column:status" json:"status"`
	Perms     string `gorm:"column:perms" json:"perms"`
	Icon      string `gorm:"column:icon" json:"icon"`
	BaseModel

	// Children is transient, used to render the menu tree.
	Children []SysMenu `gorm:"-" json:"children,omitempty"`
}

// TableName sets the table name of the model.
func (r *SysMenu) TableName() string {
	return "sys_menu"
}

// SysUserRole is the RuoYi `sys_user_role` join table.
type SysUserRole struct {
	UserID uint `gorm:"column:user_id;primaryKey" json:"userId"`
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
}

// TableName sets the table name of the model.
func (r *SysUserRole) TableName() string {
	return "sys_user_role"
}

// SysRoleMenu is the RuoYi `sys_role_menu` join table.
type SysRoleMenu struct {
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
	MenuID uint `gorm:"column:menu_id;primaryKey" json:"menuId"`
}

// TableName sets the table name of the model.
func (r *SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
