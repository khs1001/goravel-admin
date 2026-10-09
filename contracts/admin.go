package contracts

import contractsroute "github.com/goravel/framework/contracts/route"

// PluginKey 是插件的容器绑定 key 类型。
//
// 业务扩展包在自己的 ServiceProvider.Register 中，以 PluginKey 作为绑定 key
// 将 Plugin 实现注册进容器：
//
//	app.Bind(admin.PluginKey("blog"), func(app foundation.Application) (any, error) {
//	    return &blog.Plugin{}, nil
//	})
//
// admin 包在 Boot 阶段枚举容器中所有 PluginKey 类型的绑定并解析收集，
// 从而实现零修改挂载业务菜单、页面与路由。
type PluginKey string

// Admin 是 admin 包的容器绑定契约。
type Admin interface {
	// Name 后台应用名称 (config admin.name)
	Name() string
	// Logo 后台 logo 地址 (config admin.logo)
	Logo() string
	// RoutePrefix 后台 API 路由前缀 (config admin.route.prefix)
	RoutePrefix() string
	// RegisterPlugin 手动注册插件（除容器绑定外的补充方式）
	RegisterPlugin(plugin Plugin)
	// Plugins 已注册的全部插件
	Plugins() []Plugin
}

// Plugin 是业务扩展包接入 admin 的注册契约。
//
// 插件本身也应当是一个 goravel 扩展包（有自己的 setup 安装机制），
// 通过容器绑定（推荐）或 Admin.RegisterPlugin 注册。
type Plugin interface {
	// Name 插件唯一标识
	Name() string
	// Menus 追加到菜单树的代码级菜单（返回完整的 MenuItem，含 name/path/meta）
	Menus() []MenuItem
	// Pages 注册内置页面：sign → amis 页面 schema
	Pages() map[string]any
	// Routes 在 admin-api 前缀下、鉴权中间件之后注册业务路由
	Routes(router contractsroute.Router)
}

// MenuItem 对应 owl-admin 前端 IRoute（/menus 端点返回的菜单树节点）。
type MenuItem struct {
	Name      string     `json:"name"`
	Path      string     `json:"path,omitempty"`
	Component string     `json:"component,omitempty"`
	IsHome    int        `json:"is_home"`
	IframeURL string     `json:"iframe_url,omitempty"`
	UrlType   int        `json:"url_type"`
	KeepAlive *int       `json:"keep_alive"`
	IsFull    int        `json:"is_full"`
	IsLink    bool       `json:"is_link"`
	PageSign  string     `json:"page_sign,omitempty"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuItem `json:"children,omitempty"`
}

// MenuMeta 菜单元信息。
type MenuMeta struct {
	Title        string `json:"title"`
	Icon         string `json:"icon,omitempty"`
	Hide         bool   `json:"hide,omitempty"`
	CustomOrder  int    `json:"custom_order,omitempty"`
	SingleLayout string `json:"singleLayout,omitempty"`
}

// 菜单 url_type 常量，与 owl-admin AdminMenu 模型对齐。
const (
	MenuTypeRoute  = 1 // 站内路由
	MenuTypeLink   = 2 // 外链
	MenuTypeIframe = 3 // iframe
	MenuTypePage   = 4 // 页面（page_schema 下发）
)
