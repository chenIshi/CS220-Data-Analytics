package stats

import (
    "math"
    "sort"

    "github.com/chenIshi/CS220-Data-Analytics/assignment3/datafactory"
)

// ModeStudentID returns the most frequent (non-missing) StudentID and its count.
// On ties, returns the smallest StudentID.
func ModeStudentID(recs []datafactory.StudentRecord) (int64, int) {
    freq := map[int64]int{}
    var modeID int64
    var modeCount int
    init := false
    for _, r := range recs {
        if r.StudentID < 0 { // treat negative as missing
            continue
        }
        freq[r.StudentID]++
    }
    for id, c := range freq {
        if !init || c > modeCount || (c == modeCount && id < modeID) {
            modeID = id
            modeCount = c
            init = true
        }
    }
    return modeID, modeCount
}

// GenderFrequency returns counts for male=0 and female=1, ignoring missing (-1).
func GenderFrequency(recs []datafactory.StudentRecord) (male0 int, female1 int) {
    for _, r := range recs {
        switch r.Gender {
        case 0:
            male0++
        case 1:
            female1++
        default:
            // ignore -1 or others
        }
    }
    return
}

// PreReqFrequency returns counts for false=0 and true=1, ignoring missing (-1).
func PreReqFrequency(recs []datafactory.StudentRecord) (zero int, one int) {
    for _, r := range recs {
        switch r.PrereqTaken {
        case 0:
            zero++
        case 1:
            one++
        default:
            // ignore -1 or others
        }
    }
    return
}

// percentiles implements nearest-rank percentiles for a sorted slice.
// p in [0,1]. Returns math.NaN() if values is empty.
func percentiles(sorted []float64, p float64) float64 {
    n := len(sorted)
    if n == 0 {
        return math.NaN()
    }
    if p <= 0 {
        return sorted[0]
    }
    if p >= 1 {
        return sorted[n-1]
    }
    // Nearest-rank: k = ceil(p*n) - 1 for 0-indexed
    // For small n, type-7 interpolation is also common, but nearest-rank is simple and robust.
    k := int(math.Ceil(p*float64(n))) - 1
    if k < 0 {
        k = 0
    }
    if k >= n {
        k = n - 1
    }
    return sorted[k]
}

// GPAQuantiles returns P25, P50, P75 for AverageGPA ignoring -1.
func GPAQuantiles(recs []datafactory.StudentRecord) (p25, p50, p75 float64) {
    vals := make([]float64, 0, len(recs))
    for _, r := range recs {
        if r.AverageGPA >= 0 { // ignore -1
            vals = append(vals, r.AverageGPA)
        }
    }
    sort.Float64s(vals)
    return percentiles(vals, 0.25), percentiles(vals, 0.50), percentiles(vals, 0.75)
}

// PreTestQuantiles returns P25, P50, P75 for PreTestScore ignoring -1.
func PreTestQuantiles(recs []datafactory.StudentRecord) (p25, p50, p75 float64) {
    vals := make([]float64, 0, len(recs))
    for _, r := range recs {
        if r.PreTestScore >= 0 { // ignore -1
            vals = append(vals, float64(r.PreTestScore))
        }
    }
    sort.Float64s(vals)
    return percentiles(vals, 0.25), percentiles(vals, 0.50), percentiles(vals, 0.75)
}
