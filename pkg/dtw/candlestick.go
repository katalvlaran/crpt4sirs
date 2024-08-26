package dtw

// Candlestick structure representing OHLC data
type Candlestick struct {
	Open, High, Low, Close float64
}

// FetchCandlestickData simulates fetching candlestick data, which in a real application would come from an API like Binance or CoinGecko.
func FetchCandlestickData() []Candlestick {
	// Simulating fetched candlestick data
	return []Candlestick{
		{Open: 30000, High: 31000, Low: 29000, Close: 30500},
		{Open: 30600, High: 31500, Low: 30000, Close: 31000},
		{Open: 31100, High: 32000, Low: 30500, Close: 31500},
	}
}

// ExtractClosePrices extracts the close prices from a slice of candlesticks
func ExtractClosePrices(candles []Candlestick) []float64 {
	closePrices := make([]float64, len(candles))
	for i, candle := range candles {
		closePrices[i] = candle.Close
	}
	return closePrices
}
