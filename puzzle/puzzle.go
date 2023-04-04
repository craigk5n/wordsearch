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
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	rand.Seed(time.Now().UnixNano())
	logger = logrus.New()
}

type placedSearchWord struct {
	word string
	x    int
	y    int
	dx   int
	dy   int
}
type Puzzle struct {
	grid        [][]rune
	solution    [][]rune
	placedWords []placedSearchWord
}

type Grid [][]rune

func GeneratePuzzle(gridSize int, words []string, columns int, difficulty int, dictionaryPath string, verbose bool) (Puzzle, error) {
	// Set the default log level
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	logger.Logf(logrus.InfoLevel, "Generating puzzle...")

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
			return puzzle, fmt.Errorf("invalid word '%s' (%d is too long)", word, len(word))
		}
	}

	err = insertWordsIntoGrid(&puzzle, validWords, difficulty, dictionaryPath, verbose)
	if err != nil {
		return puzzle, err
	}

	err = fillEmptyCells(puzzle.grid)
	if err != nil {
		return puzzle, err
	}

	for _, placedWord := range puzzle.placedWords {
		logger.Logf(logrus.DebugLevel, "Word: %s, X: %d, Y: %d, dX: %d, dY: %d\n",
			placedWord.word, placedWord.x, placedWord.y, placedWord.dx, placedWord.dy)
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

func insertWordsIntoGrid(puzzle *Puzzle, words []string, difficulty int, dictionaryPath string, verbose bool) error {
	dictionary, err := LoadDictionary(dictionaryPath)
	if err != nil {
		return err
	}

	randomWords := dictionary.RandomWords(numberOfRandomWords(difficulty))

	logger.Logf(logrus.DebugLevel, "Inserting search words.\n")
	for _, word := range words {
		logger.Logf(logrus.DebugLevel, "Attempting to insert word: %s\n", word)
		if !tryInsertWord(puzzle, word, true, verbose) {
			return errors.New("Failed to insert word into the grid: " + word)
		}
		logger.Logf(logrus.DebugLevel, "Successfully inserted word: %s\n", word)
	}

	logger.Logf(logrus.DebugLevel, "Inserting random words.\n")
	for _, word := range randomWords {
		logger.Logf(logrus.DebugLevel, "Attempting to insert random word: %s\n", word)
		if !tryInsertWord(puzzle, word, false, verbose) {
			logger.Logf(logrus.DebugLevel, "Failed to insert random word: %s\n", word)
		}
	}

	logger.Logf(logrus.DebugLevel, "Inserting close words.\n")
	adjustedWords := adjustWordsForDifficulty(words, difficulty, dictionary)
	for _, word := range adjustedWords {
		for _, closeMatch := range dictionary.CloseMatches(word) {
			logger.Logf(logrus.DebugLevel, "Attempging to insert close word: %s\n", word)
			tryInsertWord(puzzle, closeMatch, false, verbose)
		}
	}

	return nil
}

func tryInsertWord(puzzle *Puzzle, word string, isSearchWord bool, verbose bool) bool {
	gridSize := len(puzzle.grid)

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
		placeWord(puzzle, word, bestX, bestY, bestDx, bestDy, isSearchWord)
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
			logger.Logf(logrus.DebugLevel, "(%d, %d) is out of bounds for grid\n", newX, newY)
			return false
		}

		if !isEmptyCell(grid, newX, newY) && grid[newX][newY] != r {
			logger.Logf(logrus.DebugLevel, "(%d, %d) is not an empty cell\n", newX, newY)
			return false
		}
	}
	logger.Logf(logrus.DebugLevel, "(%d, %d) can be used to place %s\n", x, y, word)
	return true
}

func placeWord(puzzle *Puzzle, word string, x, y, dx, dy int, isSearchWord bool) {
	for i, r := range word {
		newX := x + i*dx
		newY := y + i*dy
		puzzle.grid[newX][newY] = r
	}
	if isSearchWord {
		for i, r := range word {
			newX := x + i*dx
			newY := y + i*dy
			puzzle.solution[newX][newY] = r
		}
		puzzle.placedWords = append(puzzle.placedWords, placedSearchWord{word: word, x: x, y: y, dx: dx, dy: dy})
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
			_, err = writer.WriteRune(puzzle.grid[x][y])
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
				_, err = writer.WriteRune(puzzle.solution[x][y])
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

func GeneratePDF(puzzle Puzzle, title string, words []string, columns int, outputFile string) (Puzzle, error) {
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
		return puzzle, err
	}

	return puzzle, nil
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

	var gr [][]rune
	if isSolution {
		gr = puzzle.solution
		if isSolution {
			pageX, pageY := pdf.GetXY()
			pdf.SetLineWidth(1)             // Set the line width
			pdf.SetDrawColor(192, 192, 192) // Set the line color (black)
			for _, word := range puzzle.placedWords {
				startX := float64(word.x)*cellSize + (cellSize / 2)
				startY := float64(word.y)*cellSize + (cellSize / 2)
				endX := startX + float64(word.dx*(len(word.word)-1))*cellSize
				endY := startY + float64(word.dy*(len(word.word)-1))*cellSize
				pdf.Line(pageX+startX, pageY+startY, pageX+endX, pageY+endY)
			}
		}
	} else {
		gr = puzzle.grid
	}
	for y := 0; y < len(gr); y++ {
		for x := 0; x < len(puzzle.grid); x++ {
			pdf.CellFormat(cellSize, cellSize, string(gr[x][y]), "0", 0, "C", false, 0, "")
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
