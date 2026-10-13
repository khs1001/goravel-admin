package owl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khs1001/goravel-admin/models"
)

func sampleMenus() []models.SysMenu {
	return []models.SysMenu{
		{MenuID: 1, MenuName: "首页", ParentID: 0, Path: "/dashboard", Component: "dashboard", MenuType: models.MenuTypeMenu, OrderNum: 1, Status: "0"},
		{MenuID: 2, MenuName: "系统管理", ParentID: 0, Path: "/system", MenuType: models.MenuTypeDir, OrderNum: 2, Status: "0"},
		{MenuID: 3, MenuName: "用户管理", ParentID: 2, Path: "/system/user", Component: "system_user", Perms: "system:user:list", MenuType: models.MenuTypeMenu, OrderNum: 1, Status: "0"},
		{MenuID: 4, MenuName: "用户新增", ParentID: 3, Perms: "system:user:add", MenuType: models.MenuTypeButton, OrderNum: 1, Status: "0"},
	}
}

func TestBuildMenuTreeSkipsButtons(t *testing.T) {
	tree := BuildMenuTree(sampleMenus(), 0)

	assert.Len(t, tree, 2)
	assert.Equal(t, uint(2), tree[1].MenuID)
	assert.Len(t, tree[1].Children, 1)
	assert.Equal(t, uint(3), tree[1].Children[0].MenuID)
	assert.Empty(t, tree[1].Children[0].Children)
}

func TestBuildAllTreeKeepsButtons(t *testing.T) {
	tree := BuildAllTree(sampleMenus(), 0)

	assert.Len(t, tree, 2)
	assert.Len(t, tree[1].Children, 1)
	assert.Len(t, tree[1].Children[0].Children, 1)
}

func TestToRoutesFlagsHomeAndPageSign(t *testing.T) {
	tree := BuildMenuTree(sampleMenus(), 0)
	routes := ToRoutes(tree, 1)

	// Home route.
	assert.Equal(t, 1, routes[0]["is_home"])
	assert.Equal(t, UrlTypePage, routes[0]["url_type"])
	assert.Equal(t, "dashboard", routes[0]["page_sign"])

	// Directory with children.
	system := routes[1]
	assert.Equal(t, "amis", system["component"])
	children, ok := system["children"].([]map[string]any)
	assert.True(t, ok)
	assert.Len(t, children, 1)
	assert.Equal(t, "system_user", children[0]["page_sign"])
	assert.Equal(t, UrlTypePage, children[0]["url_type"])
}
