package TelegramBot

import (
	"fmt"
	"library/db"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username string
}

func TelegramBot() {
	fmt.Println("TelegramBot is starting...")
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
		rows, err := DB.Query("Select username FROM user WHERE username = (?)", update.Message.From.UserName)
		if err != nil {
			log.Error("Failed to select certain data in the database! ", err)
		}

		for rows.Next() {
			err := rows.Scan(&data.Username)
			if err != nil {
				log.Error("The structures does not match! ", err)
			}
		}

		if data.Username == update.Message.From.UserName {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are registered")
			bot.Send(msg)
		} else if data.Username != update.Message.From.UserName {

			_, err := DB.Query("INSERT INTO user(username, user_type) VALUES (?, ?)", update.Message.From.UserName, client)
			if err != nil {
				log.Error("Failed to insert data in the database! ", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You subscribed")
			bot.Send(msg)

		}
	}

}
