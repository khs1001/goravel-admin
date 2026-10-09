package services

import (
	"encoding/json"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/models"
)

// PageService 页面 schema 服务（双轨）：
//  1. 代码注册表：内置页面与插件 Pages() 声明，开箱即用；
//  2. DB 覆盖层：admin_pages 表（可视化编辑器），sign 命中时优先返回。
//
// 对齐 owl-admin src/Services/AdminPageService.php::get。
type PageService struct {
	builtin map[string]any
}

func NewPageService(builtin map[string]any, plugins []contracts.Plugin) *PageService {
	registry := map[string]any{}
	for sign, schema := range builtin {
		registry[sign] = schema
	}
	for _, plugin := range plugins {
		for sign, schema := range plugin.Pages() {
			registry[sign] = schema
		}
	}

	return &PageService{builtin: registry}
}

// Get 按 sign 获取页面 schema：先查 DB 覆盖层，再查代码注册表。
func (r *PageService) Get(sign string) any {
	if sign == "" {
		return nil
	}

	var page models.AdminPage
	if err := facades.Orm().Query().Where("sign", sign).First(&page); err == nil && page.Sign != "" && page.Schema != "" {
		var schema any
		if err := json.Unmarshal([]byte(page.Schema), &schema); err == nil {
			return schema
		}

		return page.Schema
	}

	if schema, ok := r.builtin[sign]; ok {
		return schema
	}

	return nil
}
