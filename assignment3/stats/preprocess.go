package stats

import (
	"math"
	"sort"
)

// MedianFloat64 computes median of a float64 slice (does not modify input).
func MedianFloat64(vals []float64) float64 {
	if len(vals) == 0 {
		return math.NaN()
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2.0
}

// ModeInt computes mode of an int slice. On ties, returns smallest value.
func ModeInt(vals []int) int {
	if len(vals) == 0 {
		return 0
	}
	freq := map[int]int{}
	for _, v := range vals {
		freq[v]++
	}
	var modeVal int
	var modeCount int
	init := false
	for v, c := range freq {
		if !init || c > modeCount || (c == modeCount && v < modeVal) {
			modeVal = v
			modeCount = c
			init = true
		}
	}
	return modeVal
}

// ZScoreNormalize normalizes a slice to zero mean and unit std deviation.
// Returns normalized slice. If stddev is zero, returns zeros.
func ZScoreNormalize(vals []float64) []float64 {
	if len(vals) == 0 {
		return vals
	}
	// Compute mean
	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	
	// Compute stddev
	var sumSq float64
	for _, v := range vals {
		d := v - mean
		sumSq += d * d
	}
	stddev := math.Sqrt(sumSq / float64(len(vals)))
	
	// Normalize
	out := make([]float64, len(vals))
	if stddev == 0 {
		return out // all zeros if no variance
	}
	for i, v := range vals {
		out[i] = (v - mean) / stddev
	}
	return out
}
