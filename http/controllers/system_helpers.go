package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"

	"github.com/khs1001/goravel-admin/support"
)

// 分页默认值（对齐 owl-admin AdminService::list）。
const (
	defaultPerPage = 20
	maxPerPage     = 100
)

// pageParams 解析 page/perPage 查询参数。
func pageParams(ctx contractshttp.Context) (page, perPage int) {
	page = ctx.Request().QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	perPage = ctx.Request().QueryInt("perPage", defaultPerPage)
	if perPage < 1 {
		perPage = defaultPerPage
	}
	if perPage > maxPerPage {
		perPage = maxPerPage
	}

	return page, perPage
}

// listEnvelope 列表数据响应（amis CRUD 消费 {items, total}）。
func listEnvelope(items any, total int64) support.Envelope {
	return support.Success(map[string]any{"items": items, "total": total}, "")
}

// isGetData 请求是否为列表数据（对齐 CheckActionTrait::actionOfGetData）。
func isGetData(ctx contractshttp.Context) bool {
	return ctx.Request().Query("_action", "") == "getData"
}

// idList 解析批量删除的 id 列表（query `ids` 逗号分隔）。
func idList(ctx contractshttp.Context) []int {
	raw := ctx.Request().Query("ids", "")
	if raw == "" {
		return nil
	}

	out := []int{}
	for _, part := range strings.Split(raw, ",") {
		if id, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && id > 0 {
			out = append(out, id)
		}
	}

	return out
}

// intSlice 从任意输入提取 int 数组（amis extractValue 的 roles/permissions/menus）。
func intSlice(value any) []int {
	switch v := value.(type) {
	case []any:
		out := make([]int, 0, len(v))
		for _, item := range v {
			switch n := item.(type) {
			case float64:
				out = append(out, int(n))
			case string:
				if id, err := strconv.Atoi(n); err == nil {
					out = append(out, id)
				}
			case int:
				out = append(out, n)
			}
		}

		return out
	case []int:
		return v
	case string:
		var arr []int
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
		if id, err := strconv.Atoi(v); err == nil {
			return []int{id}
		}
	case float64:
		return []int{int(v)}
	}

	return nil
}

// strSlice 提取字符串数组（http_method 多选）。
func strSlice(value any) []string {
	switch v := value.(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}

		return out
	case []string:
		return v
	case string:
		if v == "" {
			return nil
		}

		return []string{v}
	}

	return nil
}
