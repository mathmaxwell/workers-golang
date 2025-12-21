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
	router.HandleFunc("/users/createUser", handler.createUser())
	router.HandleFunc("/users/deleteUser", handler.deleteUser())
	router.HandleFunc("/users/updateUser", handler.updateUser())
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
func (handler *AuthHandler) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createRequest](&w, r)
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
		token := body.UserId
		isAdmin := body.Login == "testadmin" && body.Password == "tadi123$"
		userRole := 99
		if isAdmin {
			userRole = 1
		}
		data := User{
			Login:    email,
			Password: password,
			Token:    token,
			UserRole: userRole,
		}
		handler.AuthRepository.DataBase.Create(&data)
		res.Json(w, data, 200)
	}
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
		isAdmin := body.Login == "testadmin" && body.Password == "tadi123$"
		userRole := 99
		if isAdmin {
			userRole = 1
		}
		data := User{
			Login:    email,
			Password: password,
			Token:    token,
			UserRole: userRole,
		}
		handler.AuthRepository.DataBase.Create(&data)
		res.Json(w, data, 200)
	}
}
func (handler *AuthHandler) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DeleteRequest](&w, r)
		if err != nil {
			return
		}
		admin, err := handler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if admin.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		db := handler.AuthRepository.DataBase
		result := db.Delete(&User{}, "token = ?", body.UserId)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}
		res.Json(w, "user deleted", 200)
	}
}
func (handler *AuthHandler) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[updateRequest](&w, r)
		if err != nil {
			return
		}

		if body.UserId == "" {
			res.Json(w, "userId is required", 400)
			return
		}

		updates := map[string]interface{}{
			"login":     body.Login,
			"password":  body.Password,
			"user_role": body.UserRole,
		}

		if err := handler.AuthRepository.DataBase.
			Model(&User{}).
			Where("token = ?", body.UserId).
			Updates(updates).Error; err != nil {

			res.Json(w, err.Error(), 500)
			return
		}

		res.Json(w, "ok", 200)
	}
}
