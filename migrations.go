package admin

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/khs1001/goravel-admin/database/migrations"
)

// Migrations 扩展包全部迁移（单一注册源）。
//
// 迁移注册规范：数据迁移只在扩展包内部的 ServiceProvider.Boot 中注册
//（admin 包自身 + 实现 contracts.Migrator 的二级插件），
// 宿主应用不应在 bootstrap/migrations.go 中罗列扩展包/插件迁移。
func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20260101000001CreateAdminTables{},
		&migrations.M20260101000002SeedRbac{},
		&migrations.M20260101000003CreateOrgTables{},
		&migrations.M20260101000004SeedOrgRbac{},
		&migrations.M20260101000005CreateOperationLogs{},
		&migrations.M20260101000007AddChangesToOperationLogs{},
		&migrations.M20260101000008CreateLoginLogs{},
	}
}
