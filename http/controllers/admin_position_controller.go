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

// AdminPositionController 岗位管理。
type AdminPositionController struct{}

func NewAdminPositionController() *AdminPositionController { return &AdminPositionController{} }

// Index GET /system/admin_positions（无 _action 返回页面 schema；_action=getData 返回列表数据）。
func (r *AdminPositionController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.PositionsPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminPosition{})
	if name := ctx.Request().Query("name", ""); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var positions []models.AdminPosition
	if err := query.OrderBy("custom_order").OrderBy("id").
		Offset((page - 1) * perPage).Limit(perPage).Find(&positions); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(positions))
	for _, position := range positions {
		items = append(items, map[string]any{
			"id":           position.ID,
			"name":         position.Name,
			"custom_order": position.CustomOrder,
			"remark":       position.Remark,
			"created_at":   position.CreatedAt,
			"updated_at":   position.UpdatedAt,
		})
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

// Options GET /system/admin_positions/options（下拉选项）。
func (r *AdminPositionController) Options(ctx contractshttp.Context) contractshttp.Response {
	var positions []models.AdminPosition
	if err := facades.Orm().Query().OrderBy("custom_order").OrderBy("id").Find(&positions); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success([]any{}, ""))
	}

	out := make([]map[string]any, 0, len(positions))
	for _, position := range positions {
		out = append(out, map[string]any{"id": position.ID, "name": position.Name})
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(out, ""))
}

// Store POST /system/admin_positions。
func (r *AdminPositionController) Store(ctx contractshttp.Context) contractshttp.Response {
	position, errResp := r.fill(ctx, &models.AdminPosition{})
	if errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Create(position); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_positions/{id}/edit。
func (r *AdminPositionController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var position models.AdminPosition
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&position); err != nil || position.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"id":           position.ID,
		"name":         position.Name,
		"custom_order": position.CustomOrder,
		"remark":       position.Remark,
	}, ""))
}

// Update PUT /system/admin_positions/{id}。
func (r *AdminPositionController) Update(ctx contractshttp.Context) contractshttp.Response {
	var position models.AdminPosition
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&position); err != nil || position.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	old := position

	if _, errResp := r.fill(ctx, &position); errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Save(&position); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	services.DispatchChangeLog(ctx, "admin_positions", strconv.FormatUint(uint64(position.ID), 10),
		services.DiffFields(&old, &position, []string{"name", "custom_order", "remark"}))

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// fill 从请求提取字段（Store/Update 共用）。
func (r *AdminPositionController) fill(ctx contractshttp.Context, position *models.AdminPosition) (*models.AdminPosition, func(contractshttp.Context) contractshttp.Response) {
	name := ctx.Request().Input("name")
	if name == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("岗位名称不能为空"))
		}
	}

	position.Name = name
	position.CustomOrder = ctx.Request().InputInt("custom_order", position.CustomOrder)

	remark := ctx.Request().Input("remark")
	position.Remark = &remark

	return position, nil
}

// Destroy DELETE /system/admin_positions/{id}（支持 ?ids= 批量）。
func (r *AdminPositionController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminPosition{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}
