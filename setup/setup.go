
package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	serviceProvider := "&admin.ServiceProvider{}"
	moduleImport := setup.Paths().Module().Import()
	configPath := path.Config("admin.go")

	setup.Install(
		// Register the service provider
		modify.RegisterProvider(moduleImport, serviceProvider),

		// Add config
		modify.File(configPath).Overwrite(config(setup.Paths().Config().Package(), setup.Paths().Facades().Import(), setup.Paths().Facades().Package())),
	).Uninstall(
		// Remove config/cache.go
		modify.File(configPath).Remove(),

		// Unregister the service provider
		modify.UnregisterProvider(moduleImport, serviceProvider),
	).Execute()
}
