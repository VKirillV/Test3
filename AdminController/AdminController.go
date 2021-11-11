package AdminController

import (
	Error "library/JsonError"
	Select "library/SelectMethod"
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type DataAdmin struct {
	Username string `json:"username"`
}

type UserType string

const (
	client UserType = "Client"
	admin  UserType = "Admin"
)

func PostAdminController(c *gin.Context) {
	var adminname string = c.Param("adminname")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()
	id := Select.ID(tx, c, adminname)
	if id == nil {
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
	} else if id != nil {
		update, err := tx.Prepare("UPDATE user SET user_type = (?) WHERE username = (?)")
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
	}
	tx.Commit()
}

func DeleteAdminController(c *gin.Context) {
	var adminname string = c.Param("adminname")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()
	update, err := tx.Prepare("UPDATE user SET user_type = (?) WHERE username = (?)")
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
	var dataAdmin DataAdmin
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	rows, err := tx.Query("Select username FROM user WHERE user_type = (?)", admin)
	if Error.Error(c, err) {
		log.Error("Failed to select certain data in the database! ", err)
		return
	}
	var allAdmin []DataAdmin
	for rows.Next() {
		err := rows.Scan(&dataAdmin.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		allAdmin = append(allAdmin, DataAdmin{Username: dataAdmin.Username})
	}
	c.JSON(http.StatusOK, allAdmin)
}
