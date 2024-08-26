package main

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"time"

	_ "github.com/adshao/go-binance/v2"
	_ "github.com/adshao/go-binance/v2/futures"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var rdb *redis.Client

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegramToken := os.Getenv("telegramToken")

	// Подключение к Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	var openPositions []string

	initDB()

	//ParseKlines()
	//os.Exit(1)

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Failed to get updates:", err)
	}

	for update := range updates {
		result := processUpdate(update, bot)

		// разбор сообщения с тг
		sign := parseSignal(result)
		if sign == nil {
			log.Println("Error parsing signal")
			return
		}

		lastSymbol := sign.Symbol
		err := rdb.Set(context.Background(), "last_open_time", lastSymbol, 0).Err()
		if err != nil {
			log.Fatal(err)
		}

		openPositions = append(openPositions, lastSymbol)

		// создание нового ордера
		NewMarketOrder(*sign)

		time.Sleep(1 * time.Second)

		AddOrders(*sign)

		time.Sleep(1 * time.Second)

		TakeProfitOrders(*sign)

		time.Sleep(1 * time.Second)

		StopLossOrder(*sign)

		fmt.Println(lastSymbol)

		fmt.Println(openPositions)
	}
}
