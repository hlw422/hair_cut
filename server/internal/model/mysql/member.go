package mysql

import (
	"time"
)

// Member 会员信息表 - 用户会员等级与权益
type Member struct {
	BaseModel
	UserID           uint64      `json:"user_id" gorm:"uniqueIndex;not null;comment:关联用户ID"`
	Level            int8        `json:"level" gorm:"type:tinyint;default:1;index;comment:会员等级: 1普通 2银卡 3金卡 4黑金 5钻石"`
	Points           int         `json:"points" gorm:"type:int;default:0;comment:当前积分余额"`
	Balance          float64     `json:"balance" gorm:"type:decimal(10,2);default:0.00;comment:余额(元)"`
	TotalSpent       float64     `json:"total_spent" gorm:"type:decimal(12,2);default:0.00;comment:累计消费金额"`
	OrderCount       int         `json:"order_count" gorm:"type:int;default:0;comment:总订单数"`
	ExpireDate       *time.Time `json:"expire_date" gorm:"type:date;comment:会员有效期至"`
	UpgradeThreshold float64    `json:"upgrade_threshold" gorm:"type:decimal(12,2);default:0.00;comment:升级下一等级所需消费额"`

	// 关联
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PointsRecords []PointsRecord `json:"points_records,omitempty" gorm:"foreignKey:MemberID"` // 积分记录

	// 等级名称（虚拟字段，不从DB读取）
	LevelName string `json:"level_name" gorm:"-"`
}

func (Member) TableName() string {
	return "members"
}

// GetLevelName 获取会员等级中文名称
func (m *Member) GetLevelName() string {
	levelNames := map[int8]string{
		1: "普通会员",
		2: "银卡会员",
		3: "金卡会员",
		4: "黑金会员",
		5: "钻石会员",
	}
	if name, ok := levelNames[m.Level]; ok {
		return name
	}
	return "未知"
}

// CanUpgrade 检查是否满足升级条件
func (m *Member) CanUpgrade(nextThreshold float64) bool {
	if m.Level >= 5 { // 钻石已是最高级
		return false
	}
	return m.TotalSpent >= nextThreshold
}

// PointsRecord 积分变动记录表
type PointsRecord struct {
	BaseModel
	MemberID uint64  `json:"member_id" gorm:"index;not null;comment:会员ID"`
	Type     int8    `json:"type" gorm:"type:tinyint;not null;comment:类型: 1获得 2消耗(兑换) 3过期 4管理员调整"`
	Points   int     `json:"points" gorm:"type:int;not null;comment:积分数(正数增加/负数减少)"`
	Balance  int     `json:"balance" gorm:"type:int;not null;comment:变动后余额"`
	Source   string  `json:"source" gorm:"type:varchar(50);comment:来源描述"`
	OrderID  *uint64 `json:"order_id,omitempty" gorm:"comment:关联订单ID"`
	Remark   string  `json:"remark" gorm:"type:varchar(200);comment:备注"`

	// 关联
	Member *Member `json:"member,omitempty" gorm:"foreignKey:MemberID"`
}

func (PointsRecord) TableName() string {
	return "points_records"
}
