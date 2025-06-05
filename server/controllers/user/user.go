package user

import (
	userModel "TaskHub/models/user"
	userService "TaskHub/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var req userModel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	user, err := userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Register(c *gin.Context) {
	var req userModel.RegisterRequest
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	user, err := userService.Register(&req)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}