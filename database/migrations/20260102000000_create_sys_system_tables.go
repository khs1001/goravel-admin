package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

// CreateSysSystemTables adds the RuoYi system-management tables (dept, post,
// dict) and the `sys_user.dept_id` column.
//
// It is written to apply on both fresh installs (after CreateSysTables) and
// existing installs, hence the HasTable/HasColumn guards.
type CreateSysSystemTables struct{}

// Signature The unique signature for the migration.
func (r *CreateSysSystemTables) Signature() string {
	return "20260102000000_create_sys_system_tables"
}

// Up Run the migrations.
func (r *CreateSysSystemTables) Up() error {
	if !facades.Schema().HasColumn("sys_user", "dept_id") {
		if err := facades.Schema().Table("sys_user", func(table schema.Blueprint) {
			table.Integer("dept_id").Default(0)
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_dept") {
		if err := facades.Schema().Create("sys_dept", func(table schema.Blueprint) {
			table.Increments("dept_id")
			table.Integer("parent_id").Default(0)
			table.String("ancestors", 50).Default("")
			table.String("dept_name", 50).Default("")
			table.Integer("order_num").Default(0)
			table.String("leader", 20).Default("")
			table.String("phone", 11).Default("")
			table.String("email", 50).Default("")
			table.String("status", 1).Default("0")
			table.String("del_flag", 1).Default("0")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_post") {
		if err := facades.Schema().Create("sys_post", func(table schema.Blueprint) {
			table.Increments("post_id")
			table.String("post_code", 64).Default("")
			table.String("post_name", 50).Default("")
			table.Integer("post_sort").Default(0)
			table.String("status", 1).Default("0")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_user_post") {
		if err := facades.Schema().Create("sys_user_post", func(table schema.Blueprint) {
			table.Integer("user_id")
			table.Integer("post_id")
			table.Index("user_id")
			table.Index("post_id")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_dict_type") {
		if err := facades.Schema().Create("sys_dict_type", func(table schema.Blueprint) {
			table.Increments("dict_id")
			table.String("dict_name", 100).Default("")
			table.String("dict_type", 100).Default("")
			table.String("status", 1).Default("0")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
			table.Unique("dict_type")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_dict_data") {
		if err := facades.Schema().Create("sys_dict_data", func(table schema.Blueprint) {
			table.Increments("dict_code")
			table.Integer("dict_sort").Default(0)
			table.String("dict_label", 100).Default("")
			table.String("dict_value", 100).Default("")
			table.String("dict_type", 100).Default("")
			table.String("css_class", 100).Nullable()
			table.String("list_class", 100).Nullable()
			table.String("is_default", 1).Default("N")
			table.String("status", 1).Default("0")
			table.String("create_by", 64).Default("")
			table.DateTime("create_time").Nullable()
			table.String("update_by", 64).Default("")
			table.DateTime("update_time").Nullable()
			table.String("remark", 500).Nullable()
			table.Index("dict_type")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_role_dept") {
		if err := facades.Schema().Create("sys_role_dept", func(table schema.Blueprint) {
			table.Integer("role_id")
			table.Integer("dept_id")
			table.Index("role_id")
			table.Index("dept_id")
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *CreateSysSystemTables) Down() error {
	for _, table := range []string{"sys_dept", "sys_post", "sys_user_post", "sys_dict_type", "sys_dict_data", "sys_role_dept"} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}

	if facades.Schema().HasColumn("sys_user", "dept_id") {
		if err := facades.Schema().Table("sys_user", func(table schema.Blueprint) {
			table.DropColumn("dept_id")
		}); err != nil {
			return err
		}
	}

	return nil
}
