package pagination

import (
	"math"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Params 分页参数（从前端Query参数解析）
type Params struct {
	Page     int    `json:"page"`              // 当前页码（默认1）
	PerPage  int    `json:"per_page"`          // 每页数量（默认20，最大100）
	SortBy   string `json:"sort_by"`           // 排序字段
	SortDir  string `json:"sort_dir"`          // 排序方向: asc/desc
	Keyword   string `json:"keyword"`          // 搜索关键词
	Filters  map[string]string `json:"filters"` // 额外筛选条件
}

// Result 分页结果（包含数据和元信息）
type Result struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`      // 总记录数
	Page     int         `json:"page"`       // 当前页码
	PerPage  int         `json:"per_page"`   // 每页数量
	TotalPages int       `json:"total_pages"` // 总页数
}

// Meta 转换为response.Meta格式（用于API响应）
func (r *Result) ToMeta() map[string]interface{} {
	return map[string]interface{}{
		"page":        r.Page,
		"per_page":    r.PerPage,
		"total":       r.Total,
		"total_pages": r.TotalPages,
	}
}

// Parse 从URL Query参数解析分页参数
func Parse(query url.Values) *Params {
	p := &Params{
		Page:    1,
		PerPage: 20, // 默认每页20条
		SortBy:  "id",
		SortDir: "desc",
	}

	// 页码
	if page := query.Get("page"); page != "" {
		if v, err := strconv.Atoi(page); err == nil && v > 0 {
			p.Page = v
		}
	}

	// 每页数量
	if perPage := query.Get("per_page"); perPage != "" {
		if v, err := strconv.Atoi(perPage); err == nil && v > 0 && v <= 100 { // 最大100条防滥用
			p.PerPage = v
		}
	}

	// 排序字段
	if sortBy := query.Get("sort_by"); sortBy != "" {
		p.SortBy = sanitizeField(sortBy)
	}

	// 排序方向
	if sortDir := query.Get("sort_dir"); sortDir != "" {
		sortDir = strings.ToLower(sortDir)
		if sortDir == "asc" || sortDir == "desc" {
			p.SortDir = sortDir
		}
	}

	// 搜索关键词
	if keyword := query.Get("keyword"); keyword != "" {
		p.Keyword = strings.TrimSpace(keyword)
	}

	return p
}

// ApplyToGORM 将分页参数应用到GORM查询
func (p *Params) ApplyToGORM(db *gorm.DB) (*gorm.DB, *Result) {
	var total int64

	// 1. 先执行Count查询获取总数
	countDB := db.Session(&gorm.Session{})
	if p.Keyword != "" {
		// TODO: 根据具体模型添加搜索条件（需在Service层处理）
	}
	countModel := countDB.Model(nil) // 实际使用时应传入具体模型
	countModel.Count(&total)

	// 2. 计算偏移量和总页数
	offset := (p.Page - 1) * p.PerPage
	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	// 3. 应用排序和分页
	queryDB := db.Order(p.SortBy + " " + p.SortDir).Offset(offset).Limit(p.PerPage)

	// 构建结果对象
	result := &Result{
		Total:      total,
		Page:       p.Page,
		PerPage:    p.PerPage,
		TotalPages: totalPages,
	}

	return queryDB, result
}

// sanitizeField 清理排序字段名，防止SQL注入
func sanitizeField(field string) string {
	// 白名单方式：只允许字母数字下划线的字段名
	allowedChars := make([]rune, 0)
	for _, r := range field {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			allowedChars = append(allowedChars, r)
		}
	}
	return string(allowedChars)
}
