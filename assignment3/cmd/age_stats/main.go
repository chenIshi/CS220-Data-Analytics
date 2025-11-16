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

    avg := stats.AverageAge(fx.X)
    med := stats.MedianAge(fx.X)
    fmt.Printf("Age average (ignoring -1): %v\n", avg)
    fmt.Printf("Age median  (ignoring -1): %v\n", med)
}
