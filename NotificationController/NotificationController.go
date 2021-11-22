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
	c.JSON(http.StatusOK, Notification{notification})
	go TelegramBot.SendMessage(notification, admin)
}

func ClientNotificationController(c *gin.Context) {
	client := "Client"
	var notification string = c.Param("notification")
	c.JSON(http.StatusOK, Notification{notification})
	go TelegramBot.SendMessage(notification, client)
}
