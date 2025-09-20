package main

import (
	"fmt"
	"os"
	"strconv"
)

var count int

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: ./nqueen <board size>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		fmt.Println("Please provide a positive integer for board size.")
		return
	}

	solveNQueens(n)
	fmt.Printf("Total solutions for %d-Queens: %d\n", n, count)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Check if the current board state is valid (no two queens threaten each other in diagonals)
// return true if valid, false otherwise
func isValid(locs []int) bool {
	for i := 0; i < len(locs); i++ {
		for j := i + 1; j < len(locs); j++ {
			// diagonal check: for any two queens at (i, locs[i]) and (j, locs[j]),
			// they are on the same diagonal if abs(i-j) == abs(locs[i]-locs[j])
			if abs(locs[i]-locs[j]) == j-i {
				return false
			}
		}
	}
	return true
}



// DFS-based backtracking to place queens
// input: locs - current board state, row - current row to place a queen
// states: count - number of valid solutions from this state
func backtrack(locs []int, row int) {
	// Base case: all queen locations are swapped at least once
	// Check if the current configuration is valid
	if row == len(locs) {
		if isValid(locs) {
			count++
			// printlocs(locs)
		}
		return
	}

	// Try swapping the current row with each row below it
	for i := row; i < len(locs); i++ {
		// Swap to place a queen at (row, locs[i])
		locs[row], locs[i] = locs[i], locs[row]
		// Recurse to swap queens in the next row
		backtrack(locs, row+1)
		// Backtrack: swap back
		locs[row], locs[i] = locs[i], locs[row]
	}
}

func solveNQueens(n int) {
	locs := make([]int, n)
	// Initialize the locs with column indices
	for i := range locs {
		locs[i] = i
	}
	backtrack(locs, 0)
	// printlocs(locs)
}

func printlocs(locs []int) {
	for _, col := range locs {
		for j := 0; j < len(locs); j++ {
			if j == col {
				fmt.Print(" Q ")
			} else {
				fmt.Print(" . ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}