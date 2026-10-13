package commands

import (
	"errors"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// ResetPassword resets an admin user's password.
type ResetPassword struct{}

// Signature The name and signature of the console command.
func (r *ResetPassword) Signature() string {
	return "admin:reset-password"
}

// Description The console command description.
func (r *ResetPassword) Description() string {
	return "重置后台管理员密码"
}

// Extend The console command extend.
func (r *ResetPassword) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "登录用户名"},
			&command.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "新密码"},
		},
	}
}

// Handle Execute the console command.
func (r *ResetPassword) Handle(ctx console.Context) error {
	username := ctx.Option("username")
	if username == "" {
		username = ctx.Argument(0)
	}
	if username == "" {
		return errors.New("用户名不能为空")
	}

	password := ctx.Option("password")
	if password == "" {
		return errors.New("新密码不能为空")
	}

	var user models.SysUser
	if err := facades.Orm().Query().Where("user_name", username).First(&user); err != nil {
		return errors.New("用户不存在: " + username)
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return err
	}
	user.Password = hashed

	if err := facades.Orm().Query().Save(&user); err != nil {
		return err
	}

	ctx.Success("密码已重置: " + username)

	return nil
}
