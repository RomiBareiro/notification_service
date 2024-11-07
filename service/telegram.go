package service

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getToken() string {
	return os.Getenv("TELEGRAM_TOKEN")
}
func sendTelegramMessage(chatID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		return fmt.Errorf("could not initialize Telegram bot: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)

	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	log.Printf("Message sent to chat ID: %d", chatID)
	return nil
}
