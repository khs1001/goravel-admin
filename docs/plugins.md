# 插件开发指南

admin 包提供二级插件系统：业务扩展包实现 `contracts.Plugin` 接口，即可自动挂载菜单、页面与路由，宿主应用零修改。

## 快速接入

### 1. 实现 Plugin 接口

```go
package blog

import (
    contractsroute "github.com/goravel/framework/contracts/route"

    "github.com/khs1001/goravel-admin/contracts"
)

type Plugin struct{}

func (r *Plugin) Name() string { return "blog" }

// Menus 追加到菜单树的代码级菜单
func (r *Plugin) Menus() []contracts.MenuItem {
    return []contracts.MenuItem{{
        Name:     "blog",
        Path:     "/blog",
        PageSign: "blog-posts",
        UrlType:  contracts.MenuTypePage,
        Meta:     contracts.MenuMeta{Title: "博客管理", Icon: "iconoir-book"},
    }}
}

// Pages 注册页面：sign → amis 页面 schema
func (r *Plugin) Pages() map[string]any {
    return map[string]any{
        "blog-posts": map[string]any{
            "type": "crud",
            // ... amis schema
        },
    }
}

// Routes 在 admin-api 前缀下、鉴权 + 权限中间件之后注册业务路由
func (r *Plugin) Routes(router contractsroute.Router) {
    router.Get("/blog/posts", controller.Index)
    router.Post("/blog/posts", controller.Store)
}
```

### 2. 容器绑定注册

推荐方式：在自己的 `ServiceProvider.Register` 中，以 `admin.PluginKey("xxx")` 作为绑定 key：

```go
func (r *ServiceProvider) Register(app foundation.Application) {
    app.Bind(admin.PluginKey("blog"), func(app foundation.Application) (any, error) {
        return &blog.Plugin{}, nil
    })
}
```

补充方式：也可通过 `admin.Facades.Admin().RegisterPlugin(plugin)` 手动注册。

admin 包在 Boot 阶段枚举容器中所有 `PluginKey` 类型的绑定并解析收集，自动完成菜单、页面、路由挂载。

### 3. 插件自带数据迁移（可选）

实现可选的 `contracts.Migrator` 接口，迁移由 admin 包在 Boot 阶段统一注册进宿主，**宿主的 `bootstrap/migrations.go` 不罗列**：

```go
func (r *Plugin) Migrations() []schema.Migration {
    return []schema.Migration{&migrations.M20260101000006CreateNotices{}}
}
```

## 菜单与页面

### MenuItem 结构

`MenuItem` 对应 owl-admin 前端 `IRoute`（`/menus` 端点返回的菜单树节点）：

```go
type MenuItem struct {
    Name      string     // 唯一标识
    Path      string     // 前端路由路径
    Component string
    IsHome    int
    IframeURL string
    UrlType   int        // 见下方常量
    KeepAlive *int
    IsFull    int
    IsLink    bool
    PageSign  string     // 页面标识（url_type=4 时前端按此拉取 /page_schema）
    Meta      MenuMeta   // title/icon/hide/custom_order 等
    Children  []MenuItem // 子菜单
}
```

### url_type 常量

| 常量 | 值 | 含义 |
| --- | --- | --- |
| `MenuTypeRoute` | 1 | 站内路由（前端按路径拉取 schema） |
| `MenuTypeLink` | 2 | 外链 |
| `MenuTypeIframe` | 3 | iframe |
| `MenuTypePage` | 4 | 页面（`page_schema` 下发，配合 `PageSign`） |

### 页面双轨机制

- **代码注册**：`Plugin.Pages()` 返回 `sign → amis schema`；
- **DB 覆盖**：`admin_pages` 表中同 sign 的记录可覆盖代码注册的页面，实现后台可视化改页面而无需发版。

## 注意事项

### ServiceProvider 命名冲突

goravel 框架的 ProviderRepository 按 `fmt.Sprintf("%T", provider)`（类型全名）对 ServiceProvider 去重。**插件包名不要与框架包同名**（例如包名 `schedule` 会与 `github.com/goravel/framework/services/schedule` 冲突导致整个插件被静默丢弃），建议使用不易撞名的包名，如 `scheduler`。

### 诊断插件是否挂载

若插件路由/菜单未出现，可临时编写测试遍历 `app.Bindings()` 断言 `admin.PluginKey("xxx")` 绑定是否存在。

### 路由上下文

`Plugin.Routes()` 注册的路由挂在三层中间件之内：JWT 鉴权（`admin` guard）→ 权限拦截 → 操作日志。业务控制器中可正常通过 `facades.Auth(ctx).Guard("admin")` 获取操作人。

## 参考实现

三个不同形态的完整插件见 [goravel 主项目](https://github.com/khs1001/goravel)：

| 插件 | 形态 | 特性 |
| --- | --- | --- |
| `packages/notice` | 数据库型 | 自带迁移（Migrator）、公告 CRUD |
| `packages/filemanager` | 纯文件型 | 无数据库表、safeJoin 路径穿越防护、图片预览 |
| `packages/scheduler` | 系统服务型 | 汇总 WithSchedule 定时任务、可视化手动执行 |
