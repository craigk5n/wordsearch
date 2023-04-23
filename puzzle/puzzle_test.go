package puzzle

import (
	"testing"
)

func TestFillEmptyCells(t *testing.T) {
	// Create a test grid containing some empty cells
	grid := Grid{
		{'A', 'B', 'C', ' '},
		{'D', ' ', 'F', 'G'},
		{'H', 'I', 'J', 'K'},
		{' ', 'M', 'N', 'O'},
	}

	// Copy the original grid to compare after filling
	originalGrid := make(Grid, len(grid))
	copy(originalGrid, grid)

	err := fillEmptyCells(grid)
	if err != nil {
		t.Errorf("fillEmptyCells returned an unexpected error: %v", err)
	}

	for i := range grid {
		for j := range grid[i] {
			if originalGrid[i][j] == ' ' && grid[i][j] == ' ' {
				t.Errorf("fillEmptyCells did not fill an empty cell at (%d,%d)", i, j)
			}
			if originalGrid[i][j] != ' ' && grid[i][j] != originalGrid[i][j] {
				t.Errorf("fillEmptyCells overwrote a non-empty cell at (%d,%d)", i, j)
			}
			if grid[i][j] != ' ' && !isValidLetter(grid[i][j]) {
				t.Errorf("fillEmptyCells filled an invalid letter at (%d,%d): %c", i, j, grid[i][j])
			}
		}
	}
}

func isValidLetter(r rune) bool {
	// Check if the given rune is a valid uppercase letter
	return r >= 'A' && r <= 'Z'
}

func TestRandomLetter(t *testing.T) {
	r := randomLetter()
	if r < 'A' || r > 'Z' {
		t.Errorf("randomLetter() = %v; want a value between 'A' and 'Z'", r)
	}
}
