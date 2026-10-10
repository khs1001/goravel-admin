package admin

import (
	"embed"
	"io/fs"
	"net/http"
)

// 前端构建产物（owl-admin admin-views）随扩展包分发，经 go:embed 嵌入二进制。
// 前端资源注册规范：由扩展包内部在 Boot 阶段经 StaticFS 注册（/admin-assets/*），
// 宿主应用的 public/ 不携带、不注册后台前端资源。
//
//go:embed all:admin-assets
var adminAssetsFS embed.FS

// AdminAssets 前端构建产物的嵌入式文件系统。
func AdminAssets() http.FileSystem {
	sub, err := fs.Sub(adminAssetsFS, "admin-assets")
	if err != nil {
		return http.FS(adminAssetsFS)
	}

	return http.FS(sub)
}
