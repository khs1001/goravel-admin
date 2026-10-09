package services

import (
	"encoding/json"
	"strings"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// SettingService 系统设置服务（admin_settings 表，key → JSON values）。
// 逻辑对齐 owl-admin src/Services/AdminSettingService.php。
type SettingService struct{}

func NewSettingService() *SettingService { return &SettingService{} }

// Get 读取设置项并反序列化；不存在时返回 nil。
func (r *SettingService) Get(key string) any {
	var setting models.AdminSetting
	if err := facades.Orm().Query().Where("`key`", key).First(&setting); err != nil {
		// 记录不存在属正常情况；仅记录真实错误
		if !strings.Contains(err.Error(), "record not found") {
			facades.Log().Errorf("[admin] get setting %s failed: %v", key, err)
		}

		return nil
	}

	if setting.Key == "" || setting.Values == nil || *setting.Values == "" {
		return nil
	}

	var value any
	if err := json.Unmarshal([]byte(*setting.Values), &value); err != nil {
		return *setting.Values
	}

	return value
}

// GetByModule 读取设置项，不存在时返回默认值。
func (r *SettingService) GetByModule(key string, def any) any {
	value := r.Get(key)
	if value == nil {
		return def
	}

	return value
}

// SetMany 批量保存设置项（值序列化为 JSON 文本存储）。
func (r *SettingService) SetMany(data map[string]any) error {
	for key, value := range data {
		raw, err := json.Marshal(value)
		if err != nil {
			return err
		}

		encoded := string(raw)
		if err := facades.Orm().Query().UpdateOrCreate(
			&models.AdminSetting{},
			map[string]any{"key": key},
			map[string]any{"values": encoded},
		); err != nil {
			return err
		}
	}

	return nil
}
