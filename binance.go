package main

import (
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func createBinanceClient(apiKey, secretKey string) *futures.Client {
	client := binance.NewFuturesClient(apiKey, secretKey)
	client.BaseURL = TestnetBaseURL
	return client
}
