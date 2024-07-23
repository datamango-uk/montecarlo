package montecarlo_test

import (
	"testing"

	"github.com/datamango-uk/montecarlo"
)

func TestNormalDistribution(t *testing.T) {
	var (
		mean       = 0.0
		stddev     = 1.0
		sampleSize = 10000
		n1         int
		n2         int
		n3         int
	)
	for i := 0; i < sampleSize; i++ {
		rnd := montecarlo.NormalFloat64(mean, stddev)
		if inRange(rnd, -1, 1) {
			n1++
			n2++
			n3++
			continue
		}
		if inRange(rnd, -2, 2) {
			n2++
			n3++
			continue
		}
		if inRange(rnd, -3, 3) {
			n3++
			continue
		}
	}

	var (
		p1 = float64(n1) / float64(sampleSize) * 100
		p2 = float64(n2) / float64(sampleSize) * 100
		p3 = float64(n3) / float64(sampleSize) * 100
		e1 = 68.0
		e2 = 95.0
		e3 = 99.7
	)
	const tolerance = 1.0

	if p1 < e1-tolerance || p1 > e1+tolerance {
		t.Errorf("Value of p1 %f is not within the tolerance range of the expected value %f", p1, e1)
	}
	if p2 < e2-tolerance || p2 > e2+tolerance {
		t.Errorf("Value of p2 %f is not within the tolerance range of the expected value %f", p2, e2)
	}
	if p3 < e3-tolerance || p3 > e3+tolerance {
		t.Errorf("Value of p3 %f is not within the tolerance range of the expected value %f", p3, e3)
	}

}

func inRange(n, a, b float64) bool {
	if n > a && n < b {
		return true
	}
	return false
}
