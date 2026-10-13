package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

// CreateSysTables creates the RuoYi-style admin tables.
type CreateSysTables struct{}

// Signature The unique signature for the migration.
func (r *CreateSysTables) Signature() string {
	return "20260101000000_create_sys_tables"
}

// Up Run the migrations.
func (r *CreateSysTables) Up() error {
	if !facades.Schema().HasTable("sys_user") {
		if err := facades.Schema().Create("sys_user", func(table schema.Blueprint) {
			table.Increments("user_id")
			table.String("user_name", 30).Default("")
			table.String("nick_name", 30).Default("")
			table.String("user_type", 2).Default("00")
			table.String("email", 50).Default("")
			table.String("phonenumber", 11).Default("")
			table.String("sex", 1).Default("0")
			table.String("avatar", 100).Default("")
			table.String("password", 100).Default("")
			table.String("status", 1).Default("0")
			table.String("del_flag", 1).Default("0")
			table.String("login_ip", 128).Default("")
			table.DateTime("login_date").Nullable()
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
			table.Unique("user_name")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_role") {
		if err := facades.Schema().Create("sys_role", func(table schema.Blueprint) {
			table.Increments("role_id")
			table.String("role_name", 30).Default("")
			table.String("role_key", 100).Default("")
			table.Integer("role_sort").Default(0)
			table.String("data_scope", 1).Default("1")
			table.Boolean("menu_check_strictly").Default(true)
			table.Boolean("dept_check_strictly").Default(true)
			table.String("status", 1).Default("0")
			table.String("del_flag", 1).Default("0")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_menu") {
		if err := facades.Schema().Create("sys_menu", func(table schema.Blueprint) {
			table.Increments("menu_id")
			table.String("menu_name", 50).Default("")
			table.Integer("parent_id").Default(0)
			table.Integer("order_num").Default(0)
			table.String("path", 200).Default("")
			table.String("component", 255).Nullable()
			table.String("query", 255).Nullable()
			table.Integer("is_frame").Default(1)
			table.Integer("is_cache").Default(0)
			table.String("menu_type", 1).Nullable()
			table.String("visible", 1).Default("0")
			table.String("status", 1).Default("0")
			table.String("perms", 100).Nullable()
			table.String("icon", 100).Default("#")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_user_role") {
		if err := facades.Schema().Create("sys_user_role", func(table schema.Blueprint) {
			table.Integer("user_id")
			table.Integer("role_id")
			table.Index("user_id")
			table.Index("role_id")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_role_menu") {
		if err := facades.Schema().Create("sys_role_menu", func(table schema.Blueprint) {
			table.Integer("role_id")
			table.Integer("menu_id")
			table.Index("role_id")
			table.Index("menu_id")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_config") {
		if err := facades.Schema().Create("sys_config", func(table schema.Blueprint) {
			table.Increments("config_id")
			table.String("config_name", 100).Default("")
			table.String("config_key", 100).Default("")
			table.String("config_value", 500).Default("")
			table.String("config_type", 1).Default("N")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *CreateSysTables) Down() error {
	for _, table := range []string{"sys_user", "sys_role", "sys_menu", "sys_user_role", "sys_role_menu", "sys_config"} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}

	return nil
}
