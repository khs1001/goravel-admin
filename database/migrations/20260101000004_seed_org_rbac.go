package migrations

import (
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000004SeedOrgRbac struct{}

// Signature 迁移唯一签名。
func (r *M20260101000004SeedOrgRbac) Signature() string {
	return "20260101000004_seed_org_rbac"
}

// Up 在"系统管理"菜单下追加 部门/岗位/字典 菜单，并写入对应权限与挂载（幂等：有权限时跳过）。
func (r *M20260101000004SeedOrgRbac) Up() error {
	count, err := facades.Orm().Query().Model(&models.AdminPermission{}).Where("slug LIKE ?", "system.admin_departments%").Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 找到 M2 种子创建的"系统管理"父菜单
	var systemMenu models.AdminMenu
	if err := facades.Orm().Query().
		Where("parent_id", 0).Where("title", "系统管理").First(&systemMenu); err != nil || systemMenu.ID == 0 {
		// 主项目未初始化后台菜单时直接跳过（表迁移已保证表存在）
		return nil
	}

	subMenus := []models.AdminMenu{
		{ParentID: int(systemMenu.ID), CustomOrder: 5, Title: "部门", URL: strPtr("/system/admin_departments"), UrlType: 1, Visible: 1},
		{ParentID: int(systemMenu.ID), CustomOrder: 6, Title: "岗位", URL: strPtr("/system/admin_positions"), UrlType: 1, Visible: 1},
		{ParentID: int(systemMenu.ID), CustomOrder: 7, Title: "字典", URL: strPtr("/system/admin_dictionaries"), UrlType: 1, Visible: 1},
	}
	if err := facades.Orm().Query().Create(&subMenus); err != nil {
		return err
	}

	httpPath := func(paths ...string) *string {
		joined := ""
		for i, p := range paths {
			if i > 0 {
				joined += "\n"
			}
			joined += p
		}

		return &joined
	}

	type permDef struct {
		name, slug string
		path       *string
		menuIdx    []int // 挂载的子菜单索引
	}

	allMenuIdx := []int{0, 1, 2}
	defs := []permDef{
		{"部门查看", "system.admin_departments.index", httpPath("GET /admin-api/system/admin_departments*"), allMenuIdx},
		{"部门写入", "system.admin_departments.store", httpPath("POST /admin-api/system/admin_departments", "PUT /admin-api/system/admin_departments/*", "DELETE /admin-api/system/admin_departments*"), []int{0}},
		{"岗位查看", "system.admin_positions.index", httpPath("GET /admin-api/system/admin_positions*"), allMenuIdx},
		{"岗位写入", "system.admin_positions.store", httpPath("POST /admin-api/system/admin_positions", "PUT /admin-api/system/admin_positions/*", "DELETE /admin-api/system/admin_positions*"), []int{1}},
		{"字典查看", "system.admin_dictionaries.index", httpPath("GET /admin-api/system/admin_dictionaries*", "GET /admin-api/system/admin_dictionary_items*"), allMenuIdx},
		{"字典写入", "system.admin_dictionaries.store", httpPath("POST /admin-api/system/admin_dictionaries", "PUT /admin-api/system/admin_dictionaries/*", "DELETE /admin-api/system/admin_dictionaries*",
			"POST /admin-api/system/admin_dictionary_items", "PUT /admin-api/system/admin_dictionary_items/*", "DELETE /admin-api/system/admin_dictionary_items*"), []int{2}},
	}

	for _, def := range defs {
		permission := &models.AdminPermission{
			Name:     def.name,
			Slug:     def.slug,
			HTTPPath: def.path,
		}
		if err := facades.Orm().Query().Create(permission); err != nil {
			return err
		}

		links := make([]models.AdminPermissionMenu, 0, len(def.menuIdx))
		for _, idx := range def.menuIdx {
			links = append(links, models.AdminPermissionMenu{
				PermissionID: int(permission.ID),
				MenuID:       int(subMenus[idx].ID),
			})
		}
		if err := facades.Orm().Query().Create(&links); err != nil {
			return err
		}
	}

	return nil
}

// Down 删除种子写入的权限与挂载（菜单保留）。
func (r *M20260101000004SeedOrgRbac) Down() error {
	_, err := facades.Orm().Query().Where("slug LIKE ?", "system.admin_departments%").
		OrWhere("slug LIKE ?", "system.admin_positions%").
		OrWhere("slug LIKE ?", "system.admin_dictionaries%").
		Delete(&models.AdminPermission{})

	return err
}
