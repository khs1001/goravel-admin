package commands

import (
	"errors"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// CreateUser creates an admin user (sys_user).
type CreateUser struct{}

// Signature The name and signature of the console command.
func (r *CreateUser) Signature() string {
	return "admin:create-user"
}

// Description The console command description.
func (r *CreateUser) Description() string {
	return "创建一个后台管理员账号"
}

// Extend The console command extend.
func (r *CreateUser) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "登录用户名"},
			&command.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "登录密码"},
			&command.StringFlag{Name: "name", Usage: "昵称"},
		},
	}
}

// Handle Execute the console command.
func (r *CreateUser) Handle(ctx console.Context) error {
	username := ctx.Option("username")
	if username == "" {
		username = ctx.Argument(0)
	}
	if username == "" {
		return errors.New("用户名不能为空")
	}

	password := ctx.Option("password")
	if password == "" {
		password = "admin123"
	}

	var existing int64
	existing, err := facades.Orm().Query().Model(&models.SysUser{}).Where("user_name", username).Count()
	if err != nil {
		return err
	}
	if existing > 0 {
		return errors.New("用户已存在: " + username)
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return err
	}

	name := ctx.Option("name")
	if name == "" {
		name = username
	}

	user := models.SysUser{
		UserName: username,
		NickName: name,
		UserType: "00",
		Status:   "0",
		DelFlag:  "0",
		Password: hashed,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return err
	}

	ctx.Success("管理员创建成功: " + username)

	return nil
}
