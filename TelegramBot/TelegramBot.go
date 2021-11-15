package TelegramBot

import (
	Error "library/JsonError"

	"library/db"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username       string
	TelegramChatId int
}

func Start(startBot *tgbotapi.BotAPI) (c *gin.Context) {
	log.Info("GoGinBot is starting...")
	var data Data

	tx, err := db.Connect().Begin()
	if Error.Error(c, err) {
		log.Error("Failed to connect to database! ", err)
		return
	}

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
		messageChatId := update.Message.Chat.ID
		botSelfId := startBot.Self.ID

		rows, err := tx.Query("Select username, telegram_chat_id FROM user WHERE username = (?)", telegramUser)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)
		}

		for rows.Next() {
			err := rows.Scan(&data.Username, &data.TelegramChatId)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
		}

		if data.Username == telegramUser {
			if data.TelegramChatId == 0 {
				update, err := tx.Prepare("UPDATE user SET telegram_chat_id = (?) WHERE username = (?)")
				if err != nil {
					log.Error("Failed to update data in the database! ", err)
				}
				_, err = update.Exec(botSelfId, data.Username)
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
				continue
			}

			msg := tgbotapi.NewMessage(messageChatId, "You are registered")
			_, err = startBot.Send(msg)
			if err == nil {
				log.Infof("Message successfully delivered to %s", telegramUser)
			} else if err != nil {
				log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
			}
		}
	}

	return
}
