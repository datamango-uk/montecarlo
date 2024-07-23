package montecarlo

import (
	"math"
	"sync"
)

// RunFunc is a function type that defines the operation to be performed in each iteration of the simulation.
type RunFunc func(inputValues map[string]float64) float64

// Simulation represents a Monte Carlo simulation.
type Simulation struct {
	Iterations  int
	InputValues map[string]float64
	RunFunc     RunFunc
	Workers     int
}

// Worker function
func (s *Simulation) worker(taskChan <-chan struct{}, resultsChan chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for range taskChan {
		resultsChan <- s.RunFunc(s.InputValues)
	}
}

// Run executes the simulation and returns the results.
func (s *Simulation) Run() []float64 {
	if s.Workers == 0 {
		s.Workers = 100
	}
	taskChan := make(chan struct{})
	resultsChan := make(chan float64, s.Iterations)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go s.worker(taskChan, resultsChan, &wg)
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

// Stats represents the statistical results of a simulation.
type Statistics struct {
	Mean              float64
	StandardDeviation float64
	Min               float64
	Max               float64
}

// Stats calculates and returns the statistics from a slice of float64 results.
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

// StatsResults holds the results and statistics of a single simulation run.
type StatsResults struct {
	Results []float64
	Stats   Statistics
}

// RunMultiple executes the simulation multiple times and returns the results and statistics for each run.
func (s *Simulation) RunMultiple(n int) []StatsResults {
	results := make([]StatsResults, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r := s.Run()
			results[idx].Results = r
			results[idx].Stats = Stats(r)
		}(i)
	}
	wg.Wait()

	return results
}
