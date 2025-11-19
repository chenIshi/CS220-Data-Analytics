# Assignment 3 – Data Analytics Toolkit

This folder contains:
- A small Go utility to convert the MATLAB-style data file `HW3_data.m` into CSV files (`X.csv`, `Y.csv`).
- A reusable data factory package to load `X.csv` into typed records.
- A stats package with statistical functions, k-NN classifier, and correlation analysis.
- Multiple CLI tools for data analysis, preprocessing, and machine learning.

## Contents
- `HW3_data.m`: MATLAB-like definitions for `X` (matrix) and `Y` (label vector).
- `parse_matlab.go`: Parser that reads `X=[...]` and `Y=[...]` and writes CSVs.
- `datafactory/`: Shared package with `StudentRecord`, `LoadX`, and a `Factory`.
- `stats/`: Shared package with statistical functions, preprocessing, k-NN, and correlation.
- `cmd/age_stats/`: CLI that computes age statistics.
- `cmd/summary/`: CLI that computes comprehensive statistics (mode, frequencies, quantiles).
- `cmd/correlation/`: CLI that computes Pearson correlations between X columns and Y.
- `cmd/missing/`: CLI that analyzes missing values (-1) in CSV files.
- `cmd/knn/`: CLI that performs k-NN classification with preprocessing and cross-validation.
- `cmd/plot_knn/`: CLI that generates error vs k visualization for k-NN analysis.
- `go.mod`: Go module for running code in this folder.

## Quick Start
```bash
cd assignment3

# 1. Parse the MATLAB-like file and write X.csv and Y.csv
go run parse_matlab.go HW3_data.m

# 2. Compute basic stats (average/median age) from X.csv
go run ./cmd/age_stats -file X.csv

# 3. Compute comprehensive statistics (mode, frequencies, quantiles)
go run ./cmd/summary

# 4. Analyze missing values
go run ./cmd/missing X.csv

# 5. Compute Pearson correlations between X columns and Y
go run ./cmd/correlation

# 6. Run k-NN classification with preprocessing and cross-validation
# Default (Euclidean distance)
go run ./cmd/knn

# Use cosine distance instead of Euclidean
go run ./cmd/knn -metric cosine

# Run a single-K evaluation (e.g., k=5)
go run ./cmd/knn -k 5

# Combine flags (single-K with cosine)
go run ./cmd/knn -k 5 -metric cosine

# 7. Generate error vs k plot (Euclidean distance)
go run ./cmd/plot_knn
```

If successful, you’ll see a summary like:
```
Parsed X: 60 rows x 6 cols
Parsed Y: 60 rows
Wrote X.csv and Y.csv
```

## Outputs
- `X.csv`: Rows correspond to samples; columns to features (for the provided file: 60×6).
- `Y.csv`: Single-column vector of length 60.

Values are written as plain numeric strings. The parser preserves the data as-is from `HW3_data.m`.

## Packages
- `github.com/chenIshi/CS220-Data-Analytics/assignment3/datafactory`
    - Types: `StudentRecord`, `Factory` (field `X []StudentRecord`)
    - Functions: `LoadX(path) ([]StudentRecord, error)`, `NewFromCSV(path) (*Factory, error)`
- `github.com/chenIshi/CS220-Data-Analytics/assignment3/stats`
    - **Statistical Functions**: `AverageAge`, `MedianAge`, `ModeStudentID`, `GenderFrequency`, `PreReqFrequency`, `GPAQuantiles`, `PreTestQuantiles`
    - **Correlation**: `PearsonCorrelation(x, y []float64) float64`
    - **Preprocessing**: `MedianFloat64`, `ModeInt`, `ZScoreNormalize`
    - **Machine Learning**: `KNNClassifier` with `Predict` and `ErrorRate` methods

## Available Commands

| Command | Purpose | Sample Output |
|---------|---------|---------------|
| `cmd/age_stats` | Compute average and median age | Average: 25.3, Median: 24.0 |
| `cmd/summary` | Comprehensive statistics (mode, frequencies, quantiles) | Mode StudentID, Gender/Prereq frequencies, GPA/PreTest quantiles |
| `cmd/missing` | Analyze missing values (-1) in CSV | 37/360 cells (10.28%), 31/60 rows (51.67%) |
| `cmd/correlation` | Pearson correlations between X columns and Y | GPA: 0.6339, PreTest: 0.4246 |
| `cmd/knn` | k-NN classification with preprocessing | Sweep k∈{2,4,6,8,10,12,14} over 5 runs by default; or single-K via `-k N`. Distance selectable via `-metric euclidean|cosine`. |
| `cmd/plot_knn` | Generate error vs k visualization | Creates `knn_error_vs_k.png` showing training and test error curves for bias-variance analysis |

## Data Analysis Pipeline

1. **Parse**: Convert MATLAB format to CSV (`parse_matlab.go`)
2. **Load**: Read data with `datafactory` package
3. **Explore**: Compute statistics (`cmd/summary`, `cmd/missing`)
4. **Correlate**: Find feature-outcome relationships (`cmd/correlation`)
5. **Predict**: Train k-NN classifier with preprocessing (`cmd/knn`)
6. **Visualize**: Generate error vs k plots for model selection (`cmd/plot_knn`)

## k-NN Classification

The k-NN implementation uses:
- **Dataset**: N=60 samples (48 train, 12 test with 80/20 split)
- **Features**: Average GPA (col 3), Prereq Taken (col 4), Pre-test Score (col 5)
- **Preprocessing**:
  - Imputation: Median for continuous (GPA, PreTest), mode for categorical (Prereq)
  - Normalization: Z-score (zero mean, unit variance) per feature
- **Distance**: `euclidean` (default) or `cosine` selected via `-metric` flag
- **Training**: 80/20 random train-test split, 5 runs per k value
- **Evaluation**: Error rate on both train and test sets

### Model Selection and Bias-Variance Tradeoff

Use `cmd/plot_knn` to visualize training vs test error across k values:
- **Low k (2-4)**: Low training error but high test error → overfitting
- **Optimal k (4-8)**: Best balance, following √N ≈ 7 rule for N=48 training samples
- **High k (>10)**: Both errors increase → underfitting (averaging over >20% of training data)

The plot shows the classic U-shaped test error curve demonstrating the bias-variance tradeoff. For this dataset, k=4 to k=8 typically yields the best generalization.

## Data Notes
- The `.m` file uses bracketed rows for `X` and a column vector for `Y`, with commas/spaces allowed.
- **Missing Values**: In `X.csv`, `-1` denotes missing values. In `Y.csv`, `-1` and `1` are valid binary labels (pass/fail).
- The k-NN pipeline handles missing values through imputation (median for continuous features, mode for categorical).
- If you modify the `.m` file, keep `X=` and `Y=` definitions in the same basic structure (rows for `X`, one value per line for `Y`).

## Troubleshooting
- Ensure you run from `assignment3` so the local `go.mod` is used.
- If you see an error like “did not find non-empty X and Y definitions,” verify `HW3_data.m` still has `X=[ ... ];` and `Y=[ ... ];` blocks.
- Go version: this module targets the version declared in `go.mod` and works with a recent Go toolchain.
