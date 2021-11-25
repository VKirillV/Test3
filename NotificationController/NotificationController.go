package notificationcontroller

import (
	"io/ioutil"
	db "library/ConnectionDatabase"
	Error "library/JsonError"
	telegrambot "library/TelegramBot"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Notification struct {
	Message string `json:"message"`
}

func AdminNotificationController(ctx *gin.Context) {
	admin := "Admin"

	message, err := ioutil.ReadAll(ctx.Request.Body)
	if Error.Error(ctx, err) {
		log.Error("Failed to select certain data in the database! ", err)

		return
	}

	rows, err := db.Connect().Query("Select telegram_chat_id, username FROM user "+
		"WHERE user_type = (?) AND telegram_chat_id IS NOT NULL", admin)
	if err != nil {
		log.Error("Failed to select certain data in the database", err)
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var messageChatID int64

		var telegramUser string

		err := rows.Scan(&messageChatID, &telegramUser)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}

		processedNotification := telegrambot.EscapeMessage(string(message))

		go telegrambot.SendMessage(processedNotification, telegramUser, messageChatID)
	}

	processedNotification := telegrambot.EscapeMessage(string(message))

	ctx.JSON(http.StatusOK, Notification{processedNotification})
}

func ClientNotificationController(ctx *gin.Context) {
	client := "client"
	clientGUID := ctx.Param("client_guid")

	message, err := ioutil.ReadAll(ctx.Request.Body)
	if Error.Error(ctx, err) {
		log.Error("Failed to select certain data in the database! ", err)

		return
	}

	processedNotification := telegrambot.EscapeMessage(string(message))

	rows, err := db.Connect().Query("Select username, telegram_chat_id FROM user "+
		"LEFT JOIN client_user ON user.id = client_user.user_fk "+
		"WHERE client_user.client_guid = (?) AND user_type = (?)", clientGUID, client)
	if Error.Error(ctx, err) {
		log.Error("Failed to select certain data in the database! ", err)

		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		var messageChatID int64

		var telegramUser string

		err := rows.Scan(&telegramUser, &messageChatID)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}

		go telegrambot.SendMessage(processedNotification, telegramUser, messageChatID)
	}
	ctx.JSON(http.StatusOK, Notification{processedNotification})
}
