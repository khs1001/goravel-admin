package jobs

import (
	"encoding/json"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// LoginLogPayload 登录日志队列载荷（与操作日志同款：整体 JSON 序列化传递）。
type LoginLogPayload struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

// LoginLogJob 异步落库一条登录日志。
type LoginLogJob struct{}

// Signature 队列任务唯一签名。
func (r *LoginLogJob) Signature() string { return "admin.login_log" }

// Handle 执行落库。
func (r *LoginLogJob) Handle(args ...any) error {
	if len(args) == 0 {
		return nil
	}

	raw, ok := args[0].(string)
	if !ok {
		return nil
	}

	var payload LoginLogPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return err
	}

	log := &models.AdminLoginLog{
		UserID:    payload.UserID,
		Username:  payload.Username,
		Status:    payload.Status,
		Message:   payload.Message,
		IP:        payload.IP,
		UserAgent: payload.UserAgent,
	}

	if err := facades.Orm().Query().Create(log); err != nil {
		facades.Log().Warningf("[admin] save login log failed: %v", err)

		return err
	}

	return nil
}
