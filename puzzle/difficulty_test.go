package puzzle

import (
	"testing"
)

func TestAdjustWordsForDifficulty(t *testing.T) {
	words := []string{"apple", "banana", "orange"}
	difficulty := 5
	dictionary := &Dictionary{}

	adjustedWords := adjustWordsForDifficulty(words, difficulty, dictionary)

	if len(adjustedWords) != len(words) {
		t.Errorf("Expected length of adjustedWords to be %d, got %d", len(words), len(adjustedWords))
	}
}

func TestShouldReverseWord(t *testing.T) {
	difficulty := 5
	reversed := 0

	for i := 0; i < 1000; i++ {
		if shouldReverseWord(difficulty) {
			reversed++
		}
	}

	// It's a probabilistic test, so we just check if the number of reversed words is in a reasonable range.
	if reversed < 400 || reversed > 600 {
		t.Errorf("Expected around 50%% of words reversed with difficulty %d, got %d%%", difficulty, 100*reversed/1000)
	}
}

func TestReverseWord(t *testing.T) {
	word := "apple"
	expected := "elppa"
	reversed := reverseWord(word)

	if reversed != expected {
		t.Errorf("Expected reversed word to be %q, got %q", expected, reversed)
	}
}

func TestNumberOfRandomWords(t *testing.T) {
	difficulty := 5
	expected := 10
	result := numberOfRandomWords(difficulty)

	if result != expected {
		t.Errorf("Expected numberOfRandomWords with difficulty %d to be %d, got %d", difficulty, expected, result)
	}
}

func TestNumberOfCloseMatches(t *testing.T) {
	difficulty := 5
	expected := 5
	result := numberOfCloseMatches(difficulty)

	if result != expected {
		t.Errorf("Expected numberOfCloseMatches with difficulty %d to be %d, got %d", difficulty, expected, result)
	}
}
