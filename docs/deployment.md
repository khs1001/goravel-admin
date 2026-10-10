# 部署

## 单二进制部署

前端构建产物经 `go:embed` 嵌入 admin 包，**编译出的宿主二进制即包含完整后台前端**，无需 nginx 托管静态文件或另行部署前端：

```bash
# 宿主项目
go build -o app .
./app
```

后台入口：`http://<host>:<port>/admin-assets/`。

## 运行时依赖

| 依赖 | 说明 |
| --- | --- |
| MySQL | 主数据库 |
| storage/app/public | 上传文件目录（local 磁盘），需保证可写且持久化 |
| Queue | 操作日志异步落库依赖队列；sync 驱动亦可运行，生产建议独立队列驱动 |

## 生产检查清单

- [ ] 修改默认管理员密码（admin / admin）
- [ ] `config admin.auth.permission` 按需开启接口级权限拦截
- [ ] `FILESYSTEM_DISK=public`（local 磁盘场景），确认 `/storage/*` 可访问
- [ ] `bootstrap/migrations.go` 只含应用自身迁移（扩展包迁移由包内自动注册）
- [ ] 操作日志定时清理任务已注册（参考 `app/schedules/kernel.go`，配合 `WithSchedule` builder）
- [ ] JWT_SECRET 已配置且不与前台共用弱密钥
- [ ] 定时任务经 `go Run()` 启动并配置信号优雅停机

## 容器化注意

- 二进制自包含前端资源，镜像内只需应用本体 + 可写存储卷；
- 持久化 `storage/` 目录（上传文件）；
- 首次启动前执行 `go run . artisan migrate`（或以初始化容器/job 执行）。

## 升级扩展包

```bash
go get github.com/khs1001/goravel-admin@latest
go mod tidy
go run . artisan migrate   # 新版本如带迁移会自动注册
go build -o app .          # 重新编译即携带最新前端产物
```

版本历史见 [CHANGELOG](../CHANGELOG.md)。
