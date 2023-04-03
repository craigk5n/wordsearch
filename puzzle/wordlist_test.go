package puzzle

import (
	"os"
	"testing"
)

func TestReadWordsFromFile(t *testing.T) {
	// Create a temporary file with test words for testing purposes
	wordsContent := []byte("apple\nbanana\norange\n")
	tmpFile, err := os.CreateTemp("", "test_words_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary test words file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file after test

	if _, err := tmpFile.Write(wordsContent); err != nil {
		t.Fatalf("Failed to write temporary test words file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary test words file: %v", err)
	}

	words, err := ReadWordsFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadWordsFromFile returned error: %v", err)
	}

	expectedWords := []string{"apple", "banana", "orange"}
	if len(words) != len(expectedWords) {
		t.Fatalf("Expected %d words, got %d", len(expectedWords), len(words))
	}
	for i, word := range words {
		if word != expectedWords[i] {
			t.Errorf("Expected word %d to be %q, got %q", i, expectedWords[i], word)
		}
	}
}

func TestIsValidWord(t *testing.T) {
	testCases := []struct {
		word     string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"123", false},
		{"apple123", false},
		{"", false},
		{"HelloWorld", true},
	}

	for _, tc := range testCases {
		result := isValidWord(tc.word)
		if result != tc.expected {
			t.Errorf("Expected isValidWord(%q) to be %v, got %v", tc.word, tc.expected, result)
		}
	}
}
