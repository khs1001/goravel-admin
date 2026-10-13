package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysLogininforController implements the RuoYi `/system/logininfor` endpoints.
type SysLogininforController struct {
	BaseController
}

// NewSysLogininforController creates a SysLogininforController.
func NewSysLogininforController() *SysLogininforController {
	return &SysLogininforController{}
}

// List returns a paginated login log (`GET /system/logininfor/list`).
func (c *SysLogininforController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:logininfor:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if name := strings.TrimSpace(ctx.Request().Query("userName")); name != "" {
		query = query.Where("user_name like ?", "%"+name+"%")
	}
	if ip := strings.TrimSpace(ctx.Request().Query("ipaddr")); ip != "" {
		query = query.Where("ipaddr like ?", "%"+ip+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var items []models.SysLogininfor
	var total int64
	if err := query.OrderByDesc("info_id").Paginate(page, perPage, &items, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": total}))
}

// Remove deletes a login log (`DELETE /system/logininfor/{id}`).
func (c *SysLogininforController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:logininfor:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("info_id", ctx.Request().RouteInt("id")).Delete(&models.SysLogininfor{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// Clean clears all login logs (`DELETE /system/logininfor/clean`).
func (c *SysLogininforController) Clean(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:logininfor:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("info_id > 0").Delete(&models.SysLogininfor{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "清空成功"))
}
