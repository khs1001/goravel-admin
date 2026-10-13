package seeders

import (
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// SysSystemSeeder seeds the additional RuoYi system data: a root department,
// default posts, sample dictionaries and the menus/permissions for the new
// system-management pages.
//
// Unlike SysInitSeeder it is not gated on an empty user table; every step is
// idempotent, so it also upgrades an already-seeded database.
type SysSystemSeeder struct{}

// Signature The name and signature of the seeder.
func (s *SysSystemSeeder) Signature() string {
	return "SysSystemSeeder"
}

// Run executes the seeder logic.
func (s *SysSystemSeeder) Run() error {
	if err := s.seedDept(); err != nil {
		return err
	}
	if err := s.seedPosts(); err != nil {
		return err
	}
	if err := s.seedDict(); err != nil {
		return err
	}

	if err := s.seedMenus(); err != nil {
		return err
	}

	return s.seedLogMenus()
}

// seedLogMenus adds the log-management menus (idempotent).
func (s *SysSystemSeeder) seedLogMenus() error {
	existing, err := facades.Orm().Query().Model(&models.SysMenu{}).Where("perms", "monitor:operlog:list").Count()
	if err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}

	create := func(menu models.SysMenu) (uint, error) {
		if err := facades.Orm().Query().Create(&menu); err != nil {
			return 0, err
		}

		return menu.MenuID, nil
	}

	dirID, err := create(models.SysMenu{
		MenuName: "日志管理", ParentID: 0, OrderNum: 3,
		Path: "/log", MenuType: models.MenuTypeDir,
		Icon: "mdi:file-document-outline", Visible: "0", Status: "0",
	})
	if err != nil {
		return err
	}
	_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: 1, MenuID: dirID})

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
			name: "操作日志", path: "/system/operlog", component: "system_oper_log",
			perm: "monitor:operlog:list", icon: "mdi:history", order: 1,
			buttons: []struct{ name, perm string }{
				{"操作查询", "monitor:operlog:query"},
				{"操作删除", "monitor:operlog:remove"},
			},
		},
		{
			name: "登录日志", path: "/system/logininfor", component: "system_login_log",
			perm: "monitor:logininfor:list", icon: "mdi:login-variant", order: 2,
			buttons: []struct{ name, perm string }{
				{"登录查询", "monitor:logininfor:query"},
				{"登录删除", "monitor:logininfor:remove"},
			},
		},
	}

	for _, resource := range resources {
		menuID, err := create(models.SysMenu{
			MenuName: resource.name, ParentID: dirID, OrderNum: resource.order,
			Path: resource.path, Component: resource.component, Perms: resource.perm,
			MenuType: models.MenuTypeMenu, Icon: resource.icon, Visible: "0", Status: "0",
		})
		if err != nil {
			return err
		}
		_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: 1, MenuID: menuID})

		for i, button := range resource.buttons {
			buttonID, err := create(models.SysMenu{
				MenuName: button.name, ParentID: menuID, OrderNum: i + 1,
				Perms: button.perm, MenuType: models.MenuTypeButton, Visible: "0", Status: "0",
			})
			if err != nil {
				return err
			}
			_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: 1, MenuID: buttonID})
		}
	}

	return nil
}

func (s *SysSystemSeeder) seedDept() error {
	var count int64
	count, err := facades.Orm().Query().Model(&models.SysDept{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	root := models.SysDept{
		ParentID:  0,
		Ancestors: "0",
		DeptName:  "总公司",
		OrderNum:  1,
		Status:    "0",
		DelFlag:   "0",
	}
	if err := facades.Orm().Query().Create(&root); err != nil {
		return err
	}

	// Attach the built-in administrator to the root department.
	var admin models.SysUser
	if err := facades.Orm().Query().Where("user_id", 1).First(&admin); err == nil {
		admin.DeptID = root.DeptID
		_ = facades.Orm().Query().Save(&admin)
	}

	return nil
}

func (s *SysSystemSeeder) seedPosts() error {
	var count int64
	count, err := facades.Orm().Query().Model(&models.SysPost{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	posts := []models.SysPost{
		{PostCode: "ceo", PostName: "董事长", PostSort: 1, Status: "0"},
		{PostCode: "se", PostName: "项目经理", PostSort: 2, Status: "0"},
		{PostCode: "hr", PostName: "人力资源", PostSort: 3, Status: "0"},
		{PostCode: "user", PostName: "普通员工", PostSort: 4, Status: "0"},
	}

	for i := range posts {
		if err := facades.Orm().Query().Create(&posts[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *SysSystemSeeder) seedDict() error {
	types := []struct {
		name string
		typ  string
		data []models.SysDictData
	}{
		{
			name: "用户性别", typ: "sys_user_sex",
			data: []models.SysDictData{
				{DictSort: 1, DictLabel: "男", DictValue: "0", DictType: "sys_user_sex", IsDefault: "Y", Status: "0"},
				{DictSort: 2, DictLabel: "女", DictValue: "1", DictType: "sys_user_sex", IsDefault: "N", Status: "0"},
				{DictSort: 3, DictLabel: "未知", DictValue: "2", DictType: "sys_user_sex", IsDefault: "N", Status: "0"},
			},
		},
		{
			name: "系统开关", typ: "sys_normal_disable",
			data: []models.SysDictData{
				{DictSort: 1, DictLabel: "正常", DictValue: "0", DictType: "sys_normal_disable", IsDefault: "Y", Status: "0"},
				{DictSort: 2, DictLabel: "停用", DictValue: "1", DictType: "sys_normal_disable", IsDefault: "N", Status: "0"},
			},
		},
	}

	for _, item := range types {
		var count int64
		count, err := facades.Orm().Query().Model(&models.SysDictType{}).Where("dict_type", item.typ).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		if err := facades.Orm().Query().Create(&models.SysDictType{
			DictName: item.name,
			DictType: item.typ,
			Status:   "0",
		}); err != nil {
			return err
		}
		for i := range item.data {
			if err := facades.Orm().Query().Create(&item.data[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SysSystemSeeder) seedMenus() error {
	// Idempotency marker: the department menu.
	existing, err := facades.Orm().Query().Model(&models.SysMenu{}).Where("perms", "system:dept:list").Count()
	if err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}

	var system models.SysMenu
	if err := facades.Orm().Query().Where("path", "/system").Where("menu_type", models.MenuTypeDir).First(&system); err != nil {
		return nil
	}

	create := func(menu models.SysMenu) (uint, error) {
		if err := facades.Orm().Query().Create(&menu); err != nil {
			return 0, err
		}

		return menu.MenuID, nil
	}

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
			name: "部门管理", path: "/system/dept", component: "system_dept",
			perm: "system:dept:list", icon: "mingcute:department-line", order: 4,
			buttons: []struct{ name, perm string }{
				{"部门查询", "system:dept:query"},
				{"部门新增", "system:dept:add"},
				{"部门修改", "system:dept:edit"},
				{"部门删除", "system:dept:remove"},
			},
		},
		{
			name: "岗位管理", path: "/system/post", component: "system_post",
			perm: "system:post:list", icon: "mdi:badge-account-horizontal-outline", order: 5,
			buttons: []struct{ name, perm string }{
				{"岗位查询", "system:post:query"},
				{"岗位新增", "system:post:add"},
				{"岗位修改", "system:post:edit"},
				{"岗位删除", "system:post:remove"},
			},
		},
		{
			name: "字典管理", path: "/system/dict", component: "system_dict_type",
			perm: "system:dict:list", icon: "mdi:book-open-variant", order: 6,
			buttons: []struct{ name, perm string }{
				{"字典查询", "system:dict:query"},
				{"字典新增", "system:dict:add"},
				{"字典修改", "system:dict:edit"},
				{"字典删除", "system:dict:remove"},
			},
		},
		{
			name: "参数设置", path: "/system/config", component: "system_config",
			perm: "system:config:list", icon: "mdi:cog-outline", order: 7,
			buttons: []struct{ name, perm string }{
				{"参数查询", "system:config:query"},
				{"参数新增", "system:config:add"},
				{"参数修改", "system:config:edit"},
				{"参数删除", "system:config:remove"},
			},
		},
	}

	for _, resource := range resources {
		menuID, err := create(models.SysMenu{
			MenuName: resource.name, ParentID: system.MenuID, OrderNum: resource.order,
			Path: resource.path, Component: resource.component, Perms: resource.perm,
			MenuType: models.MenuTypeMenu, Icon: resource.icon, Visible: "0", Status: "0",
		})
		if err != nil {
			return err
		}

		// Bind the menu to the administrator role.
		_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: 1, MenuID: menuID})

		for i, button := range resource.buttons {
			buttonID, err := create(models.SysMenu{
				MenuName: button.name, ParentID: menuID, OrderNum: i + 1,
				Perms: button.perm, MenuType: models.MenuTypeButton, Visible: "0", Status: "0",
			})
			if err != nil {
				return err
			}
			_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: 1, MenuID: buttonID})
		}
	}

	return nil
}
