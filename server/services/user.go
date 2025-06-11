package services

import (
	"TaskHub/global"
	"TaskHub/models"
	"TaskHub/pkg/auth"
	"TaskHub/pkg/logger"
	"TaskHub/repositories"

	"go.uber.org/zap"
)

func LoginService(req *models.LoginRequest) (string, error) {

	user, err := repositories.Login(req)
	if err != nil {
		return "", err
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	logger.Info("用户登录成功", zap.String("username", user.Username))

	return token, nil
}

func RegisterService(req *models.RegisterRequest) (*models.User, error) {

	if err := global.Validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := repositories.Register(req)
	if err != nil {
		return nil, err
	}

	logger.Info("用户注册成功", zap.String("username", user.Username))

	return user, nil
}

func UpdateService(req *models.UpdateUserRequest, userID uint) error {

	user, err := repositories.GetUserByID(userID) // 从数据库中获取用户信息
	if err != nil {
		return err
	}

	// 校验用户名与邮箱是否与数据库中一致，若一致则将请求体中的数据置空，防止唯一校验不通过
	if req.Username == user.Username {
		req.Username = ""
	}
	if req.Email == user.Email {
		req.Email = ""
	}

	if err := global.Validate.Struct(req); err != nil {
		return err
	}

	if _, err := repositories.UpdateUser(userID, req); err != nil {
		return err
	}

	logger.Info("用户信息更新成功", zap.Uint("userID", userID))

	return nil
}
