package telegrambot

import (
	"database/sql"
	"fmt"
	db "library/ConnectionDatabase"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username       string
	TelegramChatID int
}

func Start(startBot *tgbotapi.BotAPI, tx *sql.Tx, c *gin.Context) {
	log.Info("GoGinBot is starting...")

	u := tgbotapi.NewUpdate(0)

	updates, err := startBot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		telegramUser := update.Message.From.UserName

		if len(telegramUser) == 0 {
			continue
		}

		messageChatID := update.Message.Chat.ID

		username, telegramChatID := PrepareChatID(telegramUser)
		if username == telegramUser {
			if telegramChatID == 0 {
				update, err := db.Connect().Prepare("UPDATE user SET telegram_chat_id = (?) WHERE username = (?)")
				if err != nil {
					log.Error("Failed to update data in the database! ", err)
				}

				_, err = update.Exec(messageChatID, username)
				if err != nil {
					log.Error("Failed to update data in the database! ", err)
				}

				update.Close()

				msg := tgbotapi.NewMessage(messageChatID, "Successfully subscribed on updates")

				_, err = startBot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", telegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
				}
			} else {
				msg := tgbotapi.NewMessage(messageChatID, "You are registered")
				_, err = startBot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", telegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
				}
			}
		}
	}
}

func SendMessage(notification string, telegramUser string, messageChatID int64) {
	token := os.Getenv("TOKEN")

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", token, messageChatID, notification)

	_, err := http.Get(url)
	if err == nil {
		log.Infof("Message successfully delivered to %s", telegramUser)
	} else if err != nil {
		log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
	}
}

func PrepareChatID(telegramUser string) (username string, telegramChatID int) {
	var data Data

	rows, err := db.Connect().Query("Select username, telegram_chat_id FROM user WHERE username = (?)", telegramUser)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		err := rows.Scan(&data.Username, &data.TelegramChatID)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
	}

	return data.Username, data.TelegramChatID
}

func EscapeMessage(notification string) (newNotification string) {
	re := regexp.MustCompile(`[[:punct:]]`)
	newNotification = re.ReplaceAllString(notification, "")

	return newNotification
}
