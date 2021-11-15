package NotificationController

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	Message string `json:"message"`
}

func GetAdminNotificationController(c *gin.Context) {
	var notification string = c.Param("notification")
	c.JSON(http.StatusOK, Notification{notification})
}

func GetClientNotificationController(c *gin.Context) {
	var notification string = c.Param("notification")
	c.JSON(http.StatusOK, Notification{Message: notification})
}
