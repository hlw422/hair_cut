package mysql

import (
	"time"
)

// Employee 员工表 - 门店员工（含理发师、助理、前台等）
type Employee struct {
	BaseModel
	// 归属
	StoreID  uint64 `json:"store_id" gorm:"not null;index;comment:所属门店ID"`
	UserID   *uint64 `json:"user_id,omitempty" gorm:"uniqueIndex;comment:关联系统用户ID(如已注册)"`
	OrgNodeID uint64 `json:"org_node_id" gorm:"index;comment:组织架构节点ID"`

	// 基本信息
	Name      string `json:"name" gorm:"type:varchar(50);not null;comment:姓名"`
	Phone     string `json:"phone" gorm:"type:varchar(20);index;comment:手机号"`
	IDNumber  string `json:"-" gorm:"type:varchar(18);uniqueIndex;comment:身份证号(加密存储)"`
	Gender    int8   `json:"gender" gorm:"type:tinyint;default:0;comment:性别"`
	Birthday *time.Time `json:"birthday,omitempty" gorm:"type:date;comment:出生日期"`
	AvatarURL string `json:"avatar_url" gorm:"type:text;comment:头像"`

	// 职位信息
	Position    string `json:"position" gorm:"type:varchar(50);comment:职位名称(理发师/助理/前台/店长)"`
	PositionType int8  `json:"position_type" gorm:"type:tinyint;index;comment:职位类型: 1理发师 2助理 3前台 4店长 5其他"`
	Title       string `json:"title" gorm:"type:varchar(50);comment:职称(首席/高级/资深)"`
	Level       int8   `json:"level" gorm:"type:tinyint;default:1;comment:级别(1-5)"`

	// 入离职信息
	HireDate      *time.Time `json:"hire_date,omitempty" gorm:"type:date;comment:入职日期"`
	TerminationDate *time.Time `json:"termination_date,omitempty" gorm:"type:date;comment:离职日期"`
	EmploymentStatus int8     `json:"employment_status" gorm:"type:tinyint;default:1;index;comment:状态: 0离职 1在职 2休假 3试用期"`

	// 薪资信息（脱敏存储，敏感字段需权限控制）
	SalaryType   int8    `json:"salary_type" gorm:"type:tinyint;default:1;comment:薪资类型: 1固定工资 2底薪+提成 3纯提成"`
	BaseSalary   float64 `json:"-" gorm:"type:decimal(10,2);default:0.00;comment:底薪(元)"`
	CommissionRate float64 `json:"commission_rate" gorm:"type:decimal(4,3);default:0.000;comment:提成比例(如0.300表示30%)"`

	// 银行卡信息（用于发薪，加密存储）
	BankName   string `json:"-" gorm:"type:varchar(50);comment:开户银行"`
	BankAccount string `json:"-" gorm:"type:varchar(30);comment:银行账号"`

	// 紧急联系人
	EmergencyContactName  string `json:"emergency_contact_name,omitempty" gorm:"type:varchar(50);comment:紧急联系人姓名"`
	EmergencyContactPhone string `json:"emergency_contact_phone,omitempty" gorm:"type:varchar(20);comment:紧急联系电话"`

	// 备注
	Remark string `json:"remark" gorm:"type:text;comment:备注"`

	// 关联关系
	Store        *Store          `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	User         *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Attendances  []Attendance    `json:"attendances,omitempty" gorm:"foreignKey:EmployeeID"` // 考勤记录
}

func (Employee) TableName() string {
	return "employees"
}

// Organization 组织架构节点表 - 树形结构（集团>区域>城市>门店）
type Organization struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(100);not null;comment:节点名称"`
	Type     int8   `json:"type" gorm:"type:tinyint;not null;comment:类型: 1集团总部 2区域公司 3城市分公司 4门店 5部门"`
	ParentID *uint64 `json:"parent_id,omitempty" gorm:"index;comment:父节点ID(NULL=根节点)"`
	Path     string `json:"path" gorm:"type:varchar(500);comment:路径(如/1/5/12/)"`
	Level    int    `json:"level" gorm:"type:int;default:1;comment:层级深度"`
	SortOrder int   `json:"sort_order" gorm:"type:int;default:0;comment:排序权重"`
	Status   int8   `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0禁用 1启用"`
	LeaderID *uint64 `json:"leader_id,omitempty" gorm:"comment:负责人ID(Employee ID)"`

	// 关联
	Parent  *Organization `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Organization `json:"children,omitempty" gorm:"foreignKey:ParentID"` // 子节点
	Stores  []Store        `json:"stores,omitempty" gorm:"foreignKey:OrgNodeID"`  // 下属门店
	Employees []Employee    `json:"employees,omitempty" gorm:"foreignKey:OrgNodeID"` // 下属员工
}

func (Organization) TableName() string {
	return "organizations"
}
