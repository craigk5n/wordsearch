package puzzle

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"unicode"
)

func ReadWordsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		if !isValidWord(word) {
			return nil, errors.New("Invalid word found in the input file: " + word)
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func isValidWord(word string) bool {
	if word == "" {
		return false
	}

	for _, r := range word {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
