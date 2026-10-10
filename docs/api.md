# API 接口

所有接口默认挂在 `/admin-api` 前缀下（`config admin.route.prefix` 可改）。

中间件分三层：

1. **基座路由**：免鉴权；
2. **鉴权路由**：JWT（`admin` guard）；
3. **权限路由**：RBAC 接口级拦截 + 操作日志（`auth.permission` 开启时生效）。

## 基座路由（免鉴权）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/login` | 登录 |
| GET | `/logout` | 退出登录 |
| GET | `/captcha` | 登录验证码（预留端点） |
| GET | `/no-content` | 空页面占位 |
| GET | `/_settings` | 读取站点设置 |

## 鉴权路由

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/_settings` | 保存站点设置 |
| GET | `/menus` | 当前用户菜单树（代码菜单 + DB 菜单 + 插件菜单，按权限过滤） |
| GET | `/current-user` | 当前登录用户 |
| GET | `/user_setting` | 读取个人设置 |
| PUT | `/user_setting` | 保存个人设置 |
| GET | `/page_schema` | 按 sign 拉取页面 amis schema |
| POST | `/upload` | 统一上传（走 filesystems 驱动） |

## 权限路由

### 仪表盘

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/dashboard` | 仪表盘页面 schema |
| GET | `/dashboard_stats` | 仪表盘统计数据 |

### 系统管理（对齐 owl-admin /system 资源路由）

| 资源 | 端点模式 | 额外端点 |
| --- | --- | --- |
| 用户 `admin_users` | GET / POST / GET `{id}/edit` / PUT `{id}` / DELETE `{id}`、DELETE（批量） | — |
| 角色 `admin_roles` | 同上 | GET `options`、POST `save_permissions` |
| 权限 `admin_permissions` | 同上 | GET `tree` |
| 菜单 `admin_menus` | 同上 | GET `tree`、POST `save_order` |
| 部门 `admin_departments` | 同上 | GET `tree` |
| 岗位 `admin_positions` | 同上 | GET `options` |
| 字典 `admin_dictionaries` | 同上 | — |
| 字典项 `admin_dictionary_items` | 同上 | — |
| 操作日志 `admin_operation_logs` | GET / GET `{id}/edit`（只读） | — |

另有 `GET /system` 返回系统管理入口页 schema。

资源 DELETE 端点同时支持单条（`{id}`）与批量（body 传 id 数组）。

### 插件路由

各插件在 `Plugin.Routes()` 中注册的业务路由，同样位于权限层之下。内置插件示例见主项目 `packages/notice`、`packages/filemanager`、`packages/scheduler`。

## 静态资源

| 路径 | 来源 | 说明 |
| --- | --- | --- |
| `/admin-assets/*` | 包内 `go:embed` 嵌入的 admin-views 构建产物 | 后台前端（单二进制自带） |
| `/storage/*` | `storage/app/public` | 上传文件访问（public 磁盘） |
