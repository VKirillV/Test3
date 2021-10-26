package main

import (
	"log"
	"os"

	"library/AdminController"
	"library/ClientController"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func init() {

	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}
}

func main() {

	port_server := ":" + os.Getenv("port_server")

	r := gin.Default()

	r.POST("/kript/:adminname/admin", AdminController.PostAdminController)
	r.DELETE("/kript/:adminname/admin", AdminController.DeleteAdminController)
	r.POST("/:username/client/:guid", ClientController.PostClientController)
	r.DELETE("/:username/client/:guid", ClientController.DeleteClientController)
	r.GET("/admin", AdminController.GetAdminController)
	r.GET("/client", ClientController.GetClientController)

	r.Run(port_server) // listen and serve	return DB
}
