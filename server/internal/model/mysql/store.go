package mysql

import (
	"encoding/json"
)

// Store 门店表 - 全国连锁门店信息
type Store struct {
	BaseModel
	// 组织归属
	OrgID uint64 `json:"org_id" gorm:"index;comment:组织架构ID"`

	// 基本信息
	Name        string  `json:"name" gorm:"type:varchar(100);not null;comment:门店名称"`
	LogoURL     string  `json:"logo_url" gorm:"type:text;comment:Logo图片URL"`
	CoverImages string  `json:"cover_images" gorm:"type:json;comment:封面图片列表(JSON数组)"`

	// 地址位置（腾讯地图坐标）
	Province string   `json:"province" gorm:"type:varchar(50);index;comment:省"`
	City     string   `json:"city" gorm:"type:varchar(50);index;comment:市"`
	District string   `json:"district" gorm:"type:varchar(50);index;comment:区"`
	Address  string   `json:"address" gorm:"type:varchar(200);not null;comment:详细地址"`
	Latitude float64  `json:"latitude" gorm:"type:decimal(10,7);index;comment:纬度"`
	Longitude float64 `json:"longitude" gorm:"type:decimal(10,7);index;comment:经度"`

	// 联系与营业
	Phone      string `json:"phone" gorm:"type:varchar(20);comment:联系电话"`
	OpenTime   string `json:"open_time" gorm:"type:varchar(20);default:'09:00';comment:营业开始时间"`
	CloseTime  string `json:"close_time" gorm:"type:varchar(20);default:'21:00';comment:营业结束时间"`
	ParkingInfo string `json:"parking_info" gorm:"type:varchar(200);comment:停车信息"`

	// 描述与评分
	Description string  `json:"description" gorm:"type:text;comment:门店描述/环境介绍"`
	AvgPrice    float64 `json:"avg_price" gorm:"type:decimal(10,2);default:0.00;comment:人均消费价格"`
	Rating      float32 `json:"rating" gorm:"type:decimal(2,1);default:5.0;comment:综合评分(1-5)"`
	ReviewCount int     `json:"review_count" gorm:"type:int;default:0;comment:评价数量"`
	StarLevel   int8    `json:"star_level" gorm:"type:tinyint;default:1;comment:星级等级(1-5)"`

	// 状态与标签
	Status    int8 `json:"status" gorm:"type:tinyint;default:1;index;comment:状态: 0停业 1营业 2筹备中 3已关闭"`
	IsFeatured bool `json:"is_featured" gorm:"type:tinyint;default:false;comment:是否推荐门店"`

	// 关联关系
	Photos        []StorePhoto       `json:"photos,omitempty" gorm:"foreignKey:StoreID"`         // 门店照片
	Stylists      []Stylist          `json:"stylists,omitempty" gorm:"foreignKey:StoreID"`        // 理发师团队
	ServiceItems  []ServiceItem       `json:"service_items,omitempty" gorm:"foreignKey:StoreID"`   // 服务项目
	Appointments  []Appointment       `json:"appointments,omitempty" gorm:"foreignKey:StoreID"`    // 预约记录
	Orders        []Order             `json:"orders,omitempty" gorm:"foreignKey:StoreID"`          // 订单
	Employees     []Employee          `json:"employees,omitempty" gorm:"foreignKey:StoreID"`       // 员工
	Inventories   []Inventory         `json:"inventories,omitempty" gorm:"foreignKey:StoreID"`     // 库存
}

func (Store) TableName() string {
	return "stores"
}

// GetCoverImageList 获取封面图片列表（JSON解析）
func (s *Store) GetCoverImageList() []string {
	var images []string
	if s.CoverImages != "" {
		json.Unmarshal([]byte(s.CoverImages), &images)
	}
	return images
}

// IsOpenNow 检查当前是否营业中（简化版，实际需结合时区）
func (s *Store) IsOpenNow() bool {
	return s.Status == 1 // 实际应判断当前时间段
}

// StorePhoto 门店照片表 - 支持多张照片展示
type StorePhoto struct {
	BaseModel
	StoreID uint64 `json:"store_id" gorm:"not null;index;comment:门店ID"`
	URL     string `json:"url" gorm:"type:text;not null;comment:照片URL"`
	Type    int8   `json:"type" gorm:"type:tinyint;default:0;comment:类型: 0环境照 1门头照 2内部设施 3团队照"`
	SortOrder int   `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`

	// 关联
	Store *Store `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

func (StorePhoto) TableName() string {
	return "store_photos"
}
