package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000008CreateLoginLogs struct{}

// Signature 迁移唯一签名。
func (r *M20260101000008CreateLoginLogs) Signature() string {
	return "20260101000008_create_login_logs"
}

// Up 建登录日志表，并在"系统管理"下追加"登录日志"菜单与查看权限（幂等）。
func (r *M20260101000008CreateLoginLogs) Up() error {
	if !facades.Schema().HasTable("admin_login_logs") {
		if err := facades.Schema().Create("admin_login_logs", func(table schema.Blueprint) {
			table.ID()
			table.Integer("user_id").Default(0)
			table.String("username", 64).Default("")
			table.String("status", 16)
			table.String("message", 255).Default("")
			table.String("ip", 64).Default("")
			table.String("user_agent", 255).Default("")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Index("user_id", "status", "created_at")
		}); err != nil {
			return err
		}
	}

	return r.seed()
}

// seed 登录日志菜单 + 查看权限（幂等）。
func (r *M20260101000008CreateLoginLogs) seed() error {
	count, err := facades.Orm().Query().Model(&models.AdminPermission{}).
		Where("slug", "system.admin_login_logs.index").Count()
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
		ParentID: int(systemMenu.ID), CustomOrder: 9, Title: "登录日志",
		URL: strPtr("/system/admin_login_logs"), UrlType: 1, Visible: 1,
	}
	if err := facades.Orm().Query().Create(menu); err != nil {
		return err
	}

	path := "GET /admin-api/system/admin_login_logs*"
	permission := &models.AdminPermission{
		Name:     "登录日志查看",
		Slug:     "system.admin_login_logs.index",
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
func (r *M20260101000008CreateLoginLogs) Down() error {
	return facades.Schema().DropIfExists("admin_login_logs")
}
