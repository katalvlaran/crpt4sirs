//### 1. Understanding DTW: Theory
//Dynamic Time Warping (DTW) is a method for measuring the similarity between two temporal sequences, which may vary in speed. It finds an optimal alignment between the sequences, allowing for sequences that are not perfectly aligned in time. This flexibility makes DTW particularly effective in pattern recognition and time-series analysis, like matching candlestick patterns.
//#### Key Concepts:
//- Warping Path: DTW calculates the minimum cumulative distance (cost) to align two sequences by mapping elements from one sequence to the other.
//- Distance Measure: Typically, the Euclidean distance is used, but DTW can work with other distance measures.
//- Cost Matrix: A matrix that stores the cumulative distance for each element, representing how well one sequence matches another at each point.
//#### Key Concepts:
//- Warping Path: DTW calculates the minimum cumulative distance (cost) to align two sequences by mapping elements from one sequence to the other.
//- Distance Measure: Typically, the Euclidean distance is used, but DTW can work with other distance measures.
//- Cost Matrix: A matrix that stores the cumulative distance for each element, representing how well one sequence matches another at each point.

package dtw

import (
	"fmt"
	"math"
)

// Min function returns the minimum value among three floats
// Finds the minimum of three numbers, which is used to find the optimal path in the DTW algorithm.
func Min(a, b, c float64) float64 {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// DTW Function calculates the Dynamic Time Warping distance between two sequences.
func DTW(seq1, seq2 []float64) float64 {
	n, m := len(seq1), len(seq2)
	dtw := make([][]float64, n+1) // slice dtw stores the cumulative distances.
	for i := range dtw {
		dtw[i] = make([]float64, m+1)
		for j := range dtw[i] {
			dtw[i][j] = math.Inf(1)
		}
	}
	dtw[0][0] = 0

	// The main loop calculates the cost for each point in the sequences, updating the matrix with the minimum cumulative cost.
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := math.Abs(seq1[i-1] - seq2[j-1])
			dtw[i][j] = cost + Min(dtw[i-1][j], dtw[i][j-1], dtw[i-1][j-1])
		}
	}

	return dtw[n][m]
}

func Analysis() {
	//func main() {
	// Fetch real-time candlestick data
	realTimeData := FetchCandlestickData()

	// Define a pattern you want to match against
	pattern := []Candlestick{
		{Open: 30500, High: 31000, Low: 30000, Close: 30800},
		{Open: 30850, High: 31500, Low: 30500, Close: 31200},
		{Open: 31300, High: 32000, Low: 31000, Close: 31800},
	}

	// Convert to close prices
	realTimeClosePrices := ExtractClosePrices(realTimeData)
	patternClosePrices := ExtractClosePrices(pattern)

	// Calculate DTW distance. Measures how close the real-time data is to a known pattern.
	distance := DTW(realTimeClosePrices, patternClosePrices)
	fmt.Printf("DTW Distance with Real-Time Data: %.2f\n", distance)
}

//func ConvertForAnal()  {}
