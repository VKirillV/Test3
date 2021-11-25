package take

import (
	"database/sql"
	Error "library/JsonError"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ID(txn *sql.Tx, ctx *gin.Context, username string) (id *int) {
	userID, err := txn.Query("Select id FROM user WHERE username = (?)", username)
	if Error.Error(ctx, err) {
		log.Error("Failed to select certain data in the database! ", err)

		return
	}

	defer func() {
		_ = userID.Close()
		_ = userID.Err()
	}()

	for userID.Next() {
		err := userID.Scan(&id)
		if Error.Error(ctx, err) {
			log.Error("The structuesdoes not match! ", err)

			return
		}
	}

	return id
}
