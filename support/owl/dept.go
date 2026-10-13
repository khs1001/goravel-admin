package owl

import "github.com/khs1001/goravel-admin/models"

// BuildDeptTree assembles the RuoYi `sys_dept` rows into a tree, skipping
// soft-deleted rows.
func BuildDeptTree(list []models.SysDept, parentID uint) []models.SysDept {
	tree := make([]models.SysDept, 0)
	for i := range list {
		if list[i].ParentID != parentID || list[i].DelFlag == "2" {
			continue
		}

		list[i].Children = BuildDeptTree(list, list[i].DeptID)
		tree = append(tree, list[i])
	}

	return tree
}

// BuildDeptTreeFiltered assembles a tree from an arbitrary subset of
// departments. A node whose parent is not part of the subset is treated as a
// root, which is required when a data scope returns only part of the tree.
func BuildDeptTreeFiltered(list []models.SysDept) []models.SysDept {
	present := make(map[uint]bool, len(list))
	for _, dept := range list {
		present[dept.DeptID] = true
	}

	roots := make([]models.SysDept, 0)
	for i := range list {
		if list[i].DelFlag == "2" || present[list[i].ParentID] {
			continue
		}

		list[i].Children = buildDeptChildren(list, list[i].DeptID)
		roots = append(roots, list[i])
	}

	return roots
}

func buildDeptChildren(list []models.SysDept, parentID uint) []models.SysDept {
	result := make([]models.SysDept, 0)
	for i := range list {
		if list[i].ParentID != parentID || list[i].DelFlag == "2" {
			continue
		}

		list[i].Children = buildDeptChildren(list, list[i].DeptID)
		result = append(result, list[i])
	}

	return result
}
