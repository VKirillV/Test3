package AdminController

import (
	"database/sql"
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
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

var DB *sql.DB

func PostAdminController(c *gin.Context) {
	var data_admin Data_Admin
	DB := db.InitDB()
	var adminname string = c.Param("adminname")
	tx, err := DB.Begin()
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	rows, err := DB.Query("Select username FROM user WHERE username = (?)", adminname)
	if err != nil {
		log.Error("Failed to connect to database!", err)
	}

	for rows.Next() {
		err := rows.Scan(&data_admin.Username)
		if err != nil {
			log.Error("Failed to connect to database!", err)
		}
	}

	if data_admin.Username == adminname {
		// change usertype
		insert, err := DB.Prepare("UPDATE test.user SET user_type = (?) WHERE username = (?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}
		defer insert.Close()

		_, err = insert.Exec(admin, adminname)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}

	} else if data_admin.Username != adminname {

		insert, err := tx.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}

		defer insert.Close()

		_, err = insert.Exec(adminname, admin)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}
		tx.Commit()

	}
}

func DeleteAdminController(c *gin.Context) {
	var adminname string = c.Param("adminname")
	DB := db.InitDB()
	tx, err := DB.Begin()
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()
	insert, err := tx.Prepare("UPDATE test2.user SET user_type = (?) WHERE username = (?)")
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	defer insert.Close()

	_, err = insert.Exec(client, adminname)
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	tx.Commit()
}

func GetAdminController(c *gin.Context) {
	var data_admin Data_Admin
	DB := db.InitDB()

	rows, err := DB.Query("Select username FROM user WHERE user_type = (?)", admin)
	if err != nil {
		log.Error("Failed to connect to database! ", err)
	}
	var datas []Data_Admin
	for rows.Next() {
		err := rows.Scan(&data_admin.Username)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
		}
		datas = append(datas, Data_Admin{Username: data_admin.Username})
	}
	c.JSON(http.StatusOK, datas)

}
