package bot

import "fmt"

// TelegramBot represents a Telegram bot.
type TelegramBot struct {
	Token string
}

// NewTelegramBot creates a new instance of TelegramBot.
func NewTelegramBot(token string) *TelegramBot {
	return &TelegramBot{Token: token}
}

// SendMessage sends a message to a Telegram user.
func (b *TelegramBot) SendMessage(chatID, message string) {
	// Implement the logic to send a message using the Telegram Bot API
	fmt.Printf("Sending message to %s: %s\n", chatID, message)
}
