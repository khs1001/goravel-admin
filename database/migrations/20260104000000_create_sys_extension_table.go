package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

// CreateSysExtensionTable creates the plugin registry table.
type CreateSysExtensionTable struct{}

// Signature The unique signature for the migration.
func (r *CreateSysExtensionTable) Signature() string {
	return "20260104000000_create_sys_extension_table"
}

// Up Run the migrations.
func (r *CreateSysExtensionTable) Up() error {
	if facades.Schema().HasTable("sys_extension") {
		return nil
	}

	return facades.Schema().Create("sys_extension", func(table schema.Blueprint) {
		table.Increments("extension_id")
		table.String("name", 100).Default("")
		table.Integer("is_enabled").Default(0)
		table.String("version", 50).Default("")
		table.Text("options").Nullable()
		table.String("create_by", 64).Default("")
		table.DateTime("create_time").Nullable()
		table.String("update_by", 64).Default("")
		table.DateTime("update_time").Nullable()
		table.Unique("name")
	})
}

// Down Reverse the migrations.
func (r *CreateSysExtensionTable) Down() error {
	return facades.Schema().DropIfExists("sys_extension")
}
