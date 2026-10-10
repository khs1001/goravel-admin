# 更新日志

所有 notable changes 记录于此。格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)。

## [v0.1.11] — 修复

- 修复 `LoginLogJob` 未注册进 Queue JobStorer：未注册 signature 的任务被 worker 认领后永久卡在 reserved 状态（不消费、不进 failed_jobs、无错误日志）。新增队列任务必须 `facades.Queue().Register`

## [v0.1.10] — 登录日志

- 新增 `admin_login_logs` 表（迁移 20260101000008，含"登录日志"菜单与查看权限的幂等 seed）
- `LoginLogJob`（Queue 异步）：登录成功/密码错误/限流锁定/账号禁用全量留痕，记录 IP 与 User-Agent
- 只读列表页 `/system/admin_login_logs`（用户名/IP/状态筛选 + 详情抽屉），菜单挂在"系统管理"下
- 与操作日志同款异步模式，队列派发失败不影响登录主流程

## [v0.1.8] — 登录验证码

- 实现 `GET /captcha`：`base64Captcha` 数字图形验证码，返回 `{sys_captcha, captcha_img}`（data URI，前端直接渲染）
- 答案经 `facades.Cache` 存取（key 由答案 md5 派生、不下发），5 分钟过期，`Cache.Pull` 一次性消费；走宿主 cache 驱动，支持多实例部署
- `Login` 在开关（`admin.auth.login_captcha`）开启时强制校验验证码，失败不计入登录限流
- 前端零改动：admin-views 登录页原生支持（`settings.login_captcha` 触发拉取，提交带 `captcha` + `sys_captcha`）

## [v0.1.9] — 测试支持

- `GenerateLoginCaptcha` 同时返回答案，供宿主 feature 测试走真实 cache 门面验证全链路（生成→校验→一次性消费→拒绝）

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
