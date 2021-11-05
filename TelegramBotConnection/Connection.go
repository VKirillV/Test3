package start

import (
	"os"

	"library/TelegramBot"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConnectBot() {
	tokenNotificatioin := os.Getenv("NOTIFICATION_TOKEN")
	token := os.Getenv("TOKEN")

	notificationBot, err := tgbotapi.NewBotAPI(tokenNotificatioin)
	if err != nil {
		log.Error(err)
	}
	startBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	startBot.Debug, notificationBot.Debug = true, true
	TelegramBot.Start(startBot)

}
