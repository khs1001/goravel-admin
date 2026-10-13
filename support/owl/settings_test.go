package owl

import (
	"testing"

	"github.com/goravel/framework/contracts/http"
	"github.com/stretchr/testify/assert"
)

func TestSuccessEnvelope(t *testing.T) {
	envelope := Success("payload", "ok")

	assert.Equal(t, 0, envelope["status"])
	assert.Equal(t, "ok", envelope["msg"])
	assert.Equal(t, "payload", envelope["data"])
}

func TestSuccessDefaultsEmptyData(t *testing.T) {
	envelope := Success(nil)

	assert.Equal(t, 0, envelope["status"])
	assert.Equal(t, http.Json{}, envelope["data"])
}

func TestUnauthorizedCarriesCode(t *testing.T) {
	envelope := Unauthorized("请先登录")

	assert.Equal(t, 1, envelope["status"])
	assert.Equal(t, 401, envelope["code"])
	assert.Equal(t, 1, envelope["doNotDisplayToast"])
}

func TestLocaleOptionsPassThrough(t *testing.T) {
	options := LocaleOptions([]any{
		map[string]any{"label": "简体中文", "value": "zh_CN"},
		map[string]any{"label": "English", "value": "en"},
	})

	assert.Equal(t, []map[string]any{
		{"label": "简体中文", "value": "zh_CN"},
		{"label": "English", "value": "en"},
	}, options)
}

func TestLocaleOptionsFallback(t *testing.T) {
	options := LocaleOptions(nil)

	assert.Len(t, options, 2)
	assert.Equal(t, "zh_CN", options[0]["value"])
}
