package puzzle

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

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

// GeneratePuzzle creates a Word Search puzzle based on the provided gridSize, words, columns, and difficulty.
// The dictionaryPath is used to load the dictionary for generating random letters in the grid.
func GeneratePuzzle(gridSize int, words []string, columns int, difficulty int, dictionaryPath string,
	verbose bool) (Puzzle, error) {
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
