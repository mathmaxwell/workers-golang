package auth

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"errors"
	"net/http"
)

func NewUserRepository(dataBase *db.Db) *AuthRepository {
	return &AuthRepository{
		DataBase: dataBase,
	}
}

func NewAuthHandler(router *http.ServeMux, deps AuthhandlerDeps) *AuthHandler {
	handler := &AuthHandler{
		Config:         deps.Config,
		AuthRepository: *deps.AuthRepository,
	}

	router.HandleFunc("/users/login", handler.login())
	router.HandleFunc("/users/register", handler.register())

	return handler
}

func (handler *AuthHandler) GetUserByToken(token string) (User, error) {
	var user User
	err := handler.AuthRepository.DataBase.Where("token = ?", token).First(&user).Error
	if err != nil {
		return user, errors.New("user is not found")
	}
	return user, nil
}
func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("login = ?", body.Login).First(&user).Error

		if err != nil {
			res.Json(w, "user is not found", 400)
			return
		}
		if user.Password != body.Password {
			res.Json(w, "password is not correct", 400)
			return
		}
		res.Json(w, user, 200)
	}
}
func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		var user User
		err = handler.AuthRepository.DataBase.Where("login = ?", body.Login).First(&user).Error
		if err == nil {
			res.Json(w, "login is alredy exist", 400)
			return
		}
		email := body.Login
		password := body.Password
		token := token.CreateId()
		data := User{
			Login:    email,
			Password: password,
			Token:    token,
			UserRole: 1,
		}
		handler.AuthRepository.DataBase.Create(&data)
		res.Json(w, data, 200)
	}
}
