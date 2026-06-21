package minio

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"haircut-server/internal/config"
	"haircut-server/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client MinIO对象存储客户端
var Client *minio.Client

// InitMinIO 初始化MinIO连接
func InitMinIO() error {
	var err error
	
	Client, err = minio.New(config.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIO.AccessKeyID, config.MinIO.SecretKeyID, ""),
		Secure:  config.MinIO.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("MinIO连接失败: %w", err)
	}

	// 确保存储桶存在（不存在则自动创建）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := Client.BucketExists(ctx, config.MinIO.BucketName)
	if err != nil {
		return fmt.Errorf("检查Bucket失败: %w", err)
	}

	if !exists {
		if err := Client.MakeBucket(ctx, config.MinIO.BucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("创建Bucket失败: %w", err)
		}
		logger.Info("✅ MinIO Bucket创建成功: %s", config.MinIO.BucketName)
	}

	logger.Info("✅ MinIO连接成功")
	return nil
}

// UploadFile 上传文件到MinIO
// 参数：文件流、文件名、内容类型、文件大小
// 返回：访问URL、错误信息
func UploadFile(ctx context.Context, reader io.Reader, objectName string, contentType string, size int64) (string, error) {
	// 构建对象路径：按日期分目录，避免单目录文件过多
	prefix := time.Now().Format("2006/01/02/")
	objectPath := fmt.Sprintf("%s%s", prefix, objectName)

	// 设置上传选项
	opts := minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-meta-uploaded-at": time.Now().Format(time.RFC3339),
		},
	}

	// 执行上传
	info, err := Client.PutObject(ctx, config.MinIO.BucketName, objectPath, reader, size, opts)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	// 生成访问URL（MinIO默认不公开，需通过Presigned URL或Nginx代理访问）
	url := GenerateURL(objectPath)

	logger.Debug("文件上传成功: %s (大小: %dB)", objectPath, info.Size)
	return url, nil
}

// GenerateURL 生成文件访问URL
func GenerateURL(objectPath) string {
	// 方式1: Presigned URL（临时有效，适合私有文件）
	// ctx := context.Background()
	// url, _ = Client.PresignedGetObject(ctx, bucketName, objectPath, 24*time.Hour, reqParams)
	
	// 方式2: 直接拼接（需要配置MinIO或Nginx公开访问策略）
	baseURL := strings.TrimRight(config.MinIO.Endpoint, "/")
	if !config.MinIO.UseSSL {
		baseURL = "http://" + baseURL
	} else {
		baseURL = "https://" + baseURL
	}
	
	return fmt.Sprintf("%s/%s/%s", baseURL, config.MinIO.BucketName, objectPath)
}

// GetPresignedURL 获取临时下载链接（用于私有文件）
func GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := Client.PresignedGetObject(ctx, config.MinIO.BucketName, objectPath, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("生成预签名链接失败: %w", err)
	}
	return presignedURL.String(), nil
}

// DeleteFile 删除文件
func DeleteFile(ctx context.Context, objectPath string) error {
	return Client.RemoveObject(ctx, config.MinIO.BucketName, objectPath, minio.RemoveObjectOptions{})
}

// ListFiles 列出指定前缀的文件（如某用户的所有头像）
func ListFiles(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	var objects []min64.ObjectInfo
	
	objectCh := Client.ListObjects(ctx, config.MinIO.BucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return objects, obj.Err
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// Close 关闭连接（MinIO HTTP client通常不需要显式关闭）
func Close() {
	Client = nil
}
