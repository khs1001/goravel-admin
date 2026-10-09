package migrations

import (
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

type M20260101000002SeedRbac struct{}

// Signature 迁移唯一签名。
func (r *M20260101000002SeedRbac) Signature() string {
	return "20260101000002_seed_rbac"
}

// Up 写入系统管理菜单、权限集与权限-菜单挂载（幂等：有权限数据时跳过）。
func (r *M20260101000002SeedRbac) Up() error {
	count, err := facades.Orm().Query().Model(&models.AdminPermission{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 系统管理菜单树
	systemMenu := &models.AdminMenu{
		ParentID: 0, CustomOrder: 200, Title: "系统管理",
		URL: strPtr("/system"), UrlType: 1, Visible: 1,
	}
	if err := facades.Orm().Query().Create(systemMenu); err != nil {
		return err
	}

	subMenus := []models.AdminMenu{
		{ParentID: int(systemMenu.ID), CustomOrder: 1, Title: "管理员", URL: strPtr("/system/admin_users"), UrlType: 1, Visible: 1},
		{ParentID: int(systemMenu.ID), CustomOrder: 2, Title: "角色", URL: strPtr("/system/admin_roles"), UrlType: 1, Visible: 1},
		{ParentID: int(systemMenu.ID), CustomOrder: 3, Title: "权限", URL: strPtr("/system/admin_permissions"), UrlType: 1, Visible: 1},
		{ParentID: int(systemMenu.ID), CustomOrder: 4, Title: "菜单", URL: strPtr("/system/admin_menus"), UrlType: 1, Visible: 1},
	}
	if err := facades.Orm().Query().Create(&subMenus); err != nil {
		return err
	}

	// 权限集：每张资源表的查看/写入权限 + 挂载对应菜单
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
		menuIdx    []int // 挂载的子菜单索引（-1 = 系统管理父菜单）
	}

	allMenuIdx := []int{0, 1, 2, 3}
	defs := []permDef{
		{"系统管理", "system-manage", nil, []int{-1}},
		{"管理员查看", "system.admin_users.index", httpPath("GET /admin-api/system/admin_users*"), allMenuIdx},
		{"管理员写入", "system.admin_users.store", httpPath("POST /admin-api/system/admin_users", "PUT /admin-api/system/admin_users/*", "DELETE /admin-api/system/admin_users*"), []int{0}},
		{"角色查看", "system.admin_roles.index", httpPath("GET /admin-api/system/admin_roles*"), allMenuIdx},
		{"角色写入", "system.admin_roles.store", httpPath("POST /admin-api/system/admin_roles*", "PUT /admin-api/system/admin_roles/*", "DELETE /admin-api/system/admin_roles*"), []int{1}},
		{"权限查看", "system.admin_permissions.index", httpPath("GET /admin-api/system/admin_permissions*"), allMenuIdx},
		{"权限写入", "system.admin_permissions.store", httpPath("POST /admin-api/system/admin_permissions*", "PUT /admin-api/system/admin_permissions/*", "DELETE /admin-api/system/admin_permissions*"), []int{2}},
		{"菜单查看", "system.admin_menus.index", httpPath("GET /admin-api/system/admin_menus*"), allMenuIdx},
		{"菜单写入", "system.admin_menus.store", httpPath("POST /admin-api/system/admin_menus*", "PUT /admin-api/system/admin_menus/*", "DELETE /admin-api/system/admin_menus*"), []int{3}},
	}

	menuID := func(idx int) int {
		if idx == -1 {
			return int(systemMenu.ID)
		}

		return int(subMenus[idx].ID)
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
				MenuID:       menuID(idx),
			})
		}
		if err := facades.Orm().Query().Create(&links); err != nil {
			return err
		}
	}

	return nil
}

// Down 删除种子写入的权限与挂载（菜单保留）。
func (r *M20260101000002SeedRbac) Down() error {
	_, err := facades.Orm().Query().Where("slug LIKE ?", "system%").Delete(&models.AdminPermission{})

	return err
}
