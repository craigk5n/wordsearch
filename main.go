package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigk5n/wordsearch/puzzle"
)

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

	p, err := puzzle.GeneratePuzzle(config.Size, config.Words, config.Columns, config.Difficulty, *dictionaryPath, *verbose)
	if err != nil {
		fmt.Printf("Error: Failed to generate puzzle: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(config.Title)
	puzzle.PrintPuzzle(p)
	outputFilename := config.OutputBasename + ".txt"
	err = puzzle.SavePuzzleToFile(p, outputFilename)
	if err != nil {
		fmt.Printf("Error: Failed to save puzzle to file: %v\n", err)
		os.Exit(1)
	}

	// Generate PDF
	err = puzzle.GeneratePDF(p, config.Title, config.Words, config.Columns, config.OutputBasename+".pdf")
	if err != nil {
		fmt.Printf("Error generating PDF: %v", err)
		os.Exit(1)
	}

}
