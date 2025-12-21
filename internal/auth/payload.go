package auth

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/pkg/db"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	UserRole int    `json:"userRole"`
}

type LoginResponse struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	UserRole int    `json:"userRole"`
}
type AuthHandler struct {
	*configs.Config
	AuthRepository AuthRepository
}

type AuthhandlerDeps struct {
	*configs.Config
	AuthRepository *AuthRepository
}
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type createRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserId   string `json:"userId"`
}
type updateRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	UserId   string `json:"userId"`
	UserRole int    `json:"userRole"`
}
type DeleteRequest struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}
type AuthRepository struct {
	DataBase *db.Db
}

type AuthRepositoryDeps struct {
	DataBase *db.Db
}
