package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysMenuController implements the RuoYi `/system/menu` action endpoints. In the
// RuoYi permission model the menu tree *is* the permission tree: nodes of type
// F carry the permission identifiers consumed by the API authorization.
type SysMenuController struct {
	BaseController
}

// NewSysMenuController creates a SysMenuController.
func NewSysMenuController() *SysMenuController {
	return &SysMenuController{}
}

// List returns the full menu tree (`GET /system/menu/list`).
func (c *SysMenuController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:menu:list"); resp != nil {
		return resp
	}

	var menus []models.SysMenu
	if err := facades.Orm().Query().OrderBy("order_num").Get(&menus); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(owl.BuildAllTree(menus, 0)))
}

// Tree returns the menu tree used as a parent selector (`GET /system/menu/tree`).
func (c *SysMenuController) Tree(ctx http.Context) http.Response {
	var menus []models.SysMenu
	if err := facades.Orm().Query().Where("menu_type != ?", models.MenuTypeButton).OrderBy("order_num").Get(&menus); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(owl.BuildAllTree(menus, 0)))
}

// Get returns a single menu (`GET /system/menu/{id}`).
func (c *SysMenuController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:menu:query"); resp != nil {
		return resp
	}

	var menu models.SysMenu
	if err := facades.Orm().Query().Where("menu_id", ctx.Request().RouteInt("id")).FirstOrFail(&menu); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("菜单不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(menu))
}

// Add creates a menu (`POST /system/menu`).
func (c *SysMenuController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:menu:add"); resp != nil {
		return resp
	}

	menu := c.fromRequest(ctx)
	if menu.MenuName == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("菜单名称不能为空"))
	}

	if err := facades.Orm().Query().Create(&menu); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a menu (`PUT /system/menu`).
func (c *SysMenuController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:menu:edit"); resp != nil {
		return resp
	}

	var menu models.SysMenu
	if err := facades.Orm().Query().Where("menu_id", ctx.Request().InputInt("menuId")).FirstOrFail(&menu); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("菜单不存在"))
	}

	menu.MenuName = ctx.Request().Input("menuName")
	menu.ParentID = uint(ctx.Request().InputInt("parentId"))
	menu.OrderNum = ctx.Request().InputInt("orderNum")
	menu.Path = ctx.Request().Input("path")
	menu.Component = ctx.Request().Input("component")
	menu.Query = ctx.Request().Input("query")
	menu.IsFrame = ctx.Request().InputInt("isFrame")
	menu.IsCache = ctx.Request().InputInt("isCache")
	menu.MenuType = defaultValue(ctx.Request().Input("menuType"), models.MenuTypeMenu)
	menu.Visible = defaultValue(ctx.Request().Input("visible"), "0")
	menu.Status = defaultValue(ctx.Request().Input("status"), "0")
	menu.Perms = ctx.Request().Input("perms")
	menu.Icon = ctx.Request().Input("icon")

	if err := facades.Orm().Query().Save(&menu); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a menu (`DELETE /system/menu/{id}`).
func (c *SysMenuController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:menu:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")

	var children []models.SysMenu
	if err := facades.Orm().Query().Where("parent_id", id).Get(&children); err == nil && len(children) > 0 {
		return ctx.Response().Json(http.StatusOK, owl.Fail("存在子菜单，不允许删除"))
	}

	if _, err := facades.Orm().Query().Where("menu_id", id).Delete(&models.SysRoleMenu{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("menu_id", id).Delete(&models.SysMenu{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

func (c *SysMenuController) fromRequest(ctx http.Context) models.SysMenu {
	return models.SysMenu{
		MenuName:  ctx.Request().Input("menuName"),
		ParentID:  uint(ctx.Request().InputInt("parentId")),
		OrderNum:  ctx.Request().InputInt("orderNum"),
		Path:      ctx.Request().Input("path"),
		Component: ctx.Request().Input("component"),
		Query:     ctx.Request().Input("query"),
		IsFrame:   ctx.Request().InputInt("isFrame"),
		IsCache:   ctx.Request().InputInt("isCache"),
		MenuType:  defaultValue(ctx.Request().Input("menuType"), models.MenuTypeMenu),
		Visible:   defaultValue(ctx.Request().Input("visible"), "0"),
		Status:    defaultValue(ctx.Request().Input("status"), "0"),
		Perms:     ctx.Request().Input("perms"),
		Icon:      ctx.Request().Input("icon"),
	}
}
