package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("admin", map[string]any{})
}
