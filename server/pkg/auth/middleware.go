package auth

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "请先登录"})
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效的token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
