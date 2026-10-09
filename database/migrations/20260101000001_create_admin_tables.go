package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000001CreateAdminTables struct{}

// Signature 迁移唯一签名。
func (r *M20260101000001CreateAdminTables) Signature() string {
	return "20260101000001_create_admin_tables"
}

// Up 建表并写入初始数据（对齐 owl-admin src/Support/Cores/Database.php）。
func (r *M20260101000001CreateAdminTables) Up() error {
	if !facades.Schema().HasTable("admin_users") {
		if err := facades.Schema().Create("admin_users", func(table schema.Blueprint) {
			table.ID()
			table.String("username", 120)
			table.String("password", 80)
			table.UnsignedTinyInteger("enabled").Default(1)
			table.String("name").Default("")
			table.String("avatar").Nullable()
			table.String("remember_token", 100).Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("username")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_roles") {
		if err := facades.Schema().Create("admin_roles", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 50)
			table.String("slug", 50)
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
			table.Unique("slug")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_permissions") {
		if err := facades.Schema().Create("admin_permissions", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 50)
			table.String("slug", 50)
			table.Text("http_method").Nullable()
			table.Text("http_path").Nullable()
			table.Integer("custom_order").Default(0)
			table.Integer("parent_id").Default(0)
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
			table.Unique("slug")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_menus") {
		if err := facades.Schema().Create("admin_menus", func(table schema.Blueprint) {
			table.ID()
			table.Integer("parent_id").Default(0)
			table.Integer("custom_order").Default(0)
			table.String("title", 100)
			table.String("icon", 255).Nullable()
			table.String("url").Nullable()
			table.UnsignedTinyInteger("url_type").Default(1)
			table.UnsignedTinyInteger("visible").Default(1)
			table.UnsignedTinyInteger("is_home").Default(0)
			table.UnsignedTinyInteger("keep_alive").Nullable()
			table.String("iframe_url").Nullable()
			table.String("component").Nullable()
			table.UnsignedTinyInteger("is_full").Default(0)
			table.String("extension").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_role_users") {
		if err := facades.Schema().Create("admin_role_users", func(table schema.Blueprint) {
			table.Integer("role_id")
			table.Integer("user_id")
			table.Index("role_id", "user_id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_role_permissions") {
		if err := facades.Schema().Create("admin_role_permissions", func(table schema.Blueprint) {
			table.Integer("role_id")
			table.Integer("permission_id")
			table.Index("role_id", "permission_id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_permission_menu") {
		if err := facades.Schema().Create("admin_permission_menu", func(table schema.Blueprint) {
			table.Integer("permission_id")
			table.Integer("menu_id")
			table.Index("permission_id", "menu_id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_settings") {
		if err := facades.Schema().Create("admin_settings", func(table schema.Blueprint) {
			table.String("key")
			table.LongText("values").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Primary("key")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_extensions") {
		if err := facades.Schema().Create("admin_extensions", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 100)
			table.UnsignedTinyInteger("is_enabled").Default(0)
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_pages") {
		if err := facades.Schema().Create("admin_pages", func(table schema.Blueprint) {
			table.ID()
			table.String("title")
			table.String("sign")
			table.LongText("schema")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
		}); err != nil {
			return err
		}
	}

	return r.seed()
}

// Down 回滚全部 admin 表。
func (r *M20260101000001CreateAdminTables) Down() error {
	for _, table := range []string{
		"admin_users", "admin_roles", "admin_permissions", "admin_menus",
		"admin_role_users", "admin_role_permissions", "admin_permission_menu",
		"admin_settings", "admin_extensions", "admin_pages",
	} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}

	return nil
}

// seed 初始数据：仅当 admin_users 为空时写入一次。
func (r *M20260101000001CreateAdminTables) seed() error {
	count, err := facades.Orm().Query().Model(&models.AdminUser{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashed, err := facades.Hash().Make("admin")
	if err != nil {
		return err
	}

	// 超管角色
	role := &models.AdminRole{Name: "Administrator", Slug: "administrator"}
	if err := facades.Orm().Query().Create(role); err != nil {
		return err
	}

	// 超管账号 admin / admin
	user := &models.AdminUser{
		Username: "admin",
		Password: hashed,
		Enabled:  1,
		Name:     "Administrator",
	}
	if err := facades.Orm().Query().Create(user); err != nil {
		return err
	}

	if err := facades.Orm().Query().Create(&models.AdminRoleUsers{
		RoleID: int(role.ID),
		UserID: int(user.ID),
	}); err != nil {
		return err
	}

	// 基础菜单树（第一阶段：仅控制台）
	menus := []models.AdminMenu{
		{ParentID: 0, CustomOrder: 100, Title: "仪表盘", URL: strPtr("/dashboard"), UrlType: 1, Visible: 1, IsHome: 1},
	}
	if err := facades.Orm().Query().Create(&menus); err != nil {
		return err
	}

	return nil
}

func strPtr(s string) *string { return &s }
