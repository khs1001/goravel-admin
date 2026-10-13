package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysDictController implements the RuoYi `/system/dict` action endpoints, both
// the dictionary types and their data.
type SysDictController struct {
	BaseController
}

// NewSysDictController creates a SysDictController.
func NewSysDictController() *SysDictController {
	return &SysDictController{}
}

// TypeList returns a paginated dictionary type list (`GET /system/dict/type/list`).
func (c *SysDictController) TypeList(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if keyword := strings.TrimSpace(ctx.Request().Query("dictName")); keyword != "" {
		query = query.Where("dict_name like ?", "%"+keyword+"%")
	}
	if dictType := strings.TrimSpace(ctx.Request().Query("dictType")); dictType != "" {
		query = query.Where("dict_type like ?", "%"+dictType+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var items []models.SysDictType
	var total int64
	if err := query.OrderBy("dict_id").Paginate(page, perPage, &items, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": total}))
}

// TypeGet returns a single dictionary type (`GET /system/dict/type/{id}`).
func (c *SysDictController) TypeGet(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:query"); resp != nil {
		return resp
	}

	var item models.SysDictType
	if err := facades.Orm().Query().Where("dict_id", ctx.Request().RouteInt("id")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典类型不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(item))
}

// TypeAdd creates a dictionary type (`POST /system/dict/type`).
func (c *SysDictController) TypeAdd(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:add"); resp != nil {
		return resp
	}

	dictType := strings.TrimSpace(ctx.Request().Input("dictType"))
	dictName := strings.TrimSpace(ctx.Request().Input("dictName"))
	if dictType == "" || dictName == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典名称和类型不能为空"))
	}

	item := models.SysDictType{
		DictName: dictName,
		DictType: dictType,
		Status:   defaultValue(ctx.Request().Input("status"), "0"),
		BaseModel: models.BaseModel{
			Remark: ctx.Request().Input("remark"),
		},
	}
	if err := facades.Orm().Query().Create(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// TypeEdit updates a dictionary type (`PUT /system/dict/type`).
func (c *SysDictController) TypeEdit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:edit"); resp != nil {
		return resp
	}

	var item models.SysDictType
	if err := facades.Orm().Query().Where("dict_id", ctx.Request().InputInt("dictId")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典类型不存在"))
	}

	item.DictName = ctx.Request().Input("dictName")
	item.DictType = ctx.Request().Input("dictType")
	if status := ctx.Request().Input("status"); status != "" {
		item.Status = status
	}
	item.Remark = ctx.Request().Input("remark")

	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// TypeRemove deletes a dictionary type (`DELETE /system/dict/type/{id}`).
func (c *SysDictController) TypeRemove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")

	var item models.SysDictType
	if err := facades.Orm().Query().Where("dict_id", id).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典类型不存在"))
	}

	if _, err := facades.Orm().Query().Where("dict_type", item.DictType).Delete(&models.SysDictData{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("dict_id", id).Delete(&models.SysDictType{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// DataList returns a paginated dictionary data list (`GET /system/dict/data/list`).
func (c *SysDictController) DataList(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if dictType := strings.TrimSpace(ctx.Request().Query("dictType")); dictType != "" {
		query = query.Where("dict_type", dictType)
	}
	if label := strings.TrimSpace(ctx.Request().Query("dictLabel")); label != "" {
		query = query.Where("dict_label like ?", "%"+label+"%")
	}

	var items []models.SysDictData
	var total int64
	if err := query.OrderBy("dict_sort").Paginate(page, perPage, &items, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": items, "total": total}))
}

// DataByType returns the dictionary entries of a type as selector options
// (`GET /system/dict/data/type/{dictType}`).
func (c *SysDictController) DataByType(ctx http.Context) http.Response {
	dictType := ctx.Request().Route("dictType")

	var items []models.SysDictData
	if err := facades.Orm().Query().Where("dict_type", dictType).Where("status", "0").OrderBy("dict_sort").Get(&items); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(items))
}

// DataGet returns a single dictionary entry (`GET /system/dict/data/{id}`).
func (c *SysDictController) DataGet(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:query"); resp != nil {
		return resp
	}

	var item models.SysDictData
	if err := facades.Orm().Query().Where("dict_code", ctx.Request().RouteInt("id")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典数据不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(item))
}

// DataAdd creates a dictionary entry (`POST /system/dict/data`).
func (c *SysDictController) DataAdd(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:add"); resp != nil {
		return resp
	}

	item := c.dictDataFromRequest(ctx)
	if item.DictLabel == "" || item.DictValue == "" || item.DictType == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典标签、键值和类型不能为空"))
	}

	if err := facades.Orm().Query().Create(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// DataEdit updates a dictionary entry (`PUT /system/dict/data`).
func (c *SysDictController) DataEdit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:edit"); resp != nil {
		return resp
	}

	var item models.SysDictData
	if err := facades.Orm().Query().Where("dict_code", ctx.Request().InputInt("dictCode")).First(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("字典数据不存在"))
	}

	item.DictSort = ctx.Request().InputInt("dictSort")
	item.DictLabel = ctx.Request().Input("dictLabel")
	item.DictValue = ctx.Request().Input("dictValue")
	item.DictType = ctx.Request().Input("dictType")
	item.CssClass = ctx.Request().Input("cssClass")
	item.ListClass = ctx.Request().Input("listClass")
	item.IsDefault = defaultValue(ctx.Request().Input("isDefault"), "N")
	if status := ctx.Request().Input("status"); status != "" {
		item.Status = status
	}
	item.Remark = ctx.Request().Input("remark")

	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// DataRemove deletes a dictionary entry (`DELETE /system/dict/data/{id}`).
func (c *SysDictController) DataRemove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:dict:remove"); resp != nil {
		return resp
	}

	if _, err := facades.Orm().Query().Where("dict_code", ctx.Request().RouteInt("id")).Delete(&models.SysDictData{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

func (c *SysDictController) dictDataFromRequest(ctx http.Context) models.SysDictData {
	return models.SysDictData{
		DictSort:  ctx.Request().InputInt("dictSort"),
		DictLabel: ctx.Request().Input("dictLabel"),
		DictValue: ctx.Request().Input("dictValue"),
		DictType:  ctx.Request().Input("dictType"),
		CssClass:  ctx.Request().Input("cssClass"),
		ListClass: ctx.Request().Input("listClass"),
		IsDefault: defaultValue(ctx.Request().Input("isDefault"), "N"),
		Status:    defaultValue(ctx.Request().Input("status"), "0"),
		BaseModel: models.BaseModel{
			Remark: ctx.Request().Input("remark"),
		},
	}
}
