package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/support"
)

// AdminDepartmentController 部门管理（树形组织架构）。
type AdminDepartmentController struct{}

func NewAdminDepartmentController() *AdminDepartmentController { return &AdminDepartmentController{} }

// Index GET /system/admin_departments（无 _action 返回页面 schema；_action=getData 返回树形数据）。
func (r *AdminDepartmentController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.DepartmentsPage(), ""))
	}

	var depts []models.AdminDepartment
	if err := facades.Orm().Query().OrderBy("custom_order").OrderBy("id").Find(&depts); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	tree := departmentTree(depts, 0)

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(tree, int64(len(tree))))
}

func departmentTree(depts []models.AdminDepartment, parentID int) []map[string]any {
	out := []map[string]any{}
	for _, dept := range depts {
		if dept.ParentID != parentID {
			continue
		}

		node := map[string]any{
			"id":           dept.ID,
			"parent_id":    dept.ParentID,
			"name":         dept.Name,
			"custom_order": dept.CustomOrder,
			"remark":       dept.Remark,
			"created_at":   dept.CreatedAt,
			"updated_at":   dept.UpdatedAt,
		}
		if children := departmentTree(depts, int(dept.ID)); len(children) > 0 {
			node["children"] = children
		}
		out = append(out, node)
	}

	return out
}

// TreeOptions GET /system/admin_departments/tree（表单 TreeSelect 选项）。
func (r *AdminDepartmentController) TreeOptions(ctx contractshttp.Context) contractshttp.Response {
	var depts []models.AdminDepartment
	if err := facades.Orm().Query().OrderBy("custom_order").OrderBy("id").Find(&depts); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success([]any{}, ""))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(departmentTreeOptions(depts, 0), ""))
}

func departmentTreeOptions(depts []models.AdminDepartment, parentID int) []map[string]any {
	out := []map[string]any{}
	for _, dept := range depts {
		if dept.ParentID != parentID {
			continue
		}

		node := map[string]any{"id": dept.ID, "name": dept.Name}
		if children := departmentTreeOptions(depts, int(dept.ID)); len(children) > 0 {
			node["children"] = children
		}
		out = append(out, node)
	}

	return out
}

// Store POST /system/admin_departments。
func (r *AdminDepartmentController) Store(ctx contractshttp.Context) contractshttp.Response {
	dept, errResp := r.fill(ctx, &models.AdminDepartment{})
	if errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Create(dept); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_departments/{id}/edit。
func (r *AdminDepartmentController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var dept models.AdminDepartment
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&dept); err != nil || dept.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"id":           dept.ID,
		"parent_id":    dept.ParentID,
		"name":         dept.Name,
		"custom_order": dept.CustomOrder,
		"remark":       dept.Remark,
	}, ""))
}

// Update PUT /system/admin_departments/{id}。
func (r *AdminDepartmentController) Update(ctx contractshttp.Context) contractshttp.Response {
	var dept models.AdminDepartment
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&dept); err != nil || dept.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	id := int(dept.ID)
	if _, errResp := r.fill(ctx, &dept); errResp != nil {
		return errResp(ctx)
	}
	// 不允许把父级设为自己或自己的后代
	if dept.ParentID == id || isDescendant(id, dept.ParentID) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("父级不能选择自己或自己的子部门"))
	}

	if err := facades.Orm().Query().Save(&dept); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// isDescendant 检查 ancestor 的后代中是否包含 target（用于防止树成环）。
func isDescendant(target, ancestor int) bool {
	if target <= 0 || ancestor <= 0 {
		return false
	}

	var depts []models.AdminDepartment
	if err := facades.Orm().Query().Find(&depts); err != nil {
		return false
	}

	parentOf := map[int]int{}
	for _, dept := range depts {
		parentOf[int(dept.ID)] = dept.ParentID
	}

	current := ancestor
	visited := map[int]bool{}
	for current > 0 && !visited[current] {
		if current == target {
			return true
		}
		visited[current] = true
		current = parentOf[current]
	}

	return false
}

// fill 从请求提取字段（Store/Update 共用）。
func (r *AdminDepartmentController) fill(ctx contractshttp.Context, dept *models.AdminDepartment) (*models.AdminDepartment, func(contractshttp.Context) contractshttp.Response) {
	name := ctx.Request().Input("name")
	if name == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("部门名称不能为空"))
		}
	}

	dept.Name = name
	dept.ParentID = ctx.Request().InputInt("parent_id", dept.ParentID)
	dept.CustomOrder = ctx.Request().InputInt("custom_order", dept.CustomOrder)

	remark := ctx.Request().Input("remark")
	dept.Remark = &remark

	return dept, nil
}

// Destroy DELETE /system/admin_departments/{id}（有子部门时拒绝）。
func (r *AdminDepartmentController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	for _, id := range ids {
		var count int64
		count, _ = facades.Orm().Query().Model(&models.AdminDepartment{}).Where("parent_id", id).Count()
		if count > 0 {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("存在子部门，请先删除子部门"))
		}
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminDepartment{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}
