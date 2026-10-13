package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysConfigController implements the RuoYi `/system/config` (parameter settings)
// action endpoints. The same table backs the Owl `/_settings` persistence.
type SysConfigController struct {
	BaseController
}

// NewSysConfigController creates a SysConfigController.
func NewSysConfigController() *SysConfigController {
	return &SysConfigController{}
}

// List returns a paginated parameter list (`GET /system/config/list`).
func (c *SysConfigController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:config:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if name := strings.TrimSpace(ctx.Request().Query("configName")); name != "" {
		query = query.Where("config_name like ?", "%"+name+"%")
	}
	if key := strings.TrimSpace(ctx.Request().Query("configKey")); key != "" {
		query = query.Where("config_key like ?", "%"+key+"%")
	}

	var items []models.SysConfig
	var total int64
	if err := query.OrderBy("config_id").Paginate(page, perPage, &items, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": total}))
}

// Get returns a single parameter (`GET /system/config/{id}`).
func (c *SysConfigController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:config:query"); resp != nil {
		return resp
	}

	var item models.SysConfig
	if err := facades.Orm().Query().Where("config_id", ctx.Request().RouteInt("id")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("参数不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(item))
}

// Add creates a parameter (`POST /system/config`).
func (c *SysConfigController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:config:add"); resp != nil {
		return resp
	}

	key := strings.TrimSpace(ctx.Request().Input("configKey"))
	if key == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("参数键名不能为空"))
	}

	item := models.SysConfig{
		ConfigName:  ctx.Request().Input("configName"),
		ConfigKey:   key,
		ConfigValue: ctx.Request().Input("configValue"),
		ConfigType:  defaultValue(ctx.Request().Input("configType"), "N"),
		BaseModel: models.BaseModel{
			Remark: ctx.Request().Input("remark"),
		},
	}
	if err := facades.Orm().Query().Create(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a parameter (`PUT /system/config`).
func (c *SysConfigController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:config:edit"); resp != nil {
		return resp
	}

	var item models.SysConfig
	if err := facades.Orm().Query().Where("config_id", ctx.Request().InputInt("configId")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("参数不存在"))
	}

	item.ConfigName = ctx.Request().Input("configName")
	item.ConfigKey = ctx.Request().Input("configKey")
	item.ConfigValue = ctx.Request().Input("configValue")
	if configType := ctx.Request().Input("configType"); configType != "" {
		item.ConfigType = configType
	}
	item.Remark = ctx.Request().Input("remark")

	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a parameter (`DELETE /system/config/{id}`).
func (c *SysConfigController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:config:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("config_id", ctx.Request().RouteInt("id")).Delete(&models.SysConfig{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}
