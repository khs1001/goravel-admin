package models

// SysExtension records a discovered admin plugin and its enable state.
type SysExtension struct {
	ExtensionID uint   `gorm:"column:extension_id;primaryKey;autoIncrement" json:"extensionId"`
	Name        string `gorm:"column:name" json:"name"`
	IsEnabled   int    `gorm:"column:is_enabled" json:"isEnabled"`
	Version     string `gorm:"column:version" json:"version"`
	Options     string `gorm:"column:options" json:"options"`
	AuditModel
}

// TableName sets the table name of the model.
func (r *SysExtension) TableName() string {
	return "sys_extension"
}
