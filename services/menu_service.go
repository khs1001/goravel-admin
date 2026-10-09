package services

import (
	"sort"
	"strings"

	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/contracts"
	"github.com/khs1001/goravel-admin/models"
)

// MenuService 菜单服务：DB 菜单表 + 代码声明菜单 + 插件菜单 → 前端 IRoute 树。
// 逻辑对齐 owl-admin src/Support/Cores/Menu.php。
type MenuService struct{}

func NewMenuService() *MenuService { return &MenuService{} }

// All 构建完整菜单树（userID 用于按权限过滤 DB 菜单）。
func (r *MenuService) All(userID int, extraMenus []contracts.MenuItem, pluginMenus []contracts.MenuItem) []contracts.MenuItem {
	list, err := r.userMenus(userID)
	if err != nil {
		facades.Log().Errorf("[admin] load menus failed: %v", err)
	}

	items := make([]menuRow, 0, len(list)+len(pluginMenus))
	for _, m := range list {
		items = append(items, menuRow{menu: &m})
	}
	for i := range pluginMenus {
		items = append(items, menuRow{extra: &pluginMenus[i]})
	}

	// custom_order 升序，代码声明菜单排在其声明顺序
	sort.SliceStable(items, func(i, j int) bool {
		return r.orderOf(items[i]) < r.orderOf(items[j])
	})

	withoutExtra := r.withoutExtraRoutes()

	tree := r.buildTree(items, 0, "", withoutExtra)

	return append(tree, extraMenus...)
}

type menuRow struct {
	menu  *models.AdminMenu // DB 菜单
	extra *contracts.MenuItem // 代码/插件声明的成品菜单
}

func (r *MenuService) orderOf(row menuRow) int {
	if row.menu != nil {
		return row.menu.CustomOrder
	}

	return row.extra.Meta.CustomOrder
}

// userMenus 权限过滤：超管/未开启权限过滤返回全部；
// 普通用户仅返回其权限经 admin_permission_menu 挂载的菜单（对齐 owl-admin Menu::userMenus）。
func (r *MenuService) userMenus(userID int) ([]models.AdminMenu, error) {
	permissionService := NewPermissionService()
	all, ids := permissionService.VisibleMenus(userID)
	if all {
		var menus []models.AdminMenu
		err := facades.Orm().Query().Model(&models.AdminMenu{}).OrderBy("custom_order").Find(&menus)

		return menus, err
	}

	if len(ids) == 0 {
		return []models.AdminMenu{}, nil
	}

	var menus []models.AdminMenu
	err := facades.Orm().Query().
		Model(&models.AdminMenu{}).
		OrderBy("custom_order").
		WhereIn("id", intSliceToAny(ids)).
		Find(&menus)

	return menus, err
}

// intSliceToAny map keys → []any（WhereIn 参数）。
func intSliceToAny(ids map[int]bool) []any {
	out := make([]any, 0, len(ids))
	for id := range ids {
		out = append(out, id)
	}

	return out
}

// withoutExtraRoutes 不追加新增/详情/编辑额外路由的路径。
func (r *MenuService) withoutExtraRoutes() map[string]bool {
	set := map[string]bool{}
	for _, p := range facades.Config().GetStringSlice("admin.route.without_extra_routes", []string{}) {
		set[p] = true
	}

	return set
}

func (r *MenuService) buildTree(rows []menuRow, parentID int, parentName string, withoutExtra map[string]bool) []contracts.MenuItem {
	grouped := map[int][]menuRow{}
	for _, row := range rows {
		pid := 0
		if row.menu != nil {
			pid = row.menu.ParentID
		}

		grouped[pid] = append(grouped[pid], row)
	}

	return r.buildLevel(grouped, parentID, parentName, withoutExtra)
}

func (r *MenuService) buildLevel(grouped map[int][]menuRow, parentID int, parentName string, withoutExtra map[string]bool) []contracts.MenuItem {
	data := []contracts.MenuItem{}
	for _, row := range grouped[parentID] {
		var item contracts.MenuItem
		if row.menu != nil {
			item = r.formatDBMenu(row.menu, parentName)
		} else {
			item = *row.extra
		}

		children := r.buildLevel(grouped, rowID(row), item.Name, withoutExtra)
		if len(children) > 0 {
			item.Component = "amis"
			item.Children = children
		}

		data = append(data, item)

		// 额外路由：非页面菜单且未排除时，追加 create/show/edit 隐藏路由
		if !withoutExtra[item.Path] && item.UrlType != contracts.MenuTypePage && item.Path != "" && len(children) == 0 {
			data = append(data, r.generateRoutes(item)...)
		}
	}

	return data
}

func rowID(row menuRow) int {
	if row.menu != nil {
		return int(row.menu.ID)
	}

	return -1 // 代码菜单无子级挂载点
}

func (r *MenuService) formatDBMenu(m *models.AdminMenu, parentName string) contracts.MenuItem {
	idStr := "[" + itoa(int(m.ID)) + "]"
	name := idStr
	if parentName != "" {
		name = parentName + "-" + idStr
	}

	component := m.GetComponent()
	pageSign := ""
	switch m.UrlType {
	case contracts.MenuTypeIframe:
		component = "iframe"
	case contracts.MenuTypePage:
		component = "amis"
		pageSign = m.GetComponent()
	default:
		if component == "" {
			component = "amis"
		}
	}

	return contracts.MenuItem{
		Name:      name,
		Path:      m.GetURL(),
		Component: component,
		IsHome:    m.IsHome,
		UrlType:   m.UrlType,
		KeepAlive: m.KeepAlive,
		IsFull:    m.IsFull,
		IsLink:    m.UrlType == contracts.MenuTypeLink,
		PageSign:  pageSign,
		Meta: contracts.MenuMeta{
			Title:       m.Title,
			Icon:        m.GetIcon(),
			Hide:        m.Visible == 0,
			CustomOrder: m.CustomOrder,
		},
	}
}

// generateRoutes 追加 create/show/edit 隐藏路由（对齐 Menu::generateRoute）。
func (r *MenuService) generateRoutes(item contracts.MenuItem) []contracts.MenuItem {
	base := item.Path
	if i := strings.Index(base, "?"); i >= 0 {
		base = base[:i]
	}

	actionTitles := map[string]string{"create": "新增", "show": "详情", "edit": "编辑"}
	makeRoute := func(action, path string) contracts.MenuItem {
		return contracts.MenuItem{
			Name:      item.Name + "-" + action,
			Path:      base + path,
			Component: "amis",
			Meta: contracts.MenuMeta{
				Hide:  true,
				Icon:  item.Meta.Icon,
				Title: item.Meta.Title + " - " + actionTitles[action],
			},
		}
	}

	return []contracts.MenuItem{
		makeRoute("create", "/create"),
		makeRoute("show", "/:id"),
		makeRoute("edit", "/:id/edit"),
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	neg := v < 0
	if neg {
		v = -v
	}

	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}

	return string(b[i:])
}
