package controllers

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysDeptController implements the RuoYi `/system/dept` action endpoints. The
// department tree is also the basis of the data-scope feature.
type SysDeptController struct {
	BaseController
}

// NewSysDeptController creates a SysDeptController.
func NewSysDeptController() *SysDeptController {
	return &SysDeptController{}
}

// List returns the department tree (`GET /system/dept/list`).
func (c *SysDeptController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dept:list"); resp != nil {
		return resp
	}

	depts, err := c.visibleDepts(ctx)
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	if keyword := strings.TrimSpace(ctx.Request().Query("deptName")); keyword != "" {
		filtered := make([]models.SysDept, 0, len(depts))
		for _, dept := range depts {
			if strings.Contains(dept.DeptName, keyword) {
				filtered = append(filtered, dept)
			}
		}
		depts = filtered
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(owl.BuildDeptTreeFiltered(depts)))
}
func (c *SysDeptController) Tree(ctx http.Context) http.Response {
	var depts []models.SysDept
	if err := facades.Orm().Query().Where("del_flag", "0").OrderBy("order_num").Get(&depts); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(owl.BuildDeptTree(depts, 0)))
}

// Get returns a single department (`GET /system/dept/{id}`).
func (c *SysDeptController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dept:query"); resp != nil {
		return resp
	}

	var dept models.SysDept
	if err := facades.Orm().Query().Where("dept_id", ctx.Request().RouteInt("id")).FirstOrFail(&dept); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("部门不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(dept))
}

// Add creates a department (`POST /system/dept`).
func (c *SysDeptController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dept:add"); resp != nil {
		return resp
	}

	parentID := uint(ctx.Request().InputInt("parentId"))
	dept := models.SysDept{
		ParentID:  parentID,
		Ancestors: c.ancestors(parentID),
		DeptName:  ctx.Request().Input("deptName"),
		OrderNum:  ctx.Request().InputInt("orderNum"),
		Leader:    ctx.Request().Input("leader"),
		Phone:     ctx.Request().Input("phone"),
		Email:     ctx.Request().Input("email"),
		Status:    defaultValue(ctx.Request().Input("status"), "0"),
		DelFlag:   "0",
	}
	if dept.DeptName == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("部门名称不能为空"))
	}

	if err := facades.Orm().Query().Create(&dept); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a department (`PUT /system/dept`).
func (c *SysDeptController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dept:edit"); resp != nil {
		return resp
	}

	var dept models.SysDept
	if err := facades.Orm().Query().Where("dept_id", ctx.Request().InputInt("deptId")).FirstOrFail(&dept); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("部门不存在"))
	}

	parentID := uint(ctx.Request().InputInt("parentId"))
	dept.ParentID = parentID
	dept.Ancestors = c.ancestors(parentID)
	dept.DeptName = ctx.Request().Input("deptName")
	dept.OrderNum = ctx.Request().InputInt("orderNum")
	dept.Leader = ctx.Request().Input("leader")
	dept.Phone = ctx.Request().Input("phone")
	dept.Email = ctx.Request().Input("email")
	if status := ctx.Request().Input("status"); status != "" {
		dept.Status = status
	}

	if err := facades.Orm().Query().Save(&dept); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a department (`DELETE /system/dept/{id}`).
func (c *SysDeptController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dept:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")

	var children []models.SysDept
	if err := facades.Orm().Query().Where("parent_id", id).Where("del_flag", "0").Get(&children); err == nil && len(children) > 0 {
		return ctx.Response().Json(http.StatusOK, owl.Fail("存在下级部门，不允许删除"))
	}

	var users []models.SysUser
	if err := facades.Orm().Query().Where("dept_id", id).Where("del_flag", "0").Get(&users); err == nil && len(users) > 0 {
		return ctx.Response().Json(http.StatusOK, owl.Fail("部门存在用户，不允许删除"))
	}

	if _, err := facades.Orm().Query().Where("dept_id", id).Delete(&models.SysDept{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// visibleDepts returns the departments the current user may see, honouring the
// data scope.
func (c *SysDeptController) visibleDepts(ctx http.Context) ([]models.SysDept, error) {
	query := facades.Orm().Query().Where("del_flag", "0")

	scope := services.NewDataScopeService().Resolve(current.User(ctx))
	if !scope.Unrestricted {
		if len(scope.DeptIDs) == 0 {
			return []models.SysDept{}, nil
		}
		ids := make([]any, len(scope.DeptIDs))
		for i, id := range scope.DeptIDs {
			ids[i] = id
		}
		query = query.WhereIn("dept_id", ids)
	}

	var depts []models.SysDept
	if err := query.OrderBy("order_num").Get(&depts); err != nil {
		return nil, err
	}

	return depts, nil
}

func (c *SysDeptController) ancestors(parentID uint) string {
	if parentID == 0 {
		return "0"
	}

	var parent models.SysDept
	if err := facades.Orm().Query().Where("dept_id", parentID).FirstOrFail(&parent); err != nil {
		return "0"
	}

	if parent.Ancestors == "" {
		return strconv.FormatUint(uint64(parentID), 10)
	}

	return parent.Ancestors + "," + strconv.FormatUint(uint64(parentID), 10)
}
