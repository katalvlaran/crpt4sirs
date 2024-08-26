package hmm

import (
	"fmt"
	"math/rand"
)

// HiddenMarkovModel represents an HMM
type HiddenMarkovModel struct {
	States       []string
	Observations []float64
	Transition   [][]float64
	Emission     [][]float64
	Initial      []float64
}

// Initialize HMM with random probabilities
func (hmm *HiddenMarkovModel) Initialize(numStates int, numObservations int) {
	hmm.States = []string{"Bullish", "Bearish"}
	hmm.Transition = make([][]float64, numStates)
	hmm.Emission = make([][]float64, numStates)
	hmm.Initial = make([]float64, numStates)

	// Initialize transition probabilities randomly
	for i := range hmm.Transition {
		hmm.Transition[i] = make([]float64, numStates)
		sum := 0.0
		for j := range hmm.Transition[i] {
			hmm.Transition[i][j] = rand.Float64()
			sum += hmm.Transition[i][j]
		}
		// Normalize to sum to 1
		for j := range hmm.Transition[i] {
			hmm.Transition[i][j] /= sum
		}
	}

	// Initialize emission probabilities randomly
	for i := range hmm.Emission {
		hmm.Emission[i] = make([]float64, numObservations)
		sum := 0.0
		for j := range hmm.Emission[i] {
			hmm.Emission[i][j] = rand.Float64()
			sum += hmm.Emission[i][j]
		}
		// Normalize to sum to 1
		for j := range hmm.Emission[i] {
			hmm.Emission[i][j] /= sum
		}
	}

	// Initialize initial state probabilities
	sum := 0.0
	for i := range hmm.Initial {
		hmm.Initial[i] = rand.Float64()
		sum += hmm.Initial[i]
	}
	// Normalize to sum to 1
	for i := range hmm.Initial {
		hmm.Initial[i] /= sum
	}
}

// Forward algorithm: computes the probability of the observation sequence given the HMM
func (hmm *HiddenMarkovModel) Forward(observations []float64) float64 {
	T := len(observations)
	N := len(hmm.States)

	// Initialize the forward probabilities
	forward := make([][]float64, T)
	for i := range forward {
		forward[i] = make([]float64, N)
	}

	// Initial step
	for i := 0; i < N; i++ {
		forward[0][i] = hmm.Initial[i] * hmm.Emission[i][0] // Assuming the first observation is the first in the list
	}

	// Recursive step
	for t := 1; t < T; t++ {
		for j := 0; j < N; j++ {
			for i := 0; i < N; i++ {
				forward[t][j] += forward[t-1][i] * hmm.Transition[i][j] * hmm.Emission[j][0]
			}
		}
	}

	// Termination step
	prob := 0.0
	for i := 0; i < N; i++ {
		prob += forward[T-1][i]
	}
	return prob
}

func Analysis() {
	// Define observations (candlestick closing prices)
	observations := []float64{101.0, 100.5, 101.5, 102.0, 101.8, 102.3, 102.5, 103.0, 103.5, 102.8}

	// Initialize the HMM
	hmm := &HiddenMarkovModel{}
	hmm.Initialize(2, 1) // 2 states, 1 observation feature (e.g., closing price)

	// Compute the probability of the observation sequence
	prob := hmm.Forward(observations)

	// Output the result
	fmt.Printf("Probability of the observation sequence: %f\n", prob)
}

//func ConvertForAnal()  {}
