package TelegramBot

import (
	"fmt"
	Error "library/JsonError"
	Notification "library/NotificationService"
	"library/db"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username         string
	Telegram_chat_id int
}

func TelegramBot() (c *gin.Context) {
	fmt.Println("GoGinBot is starting...")
	var data Data

	client := "Client"
	token := os.Getenv("TOKEN")

	DB := db.InitDB()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		TelegramUser := update.Message.From.UserName
		chat_ID := update.Message.Chat.ID
		bot_self_id := bot.Self.ID

		rows, err := DB.Query("Select username, telegram_chat_id FROM user WHERE username = (?)", TelegramUser)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)
		}

		for rows.Next() {
			err := rows.Scan(&data.Username, &data.Telegram_chat_id)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
		}

		if data.Username == TelegramUser {
			if data.Telegram_chat_id == 0 {
				update, err := DB.Prepare("UPDATE test2.user SET telegram_chat_id = (?) WHERE username = (?)")
				if Error.Error(c, err) {
					log.Error("Failed to update data in the database! ", err)
					return
				}

				_, err = update.Exec(bot_self_id, data.Username)
				if Error.Error(c, err) {
					log.Error("Failed to execute data in the database! ", err)
					return
				}
				msg := tgbotapi.NewMessage(chat_ID, "Successfully subscribed on updates")
				_, err = bot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", TelegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", TelegramUser, err)
				}
				Notification.Notification(chat_ID, TelegramUser, bot_self_id)
			} else {
				msg := tgbotapi.NewMessage(chat_ID, "You are registered")
				_, err = bot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", TelegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", TelegramUser, err)
				}
			}

		} else if data.Username != TelegramUser {

			_, err := DB.Query("INSERT INTO user(username, user_type, telegram_chat_id) VALUES (?, ?, ?)", TelegramUser, client, bot_self_id)
			if err != nil {
				log.Error("Failed to insert data in the database! ", err)
			}
			msg := tgbotapi.NewMessage(chat_ID, "Successfully subscribed on updates")
			_, err = bot.Send(msg)
			if err == nil {
				log.Infof("Message successfully delivered to %s", TelegramUser)
			} else if err != nil {
				log.Errorf("Message delivery failed to user %s with error: %s", TelegramUser, err)
			}
			Notification.Notification(chat_ID, TelegramUser, bot_self_id)
		}
	}
	return
}
