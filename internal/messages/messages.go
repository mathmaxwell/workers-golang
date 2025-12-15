package messages

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewMessageRepository(dataBase *db.Db) *MessageRepository {
	return &MessageRepository{
		DataBase: dataBase,
	}
}
func NewMessagesHandler(router *http.ServeMux, deps MessagehandlerDeps) {
	handler := &MessageHandler{
		Config:            deps.Config,
		MessageRepository: deps.MessageRepository,
		AuthHandler:       deps.AuthHandler,
	}
	router.HandleFunc("/messages/createMessage", handler.createMessage())
	router.HandleFunc("/messages/getUnreadMessages", handler.getUnreadMessages())
	router.HandleFunc("/messages/getLastMessages", handler.getLastMessages())
	router.HandleFunc("/messages/getMessageWithUser", handler.getMessageWithUser())
	router.HandleFunc("/messages/readedMessage", handler.readedMessage())

}

func (handler *MessageHandler) createMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createResponse](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		newMessage := Message{
			Id:           token.CreateId(),
			SenderId:     body.Token,
			RecipientIds: body.RecipientIds,
			Content:      body.Content,
			SentMinute:   body.SentMinute,
			SentHour:     body.SentHour,
			SentDay:      body.SentDay,
			SentMonth:    body.SentMonth,
			SentYear:     body.SentYear,
			IsRead:       false,
		}
		db := handler.MessageRepository.DataBase
		err = db.Create(&newMessage).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, newMessage, 200)
	}
}
func (handler *MessageHandler) readedMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[readedMessage](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		db := handler.MessageRepository.DataBase
		var oldMessage Message
		err = db.Where("id = ?", body.Id).First(&oldMessage).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		newMessage := Message{
			Id:           oldMessage.Id,
			SenderId:     oldMessage.SenderId,
			RecipientIds: oldMessage.RecipientIds,
			Content:      oldMessage.Content,
			SentMinute:   oldMessage.SentMinute,
			SentHour:     oldMessage.SentHour,
			SentDay:      oldMessage.SentDay,
			SentMonth:    oldMessage.SentMonth,
			SentYear:     oldMessage.SentYear,
			IsRead:       true,
		}
		err = db.Model(&newMessage).Updates(newMessage).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, newMessage, 200)
	}
}
func (handler *MessageHandler) getUnreadMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getUnreadMessagesRequest](&w, r)
		if err != nil {
			return
		}
		data := []Message{}
		db := handler.MessageRepository.DataBase
		err = db.Where("recipient_ids = ?", body.Token).Where("is_read = ?", false).Find(&data).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, data, 200)
	}
}
func (handler *MessageHandler) getLastMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getUnreadMessagesRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		token := body.Token

		var allMessages []Message
		err = handler.MessageRepository.DataBase.
			Where("sender_id = ? OR recipient_ids = ?", token, token).
			Find(&allMessages).Error

		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}

		lastMessages := make(map[string]Message)

		for _, msg := range allMessages {
			var other string
			if msg.SenderId == token {
				other = msg.RecipientIds
			} else {
				other = msg.SenderId
			}

			if existing, ok := lastMessages[other]; !ok || isAfter(msg, existing) {
				lastMessages[other] = msg
			}
		}
		data := make([]Message, 0)
		for _, msg := range lastMessages {
			data = append(data, msg)
		}

		res.Json(w, data, 200)
	}
}
func (handler *MessageHandler) getMessageWithUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getMessageWithUserRequest](&w, r)
		if err != nil {
			return
		}
		var data []Message
		db := handler.MessageRepository.DataBase

		err = db.
			Where(
				"(sender_id = ? AND recipient_ids = ?) OR (sender_id = ? AND recipient_ids = ?)",
				body.Token, body.EmployeeId,
				body.EmployeeId, body.Token,
			).
			Order("sent_year DESC, sent_month DESC, sent_day DESC, sent_hour DESC, sent_minute DESC").
			Find(&data).Error

		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		res.Json(w, data, 200)
	}
}
