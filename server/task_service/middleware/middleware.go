package middleware

import (
	"TaskHub/task_service/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "请先登录"})
			return
		}

		claims, err := services.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效的token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "需要管理员权限"})
			return
		}
		c.Next()
	}
}