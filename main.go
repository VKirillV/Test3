package main

import (
	"fmt"
	"library/AdminController"
	"library/ClientController"
	start "library/TelegramBotConnection"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
)

func main() {
	go start.ConnectBot()
	fmt.Println("Server is starting...")

	port_server := ":" + os.Getenv("port_server")

	r := gin.Default()

	r.POST("/kript/:adminname/admin", AdminController.PostAdminController)
	r.DELETE("/kript/:adminname/admin", AdminController.DeleteAdminController)
	r.POST("/:username/client/:guid", ClientController.PostClientController)
	r.DELETE("/:username/client/:guid", ClientController.DeleteClientController)
	r.GET("/admin", AdminController.GetAdminController)
	r.GET("/client", ClientController.GetClientController)

	r.Run(port_server)

}
