package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysPostController implements the RuoYi `/system/post` action endpoints.
type SysPostController struct {
	BaseController
}

// NewSysPostController creates a SysPostController.
func NewSysPostController() *SysPostController {
	return &SysPostController{}
}

// List returns a paginated post list (`GET /system/post/list`).
func (c *SysPostController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:post:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query()
	if keyword := strings.TrimSpace(ctx.Request().Query("postName")); keyword != "" {
		query = query.Where("post_name like ?", "%"+keyword+"%")
	}
	if code := strings.TrimSpace(ctx.Request().Query("postCode")); code != "" {
		query = query.Where("post_code like ?", "%"+code+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	var posts []models.SysPost
	var total int64
	if err := query.OrderBy("post_sort").Paginate(page, perPage, &posts, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": posts, "total": total}))
}

// Options returns all posts for selectors (`GET /system/post/options`).
func (c *SysPostController) Options(ctx http.Context) http.Response {
	var posts []models.SysPost
	if err := facades.Orm().Query().Where("status", "0").OrderBy("post_sort").Get(&posts); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	options := make([]map[string]any, 0, len(posts))
	for _, post := range posts {
		options = append(options, map[string]any{"label": post.PostName, "value": post.PostID})
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(options))
}

// Get returns a single post (`GET /system/post/{id}`).
func (c *SysPostController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:post:query"); resp != nil {
		return resp
	}

	var post models.SysPost
	if err := facades.Orm().Query().Where("post_id", ctx.Request().RouteInt("id")).First(&post); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("岗位不存在"))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(post))
}

// Add creates a post (`POST /system/post`).
func (c *SysPostController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:post:add"); resp != nil {
		return resp
	}

	post := c.fromRequest(ctx)
	if post.PostName == "" || post.PostCode == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("岗位名称和编码不能为空"))
	}

	if err := facades.Orm().Query().Create(&post); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a post (`PUT /system/post`).
func (c *SysPostController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:post:edit"); resp != nil {
		return resp
	}

	var post models.SysPost
	if err := facades.Orm().Query().Where("post_id", ctx.Request().InputInt("postId")).First(&post); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("岗位不存在"))
	}

	post.PostCode = ctx.Request().Input("postCode")
	post.PostName = ctx.Request().Input("postName")
	post.PostSort = ctx.Request().InputInt("postSort")
	if status := ctx.Request().Input("status"); status != "" {
		post.Status = status
	}
	post.Remark = ctx.Request().Input("remark")

	if err := facades.Orm().Query().Save(&post); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a post (`DELETE /system/post/{id}`).
func (c *SysPostController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:post:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")

	if _, err := facades.Orm().Query().Where("post_id", id).Delete(&models.SysUserPost{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("post_id", id).Delete(&models.SysPost{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

func (c *SysPostController) fromRequest(ctx http.Context) models.SysPost {
	return models.SysPost{
		PostCode: ctx.Request().Input("postCode"),
		PostName: ctx.Request().Input("postName"),
		PostSort: ctx.Request().InputInt("postSort"),
		Status:   defaultValue(ctx.Request().Input("status"), "0"),
		BaseModel: models.BaseModel{
			Remark: ctx.Request().Input("remark"),
		},
	}
}
