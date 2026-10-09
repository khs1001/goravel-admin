package controllers

import (
	"path/filepath"
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/support"
)

// UploadController 统一上传接口（对齐 owl-admin 前端 amis receiver 响应格式）。
//
// 存储走 config/filesystems.go 默认磁盘（driver 可切换 local/cloud）；
// local/public 磁盘的文件由 /storage 静态路由访问，云盘由 Storage.Url() 返回远端地址。
type UploadController struct{}

func NewUploadController() *UploadController { return &UploadController{} }

// allowedExtensions 上传白名单（不含点）。
var allowedExtensions = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true, "webp": true,
	"pdf": true, "doc": true, "docx": true, "xls": true, "xlsx": true,
	"txt": true, "md": true, "zip": true,
}

// maxUploadSize 10MB。
const maxUploadSize = 10 << 20

// Store POST /admin-api/upload（multipart 字段名 file）。
// 响应 {status:0, data:{url}} 供 amis input-image / input-file receiver 消费。
func (r *UploadController) Store(ctx contractshttp.Context) contractshttp.Response {
	file, err := ctx.Request().File("file")
	if err != nil || file == nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("未接收到上传文件"))
	}

	size, err := file.Size()
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("读取文件失败："+err.Error()))
	}
	if size > maxUploadSize {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("文件不能超过 10MB"))
	}

	ext := strings.ToLower(strings.TrimPrefix(file.GetClientOriginalExtension(), "."))
	if ext == "" {
		ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(file.File()), "."))
	}
	if !allowedExtensions[ext] {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("不支持的文件类型："+ext))
	}

	disk := facades.Config().GetString("filesystems.default", "local")
	storage := facades.Storage().Disk(disk)

	relPath, err := storage.PutFile("uploads", file)
	if err != nil {
		return ctx.Response().Json(contractshttp.StatusOK, support.Fail("保存失败："+err.Error()))
	}
	// Windows 下 PutFile 可能返回 \ 分隔路径，统一为 URL 的 /
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	// local 驱动：/storage/* 静态路由映射磁盘根目录；
	// 云驱动：使用 Storage.Url() 返回的远端地址。
	url := ""
	if facades.Config().GetString("filesystems.disks."+disk+".driver", "local") == "local" {
		url = "/storage/" + relPath
	} else {
		url = storage.Url(relPath)
	}
	if url == "" {
		url = "/storage/" + relPath
	}

	return ctx.Response().Json(contractshttp.StatusOK, support.Success(map[string]any{
		"url": url,
	}, ""))
}
