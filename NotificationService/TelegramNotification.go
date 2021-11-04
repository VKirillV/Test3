package Notification

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Notification(chat_ID int64, TelegramUser string, bot_self_id int) {

	fmt.Println("NotificationBot is starting...")
	token := os.Getenv("NOTIFIC_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	reply := fmt.Sprintf("Registering Telegram User %s from chat %d", TelegramUser, bot_self_id)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, strconv.Itoa(int(chat_ID)), reply)
	_, err = http.Get(url)
	if err == nil {
		log.Infof("Message successfully delivered to %s", TelegramUser)
	} else if err != nil {
		log.Errorf("Message delivery failed to user %s with error: %s", TelegramUser, err)
	}
}
