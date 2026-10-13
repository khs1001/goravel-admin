// Package plugin implements the in-process plugin registry: discovery via the
// container bindings, contribution aggregation and hook dispatch.
//
// It is a leaf package (it does not import the admin root package, routes,
// controllers or services) so every other layer can depend on it without
// creating an import cycle.
package plugin

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/models"
)

// KeyPrefix is the container-binding key prefix a plugin must use to be
// discovered. Plugins do not need to import this package to use it.
const KeyPrefix = "admin.plugin."

// Binding returns the container key for a plugin name.
func Binding(name string) string {
	return KeyPrefix + name
}

// Plugin is a discovered plugin and its resolved state.
type Plugin struct {
	Name        string
	Version     string
	Description string
	Author      string
	Enabled     bool

	impl any
}

// Registry holds the discovered plugins.
type Registry struct {
	mu      sync.RWMutex
	plugins []*Plugin
	byName  map[string]*Plugin
}

var registry = &Registry{byName: map[string]*Plugin{}}

// R returns the process-wide registry.
func R() *Registry {
	return registry
}

// Discover scans the container bindings for `admin.plugin.*` keys, resolves
// them, records them in `sys_extension` (created disabled) and registers them.
//
// It is called from the admin ServiceProvider Boot, after all plugin providers
// have run their Register.
func Discover(app foundation.Application) {
	for _, key := range app.Bindings() {
		name, ok := key.(string)
		if !ok || !strings.HasPrefix(name, KeyPrefix) {
			continue
		}

		instance, err := app.Make(name)
		if err != nil {
			facades.Log().Errorf("admin plugin %s resolve failed: %s", name, err.Error())
			continue
		}

		base, ok := instance.(contracts.Plugin)
		if !ok {
			facades.Log().Errorf("admin plugin %s does not implement Plugin (Name)", name)
			continue
		}

		registry.add(base, instance)
	}

	registry.syncStates()
}

func (r *Registry) add(base contracts.Plugin, impl any) {
	defer recoverLog("plugin metadata " + safeName(base))

	name := strings.TrimSpace(base.Name())
	if name == "" {
		facades.Log().Error("admin plugin with empty Name ignored")
		return
	}

	entry := &Plugin{Name: name, impl: impl}
	if meta, ok := impl.(contracts.MetadataProvider); ok {
		entry.Version = meta.Version()
		entry.Description = meta.Description()
		entry.Author = meta.Author()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byName[name]; exists {
		return
	}
	r.plugins = append(r.plugins, entry)
	r.byName[name] = entry
}

// syncStates upserts every discovered plugin into sys_extension (default
// disabled) and loads the enabled flag.
func (r *Registry) syncStates() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.plugins {
		var record models.SysExtension
		err := facades.Orm().Query().Model(&models.SysExtension{}).Where("name", entry.Name).FirstOrFail(&record)
		if err != nil {
			record = models.SysExtension{Name: entry.Name, IsEnabled: 0, Version: entry.Version}
			if createErr := facades.Orm().Query().Create(&record); createErr != nil {
				facades.Log().Errorf("admin plugin %s record create failed: %s", entry.Name, createErr.Error())
				continue
			}
		} else if entry.Version != "" && record.Version != entry.Version {
			record.Version = entry.Version
			_ = facades.Orm().Query().Save(&record)
		}

		entry.Enabled = record.IsEnabled == 1
	}
}

// All returns every discovered plugin (in discovery order).
func (r *Registry) All() []*Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]*Plugin{}, r.plugins...)
}

// Enabled returns the enabled plugins.
func (r *Registry) Enabled() []*Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Plugin, 0, len(r.plugins))
	for _, entry := range r.plugins {
		if entry.Enabled {
			result = append(result, entry)
		}
	}

	return result
}

// IsEnabled reports whether a plugin is currently enabled.
func (r *Registry) IsEnabled(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, ok := r.byName[name]

	return ok && entry.Enabled
}

// List returns plugin descriptors for the extension-management page.
func (r *Registry) List() []map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]map[string]any, 0, len(r.plugins))
	for _, entry := range r.plugins {
		result = append(result, map[string]any{
			"name":         entry.Name,
			"version":      entry.Version,
			"description":  entry.Description,
			"author":       entry.Author,
			"enabled":      enabledInt(entry.Enabled),
			"capabilities": capabilities(entry.impl),
		})
	}

	return result
}

// RouteProviders returns the enabled-or-not route providers together with the
// owning plugin name (routes are always registered; the guard gates them).
func (r *Registry) RouteProviders() []struct {
	Name     string
	Provider contracts.RouteProvider
} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]struct {
		Name     string
		Provider contracts.RouteProvider
	}, 0)
	for _, entry := range r.plugins {
		if provider, ok := entry.impl.(contracts.RouteProvider); ok {
			result = append(result, struct {
				Name     string
				Provider contracts.RouteProvider
			}{Name: entry.Name, Provider: provider})
		}
	}

	return result
}

// MenuNodes aggregates menu / permission nodes from enabled plugins, with page
// signs namespaced as "<plugin>.<short>".
func (r *Registry) MenuNodes() []map[string]any {
	result := make([]map[string]any, 0)
	for _, entry := range r.Enabled() {
		provider, ok := entry.impl.(contracts.MenuProvider)
		if !ok {
			continue
		}

		nodes := safeCall("plugin menus "+entry.Name, func() []map[string]any { return provider.Menus() })
		for _, node := range nodes {
			normalized := normalizeMenuNode(entry.Name, node)
			if normalized != nil {
				result = append(result, normalized)
			}
		}
	}

	return result
}

// Pages aggregates amis pages from enabled plugins, keyed by namespaced sign.
func (r *Registry) Pages() map[string]map[string]any {
	result := map[string]map[string]any{}
	for _, entry := range r.Enabled() {
		provider, ok := entry.impl.(contracts.PageProvider)
		if !ok {
			continue
		}

		pages := safeCall("plugin pages "+entry.Name, func() map[string]map[string]any { return provider.Pages() })
		for sign, schema := range pages {
			result[entry.Name+"."+sign] = schema
		}
	}

	return result
}

// Settings aggregates default settings from enabled plugins.
func (r *Registry) Settings() map[string]any {
	result := map[string]any{}
	for _, entry := range r.Enabled() {
		provider, ok := entry.impl.(contracts.SettingProvider)
		if !ok {
			continue
		}

		settings := safeCall("plugin settings "+entry.Name, func() map[string]any { return provider.Settings() })
		for key, value := range settings {
			result[key] = value
		}
	}

	return result
}

// Assets aggregates js/css from enabled plugins (discovery order, de-duplicated).
func (r *Registry) Assets() map[string]any {
	js := make([]string, 0)
	css := make([]string, 0)
	seenJS := map[string]bool{}
	seenCSS := map[string]bool{}

	for _, entry := range r.Enabled() {
		provider, ok := entry.impl.(contracts.AssetProvider)
		if !ok {
			continue
		}

		assets := safeCall("plugin assets "+entry.Name, func() map[string]any { return provider.Assets() })
		for _, url := range toStringSlice(assets["js"]) {
			if !seenJS[url] {
				seenJS[url] = true
				js = append(js, url)
			}
		}
		for _, url := range toStringSlice(assets["css"]) {
			if !seenCSS[url] {
				seenCSS[url] = true
				css = append(css, url)
			}
		}
	}

	return map[string]any{"js": js, "css": css}
}

// Nav aggregates navbar schemas from enabled plugins.
func (r *Registry) Nav() map[string]any {
	result := map[string]any{}
	for _, entry := range r.Enabled() {
		provider, ok := entry.impl.(contracts.NavProvider)
		if !ok {
			continue
		}

		nav := safeCall("plugin nav "+entry.Name, func() map[string]any { return provider.Nav() })
		for key, value := range nav {
			result[key] = value
		}
	}

	return result
}

// Migrations aggregates migrations from every plugin (not gated; install-time).
func (r *Registry) Migrations() []schema.Migration {
	result := make([]schema.Migration, 0)
	for _, entry := range r.All() {
		provider, ok := entry.impl.(contracts.MigrationProvider)
		if !ok {
			continue
		}
		result = append(result, safeCall("plugin migrations "+entry.Name, func() []schema.Migration { return provider.Migrations() })...)
	}

	return result
}

// Seeders aggregates seeders from every plugin.
func (r *Registry) Seeders() []seeder.Seeder {
	result := make([]seeder.Seeder, 0)
	for _, entry := range r.All() {
		provider, ok := entry.impl.(contracts.SeederProvider)
		if !ok {
			continue
		}
		result = append(result, safeCall("plugin seeders "+entry.Name, func() []seeder.Seeder { return provider.Seeders() })...)
	}

	return result
}

// Commands aggregates console commands from every plugin.
func (r *Registry) Commands() []console.Command {
	result := make([]console.Command, 0)
	for _, entry := range r.All() {
		provider, ok := entry.impl.(contracts.CommandProvider)
		if !ok {
			continue
		}
		result = append(result, safeCall("plugin commands "+entry.Name, func() []console.Command { return provider.Commands() })...)
	}

	return result
}

// RegisterRoutes registers every plugin's routes on a sub-router already
// wrapped with auth + operation log + plugin guard.
func (r *Registry) RegisterRoutes(router route.Router, guard func(name string) http.Middleware) {
	for _, item := range r.RouteProviders() {
		item := item
		func() {
			defer recoverLog("plugin routes " + item.Name)
			router.Prefix("/ext/" + item.Name).Middleware(guard(item.Name)).Group(func(sub route.Router) {
				item.Provider.Routes(sub)
			})
		}()
	}
}

// ---- hooks ----

func (r *Registry) RunUserAuthenticated(ctx http.Context, user any) {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.UserAuthenticatedHook); ok {
			safeRun("plugin OnUserAuthenticated "+entry.Name, func() { hook.OnUserAuthenticated(ctx, user) })
		}
	}
}

func (r *Registry) RunLoginSucceeded(ctx http.Context, username string) {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.LoginHook); ok {
			safeRun("plugin OnLoginSucceeded "+entry.Name, func() { hook.OnLoginSucceeded(ctx, username) })
		}
	}
}

func (r *Registry) RunLoginFailed(ctx http.Context, username, msg string) {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.LoginHook); ok {
			safeRun("plugin OnLoginFailed "+entry.Name, func() { hook.OnLoginFailed(ctx, username, msg) })
		}
	}
}

func (r *Registry) RunLogout(ctx http.Context, username string) {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.LogoutHook); ok {
			safeRun("plugin OnLogout "+entry.Name, func() { hook.OnLogout(ctx, username) })
		}
	}
}

func (r *Registry) RunOperationLogged(ctx http.Context, entry map[string]any) {
	for _, item := range r.Enabled() {
		if hook, ok := item.impl.(contracts.OperationLoggedHook); ok {
			safeRun("plugin OnOperationLogged "+item.Name, func() { hook.OnOperationLogged(ctx, entry) })
		}
	}
}

func (r *Registry) RunMenusBuilt(ctx http.Context, routes []map[string]any) []map[string]any {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.MenusBuiltHook); ok {
			routes = safeCallRaw("plugin OnMenusBuilt "+entry.Name, routes, func() []map[string]any {
				return hook.OnMenusBuilt(ctx, routes)
			})
		}
	}

	return routes
}

func (r *Registry) RunSettingsBuilt(ctx http.Context, payload map[string]any) {
	for _, entry := range r.Enabled() {
		if hook, ok := entry.impl.(contracts.SettingsBuiltHook); ok {
			safeRun("plugin OnSettingsBuilt "+entry.Name, func() { hook.OnSettingsBuilt(ctx, payload) })
		}
	}
}

// ---- helpers ----

func normalizeMenuNode(pluginName string, raw map[string]any) map[string]any {
	if raw == nil {
		return nil
	}

	node := map[string]any{
		"title":    stringValue(raw["title"]),
		"path":     stringValue(raw["path"]),
		"parent":   stringValue(raw["parent"]),
		"icon":     stringValue(raw["icon"]),
		"perms":    stringValue(raw["perms"]),
		"menuType": defaultValue(stringValue(raw["menuType"]), "C"),
		"order":    intValue(raw["order"]),
		"plugin":   pluginName,
	}

	if path := node["path"].(string); path == "" {
		return nil
	}

	component := stringValue(raw["component"])
	if component != "" && node["menuType"].(string) == "C" {
		component = pluginName + "." + component
	}
	node["component"] = component

	return node
}

func capabilities(impl any) map[string]bool {
	_, menu := impl.(contracts.MenuProvider)
	_, routes := impl.(contracts.RouteProvider)
	_, pages := impl.(contracts.PageProvider)
	_, settings := impl.(contracts.SettingProvider)
	_, assets := impl.(contracts.AssetProvider)
	_, nav := impl.(contracts.NavProvider)
	_, commands := impl.(contracts.CommandProvider)
	_, migrations := impl.(contracts.MigrationProvider)
	_, seeders := impl.(contracts.SeederProvider)
	_, hook := impl.(contracts.MenusBuiltHook)

	return map[string]bool{
		"menu": menu, "route": routes, "page": pages, "setting": settings,
		"asset": assets, "nav": nav, "command": commands, "migration": migrations,
		"seeder": seeders, "hook": hook,
	}
}

func safeName(base contracts.Plugin) string {
	defer func() { _ = recover() }()

	return base.Name()
}

func safeRun(label string, fn func()) {
	defer recoverLog(label)
	fn()
}

func safeCall[T any](label string, fn func() T) T {
	defer recoverLog(label)

	return fn()
}

func safeCallRaw[T any](label string, fallback T, fn func() T) T {
	defer func() {
		if r := recover(); r != nil {
			facades.Log().Errorf("%s panicked: %v", label, r)
		}
	}()

	return fn()
}

func recoverLog(label string) {
	if r := recover(); r != nil {
		facades.Log().Errorf("%s panicked: %v", label, r)
	}
}

func enabledInt(enabled bool) int {
	if enabled {
		return 1
	}

	return 0
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", value)
}

func intValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}

	return 0
}

func defaultValue(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func toStringSlice(value any) []string {
	switch v := value.(type) {
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			result = append(result, stringValue(item))
		}
		return result
	case string:
		if v == "" {
			return nil
		}
		return []string{v}
	}

	return nil
}

// SortByOrder is a small helper used by the menu merge.
func SortByOrder(nodes []map[string]any) {
	sort.SliceStable(nodes, func(i, j int) bool {
		return intValue(nodes[i]["order"]) < intValue(nodes[j]["order"])
	})
}

// Has reports whether a plugin with the given name was discovered.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.byName[name]

	return ok
}

// Names returns the discovered plugin names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for _, entry := range r.plugins {
		names = append(names, entry.Name)
	}

	return names
}

// SetEnabled updates the in-memory state and persists it in sys_extension.
func (r *Registry) SetEnabled(name string, enabled bool) error {
	r.mu.RLock()
	entry, ok := r.byName[name]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}

	value := 0
	if enabled {
		value = 1
	}

	var record models.SysExtension
	if err := facades.Orm().Query().Model(&models.SysExtension{}).Where("name", name).FirstOrFail(&record); err != nil {
		record = models.SysExtension{Name: name, IsEnabled: value, Version: entry.Version}
		if createErr := facades.Orm().Query().Create(&record); createErr != nil {
			return createErr
		}
	} else {
		record.IsEnabled = value
		if err := facades.Orm().Query().Save(&record); err != nil {
			return err
		}
	}

	r.mu.Lock()
	entry.Enabled = enabled
	r.mu.Unlock()

	return nil
}

// Clean removes sys_extension rows for plugins that are no longer discovered.
func (r *Registry) Clean() (*int64, error) {
	names := r.Names()
	query := facades.Orm().Query()
	if len(names) > 0 {
		args := make([]any, len(names))
		for i, name := range names {
			args[i] = name
		}
		query = query.WhereNotIn("name", args)
	} else {
		query = query.Where("extension_id > 0")
	}

	result, err := query.Delete(&models.SysExtension{})
	if err != nil {
		return nil, err
	}

	affected := result.RowsAffected

	return &affected, nil
}
