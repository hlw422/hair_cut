package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"haircut-server/internal/config"
	"haircut-server/internal/middleware"
	_ "haircut-server/internal/model/mysql" // 自动注册GORM模型
	"haircut-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 2. 初始化日志
	logger.Init(config.App.LogLevel, config.App.LogPath)

	// 3. 创建 Gin 引擎（根据环境切换模式）
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 4. 注册全局中间件（按顺序执行）
	r.Use(middleware.Recovery())      // 异常恢复
	r.Use(middleware.RequestID())     // 请求追踪ID
	r.Use(middleware.Logger())        // 日志记录
	rUse := middleware.CORS()         // 跨域处理
	r.Use(rUse)
	r.Use(middleware.RateLimit())     // 限流

	// 5. 健康检查路由（无需认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   config.App.Version,
			"service":   "HairCut API Server",
		})
	})

	// 6. API 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		public := apiV1.Group("/public")
		{
			public.GET("/stores/nearby", nil) // TODO: 附近门店查询
			public.GET("/stores/:id", nil)   // TODO: 门店详情
			public.POST("/auth/login", nil)  // TODO: 登录
			public.POST("/auth/wechat", nil) // TODO: 微信登录
		}

		// 需要认证的路由
		auth := apiV1.Group("")
		auth.Use(middleware.JWTAuth()) // JWT 认证中间件
		auth.Use(middleware.RBAC())    // Casbin 权限校验

		{
			// 用户模块
			userGroup := auth.Group("/user")
			{
				userGroup.GET("/profile", nil)
				userGroup.PUT("/profile", nil)
				userGroup.GET("/orders", nil)
			}

			// 门店模块
			storeGroup := auth.Group("/stores")
			{
				storeGroup.GET("", nil)
				storeGroup.POST("", nil)
				storeGroup.PUT("/:id", nil)
				storeGroup.DELETE("/:id", nil)
			}

			// 预约模块
			appointmentGroup := auth.Group("/appointments")
			{
				appointmentGroup.GET("", nil)
				appointmentGroup.POST("", nil)
				appointmentGroup.PUT("/:id", nil)
				appointmentGroup.POST("/:id/cancel", nil)
			}

			// 订单模块
			orderGroup := auth.Group("/orders")
			{
				orderGroup.GET("", nil)
				orderGroup.GET("/:id", nil)
				orderGroup.POST("/:id/pay", nil)
			}

			// 理发师端 API
			stylistGroup := auth.Group("/stylist")
			{
				stylistGroup.GET("/dashboard", nil)
				stylistGroup.GET("/appointments/today", nil)
				stylistGroup.GET("/customers", nil)
				stylistGroup.POST("/portfolio", nil)
			}

			// 店长端 API
			managerGroup := auth.Group("/manager")
			{
				managerGroup.GET("/dashboard/stats", nil)
				managerGroup.GET("/employees", nil)
				managerGroup.GET("/finance/reports", nil)
				managerGroup.GET("/inventory", nil)
			}

			// 总部后台管理 API
			adminGroup := auth.Group("/admin")
			{
				// 用户管理
				adminGroup.GET("/users", nil)
				adminGroup.PUT("/users/:id/status", nil)
				// 门店管理
				adminGroup.GET("/stores/all", nil)
				adminGroup.POST("/stores/approve", nil)
				// 数据分析
				adminGroup.GET("/analytics/dashboard", nil)
				adminGroup.GET("/analytics/gmv", nil)
				// 系统设置
				adminGroup.Get("/system/config", nil)
				adminGroup.Put("/system/config", nil)
			}
		}
	}

	// 7. 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", config.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在 goroutine 中启动服务以支持优雅关闭
	go func() {
		logger.Info(fmt.Sprintf("HairCut API Server 正在启动... %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("服务器启动失败: %v", err))
		}
	}()

	// 8. 优雅关闭处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("服务器强制关闭: %v", err))
	} else {
		logger.Info("服务器已优雅关闭")
	}
}
