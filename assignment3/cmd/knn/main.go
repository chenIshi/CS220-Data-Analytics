package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

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
	// Flags
	metric := flag.String("metric", "euclidean", "distance metric: euclidean or cosine")
	kFlag := flag.Int("k", 0, "number of neighbors; if > 0, run a single-K evaluation; default 0 runs sweep {2,4,6,8,10,12,14}")
	flag.Parse()

	// Load X and Y
	X, err := readCSV("X.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading X: %v\n", err)
		os.Exit(1)
	}
	Y, err := readCSV("Y.csv")
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
	for i := 0; i < len(Y); i++ {
		yVals[i] = Y[i][0]
	}

	// Extract selected features: Avg GPA (col 3), Prereq Taken (col 4), Pre-test Score (col 5)
	// Column indices: 0=StudentID, 1=Gender, 2=Age, 3=GPA, 4=Prereq, 5=PreTest
	colIndices := []int{3, 4, 5}

	// Step 1: Impute missing values (-1)
	// GPA (col 3): median, Prereq (col 4): mode, PreTest (col 5): median
	for _, colIdx := range colIndices {
		var validVals []float64
		var validInts []int
		for i := 0; i < len(X); i++ {
			if X[i][colIdx] != -1 {
				validVals = append(validVals, X[i][colIdx])
				validInts = append(validInts, int(X[i][colIdx]))
			}
		}

		var fillVal float64
		if colIdx == 3 { // GPA: median
			fillVal = stats.MedianFloat64(validVals)
		} else if colIdx == 4 { // Prereq: mode
			fillVal = float64(stats.ModeInt(validInts))
		} else if colIdx == 5 { // PreTest: median
			fillVal = stats.MedianFloat64(validVals)
		}

		// Fill missing values
		for i := 0; i < len(X); i++ {
			if X[i][colIdx] == -1 {
				X[i][colIdx] = fillVal
			}
		}
	}

	// Step 2: Extract selected features into a new matrix
	features := make([][]float64, len(X))
	for i := 0; i < len(X); i++ {
		features[i] = make([]float64, len(colIndices))
		for j, colIdx := range colIndices {
			features[i][j] = X[i][colIdx]
		}
	}

	// Step 3: Normalize each feature to zero mean and unit std deviation
	numFeatures := len(colIndices)
	for j := 0; j < numFeatures; j++ {
		// Extract column
		col := make([]float64, len(features))
		for i := 0; i < len(features); i++ {
			col[i] = features[i][j]
		}
		// Normalize
		normalized := stats.ZScoreNormalize(col)
		// Write back
		for i := 0; i < len(features); i++ {
			features[i][j] = normalized[i]
		}
	}

	// Step 4: Run k-NN for either a single k or the sweep {2,4,6,8,10,12,14}
	var kValues []int
	if *kFlag > 0 {
		kValues = []int{*kFlag}
	} else {
		kValues = []int{2, 4, 6, 8, 10, 12, 14}
	}
	numRuns := 5

	fmt.Println("k-NN Classification Results")
	fmt.Println("============================")
	fmt.Printf("Features used: Avg GPA, Prereq Taken, Pre-test Score\n")
	fmt.Printf("Train/Test split: 80%%/20%%\n")
	fmt.Printf("Number of runs: %d\n\n", numRuns)
	fmt.Printf("%-5s  %-15s  %-15s\n", "k", "Avg Train Error", "Avg Test Error")
	fmt.Println("-----------------------------------------------")

	rand.Seed(time.Now().UnixNano())

	for _, k := range kValues {
		var trainErrors []float64
		var testErrors []float64

		for run := 0; run < numRuns; run++ {
			// Random 80/20 split
			n := len(features)
			indices := rand.Perm(n)
			trainSize := int(0.8 * float64(n))

			trainX := make([][]float64, trainSize)
			trainY := make([]float64, trainSize)
			testX := make([][]float64, n-trainSize)
			testY := make([]float64, n-trainSize)

			for i := 0; i < trainSize; i++ {
				idx := indices[i]
				trainX[i] = features[idx]
				trainY[i] = yVals[idx]
			}
			for i := trainSize; i < n; i++ {
				idx := indices[i]
				testX[i-trainSize] = features[idx]
				testY[i-trainSize] = yVals[idx]
			}

			// Train k-NN (just store training data)
			knn := &stats.KNNClassifier{
				K:       k,
				TrainX:  trainX,
				TrainY:  trainY,
				Metric:  *metric,
			}

			// Compute errors
			trainErr := knn.ErrorRate(trainX, trainY)
			testErr := knn.ErrorRate(testX, testY)

			trainErrors = append(trainErrors, trainErr)
			testErrors = append(testErrors, testErr)
		}

		// Compute averages
		avgTrainErr := 0.0
		avgTestErr := 0.0
		for i := 0; i < numRuns; i++ {
			avgTrainErr += trainErrors[i]
			avgTestErr += testErrors[i]
		}
		avgTrainErr /= float64(numRuns)
		avgTestErr /= float64(numRuns)

		fmt.Printf("%-5d  %.4f (%.2f%%)   %.4f (%.2f%%)\n", 
			k, avgTrainErr, avgTrainErr*100, avgTestErr, avgTestErr*100)
	}
}
