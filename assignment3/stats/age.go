package stats

import (
    "math"
    "sort"

    "github.com/chenIshi/CS220-Data-Analytics/assignment3/datafactory"
)

// AverageAge returns the mean age ignoring -1 entries. Returns NaN if no valid ages.
func AverageAge(recs []datafactory.StudentRecord) float64 {
    var sum float64
    var cnt int
    for _, r := range recs {
        if r.Age >= 0 {
            sum += float64(r.Age)
            cnt++
        }
    }
    if cnt == 0 {
        return math.NaN()
    }
    return sum / float64(cnt)
}

// MedianAge returns the median age ignoring -1 entries. Returns NaN if no valid ages.
func MedianAge(recs []datafactory.StudentRecord) float64 {
    ages := make([]int, 0, len(recs))
    for _, r := range recs {
        if r.Age >= 0 {
            ages = append(ages, r.Age)
        }
    }
    if len(ages) == 0 {
        return math.NaN()
    }
    sort.Ints(ages)
    n := len(ages)
    if n%2 == 1 {
        return float64(ages[n/2])
    }
    return (float64(ages[n/2-1]) + float64(ages[n/2])) / 2.0
}
