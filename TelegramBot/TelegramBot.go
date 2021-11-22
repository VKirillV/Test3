package TelegramBot

import (
	"database/sql"
	"fmt"
	"library/db"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username       string
	TelegramChatId int
}

func Start(startBot *tgbotapi.BotAPI, tx *sql.Tx, c *gin.Context) {
	log.Info("GoGinBot is starting...")
	var telegramUser string
	var messageChatId int64
	u := tgbotapi.NewUpdate(0)
	updates, err := startBot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		telegramUser = update.Message.From.UserName
		if len(telegramUser) == 0 {
			continue
		}
		messageChatId = update.Message.Chat.ID
		username, telegramChatId := PrepareChatId(telegramUser)
		if username == telegramUser {
			if telegramChatId == 0 {
				update, err := db.Connect().Prepare("UPDATE user SET telegram_chat_id = (?) WHERE username = (?)")
				if err != nil {
					log.Error("Failed to update data in the database! ", err)
				}
				_, err = update.Exec(messageChatId, username)
				if err != nil {
					log.Error("Failed to update data in the database! ", err)
				}
				msg := tgbotapi.NewMessage(messageChatId, "Successfully subscribed on updates")
				_, err = startBot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", telegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
				}

			} else {
				msg := tgbotapi.NewMessage(messageChatId, "You are registered")
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

func SendMessage(notification string, usertype string) {
	token := os.Getenv("TOKEN")
	rows, err := db.Connect().Query("Select telegram_chat_id FROM user WHERE user_type = (?) AND telegram_chat_id IS NOT NULL", usertype)
	if err != nil {
		log.Error("Failed to select certain data in the database", err)
	}

	for rows.Next() {
		var messageChatId int64
		var url string
		err := rows.Scan(&messageChatId)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		url = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", token, messageChatId, notification)
		_, err = http.Get(url)
		if err == nil {
			log.Infof("Message successfully delivered to ")
		} else if err != nil {
			log.Errorf("Message delivery failed to user %s with error: %s") //!!!!!!!!!!
		}
	}
}

func PrepareChatId(telegramUser string) (username string, telegramChatId int) {
	var data Data
	rows, err := db.Connect().Query("Select username, telegram_chat_id FROM user WHERE username = (?)", telegramUser)
	if err != nil {
		log.Error("Failed to select certain data in the database! ", err)
	}
	for rows.Next() {
		err := rows.Scan(&data.Username, &data.TelegramChatId)
		if err != nil {
			log.Error("The structures does not match! ", err)
		}
		fmt.Println(data.Username, data.TelegramChatId)
	}
	return data.Username, data.TelegramChatId
}

func EscapeMessage(notification string) (newNotification string) {
	var re = regexp.MustCompile(`[[:punct:]]`)
	newNotification = re.ReplaceAllString(notification, "")
	return newNotification
}
