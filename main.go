package main

import (
	admincontroller "library/AdminController"
	clientcontroller "library/ClientController"
	notificationcontroller "library/NotificationController"
	telegrambot "library/TelegramBot"
	"os"

	//"unclass_console\internal\controllers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {
	telegrambot.Bot = telegrambot.InitBot()

	go telegrambot.Listen()

	log.Info("Server is starting...")

	port_server := ":" + os.Getenv("port_server")

	r := gin.Default()

	r.POST("/kript/:adminname/admin", admincontroller.PostAdminController)
	r.DELETE("/kript/:adminname/admin", admincontroller.DeleteAdminController)
	r.POST("/:username/client/:guid", clientcontroller.PostClientController)
	r.DELETE("/:username/client/:guid", clientcontroller.DeleteClientController)
	r.GET("/admin", admincontroller.GetAdminController)
	r.GET("/client", clientcontroller.GetClientController)
	r.POST("/admin", notificationcontroller.AdminNotificationController)
	r.POST("/client/:client_guid", notificationcontroller.ClientNotificationController)
	errServer := r.Run(port_server)
	if errServer != nil {
		log.Error(errServer)
	}

}
