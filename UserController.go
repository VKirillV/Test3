package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

//const (
//	db_name = "tests"
//	db_host = "127.0.0.1"
//	db_user = "root"
//	db_pass = "!KV54691123s"
//	db_port = 3306
//)

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Println(e)
	}
}

func main() {

	var data_client Data_Client
	var data_admin Data_Admin
	db_name := os.Getenv("db_name")
	db_host := os.Getenv("db_host")
	db_user := os.Getenv("db_user")
	db_pass := os.Getenv("db_pass")
	db_port := os.Getenv("db_port")
	port_server := ":" + os.Getenv("port_server")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", string(db_user), string(db_pass), string(db_host), string(db_port), string(db_name))
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}

	//overtime times
	DB.SetConnMaxLifetime(3 * time.Minute)
	// maximum connection number
	DB.SetMaxOpenConns(100)
	// Set the number of idle connections
	DB.SetMaxIdleConns(16)
	if err := DB.Ping(); err != nil {
		log.Error("DB.Ping = ", err)
	}

	r := gin.Default()

	r.POST("/kript/:adminname/admin", func(c *gin.Context) {
		var adminname string = c.Param("adminname")

		tx, err := DB.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
			return
		}
		defer tx.Rollback()
		insert, err := DB.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", adminname, admin)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}
		defer insert.Close()

	})

	r.DELETE("/kript/:adminname/admin", func(c *gin.Context) {
		var adminname string = c.Param("adminname")
		_, err := DB.Query("UPDATE test.user SET usertype = (?) WHERE username = (?)", client, adminname)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}

	})

	r.POST("/:username/client/:guid", func(c *gin.Context) {
		var username string = c.Param("username")
		var guid string = c.Param("guid")
		fmt.Println(username)
		fmt.Println(guid)
		tx, err := DB.Begin()
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}
		defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

		insert, err := DB.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", username, client)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}
		defer insert.Close()

		insert, err = DB.Query("INSERT INTO client_user(client_guid) VALUES(?)", guid)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}
		defer insert.Close()
	})

	r.DELETE("/:username/client/:guid", func(c *gin.Context) {

	})

	r.GET("/admin", func(c *gin.Context) {

		rows, err := DB.Query("Select username FROM user WHERE usertype = (?)", admin)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error.",
			})
		}
		var datas []Data_Admin
		for rows.Next() {
			err := rows.Scan(&data_admin.Username)
			if err != nil {
				log.Error("Failed to connect to database !")
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error.",
				})
			}
			datas = append(datas, Data_Admin{Username: data_admin.Username})
		}
		c.JSON(http.StatusOK, datas)
	})

	r.GET("/client", func(c *gin.Context) {

		rows, err := DB.Query("Select username FROM user WHERE usertype = (?)", client)
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Service Error.",
			})

		}
		rows2, err := DB.Query("Select client_guid FROM client_user")
		if err != nil {
			log.Error("Failed to connect to database !")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Service Error.",
			})
			return
		}

		var datas []Data_Client
		for rows.Next() && rows2.Next() {
			err := rows.Scan(&data_client.Username)
			if err != nil {

				log.Error("Failed to connect to database !")

				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Service Error.",
				})
			}
			err2 := rows2.Scan(&data_client.Clients)
			if err2 != nil {
				log.Error("Failed to connect to database !")
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Service Error.",
				})

			}
			datas = append(datas, Data_Client{Username: data_client.Username, Usertype: string(client), Clients: data_client.Clients})

		}
		c.JSON(http.StatusOK, datas)
	})

	r.Run(port_server) // listen and serve	return DB
}

//func getDatabaseUser(username, guid string) (s *Server) {

// Add user in database
//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
//DB, err := sql.Open("mysql", dsn)
//tx, err := s.DB.Begin()
//if err != nil {
//	log.Fatal(err)
//}
//defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
//
//insert, err := s.DB.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", username, client)
//if err != nil {
//	log.Fatal(err)
//}
//defer insert.Close()
//
//insert, err = tx.Query("INSERT INTO client_user(client_guid) VALUES(?)", guid)
//if err != nil {
//	log.Fatal(err)
//}
//defer insert.Close()
//
//fmt.Println("Succesfully")

//}

//func getDatabaseAdmin(adminname string) (DB *sql.DB) {
//
//	insert, err := DB.Query("INSERT INTO user(username, usertype) VALUES(?, ?)", adminname, admin)
//	if err != nil {
//		panic(err.Error())
//	}
//	defer insert.Close()
//
//	fmt.Println("Succesfully")
//	return DB
//
//}
//
