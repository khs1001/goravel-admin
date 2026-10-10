package console

import (
	"os"

	"github.com/goravel/framework/contracts/console"
	consolecommand "github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// InstallCommand 一键安装：迁移 → 数据核验 → 配置体检。
//
//	php artisan admin:install
//
// 幂等：迁移重复执行为 no-op，初始数据由迁移内 seed 保障（admin_users 为空时写入），
// 体检项只报告不修改（storage 上传目录缺失会自动创建）。
type InstallCommand struct{}

func (r *InstallCommand) Signature() string { return "admin:install" }

func (r *InstallCommand) Description() string {
	return "Install the goravel-admin extension: run migrations, verify seed data and check configuration"
}

func (r *InstallCommand) Extend() consolecommand.Extend {
	return consolecommand.Extend{Category: "admin"}
}

func (r *InstallCommand) Handle(ctx console.Context) error {
	ctx.Info("goravel-admin installing...")

	// 1. 迁移（admin 包自身 + 二级插件迁移已在 Boot 阶段注册，此处统一执行）
	if err := facades.Artisan().Call("migrate"); err != nil {
		ctx.Error("migration failed: " + err.Error())

		return err
	}
	ctx.NewLine()
	ctx.Info("✔ migrations done")

	// 2. 数据核验：默认超管账号（由迁移 seed 写入，此处核验并提示）
	userCount, err := facades.Orm().Query().Model(&models.AdminUser{}).Count()
	if err != nil {
		ctx.Error("✘ query admin_users failed: " + err.Error())

		return err
	}
	switch {
	case userCount == 0:
		ctx.Error("✘ no admin user found — migration seed failed, please check database logs")
	case userCount == 1:
		ctx.Comment("✔ default super admin ready: admin / admin (change the password after first login)")
	default:
		ctx.Comment("✔ admin users exist")
	}

	// 3. 配置体检（只报告，不修改）
	ctx.NewLine()
	ctx.Info("configuration checks:")

	checkEnvSecrets(ctx)
	checkQueueDriver(ctx)
	checkFilesystemDisk(ctx)
	checkCaptchaFlag(ctx)

	ctx.NewLine()
	ctx.Info("done. admin panel: {APP_URL}/admin-assets/")

	return nil
}

// checkEnvSecrets 密钥强度（生产环境由宿主启动逻辑强制，此处给出提醒）。
func checkEnvSecrets(ctx console.Context) {
	if len(facades.Config().EnvString("JWT_SECRET", "")) < 32 {
		ctx.Error("✘ JWT_SECRET is weak or empty — set a strong random secret before production")
	} else {
		ctx.Comment("✔ JWT_SECRET ok")
	}

	if facades.Config().EnvString("APP_KEY", "") == "" {
		ctx.Error("✘ APP_KEY is empty — run `go run . artisan key:generate`")
	} else {
		ctx.Comment("✔ APP_KEY ok")
	}
}

// checkQueueDriver 队列驱动与 worker 提示。
func checkQueueDriver(ctx console.Context) {
	driver := facades.Config().GetString("queue.default", "sync")
	if driver == "database" {
		ctx.Comment("✔ queue driver: database (embedded worker required — host main.go already starts one)")
	} else {
		ctx.Comment("✔ queue driver: " + driver)
	}
}

// checkFilesystemDisk 上传磁盘：public 上传目录缺失时自动创建（/storage/* 静态访问依赖）。
func checkFilesystemDisk(ctx console.Context) {
	disk := facades.Config().GetString("filesystems.default", "local")
	ctx.Comment("✔ filesystem disk: " + disk)

	if disk != "local" && disk != "public" {
		return
	}

	dir := "storage/app/public"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			ctx.Error("✘ cannot create upload directory " + dir + ": " + err.Error())

			return
		}
		ctx.Comment("✔ created upload directory: " + dir)
	} else {
		ctx.Comment("✔ upload directory ok: " + dir)
	}
}

// checkCaptchaFlag 验证码开关状态提示。
func checkCaptchaFlag(ctx console.Context) {
	if facades.Config().GetBool("admin.auth.login_captcha", false) {
		ctx.Comment("✔ login captcha enabled")
	} else {
		ctx.Comment("○ login captcha disabled (config admin.auth.login_captcha)")
	}
}

// Ensure InstallCommand implements the contract.
var _ console.Command = (*InstallCommand)(nil)
