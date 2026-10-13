package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

// CreateSysLogTables adds the RuoYi operation/login log tables.
type CreateSysLogTables struct{}

// Signature The unique signature for the migration.
func (r *CreateSysLogTables) Signature() string {
	return "20260103000000_create_sys_log_tables"
}

// Up Run the migrations.
func (r *CreateSysLogTables) Up() error {
	if !facades.Schema().HasTable("sys_oper_log") {
		if err := facades.Schema().Create("sys_oper_log", func(table schema.Blueprint) {
			table.Increments("oper_id")
			table.String("title", 50).Default("")
			table.Integer("business_type").Default(0)
			table.String("method", 200).Default("")
			table.String("request_method", 10).Default("")
			table.Integer("operator_type").Default(0)
			table.String("oper_name", 50).Default("")
			table.String("oper_url", 255).Default("")
			table.String("oper_ip", 128).Default("")
			table.Text("oper_param").Nullable()
			table.Text("json_result").Nullable()
			table.Integer("status").Default(0)
			table.String("error_msg", 2000).Nullable()
			table.DateTime("oper_time").Nullable()
			table.Integer("cost_time").Default(0)
			table.Index("oper_time")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("sys_logininfor") {
		if err := facades.Schema().Create("sys_logininfor", func(table schema.Blueprint) {
			table.Increments("info_id")
			table.String("user_name", 50).Default("")
			table.String("ipaddr", 128).Default("")
			table.String("browser", 50).Default("")
			table.String("os", 50).Default("")
			table.Integer("status").Default(0)
			table.String("msg", 255).Default("")
			table.DateTime("login_time").Nullable()
			table.Index("login_time")
		}); err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *CreateSysLogTables) Down() error {
	for _, table := range []string{"sys_oper_log", "sys_logininfor"} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}

	return nil
}
