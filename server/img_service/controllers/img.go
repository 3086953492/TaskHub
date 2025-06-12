package controllers

import (
	"TaskHub/img_service/services"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// 验证文件类型是否为图片
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只允许上传JPEG, PNG或GIF图片"})
		return
	}

	// 打开文件并调用上传服务
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 转换为os.File以便上传
	dst, err := os.Create(file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调用上传服务
	filePath, err := services.UploadFile(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "图片上传成功",
		"path":    filePath,
	})
}
