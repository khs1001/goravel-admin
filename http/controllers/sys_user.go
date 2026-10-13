package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/owl"
)

// SysUserController implements the RuoYi `/system/user` action endpoints.
type SysUserController struct {
	BaseController
}

// NewSysUserController creates a SysUserController.
func NewSysUserController() *SysUserController {
	return &SysUserController{}
}

// List returns a paginated user list (`GET /system/user/list`).
func (c *SysUserController) List(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:user:list"); resp != nil {
		return resp
	}

	page, perPage := c.page(ctx)
	query := facades.Orm().Query().Where("del_flag", "0")
	if keyword := strings.TrimSpace(ctx.Request().Query("userName")); keyword != "" {
		query = query.Where("user_name like ?", "%"+keyword+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	if scope := services.NewDataScopeService().Resolve(c.user(ctx)); !scope.Unrestricted {
		parts := make([]string, 0, 2)
		args := make([]any, 0, 2)
		if len(scope.DeptIDs) > 0 {
			ids := make([]any, len(scope.DeptIDs))
			for i, deptID := range scope.DeptIDs {
				ids[i] = deptID
			}
			parts = append(parts, "dept_id in ?")
			args = append(args, ids)
		}
		if scope.SelfOnly {
			parts = append(parts, "user_id = ?")
			args = append(args, scope.UserID)
		}
		if len(parts) == 0 {
			query = query.Where("1 = 0")
		} else {
			query = query.Where(strings.Join(parts, " or "), args...)
		}
	}

	var users []models.SysUser
	var total int64
	if err := query.OrderByDesc("user_id").Paginate(page, perPage, &users, &total); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	c.attachRoleIDs(users)
	c.attachPostIDs(users)

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"items": users, "total": total}))
}

// Get returns a single user (`GET /system/user/{id}`).
func (c *SysUserController) Get(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:user:query"); resp != nil {
		return resp
	}

	var user models.SysUser
	if err := facades.Orm().Query().Where("user_id", ctx.Request().RouteInt("id")).First(&user); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("用户不存在"))
	}

	users := []models.SysUser{user}
	c.attachRoleIDs(users)
	c.attachPostIDs(users)

	return ctx.Response().Json(http.StatusOK, owl.Success(users[0]))
}

// Add creates a user (`POST /system/user`).
func (c *SysUserController) Add(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:user:add"); resp != nil {
		return resp
	}

	userName := strings.TrimSpace(ctx.Request().Input("userName"))
	if userName == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("用户名不能为空"))
	}

	password := ctx.Request().Input("password")
	if password == "" {
		password = "123456"
	}
	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	user := models.SysUser{
		DeptID:      uint(ctx.Request().InputInt("deptId")),
		UserName:    userName,
		NickName:    ctx.Request().Input("nickName"),
		Email:       ctx.Request().Input("email"),
		Phonenumber: ctx.Request().Input("phonenumber"),
		Sex:         defaultValue(ctx.Request().Input("sex"), "0"),
		Status:      defaultValue(ctx.Request().Input("status"), "0"),
		UserType:    "00",
		DelFlag:     "0",
		Password:    hashed,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	c.assignRoles(user.UserID, ctx.Request().InputArray("roleIds"))
	c.assignPosts(user.UserID, ctx.Request().InputArray("postIds"))

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "新增成功"))
}

// Edit updates a user (`PUT /system/user`).
func (c *SysUserController) Edit(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:user:edit"); resp != nil {
		return resp
	}

	var user models.SysUser
	if err := facades.Orm().Query().Where("user_id", ctx.Request().InputInt("userId")).First(&user); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("用户不存在"))
	}

	user.NickName = ctx.Request().Input("nickName")
	user.Email = ctx.Request().Input("email")
	user.Phonenumber = ctx.Request().Input("phonenumber")
	user.DeptID = uint(ctx.Request().InputInt("deptId"))
	if sex := ctx.Request().Input("sex"); sex != "" {
		user.Sex = sex
	}
	if status := ctx.Request().Input("status"); status != "" {
		user.Status = status
	}
	if password := ctx.Request().Input("password"); password != "" {
		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
		}
		user.Password = hashed
	}

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	c.assignRoles(user.UserID, ctx.Request().InputArray("roleIds"))
	c.assignPosts(user.UserID, ctx.Request().InputArray("postIds"))

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "修改成功"))
}

// Remove deletes a user (`DELETE /system/user/{id}`).
func (c *SysUserController) Remove(ctx http.Context) http.Response {
	if resp := c.denied(ctx, "system:user:remove"); resp != nil {
		return resp
	}

	id := ctx.Request().RouteInt("id")
	if id == 1 {
		return ctx.Response().Json(http.StatusOK, owl.Fail("不允许删除超级管理员"))
	}

	if _, err := facades.Orm().Query().Where("user_id", id).Delete(&models.SysUserRole{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}
	if _, err := facades.Orm().Query().Where("user_id", id).Delete(&models.SysUser{}); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "删除成功"))
}

// RoleOptions returns the roles available for assignment.
func (c *SysUserController) RoleOptions(ctx http.Context) http.Response {
	var roles []models.SysRole
	if err := facades.Orm().Query().Where("del_flag", "0").OrderBy("role_sort").Get(&roles); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	options := make([]map[string]any, 0, len(roles))
	for _, role := range roles {
		options = append(options, map[string]any{"label": role.RoleName, "value": role.RoleID})
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(options))
}

func (c *SysUserController) attachRoleIDs(users []models.SysUser) {
	if len(users) == 0 {
		return
	}

	ids := make([]any, len(users))
	index := make(map[uint]int, len(users))
	for i := range users {
		ids[i] = users[i].UserID
		index[users[i].UserID] = i
	}

	var bindings []models.SysUserRole
	if err := facades.Orm().Query().WhereIn("user_id", ids).Get(&bindings); err != nil {
		return
	}

	for _, binding := range bindings {
		i, ok := index[binding.UserID]
		if !ok {
			continue
		}
		users[i].RoleIDs = append(users[i].RoleIDs, binding.RoleID)
	}
}

func (c *SysUserController) assignRoles(userID uint, roleIDs []string) {
	if _, err := facades.Orm().Query().Where("user_id", userID).Delete(&models.SysUserRole{}); err != nil {
		return
	}

	for _, raw := range roleIDs {
		roleID := parseUint(raw)
		if roleID == 0 {
			continue
		}
		_ = facades.Orm().Query().Create(&models.SysUserRole{UserID: userID, RoleID: roleID})
	}
}

func (c *SysUserController) attachPostIDs(users []models.SysUser) {
	if len(users) == 0 {
		return
	}

	ids := make([]any, len(users))
	index := make(map[uint]int, len(users))
	for i := range users {
		ids[i] = users[i].UserID
		index[users[i].UserID] = i
	}

	var bindings []models.SysUserPost
	if err := facades.Orm().Query().WhereIn("user_id", ids).Get(&bindings); err != nil {
		return
	}

	for _, binding := range bindings {
		i, ok := index[binding.UserID]
		if !ok {
			continue
		}
		users[i].PostIDs = append(users[i].PostIDs, binding.PostID)
	}
}

func (c *SysUserController) assignPosts(userID uint, postIDs []string) {
	if _, err := facades.Orm().Query().Where("user_id", userID).Delete(&models.SysUserPost{}); err != nil {
		return
	}

	for _, raw := range postIDs {
		postID := parseUint(raw)
		if postID == 0 {
			continue
		}
		_ = facades.Orm().Query().Create(&models.SysUserPost{UserID: userID, PostID: postID})
	}
}
