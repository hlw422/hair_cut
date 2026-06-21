package mysql

import (
	"time"
)

// Message 消息表 - 系统消息与通知
type Message struct {
	BaseModel
	// 接收者
	UserID   uint64 `json:"user_id" gorm:"not null;index;comment:接收用户ID"`
	UserType int8   `json:"user_type" gorm:"type:tinyint;default:1;comment:用户类型: 1会员用户 2理发师 3店长"`

	// 消息内容
	Title     string `json:"title" gorm:"type:varchar(100);not null;comment:消息标题"`
	Content   string `json:"content" gorm:"type:text;not null;comment:消息内容"`
	Type      int8   `json:"type" gorm:"type:tinyint;not null;index;comment:类型: 1系统通知 2预约提醒 3支付通知 4活动推送 5营销消息"`

	// 关联业务（可选）
	RelatedType string  `json:"related_type,omitempty" gorm:"type:varchar(50);comment:关联业务类型(order/appointment/store等)"`
	RelatedID   *uint64 `json:"related_id,omitempty" gorm:"comment:关联业务ID"`

	// 状态
	IsRead    bool       `json:"is_read" gorm:"type:tinyint;default:false;comment:是否已读"`
	ReadAt    *time.Time `json:"read_at,omitempty" gorm:"comment:读取时间"`
	IsDeleted bool       `json:"is_deleted" gorm:"type:tinyint;default:false;comment:是否已删除(软删除)"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Message) TableName() string {
	return "messages"
}

// Notification 通知表 - WebSocket实时推送记录
type Notification struct {
	BaseModel
	// 接收者信息
	ReceiverID   uint64 `json:"receiver_id" gorm:"not null;index;comment:接收者ID"`
	ReceiverType int8   `json:"receiver_type" gorm:"type:tinyint;not null;comment:接收者类型: 1用户 2理发师 3店长 4管理员"`
	
	// 发送者信息（可为系统）
	SenderID   *uint64 `json:"sender_id,omitempty" gorm:"comment:发送者ID(空=系统)"`
	SenderType int8    `json:"sender_type" gorm:"type:tinyint;default:0;comment:发送者类型: 0系统 1用户 2理发师"`

	// 内容与分类
	Title   string `json:"title" gorm:"type:varchar(100);not null;comment:通知标题"`
	Content string `json:"content" gorm:"type:text;not null;comment:通知内容"`
	Category int8   `json:"category" gorm:"type:tinyint;not null;comment:分类: 1预约变更 2订单状态 3新粉丝 4系统公告 5营销"`

	// 业务关联
	BusinessType string  `json:"business_type" gorm:"type:varchar(50);index;comment:业务类型"`
	BusinessID   *uint64 `json:"business_id,omitempty" gorm:"comment:业务ID"`

	// 状态追踪
	Status     int8       `json:"status" gorm:"type:tinyint;default:0;comment:状态: 0待推送 1已推送 2已送达 3已阅读"`
	PushedAt   *time.Time `json:"pushed_at,omitempty" gorm:"comment:推送时间"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" gorm:"comment:送达时间"`
	ReadAt     *time.Time `json:"read_at,omitempty" gorm:"comment:阅读时间"`
}

func (Notification) TableName() string {
	return "notifications"
}

// FanRelation 粉丝关注表 - 用户关注理发师的社交关系
type FanRelation struct {
	BaseModel
	UserID    uint64 `json:"user_id" gorm:"not null;index;comment:用户ID(粉丝)"`
	StylistID uint64 `json:"stylist_id" gorm:"not null;index;comment:理发师ID(被关注)"`
	Status    int8   `json:"status" gorm:"type:tinyint;default:1;index;comment:状态: 0取关 1关注"`

	// 关联关系
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Stylist *Stylist `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
}

func (FanRelation) TableName() string {
	return "fan_relations"

// CustomerProfile 客户档案表 - 理发师维护的客户详细信息
type CustomerProfile struct {
	BaseModel
	// 归属
	StylistID uint64 `json:"stylist_id" gorm:"not null;index;comment:所属理发师ID"`
	UserID    uint64 `json:"user_id" gorm:"not null;uniqueIndex:idx_stylist_user;comment:客户用户ID"`

	// 偏好记录（理发师手动维护）
	HairStylePreference string `json:"hair_style_preference" gorm:"type:varchar(200);comment:发型偏好(如:短发/中分/偏分等)"`
	ColorPreference     string `json:"color_preference" gorm:"type:varchar(100);comment:染发颜色偏好"`
	AllergyInfo         string `json:"allergy_info" gorm:"type:varchar(200);comment:过敏史/皮肤敏感情况"`
	ScalpCondition      string `json:"scalp_condition" gorm:"type:varchar(100);comment:头皮状况(油性/干性/敏感)"`

	// 发型记录（最近一次）
	LastServiceDate    *time.Time `json:"last_service_date,omitempty" gorm:"comment:最近服务日期"`
	LastHairstyleDesc  string     `json:"last_hairstyle_desc" gorm:"type:varchar(200);comment:最近发型描述"`
	LastServiceItems   string     `json:"last_service_items" gorm:"type:varchar(200);comment:最近服务项目"`

	// 标签与备注
	Tags   string `json:"tags" gorm:"type:varchar(200);comment:自定义标签(VIP/常客/难搞等,逗号分隔)"`
	Remark string `json:"remark" gorm:"type:text;comment:备注"`

	// 统计数据（冗余，提升查询性能）
	TotalVisits   int     `json:"total_visits" gorm:"type:int;default:0;comment:总到访次数"`
	TotalSpent    float64 `json:"total_spent" gorm:"type:decimal(12,2);default:0.00;comment:累计消费金额"`
	FirstVisitAt  *time.Time `json:"first_visit_at,omitempty" gorm:"comment:首次到访时间"`

	// 关联关系
	Stylist *Stylist `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (CustomerProfile) TableName() string {
	return "customer_profiles"
}
