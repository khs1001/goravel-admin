package owl

import (
	"fmt"
	"sort"
)

// PluginNodeToRoute converts a normalized plugin menu node (produced by
// support/plugin) into an owl route map.
func PluginNodeToRoute(node map[string]any) map[string]any {
	path := asString(node["path"])
	component := asString(node["component"])
	pluginName := asString(node["plugin"])
	title := asString(node["title"])
	icon := asString(node["icon"])
	menuType := asString(node["menuType"])

	urlType := UrlTypeRoute
	if menuType == "C" && component != "" {
		urlType = UrlTypePage
	}

	pageSign := ""
	if urlType == UrlTypePage {
		pageSign = component
	}

	return map[string]any{
		"name":       "plugin:" + pluginName + ":" + path,
		"path":       path,
		"component":  "amis",
		"is_home":    0,
		"iframe_url": "",
		"url_type":   urlType,
		"keep_alive": 0,
		"is_full":    0,
		"is_link":    false,
		"page_sign":  pageSign,
		"meta": map[string]any{
			"title":        title,
			"icon":         iconOrDash(icon),
			"hide":         false,
			"custom_order": asInt(node["order"]),
		},
		"_plugin": pluginName,
		"_parent": asString(node["parent"]),
	}
}

// MergePluginRoutes merges plugin route maps into the database-derived route
// tree. A plugin node with a resolvable `parent` (by path) is inserted under it;
// otherwise it becomes a top level node. Siblings are sorted by `custom_order`
// (ties keep host routes before plugin routes).
func MergePluginRoutes(routes []map[string]any, nodes []map[string]any) []map[string]any {
	index := map[string]map[string]any{}
	indexRoutes(routes, index)

	for _, node := range nodes {
		route := PluginNodeToRoute(node)
		if route["path"] == "" {
			continue
		}

		parentPath := asString(node["parent"])
		if parentPath != "" {
			if parent, ok := index[parentPath]; ok {
				children, _ := parent["children"].([]map[string]any)
				children = append(children, route)
				parent["children"] = children
				parent["component"] = "amis"
				parent["page_sign"] = ""
				index[asString(route["path"])] = route
				continue
			}
		}

		routes = append(routes, route)
		index[asString(route["path"])] = route
	}

	sortRoutes(routes)
	cleanInternalKeys(routes)

	return routes
}

func indexRoutes(routes []map[string]any, index map[string]map[string]any) {
	for _, route := range routes {
		if path := asString(route["path"]); path != "" {
			index[path] = route
		}
		if children, ok := route["children"].([]map[string]any); ok {
			indexRoutes(children, index)
		}
	}
}

func sortRoutes(routes []map[string]any) {
	sort.SliceStable(routes, func(i, j int) bool {
		oi, oj := routeOrder(routes[i]), routeOrder(routes[j])
		if oi != oj {
			return oi < oj
		}
		pi, pj := isPluginRoute(routes[i]), isPluginRoute(routes[j])
		if pi != pj {
			return !pi
		}
		return false
	})

	for _, route := range routes {
		if children, ok := route["children"].([]map[string]any); ok {
			sortRoutes(children)
		}
	}
}

func cleanInternalKeys(routes []map[string]any) {
	for _, route := range routes {
		delete(route, "_plugin")
		delete(route, "_parent")
		if children, ok := route["children"].([]map[string]any); ok {
			cleanInternalKeys(children)
		}
	}
}

func routeOrder(route map[string]any) int {
	meta, ok := route["meta"].(map[string]any)
	if !ok {
		return 0
	}

	return asInt(meta["custom_order"])
}

func isPluginRoute(route map[string]any) bool {
	_, ok := route["_plugin"]

	return ok
}

func asString(value any) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", value)
}

func asInt(value any) int {
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
