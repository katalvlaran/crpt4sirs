package main

import (
	_ "context"
	"fmt"
	"log"
	"time"

	_ "github.com/adshao/go-binance/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	telegramToken = "6430281663:AAGR0sFbhnJzFxDiKc4NC1VYGQ-iHejJXf8"
	// test binance:
	TestnetBaseURL = "https://testnet.binancefuture.com"
	APIKey         = "a365f1b5f1e4bba5f55e8156139b73c841e4a5578aea7c7496e4b3da16bd0957"
	SecretKey      = "31e5e249781567a1dd62f1db80e0bbd5cc5a058b68331a89d5a71058412d6436"
)

func main() {
	// Initialize DB
	initDB()

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	// request ip: 188.163.12.53
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Failed to get updates:", err)
	}

	client := createBinanceClient(APIKey, SecretKey)
	fmt.Println(client)

	for update := range updates {
		result := processUpdate(update, bot)

		// разбор сообщения с тг
		sign := parseSignal(result)
		if sign == nil {
			log.Println("Error parsing signal")
			return
		}

		// создание нового ордера
		NewMarketOrder(*sign)

		time.Sleep(1 * time.Second)

		AddOrders(*sign)

		time.Sleep(1 * time.Second)

		TakeProfitOrders(*sign)

		time.Sleep(1 * time.Second)

		StopLossOrder(*sign)

	}
}
