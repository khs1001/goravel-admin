package models

// SysConfig is the RuoYi `sys_config` parameter table. It also backs the
// Owl-style `/_settings` persistence (key -> JSON string value).
type SysConfig struct {
	ConfigID    uint   `gorm:"column:config_id;primaryKey;autoIncrement" json:"configId"`
	ConfigName  string `gorm:"column:config_name" json:"configName"`
	ConfigKey   string `gorm:"column:config_key" json:"configKey"`
	ConfigValue string `gorm:"column:config_value" json:"configValue"`
	ConfigType  string `gorm:"column:config_type" json:"configType"`
	BaseModel
}

// TableName sets the table name of the model.
func (r *SysConfig) TableName() string {
	return "sys_config"
}
