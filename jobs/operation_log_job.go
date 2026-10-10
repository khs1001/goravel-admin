package jobs

import (
	"encoding/json"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// OperationLogPayload 操作日志队列载荷（Args 仅支持基本类型，故整体 JSON 序列化传递）。
type OperationLogPayload struct {
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
	ResourceID string `json:"resource_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	IP         string `json:"ip"`
	Input      string `json:"input"`
	// Changes 字段级变更 JSON（{col: {old, new}}），仅 ORM updated 动作记录
	Changes string `json:"changes"`
}

// OperationLogJob 异步落库一条操作日志（经 facades.Queue() 调度，
// sync 驱动下立即执行，database 驱动下由 worker 消费）。
type OperationLogJob struct{}

// Signature 队列任务唯一签名。
func (r *OperationLogJob) Signature() string { return "admin.operation_log" }

// Handle 执行落库。
func (r *OperationLogJob) Handle(args ...any) error {
	if len(args) == 0 {
		return nil
	}

	raw, ok := args[0].(string)
	if !ok {
		return nil
	}

	var payload OperationLogPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return err
	}

	log := &models.AdminOperationLog{
		UserID:     payload.UserID,
		Username:   payload.Username,
		Action:     payload.Action,
		Resource:   payload.Resource,
		ResourceID: payload.ResourceID,
		Method:     payload.Method,
		Path:       payload.Path,
		IP:         payload.IP,
	}
	if payload.Input != "" {
		log.Input = &payload.Input
	}
	if payload.Changes != "" {
		log.Changes = &payload.Changes
	}

	if err := facades.Orm().Query().Create(log); err != nil {
		facades.Log().Warningf("[admin] save operation log failed: %v", err)

		return err
	}

	return nil
}
