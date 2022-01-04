package main

import (
	admincontroller "library/AdminController"
	clientcontroller "library/ClientController"
	notificationcontroller "library/NotificationController"
	telegrambot "library/TelegramBot"
	"os"
	"time"

	//"unclass_console\internal\controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {
	telegrambot.MainBot = telegrambot.InitMainBot()
	telegrambot.HealthBot = telegrambot.InitHealthBot()
	go telegrambot.Listen()
	log.SetReportCaller(true)
	log.Info("Server is starting...")

	port_server := ":" + os.Getenv("port_server")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/user/:name/admin", admincontroller.PostAdminController)
	r.DELETE("/user/:name/admin", admincontroller.DeleteAdminController)
	r.POST("/user/:name/client/:guid", clientcontroller.PostClientController)
	r.DELETE("/user/:name/client/:guid", clientcontroller.DeleteClientController)
	r.GET("/user/admin", admincontroller.GetAdminController)
	r.GET("/user/client", clientcontroller.GetClientController)
	r.POST("/notification/admin", notificationcontroller.AdminNotificationController)
	r.POST("/notification/client/:client_guid", notificationcontroller.ClientNotificationController)
	r.POST("/notification/health", notificationcontroller.HealthBotController)
	errServer := r.Run(port_server)
	if errServer != nil {
		log.Error(errServer)
	}

}
