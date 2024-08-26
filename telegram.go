package main

import (
	_ "fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

func processUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) string {
	//if update.Message == nil {
	//	//return nil
	//}

	text := update.Message.Text
	chatID := update.Message.Chat.ID
	// kinda validation
	if chatID != 5801485139 && chatID != 470937085 {
		log.Println("Invalid sender:", chatID)
		//return
	}
	// response in telegram bot
	reply := "Message received: " + text
	msg := tgbotapi.NewMessage(chatID, reply)

	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
		//return
	}
	return text
}

func parseSignal(text string) *Signal {
	dataSlice := strings.Split(text, ",")
	if len(dataSlice) < 3 {
		return nil
	}

	s := &Signal{
		PositionSide: dataSlice[0],
		Symbol:       dataSlice[1] + "USDT",
		Chat:         dataSlice[2],
		Time:         time.Now(),
	}

	_, err := saveToTable(db, *s, "signal")
	if err != nil {
		log.Printf("Error saving to database: %v", err)
	}

	return s

}
