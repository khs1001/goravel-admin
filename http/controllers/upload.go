package controllers

import (
	neturl "net/url"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/support/owl"
)

// UploadController implements the Owl upload endpoints (`/upload_image`,
// `/upload_file`, `/upload_rich`). Files are stored on the configured disk.
type UploadController struct {
	BaseController
}

// NewUploadController creates an UploadController.
func NewUploadController() *UploadController {
	return &UploadController{}
}

// Image handles image uploads (`POST /upload_image`).
func (c *UploadController) Image(ctx http.Context) http.Response {
	url, resp := c.store(ctx, facades.Config().GetString("admin.upload.directory.image", "admin/images"))
	if resp != nil {
		return resp
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"value": url}))
}

// File handles generic file uploads (`POST /upload_file`).
func (c *UploadController) File(ctx http.Context) http.Response {
	url, resp := c.store(ctx, facades.Config().GetString("admin.upload.directory.file", "admin/files"))
	if resp != nil {
		return resp
	}

	return ctx.Response().Json(http.StatusOK, owl.Success(http.Json{"value": url}))
}

// Rich handles rich-text uploads (`POST /upload_rich`). The response follows the
// wangEditor protocol (`{errno, data:{url}}`).
func (c *UploadController) Rich(ctx http.Context) http.Response {
	url, resp := c.store(ctx, facades.Config().GetString("admin.upload.directory.rich", "admin/rich"))
	if resp != nil {
		return resp
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"errno":             0,
		"status":            0,
		"msg":               "",
		"doNotDisplayToast": 0,
		"data":              http.Json{"url": url},
	})
}

// store saves the request file into dir on the configured disk and returns its
// public URL, or an error response.
func (c *UploadController) store(ctx http.Context, dir string) (string, http.Response) {
	disk := facades.Config().GetString("admin.upload.disk", "public")

	file, err := ctx.Request().File("file")
	if err != nil || file == nil {
		return "", ctx.Response().Json(http.StatusOK, owl.Fail("上传文件不能为空"))
	}

	path, err := facades.Storage().Disk(disk).PutFile(dir, file)
	if err != nil {
		return "", ctx.Response().Json(http.StatusOK, owl.Fail(err.Error()))
	}

	return relativeURL(facades.Storage().Disk(disk).Url(path)), nil
}

// relativeURL keeps only the path of a stored file URL so it works regardless
// of the configured host/port (the SPA is served from the same origin).
func relativeURL(url string) string {
	if parsed, err := neturl.Parse(url); err == nil && parsed.Path != "" {
		return parsed.Path
	}

	return url
}
