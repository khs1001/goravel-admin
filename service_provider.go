package admin

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/foundation"

	"github.com/khs1001/goravel-admin/console/commands"
	"github.com/khs1001/goravel-admin/database/migrations"
	"github.com/khs1001/goravel-admin/database/seeders"
	"github.com/khs1001/goravel-admin/routes"
	"github.com/khs1001/goravel-admin/support/plugin"
)

// Binding is the container key of the package service. It is also the config
// key and the name of the facade: config/admin.go <-> Binding "admin" <->
// facades.Admin().
const Binding = "admin"

// ModulePath is the Go import path of this module. It is used to register the
// package assets with the vendor:publish command.
const ModulePath = "github.com/khs1001/goravel-admin"

// App is the application handle captured at Register time, so the facade can
// resolve the binding later.
var App foundation.Application

// ServiceProvider wires the admin package into a Goravel application.
type ServiceProvider struct {
}

// Relationship returns the relationship of the service provider.
func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{Binding},
		Dependencies: []string{
			binding.Config,
			binding.Log,
			binding.Cache,
			binding.Http,
			binding.Route,
			binding.Orm,
			binding.Schema,
			binding.Seeder,
			binding.Auth,
			binding.Hash,
		},
		ProvideFor: []string{},
	}
}

// Register registers the service provider. Only container bindings and
// configuration live here.
func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return NewAdmin(app), nil
	})

	registerAdminGuard(app)
}

// Boot wires every resource layer at runtime through the foundation
// Application, so package install/uninstall stays symmetric and idempotent.
func (r *ServiceProvider) Boot(app foundation.Application) {
	// Discover plugins from the container bindings first, so routes / menus /
	// migrations / commands can include their contributions.
	plugin.Discover(app)

	if router := app.MakeRoute(); router != nil {
		routes.Admin(router)
	}

	allCommands := []console.Command{
		&commands.CreateUser{},
		&commands.ResetPassword{},
		&commands.PluginEnable{},
		&commands.PluginDisable{},
		&commands.PluginClean{},
	}
	allCommands = append(allCommands, plugin.R().Commands()...)
	app.Commands(allCommands)

	if schemaInstance := app.MakeSchema(); schemaInstance != nil {
		allMigrations := []schema.Migration{
			&migrations.CreateSysTables{},
			&migrations.CreateSysSystemTables{},
			&migrations.CreateSysLogTables{},
			&migrations.CreateSysExtensionTable{},
		}
		allMigrations = append(allMigrations, plugin.R().Migrations()...)
		schemaInstance.Register(allMigrations)
	}

	if seederInstance := app.MakeSeeder(); seederInstance != nil {
		allSeeders := []seeder.Seeder{
			&seeders.SysInitSeeder{},
			&seeders.SysSystemSeeder{},
		}
		allSeeders = append(allSeeders, plugin.R().Seeders()...)
		seederInstance.Register(allSeeders)
	}

	// Publish the config template so the host can override it.
	app.Publishes(ModulePath, map[string]string{
		"config/admin.go": app.ConfigPath("admin.go"),
	})
}

// registerAdminGuard adds the dedicated `admin` JWT guard and its ORM user
// provider to the host auth configuration.
//
// It runs at Register time (after the configuration is fully loaded) and only
// fills in missing keys, so any host-defined admin guard wins.
func registerAdminGuard(app foundation.Application) {
	config := app.MakeConfig()
	if config == nil {
		return
	}

	authConfig, ok := config.Get("auth").(map[string]any)
	if !ok || authConfig == nil {
		authConfig = map[string]any{}
	}

	guards, ok := authConfig["guards"].(map[string]any)
	if !ok || guards == nil {
		guards = map[string]any{}
	}
	if _, exists := guards["admin"]; !exists {
		guards["admin"] = map[string]any{"driver": "jwt", "provider": "admin"}
	}
	authConfig["guards"] = guards

	providers, ok := authConfig["providers"].(map[string]any)
	if !ok || providers == nil {
		providers = map[string]any{}
	}
	if _, exists := providers["admin"]; !exists {
		providers["admin"] = map[string]any{"driver": "orm"}
	}
	authConfig["providers"] = providers

	config.Add("auth", authConfig)
}
