package datafactory

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
)

// StudentRecord represents one row from X.csv
// Columns: [Student ID, Gender, Age, Average GPA, Prereq Taken, Pre-test Score]
// Missing values are encoded as -1.
type StudentRecord struct {
    StudentID    int64
    Gender       int
    Age          int
    AverageGPA   float64
    PrereqTaken  int
    PreTestScore int
}

// LoadX reads X.csv (no header) and returns parsed records.
func LoadX(path string) ([]StudentRecord, error) {
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
    out := make([]StudentRecord, 0, len(rows))
    for idx, row := range rows {
        if len(row) != 6 {
            return nil, fmt.Errorf("row %d: expected 6 columns, got %d", idx+1, len(row))
        }
        vals := make([]float64, 6)
        for i := 0; i < 6; i++ {
            v, err := strconv.ParseFloat(row[i], 64)
            if err != nil {
                return nil, fmt.Errorf("row %d col %d: %v", idx+1, i+1, err)
            }
            vals[i] = v
        }
        rec := StudentRecord{
            StudentID:    int64(vals[0]),
            Gender:       int(vals[1]),
            Age:          int(vals[2]),
            AverageGPA:   vals[3],
            PrereqTaken:  int(vals[4]),
            PreTestScore: int(vals[5]),
        }
        out = append(out, rec)
    }
    return out, nil
}

// Factory stores parsed datasets and can be shared across programs.
type Factory struct {
    X []StudentRecord
}

// NewFromCSV constructs a Factory by loading records from X.csv path.
func NewFromCSV(path string) (*Factory, error) {
    recs, err := LoadX(path)
    if err != nil {
        return nil, err
    }
    return &Factory{X: recs}, nil
}
