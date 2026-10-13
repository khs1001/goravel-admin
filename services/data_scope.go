package services

import (
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/models"
)

// Scope describes the data a user may access, resolved from the `data_scope` of
// their roles (following RuoYi semantics).
type Scope struct {
	// Unrestricted means the user can see all data (scope "1" / super admin).
	Unrestricted bool
	// DeptIDs is the set of department ids the user may access.
	DeptIDs []uint
	// SelfOnly means the user may only access their own rows (scope "5").
	SelfOnly bool
	// UserID is the current user id, used for the "self only" clause.
	UserID uint
}

// DataScopeService resolves data-scope restrictions from a user's roles.
type DataScopeService struct{}

// NewDataScopeService creates a DataScopeService.
func NewDataScopeService() *DataScopeService {
	return &DataScopeService{}
}

// Resolve computes the data scope of a user.
//
// Multiple roles are combined as a union; if any role grants "all data" the
// result is unrestricted.
func (s *DataScopeService) Resolve(user *models.SysUser) Scope {
	if user == nil {
		return Scope{}
	}
	if NewMenuService().IsSuperAdmin(user) {
		return Scope{Unrestricted: true, UserID: user.UserID}
	}

	var bindings []models.SysUserRole
	if err := facades.Orm().Query().Where("user_id", user.UserID).Get(&bindings); err != nil {
		return Scope{UserID: user.UserID}
	}

	roleIDs := make([]any, 0, len(bindings))
	for _, binding := range bindings {
		roleIDs = append(roleIDs, binding.RoleID)
	}
	if len(roleIDs) == 0 {
		return Scope{UserID: user.UserID}
	}

	var roles []models.SysRole
	if err := facades.Orm().Query().WhereIn("role_id", roleIDs).Get(&roles); err != nil {
		return Scope{UserID: user.UserID}
	}

	deptSet := map[uint]bool{}
	selfOnly := false

	for _, role := range roles {
		switch role.DataScope {
		case "1":
			return Scope{Unrestricted: true, UserID: user.UserID}
		case "2":
			for _, deptID := range s.roleDeptIDs(role.RoleID) {
				deptSet[deptID] = true
			}
		case "3":
			deptSet[user.DeptID] = true
		case "4":
			deptSet[user.DeptID] = true
			for _, deptID := range s.descendantDeptIDs(user.DeptID) {
				deptSet[deptID] = true
			}
		case "5":
			selfOnly = true
		}
	}

	deptIDs := make([]uint, 0, len(deptSet))
	for deptID := range deptSet {
		if deptID != 0 {
			deptIDs = append(deptIDs, deptID)
		}
	}

	return Scope{DeptIDs: deptIDs, SelfOnly: selfOnly, UserID: user.UserID}
}

func (s *DataScopeService) roleDeptIDs(roleID uint) []uint {
	var bindings []models.SysRoleDept
	if err := facades.Orm().Query().Where("role_id", roleID).Get(&bindings); err != nil {
		return nil
	}

	ids := make([]uint, 0, len(bindings))
	for _, binding := range bindings {
		ids = append(ids, binding.DeptID)
	}

	return ids
}

// descendantDeptIDs returns the department subtree rooted at deptID (excluding
// deptID itself). The tree is resolved in Go to stay database-agnostic.
func (s *DataScopeService) descendantDeptIDs(deptID uint) []uint {
	if deptID == 0 {
		return nil
	}

	var depts []models.SysDept
	if err := facades.Orm().Query().Where("del_flag", "0").Get(&depts); err != nil {
		return nil
	}

	children := map[uint][]uint{}
	for _, dept := range depts {
		children[dept.ParentID] = append(children[dept.ParentID], dept.DeptID)
	}

	result := make([]uint, 0)
	queue := append([]uint{}, children[deptID]...)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		queue = append(queue, children[current]...)
	}

	return result
}
