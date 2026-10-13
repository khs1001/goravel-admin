// Package services holds the business logic for the admin package, organised
// around the RuoYi data model.
package services

import (
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// MenuService resolves menus and permissions for the RuoYi model.
type MenuService struct{}

// NewMenuService creates a MenuService.
func NewMenuService() *MenuService {
	return &MenuService{}
}

// All returns every menu, ordered for tree rendering.
func (s *MenuService) All() ([]models.SysMenu, error) {
	var menus []models.SysMenu
	if err := facades.Orm().Query().OrderBy("order_num").Get(&menus); err != nil {
		return nil, err
	}

	return menus, nil
}

// Enabled returns every enabled menu (status = 0).
func (s *MenuService) Enabled() ([]models.SysMenu, error) {
	var menus []models.SysMenu
	if err := facades.Orm().Query().Where("status", "0").OrderBy("order_num").Get(&menus); err != nil {
		return nil, err
	}

	return menus, nil
}

// Routable returns the enabled menus a user can navigate to.
//
// The super administrator (id 1) sees every enabled menu; other users see the
// union of the menus granted to their roles.
func (s *MenuService) Routable(user *models.SysUser) ([]models.SysMenu, error) {
	all, err := s.Enabled()
	if err != nil {
		return nil, err
	}

	if s.IsSuperAdmin(user) {
		return all, nil
	}

	allowed := make(map[uint]bool)
	roleIDs, err := s.roleIDs(user.UserID)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []models.SysMenu{}, nil
	}

	var roleMenus []models.SysRoleMenu
	if err := facades.Orm().Query().WhereIn("role_id", uintSliceToAny(roleIDs)).Get(&roleMenus); err != nil {
		return nil, err
	}
	for _, rm := range roleMenus {
		allowed[rm.MenuID] = true
	}

	result := make([]models.SysMenu, 0, len(all))
	for _, menu := range all {
		if allowed[menu.MenuID] {
			result = append(result, menu)
		}
	}

	return result, nil
}

// PermsForUser returns the permission identifiers granted to a user.
func (s *MenuService) PermsForUser(user *models.SysUser) ([]string, error) {
	menus, err := s.Routable(user)
	if err != nil {
		return nil, err
	}

	perms := make([]string, 0, len(menus))
	seen := make(map[string]bool)
	for _, menu := range menus {
		if menu.Perms == "" || seen[menu.Perms] {
			continue
		}
		seen[menu.Perms] = true
		perms = append(perms, menu.Perms)
	}

	return perms, nil
}

// HasPermission reports whether the user holds the given RuoYi permission
// identifier (e.g. "system:user:add").
func (s *MenuService) HasPermission(user *models.SysUser, perm string) bool {
	if perm == "" || s.IsSuperAdmin(user) {
		return true
	}

	perms, err := s.PermsForUser(user)
	if err != nil {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}

	return false
}

// IsSuperAdmin reports whether the user is the built-in administrator.
func (s *MenuService) IsSuperAdmin(user *models.SysUser) bool {
	if user == nil {
		return false
	}
	if user.UserID == 1 {
		return true
	}

	roleIDs, err := s.roleIDs(user.UserID)
	if err != nil {
		return false
	}
	if len(roleIDs) == 0 {
		return false
	}

	var roles []models.SysRole
	if err := facades.Orm().Query().WhereIn("role_id", uintSliceToAny(roleIDs)).Get(&roles); err != nil {
		return false
	}
	for _, role := range roles {
		if role.RoleKey == "admin" {
			return true
		}
	}

	return false
}

// HomeMenuID returns the menu id that should be flagged as the app home
// (the first top-level menu of type C).
func (s *MenuService) HomeMenuID(menus []models.SysMenu) uint {
	var homeID uint
	for _, menu := range menus {
		if menu.ParentID == 0 && menu.MenuType == models.MenuTypeMenu && menu.Status == "0" {
			if homeID == 0 || menu.OrderNum < homeIDOrder(menus, homeID) {
				homeID = menu.MenuID
			}
		}
	}

	return homeID
}

func homeIDOrder(menus []models.SysMenu, id uint) int {
	for _, menu := range menus {
		if menu.MenuID == id {
			return menu.OrderNum
		}
	}

	return 1 << 30
}

func (s *MenuService) roleIDs(userID uint) ([]uint, error) {
	var bindings []models.SysUserRole
	if err := facades.Orm().Query().Where("user_id", userID).Get(&bindings); err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(bindings))
	for _, b := range bindings {
		ids = append(ids, b.RoleID)
	}

	return ids, nil
}

func uintSliceToAny(in []uint) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = v
	}

	return out
}
