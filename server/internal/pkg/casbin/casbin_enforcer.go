package casbin

import (
	"fmt"
	"sync"

	"haircut-server/pkg/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var (
	Enforcer *casbin.Enforcer
	once     sync.Once
)

// InitCasbin 初始化Casbin权限引擎（应用启动时调用一次）
// 参数：GORM DB实例（用于从数据库加载策略）
func InitCasbin(db *gorm.DB) error {
	var initErr error
	once.Do(func() {
		// 使用GORM适配器（策略存储在数据库的casbin_rule表中）
		adapter, err := gormadapter.NewAdapterByDB(db)
		if err != nil {
			initErr = fmt.Errorf("创建Casbin GORM适配器失败: %w", err)
			return
		}

		// 创建Enforcer实例（使用RBAC模型文件）
		enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
		if err != nil {
			initErr = fmt.Errorf("创建Casbin Enforcer失败: %w", err)
			return
		}

		// 加载策略（从数据库）
		if err := enforcer.LoadPolicy(); err != nil {
			initErr = fmt.Errorf("加载Casbin策略失败: %w", err)
			return
		}

		Enforcer = enforcer
		logger.Info("✅ Casbin权限引擎初始化成功")
	})

	return initErr
}

// AddPolicy 为用户/角色添加权限策略
func AddPolicy(role, path, method string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}
	
	_, err := Enforcer.AddPolicy(role, path, method)
	if err != nil {
		return fmt.Errorf("添加权限策略失败: %w", err)
	}
	
	// 保存到数据库（自动持久化）
	return Enforcer.SavePolicy()
}

// RemovePolicy 删除权限策略
func RemovePolicy(role, path, method string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}
	
	_, err := Enforcer.RemovePolicy(role, path, method)
	if err != nil {
		return fmt.Errorf("删除权限策略失败: %w", err)
	}
	
	return Enforcer.SavePolicy()
}

// HasPermission 检查是否有指定权限（便捷方法）
func HasPermission(role, path, method string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}
	return Enforcer.Enforce(role, path, method)
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(userID uint64) ([]string, error) {
	if Enforcer == nil {
		return nil, fmt.Errorf("Casbin未初始化")
	}
	userKey := fmt.Sprintf("user_%d", userID)
	return Enforcer.GetRolesForUser(userKey)
}
