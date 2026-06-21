package jwt

import (
	"errors"
	"fmt"
	"time"

	"haircut-server/internal/config"
	
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT载荷（存储用户身份信息）
type Claims struct {
	UserID   uint64   `json:"user_id"`    // 用户ID
	Username string   `json:"username"`    // 用户名
	Roles    []string `json:"roles"`       // 角色列表
	TenantID uint64   `json:"tenant_id"`   // 租户ID（多租户）
	jwt.RegisteredClaims                 // 标准声明（过期时间、签发时间等）
}

// GenerateToken 生成JWT Token
// 参数：用户ID、用户名、角色列表、租户ID
// 返回：Token字符串、过期时间、错误
func GenerateToken(userID uint64, username string, roles []string, tenantID uint64) (string, time.Time, error) {
	// Token有效期（从配置读取，默认24小时)
	expireTime := time.Now().Add(time.Duration(config.Server.JWTExpire) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间
			Issuer:    "haircut-server",               // 签发者
			Subject:   fmt.Sprintf("%d", userID),     // 主题（用户ID）
		},
	}

	// 使用HS256算法签名（生产环境建议RS256）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Server.JWTSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成Token失败: %w", err)
	}

	return tokenString, expireTime, nil
}

// ParseToken 解析并验证JWT Token
// 返回：Claims对象、错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token（验证签名和过期时间）
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		
		// 返回签名密钥（从配置读取）
		return []byte(config.Server.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Token解析失败: %w", err)
	}

	// 类型断言获取Claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的Token")
	}

	return claims, nil
}

// RefreshToken 刷新Token（延长有效期）
// 用于"记住我"或长时间操作场景
func RefreshToken(oldToken string) (string, time.Time, error) {
	// 1. 解析旧Token（允许即将过期的Token刷新）
	claims, err := ParseToken(oldToken)
	if err != nil {
		return "", time.Time{}, err
	}

	// 2. 检查旧Token是否在可刷新窗口内（如剩余1小时有效期内才允许刷新）
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry > time.Hour {
		return "", time.Time{}, errors.New("Token未到刷新时间")
	}

	// 3. 使用原用户信息生成新Token
	return GenerateToken(claims.UserID, claims.Username, claims.Roles, claims.TenantID)
}
