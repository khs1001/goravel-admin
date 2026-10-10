package services

import (
	"reflect"
	"testing"

	"github.com/khs1001/goravel-admin/models"
)

func perm(id uint, parentID int, name string) models.AdminPermission {
	p := models.AdminPermission{Name: name, ParentID: parentID}
	p.ID = id

	return p
}

// buildPermissionTree：两级嵌套、兄弟顺序按输入序、孤儿节点不出现。
func TestBuildPermissionTree(t *testing.T) {
	permissions := []models.AdminPermission{
		perm(2, 0, "角色管理"),
		perm(1, 0, "用户管理"),
		perm(3, 1, "用户查看"),
		perm(4, 1, "用户编辑"),
		perm(5, 3, "深层子项"),
		perm(6, 99, "孤儿（父不存在）"),
	}

	tree := buildPermissionTree(permissions, 0)

	if len(tree) != 2 {
		t.Fatalf("root count should be 2, got %d", len(tree))
	}
	// 输入序即兄弟序（DB 查询已按 custom_order 排序）
	if tree[0].Name != "角色管理" || tree[1].Name != "用户管理" {
		t.Fatalf("sibling order mismatch: %v", tree)
	}

	userNode := tree[1]
	if len(userNode.Children) != 2 {
		t.Fatalf("user node should have 2 children")
	}
	if userNode.Children[0].Name != "用户查看" {
		t.Fatalf("child order mismatch")
	}
	if len(userNode.Children[0].Children) != 1 || userNode.Children[0].Children[0].Name != "深层子项" {
		t.Fatalf("grandchild missing")
	}

	// 孤儿节点不应出现在树中
	for _, node := range tree {
		if node.Name == "孤儿（父不存在）" {
			t.Fatalf("orphan node should not appear in tree")
		}
	}
}

// buildPermissionTree：空输入返回空切片（非 nil，保证 JSON 输出 []）。
func TestBuildPermissionTreeEmpty(t *testing.T) {
	tree := buildPermissionTree(nil, 0)
	if tree == nil {
		t.Fatalf("should return empty slice, not nil (json must be [])")
	}
	if len(tree) != 0 {
		t.Fatalf("length should be 0")
	}
}

// toAnySlice：WhereIn 参数转换。
func TestToAnySlice(t *testing.T) {
	got := toAnySlice([]int{1, 2, 3})
	if !reflect.DeepEqual(got, []any{1, 2, 3}) {
		t.Fatalf("toAnySlice mismatch: %v", got)
	}

	if out := toAnySlice(nil); len(out) != 0 {
		t.Fatalf("nil input should yield empty slice")
	}
}
