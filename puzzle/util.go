package puzzle

import (
	"math/rand"
)

func randomDirection() (int, int) {
	directions := [][2]int{
		{1, 0},   // Right
		{0, 1},   // Down
		{1, 1},   // Diagonal down-right
		{-1, 0},  // Left
		{0, -1},  // Up
		{-1, -1}, // Diagonal up-left
		{1, -1},  // Diagonal up-right
		{-1, 1},  // Diagonal down-left
	}

	index := rand.Intn(len(directions))
	return directions[index][0], directions[index][1]
}

func inBounds(grid Grid, x, y int) bool {
	return x >= 0 && x < len(grid) && y >= 0 && y < len(grid)
}

func isEmptyCell(grid Grid, x, y int) bool {
	return grid[x][y] == ' '
}
