package main

import (
    "bufio"
    "encoding/csv"
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

// parseMatlabArrays parses a simple MATLAB-like .m file that defines X=[...]; and Y=[...];
// It supports row-wise brackets for X and a column vector for Y, with commas or spaces.
func parseMatlabArrays(path string) ([][]float64, []float64, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, nil, err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    inX, inY := false, false
    var X [][]float64
    var Y []float64

    // Regex helpers
    // Remove surrounding brackets and trailing commas
    trimBrackets := regexp.MustCompile(`^[\[\s]*|(\]|;)[\s]*$`)
    multiSpace := regexp.MustCompile(`\s+`)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" { // skip empty lines
            continue
        }

        // Entry points
        if strings.HasPrefix(line, "X=") && strings.Contains(line, "[") {
            inX, inY = true, false
            // If content on same line after X=[..., handle that by continuing
            continue
        }
        if strings.HasPrefix(line, "Y=") && strings.Contains(line, "[") {
            inX, inY = false, true
            continue
        }

        // Exit points
        if strings.HasPrefix(line, "]") || strings.HasSuffix(line, "];") || line == "]];" || line == "];" {
            inX, inY = false, false
            continue
        }

        if inX {
            // Expect row lines like: [ 450883.0, -1.0, 29.0, 3.55,    0,  5.0]
            l := trimBrackets.ReplaceAllString(line, "")
            l = strings.TrimSpace(l)
            if l == "" { // skip
                continue
            }
            // Replace commas with spaces, collapse spaces
            l = strings.ReplaceAll(l, ",", " ")
            l = multiSpace.ReplaceAllString(l, " ")
            parts := strings.Fields(l)
            var row []float64
            for _, p := range parts {
                v, err := strconv.ParseFloat(p, 64)
                if err != nil {
                    return nil, nil, fmt.Errorf("parse X value '%s': %w", p, err)
                }
                row = append(row, v)
            }
            if len(row) > 0 {
                X = append(X, row)
            }
            continue
        }
        if inY {
            // Expect scalar per line, possibly with trailing commas or spaces
            l := trimBrackets.ReplaceAllString(line, "")
            l = strings.TrimSpace(strings.TrimSuffix(l, ","))
            if l == "" {
                continue
            }
            v, err := strconv.ParseFloat(l, 64)
            if err != nil {
                return nil, nil, fmt.Errorf("parse Y value '%s': %w", l, err)
            }
            Y = append(Y, v)
            continue
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, nil, err
    }

    if len(X) == 0 || len(Y) == 0 {
        return nil, nil, errors.New("did not find non-empty X and Y definitions in file")
    }
    return X, Y, nil
}

func writeCSV(path string, rows [][]float64) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    for _, r := range rows {
        rec := make([]string, len(r))
        for i, v := range r {
            rec[i] = strconv.FormatFloat(v, 'f', -1, 64)
        }
        if err := w.Write(rec); err != nil {
            return err
        }
    }
    return w.Error()
}

func writeCSVCol(path string, col []float64) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()
    for _, v := range col {
        if err := w.Write([]string{strconv.FormatFloat(v, 'f', -1, 64)}); err != nil {
            return err
        }
    }
    return w.Error()
}

func main() {
    outDir := flag.String("out", "", "output directory for CSVs (default: alongside input file)")
    flag.Parse()
    if flag.NArg() < 1 {
        fmt.Println("Usage: go run parse_matlab.go <path/to/HW3_data.m> [-out <dir>]")
        os.Exit(2)
    }
    in := flag.Arg(0)
    X, Y, err := parseMatlabArrays(in)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

    baseDir := *outDir
    if baseDir == "" {
        baseDir = filepath.Dir(in)
    }
    if err := os.MkdirAll(baseDir, 0o755); err != nil {
        fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
        os.Exit(1)
    }

    xPath := filepath.Join(baseDir, "X.csv")
    yPath := filepath.Join(baseDir, "Y.csv")
    if err := writeCSV(xPath, X); err != nil {
        fmt.Fprintf(os.Stderr, "write X.csv: %v\n", err)
        os.Exit(1)
    }
    if err := writeCSVCol(yPath, Y); err != nil {
        fmt.Fprintf(os.Stderr, "write Y.csv: %v\n", err)
        os.Exit(1)
    }

    // Quick summary
    cols := 0
    if len(X) > 0 {
        cols = len(X[0])
    }
    fmt.Printf("Parsed X: %d rows x %d cols\n", len(X), cols)
    fmt.Printf("Parsed Y: %d rows\n", len(Y))
    fmt.Printf("Wrote %s and %s\n", xPath, yPath)
}
