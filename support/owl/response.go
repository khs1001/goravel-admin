// Package owl implements the Owl Admin wire contract: the response envelope,
// the dynamic-menu route tree and the settings payload consumed by the shipped
// `admin-views` frontend.
//
// It is the single place in this package that "remembers" the Owl contract.
// Everything else is organised around the RuoYi data model.
package owl

import (
	"github.com/goravel/framework/contracts/http"
)

// Success builds a successful envelope: {status:0, msg, doNotDisplayToast, data}.
func Success(data any, msg ...string) http.Json {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if data == nil {
		data = http.Json{}
	}

	return http.Json{
		"status":            0,
		"msg":               message,
		"doNotDisplayToast": 0,
		"data":              data,
	}
}

// Fail builds a failed envelope: {status:1, msg, doNotDisplayToast:0, data:null}.
func Fail(msg string) http.Json {
	return http.Json{
		"status":            1,
		"msg":               msg,
		"doNotDisplayToast": 0,
		"data":              nil,
	}
}

// FailSilent behaves like Fail but suppresses the frontend toast.
func FailSilent(msg string) http.Json {
	return http.Json{
		"status":            1,
		"msg":               msg,
		"doNotDisplayToast": 1,
		"data":              nil,
	}
}

// Unauthorized is the envelope the frontend uses to detect an expired token.
// The frontend switches to the login page when `code == 401`.
func Unauthorized(msg string) http.Json {
	return http.Json{
		"status":            1,
		"msg":               msg,
		"code":              401,
		"doNotDisplayToast": 1,
		"data":              nil,
	}
}
