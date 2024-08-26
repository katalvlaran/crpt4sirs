/* # Sliding Window Method for Pattern Matching. The Sliding Window method is a popular technique for time series
analysis and pattern recognition. It involves moving a fixed-size window across the data series to extract subsequences,
which are then compared with the target pattern. This method is particularly useful for finding patterns in real-time
data, such as identifying specific candlestick patterns in financial charts.*/
/*
- Sliding Window: A fixed-size window moves across the time series data, capturing a segment (or subsequence)
	of the data at each step.
- Pattern Matching: Each captured subsequence is compared to a predefined pattern using a similarity measure
	(e.g., Euclidean distance, correlation).
- Thresholding: If the similarity measure between the windowed subsequence and the target pattern exceeds
	a certain threshold, a match is identified.
The key parameters are:
1. Window Size: Determines the length of the subsequences.
2. Step Size: Determines how far the window moves at each step (usually one data point).
3. Similarity Measure: Determines how the similarity between the subsequence and the pattern is calculated.*/
/*Tuning and Enhancing the Sliding Window Method

# 1. Adjusting the Window Size and Threshold
- Window Size: The length of the pattern determines the window size. If the pattern length varies, you may need
	to implement a variable window size or use a multi-resolution approach.
- Threshold: The threshold determines the sensitivity of the pattern matching. Lowering the threshold will result
	in fewer, more precise matches, while raising it will allow for more approximate
	matches.
# 2. Implementing Advanced Similarity Measures
- Normalized Cross-Correlation: For data where amplitude variations are common, you might use normalized
	cross-correlation instead of Euclidean distance to focus on the shape of the pattern rather than its absolute values.
# 3. Multi-resolution Sliding Window
- Multi-resolution Approach: For complex patterns, apply the sliding window method at different resolutions
	(e.g., varying the window size) to capture patterns at different scales.*/

package sw

import (
	"fmt"
	"math"
)

// Example candlestick closing prices
var data = []float64{100.0, 101.5, 102.0, 101.8, 102.3, 103.5, 102.8, 104.0, 103.5, 103.0, 102.0}

// Pattern to match (e.g., a bullish pattern)
var pattern = []float64{101.8, 102.3, 103.5}

// EuclideanDistance calculates the Euclidean distance between two slices of data
func EuclideanDistance(slice1, slice2 []float64) float64 {
	if len(slice1) != len(slice2) {
		panic("Slices must have the same length")
	}
	sum := 0.0
	for i := range slice1 {
		sum += math.Pow(slice1[i]-slice2[i], 2)
	}
	return math.Sqrt(sum)
}

// SlidingWindowSearch performs a sliding window search on the data to find a matching pattern
func SlidingWindowSearch(data, pattern []float64, threshold float64) int {
	windowSize := len(pattern)
	bestMatchIndex := -1
	bestMatchDistance := math.MaxFloat64

	for i := 0; i <= len(data)-windowSize; i++ {
		window := data[i : i+windowSize]
		distance := EuclideanDistance(window, pattern)
		fmt.Printf("Window [%d:%d] -> Distance: %f\n", i, i+windowSize, distance)
		if distance < threshold && distance < bestMatchDistance {
			bestMatchIndex = i
			bestMatchDistance = distance
		}
	}
	return bestMatchIndex
}

func NormalizedCrossCorrelation(slice1, slice2 []float64) float64 {
	if len(slice1) != len(slice2) {
		panic("Slices must have the same length")
	}
	mean1, mean2 := mean(slice1), mean(slice2)
	numerator, denom1, denom2 := 0.0, 0.0, 0.0
	for i := range slice1 {
		diff1 := slice1[i] - mean1
		diff2 := slice2[i] - mean2
		numerator += diff1 * diff2
		denom1 += diff1 * diff1
		denom2 += diff2 * diff2
	}
	return numerator / math.Sqrt(denom1*denom2)
}

func mean(slice []float64) float64 {
	sum := 0.0
	for _, v := range slice {
		sum += v
	}
	return sum / float64(len(slice))
}

func Analysis() {
	// Define the threshold for a match
	threshold := 0.5

	// Perform the sliding window search
	matchIndex := SlidingWindowSearch(data, pattern, threshold)

	if matchIndex != -1 {
		fmt.Printf("Pattern matched at index: %d\n", matchIndex)
	} else {
		fmt.Println("Pattern not found.")
	}
}

//func ConvertForAnal()  {}
