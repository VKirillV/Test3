package Select

import (
	"database/sql"
	Error "library/JsonError"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ID(tx *sql.Tx, c *gin.Context, username string) (id *int) {
	userId, err := tx.Query("Select id FROM user WHERE username = (?)", username)
	if Error.Error(c, err) {
		log.Error("Failed to selct certain data in the database! ", err)
		return
	}
	for userId.Next() {
		err := userId.Scan(&id)
		if Error.Error(c, err) {
			log.Error("The structuesdoes not match! ", err)
			return
		}
	}
	return id
}
