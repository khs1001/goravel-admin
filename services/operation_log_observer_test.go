package services

import (
	"testing"

	ormmocks "github.com/goravel/framework/mocks/database/orm"
	"github.com/stretchr/testify/mock"
)

func TestBuildChanges(t *testing.T) {
	event := ormmocks.NewEvent(t)
	// username 与 name 脏、remark 干净；username 无原始值
	event.On("IsDirty", "username").Return(true)
	event.On("IsDirty", "name").Return(true)
	event.On("IsDirty", "remark").Return(false)
	event.On("GetOriginal", "username").Return("old-user")
	event.On("GetAttribute", "username").Return("new-user")
	event.On("GetOriginal", "name").Return(nil)
	event.On("GetAttribute", "name").Return("新名称")

	got := buildChanges(event, []string{"username", "name", "remark"})

	expected := `{"name":{"new":"新名称","old":null},"username":{"new":"new-user","old":"old-user"}}`
	if got != expected {
		t.Errorf("buildChanges() = %s, want %s", got, expected)
	}
}

func TestBuildChangesNoDirty(t *testing.T) {
	event := ormmocks.NewEvent(t)
	event.On("IsDirty", mock.Anything).Return(false)

	if got := buildChanges(event, []string{"name", "slug"}); got != "" {
		t.Errorf("buildChanges() = %q, want empty", got)
	}
}

func TestBuildChangesEmptyColumns(t *testing.T) {
	event := ormmocks.NewEvent(t)

	if got := buildChanges(event, nil); got != "" {
		t.Errorf("buildChanges(nil) = %q, want empty", got)
	}
}
