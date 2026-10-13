package services

import (
	"encoding/json"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// ConfigService reads and writes the RuoYi `sys_config` table.
//
// The admin settings (`/_settings`) are persisted here: Owl setting keys such
// as `system_theme_setting` and `admin_locale` map to `config_key`, their JSON
// payloads to `config_value`.
type ConfigService struct{}

// NewConfigService creates a ConfigService.
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// Get returns the raw string value of a config key, or def when missing.
func (s *ConfigService) Get(key, def string) string {
	var config models.SysConfig
	if err := facades.Orm().Query().Where("config_key", key).FirstOrFail(&config); err != nil {
		return def
	}
	if config.ConfigValue == "" {
		return def
	}

	return config.ConfigValue
}

// Set stores a string value for a config key, creating the row if needed.
func (s *ConfigService) Set(key, value string) error {
	var config models.SysConfig
	err := facades.Orm().Query().Where("config_key", key).FirstOrFail(&config)
	if err != nil {
		return facades.Orm().Query().Create(&models.SysConfig{
			ConfigName:  key,
			ConfigKey:   key,
			ConfigValue: value,
			ConfigType:  "N",
		})
	}

	config.ConfigValue = value

	return facades.Orm().Query().Save(&config)
}

// GetJSON decodes a config value as a JSON object, returning nil when absent or
// malformed.
func (s *ConfigService) GetJSON(key string) map[string]any {
	raw := s.Get(key, "")
	if raw == "" {
		return nil
	}

	out := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}

	return out
}

// SetJSON encodes a value as JSON and stores it.
func (s *ConfigService) SetJSON(key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.Set(key, string(raw))
}
