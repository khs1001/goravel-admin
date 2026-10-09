# goravel-admin

基于 [Goravel](https://github.com/goravel/goravel) 扩展包机制的后台管理扩展，对接 [owl-admin](https://github.com/slowlyo/owl-admin) 前端（admin-views），前端零改动。

## 功能

- **认证**：JWT（独立 `admin` guard）、登录/当前用户/个人设置
- **RBAC**：用户/角色/权限/菜单管理，菜单可见性过滤 + 接口级权限拦截双层闭环
- **组织架构**：部门（树形、防环校验）、岗位、数据字典（含字典项子资源）
- **操作日志**：HTTP 中间件 + ORM 模型事件 + Queue 异步落库，密码脱敏
- **统一上传**：`POST /admin-api/upload`，走 `filesystems` 配置驱动（local/OSS/COS 可切换）
- **二级插件系统**：实现 `contracts.Plugin` 接口，`app.Bind(PluginKey("xxx"), impl)` 即自动挂载菜单、页面与路由
- **页面全下发**：所有管理页面均为后端 amis JSON schema，支持代码注册 + `admin_pages` 表 DB 覆盖双轨

## 安装

```bash
go get github.com/khs1001/goravel-admin
```

```bash
# 注册 ServiceProvider 与 config/admin.go（官方扩展包安装协议）
go run github.com/khs1001/goravel-admin/setup
```

### 手动安装

`bootstrap/providers.go`：

```go
import goraveladmin "github.com/khs1001/goravel-admin"

func Providers() []foundation.ServiceProvider {
    return []foundation.ServiceProvider{
        // ...
        &goraveladmin.ServiceProvider{},
    }
}
```

`bootstrap/migrations.go`：

```go
import adminmigrations "github.com/khs1001/goravel-admin/database/migrations"

func Migrations() []schema.Migration {
    return []schema.Migration{
        // ...
        &adminmigrations.M20260101000001CreateAdminTables{},
        &adminmigrations.M20260101000002SeedRbac{},
        &adminmigrations.M20260101000003CreateOrgTables{},
        &adminmigrations.M20260101000004SeedOrgRbac{},
        &adminmigrations.M20260101000005CreateOperationLogs{},
    }
}
```

`config/auth.go` 增加 guard：

```go
"guards": map[string]any{
    "admin": map[string]any{
        "driver": "jwt",
    },
},
```

然后 `go run . artisan migrate`，默认超级管理员 **admin / admin**。

## 前端

将 owl-admin admin-views 构建产物拷贝到主项目 `public/admin-assets/`，访问 `http://localhost:3000/admin-assets/`。

## 二级插件示例

见 [goravel 主项目](https://github.com/khs1001/goravel) 中的 `packages/notice`（公告管理）：实现 `Plugin` 接口并在自己的 `ServiceProvider.Register` 中绑定即可，admin 包 Boot 时自动发现，宿主零修改。

## 许可

MIT
