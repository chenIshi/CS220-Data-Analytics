package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "os"
    "strconv"
)

type stats struct {
    rows             int
    cols             int
    cells            int
    missingCells     int
    rowsWithMissing  int
}

func analyzeCSV(path string) (stats, error) {
    f, err := os.Open(path)
    if err != nil {
        return stats{}, err
    }
    defer f.Close()
    r := csv.NewReader(f)
    recs, err := r.ReadAll()
    if err != nil {
        return stats{}, err
    }
    var s stats
    s.rows = len(recs)
    if s.rows == 0 {
        return s, nil
    }
    s.cols = len(recs[0])
    for i, row := range recs {
        if len(row) != s.cols {
            return s, fmt.Errorf("row %d has %d cols; expected %d", i+1, len(row), s.cols)
        }
        rowHasMissing := false
        for _, cell := range row {
            s.cells++
            v, err := strconv.ParseFloat(cell, 64)
            if err != nil {
                return s, fmt.Errorf("parse error at row %d: %v", i+1, err)
            }
            if v == -1 { // treat -1 as missing
                s.missingCells++
                rowHasMissing = true
            }
        }
        if rowHasMissing {
            s.rowsWithMissing++
        }
    }
    return s, nil
}

func pct(num, den int) float64 {
    if den == 0 {
        return 0
    }
    return (float64(num) / float64(den)) * 100.0
}

func main() {
    file := flag.String("file", "X.csv", "CSV file to analyze (numeric values, -1 = missing)")
    flag.Parse()

    s, err := analyzeCSV(*file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("File: %s\n", *file)
    fmt.Printf("Rows: %d, Cols: %d, Cells: %d\n", s.rows, s.cols, s.cells)
    fmt.Printf("Missing cells (-1): %d (%.2f%% of all cells)\n", s.missingCells, pct(s.missingCells, s.cells))
    fmt.Printf("Rows containing any missing: %d (%.2f%% of rows)\n", s.rowsWithMissing, pct(s.rowsWithMissing, s.rows))
}
