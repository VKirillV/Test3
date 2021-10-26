package ClientController

import (
	"database/sql"
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
	client UserType = "Client" //enum
	admin  UserType = "Admin"
)

var DB *sql.DB

func PostClientController(c *gin.Context) {
	var username string = c.Param("username")
	var guid string = c.Param("guid")
	DB := db.InitDB()
	tx, err := DB.Begin()
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	insert, err := tx.Prepare("INSERT INTO user(username, user_type) VALUES(?, ?)")
	if err != nil {
		log.Error("Failed to add the user to the data base! ", err)
		c.JSON(http.StatusInternalServerError, err)
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
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	id, err := DB.Query("Select id FROM user WHERE username = (?)", username)
	if err != nil {
		log.Error("Failed to connect to database! ", err)
	}

	for id.Next() {
		err := id.Scan(&data_client.Username)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
		}
		insert, err := tx.Prepare("DELETE client_user FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id= (?)")
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}
		defer insert.Close()
		_, err = insert.Exec(data_client.Username)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
			c.JSON(http.StatusInternalServerError, err)
		}
		tx.Commit()
	}
}

func GetClientController(c *gin.Context) {
	var data_client Data_Client
	var client_data Client_Data
	var query_parametrs Query_Parametrs
	var datas []Data_Client
	var query []Query_Parametrs
	client_guid := c.DefaultQuery("client_guid", "Guest")
	DB := db.InitDB()

	if client_guid == "Guest" {
		rows, err := DB.Query("Select id, username FROM user WHERE user_type = (?)", client)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
		}
		for rows.Next() {
			err := rows.Scan(&client_data.Clients_guid, &data_client.Username)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
			}
			rows2, err := DB.Query("Select client_guid FROM client_user INNER JOIN user ON user.id = client_user.user_fk WHERE user.id = (?)", client_data.Clients_guid)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
			}
			test := []string{}
			for rows2.Next() {
				err2 := rows2.Scan(&client_data.Guid_array)
				if err2 != nil {
					log.Error("Failed to connect to database! ", err)
				}
				test = append(test, client_data.Guid_array)
			}
			datas = append(datas, Data_Client{Username: string(data_client.Username), Usertype: string(client), Clients: test})
		}
		c.JSON(http.StatusOK, datas)
	} else {

		rows_test, err := DB.Query("Select user_fk FROM client_user WHERE client_guid = (?)", client_guid)
		if err != nil {
			log.Error("Failed to connect to database! ", err)

		}

		for rows_test.Next() {
			err := rows_test.Scan(&client_data.Clients_guid)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
			}
			rows_test2, err := DB.Query("Select username FROM user INNER JOIN client_user ON client_user.user_fk = user.id WHERE client_user.user_fk = (?)", &client_data.Clients_guid)
			if err != nil {
				log.Error("Failed to connect to database! ", err)
			}
			for rows_test2.Next() {
				err := rows_test2.Scan(&query_parametrs.Username)
				if err != nil {
					log.Error("Failed to connect to database! ", err)
				}
			}
			query = append(query, Query_Parametrs{Username: query_parametrs.Username})
		}
		c.JSON(http.StatusOK, query)
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
		log.Error("Failed to connect to database! ", err)
	}

	for id.Next() {
		err := id.Scan(&data_client.Username)
		if err != nil {
			log.Error("Failed to connect to database! ", err)
		}
	}

	insert2, err := tx.Prepare("INSERT INTO client_user(client_guid, user_FK) VALUES(?, ?)")
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}

	_, err = insert2.Exec(guid, data_client.Username)
	if err != nil {
		log.Error("Failed to connect to database! ", err)
		c.JSON(http.StatusInternalServerError, err)
	}
	tx.Commit()
	return
}
