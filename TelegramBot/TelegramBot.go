package TelegramBot

import (
	Notification "library/NotificationService"
	"library/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username       string
	TelegramChatId int
}

func Start(startBot *tgbotapi.BotAPI) {
	log.Info("GoGinBot is starting...")
	var data Data
	client := "Client"
	DB := db.Connect()

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

		rows, err := DB.Query("Select username, telegram_chat_id FROM user WHERE username = (?)", telegramUser)
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
				update, err := DB.Prepare("UPDATE test2.user SET telegram_chat_id = (?) WHERE username = (?)")
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
				Notification.SendMessage(messageChatId, telegramUser, botSelfId)
			} else {
				msg := tgbotapi.NewMessage(messageChatId, "You are registered")
				_, err = startBot.Send(msg)
				if err == nil {
					log.Infof("Message successfully delivered to %s", telegramUser)
				} else if err != nil {
					log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
				}
			}

		} else if data.Username != telegramUser {

			_, err := DB.Query("INSERT INTO user(username, user_type, telegram_chat_id) VALUES (?, ?, ?)", telegramUser, client, botSelfId)
			if err != nil {
				log.Error("Failed to insert data in the database! ", err)
			}
			msg := tgbotapi.NewMessage(messageChatId, "Successfully subscribed on updates")
			_, err = startBot.Send(msg)
			if err == nil {
				log.Infof("Message successfully delivered to %s", telegramUser)
			} else if err != nil {
				log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
			}
			Notification.SendMessage(messageChatId, telegramUser, botSelfId)
		}
	}
}
