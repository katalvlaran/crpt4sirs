package main

import (
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func BinFutCl() *futures.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	APIKey := os.Getenv("APIKey")
	SecretKey := os.Getenv("SecretKey")
	TestnetBaseURL := os.Getenv("TestnetBaseURL")

	client := binance.NewFuturesClient(APIKey, SecretKey)
	client.BaseURL = TestnetBaseURL
	return client
}

func BinClient() *binance.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	APIKey := os.Getenv("APIKey")
	SecretKey := os.Getenv("SecretKey")
	TestnetBaseURL := os.Getenv("TestnetBaseURL")

	client := binance.NewClient(APIKey, SecretKey)
	client.BaseURL = TestnetBaseURL
	return client
}
