package start

import (
	"database/sql"
	"os"

	Error "library/JsonError"

	"github.com/gin-gonic/gin"

	db "library/ConnectionDatabase"

	log "github.com/sirupsen/logrus"

	TelegramBot "library/TelegramBot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConnectBot() (startBot *tgbotapi.BotAPI, tx *sql.Tx, c *gin.Context) {
	token := os.Getenv("TOKEN")

	startBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	tx, err = db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	startBot.Debug = true
	TelegramBot.Start(startBot, tx, c)

	return startBot, tx, c
}
