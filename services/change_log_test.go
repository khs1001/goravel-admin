package services

import (
	"testing"

	"github.com/khs1001/goravel-admin/models"
)

func TestDiffFields(t *testing.T) {
	oldPosition := &models.AdminPosition{Name: "负责人", CustomOrder: 0, Remark: ptr("旧备注")}
	newPosition := &models.AdminPosition{Name: "负责人-新", CustomOrder: 3, Remark: ptr("新备注")}

	changes := DiffFields(oldPosition, newPosition, []string{"name", "custom_order", "remark"})

	if len(changes) != 3 {
		t.Fatalf("expect 3 changes, got %d: %+v", len(changes), changes)
	}

	// 字段按名称排序：custom_order < name < remark
	if changes[0].Field != "custom_order" || changes[0].Old != float64(0) || changes[0].New != float64(3) {
		t.Errorf("changes[0] = %+v, want custom_order 0 -> 3", changes[0])
	}

	nameChange := changes[1]
	if nameChange.Old != "负责人" || nameChange.New != "负责人-新" {
		t.Errorf("name change = %v -> %v, want 负责人 -> 负责人-新", nameChange.Old, nameChange.New)
	}
}

func TestDiffFieldsNoChange(t *testing.T) {
	a := &models.AdminRole{Name: "editor", Slug: "editor"}
	b := &models.AdminRole{Name: "editor", Slug: "editor"}

	if changes := DiffFields(a, b, []string{"name", "slug"}); len(changes) != 0 {
		t.Errorf("expect no changes, got %+v", changes)
	}
}

func TestChangesPayload(t *testing.T) {
	if got := changesPayload(nil); got != "" {
		t.Errorf("changesPayload(nil) = %q, want empty", got)
	}

	got := changesPayload([]FieldChange{
		{Field: "name", Old: "a", New: "b"},
		{Field: "sort", Old: float64(1), New: float64(2)},
	})
	expected := `{"name":{"new":"b","old":"a"},"sort":{"new":2,"old":1}}`
	if got != expected {
		t.Errorf("changesPayload() = %s, want %s", got, expected)
	}
}

func ptr(s string) *string { return &s }
