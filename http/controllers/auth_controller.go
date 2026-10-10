package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support"
)

// AuthController 登录/登出/当前用户/验证码，对齐 owl-admin AuthController。
type AuthController struct{}

func NewAuthController() *AuthController { return &AuthController{} }

// Login POST /admin-api/login
// 内置基于 IP 的登录失败限流：连续失败 ADMIN_LOGIN_MAX_ATTEMPTS 次
// （默认 5）锁定 ADMIN_LOGIN_LOCK_MINUTES 分钟（默认 10），成功登录后清零。
func (r *AuthController) Login(ctx contractshttp.Context) contractshttp.Response {
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")

	if username == "" {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("用户名不能为空"))
	}
	if password == "" {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("密码不能为空"))
	}

	// 验证码校验（开关开启时强制）：一次性消费，失败不计入登录限流
	if services.CaptchaEnabled() && !services.VerifyLoginCaptcha(ctx.Request().Input("captcha")) {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("验证码错误"))
	}

	throttle := services.NewLoginThrottle()
	ip := ctx.Request().Ip()
	if throttle.TooManyAttempts(ip) {
		facades.Log().Warningf("[admin] login throttled for ip %s (user %s)", ip, username)

		return ctx.Response().Json(contractshttp.StatusOK, support.Fail(services.LoginThrottleMessage(services.DefaultLoginLockMinutes)))
	}

	var user models.AdminUser
	if err := facades.Orm().Query().Where("username", username).First(&user); err != nil {
		facades.Log().Warningf("[admin] query user %s failed: %v", username, err)
	}
	if user.ID == 0 || !facades.Hash().Check(password, user.Password) {
		throttle.Hit(ip)

		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("登录失败，请检查用户名或密码"))
	}

	if !user.IsEnabled() {
		throttle.Hit(ip)

		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("用户被禁用"))
	}

	token, err := facades.Auth(ctx).Guard("admin").Login(&user)
	if err != nil {
		facades.Log().Errorf("[admin] login failed: %v", err)

		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("登录失败："+err.Error()))
	}

	throttle.Clear(ip)

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{"token": token}, "登录成功"))
}

// Logout GET /admin-api/logout
func (r *AuthController) Logout(ctx contractshttp.Context) contractshttp.Response {
	if err := facades.Auth(ctx).Guard("admin").Logout(); err != nil {
		facades.Log().Warningf("[admin] logout failed: %v", err)
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("退出成功"))
}

// CurrentUser GET /admin-api/current-user
func (r *AuthController) CurrentUser(ctx contractshttp.Context) contractshttp.Response {
	user := &models.AdminUser{}
	if err := facades.Auth(ctx).Guard("admin").User(user); err != nil || user.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Unauthorized("请先登录"))
	}

	name := user.Name
	if name == "" {
		name = user.Username
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"name":   name,
		"avatar": user.Avatar,
		"menus":  pages.UserDropdown(name, user.Avatar),
	}, ""))
}

// Captcha GET /admin-api/captcha
// 对齐 owl-admin 契约：sys_captcha 控制前端是否渲染验证码输入框，
// captcha_img 为可直接作为 <img src> 的 data URI；答案入 Cache（5 分钟过期）。
func (r *AuthController) Captcha(ctx contractshttp.Context) contractshttp.Response {
	if !services.CaptchaEnabled() {
		return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
			"sys_captcha": false,
			"captcha_img": "",
		}, ""))
	}

	img, err := services.GenerateCaptchaImage()
	if err != nil {
		facades.Log().Errorf("[admin] generate captcha failed: %v", err)

		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("验证码生成失败"))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"sys_captcha": true,
		"captcha_img": img,
	}, ""))
}
