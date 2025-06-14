package services

import (
	"TaskHub/img_service/utils/logger"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func UploadFile(src multipart.File, originalFileName string) (string, error) {
	// 创建img目录
	err := os.MkdirAll("./img", os.ModePerm)
	if err != nil {
		return "", err
	}

	// 提取文件扩展名
	ext := filepath.Ext(originalFileName)

	// 生成更短的文件名（UUID + 扩展名）
	safeName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 创建目标文件
	destPath := filepath.Join("./img", safeName)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(destPath) // 失败时清理文件
		return "", err
	}

	logger.Info("图片上传成功", zap.String("path", destPath))
	return destPath, nil
}