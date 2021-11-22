package ClientController

import (
	Error "library/JsonError"
	Select "library/SelectMethod"

	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DataClient struct {
	Username string `json:"username"`
}

type ClientData struct {
	ClientsGuid string `json:"clients"`
	GuidArray   string
}

type QueryParametrs struct {
	Username string `json:"username"`
}

type UserType string

const (
	client UserType = "Client"
	admin  UserType = "Admin"
)

func PostClientController(c *gin.Context) {
	var username string = c.Param("username")
	var guid string = c.Param("guid")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()
	id := Select.ID(tx, c, username)
	if id == nil {
		insert, err := tx.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
		if Error.Error(c, err) {
			log.Error("Failed to insert data in the database! ", err)
			return
		}
		defer insert.Close()
		_, err = insert.Exec(username, client)
		if Error.Error(c, err) {
			log.Error("Failed to execute data in the database! ", err)
			return
		}
		id = Select.ID(tx, c, username)
	}

	insert2, err := tx.Prepare("INSERT INTO client_user(client_guid, user_fk) VALUES(?, ?)")
	if Error.Error(c, err) {
		log.Error("Failed to insert data in the database! ", err)
		return
	}
	defer insert2.Close()
	_, err = insert2.Exec(guid, id)
	if Error.Error(c, err) {
		log.Error("Failed to execute data in the database! ", err)
		return
	}
	tx.Commit()
}

func DeleteClientController(c *gin.Context) {
	var username string = c.Param("username")
	var guid string = c.Param("guid")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	defer tx.Rollback()
	statement, err := db.Connect().Prepare("DELETE client_user FROM client_user " +
		"JOIN user ON user.id = client_user.user_fk " +
		"WHERE user.username = (?) AND client_user.client_guid = (?)")
	if Error.Error(c, err) {
		log.Error("Failed to delete data in the database! ", err)
		return
	}
	defer statement.Close()
	_, err = statement.Exec(username, guid)
	if Error.Error(c, err) {
		log.Error("Failed to execute data in the database! ", err)
		return
	}
	tx.Commit()
}

func GetClientController(c *gin.Context) {
	var dataClient DataClient
	var queryParametrs QueryParametrs
	var allClient []DataClient
	var allClientGuid []QueryParametrs
	clientGuid := c.Query("client_guid")

	if len(clientGuid) == 0 {

		rows2, err := db.Connect().Query(
			"Select DISTINCT username FROM user "+
				"LEFT JOIN client_user ON user.id = client_user.user_fk AND user_type = (?) "+
				"WHERE client_user.user_fk IS NOT NULL", client)
		if Error.Error(c, err) {
			log.Error("Failed to select certain data in the database! ", err)
			return
		}
		for rows2.Next() {
			err2 := rows2.Scan(&dataClient.Username)
			if Error.Error(c, err2) {
				log.Error("The structures does not match! ", err)
				return
			}
			allClient = append(allClient, DataClient{Username: string(dataClient.Username)})
		}
		c.JSON(http.StatusOK, allClient)
	} else {
		rows, err := db.Connect().Query("Select DISTINCT username FROM user "+
			"LEFT JOIN client_user ON user.id = client_user.user_fk "+
			"WHERE client_user.client_guid = (?)", clientGuid)
		if Error.Error(c, err) {
			log.Error("Failed to select certain data in the database! ", err)
			return
		}
		for rows.Next() {
			err := rows.Scan(&queryParametrs.Username)
			if Error.Error(c, err) {
				log.Error("The structures does not match! ", err)
				return
			}
			allClientGuid = append(allClientGuid, QueryParametrs{Username: queryParametrs.Username})
		}
		c.JSON(http.StatusOK, allClientGuid)
	}
}
