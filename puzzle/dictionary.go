package puzzle

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strings"
)

type Dictionary struct {
	words []string
}

func LoadDictionary(dictionaryPath string) (*Dictionary, error) {
	if dictionaryPath == "" {
		return nil, errors.New("Dictionary path not provided")
	}

	file, err := os.Open(dictionaryPath)
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
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Dictionary{words: words}, nil
}

func (d *Dictionary) RandomWord() string {
	index := rand.Intn(len(d.words))
	return d.words[index]
}

func (d *Dictionary) CloseMatches(word string) []string {
	closeMatches := make([]string, 0)
	wordRunes := []rune(word)
	wordLength := len(wordRunes)

	// Generate close matches by removing one letter from the word
	for i := 0; i < wordLength; i++ {
		partialWord := string(append(wordRunes[:i], wordRunes[i+1:]...))
		closeMatches = append(closeMatches, partialWord)
	}

	// Generate close matches by adding one letter to the word
	for i := 0; i <= wordLength; i++ {
		for _, r := range d.words {
			if len(r) == 1 {
				newWord := string(append(append(wordRunes[:i], []rune(r)...), wordRunes[i:]...))
				closeMatches = append(closeMatches, newWord)
			}
		}
	}

	// Remove duplicates and the original word from close matches
	closeMatches = removeDuplicatesAndOriginalWord(closeMatches, word)

	return closeMatches
}

func removeDuplicatesAndOriginalWord(words []string, originalWord string) []string {
	uniqueWords := make([]string, 0)
	seen := make(map[string]bool)

	for _, word := range words {
		if _, ok := seen[word]; !ok && word != originalWord {
			seen[word] = true
			uniqueWords = append(uniqueWords, word)
		}
	}

	return uniqueWords
}

func (d *Dictionary) RandomWords(count int) []string {
	randomWords := make([]string, 0, count)
	for i := 0; i < count; i++ {
		randomWords = append(randomWords, d.RandomWord())
	}
	return randomWords
}
