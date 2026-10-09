package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/support"
)

// AdminUserController 管理员管理（对齐 owl-admin AdminUserController）。
type AdminUserController struct{}

func NewAdminUserController() *AdminUserController { return &AdminUserController{} }

// Index GET /system/admin_users（无 _action 返回页面 schema；_action=getData 返回列表数据）。
func (r *AdminUserController) Index(ctx contractshttp.Context) contractshttp.Response {
	if !isGetData(ctx) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.UsersPage(), ""))
	}

	page, perPage := pageParams(ctx)
	query := facades.Orm().Query().Model(&models.AdminUser{})

	if username := ctx.Request().Query("username", ""); username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if name := ctx.Request().Query("name", ""); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if enabled := ctx.Request().Query("enabled", ""); enabled != "" {
		query = query.Where("enabled", enabled)
	}
	if roleIDs := intSlice(rawQueryValue(ctx, "roles")); len(roleIDs) > 0 {
		userIDs := r.userIDsByRoles(roleIDs)
		if len(userIDs) == 0 {
			return ctx.Response().Json(contractshttp.StatusOK, listEnvelope([]any{}, 0))
		}
		query = query.WhereIn("id", toAny(userIDs))
	}

	total, err := query.Count()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	var users []models.AdminUser
	if err := query.OrderBy("id").Offset((page - 1) * perPage).Limit(perPage).Find(&users); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("查询失败："+err.Error()))
	}

	items := make([]map[string]any, 0, len(users))
	for _, user := range users {
		items = append(items, r.formatUser(&user))
	}

	return ctx.Response().Json(contractshttp.StatusOK, listEnvelope(items, total))
}

// userIDsByRoles 角色筛选 → 用户 ID 集合。
func (r *AdminUserController) userIDsByRoles(roleIDs []int) []int {
	var links []models.AdminRoleUsers
	if err := facades.Orm().Query().WhereIn("role_id", toAny(roleIDs)).Find(&links); err != nil {
		return nil
	}

	seen := map[int]bool{}
	ids := []int{}
	for _, link := range links {
		if !seen[link.UserID] {
			seen[link.UserID] = true
			ids = append(ids, link.UserID)
		}
	}

	return ids
}

// formatUser 行数据（含角色）。
func (r *AdminUserController) formatUser(user *models.AdminUser) map[string]any {
	var links []models.AdminRoleUsers
	_ = facades.Orm().Query().Where("user_id", user.ID).Find(&links)

	roles := []map[string]any{}
	if len(links) > 0 {
		roleIDs := make([]any, 0, len(links))
		for _, link := range links {
			roleIDs = append(roleIDs, link.RoleID)
		}

		var roleRows []models.AdminRole
		if err := facades.Orm().Query().WhereIn("id", roleIDs).Find(&roleRows); err == nil {
			for _, role := range roleRows {
				roles = append(roles, map[string]any{"id": role.ID, "name": role.Name})
			}
		}
	}

	return map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"avatar":     user.Avatar,
		"enabled":    user.Enabled,
		"roles":      roles,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}

// Store POST /system/admin_users（新增）。
func (r *AdminUserController) Store(ctx contractshttp.Context) contractshttp.Response {
	username := ctx.Request().Input("username")
	name := ctx.Request().Input("name")
	password := ctx.Request().Input("password")

	if username == "" || name == "" {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("用户名和昵称不能为空"))
	}
	if password == "" {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建用户必须设置密码"))
	}
	if password != ctx.Request().Input("confirm_password") {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("两次输入的密码不一致"))
	}

	var count int64
	count, _ = facades.Orm().Query().Model(&models.AdminUser{}).Where("username", username).Count()
	if count > 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("用户名已存在"))
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("密码加密失败"))
	}

	user := &models.AdminUser{
		Username: username,
		Password: hashed,
		Name:     name,
		Enabled:  ctx.Request().InputInt("enabled", 1),
	}
	if err := facades.Orm().Query().Create(user); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("创建失败："+err.Error()))
	}

	if err := r.syncRoles(int(user.ID), intSlice(ctx.Request().All()["roles"])); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存角色失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("创建成功"))
}

// Edit GET /system/admin_users/{id}/edit（编辑表单初始数据）。
func (r *AdminUserController) Edit(ctx contractshttp.Context) contractshttp.Response {
	var user models.AdminUser
	if err := facades.Orm().Query().Where("id", ctx.Request().RouteInt("id")).First(&user); err != nil || user.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}

	data := r.formatUser(&user)
	data["roles"] = r.roleIDs(int(user.ID))

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(data, ""))
}

// roleIDs 用户拥有的角色 ID。
func (r *AdminUserController) roleIDs(userID int) []int {
	var links []models.AdminRoleUsers
	if err := facades.Orm().Query().Where("user_id", userID).Find(&links); err != nil {
		return []int{}
	}

	ids := make([]int, 0, len(links))
	for _, link := range links {
		ids = append(ids, link.RoleID)
	}

	return ids
}

// Update PUT /system/admin_users/{id}。
func (r *AdminUserController) Update(ctx contractshttp.Context) contractshttp.Response {
	id := ctx.Request().RouteInt("id")

	var user models.AdminUser
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil || user.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("记录不存在"))
	}
	if user.ID == 1 && ctx.Request().InputInt("enabled", user.Enabled) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("超级管理员不可禁用"))
	}

	if name := ctx.Request().Input("name"); name != "" {
		user.Name = name
	}
	if username := ctx.Request().Input("username"); username != "" && username != user.Username {
		var count int64
		count, _ = facades.Orm().Query().Model(&models.AdminUser{}).
			Where("username", username).Where("id", "!=", user.ID).Count()
		if count > 0 {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("用户名已存在"))
		}
		user.Username = username
	}
	all := ctx.Request().All()
	if _, ok := all["enabled"]; ok && user.ID != 1 {
		user.Enabled = ctx.Request().InputInt("enabled", user.Enabled)
	}

	if password := ctx.Request().Input("password"); password != "" {
		if password != ctx.Request().Input("confirm_password") {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("两次输入的密码不一致"))
		}

		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("密码加密失败"))
		}
		user.Password = hashed
	}

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	if user.ID != 1 {
		if err := r.syncRoles(int(user.ID), intSlice(ctx.Request().All()["roles"])); err != nil {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存角色失败："+err.Error()))
		}
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("更新成功"))
}

// syncRoles 全量覆盖用户-角色关联。
func (r *AdminUserController) syncRoles(userID int, roleIDs []int) error {
	if _, err := facades.Orm().Query().Where("user_id", userID).Delete(&models.AdminRoleUsers{}); err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}

	links := make([]models.AdminRoleUsers, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		links = append(links, models.AdminRoleUsers{RoleID: roleID, UserID: userID})
	}

	return facades.Orm().Query().Create(&links)
}

// Destroy DELETE /system/admin_users/{id}（支持 ?ids= 批量）。
func (r *AdminUserController) Destroy(ctx contractshttp.Context) contractshttp.Response {
	ids := idList(ctx)
	if id := ctx.Request().RouteInt("id"); id > 0 {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未指定删除记录"))
	}
	for _, id := range ids {
		if id == 1 {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("超级管理员不可删除"))
		}
	}

	if _, err := facades.Orm().Query().WhereIn("id", toAny(ids)).Delete(&models.AdminUser{}); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("删除失败："+err.Error()))
	}
	if _, err := facades.Orm().Query().WhereIn("user_id", toAny(ids)).Delete(&models.AdminRoleUsers{}); err != nil {
		facades.Log().Warningf("[admin] delete user roles failed: %v", err)
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("删除成功"))
}

// ---- 共享小工具 ----

// rawQueryValue 读取 query 原始值（amis 筛选可能传数组或逗号串）。
func rawQueryValue(ctx contractshttp.Context, key string) any {
	return ctx.Request().All()[key]
}

// toAny []int → []any（WhereIn 参数）。
func toAny(ids []int) []any {
	out := make([]any, 0, len(ids))
	for _, id := range ids {
		out = append(out, id)
	}

	return out
}
