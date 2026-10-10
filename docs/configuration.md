# 配置

`config/admin.go`（由 setup 命令自动生成）全部选项如下：

```go
config.Add("admin", map[string]any{
    "name":           config.Env("ADMIN_APP_NAME", "Goravel Admin"),
    "logo":           config.Env("ADMIN_LOGO", "/admin-assets/logo.png"),
    "default_avatar": config.Env("ADMIN_DEFAULT_AVATAR", "/admin-assets/default-avatar.png"),

    "route": map[string]any{
        "prefix":               config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
        "without_extra_routes": []string{"/dashboard"},
    },

    "auth": map[string]any{
        "login_captcha": false,
        "enable":        true,
        "permission":    false,
    },

    "layout": map[string]any{
        "title": "%title%",
        "locale_options": map[string]string{
            "en":    "English",
            "zh_CN": "简体中文",
        },
    },

    "show_development_tools": false,
})
```

## 选项说明

| 配置项 | 环境变量 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `name` | `ADMIN_APP_NAME` | `Goravel Admin` | 后台应用名称 |
| `logo` | `ADMIN_LOGO` | `/admin-assets/logo.png` | 后台 logo 地址 |
| `default_avatar` | `ADMIN_DEFAULT_AVATAR` | `/admin-assets/default-avatar.png` | 默认头像 |
| `route.prefix` | `ADMIN_ROUTE_PREFIX` | `admin-api` | 后台 API 路由前缀 |
| `route.without_extra_routes` | — | `["/dashboard"]` | 不做权限校验的额外路由 |
| `auth.login_captcha` | — | `false` | 登录验证码开关（预留，`/captcha` 端点已就位） |
| `auth.enable` | — | `true` | 是否启用认证 |
| `auth.permission` | — | `false` | 接口级权限拦截开关（false 时仅按菜单可见性过滤） |
| `layout.title` | — | `%title%` | 布局标题模板 |
| `layout.locale_options` | — | 见上 | 前端语言选项 |
| `show_development_tools` | — | `false` | 是否展示开发工具 |

## 关联配置

### filesystems（上传）

统一上传端点 `POST /admin-api/upload` 走 `filesystems` 配置驱动（local/OSS/COS 可切换）。使用 local 磁盘时建议将默认磁盘设为 `public`，保证上传后可通过 `/storage/*` 直接访问：

```env
FILESYSTEM_DISK=public
```

同时 admin 包会注册 `Route().Static("storage", "storage/app/public")` 静态路由，文件管理器插件的根目录也与之对齐。

### auth guards

`admin` guard 必须存在（JWT 驱动），详见[安装文档](installation.md)。

## 环境变量汇总

```env
# 后台应用
ADMIN_APP_NAME=Goravel Admin
ADMIN_LOGO=/admin-assets/logo.png
ADMIN_DEFAULT_AVATAR=/admin-assets/default-avatar.png
ADMIN_ROUTE_PREFIX=admin-api

# 上传磁盘（建议 public）
FILESYSTEM_DISK=public

# JWT（复用 goravel jwt 配置，如与前台共用请区分 ttl/jwt.secret 策略）
JWT_SECRET=
JWT_TTL=
```
