package puzzle

import (
	"testing"
)

func TestRandomDirection(t *testing.T) {
	for i := 0; i < 100; i++ {
		dx, dy := randomDirection()
		if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
			t.Errorf("Invalid direction (%d, %d)", dx, dy)
		}
	}
}

func TestInBounds(t *testing.T) {
	grid := make(Grid, 5)
	for i := range grid {
		grid[i] = make([]rune, 5)
	}

	testCases := []struct {
		x, y  int
		valid bool
	}{
		{-1, 0, false},
		{0, -1, false},
		{5, 0, false},
		{0, 5, false},
		{0, 0, true},
		{4, 4, true},
	}

	for _, tc := range testCases {
		if inBounds(grid, tc.x, tc.y) != tc.valid {
			t.Errorf("Expected inBounds to return %v for x=%d, y=%d", tc.valid, tc.x, tc.y)
		}
	}
}

func TestIsEmptyCell(t *testing.T) {
	grid := Grid{
		{'A', 'B', 'C', 'D', 'E'},
		{'F', ' ', 'G', ' ', 'H'},
		{'I', 'J', 'K', 'L', 'M'},
		{'N', 'O', ' ', 'P', 'Q'},
		{'R', 'S', 'T', 'U', 'V'},
	}

	testCases := []struct {
		x, y  int
		empty bool
	}{
		{0, 0, false},
		{1, 1, true},
		{2, 2, false},
		{3, 2, true},
	}

	for _, tc := range testCases {
		if isEmptyCell(grid, tc.x, tc.y) != tc.empty {
			t.Errorf("Expected isEmptyCell to return %v for x=%d, y=%d", tc.empty, tc.x, tc.y)
		}
	}
}
