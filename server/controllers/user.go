package controllers

import (
	"TaskHub/global"
	"TaskHub/models"
	"TaskHub/pkg/auth"
	"TaskHub/pkg/logger"
	"TaskHub/services"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Login(c *gin.Context) {

	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	user, err := services.Login(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败: " + err.Error()})
		return
	}

	auth.SetSession(c, "token", token)

	logger.Info("用户登录成功", zap.String("username", user.Username))

	c.JSON(http.StatusOK, gin.H{})
}

func Register(c *gin.Context) {

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if err := global.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不符要求: " + err.Error()})
		return
	}

	user, err := services.Register(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logger.Info("用户注册成功", zap.String("username", user.Username))

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	// 解析请求参数
	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	userID := c.GetUint("userID") // 从上下文中获取用户ID

	user,err := services.GetUserByID(userID) // 从数据库中获取用户信息
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败: " + err.Error()})
		return
	}

	// 校验用户名与邮箱是否与数据库中一致，若一致则将请求体中的数据置空，防止唯一校验不通过
	if req.Username == user.Username {
		req.Username = ""
	}
	if req.Email == user.Email {
		req.Email = ""
	}

	if err := global.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不符要求: " + err.Error()})
		return
	}
	
	if _, err := services.UpdateUser(userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	logger.Info("用户信息更新成功", zap.Uint("userID", userID))

	c.JSON(http.StatusOK, gin.H{"msg": "用户信息更新成功"})
}

func Logout(c *gin.Context) {
	auth.DelSession(c, "token")
	logger.Info("用户退出成功", zap.String("username", c.GetString("username")))
	c.JSON(http.StatusOK, gin.H{"msg": "退出成功"})
}
