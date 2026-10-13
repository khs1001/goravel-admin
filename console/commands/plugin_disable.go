package commands

import (
	"errors"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	"github.com/khs1001/goravel-admin/support/plugin"
)

// PluginDisable disables a discovered plugin.
type PluginDisable struct{}

// Signature The name and signature of the console command.
func (r *PluginDisable) Signature() string {
	return "admin:plugin:disable"
}

// Description The console command description.
func (r *PluginDisable) Description() string {
	return "禁用后台插件"
}

// Extend The console command extend.
func (r *PluginDisable) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "插件名"},
		},
	}
}

// Handle Execute the console command.
func (r *PluginDisable) Handle(ctx console.Context) error {
	name := ctx.Option("name")
	if name == "" {
		name = ctx.Argument(0)
	}
	if name == "" {
		return errors.New("插件名不能为空")
	}
	if !plugin.R().Has(name) {
		return errors.New("插件不存在或未安装: " + name)
	}

	if err := plugin.R().SetEnabled(name, false); err != nil {
		return err
	}

	ctx.Success("已禁用插件: " + name)

	return nil
}
