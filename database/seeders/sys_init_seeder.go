package seeders

import (
	"encoding/json"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysInitSeeder seeds the initial RuoYi-style data: the built-in administrator,
// the administrator role, the menu/permission tree and the default settings.
//
// It is idempotent: when `sys_user` already has rows it does nothing.
type SysInitSeeder struct{}

// Signature The name and signature of the seeder.
func (s *SysInitSeeder) Signature() string {
	return "SysInitSeeder"
}

// Run executes the seeder logic.
func (s *SysInitSeeder) Run() error {
	var total int64
	total, err := facades.Orm().Query().Model(&models.SysUser{}).Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// Administrator user.
	password, err := facades.Hash().Make("admin123")
	if err != nil {
		return err
	}
	admin := models.SysUser{
		UserName: "admin",
		NickName: "Administrator",
		UserType: "00",
		Status:   "0",
		DelFlag:  "0",
		Password: password,
	}
	if err := facades.Orm().Query().Create(&admin); err != nil {
		return err
	}

	// Administrator role.
	role := models.SysRole{
		RoleName:  "超级管理员",
		RoleKey:   "admin",
		RoleSort:  1,
		DataScope: "1",
		Status:    "0",
		DelFlag:   "0",
	}
	if err := facades.Orm().Query().Create(&role); err != nil {
		return err
	}
	if err := facades.Orm().Query().Create(&models.SysUserRole{UserID: admin.UserID, RoleID: role.RoleID}); err != nil {
		return err
	}

	menuIDs, err := s.seedMenus()
	if err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		if err := facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: role.RoleID, MenuID: menuID}); err != nil {
			return err
		}
	}

	return s.seedConfig()
}

func (s *SysInitSeeder) seedMenus() ([]uint, error) {
	create := func(menu models.SysMenu) (uint, error) {
		if err := facades.Orm().Query().Create(&menu); err != nil {
			return 0, err
		}

		return menu.MenuID, nil
	}

	var ids []uint

	// Home.
	homeID, err := create(models.SysMenu{
		MenuName: "首页", ParentID: 0, OrderNum: 1,
		Path: "/dashboard", Component: "dashboard", MenuType: models.MenuTypeMenu,
		Icon: "mdi:chart-line", Visible: "0", Status: "0",
	})
	if err != nil {
		return nil, err
	}
	ids = append(ids, homeID)

	// System directory.
	systemID, err := create(models.SysMenu{
		MenuName: "系统管理", ParentID: 0, OrderNum: 2,
		Path: "/system", MenuType: models.MenuTypeDir,
		Icon: "material-symbols:settings-outline", Visible: "0", Status: "0",
	})
	if err != nil {
		return nil, err
	}
	ids = append(ids, systemID)

	resources := []struct {
		name      string
		path      string
		component string
		perm      string
		icon      string
		order     int
		buttons   []struct{ name, perm string }
	}{
		{
			name: "用户管理", path: "/system/user", component: "system_user",
			perm: "system:user:list", icon: "ph:user-gear", order: 1,
			buttons: []struct{ name, perm string }{
				{"用户查询", "system:user:query"},
				{"用户新增", "system:user:add"},
				{"用户修改", "system:user:edit"},
				{"用户删除", "system:user:remove"},
			},
		},
		{
			name: "角色管理", path: "/system/role", component: "system_role",
			perm: "system:role:list", icon: "carbon:user-role", order: 2,
			buttons: []struct{ name, perm string }{
				{"角色查询", "system:role:query"},
				{"角色新增", "system:role:add"},
				{"角色修改", "system:role:edit"},
				{"角色删除", "system:role:remove"},
			},
		},
		{
			name: "菜单管理", path: "/system/menu", component: "system_menu",
			perm: "system:menu:list", icon: "ant-design:menu-unfold-outlined", order: 3,
			buttons: []struct{ name, perm string }{
				{"菜单查询", "system:menu:query"},
				{"菜单新增", "system:menu:add"},
				{"菜单修改", "system:menu:edit"},
				{"菜单删除", "system:menu:remove"},
			},
		},
	}

	for _, resource := range resources {
		menuID, err := create(models.SysMenu{
			MenuName: resource.name, ParentID: systemID, OrderNum: resource.order,
			Path: resource.path, Component: resource.component, Perms: resource.perm,
			MenuType: models.MenuTypeMenu, Icon: resource.icon, Visible: "0", Status: "0",
		})
		if err != nil {
			return nil, err
		}
		ids = append(ids, menuID)

		for i, button := range resource.buttons {
			buttonID, err := create(models.SysMenu{
				MenuName: button.name, ParentID: menuID, OrderNum: i + 1,
				Perms: button.perm, MenuType: models.MenuTypeButton, Visible: "0", Status: "0",
			})
			if err != nil {
				return nil, err
			}
			ids = append(ids, buttonID)
		}
	}

	return ids, nil
}

func (s *SysInitSeeder) seedConfig() error {
	theme, err := json.Marshal(owl.DefaultThemeSetting())
	if err != nil {
		return err
	}

	for _, item := range []models.SysConfig{
		{ConfigName: "系统主题设置", ConfigKey: "system_theme_setting", ConfigValue: string(theme), ConfigType: "N"},
		{ConfigName: "后台语言", ConfigKey: "admin_locale", ConfigValue: "zh_CN", ConfigType: "N"},
	} {
		if err := facades.Orm().Query().Create(&item); err != nil {
			return err
		}
	}

	return nil
}
