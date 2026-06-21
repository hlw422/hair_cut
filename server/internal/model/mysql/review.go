package mysql

import (
	"time"
)

// Review 评价表 - 用户对服务/理发师的评价
type Review struct {
	BaseModel
	// 关联信息
	UserID    uint64 `json:"user_id" gorm:"not null;index;comment:用户ID"`
	OrderID   uint64 `json:"order_id" gorm:"uniqueIndex;not null;comment:订单ID(一个订单只能评价一次)"`
	StoreID   uint64 `json:"store_id" gorm:"not null;index;comment:门店ID"`
	StylistID uint64 `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`

	// 评分（1-5分）
	Rating        int8 `json:"rating" gorm:"type:tinyint;not null;comment:综合评分(1-5)"`
	ServiceRating int8 `json:"service_rating" gorm:"type:tinyint;default:5;comment:服务评分"`
	EnvironmentRating int8 `json:"environment_rating" gorm:"type:tinyint;default:5;comment:环境评分"`
	StylistRating int8 `json:"stylist_rating" gorm:"type:tinyint;default:5;comment:理发师技术评分"`

	// 内容
	Content string `json:"content" gorm:"type:text;comment:评价文字内容"`
	IsAnonymous bool `json:"is_anonymous" gorm:"type:tinyint;default:false;comment:是否匿名评价"`

	// 状态与审核
	Status     int8       `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0待审核 1已显示 2已隐藏"`
	Reply      string     `json:"reply,omitempty" gorm:"type:text;comment:商家回复内容"`
	RepliedAt  *time.Time `json:"replied_at,omitempty" gorm:"comment:回复时间"`

	// 统计
	LikeCount int `json:"like_count" gorm:"type:int;default:0;comment:点赞数"`

	// 关联关系
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Store        *Store         `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Stylist      *Stylist       `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
	Order        *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	MediaFiles   []ReviewMedia  `json:"media_files,omitempty" gorm:"foreignKey:ReviewID"` // 评价媒体文件
}

func (Review) TableName() string {
	return "reviews"
}

// ReviewMedia 评价多媒体附件表 - 支持图片和视频
type ReviewMedia struct {
	BaseModel
	ReviewID uint64 `json:"review_id" gorm:"not null;index;comment:评价ID"`
	Type     int8   `json:"type" gorm:"type:tinyint;default:0;comment:类型: 0图片 1视频"`
	URL      string `json:"url" gorm:"type:text;not null;comment:媒体URL(MinIO)"`
	CoverURL string `json:"cover_url" gorm:"type:text;comment:视频封面图URL"`
	SortOrder int   `json:"sort_order" gorm:"type:int;default:0;comment:排序"`

	// 关联
	Review *Review `json:"review,omitempty" gorm:"foreignKey:ReviewID"`
}

func (ReviewMedia) TableName() string {
	return "review_media"
}
