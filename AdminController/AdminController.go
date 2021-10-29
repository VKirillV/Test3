package AdminController

import (
	"database/sql"
	Error "library/JsonError"
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
	client UserType = "Client"
	admin  UserType = "Admin"
)

var DB *sql.DB

func PostAdminController(c *gin.Context) {
	var data_admin Data_Admin
	DB := db.InitDB()
	var adminname string = c.Param("adminname")
	tx, err := DB.Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()

	rows, err := DB.Query("Select username FROM user WHERE username = (?)", adminname)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}

	for rows.Next() {
		err := rows.Scan(&data_admin.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
	}

	if data_admin.Username == adminname {
		update, err := DB.Prepare("UPDATE test2.user SET user_type = (?) WHERE username = (?)")
		if Error.Error(c, err) {
			log.Error("Failed to update data in the database! ", err)
			return
		}
		defer update.Close()

		_, err = update.Exec(admin, adminname)
		if Error.Error(c, err) {
			log.Error("Failed to execute data in the database! ", err)
			return
		}

	} else if data_admin.Username != adminname {

		insert, err := tx.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
		if Error.Error(c, err) {
			log.Error("Failed to insert data in the database! ", err)
			return
		}

		defer insert.Close()

		_, err = insert.Exec(adminname, admin)
		if Error.Error(c, err) {
			log.Error("Failed to execute data in the database! ", err)
			return
		}
		tx.Commit()

	}
}

func DeleteAdminController(c *gin.Context) {
	var adminname string = c.Param("adminname")
	DB := db.InitDB()
	tx, err := DB.Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()
	update, err := tx.Prepare("UPDATE test2.user SET user_type = (?) WHERE username = (?)")
	if Error.Error(c, err) {
		log.Error("Failed to update data in the database! ", err)
		return
	}
	defer update.Close()

	_, err = update.Exec(client, adminname)
	if Error.Error(c, err) {
		log.Error("Failed to execute data in the database! ", err)
		return
	}
	tx.Commit()
}

func GetAdminController(c *gin.Context) {
	var data_admin Data_Admin
	DB := db.InitDB()

	rows, err := DB.Query("Select username FROM user WHERE user_type = (?)", admin)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}
	var All_Admin []Data_Admin
	for rows.Next() {
		err := rows.Scan(&data_admin.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		All_Admin = append(All_Admin, Data_Admin{Username: data_admin.Username})
	}
	c.JSON(http.StatusOK, All_Admin)

}
