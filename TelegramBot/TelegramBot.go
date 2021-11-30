package telegrambot

import (
	db "library/ConnectionDatabase"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	log "github.com/sirupsen/logrus"
)

type Data struct {
	Username       string
	TelegramChatID int
}

const ESCAPE_RUNE = '\\'

var (
	RUNE_TO_ESCAPE = []rune{'[', ']', '(', ')', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!', '_'}
	Bot            *tgbotapi.BotAPI
)

func Listen() {
	log.Info("GoGinBot is starting...")

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60
	updates, err := Bot.GetUpdatesChan(upd)
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

		username, telegramChatID := FoundChatID(telegramUser)
		if username != telegramUser {
			continue
		}

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

			_, err = Bot.Send(msg)
			if err == nil {
				log.Infof("Message successfully delivered to %s", telegramUser)
			} else if err != nil {
				log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
			}
		} else {
			msg := tgbotapi.NewMessage(messageChatID, "You are registered")
			_, err = Bot.Send(msg)
			if err == nil {
				log.Infof("Message successfully delivered to %s", telegramUser)
			} else if err != nil {
				log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
			}
		}
	}
}

func InitBot() *tgbotapi.BotAPI {

	token := os.Getenv("TOKEN")
	debug := os.Getenv("DEBUG")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug, _ = strconv.ParseBool(debug)

	return bot
}

func SendMessage(notification string, telegramUser string, messageChatID int64) {
	msg := tgbotapi.NewMessage(messageChatID, notification)
	_, err := Bot.Send(msg)

	if err == nil {
		log.Infof("Message successfully delivered to %s", telegramUser)
	} else if err != nil {
		log.Errorf("Message delivery failed to user %s with error: %s", telegramUser, err)
	}
}

func FoundChatID(telegramUser string) (username string, telegramChatID int) {
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
	var builder strings.Builder

	for _, specialWord := range notification {
		for _, specialWord2 := range RUNE_TO_ESCAPE {
			if strings.Contains(string(specialWord), string(specialWord2)) {
				builder.WriteRune(ESCAPE_RUNE)
			}
		}
		builder.WriteRune(specialWord)
	}

	return builder.String()
}
