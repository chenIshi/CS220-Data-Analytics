# Assignment 3 – Data Parsing Helper

This folder contains:
- A small Go utility to convert the MATLAB-style data file `HW3_data.m` into CSV files (`X.csv`, `Y.csv`).
- A reusable data factory package to load `X.csv` into typed records.
- A stats package and a tiny CLI to compute simple stats (e.g., age average/median).

## Contents
- `HW3_data.m`: MATLAB-like definitions for `X` (matrix) and `Y` (label vector).
- `parse_matlab.go`: Parser that reads `X=[...]` and `Y=[...]` and writes CSVs.
- `datafactory/`: Shared package with `StudentRecord`, `LoadX`, and a `Factory`.
- `stats/`: Shared package with functions like `AverageAge` and `MedianAge`.
- `cmd/age_stats/`: CLI that uses `datafactory` + `stats` to report age stats.
- `go.mod`: Go module for running code in this folder.

## Quick Start
```bash
cd assignment3
# Parse the MATLAB-like file and write X.csv and Y.csv next to it
go run parse_matlab.go HW3_data.m

# Optional: write outputs into a specific directory
go run parse_matlab.go HW3_data.m -out out/

# Compute basic stats (average/median age) from X.csv
go run ./cmd/age_stats -file X.csv
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
    - Functions: `AverageAge([]datafactory.StudentRecord) float64`, `MedianAge([]datafactory.StudentRecord) float64`

## Data Notes
- The `.m` file uses bracketed rows for `X` and a column vector for `Y`, with commas/spaces allowed.
- `-1.0` appears to denote missing values. The parser leaves these untouched so you can decide how to impute or filter later.
- If you modify the `.m` file, keep `X=` and `Y=` definitions in the same basic structure (rows for `X`, one value per line for `Y`).

## Troubleshooting
- Ensure you run from `assignment3` so the local `go.mod` is used.
- If you see an error like “did not find non-empty X and Y definitions,” verify `HW3_data.m` still has `X=[ ... ];` and `Y=[ ... ];` blocks.
- Go version: this module targets the version declared in `go.mod` and works with a recent Go toolchain.
