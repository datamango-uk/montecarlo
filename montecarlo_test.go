package montecarlo_test

import (
	"math"
	"testing"

	"github.com/datamango-uk/montecarlo"
)

func TestSimulationValueOfPi(t *testing.T) {
	simulation := montecarlo.New(
		100000, // Iterations
		200,    // Worker pool size
		montecarlo.RunFunc(func(_ map[string]float64) float64 {
			x := montecarlo.UniformFloat64(0, 1)
			y := montecarlo.UniformFloat64(0, 1)
			if x*x+y*y <= 1 {
				return 1
			}
			return 0
		}),
	)
	summary := simulation.Run(nil)

	var insideCircle float64
	for _, result := range summary.Results {
		insideCircle += result
	}
	piEstimate := (insideCircle / float64(simulation.Iterations)) * 4

	const tolerance = 0.01
	if piEstimate < math.Pi-tolerance || piEstimate > math.Pi+tolerance {
		t.Errorf("Estimated Pi value %f is not within the tolerance range of the actual Pi value %f", piEstimate, math.Pi)
	}
}

func TestSimulationPortfolioPrediction(t *testing.T) {
	summary := montecarlo.New(
		100000, // Iterations
		200,    // Worker pool size
		montecarlo.RunFunc(func(input map[string]float64) float64 {
			// This function simulates the daily returns of a portfolio over a given number of days.
			// It uses a normal distribution to generate daily returns based on the mean and standard deviation provided.
			// The portfolio value is updated daily and the final value is returned.
			portfolioValue := input["initialValue"]
			meanReturn := input["meanReturn"]
			stddevReturn := input["stddevReturn"]
			days := int(input["days"])
			for i := 0; i < days; i++ {
				dailyReturn := montecarlo.NormalFloat64(meanReturn, stddevReturn)
				portfolioValue *= (1 + dailyReturn)
			}
			return portfolioValue
		}),
	).Run(map[string]float64{
		"initialValue": 10000.0,
		"meanReturn":   0.0005,
		"stddevReturn": 0.01,
		"days":         365,
	})

	if summary.Stats.Mean <= 0 {
		t.Errorf("Mean of the portfolio value is non-positive: %f", summary.Stats.Mean)
	}

	if summary.Stats.StandardDeviation <= 0 {
		t.Errorf("Standard deviation of the portfolio value is non-positive: %f", summary.Stats.StandardDeviation)
	}

	if summary.Stats.Min <= 0 {
		t.Errorf("Minimum portfolio value is non-positive: %f", summary.Stats.Min)
	}

	if summary.Stats.Max == -math.MaxFloat64 {
		t.Errorf("Maximum portfolio value is not updated and remains at the initial value: %f", summary.Stats.Max)
	}

	if summary.Stats.Min > summary.Stats.Max {
		t.Errorf("Minimum portfolio value %f is greater than maximum portfolio value %f", summary.Stats.Min, summary.Stats.Max)
	}
}
