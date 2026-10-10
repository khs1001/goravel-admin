# 安装

## 环境要求

- Go 1.22+
- [Goravel](https://github.com/goravel/goravel) v1.18.0+
- MySQL（其他数据库未验证）

## 自动安装（推荐）

```bash
go get github.com/khs1001/goravel-admin

# 官方扩展包安装协议：自动注册 ServiceProvider 并发布 config/admin.go
go run github.com/khs1001/goravel-admin/setup
```

也可用官方 vendor:publish 单独发布配置（模板位于包内 `setup/config/admin.go`）：

```bash
go run . artisan vendor:publish --package=github.com/khs1001/goravel-admin
```

> 配置已存在时不会静默覆盖，需加 `--force`。

## 手动安装

`bootstrap/providers.go` 注册 ServiceProvider：

```go
import goraveladmin "github.com/khs1001/goravel-admin"

func Providers() []foundation.ServiceProvider {
    return []foundation.ServiceProvider{
        // ...
        &goraveladmin.ServiceProvider{},
    }
}
```

`config/auth.go` 增加 `admin` guard（JWT 独立守卫，与前台用户隔离）：

```go
"guards": map[string]any{
    "admin": map[string]any{
        "driver": "jwt",
    },
},
```

## 迁移注册规范（重要）

> **数据迁移只在扩展包内部注册，宿主应用不罗列。**

- admin 包自身的迁移来自包根的 `admin.Migrations()`（单一注册源），由 `ServiceProvider.Boot` 经 `Schema().Register` 统一注册；
- 二级插件迁移来自各插件的可选 `contracts.Migrator` 实现，同样在 admin 包 Boot 阶段自动收集注册。

因此宿主的 `bootstrap/migrations.go` **不需要**出现任何 admin/插件迁移，只保留应用自身的迁移（如 jobs 表）：

```go
func Migrations() []schema.Migration {
    return []schema.Migration{
        &migrations.M20240101000001CreateJobsTable{},
        // 只放应用自己的迁移
    }
}
```

`Schema().Register` 是追加式、按迁移签名去重的，重复注册不会产生副作用。

## 初始化

一键安装（迁移 + 数据核验 + 配置体检，幂等可重复执行）：

```bash
go run . artisan admin:install
```

或分步执行：

```bash
go run . artisan migrate
go run .
```

默认超级管理员：**admin / admin**，登录后请立即修改密码。

后台入口：`http://127.0.0.1:3000/admin-assets/`（前端构建产物随包内嵌，无需另行部署）。

## 前端资源说明

前端构建产物（owl-admin admin-views）经 `go:embed` 嵌入扩展包二进制，由 admin 包在注册路由时通过 `Route().StaticFS("admin-assets", AdminAssets())` 提供服务：

- 宿主应用的 `public/` **不携带、不注册**后台前端资源；
- 升级前端 = 替换包内 `admin-assets/` 目录 → 发新版本 → 宿主 `go get` 升级即可；
- 单二进制部署即包含完整后台前端。
