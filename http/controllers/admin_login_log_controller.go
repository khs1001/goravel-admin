package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/support"
)

// AdminLoginLogController 登录日志（只读）。
type AdminLoginLogController struct{}

func NewAdminLoginLogController() *AdminLoginLogController {
	return &AdminLoginLogController{}
}

// Index GET /system/admin_login_logs（无 _action 返回页面 schema；_action=getData 返回列表数据）。
func (r *AdminLoginLogController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.LoginLogsPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminLoginLog{})
	if username := ctx.Request().Query("username", ""); username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if ip := ctx.Request().Query("ip", ""); ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}
	if status := ctx.Request().Query("status", ""); status != "" {
		query = query.Where("status", status)
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var logs []models.AdminLoginLog
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

func (r *AdminLoginLogController) format(log *models.AdminLoginLog) map[string]any {
	return map[string]any{
		"id":         log.ID,
		"user_id":    log.UserID,
		"username":   log.Username,
		"status":     log.Status,
		"message":    log.Message,
		"ip":         log.IP,
		"user_agent": log.UserAgent,
		"created_at": log.CreatedAt,
	}
}

// Edit GET /system/admin_login_logs/{id}/edit（详情抽屉回显）。
func (r *AdminLoginLogController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var log models.AdminLoginLog
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&log); err != nil || log.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(r.format(&log), ""))
}
