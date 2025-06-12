package services

import (
	"TaskHub/img_service/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func UploadFile(file *os.File) (string, error) {
	err := os.MkdirAll("./img", os.ModePerm)
	if err != nil {
		return "", err
	}

	// 生成安全的唯一文件名
	fileName := filepath.Base(file.Name())
	ext := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(ext)]
	safeName := fmt.Sprintf("%s_%d%s", baseName, time.Now().UnixNano(), ext)

	// 移动文件
	destPath := filepath.Join("./img", safeName)
	err = os.Rename(file.Name(), destPath)
	if err != nil {
		return "", err
	}
	logger.Info("图片上传成功", zap.String("path", destPath))
	return destPath, nil
}
