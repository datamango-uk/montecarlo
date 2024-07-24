# Monte Carlo Simulation Library

This library provides tools for running Monte Carlo simulations. It includes functionality for estimating the value of Pi and predicting portfolio values based on daily returns.

## Installation

To install the library, use the following command:

```sh
go get github.com/datamango-uk/montecarlo
```

## Usage

### Estimating the Value of Pi

The following example demonstrates how to use the Monte Carlo simulation to estimate the value of Pi.

```go
package main

import (
    "github.com/datamango-uk/montecarlo"
)

func main() {
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
    summary := simulation.Run(nil) // This simulation doesnt require any inputs, so we pass nil


    // Using the results to estimate the value of Pi
    var insideCircle float64
    for _, result := range summary.Results {
        insideCircle += result
    }
    piEstimate := (insideCircle / float64(simulation.Iterations)) * 4
    // ~= 3.141592653589793
}
```

### Predicting Portfolio Values

The following example demonstrates how to use the Monte Carlo simulation to predict the value of a portfolio based on daily returns.

```go
// This function simulates the daily returns of a portfolio over a given number of days.
// It uses a normal distribution to generate daily returns based on the mean and standard deviation provided.
// The portfolio value is updated daily and the final value is returned.
func estimatePortfolio(input map[string]float64) float64 {
    portfolioValue := input["initialValue"]
    meanReturn := input["meanReturn"]
    stddevReturn := input["stddevReturn"]
    days := int(input["days"])
    for i := 0; i < days; i++ {
        dailyReturn := montecarlo.NormalFloat64(meanReturn, stddevReturn)
        portfolioValue *= (1 + dailyReturn)
    }
    return portfolioValue
}

func main() {
    summary := montecarlo.New(
        100000, // Iterations
        200,    // Worker pool size
        montecarlo.RunFunc(estimatePortfolio),
    ).Run(map[string]float64{
        "initialValue": 10000.0,
        "meanReturn":   0.0005,
        "stddevReturn": 0.01,
        "days":         365,
    })

    fmt.Println(summary.Stats)
    // {Mean:12036.539859964061 StandardDeviation:2279.4204494720975 Min:6481.218067148653 Max:23742.014047666456}
}
```

### Running multiple simulations

If you want to run a simulation multiple times with different inputs, you can use the `RunMultiple` method.

```go
func main() {
    simulation := montecarlo.New(100000, 200, montecarlo.RunFunc(estimatePortfolio))
    inputs := []map[string]float64{
        {
            "initialValue": 100.0,
            "meanReturn":   0.0005,
            "stddevReturn": 0.01,
            "days":         365,
        },
        {
            "initialValue": 100.0,
            "meanReturn":   0.00010,
            "stddevReturn": 0.02,
            "days":         200,
        },
    }
    simulation.RunMultiple(inputs...)
    // []Summary{
    //    {Results: [...] Stats:{} InputValues: {"initialValue": 10000.0, "meanReturn": 0.0005, "stddevReturn": 0.01, "days": 365}},
    //    {Results: [...] Stats:{} InputValues: {"initialValue": 10000.0, "meanReturn": 0.00010, "stddevReturn": 0.02, "days": 200}},
    // }
}
```

## License

This project is licensed under the MIT License.