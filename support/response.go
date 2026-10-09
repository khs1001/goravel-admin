package support

// Envelope 严格兼容 owl-admin 的响应包装
// （见 owl-admin src/Support/Cores/JsonResponse.php）：
//
//	成功: {"status": 0, "msg": "...", "doNotDisplayToast": 0, "data": ...}
//	失败: {"status": 1, "msg": "...", "doNotDisplayToast": 0, "data": ...}
//	未登录: {"status": 1, "code": 401, "msg": "...", "doNotDisplayToast": 1, "data": null}
//
// 前端拦截器（admin-views src/service/request/instance.ts）依赖：
//   - status === 0 表示业务成功，取 data
//   - code === 401 表示 token 失效，跳转登录页
//   - doNotDisplayToast 为 1 时不弹提示
type Envelope struct {
	Status            int    `json:"status"`
	Msg               string `json:"msg"`
	DoNotDisplayToast int    `json:"doNotDisplayToast"`
	Code              int    `json:"code,omitempty"`
	Data              any    `json:"data"`
}

// Success 成功响应，data 为业务数据（nil 时输出 {}，对齐 PHP success(null) 的 (object) 转换）。
func Success(data any, msg string) Envelope {
	if data == nil {
		data = struct{}{}
	}

	return Envelope{Status: 0, Msg: msg, Data: data}
}

// SuccessMessage 成功响应但无业务数据（data 为空数组，对齐 PHP success([])）。
func SuccessMessage(msg string) Envelope {
	return Envelope{Status: 0, Msg: msg, Data: []any{}}
}

// Fail 失败响应（status=1）。
func Fail(msg string) Envelope {
	return Envelope{Status: 1, Msg: msg, Data: nil}
}

// Unauthorized 未登录/token 失效响应（code=401，前端据此跳转登录页）。
func Unauthorized(msg string) Envelope {
	return Envelope{Status: 1, Code: 401, Msg: msg, DoNotDisplayToast: 1}
}
