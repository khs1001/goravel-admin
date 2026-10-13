package models

// SysDictType is the RuoYi `sys_dict_type` table (dictionary categories).
type SysDictType struct {
	DictID   uint   `gorm:"column:dict_id;primaryKey;autoIncrement" json:"dictId"`
	DictName string `gorm:"column:dict_name" json:"dictName"`
	DictType string `gorm:"column:dict_type" json:"dictType"`
	Status   string `gorm:"column:status" json:"status"`
	BaseModel
}

// TableName sets the table name of the model.
func (r *SysDictType) TableName() string {
	return "sys_dict_type"
}

// SysDictData is the RuoYi `sys_dict_data` table (dictionary entries).
type SysDictData struct {
	DictCode  uint   `gorm:"column:dict_code;primaryKey;autoIncrement" json:"dictCode"`
	DictSort  int    `gorm:"column:dict_sort" json:"dictSort"`
	DictLabel string `gorm:"column:dict_label" json:"dictLabel"`
	DictValue string `gorm:"column:dict_value" json:"dictValue"`
	DictType  string `gorm:"column:dict_type" json:"dictType"`
	CssClass  string `gorm:"column:css_class" json:"cssClass"`
	ListClass string `gorm:"column:list_class" json:"listClass"`
	IsDefault string `gorm:"column:is_default" json:"isDefault"`
	Status    string `gorm:"column:status" json:"status"`
	BaseModel
}

// TableName sets the table name of the model.
func (r *SysDictData) TableName() string {
	return "sys_dict_data"
}
