package mysql

import (
	"time"
)

// OperationLog 操作日志表 - 记录用户关键操作（审计）
type OperationLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint64   `json:"user_id,omitempty" gorm:"index;comment:操作用户ID(NULL=系统)"`
	Username  string    `json:"username" gorm:"type:varchar(50);comment:操作用户名"`
	TenantID  uint64    `json:"tenant_id" gorm:"index;default:0;comment:租户ID"`

	// 操作信息
	Module    string `json:"module" gorm:"type:varchar(50);index;comment:模块(user/store/order等)"`
	Action    string `json:"action" gorm:"type:varchar(50);not null;comment:操作类型(create/update/delete/login等)"`
	TargetType string `json:"target_type" gorm:"type:varchar(50);comment:目标对象类型"`
	TargetID  *uint64 `json:"target_id,omitempty" gorm:"comment:目标对象ID"`

	// 详情
	Content     string `json:"content" gorm:"type:text;comment:操作内容描述"`
	IPAddress   string `json:"ip_address" gorm:"type:varchar(45);comment:IP地址(IPV6兼容)"`
	UserAgent   string `json:"user_agent" gorm:"type:varchar(500);comment:浏览器UA"`
	Method      string `json:"method" gorm:"type:varchar(10);comment:HTTP方法(GET/POST/PUT/DELETE)"`
	RequestURI  string `json:"request_uri" gorm:"type:varchar(500);comment:请求路径"`
	RequestData string `json:"-" grom:"type:text;comment:请求参数(JSON,敏感字段需脱敏)"`
	ResponseData string `json:"-" gorm:"type:text;comment:响应数据(JSON)"`

	// 状态与时间
	Status   int8       `json:"status" gorm:"type:tinyint;default:1;comment:结果: 0失败 1成功"`
	Duration int64      `json:"duration" gorm:"type:bigint;default:0;comment:耗时(毫秒)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:操作时间"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

// LoginLog 登录日志表 - 记录用户登录行为（安全审计）
type LoginLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint64   `json:"user_id,omitempty" gorm:"index;comment:用户ID"`
	Username  string    `json:"username" gorm:"type:varchar(50);comment:登录账号"`
	TenantID  uint64    `json:"tenant_id" gorm:"index;default:0;comment:租户ID"`

	// 登录信息
	LoginType  int8   `json:"login_type" gorm:"type:tinyint;default:1;comment:方式: 1密码登录 2微信登录 3短信验证码"`
	Platform   int8   `json:"platform" gorm:"type:tinyint;default:1;comment:平台: 1小程序 2Web后台 3App"`
	DeviceInfo string `json:"device_info" gorm:"type:varchar(200);comment:设备信息"`

	// IP与位置
	IPAddress  string `json:"ip_address" gorm:"type:varchar(45);index;comment:IP地址"`
	Location   string `json:"location" gorm:"type:varchar(100);comment:地理位置"`

	// 结果
	Status  int8       `json:"status" gorm:"type:tinyint;not null;comment:结果: 0失败 1成功"`
	Message string     `json:"message" gorm:"type:varchar(200);comment:提示信息(如:密码错误)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;index;comment:登录时间"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}

// SystemConfig 系统配置表 - 全局配置项（KV存储）
type SystemConfig struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `json:"key" gorm:"type:varchar(100);uniqueIndex;not null;comment:配置键"`
	Value     string    `json:"value" gorm:"type:text;comment:配置值"`
	ValueType int8      `json:"value_type" gorm:"type:tinyint;default:1;comment:值类型: 1字符串 2数字 3布尔 4JSON 5文本"`
	Group     string    `json:"group" gorm:"type:varchar(50);index;comment:配置分组(basic/payment/wechat/map等)"`
	Name      string    `json:"name" gorm:"type:varchar(100);comment:配置名称(中文)"`
	Description string  `json:"description" gorm:"type:varchar(200);comment:配置说明"`
	IsPublic  bool      `json:"is_public" gorm:"type:tinyint;default:false;comment:是否公开(前端可读)"`
	SortOrder int       `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`
	Status    int8      `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1启用"`
	CreatedAt time.Time `json:"created_at" gorm:"autoUpdateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
