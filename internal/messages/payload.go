package messages

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/pkg/db"
)

type LoginResponse struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Message struct {
	Id           string `json:"id"`
	SenderId     string `json:"senderId"`
	RecipientIds string `json:"recipientIds"`
	Content      string `json:"content"`
	SentMinute   int    `json:"sentMinute"`
	SentHour     int    `json:"sentHour"`
	SentDay      int    `json:"sentDay"`
	SentMonth    int    `json:"sentMonth"`
	SentYear     int    `json:"sentYear"`
	IsRead       bool   `json:"isRead"`
}
type createResponse struct {
	Token        string `json:"token"`
	RecipientIds string `json:"recipientIds"`
	Content      string `json:"content"`
	SentMinute   int    `json:"sentMinute"`
	SentHour     int    `json:"sentHour"`
	SentDay      int    `json:"sentDay"`
	SentMonth    int    `json:"sentMonth"`
	SentYear     int    `json:"sentYear"`
}
type MessageHandler struct {
	*configs.Config
	MessageRepository *MessageRepository
	AuthHandler       *auth.AuthHandler
}
type MessageRepository struct {
	DataBase *db.Db
}
type MessagehandlerDeps struct {
	*configs.Config
	MessageRepository *MessageRepository
	AuthHandler       *auth.AuthHandler
}
type getUnreadMessagesRequest struct {
	Token string `json:"token" validate:"required"`
}
type getMessageWithUserRequest struct {
	Token      string `json:"token" validate:"required"`
	EmployeeId string `json:"employeeId" validate:"required"`
}
type readedMessage struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

func isAfter(a, b Message) bool {
	if a.SentYear != b.SentYear {
		return a.SentYear > b.SentYear
	}
	if a.SentMonth != b.SentMonth {
		return a.SentMonth > b.SentMonth
	}
	if a.SentDay != b.SentDay {
		return a.SentDay > b.SentDay
	}
	if a.SentHour != b.SentHour {
		return a.SentHour > b.SentHour
	}
	if a.SentMinute != b.SentMinute {
		return a.SentMinute > b.SentMinute
	}
	return false
}
