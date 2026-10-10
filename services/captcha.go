package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/mojocn/base64Captcha"
)

// 登录验证码（对齐 owl-admin 契约）：
//   - GET /captcha 下发 {sys_captcha: true, captcha_img: "data:image/png;base64,..."}；
//   - 答案经 Cache 存取，key 由答案派生（md5），不向前端下发 key；
//   - 登录时按用户输入派生 key 做 Cache.Pull 一次性消费，命中即校验通过。

const (
	captchaTTL       = 5 * time.Minute
	captchaWidth     = 160
	captchaHeight    = 48
	captchaLength    = 4
	captchaMaxSkew   = 0.7
	captchaDotCount  = 80
	captchaCacheKeyT = "admin:captcha:"
)

// NormalizeCaptchaInput 归一化用户输入（去空白、小写）。
func NormalizeCaptchaInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

// CaptchaCacheKey 由答案派生的缓存 key（纯函数，供登录校验复用）。
func CaptchaCacheKey(normalizedAnswer string) string {
	sum := md5.Sum([]byte(normalizedAnswer))

	return captchaCacheKeyT + hex.EncodeToString(sum[:])
}

// CaptchaEnabled 验证码开关（config admin.auth.login_captcha）。
func CaptchaEnabled() bool {
	return facades.Config().GetBool("admin.auth.login_captcha", false)
}

// cacheCaptchaStore 基于 facades.Cache 的验证码存储：
// key 由答案派生（不含 id），登录校验端无需持有 id 即可定位。
// 相比库默认的进程内存储，走宿主 cache 驱动，支持多实例部署。
type cacheCaptchaStore struct{}

func (cacheCaptchaStore) Set(id, value string) error {
	answer := NormalizeCaptchaInput(value)
	if answer == "" {
		return fmt.Errorf("captcha answer is empty")
	}

	return facades.Cache().Put(CaptchaCacheKey(answer), answer, captchaTTL)
}

func (cacheCaptchaStore) Get(id string, clear bool) string {
	return ""
}

func (cacheCaptchaStore) Verify(id, answer string, clear bool) bool {
	return VerifyLoginCaptcha(answer)
}

// GenerateLoginCaptcha 生成一张数字验证码图片：答案入 cache，
// 返回 data URI 与答案（答案仅供测试/内部使用，不下发前端）。
func GenerateLoginCaptcha() (img string, answer string, err error) {
	driver := base64Captcha.NewDriverDigit(captchaHeight, captchaWidth, captchaLength, captchaMaxSkew, captchaDotCount)
	captcha := base64Captcha.NewCaptcha(driver, cacheCaptchaStore{})

	_, b64, ans, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}

	return b64, NormalizeCaptchaInput(ans), nil
}

// VerifyLoginCaptcha 校验并一次性消费验证码（Cache.Pull 原子取删）。
func VerifyLoginCaptcha(input string) bool {
	answer := NormalizeCaptchaInput(input)
	if answer == "" {
		return false
	}

	value := facades.Cache().Pull(CaptchaCacheKey(answer))

	return value != nil && fmt.Sprint(value) == answer
}
