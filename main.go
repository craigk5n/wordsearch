package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/wordsearch/puzzle"
)

const max_puzzle_size int = 1024

func main() {
	inputFile := flag.String("i", "", "YAML input file")
	dictionaryPath := flag.String("d", "", "Custom dictionary file (optional)")
	verbose := flag.Bool("v", false, "enable debug logging")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: YAML input file is required.")
		flag.Usage()
		os.Exit(1)
	}

	config, err := puzzle.ParseConfig(*inputFile)
	if err != nil {
		fmt.Printf("Error: Failed to parse YAML input file: %v\n", err)
		os.Exit(1)
	}

	// We either generate a puzzle of the specified size, or we start with the max word length
	// and keep adding 1 until we successfully generate the puzzle.
	autoSize := (config.Size == 0)
	fmt.Printf("size=%d, autoSize=%v\n", config.Size, autoSize)
	if autoSize {
		// Set puzzle size to longest word
		for _, word := range config.Words {
			if len(word) > config.Size {
				config.Size = len(word)
			}
		}
	}
	p, err := puzzle.GeneratePuzzle(config.Size, config.Words, config.Columns, config.Difficulty, *dictionaryPath, *verbose)
	for ; err != nil && autoSize && config.Size < max_puzzle_size; config.Size = config.Size + 1 {
		fmt.Printf("Generating puzzle of size %d\n", config.Size)
		p, err = puzzle.GeneratePuzzle(config.Size, config.Words, config.Columns, config.Difficulty, *dictionaryPath, *verbose)
	}
	if err != nil {
		fmt.Printf("Error: Failed to generate puzzle: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(config.Title)
	puzzle.PrintPuzzle(p)
	outputFilename := config.OutputBasename + ".txt"
	err = puzzle.SavePuzzleToFile(p, outputFilename, true)
	if err != nil {
		fmt.Printf("Error: Failed to save puzzle to file: %v\n", err)
		os.Exit(1)
	}

	// Generate PDF
	_, err = puzzle.GeneratePDF(p, config.Title, config.Words, config.Columns, config.OutputBasename+".pdf")
	if err != nil {
		fmt.Printf("Error generating PDF: %v", err)
		os.Exit(1)
	}
}
