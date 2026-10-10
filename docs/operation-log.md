# 操作日志

操作日志由三层机制组成，全部经 Queue 异步落库，不阻塞请求。

## 组成

| 层 | 触发方式 | 记录内容 |
| --- | --- | --- |
| HTTP 中间件 | 挂在权限路由层 | 请求方法/路径/参数/响应码、操作人 |
| ORM 模型事件 | tracingObserver 监听各表 | created / deleted 事件（含 TraceColumns 关键列） |
| 控制器手动上报 | Update 控制器调用 | updated 事件的字段级 diff（changes） |

## 字段级变更详情（updated）

### 为什么在控制器层做 diff

goravel ORM 的 `Save(&model)` 模式下，`Updated` 模型事件的 `Original` 快照不可用（old 值全为 null），因此 **updated 的字段级 diff 在控制器层组装**（与 owl-admin 同款方案）；observer 只负责 created/deleted。

### 控制器接入模式

```go
func (r *AdminUserController) Update(ctx contractshttp.Context) contractshttp.Response {
    // ... 省略查找与校验

    old := *model // 值拷贝快照（保存前）

    // fill & Save ...

    // 字段对比 + 异步落日志
    changes := services.DiffFields(&old, model, []string{
        "username", "nickname", "department_id", "position_id",
        // 注意排除 password 等敏感字段
    })
    services.DispatchChangeLog(ctx, "admin_users", fmt.Sprintf("%d", model.ID), changes)

    // ...
}
```

### 关键 API

```go
// DiffFields 对比同类型模型的指定列，返回变更数组。
// 列名与模型 json tag 对齐（即 gorm column 名）。
services.DiffFields(oldModel, newModel any, columns []string) []FieldChange

type FieldChange struct {
    Field string `json:"field"`
    Old   any    `json:"old"`
    New   any    `json:"new"`
}

// DispatchChangeLog 经 Queue 异步落一条 updated 日志；
// 操作人从 ctx 的 admin guard 识别，失败降级为 system；
// changes 为空时不落库。
services.DispatchChangeLog(ctx, resource, resourceID string, changes []FieldChange)
```

### changes 存储格式

`admin_operation_logs.changes` 为 JSON 列，格式 `{"字段": {"old": ..., "new": ...}}`；后台操作日志详情抽屉以 field/old/new 三列静态表格渲染。

## 敏感信息

密码字段全程脱敏，不落日志（HTTP 中间件参数与控制器 diff 均排除）。

## 清理

建议在宿主注册定时清理任务（每天凌晨删除 30 天前日志），参考主项目 `app/schedules/kernel.go`：

```go
func Schedule() []schedule.Event {
    return []schedule.Event{
        facades.Schedule().Call(func() {
            // 删除 30 天前的 admin_operation_logs
        }).DailyAt("02:00").Name("cleanup-operation-logs"),
    }
}
```
