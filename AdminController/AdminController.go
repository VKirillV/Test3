package admincontroller

import (
	db "library/ConnectionDatabase"
	Error "library/JsonError"
	take "library/SelectMethod"
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

func PostAdminController(ctx *gin.Context) {
	adminname := ctx.Param("name")

	txn, err := db.Connect().Begin()
	if Error.Error(ctx, err) {
		log.Error("Failed to connect to database! ", err)

		return
	}

	checkID := take.ID(txn, ctx, adminname)
	if checkID == nil {
		insert, err := txn.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
		if Error.Error(ctx, err) {
			log.Error("Failed to insert data in the database! ", err)

			return
		}
		defer insert.Close()

		_, err = insert.Exec(adminname, admin)
		if Error.Error(ctx, err) {
			log.Error("Failed to execute data in the database! ", err)

			return
		}
	} else {
		update, err := txn.Prepare("UPDATE user SET user_type = (?) WHERE username = (?)")
		if Error.Error(ctx, err) {
			log.Error("Failed to update data in the database! ", err)

			return
		}
		defer update.Close()
		_, err = update.Exec(admin, adminname)
		if Error.Error(ctx, err) {
			log.Error("Failed to execute data in the database! ", err)

			return
		}
	}

	err = txn.Commit()
	if err != nil {
		log.Error(err)
	}
}

func DeleteAdminController(ctx *gin.Context) {
	adminname := ctx.Param("name")

	txn, err := db.Connect().Begin()
	if Error.Error(ctx, err) {
		log.Error("Failed to connect to database! ", err)

		return
	}

	update, err := txn.Prepare("UPDATE user SET user_type = (?) WHERE username = (?)")
	if Error.Error(ctx, err) {
		log.Error("Failed to update data in the database! ", err)

		return
	}
	defer update.Close()

	_, err = update.Exec(client, adminname)
	if Error.Error(ctx, err) {
		log.Error("Failed to execute data in the database! ", err)

		return
	}

	err = txn.Commit()
	if err != nil {
		log.Error(err)
	}
}

func GetAdminController(ctx *gin.Context) {
	var dataAdmin DataAdmin

	txn, err := db.Connect().Begin()
	if Error.Error(ctx, err) {
		log.Error("Failed to connect to database! ", err)

		return
	}

	rows, err := txn.Query("Select username FROM user WHERE user_type = (?)", admin)
	if Error.Error(ctx, err) {
		log.Error("Failed to select certain data in the database! ", err)

		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var allAdmin []DataAdmin

	for rows.Next() {
		err := rows.Scan(&dataAdmin.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}

		allAdmin = append(allAdmin, DataAdmin{Username: dataAdmin.Username})
	}
	ctx.JSON(http.StatusOK, allAdmin)
}
