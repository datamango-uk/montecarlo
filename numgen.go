package montecarlo

import (
	"math/rand/v2"
)

// Normal returns a random number which falls within the normal distribution as defined by the mean and standard deviation
func NormalFloat64(mean, stdDeviation float64) float64 {
	return rand.NormFloat64()*stdDeviation + mean
}

// Uniform returns a random number which is within the specified limits
func UniformFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
