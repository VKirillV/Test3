package ClientController

import (
	Error "library/JsonError"
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DataClient struct {
	Username string   `json:"username"`
	Usertype string   `json:"usertype"`
	Clients  []string `json:"clients"`
}

type ClientData struct {
	ClientsGuid string `json:"clients"`
	GuidArray   string
}

type QueryParametrs struct {
	Username string `json:"username"`
}
type UserFk struct {
	UserFk int
}
type UserType string

const (
	client UserType = "Client"
	admin  UserType = "Admin"
)

func PostClientController(c *gin.Context) {
	var username string = c.Param("username")
	var guid string = c.Param("guid")
	var dataClient DataClient

	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

	defer tx.Rollback()

	id, err := tx.Query("Select username FROM user WHERE username = (?)", username)
	if Error.Error(c, err) {
		log.Error("Failed to select certain data in the database! ", err)
		return
	}

	for id.Next() {
		err := id.Scan(&dataClient.Username)
		if Error.Error(c, err) {
			log.Error("The structures does not match! ", err)
			return
		}
	}

	if dataClient.Username == username {

		GetId(dataClient.Username, guid)

	} else if dataClient.Username != username {

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

		defer GetId(username, guid)

		tx.Commit()

	}
}

func DeleteClientController(c *gin.Context) {
	var dataClient DataClient
	var username string = c.Param("username")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

	defer tx.Rollback()

	id, err := tx.Query("Select id FROM user WHERE username = (?)", username)
	if Error.Error(c, err) {
		log.Error("Failed to select certain data in the database! ", err)
		return
	}

	for id.Next() {
		err := id.Scan(&dataClient.Username)
		if Error.Error(c, err) {
			log.Error("The structures does not match! ", err)
			return
		}
		delete, err := db.Connect().Prepare("DELETE client_user FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id= (?)")
		if Error.Error(c, err) {
			log.Error("Failed to delete data in the database! ", err)
			return
		}
		defer delete.Close()
		_, err = delete.Exec(dataClient.Username)
		if Error.Error(c, err) {
			log.Error("Failed to execute data in the database! ", err)
			return
		}
		tx.Commit()
	}
}

func GetClientController(c *gin.Context) {
	var dataClient DataClient
	var clientData ClientData
	var queryParametrs QueryParametrs
	var allClient []DataClient
	var allClientGuid []QueryParametrs
	clientGuid := c.DefaultQuery("client_guid", "Guest")
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

	if clientGuid == "Guest" {
		rows, err := tx.Query("Select i, username FROM user WHERE user_type = (?)", client)
		if Error.Error(c, err) {
			log.Error("Failed to select certain data in the database! ", err)
			return
		}
		for rows.Next() {
			err := rows.Scan(&clientData.ClientsGuid, &dataClient.Username)
			if Error.Error(c, err) {
				log.Error("The structures does not match! ", err)
				return
			}
			rows2, err := db.Connect().Query("Select client_guid FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id = (?)", clientData.ClientsGuid)
			if Error.Error(c, err) {
				log.Error("Failed to select certain data in the database! ", err)
				return
			}
			allClientGuid := []string{}
			for rows2.Next() {
				err2 := rows2.Scan(&clientData.GuidArray)
				if Error.Error(c, err2) {
					log.Error("The structures does not match! ", err)
					return
				}
				allClientGuid = append(allClientGuid, clientData.GuidArray)
			}
			allClient = append(allClient, DataClient{Username: string(dataClient.Username), Usertype: string(client), Clients: allClientGuid})
		}
		c.JSON(http.StatusOK, allClient)
	} else {

		rowsGuid, err := tx.Query("Select user_fk FROM client_user WHERE client_guid = (?)", clientGuid)
		if Error.Error(c, err) {
			log.Error("Failed to select certain data in the database! ", err)
			return
		}

		for rowsGuid.Next() {
			err := rowsGuid.Scan(&clientData.ClientsGuid)
			if Error.Error(c, err) {
				log.Error("The structures does not match! ", err)
				return
			}
			rowsUserFk, err := db.Connect().Query("Select username FROM user INNER JOIN client_user ON client_user.user_fk = user.id WHERE client_user.user_fk = (?)", &clientData.ClientsGuid)
			if Error.Error(c, err) {
				log.Error("Failed to select certain data in the database! ", err)
				return
			}
			for rowsUserFk.Next() {
				err := rowsUserFk.Scan(&queryParametrs.Username)
				if Error.Error(c, err) {
					log.Error("The structures does not match! ", err)
					return
				}
			}
			allClientGuid = append(allClientGuid, QueryParametrs{Username: queryParametrs.Username})
		}
		c.JSON(http.StatusOK, allClientGuid)
	}

}

func GetUserfk(guid string, id int) (c *gin.Context) {

	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
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

	return
}

func GetId(username, guid string) (c *gin.Context) {
	var userFk UserFk
	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}
	id, err := tx.Query("Select id FROM user WHERE username = (?)", username)
	if Error.Error(c, err) {
		log.Error("Failed to select certain data in the database! ", err)
		return
	}

	for id.Next() {
		err := id.Scan(&userFk.UserFk)
		if Error.Error(c, err) {
			log.Error("The structures does not match! ", err)
			return
		}
		GetUserfk(guid, userFk.UserFk)
	}
	return
}
