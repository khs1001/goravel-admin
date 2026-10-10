# goravel-admin

基于 [Goravel](https://github.com/goravel/goravel) 扩展包机制的后台管理扩展，对接 [owl-admin](https://github.com/slowlyo/owl-admin) 前端（admin-views），前端零改动。

[![CI](https://github.com/khs1001/goravel-admin/actions/workflows/ci.yml/badge.svg)](https://github.com/khs1001/goravel-admin/actions)

## 功能

- **认证**：JWT（独立 `admin` guard）、登录/当前用户/个人设置、登录限流（按 IP 滑动窗口）、登录验证码（可开关）
- **RBAC**：用户/角色/权限/菜单管理，菜单可见性过滤 + 接口级权限拦截双层闭环
- **组织架构**：部门（树形、防环校验）、岗位、数据字典（含字典项子资源）
- **操作日志**：HTTP 中间件 + ORM 模型事件 + 控制器字段级 diff（old/new 对照），Queue 异步落库，密码脱敏
- **登录日志**：登录成功/失败/限流/禁用全量留痕（IP + UA），Queue 异步落库，只读审计页面
- **统一上传**：`POST /admin-api/upload`，走 `filesystems` 配置驱动（local/OSS/COS 可切换）
- **前端随包分发**：admin-views 构建产物经 `go:embed` 嵌入二进制，单二进制部署即含完整后台前端
- **迁移包内注册**：admin 自身与插件迁移均在扩展包 Boot 阶段统一注册，宿主零维护
- **二级插件系统**：实现 `contracts.Plugin` 接口，`app.Bind(PluginKey("xxx"), impl)` 即自动挂载菜单、页面与路由
- **页面全下发**：所有管理页面均为后端 amis JSON schema，支持代码注册 + `admin_pages` 表 DB 覆盖双轨

## 快速开始

```bash
go get github.com/khs1001/goravel-admin

# 注册 ServiceProvider 与 config/admin.go（官方扩展包安装协议）
go run github.com/khs1001/goravel-admin/setup
```

```bash
# 一键安装：迁移 + 数据核验 + 配置体检（幂等）
go run . artisan admin:install
```

默认超级管理员 **admin / admin**，后台入口 `http://127.0.0.1:3000/admin-assets/`。

`config/auth.go` 需增加 `admin` guard：

```go
"guards": map[string]any{
    "admin": map[string]any{
        "driver": "jwt",
    },
},
```

> ⚠️ 宿主的 `bootstrap/migrations.go` **不要**罗列扩展包/插件迁移——迁移只在扩展包内部注册（admin 包 Boot 阶段经 `Schema().Register` 统一注册），详见[迁移注册规范](docs/installation.md#迁移注册规范重要)。

## 文档

| 文档 | 内容 |
| --- | --- |
| [安装](docs/installation.md) | 自动/手动安装、auth guard、迁移注册规范、前端资源说明 |
| [配置](docs/configuration.md) | `config/admin.go` 全选项、环境变量、filesystems 关联配置 |
| [插件开发指南](docs/plugins.md) | `Plugin` 接口、容器绑定、Migrator、菜单/页面/路由、常见坑 |
| [API 接口](docs/api.md) | 三层中间件下的全部端点列表 |
| [操作日志](docs/operation-log.md) | 三层日志机制、字段级 diff 接入模式、changes 格式 |
| [部署](docs/deployment.md) | 单二进制部署、生产检查清单、升级流程 |
| [更新日志](CHANGELOG.md) | 版本历史 |

## 二级插件示例

见 [goravel 主项目](https://github.com/khs1001/goravel) 的 `packages/` 目录：

- `packages/notice`：公告管理（数据库型，自带迁移）
- `packages/filemanager`：文件管理器（纯文件型，safeJoin 防路径穿越）
- `packages/scheduler`：定时任务可视化（系统服务型）

插件接入方式与完整契约见[插件开发指南](docs/plugins.md)。

## 项目结构

```
├── contracts/        # Plugin / Admin / Migrator 契约
├── config/           # 包内配置模板（setup 安装 / vendor:publish 时发布到宿主 config/）
├── setup/            # 官方扩展包安装命令
├── console/          # artisan 命令（admin:install 一键安装）
├── http/             # 控制器与中间件（auth/permission/operation_log）
├── services/         # 业务服务（diff、登录限流、菜单/页面/权限）
├── models/           # gorm 模型
├── pages/            # 内置 amis 页面 schema + 通用列构建器
├── database/migrations/  # 包内迁移（经 Migrations() 单一注册源）
├── jobs/             # Queue 任务（操作日志落库）
└── admin-assets/     # go:embed 嵌入的前端构建产物
```

## 许可

MIT
