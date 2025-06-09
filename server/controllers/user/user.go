package user

import (
	userModel "TaskHub/models/user"
	"TaskHub/pkg/auth"
	"TaskHub/pkg/logger"
	userService "TaskHub/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username,user.Role)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败: " + err.Error()})
		return
	}

	auth.SetSession(c, "token", token)

	logger.Info("用户登录成功", zap.String("username", user.Username))

	c.JSON(http.StatusOK, gin.H{})
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

	logger.Info("用户注册成功", zap.String("username", user.Username))

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	// 解析请求参数
	var req userModel.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	userID := c.GetUint("userID") // 从上下文中获取用户ID

	if _,err := userService.UpdateUser(userID, &req); err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logger.Info("用户信息更新成功", zap.Uint("userID", userID))

	c.JSON(http.StatusOK, gin.H{"msg": "用户信息更新成功"})
}

func Logout(c *gin.Context) {
	auth.DelSession(c, "token")
	logger.Info("用户退出成功")
	c.JSON(http.StatusOK, gin.H{"msg": "退出成功"})
}