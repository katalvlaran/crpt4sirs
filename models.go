package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	//_ "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"strconv"
	"time"
)

type Signal struct {
	ID           int       `json:"id,omitempty"`
	PositionSide string    `json:"position_side"`
	Symbol       string    `json:"symbol"`
	Chat         string    `json:"chat"`
	Time         time.Time `json:"time"`
}

type Order struct {
	ID       int       `json:"id,omitempty"`
	Symbol   string    `json:"symbol"`
	Side     string    `json:"side"`  // BUY/SELL
	OType    string    `json:"otype"` // market / take profit / add to position
	Price    string    `json:"price"`
	Quantity string    `json:"quantity"`
	SignalId int       `json:"signal_id"`
	Time     time.Time `json:"time"`
}

var amount = 5.00        // $$ на ордер
var quantityOfOrders = 4 // количество ордеров Take Profit
var futuresClient = createBinanceClient(APIKey, SecretKey)

//futuresClient := binance.NewFuturesClient(APIKey, SecretKey)

// сбор ордера
func NewMarketOrder(sign Signal) {

	ord := new(Order)
	ord.Symbol = sign.Symbol

	price1 := CheckSymbolPrice(sign.Symbol, sign)

	price2 := PriceNormalizer(price1)
	price, err := strconv.ParseFloat(price2, 64)
	if err != nil {
		fmt.Println(err)
	}
	quantity := SetQuantity(price)

	ord.Quantity = quantity

	if sign.PositionSide == "LONG" {
		ord.Side = "BUY"
	} else if sign.PositionSide == "SHORT" {
		ord.Side = "SELL"
	} else {
		fmt.Println()
	}

	order, err := futuresClient.NewCreateOrderService().Symbol(ord.Symbol).
		Side(futures.SideType(ord.Side)).Type(futures.OrderType(binance.OrderTypeMarket)).
		Quantity(ord.Quantity).PriceProtect(true).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	ord.OType = "market"

	signalId, err := getMaxSignalId(db)
	if err != nil {
		//return fmt.Errorf("failed to get max signal ID: %w", err)
		return
	}
	ord.SignalId = signalId

	ord.Price = CheckOpenPositionPrice(sign)

	_, err = saveToTable(db, *ord, "orders")
	if err != nil {
		log.Printf("Error saving to database: %v", err)
	}
	fmt.Println(order, "market order выставлен")
}

// проверка стоимости монеты
func CheckSymbolPrice(str string, sign Signal) float64 {

	prices, err := futuresClient.NewListPricesService().Symbol(sign.Symbol).Do(context.Background())
	if err != nil {
		fmt.Println(prices, err)
	}
	for _, p := range prices {
		fmt.Println(p)
	}
	almostPrice := prices[0]
	price, err := strconv.ParseFloat(almostPrice.Price, 64)
	if err != nil {
		fmt.Println(err)
	}
	return price
}

// установка количества покупаемых монет в ордере
func SetQuantity(price float64) string {

	leverage := 20.00
	quant := amount / price * leverage
	fmt.Println("Количество в ордере", quant)
	price1 := PriceNormalizer(price)
	quantity := QuantityNormalizer(quant, price1)
	return quantity
}

// решение проблемы с количеством чисел после точки для количества монет
func QuantityNormalizer(quant float64, price1 string) string {
	price, err := strconv.ParseFloat(price1, 64)
	if err != nil {
		fmt.Println(err)
	}
	if price <= 1 {
		quantity := strconv.FormatFloat(quant, 'f', 0, 64)
		return quantity
	} else if price <= 15 && price > 1 {
		quantity := strconv.FormatFloat(quant, 'f', 1, 64)
		return quantity
	} else if price <= 100 && price > 15 {
		quantity := strconv.FormatFloat(quant, 'f', 2, 64)
		return quantity
	} else if price <= 1000 && price > 100 {
		quantity := strconv.FormatFloat(quant, 'f', 3, 64)
		return quantity
	} else if price <= 10000 && price > 1000 {
		quantity := strconv.FormatFloat(quant, 'f', 3, 64)
		return quantity
	} else if price <= 100000 && price > 10000 {
		quantity := strconv.FormatFloat(quant, 'f', 3, 64)
		if quantity == "0.000" {
			q := strconv.FormatFloat(quant, 'f', 4, 64)
			return q
		}
		return quantity
	} else {
		quantity := strconv.FormatFloat(quant, 'f', 5, 64)
		return quantity
	}
}

func PriceNormalizer(ordPrice float64) string {
	//price, err := strconv.ParseFloat(ordPrice, 64)
	//if err != nil {
	//	fmt.Println(err)
	//}
	price := ordPrice

	if price <= 0.02 {
		price2 := strconv.FormatFloat(price, 'f', 5, 64)
		return price2
	} else if price <= 0.1 && price > 0.02 {
		price2 := strconv.FormatFloat(price, 'f', 4, 64)
		return price2
	} else if price <= 0.5 && price > 0.1 {
		price2 := strconv.FormatFloat(price, 'f', 4, 64)
		return price2
	} else if price <= 1 && price > 0.5 {
		price2 := strconv.FormatFloat(price, 'f', 3, 64)
		return price2
	} else if price <= 15 && price > 1 {
		price2 := strconv.FormatFloat(price, 'f', 3, 64)
		return price2
	} else if price <= 100 && price > 15 {
		price2 := strconv.FormatFloat(price, 'f', 2, 64)
		return price2
	} else if price <= 1000 && price > 100 {
		price2 := strconv.FormatFloat(price, 'f', 2, 64)
		return price2
	} else if price <= 10000 && price > 1000 {
		price2 := strconv.FormatFloat(price, 'f', 2, 64)
		return price2
	} else {
		price2 := strconv.FormatFloat(price, 'f', 1, 64)
		return price2
	}
}

// проверяем цену входа открытой сделки
func CheckOpenPositionPrice(sign Signal) string {

	accountInfo, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	for _, position := range accountInfo.Positions {
		if position.Symbol == sign.Symbol {
			//fmt.Println(position.EntryPrice)
			return position.EntryPrice //+ "and" + position.PositionAmt
		}
	}
	return ""
}

// проверяем количество монет открытой сделки
func CheckOpenPositionQuantity(sign Signal) string {

	accountInfo, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	for _, position := range accountInfo.Positions {
		if position.Symbol == sign.Symbol {
			return position.PositionAmt
		}
	}
	return ""
}

// дополнительные ордера, чтоб улучшить точку входа
func AddOrders(sign Signal) {
	opPosPrice := CheckOpenPositionPrice(sign)
	oPP, err := strconv.ParseFloat(opPosPrice, 64)
	if err != nil {
		fmt.Println(err)
		//return
	}
	opPosQuantity := CheckOpenPositionQuantity(sign)
	for i := 0; i < 1; i++ {
		ord := new(Order)
		ord.Symbol = sign.Symbol
		// выбираем сторону ордера (обратную от позиции) и формируем цену ордера
		if sign.PositionSide == "LONG" {
			ord.Side = "BUY"
			ordPrice := oPP - (float64(i+1) * oPP / 100) // +- 20%
			ord.Price = PriceNormalizer(ordPrice)
		} else if sign.PositionSide == "SHORT" {
			ord.Side = "SELL"
			ordPrice := oPP + (float64(i+1) * oPP / 100) // decrease by i% of oPP
			ord.Price = PriceNormalizer(ordPrice)
		}
		quant1, err := strconv.ParseFloat(opPosQuantity, 64)
		if quant1 < 0 {
			quant1 = quant1 * (-1)
		}
		if err != nil {
			fmt.Println(err)
			continue // if there's an error, skip this iteration
		}

		quantity := quant1 / 2
		ord.Quantity = QuantityNormalizer(quantity, ord.Price)
		ord.OType = "additional orders"

		signalId, err := getMaxSignalId(db)
		if err != nil {
			//return fmt.Errorf("failed to get max signal ID: %w", err)
			return
		}
		ord.SignalId = signalId

		fmt.Println(ord)

		//if ord.Quantity == "0.00" {
		//	continue
		//}

		tPOrder, err := futuresClient.NewCreateOrderService().Symbol(ord.Symbol).Side(futures.SideType(ord.Side)).
			Type(futures.OrderTypeLimit).Quantity(ord.Quantity).Price(ord.Price).TimeInForce("GTC").
			PriceProtect(true).
			Do(context.Background())

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(tPOrder, "additional orders выставлен")
		}
		_, err = saveToTable(db, *ord, "orders")
		if err != nil {
			log.Printf("Error saving to database: %v", err)
		}
	}

}

// частичная разгрузка по мере получения профита
func TakeProfitOrders(sign Signal) {

	order := new(Order)
	// get the initial opPosPrice
	opPosPrice := CheckOpenPositionPrice(sign)
	oPP, err := strconv.ParseFloat(opPosPrice, 64)
	if err != nil {
		fmt.Println(err)
		//return
	}

	opPosQuantity := CheckOpenPositionQuantity(sign)
	for i := 0; i < quantityOfOrders; i++ {
		ord := new(Order)
		ord.Symbol = sign.Symbol

		// выбираем сторону ордера (обратную от позиции) и формируем цену ордера
		if sign.PositionSide == "LONG" {
			ord.Side = "SELL"
			ordPrice := oPP + (float64(i+1) * oPP / 100) // increase by i% of oPP
			ord.Price = PriceNormalizer(ordPrice)
		} else if sign.PositionSide == "SHORT" {
			ord.Side = "BUY"
			ordPrice := oPP - (float64(i+1) * oPP / 100) // decrease by i% of oPP
			ord.Price = PriceNormalizer(ordPrice)
			//ord.Price = strconv.FormatFloat(ordPrice, 'f', 2, 64)
		}

		// формируем количество в ордере
		quant1, err := strconv.ParseFloat(opPosQuantity, 64)
		if quant1 < 0 {
			quant1 = quant1 * (-1)
		}
		if err != nil {
			fmt.Println(err)
			continue // if there's an error, skip this iteration
		}

		qOOFl := float64(quantityOfOrders)

		quantity := quant1 / qOOFl
		ord.Quantity = QuantityNormalizer(quantity, ord.Price)
		ord.OType = "take profit"

		signalId, err := getMaxSignalId(db)
		if err != nil {
			//return fmt.Errorf("failed to get max signal ID: %w", err)
			return
		}

		ord.SignalId = signalId

		fmt.Println(ord)

		tPOrder, err := futuresClient.NewCreateOrderService().Symbol(ord.Symbol).Side(futures.SideType(ord.Side)).
			Type(futures.OrderTypeLimit).Quantity(ord.Quantity).Price(ord.Price).TimeInForce("GTC").
			PriceProtect(true).
			Do(context.Background())

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(tPOrder, "take profit order выставлен")
		}
		_, err = saveToTable(db, *ord, "orders")
		if err != nil {
			log.Printf("Error saving to database: %v", err)
		}
		fmt.Println(order)
	}
}

func StopLossOrder(sign Signal) {

	// get the initial opPosPrice
	opPosPrice := CheckOpenPositionPrice(sign)
	oPP, err := strconv.ParseFloat(opPosPrice, 64)
	if err != nil {
		fmt.Println(err)
		//return
	}

	//opPosQuantity := CheckOpenPositionQuantity(sign)

	ord := new(Order)
	ord.Symbol = sign.Symbol

	// выбираем сторону ордера (обратную от позиции) и формируем цену ордера
	if sign.PositionSide == "LONG" {
		ord.Side = "SELL"
		ordPrice := oPP - (oPP / 50)
		ord.Price = PriceNormalizer(ordPrice)
	} else if sign.PositionSide == "SHORT" {
		ord.Side = "BUY"
		ordPrice := oPP + (oPP / 50)
		ord.Price = PriceNormalizer(ordPrice)
	}

	ord.OType = "stop loss"

	signalId, err := getMaxSignalId(db)
	if err != nil {
		//return fmt.Errorf("failed to get max signal ID: %w", err)
		return
	}

	ord.SignalId = signalId

	fmt.Println(ord)

	tPOrder, err := futuresClient.NewCreateOrderService().Symbol(ord.Symbol).Side(futures.SideType(ord.Side)).
		Type(futures.OrderTypeStopMarket).StopPrice(ord.Price).TimeInForce("GTC").
		PriceProtect(true).ClosePosition(true).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tPOrder, "stop loss order выставлен")
	}
	_, err = saveToTable(db, *ord, "orders")
	if err != nil {
		log.Printf("Error saving to database: %v", err)
	}

}

//func PriceChecker(sign Signal) {
//	price := CheckSymbolPrice
//	opPrice := CheckOpenPositionPrice(sign)
//	if sign.PositionSide == "LONG" {
//		if price > opPrice {
//
//		}
//	}
//
//}
