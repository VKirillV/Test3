package start

import (
	"os"

	"library/TelegramBot"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConnectBot() {
	token := os.Getenv("TOKEN")
	startBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	startBot.Debug = true
	TelegramBot.Start(startBot)

}
