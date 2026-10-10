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

// AdminMenuController 菜单管理（对齐 owl-admin AdminMenuController）。
type AdminMenuController struct{}

func NewAdminMenuController() *AdminMenuController { return &AdminMenuController{} }

// Index GET /system/admin_menus（树形数据，loadDataOnce 一次拉全量）。
func (r *AdminMenuController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.MenusPage(), ""))
	}

	tree, err := r.tree()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(tree, int64(len(tree))))
}

// tree 菜单树（items 为嵌套 children 结构，对齐 AdminMenuService::list）。
func (r *AdminMenuController) tree() ([]map[string]any, error) {
	var menus []models.AdminMenu
	if err := facades.Orm().Query().OrderBy("custom_order").Find(&menus); err != nil {
		return nil, err
	}

	return menuTree(menus, 0), nil
}

func menuTree(menus []models.AdminMenu, parentID int) []map[string]any {
	out := []map[string]any{}
	for _, menu := range menus {
		if menu.ParentID != parentID {
			continue
		}

		node := map[string]any{
			"id":           menu.ID,
			"parent_id":    menu.ParentID,
			"title":        menu.Title,
			"icon":         menu.GetIcon(),
			"url":          menu.GetURL(),
			"url_type":     menu.UrlType,
			"component":    menu.GetComponent(),
			"visible":      menu.Visible,
			"is_home":      menu.IsHome,
			"custom_order": menu.CustomOrder,
			"created_at":   menu.CreatedAt,
			"updated_at":   menu.UpdatedAt,
		}
		if children := menuTree(menus, int(menu.ID)); len(children) > 0 {
			node["children"] = children
		}
		out = append(out, node)
	}

	return out
}

// TreeOptions GET /system/admin_menus/tree（表单 TreeSelect 选项）。
func (r *AdminMenuController) TreeOptions(ctx contractshttp.Context) contractshttp.Response {
	var menus []models.AdminMenu
	if err := facades.Orm().Query().OrderBy("custom_order").Find(&menus); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success([]any{}, ""))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(menuTreeOptions(menus, 0), ""))
}

func menuTreeOptions(menus []models.AdminMenu, parentID int) []map[string]any {
	out := []map[string]any{}
	for _, menu := range menus {
		if menu.ParentID != parentID {
			continue
		}

		node := map[string]any{"id": menu.ID, "title": menu.Title}
		if children := menuTreeOptions(menus, int(menu.ID)); len(children) > 0 {
			node["children"] = children
		}
		out = append(out, node)
	}

	return out
}

// Store POST /system/admin_menus。
func (r *AdminMenuController) Store(ctx contractshttp.Context) contractshttp.Response {
	menu, errResp := r.fill(ctx, &models.AdminMenu{})
	if errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Create(menu); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_menus/{id}/edit。
func (r *AdminMenuController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var menu models.AdminMenu
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&menu); err != nil || menu.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"id":           menu.ID,
		"parent_id":    menu.ParentID,
		"title":        menu.Title,
		"icon":         menu.GetIcon(),
		"url":          menu.GetURL(),
		"url_type":     menu.UrlType,
		"component":    menu.GetComponent(),
		"visible":      menu.Visible,
		"is_home":      menu.IsHome,
		"custom_order": menu.CustomOrder,
	}, ""))
}

// Update PUT /system/admin_menus/{id}。
func (r *AdminMenuController) Update(ctx contractshttp.Context) contractshttp.Response {
	var menu models.AdminMenu
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&menu); err != nil || menu.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	old := menu

	if _, errResp := r.fill(ctx, &menu); errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Save(&menu); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	services.DispatchChangeLog(ctx, "admin_menus", strconv.FormatUint(uint64(menu.ID), 10),
		services.DiffFields(&old, &menu, []string{"parent_id", "title", "icon", "url", "url_type", "visible", "custom_order", "is_home", "keep_alive", "iframe_url", "component", "is_full"}))

	// 设为首页时清除其他首页标记（对齐 AdminMenuService::changeHomePage）
	if menu.IsHome == 1 {
		_, _ = facades.Orm().Query().
			Model(&models.AdminMenu{}).
			Where("is_home", 1).Where("id", "!=", menu.ID).
			Update("is_home", 0)
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// fill 从请求提取字段（Store/Update 共用）。
func (r *AdminMenuController) fill(ctx contractshttp.Context, menu *models.AdminMenu) (*models.AdminMenu, func(contractshttp.Context) contractshttp.Response) {
	title := ctx.Request().Input("title")
	if title == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("标题不能为空"))
		}
	}

	menu.Title = title
	menu.ParentID = ctx.Request().InputInt("parent_id", menu.ParentID)
	menu.CustomOrder = ctx.Request().InputInt("custom_order", menu.CustomOrder)
	menu.UrlType = ctx.Request().InputInt("url_type", menu.UrlType)
	menu.Visible = ctx.Request().InputInt("visible", menu.Visible)
	menu.IsHome = ctx.Request().InputInt("is_home", menu.IsHome)

	icon := ctx.Request().Input("icon")
	menu.Icon = &icon

	url := ctx.Request().Input("url")
	menu.URL = &url

	component := ctx.Request().Input("component")
	menu.Component = &component

	return menu, nil
}

// SaveOrder POST /system/admin_menus/save_order（ids 按新顺序全量更新 custom_order）。
func (r *AdminMenuController) SaveOrder(ctx contractshttp.Context) contractshttp.Response {
	ids := intSlice(ctx.Request().All()["ids"])
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定排序记录"))
	}

	for order, id := range ids {
		if _, err := facades.Orm().Query().
			Model(&models.AdminMenu{}).
			Where("id", id).
			Update("custom_order", order); err != nil {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存排序失败："+err.Error()))
		}
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("保存成功"))
}

// Destroy DELETE /system/admin_menus/{id}。
func (r *AdminMenuController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminMenu{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}
	_, _ = facades.Orm().Query().WhereIn("menu_id", toAny(ids)).Delete(&models.AdminPermissionMenu{})

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}
