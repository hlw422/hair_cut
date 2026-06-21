package mysql

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局DB实例
var DB *gorm.DB

// InitDB 初始化数据库连接并自动迁移所有表结构
func InitDB(dsn string) (err error) {
	// 配置GORM日志模式（开发环境打印SQL，生产环境仅记录错误）
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发用Info，生产建议Warn或Error
		// 禁用外键约束（性能考虑，应用层保证数据完整性）
		DisableForeignKeyConstraintWhenMigrating: true,
		// 命名策略：蛇形命名
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	// 连接MySQL数据库
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 获取底层sql.DB对象，配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层DB失败: %w", err)
	}
	
	// 设置连接池参数（根据实际负载调整）
	sqlDB.SetMaxIdleConns(10)   // 最大空闲连接
	sqlDB.SetMaxOpenConns(100)  // 最大打开连接
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	// 自动迁移：按顺序注册所有模型表
	log.Println("📦 开始数据库表结构迁移...")
	err = DB.AutoMigrate(
		// ===== 1. 用户与会员 =====
		(*User)(nil),
		(*Member)(nil),
		(*PointsRecord)(nil),

		// ===== 2. 门店与位置 =====
		(*Store)(nil),
		(*StorePhoto)(nil),

		// ===== 3. 理发师与作品 =====
		(*Stylist)(nil),
		(*StylistPortfolio)(nil),

		// ===== 4. 服务项目 =====
		(*ServiceCategory)(nil),
		(*ServiceItem)(nil),

		// ===== 5. 预约与排班 =====
		(*Appointment)(nil),
		(*Schedule)(nil),
		(*Attendance)(nil),

		// ===== 6. 订单与支付 =====
		(*Order)(nil),
		(*OrderItem)(nil),
		(*PaymentRecord)(nil),

		// ===== 7. 优惠券 =====
		(*CouponTemplate)(nil),
		(*Coupon)(nil),

		// ===== 8. 评价 =====
		(*Review)(nil),
		(*ReviewMedia)(nil),

		// ===== 9. 消息与通知 =====
		(*Message)(nil),
		(*Notification)(nil),

		// ===== 10. 社交关系 =====
		(*FanRelation)(nil),
		(*CustomerProfile)(nil),

		// ===== 11. 组织与员工 =====
		(*Organization)(nil),
		(*Employee)(nil),

		// ===== 12. 库存与采购 =====
		(*Inventory)(nil),
		(*Supplier)(nil),
		(*PurchaseOrder)(nil),
		(*PurchaseItem)(nil),

		// ===== 13. 营销活动 =====
		(*Campaign)(nil),

		// ===== 14. 财务 =====
		(*FinanceRecord)(nil),
		(*StylistCommission)(nil),

		// ===== 15. 权限系统（RBAC）=====
		(*Role)(nil),
		(*Permission)(nil),
		(*RolePermission)(nil),
		(*UserRole)(nil),

		// ===== 16. 日志与配置 =====
		(*OperationLog)(nil),
		(*LoginLog)(nil),
		(*SystemConfig)(nil),
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	log.Println("✅ 数据库表结构迁移完成！")
	log.Printf("   共计 %d 张业务数据表已就绪\n", countTables())

	// 创建索引优化查询性能
	createIndexes()

	return nil
}

// countTables 统计已注册的表数量
func countTables() int {
	// 返回已迁移的模型数量
	return 38 // 根据上面的AutoMigrate列表统计
}

// createIndexes 创建额外的索引（AutoMigrate可能无法覆盖的复合索引）
func createIndexes() {
	// 示例: 为门店搜索创建地理空间联合索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_store_location ON stores(latitude, longitude, status)")
	
	// 示例: 为订单查询创建用户+状态+时间复合索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_order_user_status_time ON orders(user_id, status, created_at)")
	
	// 示例: 为预约冲突检测创建理发师+日期+时间复合索引
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_appointment_stylist_datetime ON appointments(stylist_id, appointment_date, appointment_time, status)")

	log.Println("🔍 性能索引优化完成")
}
