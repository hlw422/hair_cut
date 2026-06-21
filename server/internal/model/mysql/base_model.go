package mysql

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 所有业务模型的公共基类
// 包含通用字段：主键、租户ID、创建/更新时间、软删除
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	TenantID  uint64         `json:"tenant_id,omitempty" gorm:"index;not null;default:0;comment:租户ID(多租户隔离)"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`
}

// BeforeCreate GORM Hook - 创建前自动填充
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.TenantID == 0 {
		// 从上下文获取当前租户ID（中间件设置）
		if tenantID, exists := tx.Get("tenant_id"); exists {
			m.TenantID = tenantID.(uint64)
		}
	}
	return nil
}
