package commands

import (
	"fmt"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	"github.com/khs1001/goravel-admin/support/plugin"
)

// PluginClean removes sys_extension records for plugins no longer discovered.
type PluginClean struct{}

// Signature The name and signature of the console command.
func (r *PluginClean) Signature() string {
	return "admin:plugin:clean"
}

// Description The console command description.
func (r *PluginClean) Description() string {
	return "清理已卸载插件的失效记录"
}

// Extend The console command extend.
func (r *PluginClean) Extend() command.Extend {
	return command.Extend{Category: "admin"}
}

// Handle Execute the console command.
func (r *PluginClean) Handle(ctx console.Context) error {
	count, err := plugin.R().Clean()
	if err != nil {
		return err
	}

	ctx.Success(fmt.Sprintf("已清理 %d 条失效插件记录", *count))

	return nil
}
