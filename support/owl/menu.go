package owl

import (
	"github.com/khs1001/goravel-admin/models"
)

// Menu url types, matching Owl Admin's AdminMenu constants.
const (
	UrlTypeRoute  = 1
	UrlTypeLink   = 2
	UrlTypeIframe = 3
	UrlTypePage   = 4
)

// BuildMenuTree assembles the RuoYi `sys_menu` rows into a tree.
//
// Button nodes (menu_type F) are ignored: they carry permissions, not routes.
func BuildMenuTree(list []models.SysMenu, parentID uint) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	for i := range list {
		if list[i].ParentID != parentID || list[i].MenuType == models.MenuTypeButton {
			continue
		}

		list[i].Children = BuildMenuTree(list, list[i].MenuID)
		tree = append(tree, list[i])
	}

	return tree
}

// BuildAllTree assembles the full menu tree including button nodes (menu_type
// F). It is used by the menu-management page.
func BuildAllTree(list []models.SysMenu, parentID uint) []models.SysMenu {
	tree := make([]models.SysMenu, 0)
	for i := range list {
		if list[i].ParentID != parentID {
			continue
		}

		list[i].Children = BuildAllTree(list, list[i].MenuID)
		tree = append(tree, list[i])
	}

	return tree
}

// ToRoutes converts a menu tree into the Owl route tree consumed by the
// frontend router.
//
// `pageSign` semantics: for a leaf menu (menu_type C) the RuoYi `component`
// column holds the page sign, and the route is emitted as a dynamic page
// (url_type = 4) so the frontend fetches its schema from `/page_schema?sign=`.
func ToRoutes(nodes []models.SysMenu, homeID uint) []map[string]any {
	routes := make([]map[string]any, 0, len(nodes))
	for _, node := range nodes {
		routes = append(routes, toRoute(node, homeID, ""))
	}

	return routes
}

func toRoute(node models.SysMenu, homeID uint, parentName string) map[string]any {
	name := "[" + itoa(node.MenuID) + "]"
	if parentName != "" {
		name = parentName + "-" + name
	}

	urlType := UrlTypeRoute
	isLink := false
	switch {
	case node.IsFrame == 1:
		urlType = UrlTypeLink
		isLink = true
	case len(node.Children) == 0 && node.MenuType == models.MenuTypeMenu:
		urlType = UrlTypePage
	}

	pageSign := ""
	component := "amis"
	if urlType == UrlTypePage {
		pageSign = node.Component
	}

	if urlType == UrlTypeLink {
		component = "amis"
	}

	route := map[string]any{
		"name":       name,
		"path":       node.Path,
		"component":  component,
		"is_home":    boolToInt(node.MenuID == homeID),
		"iframe_url": "",
		"url_type":   urlType,
		"keep_alive": node.IsCache,
		"is_full":    0,
		"is_link":    isLink,
		"page_sign":  pageSign,
		"meta": map[string]any{
			"title":        node.MenuName,
			"icon":         iconOrDash(node.Icon),
			"hide":         node.Visible == "1",
			"custom_order": node.OrderNum,
		},
	}

	if len(node.Children) > 0 {
		children := make([]map[string]any, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, toRoute(child, homeID, name))
		}
		route["component"] = "amis"
		route["children"] = children
	}

	return route
}

func iconOrDash(icon string) string {
	if icon == "" {
		return "-"
	}

	return icon
}

func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
