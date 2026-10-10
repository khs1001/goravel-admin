package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)

	// 配置模板：包内 config/admin.go（合法 Go 文件）。
	// 原样发布到宿主 config/admin.go；vendor:publish --package=github.com/khs1001/goravel-admin 亦可发布。
	configContent, err := file.GetPackageContent(setup.Paths().Module().String(), "config/admin.go")
	if err != nil {
		panic(err)
	}

	serviceProvider := "&admin.ServiceProvider{}"
	moduleImport := setup.Paths().Module().Import()
	configPath := path.Config("admin.go")

	setup.Install(
		// Register the service provider in bootstrap/providers.go
		modify.RegisterProvider(moduleImport, serviceProvider),

		// Publish the config file to the project's config directory
		modify.File(configPath).Overwrite(configContent),
	).Uninstall(
		// Remove config/admin.go
		modify.File(configPath).Remove(),

		// Unregister the service provider
		modify.UnregisterProvider(moduleImport, serviceProvider),
	).Execute()
}
