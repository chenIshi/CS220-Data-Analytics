package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/chenIshi/CS220-Data-Analytics/assignment3/stats"
)

func readCSV(path string) ([][]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	data := make([][]float64, len(rows))
	for i, row := range rows {
		data[i] = make([]float64, len(row))
		for j, val := range row {
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("parse error at row %d col %d: %v", i+1, j+1, err)
			}
			data[i][j] = v
		}
	}
	return data, nil
}

func main() {
	xFile := flag.String("x", "X.csv", "path to X.csv")
	yFile := flag.String("y", "Y.csv", "path to Y.csv")
	flag.Parse()

	X, err := readCSV(*xFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading X: %v\n", err)
		os.Exit(1)
	}
	Y, err := readCSV(*yFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading Y: %v\n", err)
		os.Exit(1)
	}

	if len(X) != len(Y) {
		fmt.Fprintf(os.Stderr, "error: X has %d rows, Y has %d rows\n", len(X), len(Y))
		os.Exit(1)
	}

	// Extract Y as single column
	yVals := make([]float64, len(Y))
	for i, row := range Y {
		if len(row) != 1 {
			fmt.Fprintf(os.Stderr, "error: Y row %d has %d columns (expected 1)\n", i+1, len(row))
			os.Exit(1)
		}
		yVals[i] = row[0]
	}

	// Column names
	fieldNames := []string{"Student ID", "Gender", "Age", "Average GPA", "Prereq Taken", "Pre-test Score"}

	fmt.Printf("Pearson Correlations between X columns and Y:\n")
	fmt.Println("==============================================")

	if len(X) == 0 || len(X[0]) == 0 {
		fmt.Println("No data")
		return
	}

	numCols := len(X[0])
	for col := 0; col < numCols; col++ {
		// Extract column
		xCol := make([]float64, len(X))
		for row := 0; row < len(X); row++ {
			xCol[row] = X[row][col]
		}

		// Count valid pairs for debugging
		validPairs := 0
		for i := 0; i < len(xCol); i++ {
			if xCol[i] != -1 { // only x can be missing
				validPairs++
			}
		}

		corr := stats.PearsonCorrelation(xCol, yVals)
		name := fmt.Sprintf("Column %d", col+1)
		if col < len(fieldNames) {
			name = fieldNames[col]
		}
		fmt.Printf("%-20s: %7.4f  (n=%d valid pairs)\n", name, corr, validPairs)
	}
}
