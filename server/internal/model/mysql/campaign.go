package mysql

import (
	"encoding/json"
	"time"
)

// Campaign 营销活动表 - 优惠券/满减/拼团/秒杀等活动配置
type Campaign struct {
	BaseModel
	// 基本信息
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:活动名称"`
	Type        int8   `json:"type" gorm:"type:tinyint;not null;index;comment:类型: 1优惠券发放 2满减活动 3拼团 4秒杀 5会员日 6新人礼包"`
	Description string `json:"description" gorm:"type:text;comment:活动说明"`
	ImageURL    string `json:"image_url" gorm:"type:text;comment:活动封面图"`

	// 活动规则（JSON存储灵活规则）
	RuleConfig json.RawMessage `json:"rule_config" gorm:"type:json;not null;comment:活动规则配置(JSON)"`

	// 适用范围
	Scope       int8            `json:"scope" gorm:"type:tinyint;default:1;comment:适用范围: 1全国 2指定门店 3指定城市"`
	ScopeIDs    json.RawMessage `json:"scope_ids" gorm:"type:json;comment:适用范围ID列表(JSON数组)"`

	// 时间窗口
	StartTime time.Time  `json:"start_time" gorm:"not null;index;comment:活动开始时间"`
	EndTime   time.Time  `json:"end_time" gorm:"not null;index;comment:活动结束时间"`

	// 预算与限制
	Budget       float64 `json:"budget" gorm:"type:decimal(12,2);default:0.00;comment:活动预算(元)(0=不限)"`
	UsedBudget   float64 `json:"used_budget" gorm:"type:decimal(12,2);default:0.00;comment:已使用预算"`
	MaxParticipants int   `json:"max_participants" gorm:"type:int;default:0;comment:最大参与人数(0=不限)"`
	CurrentCount int     `json:"current_count" gorm:"type:int;default:0;comment:当前参与人数"`

	// 关联优惠券模板（可选）
	CouponTemplateID *uint64 `json:"coupon_template_id,omitempty" gorm:"comment:关联优惠券模板ID"`

	// 状态与统计
	Status      int8       `json:"status" gorm:"type:tinyint;default:1;index;comment:状态: 0草稿 1进行中 2已结束 3已暂停 4已取消"`
	TotalOrders int        `json:"total_orders" gorm:"type:int;default:0;comment:关联订单数"`
	TotalGMV    float64    `json:"total_gmv" gorm:"type:decimal(12,2);default:0.00;comment:活动总GMV"`

	// 关联关系
	CouponTemplate *CouponTemplate `json:"coupon_template,omitempty" gorm:"foreignKey:CouponTemplateID"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

const (
	CampaignTypeCoupon     = 1 // 优惠券发放
	CampaignTypeDiscount   = 2 // 满减活动
	CampaignTypeGroupBuy   = 3 // 拼团
	CampaignTypeFlashSale  = 4 // 秒杀
	CampaignTypeMemberDay  = 5 // 会员日
	CampaignTypeNewUserGift = 6 // 新人礼包

	CampaignStatusDraft     = 0 // 草稿
	CampaignStatusActive    = 1 // 进行中
	CampaignStatusEnded     = 2 // 已结束
	CampaignStatusPaused    = 3 // 已暂停
	CampaignStatusCancelled = 4 // 已取消
)
