package user

import (
	"TaskHub/configs"
	"TaskHub/models/user"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(req *user.RegisterRequest) (*user.User, error) {

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   1,
		Role:     "user",
	}

	if err := configs.UserDB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func Login(req *user.LoginRequest) (*user.User, error) {
	var user user.User
	if err := configs.UserDB.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
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

func GetUserByID(id uint) (*user.User, error) {
	var user user.User
	if err := configs.UserDB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id uint, req *user.UpdateUserRequest) (*user.User, error) {
	var user user.User
	if err := configs.UserDB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
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
		if err := configs.UserDB.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func DeleteUser(id uint) error {
	return configs.UserDB.Delete(&user.User{}, id).Error
}

func GetUsers(page, pageSize int) ([]*user.User, int64, error) {
	var users []*user.User
	var total int64

	offset := (page - 1) * pageSize

	if err := configs.UserDB.Model(&user.User{}).Where("status = 1").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := configs.UserDB.Where("status = 1").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUserByUsername(username string) (*user.User, error) {
	var user user.User
	if err := configs.UserDB.Where("username =? AND status = 1", username).First(&user).Error; err!= nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := configs.UserDB.Where("email =? AND status = 1", email).First(&user).Error; err!= nil {
		return nil, err
	}
	return &user, nil
}