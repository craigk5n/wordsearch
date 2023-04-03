package puzzle

import (
	"math/rand"
)

func adjustWordsForDifficulty(words []string, difficulty int, dictionary *Dictionary) []string {
	adjustedWords := make([]string, 0, len(words))

	for _, word := range words {
		if shouldReverseWord(difficulty) {
			adjustedWords = append(adjustedWords, reverseWord(word))
		} else {
			adjustedWords = append(adjustedWords, word)
		}
	}

	return adjustedWords
}

func shouldReverseWord(difficulty int) bool {
	return rand.Intn(10) < difficulty
}

func reverseWord(word string) string {
	runes := []rune(word)
	reversed := make([]rune, len(runes))
	for i, r := range runes {
		reversed[len(runes)-1-i] = r
	}
	return string(reversed)
}

func numberOfRandomWords(difficulty int) int {
	return difficulty * 2
}

func numberOfCloseMatches(difficulty int) int {
	return difficulty
}
