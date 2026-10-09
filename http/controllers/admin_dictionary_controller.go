package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/support"
)

// AdminDictionaryController 字典管理（含字典项子资源）。
type AdminDictionaryController struct{}

func NewAdminDictionaryController() *AdminDictionaryController { return &AdminDictionaryController{} }

// Index GET /system/admin_dictionaries（无 _action 返回页面 schema；_action=getData 返回列表数据）。
func (r *AdminDictionaryController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.DictionariesPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminDictionary{})
	if name := ctx.Request().Query("name", ""); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if key := ctx.Request().Query("key", ""); key != "" {
		query = query.Where("key LIKE ?", "%"+key+"%")
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var dictionaries []models.AdminDictionary
	if err := query.OrderBy("id").
		Offset((page - 1) * perPage).Limit(perPage).Find(&dictionaries); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(dictionaries))
	for _, dictionary := range dictionaries {
		items = append(items, r.format(&dictionary))
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

// format 字典行数据（含字典项数量）。
func (r *AdminDictionaryController) format(dictionary *models.AdminDictionary) map[string]any {
	count, _ := facades.Orm().Query().Model(&models.AdminDictionaryItem{}).
		Where("dictionary_id", dictionary.ID).Count()

	return map[string]any{
		"id":         dictionary.ID,
		"name":       dictionary.Name,
		"key":        dictionary.Key,
		"remark":     dictionary.Remark,
		"item_count": count,
		"created_at": dictionary.CreatedAt,
		"updated_at": dictionary.UpdatedAt,
	}
}

// Store POST /system/admin_dictionaries。
func (r *AdminDictionaryController) Store(ctx contractshttp.Context) contractshttp.Response {
	dictionary, errResp := r.fill(ctx, &models.AdminDictionary{}, 0)
	if errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Create(dictionary); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_dictionaries/{id}/edit。
func (r *AdminDictionaryController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var dictionary models.AdminDictionary
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&dictionary); err != nil || dictionary.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(r.format(&dictionary), ""))
}

// Update PUT /system/admin_dictionaries/{id}。
func (r *AdminDictionaryController) Update(ctx contractshttp.Context) contractshttp.Response {
	var dictionary models.AdminDictionary
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&dictionary); err != nil || dictionary.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	if _, errResp := r.fill(ctx, &dictionary, int(dictionary.ID)); errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Save(&dictionary); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// fill 从请求提取字段（Store/Update 共用；excludeID 用于 key 唯一校验）。
func (r *AdminDictionaryController) fill(ctx contractshttp.Context, dictionary *models.AdminDictionary, excludeID int) (*models.AdminDictionary, func(contractshttp.Context) contractshttp.Response) {
	name := ctx.Request().Input("name")
	key := ctx.Request().Input("key")
	if name == "" || key == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("字典名称和 key 不能为空"))
		}
	}

	dup := facades.Orm().Query().Model(&models.AdminDictionary{}).Where("key", key)
	if excludeID > 0 {
		dup = dup.Where("id", "!=", excludeID)
	}
	count, _ := dup.Count()
	if count > 0 {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("字典 key 已存在"))
		}
	}

	dictionary.Name = name
	dictionary.Key = key

	remark := ctx.Request().Input("remark")
	dictionary.Remark = &remark

	return dictionary, nil
}

// Destroy DELETE /system/admin_dictionaries/{id}（级联删除字典项）。
func (r *AdminDictionaryController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminDictionary{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}
	if _, err := facades.Orm().Query().WhereIn("dictionary_id", toAny(ids)).Delete(&models.AdminDictionaryItem{}); err != nil {
		facades.Log().Warningf("[admin] delete dictionary items failed: %v", err)
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}

// ---- 字典项子资源 ----

// ItemsIndex GET /system/admin_dictionary_items?dictionary_id=&_action=getData。
func (r *AdminDictionaryController) ItemsIndex(ctx contractshttp.Context) contractshttp.Response {
	dictionaryID := r.dictionaryID(ctx)
	if dictionaryID <= 0 {
		return ctx.Response().Json(contractshttp.StatusOK, listEnvelope([]any{}, 0))
	}

	var items []models.AdminDictionaryItem
	if err := facades.Orm().Query().
		Where("dictionary_id", dictionaryID).
		OrderBy("custom_order").OrderBy("id").Find(&items); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	list := make([]map[string]any, 0, len(items))
	for _, item := range items {
		list = append(list, map[string]any{
			"id":           item.ID,
			"dictionary_id": item.DictionaryID,
			"label":        item.Label,
			"value":        item.Value,
			"custom_order": item.CustomOrder,
		})
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(list, int64(len(list))))
}

// ItemsStore POST /system/admin_dictionary_items?dictionary_id=。
func (r *AdminDictionaryController) ItemsStore(ctx contractshttp.Context) contractshttp.Response {
	dictionaryID := r.dictionaryID(ctx)
	if dictionaryID <= 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("缺少字典 ID"))
	}

	var count int64
	count, _ = facades.Orm().Query().Model(&models.AdminDictionary{}).Where("id", dictionaryID).Count()
	if count == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("字典不存在"))
	}

	item, errResp := r.fillItem(ctx, &models.AdminDictionaryItem{})
	if errResp != nil {
		return errResp(ctx)
	}
	item.DictionaryID = dictionaryID

	if err := facades.Orm().Query().Create(item); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// ItemsEdit GET /system/admin_dictionary_items/{id}/edit。
func (r *AdminDictionaryController) ItemsEdit(ctx contractshttp.Context) contractshttp.Response {
	var item models.AdminDictionaryItem
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&item); err != nil || item.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"id":           item.ID,
		"dictionary_id": item.DictionaryID,
		"label":        item.Label,
		"value":        item.Value,
		"custom_order": item.CustomOrder,
	}, ""))
}

// ItemsUpdate PUT /system/admin_dictionary_items/{id}。
func (r *AdminDictionaryController) ItemsUpdate(ctx contractshttp.Context) contractshttp.Response {
	var item models.AdminDictionaryItem
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&item); err != nil || item.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	if _, errResp := r.fillItem(ctx, &item); errResp != nil {
		return errResp(ctx)
	}

	if err := facades.Orm().Query().Save(&item); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// ItemsDestroy DELETE /system/admin_dictionary_items/{id}（支持 ?ids= 批量）。
func (r *AdminDictionaryController) ItemsDestroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminDictionaryItem{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}

// dictionaryID 从 query 或 body 提取字典 ID。
func (r *AdminDictionaryController) dictionaryID(ctx contractshttp.Context) int {
	if id := ctx.Request().QueryInt("dictionary_id", 0); id > 0 {
		return id
	}

	return ctx.Request().InputInt("dictionary_id", 0)
}

// fillItem 字典项字段提取（Store/Update 共用）。
func (r *AdminDictionaryController) fillItem(ctx contractshttp.Context, item *models.AdminDictionaryItem) (*models.AdminDictionaryItem, func(contractshttp.Context) contractshttp.Response) {
	label := ctx.Request().Input("label")
	value := ctx.Request().Input("value")
	if label == "" || value == "" {
		return nil, func(c contractshttp.Context) contractshttp.Response {
			return c.Response().Json(contractshttp.StatusOK, support.Fail("标签和值不能为空"))
		}
	}

	item.Label = label
	item.Value = value
	item.CustomOrder = ctx.Request().InputInt("custom_order", item.CustomOrder)

	return item, nil
}
