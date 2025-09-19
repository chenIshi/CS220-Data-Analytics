package main

import (
	"fmt"
)

var count int

func main() {
	n := uint(8)
	solveNQueens(n)
	fmt.Printf("Total solutions for %d-Queens: %d\n", n, count)
}



func solveNQueens(n uint) {
	board := make([]uint, n)
	// Initialize the board with column indices
	for i := range board {
		board[i] = uint(i)
	}
	// swapQueens(board, 0, n)
	printBoard(board)
}

func printBoard(board []uint) {
	for _, col := range board {
		for j := uint(0); j < uint(len(board)); j++ {
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