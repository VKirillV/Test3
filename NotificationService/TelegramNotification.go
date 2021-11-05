package Notification

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func SendMessage(messageChatId int64, telegramUser string, botSelfId int) {

	log.Info("NotificationBot is starting...")
	token := os.Getenv("NOTIFICATION_TOKEN")

	reply := fmt.Sprintf("Registering Telegram User %s from chat %d", telegramUser, botSelfId)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", token, strconv.Itoa(int(messageChatId)), reply)
	_, err := http.Get(url)
	if err == nil {
		log.Infof("Message successfully delivered to %s", telegramUser)
	} else if err != nil {
		log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
	}
}
