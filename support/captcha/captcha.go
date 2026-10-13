// Package captcha wraps a base64 image captcha for the admin login page.
package captcha

import (
	"github.com/mojocn/base64Captcha"
)

// The captcha is intentionally simple: 4 digits, no external state service.
// The generated answer is kept in the library's in-memory store, keyed by the
// generated id (returned to the client as `sys_captcha`).
var instance = base64Captcha.NewCaptcha(
	base64Captcha.NewDriverDigit(40, 120, 4, 0.7, 80),
	base64Captcha.NewMemoryStore(base64Captcha.GCLimitNumber, base64Captcha.Expiration),
)

// Result is what `/captcha` returns to the frontend.
type Result struct {
	// SysCaptcha is the opaque id the client echoes back on login.
	SysCaptcha string
	// Img is a `data:image/png;base64,...` URL the frontend renders directly.
	Img string
	// Answer is the expected code; it is not returned to the client.
	Answer string
}

// Generate creates a fresh captcha image and stores its answer.
func Generate() (Result, error) {
	id, b64, answer, err := instance.Generate()
	if err != nil {
		return Result{}, err
	}

	return Result{SysCaptcha: id, Img: b64, Answer: answer}, nil
}

// Verify checks the answer for the given id and clears it (single use).
// Comparison is case-insensitive.
func Verify(id, answer string) bool {
	return instance.Verify(id, answer, true)
}
