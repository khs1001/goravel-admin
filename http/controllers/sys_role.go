package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysRoleController implements the RuoYi `/system/role` action endpoints.
type SysRoleController struct {
	BaseController
}

// NewSysRoleController creates a SysRoleController.
func NewSysRoleController() *SysRoleController {
	return &SysRoleController{}
}

// List returns a paginated role list (`GET /system/role/list`).
func (c *SysRoleController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query().Where("del_flag", "0")
	if keyword := strings.TrimSpace(ctx.Request().Query("roleName")); keyword != "" {
		query = query.Where("role_name like ?", "%"+keyword+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var roles []models.SysRole
	var total int64
	if err := query.OrderBy("role_sort").Paginate(page, perPage, &roles, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": roles, "total": total}))
}

// Get returns a single role with its menu ids (`GET /system/role/{id}`).
func (c *SysRoleController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:query"); resp != nil {
		return resp
	}

	var role models.SysRole
	if err := facades.Orm().Query().Where("role_id", ctx.Request().RouteInt("id")).FirstOrFail(&role); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("角色不存在"))
	}

	role.MenuIDs = c.menuIDs(role.RoleID)
	role.DeptIDs = c.deptIDs(role.RoleID)

	return ctx.Response().Json(http.StatusOK, owl.Success(role))
}

// Add creates a role (`POST /system/role`).
func (c *SysRoleController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:add"); resp != nil {
		return resp
	}

	role := models.SysRole{
		RoleName:  ctx.Request().Input("roleName"),
		RoleKey:   ctx.Request().Input("roleKey"),
		RoleSort:  ctx.Request().InputInt("roleSort"),
		DataScope: defaultValue(ctx.Request().Input("dataScope"), "1"),
		Status:    defaultValue(ctx.Request().Input("status"), "0"),
		DelFlag:   "0",
	}
	if role.RoleName == "" || role.RoleKey == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("角色名称和权限字符不能为空"))
	}

	if err := facades.Orm().Query().Create(&role); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	c.assignMenus(role.RoleID, ctx.Request().InputArray("menuIds"))
	c.assignDepts(role.RoleID, ctx.Request().InputArray("deptIds"))

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a role (`PUT /system/role`).
func (c *SysRoleController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:edit"); resp != nil {
		return resp
	}

	var role models.SysRole
	if err := facades.Orm().Query().Where("role_id", ctx.Request().InputInt("roleId")).FirstOrFail(&role); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("角色不存在"))
	}

	role.RoleName = ctx.Request().Input("roleName")
	role.RoleKey = ctx.Request().Input("roleKey")
	role.RoleSort = ctx.Request().InputInt("roleSort")
	if dataScope := ctx.Request().Input("dataScope"); dataScope != "" {
		role.DataScope = dataScope
	}
	if status := ctx.Request().Input("status"); status != "" {
		role.Status = status
	}

	if err := facades.Orm().Query().Save(&role); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	c.assignMenus(role.RoleID, ctx.Request().InputArray("menuIds"))
	c.assignDepts(role.RoleID, ctx.Request().InputArray("deptIds"))

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a role (`DELETE /system/role/{id}`).
func (c *SysRoleController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")
	if id == 1 {
		return ctx.Response().Json(http.StatusOK, owl.Fail("不允许删除超级管理员角色"))
	}

	if _, err := facades.Orm().Query().Where("role_id", id).Delete(&models.SysRoleMenu{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("role_id", id).Delete(&models.SysRoleDept{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("role_id", id).Delete(&models.SysUserRole{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("role_id", id).Delete(&models.SysRole{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// MenuTree returns the full menu tree plus the ids granted to a role
// (`GET /system/role/menuTree/{id}`).
func (c *SysRoleController) MenuTree(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:query"); resp != nil {
		return resp
	}

	var menus []models.SysMenu
	if err := facades.Orm().Query().OrderBy("order_num").Get(&menus); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	tree := owl.BuildAllTree(menus, 0)
	checked := c.menuIDs(uint(ctx.Request().RouteInt("id")))

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{
		"menus":       tree,
		"checkedKeys": checked,
	}))
}

func (c *SysRoleController) menuIDs(roleID uint) []uint {
	var bindings []models.SysRoleMenu
	if err := facades.Orm().Query().Where("role_id", roleID).Get(&bindings); err != nil {
		return []uint{}
	}

	ids := make([]uint, 0, len(bindings))
	for _, binding := range bindings {
		ids = append(ids, binding.MenuID)
	}

	return ids
}

func (c *SysRoleController) assignMenus(roleID uint, menuIDs []string) {
	if roleID == 0 {
		return
	}

	if _, err := facades.Orm().Query().Where("role_id", roleID).Delete(&models.SysRoleMenu{}); err != nil {
		return
	}

	for _, raw := range menuIDs {
		menuID := parseUint(raw)
		if menuID == 0 {
			continue
		}
		_ = facades.Orm().Query().Create(&models.SysRoleMenu{RoleID: roleID, MenuID: menuID})
	}
}

// DeptTree returns the department tree plus the ids granted to a role, used for
// the "custom" data scope (`GET /system/role/deptTree/{id}`).
func (c *SysRoleController) DeptTree(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:role:query"); resp != nil {
		return resp
	}

	var depts []models.SysDept
	if err := facades.Orm().Query().Where("del_flag", "0").OrderBy("order_num").Get(&depts); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{
		"depts":       owl.BuildDeptTree(depts, 0),
		"checkedKeys": c.deptIDs(uint(ctx.Request().RouteInt("id"))),
	}))
}

func (c *SysRoleController) deptIDs(roleID uint) []uint {
	var bindings []models.SysRoleDept
	if err := facades.Orm().Query().Where("role_id", roleID).Get(&bindings); err != nil {
		return []uint{}
	}

	ids := make([]uint, 0, len(bindings))
	for _, binding := range bindings {
		ids = append(ids, binding.DeptID)
	}

	return ids
}

func (c *SysRoleController) assignDepts(roleID uint, deptIDs []string) {
	if roleID == 0 {
		return
	}

	if _, err := facades.Orm().Query().Where("role_id", roleID).Delete(&models.SysRoleDept{}); err != nil {
		return
	}

	for _, raw := range deptIDs {
		deptID := parseUint(raw)
		if deptID == 0 {
			continue
		}
		_ = facades.Orm().Query().Create(&models.SysRoleDept{RoleID: roleID, DeptID: deptID})
	}
}
