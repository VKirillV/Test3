package clientcontroller

import (
	Error "library/JsonError"
	take "library/SelectMethod"
	"net/http"

	db "library/ConnectionDatabase"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DataClient struct {
	Username string `json:"username"`
}

type ClientData struct {
	ClientsGUID string `json:"clients"`
	GUIDArray   string
}

type QueryParametrs struct {
	Username string `json:"username"`
}

type UserType string

const (
	client UserType = "Client"
)

func PostClientController(ctx *gin.Context) {
	username := ctx.Param("username")
	guid := ctx.Param("guid")

	txn, err := db.Connect().Begin()
	if Error.Error(ctx, err) {
		log.Error("Failed to connect to database! ", err)

		return
	}

	defer func() {
		err = txn.Rollback()
		if err != nil {
			log.Error(err)
		}
	}()

	selectID := take.ID(txn, ctx, username)
	if selectID == nil {
		insert, err := txn.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
		if Error.Error(ctx, err) {
			log.Error("Failed to insert data in the database! ", err)

			return
		}
		defer insert.Close()

		_, err = insert.Exec(username, client)
		if Error.Error(ctx, err) {
			log.Error("Failed to execute data in the database! ", err)

			return
		}

		selectID = take.ID(txn, ctx, username)
	}

	insert2, err := txn.Prepare("INSERT INTO client_user(client_guid, user_fk) VALUES(?, ?)")
	if Error.Error(ctx, err) {
		log.Error("Failed to insert data in the database! ", err)

		return
	}
	defer insert2.Close()

	_, err = insert2.Exec(guid, selectID)
	if Error.Error(ctx, err) {
		log.Error("Failed to execute data in the database! ", err)

		return
	}

	err = txn.Commit()
	if err != nil {
		log.Error(err)
	}
}

func DeleteClientController(ctx *gin.Context) {
	username := ctx.Param("username")
	guid := ctx.Param("guid")

	txn, err := db.Connect().Begin()
	if Error.Error(ctx, err) {
		log.Error("Failed to connect to database! ", err)

		return
	}

	defer func() {
		err = txn.Rollback()
		if err != nil {
			log.Error(err)
		}
	}()

	statement, err := db.Connect().Prepare("DELETE client_user FROM client_user " +
		"JOIN user ON user.id = client_user.user_fk " +
		"WHERE user.username = (?) AND client_user.client_guid = (?)")
	if Error.Error(ctx, err) {
		log.Error("Failed to delete data in the database! ", err)

		return
	}
	defer statement.Close()

	_, err = statement.Exec(username, guid)
	if Error.Error(ctx, err) {
		log.Error("Failed to execute data in the database! ", err)

		return
	}

	err = txn.Commit()
	if err != nil {
		log.Error(err)
	}
}

func GetClientController(ctx *gin.Context) {
	var dataClient DataClient

	var queryParametrs QueryParametrs

	var allClient []DataClient

	var allClientGUID []QueryParametrs

	clientGUID, _ := ctx.GetQuery("client_guid")

	if len(clientGUID) == 0 {

		rows2, err := db.Connect().Query(
			"Select DISTINCT username FROM user "+
				"LEFT JOIN client_user ON user.id = client_user.user_fk AND user_type = (?) "+
				"WHERE client_user.user_fk IS NOT NULL", client)
		if Error.Error(ctx, err) {
			log.Error("Failed to select certain data in the database! ", err)

			return
		}

		defer func() {
			_ = rows2.Close()
			_ = rows2.Err()
		}()

		for rows2.Next() {
			err2 := rows2.Scan(&dataClient.Username)
			if Error.Error(ctx, err2) {
				log.Error("The structures does not match! ", err)

				return
			}

			allClient = append(allClient, DataClient{Username: dataClient.Username})
		}
		ctx.JSON(http.StatusOK, allClient)
	} else {
		rows, err := db.Connect().Query("Select DISTINCT username FROM user "+
			"LEFT JOIN client_user ON user.id = client_user.user_fk "+
			"WHERE client_user.client_guid = (?)", clientGUID)
		if Error.Error(ctx, err) {
			log.Error("Failed to select certain data in the database! ", err)

			return
		}
		defer func() {
			_ = rows.Close()
			_ = rows.Err()
		}()
		for rows.Next() {
			err := rows.Scan(&queryParametrs.Username)
			if Error.Error(ctx, err) {
				log.Error("The structures does not match! ", err)

				return
			}
			allClientGUID = append(allClientGUID, QueryParametrs{Username: queryParametrs.Username})
		}
		ctx.JSON(http.StatusOK, allClientGUID)
	}
}
