package services

import (
	"encoding/json"
	"fmt"
	"sort"

	contractshttp "github.com/goravel/framework/contracts/http"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"

	"github.com/khs1001/goravel-admin/jobs"
	"github.com/khs1001/goravel-admin/models"
)

// FieldChange 单个字段的变更（old/new 为渲染用原始值）。
type FieldChange struct {
	Field string `json:"field"`
	Old   any    `json:"old"`
	New   any    `json:"new"`
}

// DiffFields 对比同类型模型的指定列，返回变更数组。
// 列名与 JSON 序列化键对齐（模型 json tag 即 gorm column 名），
// 经 map 对比实现，不依赖 ORM 事件的 Original 快照（Save 模式下拿不到）。
func DiffFields(oldModel, newModel any, columns []string) []FieldChange {
	oldMap, err1 := marshalModel(oldModel)
	newMap, err2 := marshalModel(newModel)
	if err1 != nil || err2 != nil {
		return nil
	}

	changes := make([]FieldChange, 0, len(columns))
	for _, col := range columns {
		if !jsonEqual(oldMap[col], newMap[col]) {
			changes = append(changes, FieldChange{Field: col, Old: oldMap[col], New: newMap[col]})
		}
	}

	sort.Slice(changes, func(i, j int) bool { return changes[i].Field < changes[j].Field })

	return changes
}

// marshalModel 结构体序列化为 map（json tag 命名）；指针/时间等由 json 统一处理。
func marshalModel(model any) (map[string]any, error) {
	raw, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	out := map[string]any{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// jsonEqual 值相等判断（数字类型经 json 反序列化后均为 float64，直接 == 即可）。
func jsonEqual(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	return string(ab) == string(bb)
}

// changesPayload 把变更数组转为 {col: {old, new}} JSON（与日志 changes 列格式一致）。
func changesPayload(changes []FieldChange) string {
	if len(changes) == 0 {
		return ""
	}

	out := map[string]map[string]any{}
	for _, change := range changes {
		out[change.Field] = map[string]any{"old": change.Old, "new": change.New}
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return ""
	}

	return string(encoded)
}

// DispatchChangeLog 落一条带字段级 diff 的 updated 操作日志（经 Queue 异步）。
// 由控制器在更新保存后调用（控制器天然持有更新前后状态）。
func DispatchChangeLog(ctx contractshttp.Context, resource, resourceID string, changes []FieldChange) {
	payload := changesPayload(changes)
	if payload == "" {
		return
	}

	// 操作人：控制器栈内可拿登录用户；失败降级为 system
	userID, username := 0, "system"
	if ctx != nil {
		user := &models.AdminUser{}
		if err := facades.Auth(ctx).Guard("admin").User(user); err == nil && user.ID > 0 {
			userID, username = int(user.ID), user.Username
		}
	}

	encoded, err := json.Marshal(jobs.OperationLogPayload{
		UserID:     userID,
		Username:   username,
		Action:     "updated",
		Resource:   resource,
		ResourceID: resourceID,
		Changes:    payload,
	})
	if err != nil {
		facades.Log().Warningf("[admin] marshal change log failed: %v", err)

		return
	}

	if err := facades.Queue().Job(&jobs.OperationLogJob{}, []contractsqueue.Arg{
		{Value: string(encoded), Type: "string"},
	}).Dispatch(); err != nil {
		facades.Log().Warningf("[admin] dispatch change log failed: %v", err)
	}
}

// resourceIDOf 主键转字符串（nil/0 视为未知）。
func resourceIDOf(id uint) string {
	if id == 0 {
		return ""
	}

	return fmt.Sprintf("%d", id)
}
