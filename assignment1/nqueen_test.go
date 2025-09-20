package main

import "testing"

func TestNQueensSolutions(t *testing.T) {
	
	// Define test cases
	// The last case with n=2 is added to test the no-solution scenario
	tests := []struct {
		n      int
		expect int
	}{
		{4, 2},
		{8, 92},
		{2, 0},
	}

	for _, tt := range tests {
		count = 0
		solveNQueens(tt.n)
		if count != tt.expect {
			t.Errorf("For n=%d, expected %d solutions, but got %d", tt.n, tt.expect, count)
		}
	}
}