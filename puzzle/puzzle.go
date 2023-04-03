package puzzle

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const marginX = 10.0
const marginY = 10.0
const spacing = 5.0
const lineHeight = 16.0

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Puzzle struct {
	grid     [][]rune
	solution [][]rune
}

type Grid [][]rune

func GeneratePuzzle(gridSize int, words []string, columns int, difficulty int, dictionaryPath string, verbose bool) (Puzzle, error) {
	rand.Seed(time.Now().UnixNano())

	puzzle := createPuzzle(gridSize)

	// Validate difficulty
	if difficulty < 1 || difficulty > 9 {
		return puzzle, fmt.Errorf("invalid difficulty %d (only 1-9 allowed)", difficulty)
	}

	// Validate & process words
	validWords, err := processWords(words, gridSize)
	if err != nil {
		return puzzle, err
	}
	for _, word := range validWords {
		if !isValidWord(word) {
			return puzzle, fmt.Errorf("invalid word '%s'", word)
		}
		if len(word) > gridSize {
			return puzzle, fmt.Errorf("invalid word '%s' (too long)", word)
		}
	}

	err = insertWordsIntoGrid(puzzle, validWords, difficulty, dictionaryPath, verbose)
	if err != nil {
		return puzzle, err
	}

	err = fillEmptyCells(puzzle.grid)
	if err != nil {
		return puzzle, err
	}

	return puzzle, nil
}

func createPuzzle(size int) Puzzle {
	var puzzle Puzzle

	puzzle.grid = createEmptyGrid(size)
	puzzle.solution = createEmptyGrid(size)

	return puzzle

}

// createEmptyGrid initializes an empty square grid with the specified size.
func createEmptyGrid(size int) Grid {
	grid := make(Grid, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]rune, size)
		for j := 0; j < size; j++ {
			grid[i][j] = ' '
		}
	}
	return grid
}

func processWords(words []string, gridSize int) ([]string, error) {
	validWords := make([]string, 0)

	for _, word := range words {
		word = strings.ToUpper(strings.ReplaceAll(word, " ", ""))
		if len(word) == 0 {
			continue
		}

		if len(word) > gridSize {
			return nil, fmt.Errorf("word %s is too long for the grid size", word)
		}

		validWords = append(validWords, word)
	}

	// Sort by word length so longest is first.  It's easier to place long words in
	// the puzzle early.
	sort.Slice(validWords, func(i, j int) bool {
		return len(validWords[i]) > len(validWords[j])
	})

	if len(validWords) == 0 {
		return nil, errors.New("no valid words provided")
	}

	return validWords, nil
}

func insertWordsIntoGrid(puzzle Puzzle, words []string, difficulty int, dictionaryPath string, verbose bool) error {
	dictionary, err := LoadDictionary(dictionaryPath)
	if err != nil {
		return err
	}

	randomWords := dictionary.RandomWords(numberOfRandomWords(difficulty))

	logDebug(verbose, "Inserting search words.\n")
	for _, word := range words {
		logDebug(verbose, "Attempting to insert word: %s\n", word)
		if !tryInsertWord(puzzle, word, true, verbose) {
			return errors.New("Failed to insert word into the grid: " + word)
		}
		logDebug(verbose, "Successfully inserted word: %s\n", word)
	}

	logDebug(verbose, "Inserting random words.\n")
	for _, word := range randomWords {
		logDebug(verbose, "Attempting to insert random word: %s\n", word)
		if !tryInsertWord(puzzle, word, false, verbose) {
			logDebug(verbose, "Failed to insert random word: %s\n", word)
		}
	}

	logDebug(verbose, "Inserting close words.\n")
	adjustedWords := adjustWordsForDifficulty(words, difficulty, dictionary)
	for _, word := range adjustedWords {
		for _, closeMatch := range dictionary.CloseMatches(word) {
			logDebug(verbose, "Attempging to insert close word: %s\n", word)
			tryInsertWord(puzzle, closeMatch, false, verbose)
		}
	}

	return nil
}

func tryInsertWord(puzzle Puzzle, word string, isSearchWord bool, verbose bool) bool {
	gridSize := len(puzzle.grid)

	// Set a seed for the random number generator
	rand.Seed(time.Now().UnixNano())
	// Shuffle the indices randomly
	indicesX := rand.Perm(gridSize)
	indicesY := rand.Perm(gridSize)

	maxOverlap := -1
	bestX, bestY, bestDx, bestDy := -1, -1, -1, -1

	for _, x := range indicesX {
		for _, y := range indicesY {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					if canPlaceWord(puzzle.grid, word, x, y, dx, dy, verbose) {
						overlap := overlappingCells(puzzle.grid, word, x, y, dx, dy)
						if overlap > maxOverlap {
							maxOverlap = overlap
							bestX, bestY, bestDx, bestDy = x, y, dx, dy
						}
					}
				}
			}
		}
	}

	if maxOverlap >= 0 {
		placeWord(puzzle.grid, word, bestX, bestY, bestDx, bestDy, isSearchWord)
		if isSearchWord {
			placeWord(puzzle.solution, word, bestX, bestY, bestDx, bestDy, isSearchWord)
		}
		return true
	}

	return false
}

func overlappingCells(grid Grid, word string, x, y, dx, dy int) int {
	overlapCount := 0
	for i, r := range word {
		newX := x + i*dx
		newY := y + i*dy

		if inBounds(grid, newX, newY) && grid[newX][newY] == r {
			overlapCount++
		}
	}
	return overlapCount
}

func canPlaceWord(grid Grid, word string, x, y, dx, dy int, verbose bool) bool {
	for i, r := range word {
		newX := x + i*dx
		newY := y + i*dy

		if !inBounds(grid, newX, newY) {
			logDebug(verbose, "(%d, %d) is out of bounds for grid\n", newX, newY)
			return false
		}

		if !isEmptyCell(grid, newX, newY) && grid[newX][newY] != r {
			logDebug(verbose, "(%d, %d) is not an empty cell\n", newX, newY)
			return false
		}
	}
	logDebug(verbose, "(%d, %d) can be used to place %s\n", x, y, word)
	return true
}

func placeWord(grid Grid, word string, x, y, dx, dy int, isSearchWord bool) {
	for i, r := range word {
		newX := x + i*dx
		newY := y + i*dy
		grid[newX][newY] = r
	}
	if isSearchWord {
		fmt.Printf("Word=%s, Direction x=%d, y=%d, dx=%d, dy=%d\n", word, x, y, dx, dy)
	}
}

func fillEmptyCells(grid Grid) error {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == ' ' {
				grid[i][j] = randomLetter()
			}
		}
	}
	return nil
}

func randomLetter() rune {
	return rune('A' + rand.Intn(26))
}

func PrintPuzzle(puzzle Puzzle) {
	for _, row := range puzzle.grid {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Println()
	}
}

// SavePuzzleToFile saves the puzzle grid to a file with the specified filename.
func SavePuzzleToFile(puzzle Puzzle, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, row := range puzzle.grid {
		for _, cell := range row {
			_, err := writer.WriteRune(cell)
			if err != nil {
				return err
			}
			_, err = writer.WriteString(" ")
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	_, err = writer.WriteString("\n")
	for _, row := range puzzle.solution {
		for _, cell := range row {
			_, err := writer.WriteRune(cell)
			if err != nil {
				return err
			}
			_, err = writer.WriteString(" ")
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func logDebug(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format, args...)
	}
}

func GeneratePDF(puzzle Puzzle, title string, words []string, columns int, outputFile string) error {
	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a new page
	pdf.AddPage()

	// Set margins (in mm)
	pdf.SetMargins(10, 5, 10)

	// Set font, size, and styles
	pdf.SetFont("Arial", "", 14)

	// Add the title
	drawTitle(pdf, title, 10)
	pdf.Ln(10)

	// Draw the puzzle grid
	drawPuzzleGrid(pdf, puzzle, false)

	// Draw search words
	listSearchWords(pdf, words, columns)

	// Add solution page
	pdf.AddPage()

	// Set margins (in mm)
	pdf.SetMargins(10, 5, 10)

	// Set font, size, and styles
	pdf.SetFont("Arial", "", 14)

	// Add the title
	drawTitle(pdf, title, 10)
	pdf.Ln(10)

	// Draw the puzzle grid
	drawPuzzleGrid(pdf, puzzle, true)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose(outputFile)
	if err != nil {
		return err
	}

	return nil
}

func drawTitle(pdf *gofpdf.Fpdf, title string, marginY float64) {
	pdf.SetFont("Arial", "B", 24)
	titleWidth := pdf.GetStringWidth(title)
	pageWidth, _ := pdf.GetPageSize()
	titleX := (pageWidth - titleWidth) / 2

	pdf.SetXY(titleX, marginY)
	pdf.CellFormat(titleWidth, 10, title, "", 0, "C", false, 0, "")
}

func drawPuzzleGrid(pdf *gofpdf.Fpdf, puzzle Puzzle, isSolution bool) {
	// Calculate the optimal font size and cell size based on the grid size and available page width
	gridSize := len(puzzle.grid)
	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()

	maxCellWidth := (pageWidth - (leftMargin + rightMargin)) / float64(gridSize)
	//fontSize := maxCellWidth * 0.8
	fontSize := maxCellWidth * 2.0
	cellSize := maxCellWidth

	pdf.SetFont("Courier", "B", fontSize)

	/*
		for row := 0; row < len(puzzle.grid); row++ {
			for col := 0; col < len(puzzle.grid); col++ {
				if isSolution {
					pdf.CellFormat(cellSize, cellSize, string(puzzle.solution[row][col]), "0", 0, "C", false, 0, "")
				} else {
					pdf.CellFormat(cellSize, cellSize, string(puzzle.grid[row][col]), "0", 0, "C", false, 0, "")
				}
			}
		}
	*/
	var gr [][]rune
	if isSolution {
		gr = puzzle.solution
	} else {
		gr = puzzle.grid
	}
	for _, row := range gr {
		for _, cell := range row {
			//pdf.CellFormat(cellSize, cellSize, string(cell), "1", 0, "C", false, 0, "")
			pdf.CellFormat(cellSize, cellSize, string(cell), "0", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
}

func listSearchWords(pdf *gofpdf.Fpdf, words []string, columns int) {
	fontSize := float64(11)
	pdf.SetFont("Arial", "", fontSize)
	pdf.Ln(10)
	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	effectiveWidth := pageWidth - leftMargin - rightMargin

	if columns < 1 {
		columns = 5
	}
	wordCount := 0
	width := effectiveWidth / float64(columns)

	for _, word := range words {
		//pdf.CellFormat(40, 10, word, "0", 0, "L", false, 0, "")
		pdf.CellFormat(width, 10, word, "0", 0, "C", false, 0, "")
		wordCount++

		if wordCount%columns == 0 {
			//pdf.Ln(-1)
			pdf.Ln(fontSize / 2)
		}
	}
}
