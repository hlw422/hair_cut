package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig 应用程序配置结构
type AppConfig struct {
	Env      string `mapstructure:"env"`       // 环境: development, staging, production
	Name     string `mapstructure:"name"`      // 应用名称
	Version  string `mapstructure:"version"`   // 版本号
	LogLevel string `mapstructure:"log_level"` // 日志级别
	LogPath  string `mapstructure:"log_path"`  // 日志路径
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`          // HTTP 端口
	Mode         string `mapstructure:"mode"`          // Gin 运行模式: debug, release, test
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时(秒)
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时(秒)
	JWTSecret    string `mapstructure:"jwt_secret"`    // JWT 密钥
	JWTExpire    int    `mapstructure:"jwt_expire"`    // JWT 过期时间(小时)
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"db_name"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"` // GORM 日志级别: silent, error, warn, info
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// ElasticsearchConfig ES 配置
type ElasticsearchConfig struct {
	URLs []string `mapstructure:"urls"` // ES 节点地址列表
}

// MongoDBConfig MongoDB 配置
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

// MinIOConfig MinIO 对象存储配置
type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_key_id"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
}

// WechatConfig 微信小程序配置
type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

// WechatPayConfig 微信支付配置
type WechatPayConfig struct {
	MchID         string `mapstructure:"mch_id"`
	APIKey        string `mapstructure:"api_key"`
	CertPath      string `mapstructure:"cert_path"`
	KeyPath       string `mapstructure:"key_path"`
	NotifyURL     string `mapstructure:"notify_url"`
	Sandbox       bool   `mapstructure:"sandbox"`
}

// TencentMapConfig 腾讯地图配置
type TencentMapConfig struct {
	SDKKey string `mapstructure:"sdk_key"`
}

// 全局配置变量（初始化后可访问）
var (
	App           AppConfig
	Server        ServerConfig
	MySQL         MySQLConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	MongoDB       MongoDBConfig
	RabbitMQ      RabbitMQConfig
	MinIO         MinIOConfig
	Wechat        WechatConfig
	WechatPay     WechatPayConfig
	TencentMap    TencentMapConfig
)

// Init 初始化配置（从 config.yaml 和环境变量加载）
func Init() error {
	v := viper.New()

	// 设置配置文件名和类型
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 添加配置文件搜索路径（优先级从高到低）
	v.AddConfigPath("./config")  // 项目根目录/config
	v.AddConfigPath("../config") // server目录上级
	v.AddConfigPath(".")         // 当前目录

	// 设置环境变量前缀，支持通过环境变量覆盖
	v.SetEnvPrefix("HAIRCUT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值
	setDefaults(v)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("⚠️  未找到配置文件，将使用默认值和环境变量")
		} else {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 反序列化到结构体
	return unmarshalConfig(v)
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 应用配置
	v.SetDefault("app.env", "development")
	v.SetDefault("app.name", "HairCut API Server")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.log_level", "info")
	v.SetDefault("app.log_path", "./logs/app.log")

	// 服务器配置
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 15)
	v.SetDefault("server.write_timeout", 15)
	v.SetDefault("server.jwt_secret", "dev-secret-key-change-in-production-2024")
	v.SetDefault("server.jwt_expire", 24)

	// MySQL 配置
	v.SetDefault("mysql.host", "localhost")
	v.SetDefault("mysql.port", 3306)
	v.SetDefault("mysql.user", "haircut_user")
	v.SetDefault("mysql.password", "haircut_dev_2024")
	v.SetDefault("mysql.db_name", "haircut")
	v.SetDefault("mysql.max_idle_conns", 10)
	v.SetDefault("mysql.max_open_conns", 100)
	v.SetDefault("mysql.log_level", "info")

	// Redis 配置
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 25)

	// Elasticsearch 配置
	v.SetDefault("elasticsearch.urls", []string{"http://localhost:9200"})

	// MongoDB 配置
	v.SetDefault("mongodb.uri", "mongodb://admin:mongo_admin_2024@localhost:27017")
	v.SetDefault("mongodb.database", "haircut_logs")

	// RabbitMQ 配置
	vSetDefault := v.SetDefault
	vSetDefault("rabbitmq.host", "localhost")
	vSetDefault("rabbitmq.port", 5672)
	vSetDefault("rabbitmq.user", "haircut")
	vSetDefault("rabbitmq.password", "rabbitmq_2024")
	vSetDefault("rabbitmq.vhost", "/haircut")

	// MinIO 配置
	vSetDefault("minio.endpoint", "localhost:9000")
	vSetDefault("minio.access_key_id", "minioadmin")
	vSetDefault("minio.secret_key_id", "minioadmin_2024")
	vSetDefault("minio.use_ssl", false)
	vSetDefault("minio.bucket_name", "haircut-uploads")

	// 微信配置
	vSetDefault("wechat.app_id", "")
	vSetDefault("wechat.app_secret", "")

	// 微信支付配置
	vSetDefault("wechat_pay.mch_id", "")
	vSetDefault("wechat_pay.api_key", "")
	vSetDefault("wechat_pay.cert_path", "./certs/wechat/apiclient_cert.pem")
	vSetDefault("wechat_pay.key_path", "./certs/wechat/apiclient_key.pem")
	vSetDefault("wechat_pay.notify_url", "https://api.yourdomain.com/api/v1/payment/wechat/notify")
	vSetDefault("wechat_pay.sandbox", false)

	// 腾讯地图配置
	vSetDefault("tencen_map.sdk_key", "")
}

// unmarshalConfig 将 viper 配置反序列化到全局变量
func unmarshalConfig(v *viper.Viper) error {
	if err := v.Unmarshal(&App, "app"); err != nil {
		return fmt.Errorf("解析应用配置失败: %w", err)
	}
	if err := v.Unmarshal(&Server, "server"); err != nil {
		return fmt.Errorf("解析服务器配置失败: %w", err)
	}
	if err := v.Unmarshal(&MySQL, "mysql"); err != nil {
		return fmt.Errorf("解析MySQL配置失败: %w", err)
	}
	if err := v.Unmarshal(&Redis, "redis"); err != nil {
		return fmt.Errorf("解析Redis配置失败: %w", err)
	}
	if err := v.Unmarshal(&Elasticsearch, "elasticsearch"); err != nil {
		return fmt.Errorf("解析Elasticsearch配置失败: %w", err)
	}
	if err := v.Unmarshal(&MongoDB, "mongodb"); err != nil {
		return fmt.Errorf("解析MongoDB配置失败: %w", err)
	}
	if err := v.Unmarshal(&RabbitMQ, "rabbitmq"); err != nil {
		return fmt.Errorf("解析RabbitMQ配置失败: %w", err)
	}
	if err := v.Unmarshal(&MinIO, "minio"); err != nil {
		return fmt.Errorf("解析MinIO配置失败: %w", err)
	}
	if err := v.Unmarshal(&Wechat, "wechat"); err != nil {
		return fmt.Errorf("解析微信配置失败: %w", err)
	}
	if err := v.Unmarshal(&WechatPay, "wechat_pay"); err != nil {
		return fmt.Errorf("解析微信支付配置失败: %w", err)
	}
	if err := v.Unmarshal(&TencentMap, "tencent_map"); err != nil {
		return fmt.Errorf("解析腾讯地图配置失败: %w", err)
	}

	return nil
}
