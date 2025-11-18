package stats

import (
	"math"
)

// PearsonCorrelation computes Pearson correlation coefficient between x and y slices.
// For x values: ignores pairs where x[i] is -1 (missing in X).
// For y values: -1 is a valid value (e.g., binary outcome), not missing.
// Returns NaN if insufficient valid pairs.
func PearsonCorrelation(x, y []float64) float64 {
	if len(x) != len(y) {
		return math.NaN()
	}
	
	// Collect valid pairs (only skip if x is -1, y's -1 is valid)
	var pairs [][2]float64
	for i := 0; i < len(x); i++ {
		if x[i] != -1 { // only check x for missing
			pairs = append(pairs, [2]float64{x[i], y[i]})
		}
	}
	
	n := len(pairs)
	if n < 2 {
		return math.NaN()
	}
	
	// Compute means
	var sumX, sumY float64
	for _, p := range pairs {
		sumX += p[0]
		sumY += p[1]
	}
	meanX := sumX / float64(n)
	meanY := sumY / float64(n)
	
	// Compute covariance and standard deviations
	var cov, varX, varY float64
	for _, p := range pairs {
		dx := p[0] - meanX
		dy := p[1] - meanY
		cov += dx * dy
		varX += dx * dx
		varY += dy * dy
	}
	
	if varX == 0 || varY == 0 {
		return math.NaN()
	}
	
	return cov / math.Sqrt(varX*varY)
}
