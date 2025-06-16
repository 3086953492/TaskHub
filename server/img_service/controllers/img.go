package controllers

import (
	"TaskHub/img_service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// 验证文件类型是否为图片
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "图片上传失败",
			"data": gin.H{
				"path": "",
			},
		})
		return
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "只允许上传JPEG, PNG或GIF图片",
			"data": gin.H{
				"path": "",
			},
		})
		return
	}

	// 确保上传的文件可以被读取
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "图片上传失败",
			"data": gin.H{
				"path": "",
			},
		})
		return
	}
	defer src.Close()

	// 调用上传服务，直接传递文件内容和文件名
	filePath, err := services.UploadFile(src, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "图片上传失败",
			"data": gin.H{
				"path": "",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "图片上传成功",
		"data": gin.H{
			"path": filePath,
		},
	})
}
