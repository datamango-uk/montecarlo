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
    simulation := montecarlo.Simulation{
        Iterations: 100000,
        Workers: 200,
        RunFunc: func(_ map[string]float64) float64 {
            x := montecarlo.UniformFloat64(0, 1)
            y := montecarlo.UniformFloat64(0, 1)
            if xx+yy <= 1 {
                return 1
            }
            return 0
        },
    }
    results := simulation.Run()


    // Using the results of the simultion to estimate the value of Pi
    var insideCircle float64
    for _, result := range results {
        insideCircle += result
    }
    piEstimate := (insideCircle / float64(simulation.Iterations)) * 4
    // ~= 3.141592653589793
}
```

### Predicting Portfolio Values

The following example demonstrates how to use the Monte Carlo simulation to predict the value of a portfolio based on daily returns.

```go
package main

import (
    "github.com/datamango-uk/montecarlo"
)

func main() {
    simulation := montecarlo.Simulation{
        Iterations: 1000,
        InputValues: map[string]float64{
            "initialValue": 10000.0,
            "meanReturn": 0.0005,
            "stddevReturn": 0.01,
            "days": 365,
        },
        RunFunc: func(input map[string]float64) float64 {
            portfolioValue := input["initialValue"]
            meanReturn := input["meanReturn"]
            stddevReturn := input["stddevReturn"]
            days := int(input["days"])
            for i := 0; i < days; i++ {
                dailyReturn := montecarlo.NormalFloat64(meanReturn, stddevReturn)
                portfolioValue = (1 + dailyReturn)
            }
            return portfolioValue
        },
    }

    montecarlo.Stats(simulation.Run()) 
    // {Mean:12036.539859964061 StandardDeviation:2279.4204494720975 Min:6481.218067148653 Max:23742.014047666456}
}
```

### Running multiple simulations

```go
package main

import (
    "github.com/datamango-uk/montecarlo"
)

func main() {
    simulation := montecarlo.Simulation{...}
    simulation.RunMultiple(3)
    // []StatsResults{
    //    {Results: [0.0, ...] Stats:{Mean:0.0 StandardDeviation:0.0 Min:0.0 Max:0.0},
    //    {Results: [0.0, ...] Stats:{Mean:0.0 StandardDeviation:0.0 Min:0.0 Max:0.0},
    //    {Results: [0.0, ...] Stats:{Mean:0.0 StandardDeviation:0.0 Min:0.0 Max:0.0},
    // }
}
```

## License

This project is licensed under the MIT License.