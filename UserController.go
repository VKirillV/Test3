package main

import (
	"database/sql"
	"library/db"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Data_Admin struct {
	Username string `json:"username"`
}

type Data_Client struct {
	Username string `json:"username"`
	Usertype string `json:"usertype"`
	Clients  string `json:"clients"`
}

type UserType string

const (
	client UserType = "Client" //enum
	admin  UserType = "Admin"
)

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}
}

var DB *sql.DB

func main() {

	var data_client Data_Client
	var data_admin Data_Admin
	ErrorMessage := "Internal server Error"

	DB := db.InitDB()

	port_server := ":" + os.Getenv("port_server")

	r := gin.Default()

	r.POST("/kript/:adminname/admin", func(c *gin.Context) {
		var adminname string = c.Param("adminname")
		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
			return
		}
		defer tx.Rollback()

		insert, err := tx.Prepare("INSERT INTO user(username, usertype) VALUES(?, ?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}

		defer insert.Close()

		_, err = insert.Exec(adminname, admin)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		tx.Commit()
	})

	r.DELETE("/kript/:adminname/admin", func(c *gin.Context) {
		var adminname string = c.Param("adminname")

		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
			return
		}
		defer tx.Rollback()

		insert, err := tx.Prepare("UPDATE test.user SET usertype = (?) WHERE username = (?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		defer insert.Close()

		_, err = insert.Exec(client, adminname)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		tx.Commit()
	})

	r.POST("/:username/client/:guid", func(c *gin.Context) {
		var username string = c.Param("username")
		var guid string = c.Param("guid")

		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
			return
		}

		defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

		insert, err := tx.Prepare("INSERT INTO user(username, usertype) VALUES(?, ?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		defer insert.Close()

		_, err = insert.Exec(username, client)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}

		insert2, err := tx.Prepare("INSERT INTO client_user(client_guid) VALUES(?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		_, err = insert2.Exec(guid)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		tx.Commit()
	})

	r.DELETE("/:username/client/:guid", func(c *gin.Context) {

	})

	r.GET("/admin", func(c *gin.Context) {
		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
			return
		}
		defer tx.Rollback()

		rows, err := tx.Query("Select username FROM user WHERE usertype = (?)", admin)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		var datas []Data_Admin
		for rows.Next() {
			err := rows.Scan(&data_admin.Username)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
				c.JSON(http.StatusInternalServerError, ErrorMessage)
			}
			datas = append(datas, Data_Admin{Username: data_admin.Username})
		}
		tx.Commit()
		c.JSON(http.StatusOK, datas)
	})

	r.GET("/client", func(c *gin.Context) {
		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
			return
		}
		rows, err := tx.Query("Select username FROM user WHERE usertype = (?)", client)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, ErrorMessage)
		}
		var datas []Data_Client
		for rows.Next() { //&& rows2.Next() {
			err := rows.Scan(&data_client.Username)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
				c.JSON(http.StatusInternalServerError, ErrorMessage)
			}
			datas = append(datas, Data_Client{Username: data_client.Username, Usertype: string(client), Clients: data_client.Clients})
		}
		tx.Commit()
		c.JSON(http.StatusOK, datas)
	})
	r.Run(port_server) // listen and serve	return DB
}
