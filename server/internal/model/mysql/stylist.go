package mysql

import (
	"encoding/json"
)

// Stylist 理发师表 - 门店理发师/造型师信息
type Stylist struct {
	BaseModel
	// 归属
	StoreID uint64 `json:"store_id" gorm:"not null;index;comment:所属门店ID"`
	UserID  *uint64 `json:"user_id,omitempty" gorm:"uniqueIndex;comment:关联用户ID(如已注册)"`

	// 基本信息
	Name           string `json:"name" gorm:"type:varchar(50);not null;comment:姓名"`
	AvatarURL      string `json:"avatar_url" gorm:"type:text;comment:头像URL"`
	Gender         int8   `json:"gender" gorm:"type:tinyint;default:0;comment:性别: 0未知 1男 2女"`
	Phone          string `json:"phone" gorm:"type:varchar(20);index;comment:联系电话"`

	// 专业信息
	Title           string `json:"title" gorm:"type:varchar(50);comment:职称(首席/高级/资深/助理等)"`
	ExperienceYears int    `json:"experience_years" gorm:"type:int;default:0;comment:从业年限"`
	Specialties     string `json:"specialties" gorm:"type:json;comment:擅长风格标签(JSON数组)"`
	Introduction    string `json:"introduction" gorm:"type:text;comment:个人简介"`

	// 数据统计（冗余存储，提升查询性能）
	PortfolioCount   int     `json:"portfolio_count" gorm:"type:int;default:0;comment:作品数量"`
	FanCount         int     `json:"fan_count" gorm:"type:int;default:0;comment:粉丝数"`
	AppointmentCount int     `json:"appointment_count" gorm:"type:int;default:0;comment:累计预约次数"`
	Rating           float32 `json:"rating" gorm:"type:decimal(2,1);default:5.0;comment:综合评分(1-5)"`
	ReviewCount      int     `json:"review_count" gorm:"type:int;default:0;comment:评价数量"`
	Level            int8    `json:"level" gorm:"type:tinyint;default:1;index;comment:等级(1-5)"`

	// 状态
	Status int8 `json:"status" gorm:"type:tinyint;default:1;index;comment:状态: 0离职 1在职 2休假"`

	// 关联关系
	Store        *Store             `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	User         *User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Portfolios   []StylistPortfolio `json:"portfolios,omitempty" gorm:"foreignKey:StylistID"` // 作品集
	Appointments []Appointment      `json:"appointments,omitempty" gorm:"foreignKey:StylistID"` // 预约记录
	Schedules    []Schedule         `json:"schedules,omitempty" gorm:"foreignKey:StylistID"`   // 排班表
	FanFromUsers []FanRelation      `json:"-" gorm:"foreignKey:StylistID"`                    // 粉丝(不序列化)
}

func (Stylist) TableName() string {
	return "stylists"
}

// GetSpecialtiesList 获取擅长风格标签列表
func (s *Stylist) GetSpecialtiesList() []string {
	var tags []string
	if s.Specialties != "" {
		json.Unmarshal([]byte(s.Specialties), &tags)
	}
	return tags
}

// StylistPortfolio 理发师作品集表 - 展示案例与效果
type StylistPortfolio struct {
	BaseModel
	StylistID uint64 `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`
	Type      int8   `json:"type" gorm:"type:tinyint;default:0;comment:类型: 0图片 1视频"`
	MediaURL  string `json:"media_url" gorm:"type:text;not null;comment:媒体文件URL(MinIO)"`
	CoverURL  string `json:"cover_url" gorm:"type:text;comment:视频封面图URL"`
	Title     string `json:"title" gorm:"type:varchar(100);comment:作品标题"`
	Desc      string `json:"desc" gorm:"type:text;comment:作品描述"`
	Tags      string `json:"tags" gorm:"type:varchar(200);comment:标签(逗号分隔)"`
	Likes     int    `json:"likes" gorm:"type:int;default:0;comment:点赞数"`
	ViewCount int    `json:"view_count" gorm:"type:int;default:0;comment:浏览量"`
	Status    int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0隐藏 1展示"`

	// 关联
	Stylist *Stylist `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
}

func (StylistPortfolio) TableName() string {
	return "stylist_portfolios"
}
