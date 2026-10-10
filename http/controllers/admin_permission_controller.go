package controllers

import (
	"strconv"
	"encoding/json"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support"
)

func strPtr(s string) *string { return &s }

// AdminPermissionController 权限管理（对齐 owl-admin AdminPermissionController）。
type AdminPermissionController struct{}

func NewAdminPermissionController() *AdminPermissionController { return &AdminPermissionController{} }

// Index GET /system/admin_permissions。
func (r *AdminPermissionController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.PermissionsPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminPermission{})

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

	var permissions []models.AdminPermission
	if err := query.OrderBy("custom_order").Offset((page - 1) * perPage).Limit(perPage).Find(&permissions); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, r.format(&permission))
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

// Tree GET /system/admin_permissions/tree（表单/角色权限树选项）。
func (r *AdminPermissionController) Tree(ctx contractshttp.Context) contractshttp.Response {
	var permissions []models.AdminPermission
	if err := facades.Orm().Query().OrderBy("custom_order").Find(&permissions); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success([]any{}, ""))
	}

	return ctx.Response().Json(contractshttp.StatusOK,
		support.Success(permissionTreeOptions(permissions, 0), ""))
}

func permissionTreeOptions(permissions []models.AdminPermission, parentID int) []map[string]any {
	out := []map[string]any{}
	for _, permission := range permissions {
		if permission.ParentID != parentID {
			continue
		}

		node := map[string]any{"id": permission.ID, "name": permission.Name}
		if children := permissionTreeOptions(permissions, int(permission.ID)); len(children) > 0 {
			node["children"] = children
		}
		out = append(out, node)
	}

	return out
}

func (r *AdminPermissionController) format(permission *models.AdminPermission) map[string]any {
	httpPath := ""
	if permission.HTTPPath != nil {
		httpPath = *permission.HTTPPath
	}

	// http_method 存 JSON 数组文本（对齐 owl-admin），返回时反序列化为数组供 checkboxes 消费
	var methods []any
	if permission.HTTPMethod != nil && *permission.HTTPMethod != "" {
		_ = json.Unmarshal([]byte(*permission.HTTPMethod), &methods)
	}

	return map[string]any{
		"id":           permission.ID,
		"name":         permission.Name,
		"slug":         permission.Slug,
		"http_method":  methods,
		"http_path":    httpPath,
		"custom_order": permission.CustomOrder,
		"parent_id":    permission.ParentID,
		"created_at":   permission.CreatedAt,
		"updated_at":   permission.UpdatedAt,
	}
}

// Store POST /system/admin_permissions。
func (r *AdminPermissionController) Store(ctx contractshttp.Context) contractshttp.Response {
	permission, errResp := r.fill(ctx, &models.AdminPermission{})
	if errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Create(permission); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	if err := r.syncMenus(int(permission.ID), intSlice(ctx.Request().All()["menus"])); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存菜单挂载失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_permissions/{id}/edit。
func (r *AdminPermissionController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var permission models.AdminPermission
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&permission); err != nil || permission.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	data := r.format(&permission)
	data["menus"] = r.menuIDs(int(permission.ID))

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(data, ""))
}

// Update PUT /system/admin_permissions/{id}。
func (r *AdminPermissionController) Update(ctx contractshttp.Context) contractshttp.Response {
	var permission models.AdminPermission
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&permission); err != nil || permission.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	old := permission

	if _, errResp := r.fill(ctx, &permission); errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Save(&permission); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	services.DispatchChangeLog(ctx, "admin_permissions", strconv.FormatUint(uint64(permission.ID), 10),
		services.DiffFields(&old, &permission, []string{"name", "slug", "http_method", "http_path", "custom_order", "parent_id"}))

	if err := r.syncMenus(int(permission.ID), intSlice(ctx.Request().All()["menus"])); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存菜单挂载失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// fill 从请求提取字段（Store/Update 共用）。
func (r *AdminPermissionController) fill(ctx contractshttp.Context, permission *models.AdminPermission) (*models.AdminPermission, func(contractshttp.Context) contractshttp.Response) {
	name := ctx.Request().Input("name")
	slug := ctx.Request().Input("slug")
	if name == "" || slug == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("权限名和标识不能为空"))
		}
	}

	permission.Name = name
	permission.Slug = slug
	permission.ParentID = ctx.Request().InputInt("parent_id", permission.ParentID)
	permission.CustomOrder = ctx.Request().InputInt("custom_order", permission.CustomOrder)

	methods := strSlice(ctx.Request().All()["http_method"])
	if len(methods) > 0 {
		joined, err := json.Marshal(methods)
		if err == nil {
			permission.HTTPMethod = strPtr(string(joined))
		}
	} else {
		permission.HTTPMethod = nil
	}

	if path := ctx.Request().Input("http_path"); path != "" {
		permission.HTTPPath = &path
	}

	return permission, nil
}

// menuIDs 权限挂载的菜单 ID。
func (r *AdminPermissionController) menuIDs(permissionID int) []int {
	var links []models.AdminPermissionMenu
	if err := facades.Orm().Query().Where("permission_id", permissionID).Find(&links); err != nil {
		return []int{}
	}

	ids := make([]int, 0, len(links))
	for _, link := range links {
		ids = append(ids, link.MenuID)
	}

	return ids
}

// syncMenus 全量覆盖权限-菜单关联。
func (r *AdminPermissionController) syncMenus(permissionID int, menuIDs []int) error {
	if _, err := facades.Orm().Query().Where("permission_id", permissionID).Delete(&models.AdminPermissionMenu{}); err != nil {
		return err
	}
	if len(menuIDs) == 0 {
		return nil
	}

	links := make([]models.AdminPermissionMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		links = append(links, models.AdminPermissionMenu{PermissionID: permissionID, MenuID: menuID})
	}

	return facades.Orm().Query().Create(&links)
}

// Destroy DELETE /system/admin_permissions/{id}。
func (r *AdminPermissionController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminPermission{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}
	_, _ = facades.Orm().Query().WhereIn("permission_id", toAny(ids)).Delete(&models.AdminRolePermissions{})
	_, _ = facades.Orm().Query().WhereIn("permission_id", toAny(ids)).Delete(&models.AdminPermissionMenu{})

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}
