package models

import (
	"github.com/goravel/framework/support/carbon"
)

// Business types for SysOperLog, following RuoYi.
const (
	BusinessTypeOther  = 0 // 其它
	BusinessTypeInsert = 1 // 新增
	BusinessTypeUpdate = 2 // 修改
	BusinessTypeDelete = 3 // 删除
	BusinessTypeGrant  = 4 // 授权
	BusinessTypeExport = 5 // 导出
	BusinessTypeImport = 6 // 导入
	BusinessTypeClean  = 9 // 清空
)

// SysOperLog is the RuoYi `sys_oper_log` table (operation log).
type SysOperLog struct {
	OperID        uint             `gorm:"column:oper_id;primaryKey;autoIncrement" json:"operId"`
	Title         string           `gorm:"column:title" json:"title"`
	BusinessType  int              `gorm:"column:business_type" json:"businessType"`
	Method        string           `gorm:"column:method" json:"method"`
	RequestMethod string           `gorm:"column:request_method" json:"requestMethod"`
	OperatorType  int              `gorm:"column:operator_type" json:"operatorType"`
	OperName      string           `gorm:"column:oper_name" json:"operName"`
	OperURL       string           `gorm:"column:oper_url" json:"operUrl"`
	OperIP        string           `gorm:"column:oper_ip" json:"operIp"`
	OperParam     string           `gorm:"column:oper_param" json:"operParam"`
	JSONResult    string           `gorm:"column:json_result" json:"jsonResult"`
	Status        int              `gorm:"column:status" json:"status"`
	ErrorMsg      string           `gorm:"column:error_msg" json:"errorMsg"`
	OperTime      *carbon.DateTime `gorm:"column:oper_time" json:"operTime"`
	CostTime      int              `gorm:"column:cost_time" json:"costTime"`
}

// TableName sets the table name of the model.
func (r *SysOperLog) TableName() string {
	return "sys_oper_log"
}

// SysLogininfor is the RuoYi `sys_logininfor` table (login log).
type SysLogininfor struct {
	InfoID    uint             `gorm:"column:info_id;primaryKey;autoIncrement" json:"infoId"`
	UserName  string           `gorm:"column:user_name" json:"userName"`
	Ipaddr    string           `gorm:"column:ipaddr" json:"ipaddr"`
	Browser   string           `gorm:"column:browser" json:"browser"`
	Os        string           `gorm:"column:os" json:"os"`
	Status    int              `gorm:"column:status" json:"status"`
	Msg       string           `gorm:"column:msg" json:"msg"`
	LoginTime *carbon.DateTime `gorm:"column:login_time" json:"loginTime"`
}

// TableName sets the table name of the model.
func (r *SysLogininfor) TableName() string {
	return "sys_logininfor"
}
