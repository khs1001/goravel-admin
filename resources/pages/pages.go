// Package pages embeds the amis page schemas served by `/page_schema`.
//
// Page schemas are data, so they live as JSON files next to this package and
// are compiled into the binary with go:embed. A page is addressed by a "sign"
// (the RuoYi `sys_menu.component` value), e.g. "system_user".
package pages

import (
	"embed"
	"encoding/json"
	"errors"
)

//go:embed *.json
var content embed.FS

// Get returns the amis schema registered under the given sign.
func Get(sign string) (map[string]any, error) {
	if sign == "" {
		return nil, errors.New("empty sign")
	}

	raw, err := content.ReadFile(sign + ".json")
	if err != nil {
		return nil, err
	}

	schema := map[string]any{}
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil, err
	}

	return schema, nil
}
