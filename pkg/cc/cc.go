package cross_correlation

import (
	"fmt"
	"math"
)

// CrossCorrelation computes the cross-correlation of two time series
func CrossCorrelation(pattern, data []float64) []float64 {
	pLen := len(pattern)
	dLen := len(data)
	cc := make([]float64, dLen-pLen+1)

	// Calculate the mean of the pattern
	var meanPattern float64
	for _, v := range pattern {
		meanPattern += v
	}
	meanPattern /= float64(pLen)

	for i := 0; i <= dLen-pLen; i++ {
		var meanData float64
		for j := 0; j < pLen; j++ {
			meanData += data[i+j]
		}
		meanData /= float64(pLen)

		var numerator, denomPattern, denomData float64
		for j := 0; j < pLen; j++ {
			numerator += (pattern[j] - meanPattern) * (data[i+j] - meanData)
			denomPattern += math.Pow(pattern[j]-meanPattern, 2)
			denomData += math.Pow(data[i+j]-meanData, 2)
		}

		cc[i] = numerator / (math.Sqrt(denomPattern) * math.Sqrt(denomData))
	}
	return cc
}

func Analysis() {
	pattern := []float64{102.5, 103.0, 101.5, 102.0}
	data := []float64{100.0, 101.0, 102.0, 102.5, 103.0, 101.5, 102.0, 102.3, 103.5, 102.8}

	crossCorrelation := CrossCorrelation(pattern, data)
	fmt.Println("Cross-correlation values:", crossCorrelation)
}

//func ConvertForAnal()  {}
