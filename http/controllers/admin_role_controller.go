package controllers

import (
	"strconv"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support"
)

// AdminRoleController 角色管理（对齐 owl-admin AdminRoleController）。
type AdminRoleController struct {
	permissionService *services.PermissionService
}

func NewAdminRoleController() *AdminRoleController {
	return &AdminRoleController{permissionService: services.NewPermissionService()}
}

// Index GET /system/admin_roles。
func (r *AdminRoleController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.RolesPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminRole{})

	if name := ctx.Request().Query("name", ""); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if slug := ctx.Request().Query("slug", ""); slug != "" {
		query = query.Where("slug LIKE ?", "%"+slug+"%")
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var roles []models.AdminRole
	if err := query.OrderBy("id").Offset((page - 1) * perPage).Limit(perPage).Find(&roles); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(roles))
	for _, role := range roles {
		items = append(items, map[string]any{
			"id":         role.ID,
			"name":       role.Name,
			"slug":       role.Slug,
			"created_at": role.CreatedAt,
			"updated_at": role.UpdatedAt,
		})
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

// Options GET /system/admin_roles/options（用户表单的角色下拉）。
func (r *AdminRoleController) Options(ctx contractshttp.Context) contractshttp.Response {
	var roles []models.AdminRole
	if err := facades.Orm().Query().OrderBy("id").Find(&roles); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success([]any{}, ""))
	}

	options := make([]map[string]any, 0, len(roles))
	for _, role := range roles {
		options = append(options, map[string]any{"id": role.ID, "name": role.Name})
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(options, ""))
}

// Store POST /system/admin_roles。
func (r *AdminRoleController) Store(ctx contractshttp.Context) contractshttp.Response {
	name := ctx.Request().Input("name")
	slug := ctx.Request().Input("slug")
	if name == "" || slug == "" {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("角色名和标识不能为空"))
	}

	var count int64
	count, _ = facades.Orm().Query().Model(&models.AdminRole{}).Where("slug", slug).Count()
	if count > 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("标识已存在"))
	}

	role := &models.AdminRole{Name: name, Slug: slug}
	if err := facades.Orm().Query().Create(role); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_roles/{id}/edit（编辑表单 + 权限抽屉初始数据）。
func (r *AdminRoleController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var role models.AdminRole
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&role); err != nil || role.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	var links []models.AdminRolePermissions
	_ = facades.Orm().Query().Where("role_id", role.ID).Find(&links)
	permissionIDs := make([]int, 0, len(links))
	for _, link := range links {
		permissionIDs = append(permissionIDs, link.PermissionID)
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"id":          role.ID,
		"name":        role.Name,
		"slug":        role.Slug,
		"permissions": permissionIDs,
	}, ""))
}

// Update PUT /system/admin_roles/{id}。
func (r *AdminRoleController) Update(ctx contractshttp.Context) contractshttp.Response {
	var role models.AdminRole
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&role); err != nil || role.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}
	if role.Slug == services.SuperAdministratorSlug {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("超级管理员角色不可修改"))
	}

	old := role

	if name := ctx.Request().Input("name"); name != "" {
		role.Name = name
	}
	if slug := ctx.Request().Input("slug"); slug != "" && slug != role.Slug {
		var count int64
		count, _ = facades.Orm().Query().Model(&models.AdminRole{}).
			Where("slug", slug).Where("id", "!=", role.ID).Count()
		if count > 0 {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("标识已存在"))
		}
		role.Slug = slug
	}

	if err := facades.Orm().Query().Save(&role); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	services.DispatchChangeLog(ctx, "admin_roles", strconv.FormatUint(uint64(role.ID), 10),
		services.DiffFields(&old, &role, []string{"name", "slug"}))

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// SavePermissions POST /system/admin_roles/save_permissions。
func (r *AdminRoleController) SavePermissions(ctx contractshttp.Context) contractshttp.Response {
	roleID := ctx.Request().InputInt("id")
	if roleID <= 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("缺少角色 ID"))
	}

	var role models.AdminRole
	if err := facades.Orm().Query().Where("id", roleID).First(&role); err != nil || role.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	permissionIDs := intSlice(ctx.Request().All()["permissions"])
	if err := r.permissionService.SavePermissions(roleID, permissionIDs); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存权限失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("保存成功"))
}

// Destroy DELETE /system/admin_roles/{id}。
func (r *AdminRoleController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	for _, id := range ids {
		var role models.AdminRole
		if err := facades.Orm().Query().Where("id", id).First(&role); err == nil && role.Slug == services.SuperAdministratorSlug {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("超级管理员角色不可删除"))
		}
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminRole{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}
	_, _ = facades.Orm().Query().WhereIn("role_id", toAny(ids)).Delete(&models.AdminRoleUsers{})
	_, _ = facades.Orm().Query().WhereIn("role_id", toAny(ids)).Delete(&models.AdminRolePermissions{})

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}
