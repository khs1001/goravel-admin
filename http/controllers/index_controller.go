package controllers

import (
	"strconv"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/models"
	"github.com/khs1001/goravel-admin/pages"
	"github.com/khs1001/goravel-admin/services"
	"github.com/khs1001/goravel-admin/support"
)

// IndexController 菜单/设置/页面 schema 端点，对齐 owl-admin IndexController。
type IndexController struct {
	menuService    *services.MenuService
	settingService *services.SettingService
	pageService    *services.PageService
	pluginMenus    []contracts.MenuItem
}

func NewIndexController(
	menuService *services.MenuService,
	settingService *services.SettingService,
	pageService *services.PageService,
) *IndexController {
	return &IndexController{
		menuService:    menuService,
		settingService: settingService,
		pageService:    pageService,
	}
}

// SetPluginMenus 注入插件代码级菜单（RegisterRoutes 阶段调用）。
func (r *IndexController) SetPluginMenus(menus []contracts.MenuItem) {
	r.pluginMenus = menus
}

// Menus GET /admin-api/menus
func (r *IndexController) Menus(ctx contractshttp.Context) contractshttp.Response {
	userID := currentAdminUserID(ctx)

	menus := r.menuService.All(userID, extraMenus(), r.pluginMenus)

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(menus, ""))
}

// Settings GET /admin-api/_settings（免鉴权，登录页也需要）
func (r *IndexController) Settings(ctx contractshttp.Context) contractshttp.Response {
	locale := r.settingService.GetByModule("admin_locale", "zh_CN")

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"nav":      []any{},
		"assets":   []any{},
		"app_name": facades.Config().GetString("admin.name", "Goravel Admin"),
		"locale":   locale,
		"layout":   facades.Config().Get("admin.layout", map[string]any{}),
		// 默认 logo 取嵌入前端产物自带图标，零配置即有品牌展示
		"logo": facades.Config().GetString("admin.logo", "/admin-assets/logo.png"),

		"login_captcha":          facades.Config().GetBool("admin.auth.login_captcha", false),
		"locale_options":         localeOptions(),
		"show_development_tools": facades.Config().GetBool("admin.show_development_tools", false),
		"system_theme_setting":   r.settingService.Get("system_theme_setting"),
		"enabled_extensions":     []string{},
	}, ""))
}

// SaveSettings POST /admin-api/_settings
func (r *IndexController) SaveSettings(ctx contractshttp.Context) contractshttp.Response {
	raw := ctx.Request().All()

	data := map[string]any{}
	for key, value := range raw {
		if key == "" {
			continue
		}
		data[key] = value
	}

	if err := r.settingService.SetMany(data); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("保存成功"))
}

// PageSchema GET /admin-api/page_schema?sign=
func (r *IndexController) PageSchema(ctx contractshttp.Context) contractshttp.Response {
	schema := r.pageService.Get(ctx.Request().Query("sign", ""))

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(schema, ""))
}

// Dashboard GET /admin-api/dashboard（url_type=路由 的菜单直接以路径拉取页面 schema）
func (r *IndexController) Dashboard(ctx contractshttp.Context) contractshttp.Response {
	return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.Dashboard(), ""))
}

// DashboardStats GET /admin-api/dashboard_stats（仪表盘统计卡片与最近日志）。
func (r *IndexController) DashboardStats(ctx contractshttp.Context) contractshttp.Response {
	count := func(model any) int64 {
		n, _ := facades.Orm().Query().Model(model).Count()

		return n
	}

	var recent []models.AdminOperationLog
	_ = facades.Orm().Query().OrderBy("id", "desc").Limit(8).Find(&recent)
	recentLogs := make([]map[string]any, 0, len(recent))
	for _, log := range recent {
		recentLogs = append(recentLogs, map[string]any{
			"username":   log.Username,
			"action":     log.Action,
			"method":     log.Method,
			"path":       log.Path,
			"created_at": log.CreatedAt,
		})
	}

	today := time.Now().Format("2006-01-02")
	todayLogCount, _ := facades.Orm().Query().Model(&models.AdminOperationLog{}).
		Where("DATE(created_at)", today).Count()

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"user_count":      count(&models.AdminUser{}),
		"role_count":      count(&models.AdminRole{}),
		"menu_count":      count(&models.AdminMenu{}),
		"today_log_count": todayLogCount,
		"recent_logs":     recentLogs,
	}, ""))
}

// SystemIndex GET /admin-api/system（系统管理入口页）
func (r *IndexController) SystemIndex(ctx contractshttp.Context) contractshttp.Response {
	return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.SystemIndex(), ""))
}

// NoContent GET /admin-api/no-content
func (r *IndexController) NoContent(ctx contractshttp.Context) contractshttp.Response {
	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage(""))
}

// UserSetting GET /admin-api/user_setting
func (r *AuthController) UserSetting(ctx contractshttp.Context) contractshttp.Response {
	return ctx.Response().Json(contractshttp.StatusOK, support.Success(pages.UserSetting(), ""))
}

// SaveUserSetting PUT /admin-api/user_setting
func (r *AuthController) SaveUserSetting(ctx contractshttp.Context) contractshttp.Response {
	user := &models.AdminUser{}
	if err := facades.Auth(ctx).Guard("admin").User(user); err != nil || user.ID == 0 {
		return ctx.Response().Json(contractshttp.StatusOK, support.Unauthorized("请先登录"))
	}

	old := *user

	if name := ctx.Request().Input("name"); name != "" {
		user.Name = name
	}

	if avatar := ctx.Request().Input("avatar"); avatar != "" {
		user.Avatar = &avatar
	}

	password := ctx.Request().Input("password")
	if password != "" {
		oldPassword := ctx.Request().Input("old_password")
		if !facades.Hash().Check(oldPassword, user.Password) {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("旧密码错误"))
		}

		if password != ctx.Request().Input("confirm_password") {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("两次输入的密码不一致"))
		}

		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return ctx.Response().Json(contractshttp.StatusOK, support.Fail("密码加密失败"))
		}
		user.Password = hashed
	}

	if err := facades.Orm().Query().Save(user); err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}

	services.DispatchChangeLog(ctx, "admin_users", strconv.FormatUint(uint64(user.ID), 10),
		services.DiffFields(&old, user, []string{"name", "avatar"}))

	return ctx.Response().Json(contractshttp.StatusOK, support.SuccessMessage("保存成功"))
}

// ---- helpers ----

// currentAdminUserID 当前登录管理员 ID（未登录返回 0）。
func currentAdminUserID(ctx contractshttp.Context) int {
	user := &models.AdminUser{}
	if err := facades.Auth(ctx).Guard("admin").User(user); err != nil {
		return 0
	}

	return int(user.ID)
}

// extraMenus 额外菜单（对齐 owl-admin Menu::extra），第一阶段仅个人设置。
func extraMenus() []contracts.MenuItem {
	return []contracts.MenuItem{
		{
			Name:      "user_setting",
			Path:      "/user_setting",
			Component: "amis",
			UrlType:   contracts.MenuTypeRoute,
			Meta: contracts.MenuMeta{
				Hide:         true,
				Title:        "个人设置",
				Icon:         "material-symbols:manage-accounts",
				SingleLayout: "basic",
			},
		},
	}
}

func localeOptions() []map[string]string {
	return []map[string]string{
		{"label": "English", "value": "en"},
		{"label": "简体中文", "value": "zh_CN"},
	}
}
