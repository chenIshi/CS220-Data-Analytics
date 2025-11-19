package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/chenIshi/CS220-Data-Analytics/assignment3/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
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
	colIndices := []int{3, 4, 5}

	// Step 1: Impute missing values (-1)
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
		col := make([]float64, len(features))
		for i := 0; i < len(features); i++ {
			col[i] = features[i][j]
		}
		normalized := stats.ZScoreNormalize(col)
		for i := 0; i < len(features); i++ {
			features[i][j] = normalized[i]
		}
	}

	// Step 4: Run k-NN for k in {2, 4, 6, 8, 10, 12, 14} with Euclidean distance
	kValues := []int{2, 4, 6, 8, 10, 12, 14}
	numRuns := 5
	metric := "euclidean"

	fmt.Println("Running k-NN with Euclidean distance...")
	fmt.Printf("Train/Test split: 80%%/20%%, Runs: %d\n\n", numRuns)

	rand.Seed(time.Now().UnixNano())

	// Store results for plotting
	avgTrainErrors := make([]float64, len(kValues))
	avgTestErrors := make([]float64, len(kValues))

	for idx, k := range kValues {
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

			// Train k-NN
			knn := &stats.KNNClassifier{
				K:      k,
				TrainX: trainX,
				TrainY: trainY,
				Metric: metric,
			}

			// Compute errors
			trainErr := knn.ErrorRate(trainX, trainY)
			testErr := knn.ErrorRate(testX, testY)

			trainErrors = append(trainErrors, trainErr)
			testErrors = append(testErrors, testErr)
		}

		// Compute averages
		var sumTrain, sumTest float64
		for i := 0; i < numRuns; i++ {
			sumTrain += trainErrors[i]
			sumTest += testErrors[i]
		}
		avgTrainErrors[idx] = sumTrain / float64(numRuns)
		avgTestErrors[idx] = sumTest / float64(numRuns)

		fmt.Printf("k=%-2d  Train Error: %.4f (%.2f%%)  Test Error: %.4f (%.2f%%)\n",
			k, avgTrainErrors[idx], avgTrainErrors[idx]*100,
			avgTestErrors[idx], avgTestErrors[idx]*100)
	}

	// Create plot
	p := plot.New()
	p.Title.Text = "k-NN Error Rate vs k (Euclidean Distance)"
	p.X.Label.Text = "k (Number of Neighbors)"
	p.Y.Label.Text = "Error Rate"

	// Prepare data points
	trainPts := make(plotter.XYs, len(kValues))
	testPts := make(plotter.XYs, len(kValues))
	for i, k := range kValues {
		trainPts[i].X = float64(k)
		trainPts[i].Y = avgTrainErrors[i]
		testPts[i].X = float64(k)
		testPts[i].Y = avgTestErrors[i]
	}

	// Add lines
	trainLine, err := plotter.NewLine(trainPts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating train line: %v\n", err)
		os.Exit(1)
	}
	trainLine.Color = color.RGBA{R: 27, G: 153, B: 139, A: 255} // Blue
	trainLine.Width = vg.Points(2)

	testLine, err := plotter.NewLine(testPts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating test line: %v\n", err)
		os.Exit(1)
	}
	testLine.Color = color.RGBA{R: 237, G: 33, B: 124, A: 255} // Red
	testLine.Width = vg.Points(2)

	// Add scatter points
	trainScatter, err := plotter.NewScatter(trainPts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating train scatter: %v\n", err)
		os.Exit(1)
	}
	trainScatter.Color = color.RGBA{R: 27, G: 153, B: 139, A: 255}
	trainScatter.GlyphStyle.Shape = draw.CircleGlyph{}

	testScatter, err := plotter.NewScatter(testPts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating test scatter: %v\n", err)
		os.Exit(1)
	}
	testScatter.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	testScatter.GlyphStyle.Shape = draw.SquareGlyph{}

	p.Add(trainLine, trainScatter, testLine, testScatter)
	p.Legend.Add("Training Error", trainLine)
	p.Legend.Add("Test Error", testLine)
	p.Legend.Top = true
	p.Legend.Left = false

	// Save plot
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "knn_error_vs_k.png"); err != nil {
		fmt.Fprintf(os.Stderr, "error saving plot: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Plot saved to knn_error_vs_k.png")
	fmt.Println("\nObservations:")
	fmt.Println("- Training error generally increases as k increases (underfitting)")
	fmt.Println("- Test error may decrease initially then increase (bias-variance tradeoff)")
	fmt.Println("- Lower k values (e.g., k=2,4) show lower training error but may overfit")
	fmt.Println("- Optimal k balances training and test error")
}
