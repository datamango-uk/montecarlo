package montecarlo

import (
	"math"
	"sync"
)

const (
	defaultIterations = 100
	defaultWorkers    = 50
)

// RunFunc is a function type that defines the operation to be performed in each iteration of the simulation.
type Runner interface {
	Run(map[string]float64) float64
}

// Implements Runner
type RunFunc func(map[string]float64) float64

func (f RunFunc) Run(i map[string]float64) float64 {
	return f(i)
}

// Simulation represents a Monte Carlo simulation.
type Simulation struct {
	Iterations int
	Runner     Runner
	Workers    int
}

func New(iterations, workers int, runFunc Runner) *Simulation {
	return &Simulation{
		Iterations: iterations,
		Workers:    workers,
		Runner:     runFunc,
	}
}

// Run executes the simulation and returns a summary.
func (s *Simulation) Run(input map[string]float64) Summary {
	results := s.execute(input)
	return Summary{
		Results:     results,
		Stats:       Stats(results),
		InputValues: input,
	}
}

// RunMultiple executes the simulation for each map of inputs provided.
func (s *Simulation) RunMultiple(inputs ...map[string]float64) []Summary {
	results := make([]Summary, len(inputs))
	var wg sync.WaitGroup
	for i := range inputs {
		wg.Add(1)
		go func(idx int, input map[string]float64) {
			defer wg.Done()
			r := s.Run(input)
			results[idx] = r
		}(i, inputs[i])
	}
	wg.Wait()
	return results
}

// exeucte runs the core part of the simulation using a worker pool, it is not responsible for calculatins statistics,
// rather just the raw []float64 results.
func (s *Simulation) execute(input map[string]float64) []float64 {
	if s.Workers == 0 {
		s.Workers = defaultWorkers
	}
	if s.Iterations == 0 {
		s.Iterations = defaultIterations
	}
	taskChan := make(chan struct{})
	resultsChan := make(chan float64, s.Iterations)
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range taskChan {
				resultsChan <- s.Runner.Run(input)
			}
		}()
	}
	// Send tasks
	go func() {
		for i := 0; i < s.Iterations; i++ {
			taskChan <- struct{}{}
		}
		close(taskChan)
	}()
	// Collect results
	var results = make([]float64, s.Iterations)
	for i := 0; i < s.Iterations; i++ {
		results[i] = <-resultsChan
	}
	wg.Wait()
	close(resultsChan)
	return results
}

// StatsResults holds the results and statistics of a single simulation run.
type Summary struct {
	Results     []float64
	Stats       Statistics
	InputValues map[string]float64
}

// Stats represents the statistical results of a simulation.
type Statistics struct {
	Mean              float64
	StandardDeviation float64
	Min               float64
	Max               float64
}

// Stats calculates and returns Statistics from a slice of float64.
func Stats(results []float64) Statistics {
	var (
		stats Statistics
		sum   float64
		sumSq float64
		l     = float64(len(results))
	)
	if l == 0 {
		return stats
	}
	stats.Min = math.MaxFloat64
	stats.Max = -math.MaxFloat64

	for _, n := range results {
		sum += n
		sumSq += n * n
		if n < stats.Min {
			stats.Min = n
		}
		if n > stats.Max {
			stats.Max = n
		}
	}
	stats.Mean = sum / float64(l)
	stats.StandardDeviation = math.Sqrt((sumSq / l) - (stats.Mean * stats.Mean))

	return stats
}
