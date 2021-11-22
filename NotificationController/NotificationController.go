package NotificationController

import (
	"library/TelegramBot"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	Message string `json:"message"`
}

func AdminNotificationController(c *gin.Context) {
	admin := "Admin"
	var notification string = c.Param("notification")
	processedNotification := TelegramBot.EscapeMessage(notification)
	c.JSON(http.StatusOK, Notification{processedNotification})
	go TelegramBot.SendMessage(processedNotification, admin)
}

func ClientNotificationController(c *gin.Context) {
	client := "Client"
	var notification string = c.Param("notification")
	processedNotification := TelegramBot.EscapeMessage(notification)
	c.JSON(http.StatusOK, Notification{processedNotification})
	go TelegramBot.SendMessage(processedNotification, client)
}
