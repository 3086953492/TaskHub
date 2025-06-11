package controllers

import (
	"TaskHub/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RefreshHandler(c *gin.Context) {
	// 从请求头中获取旧令牌
	token := c.GetHeader("Authorization")

	// 调用刷新令牌的方法。此处不用处理错误：如果令牌还未即将过期，会返回旧令牌，前端仍可以继续使用；如果令牌无效，则返回的是空值，前端拿到空值，则重定向登录。
	newToken, _ := services.RefreshToken(token)

	// 返回新令牌
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
