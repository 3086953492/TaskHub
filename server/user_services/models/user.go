package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"unique;not null;size:50" json:"username"`
	Email     string         `gorm:"unique;not null;size:100" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Status    int            `gorm:"default:1" json:"status"`          // 1:正常 0:禁用
	Role      string         `gorm:"default:user;size:20" json:"role"` // admin, user
}

func (User) TableName() string {
	return "users"
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=15,usernameUnique"`
	Email    string `json:"email" validate:"required,email,emailUnique"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"required,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=15"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=15,usernameUnique"`
	Email    string `json:"email" validate:"omitempty,email,emailUnique"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Avatar   string `json:"avatar"`
}
