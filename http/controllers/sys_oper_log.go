package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysOperLogController implements the RuoYi `/system/operlog` endpoints.
type SysOperLogController struct {
	BaseController
}

// NewSysOperLogController creates a SysOperLogController.
func NewSysOperLogController() *SysOperLogController {
	return &SysOperLogController{}
}

// List returns a paginated operation log (`GET /system/operlog/list`).
func (c *SysOperLogController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:operlog:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if title := strings.TrimSpace(ctx.Request().Query("title")); title != "" {
		query = query.Where("title like ?", "%"+title+"%")
	}
	if name := strings.TrimSpace(ctx.Request().Query("operName")); name != "" {
		query = query.Where("oper_name like ?", "%"+name+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var items []models.SysOperLog
	var total int64
	if err := query.OrderByDesc("oper_id").Paginate(page, perPage, &items, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": total}))
}

// Get returns a single operation log (`GET /system/operlog/{id}`).
func (c *SysOperLogController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:operlog:query"); resp != nil {
		return resp
	}

	var item models.SysOperLog
	if err := facades.Orm().Query().Where("oper_id", ctx.Request().RouteInt("id")).FirstOrFail(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("日志不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(item))
}

// Remove deletes an operation log (`DELETE /system/operlog/{id}`).
func (c *SysOperLogController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:operlog:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("oper_id", ctx.Request().RouteInt("id")).Delete(&models.SysOperLog{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// Clean clears all operation logs (`DELETE /system/operlog/clean`).
func (c *SysOperLogController) Clean(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "monitor:operlog:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("oper_id > 0").Delete(&models.SysOperLog{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "清空成功"))
}
