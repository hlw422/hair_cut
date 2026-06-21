package mysql

import (
	"encoding/json"
	"time"
)

// CouponTemplate 优惠券模板表 - 定义优惠券规则
type CouponTemplate struct {
	BaseModel
	// 基本信息
	Name        string `json:"name" gorm:"type:varchar(100);not null;comment:优惠券名称"`
	Description string `json:"description" gorm:"type:text;comment:使用说明"`
	ImageURL    string `json:"image_url" gorm:"type:text;comment:优惠券展示图"`

	// 类型与面值
	Type   int8    `json:"type" gorm:"type:tinyint;not null;index;comment:类型: 1满减券 2折扣券 3新人券 4节日券 5门店券 6全国通用券"`
	Value  float64 `json:"value" gorm:"type:decimal(10,2);not null;comment:面值(满减=金额/折扣=百分比如8.5表示8.5折)"`

	// 使用条件
	MinSpend      float64 `json:"min_spend" gorm:"type:decimal(10,2);default:0.00;comment:最低消费金额(0表示无门槛)"`
	MaxDiscount   float64 `json:"max_discount" gorm:"type:decimal(10,2);default:0.00;comment:最大优惠金额(折扣券使用，0表示不限制)"`
	TargetGender  int8    `json:"target_gender" gorm:"type:tinyint;default:0;comment:适用性别: 0通用 1男 2女"`

	// 适用范围（多租户/门店隔离）
	Scope       int8            `json:"scope" gorm:"type:tinyint;default:1;comment:适用范围: 1全国通用 2指定城市 3指定门店 4指定服务"`
	ScopeConfig json.RawMessage `json:"scope_config" gorm:"type:json;comment:范围配置(JSON: 城市ID列表/门店ID列表等)"`

	// 发放控制
	TotalCount   int `json:"total_count" gorm:"type:int;default:-1;comment:发放总量(-1表示不限)"`
	IssuedCount  int `json:"issued_count" gorm:"type:int;default:0;comment:已发放数量"`
	PerUserLimit int `json:"per_user_limit" gorm:"type:int;default:1;comment:每人限领数量"`

	// 有效期设置
	ValidDays     int        `json:"valid_days" gorm:"type:int;default:30;comment:领取后有效天数"`
	FixedStartAt  *time.Time `json:"fixed_start_at,omitempty" gorm:"comment:固定开始时间(优先于ValidDays)"`
	FixedEndAt    *time.Time `json:"fixed_end_at,omitempty" gorm:"comment:固定结束时间"`

	// 状态与时间窗口
	Status    int8       `json:"status" gorm:"type:tinyint;default:1;index;comment:状态: 0停用 1启用 2已发完"`
	StartTime *time.Time `json:"start_time" gorm:"comment:活动开始时间"`
	EndTime   *time.Time `json:"end_time" gorm:"comment:活动结束时间"`
}

func (CouponTemplate) TableName() string {
	return "coupon_templates"
}

// Coupon 用户优惠券表 - 用户领取的优惠券实例
type Coupon struct {
	BaseModel
	TemplateID uint64     `json:"template_id" gorm:"not null;index;comment:模板ID"`
	UserID     uint64     `json:"user_id" gorm:"not null;index;comment:用户ID"`
	Code       string     `json:"code" gorm:"type:varchar(32);uniqueIndex;not null;comment:优惠券唯一码"`
	Status     int8       `json:"status" gorm:"type:tinyint;default:0;index;comment:状态: 0未使用 1已使用 2已过期"`
	ObtainedAt time.Time  `json:"obtained_at" gorm:"not null;comment:领取时间"`
	UsedAt     *time.Time `json:"used_at,omitempty" gorm:"comment:使用时间"`
	OrderID    *uint64    `json:"order_id,omitempty" gorm:"comment:使用订单ID"`
	ExpireAt   time.Time  `json:"expire_at" gorm:"not null;index;comment:过期时间"`
	Source     string     `json:"source" gorm:"type:varchar(50);comment:领取来源(活动/注册/签到等)"`

	// 关联关系
	Template *CouponTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	User     *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Coupon) TableName() string {
	return "coupons"
}

// IsValid 检查优惠券是否可用
func (c *Coupon) IsValid(minSpend float64, storeID uint64) bool {
	if c.Status != 0 {
		return false
	}
	if time.Now().After(c.ExpireAt) {
		return false
	}
	// TODO: 检查使用条件（最低消费、适用范围等）
	return true
}

const (
	CouponTypeFixed      = 1 // 满减券
	CouponTypeDiscount   = 2 // 折扣券
	CouponTypeNewUser    = 3 // 新人券
	CouponTypeHoliday    = 4 // 节日券
	CouponTypeStore      = 5 // 门店券
	CouponTypeNational   = 6 // 全国通用券

	CouponStatusUnused   = 0 // 未使用
	CouponStatusUsed     = 1 // 已使用
	CouponStatusExpired  = 2 // 已过期
)
