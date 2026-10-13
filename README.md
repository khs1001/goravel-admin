# goravel-admin

对标 **RuoYi** 的 Goravel 后台扩展包：**RuoYi 数据模型 + Owl Admin 接口契约**。

包本身是一个独立的 Go module（`github.com/khs1001/goravel-admin`），可通过
`./artisan package:install` 装进任意 Goravel 宿主应用。前端复用仓库内的
`admin-views/`（Owl Admin，React + amis），后端用 Go 实现其接口，因此前端零改动。

## 设计要点

- **数据模型沿用 RuoYi**：表名与列名走 `sys_*` 风格
  （`sys_user` / `sys_role` / `sys_menu` / `sys_user_role` / `sys_role_menu` / `sys_config`）。
- **权限模型沿用 RuoYi**：菜单树即权限树，`sys_menu.menu_type`（M 目录 / C 菜单 / F 按钮）
  与 `perms`（`模块:功能:操作`，如 `system:user:add`）；角色经 `sys_role_menu` 绑菜单。
- **接口对齐 Owl Admin**：统一响应信封 `{status,msg,doNotDisplayToast,data}`，
  固定端点 `/login` `/logout` `/current-user` `/menus` `/captcha` `/_settings` `/page_schema`；
  业务数据用 RuoYi 动作端点（`/system/user/list`、`/system/user/{id}`、`POST/PUT/DELETE /system/user`）。
- **适配集中在一处**：`support/owl` 是唯一"记住 Owl 契约"的地方（响应信封、菜单路由树、设置拼装）。
- **页面即数据**：系统页面是内嵌的 amis JSON（`resources/pages/*.json`，`go:embed`），
  经 `/page_schema?sign=` 下发；前端不解析业务字段，纯渲染。
- **纯运行时接线**：路由、命令、迁移、种子、静态资源都在 `ServiceProvider.Boot()` 里注册，
  `package:install` / `package:uninstall` 对称且幂等，不改宿主源码。

## 覆盖范围（v1）

| 能力 | 说明 |
|---|---|
| 登录 / 退出 / 验证码 | JWT（复用宿主 `config/jwt.go`），`base64Captcha` 图形验证码（默认开启） |
| 当前用户 / 动态菜单 | `/current-user`、`/menus`（按角色→菜单权限过滤） |
| 用户管理 | `sys_user` + `sys_user_role` + `sys_user_post` CRUD，部门/岗位/角色分配 |
| 角色管理 | `sys_role` + `sys_role_menu` + `sys_role_dept` CRUD，菜单权限树 + 数据范围 |
| 菜单管理 | `sys_menu` 树 CRUD（M/C/F + perms） |
| 部门管理 | `sys_dept` 树 CRUD（组织架构） |
| 岗位管理 | `sys_post` CRUD |
| 字典管理 | `sys_dict_type` + `sys_dict_data` CRUD |
| 参数设置 | `sys_config` CRUD（同时承载 Owl `/_settings` 的持久化） |
| 数据权限 | 角色 `data_scope`（全部/自定义/本部门/本部门及以下/仅本人），作用于用户与部门列表 |
| 文件上传 | `/upload_image`、`/upload_file`、`/upload_rich`（wangEditor 协议），存储于 `public` 磁盘并经 `/storage/*` 访问 |
| 个人设置 | 顶栏「个人设置」页，维护个人资料并可修改密码 |
| 操作日志 | `sys_oper_log`：写操作中间件自动记录（模块/类型/人员/IP/参数/状态/耗时），页面查询/详情/删除/清空 |
| 登录日志 | `sys_logininfor`：登录成功/失败/退出自动记录，页面查询/删除/清空 |
| 插件系统 | 绑定 `admin.plugin.<name>` 即可被发现；按能力接口贡献菜单/路由/页面/设置/资产/命令/迁移/钩子；默认禁用，扩展管理页与 `admin:plugin:*` 命令启停 |
| 后台命令 | `admin:create-user`、`admin:reset-password` |

暂不包含（后续迭代）：通知公告、定时任务、代码生成器。

## 安装

### 已发布的包

```shell
./artisan package:install github.com/khs1001/goravel-admin
```

安装器会：

1. `go get` 该包；
2. 运行包内 `setup/setup.go`：在 `bootstrap/providers.go` 注册
   `&admin.ServiceProvider{}`（非 bootstrap 脚手架则写入 `app.go`）；
3. 把包内 `config/admin.go` 模板发布到宿主 `config/admin.go`；
4. `go mod tidy`。

### 本地开发（go.work）

包与宿主放在同一个 `go.work` 里时，宿主 `require` 一个**真实存在的版本**
（工作区会把源码指向本地目录，无需 `replace`）：

```shell
# go.work
use (
    ./goravel
    ./goravel-admin
)
```

```go
// goravel/go.mod —— 版本只需真实存在，实际源码由 go.work 覆盖
require github.com/khs1001/goravel-admin v0.1.14
```

然后照常执行 `./artisan package:install github.com/khs1001/goravel-admin`。
（`package:install` 只接受模块路径，不支持 `./relative` 路径。）

宿主需先具备：`Auth`（JWT）、`Orm`、`Schema`、`Seeder`、`Route`、`Cache`、`Hash`。
`admin` 守卫（`auth.guards.admin` / `auth.providers.admin`）由包在运行时按需注入，
不会覆盖宿主既有配置。

## 初始化数据

```shell
./artisan migrate --seed
```

初始账号：`admin` / `admin123`。

## 访问

- 后台入口：`/admin`
- 接口前缀：`/admin-api`
- 静态资源：`/admin-assets`

## 前端（admin-views）

仓库内的 `admin-views/` 是 Owl Admin 前端（React + amis 6.13），已构建产物位于
`admin-views/dist`，由包在运行时托管，无需宿主额外配置。

重新构建：

```shell
cd admin-views
pnpm install
pnpm build      # 产物输出到 dist/（base=/admin-assets、接口前缀 /admin-api）
```

约定（与已构建产物一致）：SPA 入口 `/admin`（HashRouter，无需服务端 SPA 回退）、
接口前缀 `/admin-api`、资源前缀 `/admin-assets`。若要改前缀，需同时调整
`admin-views/.env` 的 `VITE_BASE_URL` / `VITE_API_PREFIX`、前端 `index.html` 里的
`window.$adminApiPrefix`，以及 `config/admin.go` 的 `assets` / `route.prefix`。

## 配置（`config/admin.go`）

```go
"name":   "Goravel Admin",          // 应用名
"logo":   "/admin-assets/logo.png", // 登录页 logo
"path":   "admin",                  // SPA 入口路径
"assets": "admin-assets",           // 静态资源路径
"route":  { "prefix": "admin-api" },// 接口前缀
"auth":   {
    "enable": true,             // 是否鉴权
    "permission": true,         // 是否做权限校验
    "login_captcha": true,      // 是否开启验证码
    "token_expiration": 0,      // 分钟，0 用宿主 jwt.ttl
},
"layout": { /* 标题、页脚、header 按钮、语言选项 */ },
```

## 门面

```go
import "github.com/khs1001/goravel-admin/facades"

facades.Admin().Name()        // "Goravel Admin"
facades.Admin().RoutePrefix() // "admin-api"
```

## 插件系统

其它 Goravel 扩展包可以扩展 goravel-admin 的功能，**无需依赖 goravel-admin**。

### 机制

1. 插件包提供自己的 `ServiceProvider`（由 `package:install` 装入宿主）。
2. 在 `Register` 里把插件实例绑定到约定键 `admin.plugin.<name>`：

```go
app.Bind("admin.plugin.blog", func(app foundation.Application) (any, error) {
    return &BlogPlugin{}, nil
})
```

3. goravel-admin 在 `Boot` 扫描 `app.Bindings()`，解析该前缀的键并断言**能力接口**。
4. 插件在 `sys_extension` 中登记，**默认禁用**，需在「系统管理 → 扩展管理」页或命令启用。
5. 禁用 = 其运行时贡献（菜单/权限/页面/路由/资产/nav）立即隐藏（路由常驻但经守卫返回 404，无需重启）。

### 能力接口（可选，实现即加载）

基础契约 `Plugin{ Name() string }`（必填）。其余按需实现，签名只用 **Goravel 类型或纯类型**（插件可零依赖）：

| 接口 | 作用 |
|---|---|
| `MetadataProvider` | `Version/Description/Author`，供扩展管理页展示 |
| `MenuProvider` | `Menus() []map[string]any`：菜单节点（title/path/parent/order/icon/perms/component/menuType） |
| `RouteProvider` | `Routes(router route.Router)`：子路由（前缀 `/admin-api/ext/<name>`，已套鉴权/日志/守卫） |
| `PageProvider` | `Pages() map[string]map[string]any`：amis 页面（sign 命名空间 `<name>.<short>`） |
| `SettingProvider` | `Settings() map[string]any`：默认设置（宿主 config ＜ 插件默认 ＜ DB `sys_config`） |
| `AssetProvider` / `NavProvider` | 注入 `/_settings` 的 `assets`(js/css) 与 `nav` |
| `CommandProvider` | `Commands() []console.Command` |
| `MigrationProvider` / `SeederProvider` | 迁移 / 种子（安装期，不受启停门控） |
| 钩子 | `UserAuthenticatedHook` / `LoginHook` / `LogoutHook` / `OperationLoggedHook` / `MenusBuiltHook` / `SettingsBuiltHook` |

### 命令

```shell
./artisan admin:plugin:enable  --name=blog
./artisan admin:plugin:disable --name=blog
./artisan admin:plugin:clean          # 清理已卸载插件的记录
```

### 示例

`examples/plugin/` 是一个**零依赖**示例插件（不 import goravel-admin），演示菜单（挂在「系统管理」下）、页面、设置、路由与命令。把它注册进宿主 `bootstrap/providers.go` 即可试用。

## 测试与验证

```shell
go build ./...
go vet ./...
go test ./...
```

`support/owl` 的响应信封、菜单树、语言选项等为纯单测，默认全绿。
涉及数据库的端到端验证需先装进宿主应用。

## 目录结构

```
goravel-admin/
├── admin.go / service_provider.go   # 主服务与运行时接线
├── contracts/ facades/ config/ setup/
├── routes/admin.go                  # API + SPA + 静态资源
├── models/                          # sys_user/sys_role/sys_menu/sys_dept/sys_post/sys_dict_*/sys_config 等
├── http/{controllers,middleware}/
├── services/                        # menu（权限）、config（设置）、data_scope（数据权限）
├── support/{owl,captcha,current}/   # Owl 适配、验证码、请求上下文用户
├── resources/pages/*.json           # 内嵌 amis 页面
├── database/{migrations,seeders}/
├── console/commands/                # admin:create-user / admin:reset-password
└── admin-views/                     # Owl Admin 前端（含已构建 dist）
```
