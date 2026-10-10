# 更新日志

所有 notable changes 记录于此。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)。

## [v0.1.7] — 控制器侧字段级 diff

- 新增 `services.DiffFields`（json-tag 序列化对比）与 `services.DispatchChangeLog`（Queue 异步落库、自动识别操作人）
- 全部 Update 控制器接入字段级变更记录（users 排除 password）
- 背景：goravel ORM `Save(&model)` 模式下 `Updated` 事件的 Original 快照不可用，diff 迁移至控制器层组装；observer 保留 created/deleted

## [v0.1.6] — 操作日志变更详情 + 通用列构建器

- `admin_operation_logs` 新增 `changes` JSON 列（迁移 20260101000007）
- 操作日志详情抽屉：field/old/new 三列静态表格渲染字段级变更对照
- 新增 `pages/columns.go` 9 个通用 amis 列构建器，全部列表页重构采用
- 品牌默认值调整

## [v0.1.5] — 前端产物随包分发

- 前端构建产物（owl-admin admin-views）经 `go:embed all:admin-assets` 嵌入包内，导出 `AdminAssets() http.FileSystem`
- 路由改为 `Route().StaticFS("admin-assets", AdminAssets())`；宿主 `public/` 不再携带后台前端资源
- 单二进制部署即包含完整后台前端

## [v0.1.4] — 迁移注册规范

- 包根新增 `Migrations()`（admin 自身 6 条迁移的单一注册源）
- `contracts` 新增可选 `Migrator` 接口：插件自带迁移由 admin 包 Boot 阶段统一注册
- `ServiceProvider.Boot` 内 `registerMigrations()`；宿主 `bootstrap/migrations.go` 不再罗列扩展包迁移

## [v0.1.3] — 登录限流

- 登录接口按 IP 滑动窗口限流（cache 实现）

## [v0.1.2] — 单元测试

- 响应契约、RBAC 匹配、辅助函数、日志脱敏单元测试

## [v0.1.1] — 移除包内 config

- config 归宿主应用 `config/admin.go`（由 setup 命令生成），对齐官方扩展包协议

## [v0.1.1 之前] — 初始版本

- 认证（JWT admin guard）、RBAC（用户/角色/权限/菜单 + 双层闭环）、组织架构（部门/岗位/字典）、操作日志（中间件 + 模型事件 + Queue）、统一上传、二级插件系统、页面全下发（代码注册 + DB 覆盖）
