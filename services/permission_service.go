package services

import (
	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// SuperAdministratorSlug 超管角色 slug（对齐 owl-admin AdminRole::SuperAdministrator）。
const SuperAdministratorSlug = "administrator"

// PermissionService 权限服务：角色/权限判定、权限树、菜单可见性过滤。
// 逻辑对齐 owl-admin src/Services/AdminPermissionService.php 与 AdminUser 模型的角色方法。
type PermissionService struct{}

func NewPermissionService() *PermissionService { return &PermissionService{} }

// toAnySlice []int → []any（WhereIn 参数）。
func toAnySlice(ids []int) []any {
	out := make([]any, 0, len(ids))
	for _, id := range ids {
		out = append(out, id)
	}

	return out
}

// IsAdministrator 用户是否拥有超管角色。
func (r *PermissionService) IsAdministrator(userID int) bool {
	var roleUsers []models.AdminRoleUsers
	if err := facades.Orm().Query().Where("user_id", userID).Find(&roleUsers); err != nil {
		facades.Log().Errorf("[admin] load user roles failed: %v", err)

		return false
	}
	if len(roleUsers) == 0 {
		return false
	}

	roleIDs := make([]int, 0, len(roleUsers))
	for _, ru := range roleUsers {
		roleIDs = append(roleIDs, ru.RoleID)
	}

	var count int64
	var err error
	count, err = facades.Orm().Query().
		Model(&models.AdminRole{}).
		Where("slug", SuperAdministratorSlug).
		WhereIn("id", toAnySlice(roleIDs)).
		Count()
	if err != nil {
		facades.Log().Errorf("[admin] check administrator failed: %v", err)

		return false
	}

	return count > 0
}

// UserPermissionIDs 用户经角色拥有的全部权限 ID。
func (r *PermissionService) UserPermissionIDs(userID int) []int {
	var roleUsers []models.AdminRoleUsers
	if err := facades.Orm().Query().Where("user_id", userID).Find(&roleUsers); err != nil {
		facades.Log().Errorf("[admin] load user roles failed: %v", err)

		return nil
	}
	if len(roleUsers) == 0 {
		return nil
	}

	roleIDs := make([]int, 0, len(roleUsers))
	for _, ru := range roleUsers {
		roleIDs = append(roleIDs, ru.RoleID)
	}

	var rolePermissions []models.AdminRolePermissions
	if err := facades.Orm().Query().WhereIn("role_id", toAnySlice(roleIDs)).Find(&rolePermissions); err != nil {
		facades.Log().Errorf("[admin] load role permissions failed: %v", err)

		return nil
	}

	seen := map[int]bool{}
	ids := make([]int, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		if !seen[rp.PermissionID] {
			seen[rp.PermissionID] = true
			ids = append(ids, rp.PermissionID)
		}
	}

	return ids
}

// VisibleMenus 当前用户可见的菜单。
//
// 返回 all=true 表示全部可见（超管或未开启权限过滤）；
// 否则返回可见菜单 ID 集合（权限经 admin_permission_menu 挂载的菜单）。
func (r *PermissionService) VisibleMenus(userID int) (all bool, ids map[int]bool) {
	if facades.Config().GetBool("admin.auth.permission", true) == false {
		return true, nil
	}

	if r.IsAdministrator(userID) {
		return true, nil
	}

	permissionIDs := r.UserPermissionIDs(userID)
	if len(permissionIDs) == 0 {
		return false, map[int]bool{}
	}

	var links []models.AdminPermissionMenu
	if err := facades.Orm().Query().WhereIn("permission_id", toAnySlice(permissionIDs)).Find(&links); err != nil {
		facades.Log().Errorf("[admin] load permission menus failed: %v", err)

		return false, map[int]bool{}
	}

	ids = make(map[int]bool, len(links))
	for _, link := range links {
		ids[link.MenuID] = true
	}

	return false, ids
}

// PermissionTreeNode 权限树节点（管理表单的 TreeSelect options 用）。
type PermissionTreeNode struct {
	ID       int                  `json:"id"`
	ParentID int                  `json:"parent_id"`
	Name     string               `json:"name"`
	Children []PermissionTreeNode `json:"children,omitempty"`
}

// Tree 全量权限树（按 custom_order 升序）。
func (r *PermissionService) Tree() []PermissionTreeNode {
	var permissions []models.AdminPermission
	if err := facades.Orm().Query().OrderBy("custom_order").Find(&permissions); err != nil {
		facades.Log().Errorf("[admin] load permissions failed: %v", err)

		return []PermissionTreeNode{}
	}

	return buildPermissionTree(permissions, 0)
}

func buildPermissionTree(permissions []models.AdminPermission, parentID int) []PermissionTreeNode {
	nodes := []PermissionTreeNode{}
	for _, p := range permissions {
		if p.ParentID != parentID {
			continue
		}

		nodes = append(nodes, PermissionTreeNode{
			ID:       int(p.ID),
			ParentID: p.ParentID,
			Name:     p.Name,
			Children: buildPermissionTree(permissions, int(p.ID)),
		})
	}

	return nodes
}

// SavePermissions 保存角色-权限关联（全量覆盖）。
func (r *PermissionService) SavePermissions(roleID int, permissionIDs []int) error {
	return facades.Orm().Transaction(func(tx orm.Query) error {
		if _, err := tx.Where("role_id", roleID).Delete(&models.AdminRolePermissions{}); err != nil {
			return err
		}

		if len(permissionIDs) == 0 {
			return nil
		}

		links := make([]models.AdminRolePermissions, 0, len(permissionIDs))
		for _, pid := range permissionIDs {
			links = append(links, models.AdminRolePermissions{RoleID: roleID, PermissionID: pid})
		}

		return tx.Create(&links)
	})
}
