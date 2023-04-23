package puzzle

import (
	"bufio"
	"fmt"
	"os"
)

// PrintPuzzle prints the grid of the puzzle one row at a time. It takes in a parameter
// puzzle of type Puzzle.
func PrintPuzzle(puzzle Puzzle) {
	for y := 0; y < len(puzzle.grid); y++ {
		for x := 0; x < len(puzzle.grid); x++ {
			fmt.Printf("%c ", puzzle.grid[x][y])
		}
		fmt.Println()
	}
}

// SavePuzzleToFile saves the puzzle grid to a file with the specified filename.
func SavePuzzleToFile(puzzle Puzzle, filename string, includeSolution bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for y := 0; y < len(puzzle.grid); y++ {
		for x := 0; x < len(puzzle.grid); x++ {
			_, err = writer.WriteString(fmt.Sprintf("%c ", puzzle.grid[x][y]))
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	// Write solution
	if includeSolution {
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
		for y := 0; y < len(puzzle.grid); y++ {
			for x := 0; x < len(puzzle.grid); x++ {
				_, err = writer.WriteString(fmt.Sprintf("%c ", puzzle.solution[x][y]))
				if err != nil {
					return err
				}
			}
			_, err = writer.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}
	return writer.Flush()
}
