package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/support"
)

// AdminOperationLogController 操作日志（只读）。
type AdminOperationLogController struct{}

func NewAdminOperationLogController() *AdminOperationLogController {
	return &AdminOperationLogController{}
}

// Index GET /system/admin_operation_logs（无 _action 返回页面 schema；_action=getData 返回列表数据）。
func (r *AdminOperationLogController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.OperationLogsPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminOperationLog{})
	if username := ctx.Request().Query("username", ""); username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if path := ctx.Request().Query("path", ""); path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}
	if action := ctx.Request().Query("action", ""); action != "" {
		query = query.Where("action", action)
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var logs []models.AdminOperationLog
	if err := query.OrderBy("id", "desc").
		Offset((page - 1) * perPage).Limit(perPage).Find(&logs); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(logs))
	for _, item := range logs {
		items = append(items, r.format(&item))
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

func (r *AdminOperationLogController) format(log *models.AdminOperationLog) map[string]any {
	input := ""
	if log.Input != nil {
		input = *log.Input
	}

	return map[string]any{
		"id":          log.ID,
		"user_id":     log.UserID,
		"username":    log.Username,
		"action":      log.Action,
		"resource":    log.Resource,
		"resource_id": log.ResourceID,
		"method":      log.Method,
		"path":        log.Path,
		"ip":          log.IP,
		"input":       input,
		"created_at":  log.CreatedAt,
	}
}

// Edit GET /system/admin_operation_logs/{id}/edit（详情抽屉回显）。
func (r *AdminOperationLogController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var log models.AdminOperationLog
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&log); err != nil || log.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(r.format(&log), ""))
}
