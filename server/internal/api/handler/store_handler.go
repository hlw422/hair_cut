package handler

import (
	"haircut-server/pkg/pagination"
	"haircut-server/pkg/response"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StoreHandler 门店相关API处理器
type StoreHandler struct{}

// NearbyStores 附近门店列表（腾讯地图集成）
// GET /api/v1/stores/nearby?lat=31.2304&lng=121.4737&radius=5000&page=1&per_page=20
func (h *StoreHandler) NearbyStores(c *gin.Context) {
	// 1. 解析查询参数
	lat, _ := strconv.ParseFloat(c.DefaultQuery("lat", "0"), 64)
	lng, _ := strconv.ParseFloat(c.DefaultQuery("lng", "0"), 64)
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "5000"), 64)

	if lat == 0 || lng == 0 {
		response.BadRequest(c, "请提供有效的经纬度坐标")
		return
	}

	// 2. 解析分页参数
	p := pagination.Parse(c.Request.URL.Query())

	// 3. 调用Service层查询（使用Haversine公式计算距离或MySQL空间函数）
	// TODO: storeService.GetNearbyStores(lat, lng, radius, p)
	
	// Mock数据（实际应从数据库查询）
	stores := []map[string]interface{}{
		{
			"id":          1,
			"name":        "HairCut 精品沙龙（静安店）",
			"logo_url":    "https://placehold.co/100x100/C8A882/white?text=HC",
			"address":     "上海市静安区南京西路123号",
			"distance":    850.5,
			"rating":      4.8,
			"review_count": 256,
			"avg_price":   168.00,
			"is_open":     true,
			"tags":        []string{"环境优雅", "明星理发师"},
		},
		{
			"id":          2,
			"name":        "HairCut 造型中心（徐汇店）",
			"logo_url":    "https://placehold.co/100x100/C8A882/white?text=HC",
			"address":     "上海市徐汇区淮海中路456号",
			"distance":    2350.0,
			"rating":      4.9,
			"review_count": 389,
			"avg_price":   198.00,
			"is_open":     true,
			"tags":        []string{"人气TOP1", "免费停车"},
		},
	}

	total := int64(25) // 假设总共25家门店在范围内

	// 构建分页结果
	result := pagination.Result{
		List:       stores,
		Total:      total,
		Page:       p.Page,
		PerPage:    p.PerPage,
		TotalPages: int(math.Ceil(float64(total) / float64(p.PerPage))),
	}

	response.PagedResponse(c, stores, response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      total,
		TotalPages: result.TotalPages,
	})
}

// GetStoreDetail 门店详情
// GET /api/v1/stores/:id
func (h *StoreHandler) GetStoreDetail(c *gin.Context) {
	storeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || storeID == 0 {
		response.BadRequest(c, "无效的门店ID")
		return
	}

	// TODO: 调用Service获取完整门店信息（含照片、服务、理发师、评价）
	
	storeDetail := map[string]interface{}{
		"id":            storeID,
		"name":          "HairCut 精品沙龙（静安店）",
		"cover_images": []string{"https://placehold.co/750x400/C8A882/white?text=Store+Cover"},
		"address":       "上海市静安区南京西路123号",
		"phone":         "021-12345678",
		"open_time":     "09:30",
		"close_time":    "21:30",
		"description":   "坐落于繁华南京路商圈的高端美发沙龙，拥有资深造型团队...",
		"avg_price":     168.00,
		"rating":        4.8,
		"review_count":  256,
		"star_level":    5,
		"parking_info":  "提供地下停车场，前2小时免费",
		"is_open":       true,

		// 服务项目（示例）
		"services": []map[string]interface{}{
			{"id": 1, "name": "首席设计师剪发", "price": 198.0, "original_price": 268.0, "duration": 60},
			{"id": 2, "name": "资深染发套餐", "price": 388.0, "original_price": 528.0, "duration": 120},
			{"id": 3, "name": "头皮护理SPA", "price": 268.0, "duration": 45},
		},

		// 明星理发师（示例）
		"stylists": []map[string]interface{}{
			{"id": 1, "name": "Kevin老师", "avatar": "", "title": "艺术总监", "experience_years": 12, "rating": 4.9},
			{"id": 2, "name": "Linda老师", "avatar": "", "title": "首席造型师", "experience_years": 8, "rating": 4.8},
		},

		// 用户评价（示例）
		"reviews": []map[string]interface{}{
			{"user": "匿名用户", "rating": 5, "content": "Kevin老师技术超棒！剪的发型很适合我！", "date": "2024-01-15"},
			{"user": "小**", "rating": 5, "content": "环境很好，服务贴心，会再来~", "date": "2024-01-10"},
		},
	}

	response.Success(c, storeDetail)
}

// ListStores 门店列表（支持排序和筛选）
// GET /api/v1/stores?city=上海&sort_by=rating&sort_dir=desc&keyword=
func (h *StoreHandler) ListStores(c *gin.Context) {
	p := pagination.Parse(c.Request.URL.Query())

	city := c.Query("city")
	sortBy := c.DefaultQuery("sort_by", "rating")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	// 验证排序字段白名单
	allowedSortFields := map[string]bool{"rating": true, "distance": true, "price": true, "created_at": true}
	if !allowedSortFields[sortBy] {
		sortBy = "rating" // 默认按评分排序
	}

	// TODO: 调用Service层查询
	_ = city
	_ = sortDir

	// 返回分页结果
	mockStores := []interface{}{}
	meta := response.Meta{
		Page:       p.Page,
		PerPage:    p.PerPage,
		Total:      50,
		TotalPages: 3,
	}

	response.PagedResponse(c, mockStores, meta)
}
