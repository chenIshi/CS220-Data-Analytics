package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/chenIshi/CS220-Data-Analytics/assignment3/datafactory"
    "github.com/chenIshi/CS220-Data-Analytics/assignment3/stats"
)

func main() {
    file := flag.String("file", "X.csv", "path to X.csv (no header)")
    flag.Parse()

    fx, err := datafactory.NewFromCSV(*file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Loaded %d records from %s\n", len(fx.X), *file)

    // 1) Student ID mode
    modeID, modeCount := stats.ModeStudentID(fx.X)
    fmt.Printf("Student ID mode: %d (count=%d)\n", modeID, modeCount)

    // 2) Gender frequencies (0=Male, 1=Female)
    male0, female1 := stats.GenderFrequency(fx.X)
    fmt.Printf("Gender frequency: male(0)=%d, female(1)=%d\n", male0, female1)

    // 3) Age median: already available if needed via stats.MedianAge(fx.X)

    // 4) Average GPA quantiles
    g25, g50, g75 := stats.GPAQuantiles(fx.X)
    fmt.Printf("Average GPA quantiles: P25=%.4f, P50=%.4f, P75=%.4f\n", g25, g50, g75)

    // 5) Pre-Req Taken frequency (0/1)
    pre0, pre1 := stats.PreReqFrequency(fx.X)
    fmt.Printf("Pre-Req Taken frequency: 0=%d, 1=%d\n", pre0, pre1)

    // 6) Pre-Test Score quantiles
    t25, t50, t75 := stats.PreTestQuantiles(fx.X)
    fmt.Printf("Pre-Test Score quantiles: P25=%.2f, P50=%.2f, P75=%.2f\n", t25, t50, t75)
}
