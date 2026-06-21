package mysql

// ServiceCategory 服务分类表 - 服务项目分类（剪发/染发/烫发等）
type ServiceCategory struct {
	BaseModel
	Name      string `json:"name" gorm:"type:varchar(50);not null;comment:分类名称"`
	IconURL   string `json:"icon_url" gorm:"type:text;comment:分类图标URL"`
	SortOrder int    `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`
	Status    int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1启用"`

	// 关联
	ServiceItems []ServiceItem `json:"service_items,omitempty" gorm:"foreignKey:CategoryID"`
}

func (ServiceCategory) TableName() string {
	return "service_categories"
}

// ServiceItem 服务项目表 - 具体服务项目（男士剪发/染发套餐等）
type ServiceItem struct {
	BaseModel
	// 归属
	StoreID     *uint64 `json:"store_id,omitempty" gorm:"index;comment:门店ID(NULL=通用项目，所有门店可用)"`
	CategoryID  *uint64 `json:"category_id,omitempty" gorm:"index;comment:分类ID"`

	// 基本信息
	Name        string  `json:"name" gorm:"type:varchar(100);not null;comment:服务名称"`
	Description string  `json:"description" gorm:"type:text;comment:服务详细描述"`
	ImageURL    string  `json:"image_url" gorm:"type:text;comment:服务图片URL"`

	// 价格与时长
	OriginalPrice float64 `json:"original_price" gorm:"type:decimal(10,2);not null;comment:原价(元)"`
	Price         float64 `json:"price" gorm:"type:decimal(10,2);not null;index;comment:现价(元)"`
	Duration       int     `json:"duration" gorm:"type:int;default:30;comment:服务时长(分钟)"`

	// 属性
	TargetGender int8 `json:"target_gender" gorm:"type:tinyint;default:0;comment:适用性别: 0通用 1男 2女"`
	SortOrder    int   `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`
	IsHot        bool  `json:"is_hot" gorm:"type:tinyint;default:false;index;comment:是否热门推荐"`
	Status       int8  `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0下架 1上架"`

	// 关联关系
	Store         *Store            `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Category      *ServiceCategory  `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	OrderItems    []OrderItem        `json:"-" gorm:"foreignKey:ServiceItemID"` // 订单项关联
	AppointmentServices json.RawMessage `json:"-" gorm:"-"` // 预约-服务中间表
}

func (ServiceItem) TableName() string {
	return "service_items"
}

// GetDiscount 获取折扣力度（百分比）
func (s *ServiceItem) GetDiscount() float64 {
	if s.OriginalPrice > 0 && s.Price < s.OriginalPrice {
		return (1 - s.Price/s.OriginalPrice) * 100
	}
	return 0
}
