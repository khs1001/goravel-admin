package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support/captcha"
	"github.com/khs1001/goravel-admin/support/current"
	"github.com/khs1001/goravel-admin/support/owl"
)

// AuthController implements the Owl authentication endpoints: /login, /logout,
// /current-user and /captcha.
type AuthController struct {
	BaseController
}

// NewAuthController creates an AuthController.
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Login authenticates an admin user and returns a JWT.
func (c *AuthController) Login(ctx http.Context) http.Response {
	if facades.Config().GetBool("admin.auth.login_captcha", true) {
		sysCaptcha := ctx.Request().Input("sys_captcha")
		code := ctx.Request().Input("captcha")
		if sysCaptcha == "" || code == "" {
			return ctx.Response().Json(http.StatusOK, owl.Fail("请输入验证码"))
		}
		if !captcha.Verify(sysCaptcha, code) {
			return ctx.Response().Json(http.StatusOK, owl.Fail("验证码错误"))
		}
	}

	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	if username == "" || password == "" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("请输入用户名和密码"))
	}

	var user models.SysUser
	if err := facades.Orm().Query().Where("user_name", username).Where("del_flag", "0").First(&user); err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail("用户名或密码错误"))
	}
	if user.Status == "1" {
		return ctx.Response().Json(http.StatusOK, owl.Fail("账号已被停用"))
	}
	if !facades.Hash().Check(password, user.Password) {
		return ctx.Response().Json(http.StatusOK, owl.Fail("用户名或密码错误"))
	}

	token, err := facades.Auth(ctx).Guard("admin").Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	user.LoginIP = ctx.Request().Ip()
	user.LoginDate = carbon.NewDateTime(carbon.Now())
	_ = facades.Orm().Query().Save(&user)

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"token": token}, "登录成功"))
}

// Logout is a best-effort, stateless operation: JWT tokens are not revoked
// server-side in v1.
func (c *AuthController) Logout(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, owl.Success(nil, "已退出登录"))
}

// Captcha returns a fresh image captcha.
func (c *AuthController) Captcha(ctx http.Context) http.Response {
	result, err := captcha.Generate()
	if err != nil {
		return ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{
		"captcha_img": result.Img,
		"sys_captcha": result.SysCaptcha,
	}))
}

// CurrentUser returns the profile of the authenticated user together with the
// permission set and the navbar amis schema.
func (c *AuthController) CurrentUser(ctx http.Context) http.Response {
	user := current.User(ctx)
	if user == nil {
		return ctx.Response().Json(http.StatusOK, owl.Unauthorized("请先登录"))
	}

	menuService := services.NewMenuService()
	perms, _ := menuService.PermsForUser(user)

	avatar := user.Avatar
	if avatar == "" {
		avatar = facades.Config().GetString("admin.default_avatar", "")
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{
		"id":          user.UserID,
		"username":    user.UserName,
		"name":        user.NickName,
		"avatar":      avatar,
		"roles":       c.roleKeys(user.UserID),
		"permissions": perms,
		"menus":       userMenuSchema(user),
	}))
}

func (c *AuthController) roleKeys(userID uint) []string {
	var bindings []models.SysUserRole
	if err := facades.Orm().Query().Where("user_id", userID).Get(&bindings); err != nil {
		return []string{}
	}

	ids := make([]any, 0, len(bindings))
	for _, b := range bindings {
		ids = append(ids, b.RoleID)
	}
	if len(ids) == 0 {
		return []string{}
	}

	var roles []models.SysRole
	if err := facades.Orm().Query().WhereIn("role_id", ids).Get(&roles); err != nil {
		return []string{}
	}

	keys := make([]string, 0, len(roles))
	for _, role := range roles {
		keys = append(keys, role.RoleKey)
	}

	return keys
}

// userMenuSchema is the amis schema rendered at the top-right of the layout.
func userMenuSchema(user *models.SysUser) map[string]any {
	return map[string]any{
		"type":         "dropdown-button",
		"label":        user.NickName,
		"level":        "link",
		"hideCaret":    true,
		"align":        "right",
		"className":    "h-full",
		"closeOnClick": true,
		"buttons": []map[string]any{
			{
				"type":  "button",
				"label": "退出登录",
				"level": "link",
				"onEvent": map[string]any{
					"click": map[string]any{
						"actions": []map[string]any{
							{"actionType": "custom", "script": "window.$owl.logout()"},
						},
					},
				},
			},
		},
	}
}
