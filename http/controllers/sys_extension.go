package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"

	"github.com/khs1001/goravel-admin/support/owl"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// SysExtensionController implements the plugin extension-management endpoints.
type SysExtensionController struct {
	BaseController
}

// NewSysExtensionController creates a SysExtensionController.
func NewSysExtensionController() *SysExtensionController {
	return &SysExtensionController{}
}

// List returns the discovered plugins with their enable state
// (`GET /system/extension/list`).
func (c *SysExtensionController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:extension:list"); resp != nil {
		return resp
	}

	items := plugin.R().List()

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": len(items)}))
}

// Toggle enables or disables a plugin (`POST /system/extension/toggle`).
func (c *SysExtensionController) Toggle(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:extension:toggle"); resp != nil {
		return resp
	}

	name := ctx.Request().Input("name")
	if name == "" || !plugin.R().Has(name) {
		return ctx.Response().Json(http.StatusOK, owl.Fail("插件不存在"))
	}

	enabled := ctx.Request().InputInt("enabled") == 1
	if err := plugin.R().SetEnabled(name, enabled); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "操作成功"))
}

// Clean removes records for plugins that are no longer discovered
// (`DELETE /system/extension/clean`).
func (c *SysExtensionController) Clean(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:extension:toggle"); resp != nil {
		return resp
	}

	count, err := plugin.R().Clean()
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, fmt.Sprintf("已清理 %d 条失效记录", *count)))
}
