package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	moduleImport := setup.Paths().Module().Import()
	serviceProvider := "&admin.ServiceProvider{}"

	// The publishable configuration template. The source path is relative to the
	// package root (module), the destination is the host config directory.
	config, err := file.GetPackageContent(moduleImport, "config/admin.go")
	if err != nil {
		panic(err)
	}
	configDest := path.Config("admin.go")

	setup.Install(
		// Register the service provider in bootstrap/providers.go (bootstrap scaffold).
		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.RegisterProvider(moduleImport, serviceProvider)),

		// Register the service provider in app.go (legacy scaffold).
		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(path.Config("app.go")).
			Find(match.Imports()).Modify(modify.AddImport(moduleImport)).
			Find(match.Providers()).Modify(modify.Register(serviceProvider))),

		// Publish the default configuration.
		modify.File(configDest).Overwrite(config),
	).Uninstall(
		// Remove the published configuration.
		modify.File(configDest).Remove(),

		// Unregister the service provider from bootstrap/providers.go (bootstrap scaffold).
		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.UnregisterProvider(moduleImport, serviceProvider)),

		// Unregister the service provider from app.go (legacy scaffold).
		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(path.Config("app.go")).
			Find(match.Providers()).Modify(modify.Unregister(serviceProvider)).
			Find(match.Imports()).Modify(modify.RemoveImport(moduleImport))),
	).Execute()
}
