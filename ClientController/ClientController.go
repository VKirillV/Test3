package ClientController

import (
	"database/sql"
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
type UserType string

const (
	client UserType = "Client"
	admin  UserType = "Admin"
)

var DB *sql.DB

func PostClientController(c *gin.Context) {
	var username string = c.Param("username")
	var guid string = c.Param("guid")
	DB := db.Connect()
	tx, err := DB.Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

	defer tx.Rollback()

	insert, err := tx.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
	if Error.Error(c, err) {
		log.Error("Failed to insert data in the database! ", err)
		return
	}
	defer insert.Close()
	_, _ = insert.Exec(username, client)
	defer GetUserfk(username, guid)
	tx.Commit()
}

func DeleteClientController(c *gin.Context) {
	var dataClient DataClient
	var username string = c.Param("username")
	DB := db.Connect()
	tx, err := DB.Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

	defer tx.Rollback()

	id, err := DB.Query("Select id FROM user WHERE username = (?)", username)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}

	for id.Next() {
		err := id.Scan(&dataClient.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		delete, err := tx.Prepare("DELETE client_user FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id= (?)")
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
	DB := db.Connect()

	if clientGuid == "Guest" {
		rows, err := DB.Query("Select id, username FROM user WHERE user_type = (?)", client)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)
		}
		for rows.Next() {
			err := rows.Scan(&clientData.ClientsGuid, &dataClient.Username)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
			rows2, err := DB.Query("Select client_guid FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id = (?)", clientData.ClientsGuid)
			if err != nil {
				log.Error("Failed to select certain data in the database! ", err)
			}
			allClientGuid := []string{}
			for rows2.Next() {
				err2 := rows2.Scan(&clientData.GuidArray)
				if err2 != nil {
					log.Error("The structures does not match! ", err)
				}
				allClientGuid = append(allClientGuid, clientData.GuidArray)
			}
			allClient = append(allClient, DataClient{Username: string(dataClient.Username), Usertype: string(client), Clients: allClientGuid})
		}
		c.JSON(http.StatusOK, allClient)
	} else {

		rowsGuid, err := DB.Query("Select user_fk FROM client_user WHERE client_guid = (?)", clientGuid)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)

		}

		for rowsGuid.Next() {
			err := rowsGuid.Scan(&clientData.ClientsGuid)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
			rowsUserFk, err := DB.Query("Select username FROM user INNER JOIN client_user ON client_user.user_fk = user.id WHERE client_user.user_fk = (?)", &clientData.ClientsGuid)
			if err != nil {
				log.Error("Failed to select certain data in the database! ", err)
			}
			for rowsUserFk.Next() {
				err := rowsUserFk.Scan(&queryParametrs.Username)
				if err != nil {
					log.Error("The structures does not match! ", err)
				}
			}
			allClientGuid = append(allClientGuid, QueryParametrs{Username: queryParametrs.Username})
		}
		c.JSON(http.StatusOK, allClientGuid)
	}

}

func GetUserfk(username, guid string) (c *gin.Context) {
	var dataClient DataClient
	DB := db.Connect()
	tx, err := DB.Begin()
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	id, err := DB.Query("Select id FROM user WHERE username = (?)", username)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}

	for id.Next() {
		err := id.Scan(&dataClient.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
	}

	insert2, err := tx.Prepare("INSERT INTO client_user(client_guid, user_FK) VALUES(?, ?)")
	if err != nil {
		log.Error("Failed to insert data in the database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}

	_, err = insert2.Exec(guid, dataClient.Username)
	if err != nil {
		log.Error("Failed to execute data in the database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	tx.Commit()
	return
}
