package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000005CreateOperationLogs struct{}

// Signature 迁移唯一签名。
func (r *M20260101000005CreateOperationLogs) Signature() string {
	return "20260101000005_create_operation_logs"
}

// Up 建操作日志表，并在"系统管理"下追加"操作日志"菜单与查看权限（幂等）。
func (r *M20260101000005CreateOperationLogs) Up() error {
	if !facades.Schema().HasTable("admin_operation_logs") {
		if err := facades.Schema().Create("admin_operation_logs", func(table schema.Blueprint) {
			table.ID()
			table.Integer("user_id").Default(0)
			table.String("username", 64).Default("")
			table.String("action", 16)
			table.String("resource", 64).Default("")
			table.String("resource_id", 32).Default("")
			table.String("method", 10).Default("")
			table.String("path", 255).Default("")
			table.String("ip", 64).Default("")
			table.Text("input").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Index("user_id", "resource", "created_at")
		}); err != nil {
			return err
		}
	}

	return r.seed()
}

// seed 操作日志菜单 + 查看权限（幂等）。
func (r *M20260101000005CreateOperationLogs) seed() error {
	count, err := facades.Orm().Query().Model(&models.AdminPermission{}).
		Where("slug", "system.admin_operation_logs.index").Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var systemMenu models.AdminMenu
	if err := facades.Orm().Query().
		Where("parent_id", 0).Where("title", "系统管理").First(&systemMenu); err != nil || systemMenu.ID == 0 {
		return nil
	}

	menu := &models.AdminMenu{
		ParentID: int(systemMenu.ID), CustomOrder: 8, Title: "操作日志",
		URL: strPtr("/system/admin_operation_logs"), UrlType: 1, Visible: 1,
	}
	if err := facades.Orm().Query().Create(menu); err != nil {
		return err
	}

	path := "GET /admin-api/system/admin_operation_logs*"
	permission := &models.AdminPermission{
		Name:     "日志查看",
		Slug:     "system.admin_operation_logs.index",
		HTTPPath: &path,
	}
	if err := facades.Orm().Query().Create(permission); err != nil {
		return err
	}

	links := []models.AdminPermissionMenu{
		{PermissionID: int(permission.ID), MenuID: int(systemMenu.ID)},
		{PermissionID: int(permission.ID), MenuID: int(menu.ID)},
	}

	return facades.Orm().Query().Create(&links)
}

// Down 删表。
func (r *M20260101000005CreateOperationLogs) Down() error {
	return facades.Schema().DropIfExists("admin_operation_logs")
}
