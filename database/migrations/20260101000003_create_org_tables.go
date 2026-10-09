package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000003CreateOrgTables struct{}

// Signature 迁移唯一签名。
func (r *M20260101000003CreateOrgTables) Signature() string {
	return "20260101000003_create_org_tables"
}

// Up 建 4 张表并写入示例种子（幂等：有数据时跳过）。
func (r *M20260101000003CreateOrgTables) Up() error {
	if !facades.Schema().HasTable("admin_departments") {
		if err := facades.Schema().Create("admin_departments", func(table schema.Blueprint) {
			table.ID()
			table.Integer("parent_id").Default(0)
			table.String("name", 50)
			table.Integer("custom_order").Default(0)
			table.String("remark").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Index("parent_id")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_positions") {
		if err := facades.Schema().Create("admin_positions", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 50)
			table.Integer("custom_order").Default(0)
			table.String("remark").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("name")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_dictionaries") {
		if err := facades.Schema().Create("admin_dictionaries", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 50)
			table.String("key", 50)
			table.String("remark").Nullable()
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Unique("key")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("admin_dictionary_items") {
		if err := facades.Schema().Create("admin_dictionary_items", func(table schema.Blueprint) {
			table.ID()
			table.Integer("dictionary_id")
			table.String("label", 50)
			table.String("value", 50)
			table.Integer("custom_order").Default(0)
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.Index("dictionary_id")
		}); err != nil {
			return err
		}
	}

	return r.seed()
}

// seed 示例数据：部门树、岗位、性别字典。
func (r *M20260101000003CreateOrgTables) seed() error {
	// 部门（幂等）
	count, err := facades.Orm().Query().Model(&models.AdminDepartment{}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		root := &models.AdminDepartment{Name: "总公司", CustomOrder: 1}
		if err := facades.Orm().Query().Create(root); err != nil {
			return err
		}

		depts := []models.AdminDepartment{
			{ParentID: int(root.ID), Name: "技术部", CustomOrder: 1},
			{ParentID: int(root.ID), Name: "人事部", CustomOrder: 2},
		}
		if err := facades.Orm().Query().Create(&depts); err != nil {
			return err
		}
	}

	// 岗位（幂等）
	count, err = facades.Orm().Query().Model(&models.AdminPosition{}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		positions := []models.AdminPosition{
			{Name: "负责人", CustomOrder: 1},
			{Name: "专员", CustomOrder: 2},
		}
		if err := facades.Orm().Query().Create(&positions); err != nil {
			return err
		}
	}

	// 字典（幂等）
	count, err = facades.Orm().Query().Model(&models.AdminDictionary{}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		gender := &models.AdminDictionary{Name: "性别", Key: "gender"}
		if err := facades.Orm().Query().Create(gender); err != nil {
			return err
		}

		items := []models.AdminDictionaryItem{
			{DictionaryID: int(gender.ID), Label: "男", Value: "1", CustomOrder: 1},
			{DictionaryID: int(gender.ID), Label: "女", Value: "2", CustomOrder: 2},
			{DictionaryID: int(gender.ID), Label: "未知", Value: "0", CustomOrder: 3},
		}
		if err := facades.Orm().Query().Create(&items); err != nil {
			return err
		}
	}

	return nil
}

// Down 删除 4 张表。
func (r *M20260101000003CreateOrgTables) Down() error {
	for _, table := range []string{"admin_dictionary_items", "admin_dictionaries", "admin_positions", "admin_departments"} {
		if err := facades.Schema().DropIfExists(table); err != nil {
			return err
		}
	}

	return nil
}
