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
	serviceProvider := "&admin.ServiceProvider{}"
	moduleImport := setup.Paths().Module().Import()
	configPath := path.Config("admin.go")

	config, err := file.GetPackageContent(setup.Paths().Module().String(), configPath)
	if err != nil {
		panic(err)
	}

	setup.Install(
		// Register the service provider
		modify.RegisterProvider(moduleImport, serviceProvider),

		// Add the config file to the config directory
		modify.File(configPath).Overwrite(config),
	).Uninstall(
		// Remove the config file from the config directory
		modify.File(configPath).Remove(),

		// Unregister the service provider
		modify.UnregisterProvider(moduleImport, serviceProvider),
	).Execute()
}
