package repositories

import (
	"TaskHub/user_service/global"
	"TaskHub/user_service/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(req *models.RegisterRequest) (*models.User, error) {

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
		Role:     "user",
	}

	if err := global.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func Login(req *models.LoginRequest) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		if err := global.DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func DeleteUser(id uint) error {
	return global.DB.Delete(&models.User{}, id).Error
}

func GetUsers(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	offset := (page - 1) * pageSize

	if err := global.DB.Model(&models.User{}).Where("status = 1").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := global.DB.Where("status = 1").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("username =? AND status = 1", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("email =? AND status = 1", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
