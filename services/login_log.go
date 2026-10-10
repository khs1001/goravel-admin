package services

import (
	"encoding/json"

	contractshttp "github.com/goravel/framework/contracts/http"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/jobs"
)

// DispatchLoginLog 异步落一条登录日志（经 Queue，与操作日志同款异步模式）。
// 记录字段缺失不影响登录主流程；队列派发失败仅告警。
func DispatchLoginLog(ctx contractshttp.Context, userID int, username, status, message string) {
	payload := jobs.LoginLogPayload{
		UserID:   userID,
		Username: username,
		Status:   status,
		Message:  message,
	}

	if ctx != nil {
		payload.IP = ctx.Request().Ip()
		payload.UserAgent = ctx.Request().Header("User-Agent", "")
		if payload.UserAgent == "" {
			payload.UserAgent = ctx.Request().Header("user-agent", "")
		}
		if len(payload.UserAgent) > 255 {
			payload.UserAgent = payload.UserAgent[:255]
		}
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		facades.Log().Warningf("[admin] marshal login log failed: %v", err)

		return
	}

	if err := facades.Queue().Job(&jobs.LoginLogJob{}, []contractsqueue.Arg{
		{Value: string(encoded), Type: "string"},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("[admin] dispatch login log failed: %v", err)
	}
}
