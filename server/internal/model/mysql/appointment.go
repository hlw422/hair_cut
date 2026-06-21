package mysql

import (
	"encoding/json"
	"time"
)

// Appointment 预约记录表 - 用户在线预约核心表
type Appointment struct {
	BaseModel
	// 预约编号（对外展示）
	OrderNo string `json:"order_no" gorm:"type:varchar(64);uniqueIndex;not null;comment:预约单号"`

	// 关键实体关联
	UserID    uint64 `json:"user_id" gorm:"not null;index;comment:用户ID"`
	StoreID   uint64 `json:"store_id" gorm:"not null;index;comment:门店ID"`
	StylistID uint64 `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`

	// 预约时间
	AppointmentDate time.Time `json:"appointment_date" gorm:"type:date;not null;index;comment:预约日期"`
	AppointmentTime string    `json:"appointment_time" gorm:"type:varchar(20);not null;comment:预约时间段(如14:00-15:00)"`

	// 预约服务项目（JSON数组存储ServiceItem ID列表）
	ServiceIDs json.RawMessage `json:"service_ids" gorm:"type:json;comment:服务项目ID列表(JSON)"`

	// 金额信息
	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(10,2);default:0.00;comment:总金额(元)"`
	Deposit     float64 `json:"deposit" gorm:"type:decimal(10,2);default:0.00;comment:定金(元)(可选项)"`

	// 状态机（核心字段）
	Status int8 `json:"status" gorm:"type:tinyint;default:0;index;comment:预约状态: 0待确认 1已确认 2进行中 3已完成 4已取消 5用户爽约"`

	// 备注与来源
	Remark string `json:"remark" gorm:"type:varchar(500);comment:用户备注"`
	Source int8   `json:"source" gorm:"type:tinyint;default:1;comment:来源: 1小程序 2电话预约 3到店登记 4理发师代预约"`

	// 时间戳（业务时间）
	ConfirmedAt   *time.Time `json:"confirmed_at,omitempty" gorm:"comment:确认时间"`
	StartedAt     *time.Time `json:"started_at,omitempty" gorm:"comment:服务开始时间"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" gorm:"comment:完成时间"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty" gorm:"comment:取消时间"`
	CancelReason  string     `json:"cancel_reason,omitempty" gorm:"type:varchar(200);comment:取消原因"`

	// 关联关系
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Store        *Store         `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Stylist      *Stylist       `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
	Order        *Order         `json:"order,omitempty" gorm:"foreignKey:AppointmentID"` // 关联订单
}

func (Appointment) TableName() string {
	return "appointments"
}

// 预约状态常量
const (
	AppointmentStatusPending   = 0 // 待确认（刚创建）
	AppointmentStatusConfirmed = 1 // 已确认（理发师/系统确认）
	AppointmentStatusOngoing   = 2 // 进行中（开始服务）
	AppointmentStatusCompleted = 3 // 已完成（服务结束）
	AppointmentStatusCancelled = 4 // 已取消（用户或系统取消）
	AppointmentStatusNoShow    = 5 // 爽约（用户未到店）
)

// CanCancel 检查是否可以取消
func (a *Appointment) CanCancel() bool {
	return a.Status == AppointmentStatusPending || a.Status == AppointmentStatusConfirmed
}

// GetServiceIDList 获取服务ID列表（JSON解析）
func (a *Appointment) GetServiceIDList() []uint64 {
	var ids []uint64
	if a.ServiceIDs != nil {
		json.Unmarshal(a.ServiceIDs, &ids)
	}
	return ids
}

// Schedule 排班表 - 理发师工作排班
type Schedule struct {
	BaseModel
	StylistID    uint64     `json:"stylist_id" gorm:"not null;index;comment:理发师ID"`
	Date         time.Time  `json:"date" gorm:"type:date;not null;index;comment:日期"`
	ShiftType    int8       `json:"shift_type" gorm:"type:tinyint;default:1;comment:班次: 1全天 2早班 3晚班 4休息 5请假"`
	StartTime    string     `json:"start_time" gorm:"type:varchar(10);comment:上班时间(如09:00)"`
	EndTime      string     `json:"end_time" gorm:"type:varchar(10);comment:下班时间(如18:00)"`
	MaxAppointments int      `json:"max_appointments" gorm:"type:int;default:8;comment:最大可预约数"`
	Status       int8       `json:"status" gorm:"type:tinyint;default:1;comment:状态: 0异常 1正常"`
	Remark       string     `json:"remark" gorm:"type:varchar(200);comment:备注(请假原因等)"`

	// 关联
	Stylist *Stylist `json:"stylist,omitempty" gorm:"foreignKey:StylistID"`
}

func (Schedule) TableName() string {
	return "schedules"
}

// Attendance 考勤记录表 - 员工上下班打卡
type Attendance struct {
	BaseModel
	EmployeeID uint64     `json:"employee_id" gorm:"not null;index;comment:员工ID"`
	Date       time.Time  `json:"date" gorm:"type:date;not null;index;comment:考勤日期"`
	CheckInTime  *time.Time `json:"check_in_time" gorm:"comment:签到时间"`
	CheckOutTime *time.Time `json:"check_out_time" gorm:"comment:签退时间"`
	Type        int8       `json:"type" gorm:"type:tinyint;default:1;comment:类型: 1正常 2迟到 3早退 4缺卡 5请假 6加班"`
	WorkHours   float64    `json:"work_hours" gorm:"type:decimal(4,2);default:0.00;comment:工作时长(小时)"`
	Remark      string     `json:"remark" gorm:"type:varchar(200);comment:备注"`

	// 关联
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

func (Attendance) TableName() string {
	return "attendances"
}
