package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20260101000007AddChangesToOperationLogs struct{}

// Signature 迁移唯一签名。
func (r *M20260101000007AddChangesToOperationLogs) Signature() string {
	return "20260101000007_add_changes_to_operation_logs"
}

// Up 为操作日志表增加 changes 字段（JSON：字段级变更 {col: {old, new}}，仅 updated 动作记录）。
func (r *M20260101000007AddChangesToOperationLogs) Up() error {
	if facades.Schema().HasTable("admin_operation_logs") && !facades.Schema().HasColumn("admin_operation_logs", "changes") {
		return facades.Schema().Table("admin_operation_logs", func(table schema.Blueprint) {
			table.LongText("changes").Nullable()
		})
	}

	return nil
}

// Down 回滚：删除 changes 字段。
func (r *M20260101000007AddChangesToOperationLogs) Down() error {
	if facades.Schema().HasColumn("admin_operation_logs", "changes") {
		return facades.Schema().Table("admin_operation_logs", func(table schema.Blueprint) {
			table.DropColumn("changes")
		})
	}

	return nil
}
