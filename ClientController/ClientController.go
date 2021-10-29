package ClientController

import (
	"database/sql"
	Error "library/JsonError"
	"library/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Data_Client struct {
	Username string   `json:"username"`
	Usertype string   `json:"usertype"`
	Clients  []string `json:"clients"`
}

type Client_Data struct {
	Clients_guid string `json:"clients"`
	Guid_array   string
}

type Query_Parametrs struct {
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
	DB := db.InitDB()
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
	defer User_fk(username, guid)
	tx.Commit()
}

func DeleteClientController(c *gin.Context) {
	var data_client Data_Client
	var username string = c.Param("username")
	DB := db.InitDB()
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
		err := id.Scan(&data_client.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		delete, err := tx.Prepare("DELETE client_user FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id= (?)")
		if Error.Error(c, err) {
			log.Error("Failed to delete data in the database! ", err)
			return
		}
		defer delete.Close()
		_, err = delete.Exec(data_client.Username)
		if Error.Error(c, err) {
			log.Error("Failed to execute data in the database! ", err)
			return
		}
		tx.Commit()
	}
}

func GetClientController(c *gin.Context) {
	var data_client Data_Client
	var client_data Client_Data
	var query_parametrs Query_Parametrs
	var All_Client []Data_Client
	var All_ClientGuid []Query_Parametrs
	client_guid := c.DefaultQuery("client_guid", "Guest")
	DB := db.InitDB()

	if client_guid == "Guest" {
		rows, err := DB.Query("Select id, username FROM user WHERE user_type = (?)", client)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)
		}
		for rows.Next() {
			err := rows.Scan(&client_data.Clients_guid, &data_client.Username)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
			rows2, err := DB.Query("Select client_guid FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id = (?)", client_data.Clients_guid)
			if err != nil {
				log.Error("Failed to select certain data in the database! ", err)
			}
			All_ClientGuid := []string{}
			for rows2.Next() {
				err2 := rows2.Scan(&client_data.Guid_array)
				if err2 != nil {
					log.Error("The structures does not match! ", err)
				}
				All_ClientGuid = append(All_ClientGuid, client_data.Guid_array)
			}
			All_Client = append(All_Client, Data_Client{Username: string(data_client.Username), Usertype: string(client), Clients: All_ClientGuid})
		}
		c.JSON(http.StatusOK, All_Client)
	} else {

		rows_guid, err := DB.Query("Select user_fk FROM client_user WHERE client_guid = (?)", client_guid)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)

		}

		for rows_guid.Next() {
			err := rows_guid.Scan(&client_data.Clients_guid)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
			rows_UserFK, err := DB.Query("Select username FROM user INNER JOIN client_user ON client_user.user_fk = user.id WHERE client_user.user_fk = (?)", &client_data.Clients_guid)
			if err != nil {
				log.Error("Failed to select certain data in the database! ", err)
			}
			for rows_UserFK.Next() {
				err := rows_UserFK.Scan(&query_parametrs.Username)
				if err != nil {
					log.Error("The structures does not match! ", err)
				}
			}
			All_ClientGuid = append(All_ClientGuid, Query_Parametrs{Username: query_parametrs.Username})
		}
		c.JSON(http.StatusOK, All_ClientGuid)
	}

}

func User_fk(username, guid string) (c *gin.Context) {
	var data_client Data_Client
	DB := db.InitDB()
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
		err := id.Scan(&data_client.Username)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
	}

	insert2, err := tx.Prepare("INSERT INTO client_user(client_guid, user_FK) VALUES(?, ?)")
	if err != nil {
		log.Error("Failed to insert data in the database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}

	_, err = insert2.Exec(guid, data_client.Username)
	if err != nil {
		log.Error("Failed to execute data in the database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	tx.Commit()
	return
}
