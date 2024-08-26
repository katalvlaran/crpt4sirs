package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Klines struct {
	ID                       int       `json:"id,omitempty"`
	openTime                 time.Time `json:"openTime"`
	closeTime                time.Time `json:"closeTime"`
	openPrice                float64   `json:"openPrice"`
	highPrice                float64   `json:"highPrice"`
	lowPrice                 float64   `json:"lowPrice"`
	closePrice               float64   `json:"closePrice"`
	volume                   float64   `json:"volume"`
	quoteAssetVolume         float64   `json:"quoteAssetVolume"`
	numberOfTrades           int64     `json:"numberOfTrades"`
	takerBuyBaseAssetVolume  float64   `json:"takerBuyBaseAssetVolume"`
	takerBuyQuoteAssetVolume float64   `json:"takerBuyQuoteAssetVolume"`
}

func ParseKlines() {

	timeString := "2023-07-01 00:00:00 +0200 CEST"
	timeFormat := "2006-01-02 15:04:05 -0700 MST"

	t, err := time.Parse(timeFormat, timeString)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	// Количество циклов (1440 нужно всего - 2 года по 12 часов)
	numCycles := 2

	// Продолжительность для увеличения
	duration := 12 * time.Hour

	for i := 0; i < numCycles; i++ {
		// Уdtkbxение времени на 12 часов
		t = t.Add(duration)
		time2 := t.UnixNano() / int64(time.Millisecond)
		//result := GetKlinesInfo(time2)
		GetKlinesInfo(time2, db)
	}
}

func GetKlinesInfo(timeOfKline int64, db *sql.DB) {

	klines, err := BinClient().NewKlinesService().
		Symbol("ETHUSDT").
		Interval("1m").
		StartTime(timeOfKline).
		Limit(720). // 720 минут в 12 часах
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Извлечение и вывод данных из свечей
	for _, k := range klines {
		openTime := time.Unix(0, k.OpenTime*int64(time.Millisecond))
		closeTime := time.Unix(0, k.CloseTime*int64(time.Millisecond))
		openPrice, _ := strconv.ParseFloat(k.Open, 64)
		highPrice, _ := strconv.ParseFloat(k.High, 64)
		lowPrice, _ := strconv.ParseFloat(k.Low, 64)
		closePrice, _ := strconv.ParseFloat(k.Close, 64)
		volume, _ := strconv.ParseFloat(k.Volume, 64)
		quoteAssetVolume, _ := strconv.ParseFloat(k.QuoteAssetVolume, 64)
		numberOfTrades := k.TradeNum
		takerBuyBaseAssetVolume, _ := strconv.ParseFloat(k.TakerBuyBaseAssetVolume, 64)
		takerBuyQuoteAssetVolume, _ := strconv.ParseFloat(k.TakerBuyQuoteAssetVolume, 64)

		fmt.Printf("Open Time: %s, Open Price: %f, High Price: %f, Low Price: %f, Close Price: %f, Volume: %f, Close Time: %s, Quote Asset Volume: %f, Number of Trades: %d, Taker Buy Base Asset Volume: %f, Taker Buy Quote Asset Volume: %f\n",
			openTime, openPrice, highPrice, lowPrice, closePrice, volume, closeTime, quoteAssetVolume, numberOfTrades, takerBuyBaseAssetVolume, takerBuyQuoteAssetVolume)
		kln := new(Klines)
		kln.openTime = openTime
		kln.closeTime = closeTime
		kln.openPrice = openPrice
		kln.highPrice = highPrice
		kln.lowPrice = lowPrice
		kln.closePrice = closePrice
		kln.volume = volume
		kln.quoteAssetVolume = quoteAssetVolume
		kln.numberOfTrades = numberOfTrades
		kln.takerBuyBaseAssetVolume = takerBuyBaseAssetVolume
		kln.takerBuyQuoteAssetVolume = takerBuyQuoteAssetVolume
		_, err = saveToTable(db, *kln, "testklines")
		if err != nil {
			log.Printf("Error saving to database: %v", err)
		}
	}
}
