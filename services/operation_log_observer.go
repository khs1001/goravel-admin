package services

import (
	"encoding/json"
	"fmt"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/jobs"
)

// OperationLogObserver ORM 模型事件观察者：记录后台资源表的
// created / updated / deleted 变更，经 Queue 异步落库。
// 注册方式见 ServiceProvider.Boot：facades.Orm().Observe(model, observer)。
type OperationLogObserver struct {
	Resource string
}

// Created 创建后记录。
func (r OperationLogObserver) Created(event contractsorm.Event) error {
	r.record("created", event)

	return nil
}

// Updated 更新后记录。
func (r OperationLogObserver) Updated(event contractsorm.Event) error {
	r.record("updated", event)

	return nil
}

// Deleted 删除后记录（含批量删除，逐行触发）。
func (r OperationLogObserver) Deleted(event contractsorm.Event) error {
	r.record("deleted", event)

	return nil
}

// ForceDeleted 强删后记录。
func (r OperationLogObserver) ForceDeleted(event contractsorm.Event) error {
	r.record("deleted", event)

	return nil
}

// record 组装载荷并入队（异常仅告警，不影响业务写入）。
func (r OperationLogObserver) record(action string, event contractsorm.Event) {
	defer func() {
		if err := recover(); err != nil {
			facades.Log().Warningf("[admin] operation log observer panic: %v", err)
		}
	}()

	id := fmt.Sprintf("%v", event.GetAttribute("id"))
	if id == "<nil>" || id == "0" {
		return
	}

	payload := jobs.OperationLogPayload{
		Action:     action,
		Resource:   r.Resource,
		ResourceID: id,
	}

	// 操作人：orm 事件栈里拿不到 HTTP 上下文，统一记为 system；
	// 请求级操作人由 OperationLog 中间件记录。
	payload.Username = "system"

	encoded, err := json.Marshal(payload)
	if err != nil {
		facades.Log().Warningf("[admin] marshal orm log failed: %v", err)

		return
	}

	if err := facades.Queue().Job(&jobs.OperationLogJob{}, []contractsqueue.Arg{
		{Value: string(encoded), Type: "string"},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("[admin] dispatch orm log failed: %v", err)
	}
}
