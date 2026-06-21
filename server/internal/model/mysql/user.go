package mysql

import (
	"time"
)

// User 用户表 - 存储所有平台用户（会员、理发师、店长、管理员等）
type User struct {
	BaseModel
	// 微信认证信息
	OpenID     string `json:"openid" gorm:"type:varchar(100);uniqueIndex;comment:微信OpenID"`
	UnionID    string `json:"union_id" gorm:"type:varchar(100);index;comment:微信UnionID(跨应用唯一)"`
	SessionKey string `json:"-" gorm:"type:varchar(100);comment:微信SessionKey(加密存储)"`

	// 基本信息
	Phone       string `json:"phone" gorm:"type:varchar(20);uniqueIndex;comment:手机号"`
	Password    string `json:"-" gorm:"type:varchar(200);comment:密码(BCrypt加密)"`
	Nickname    string `json:"nickname" gorm:"type:varchar(50);not null;comment:昵称"`
	AvatarURL   string `json:"avatar_url" gorm:"type:text;comment:头像URL"`
	Gender      int8   `json:"gender" gorm:"type:tinyint;default:0;comment:性别: 0未知 1男 2女"`
	Birthday    *time.Time `json:"birthday" gorm:"type:date;comment:生日"`
	CityCode    string `json:"city_code" gorm:"type:varchar(10);index;comment:城市代码"`

	// 状态与权限
	Status int8 `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1正常 2待审核"`

	// 关联关系
	Member      *Member        `json:"member,omitempty" gorm:"foreignKey:UserID"` // 一对一会员信息
	Roles       []UserRole     `json:"roles,omitempty" gorm:"foreignKey:UserID"`  // 多角色关联
	Orders      []Order        `json:"orders,omitempty" gorm:"foreignKey:UserID"` // 订单列表
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:UserID"` // 预约记录
}

func (User) TableName() string {
	return "users"
}

// UserPublic 用户公开信息（用于API响应，隐藏敏感字段）
type UserPublic struct {
	ID        uint64    `json:"id"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Gender    int8      `json:"gender"`
}
