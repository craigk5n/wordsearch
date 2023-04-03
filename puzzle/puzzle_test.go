package puzzle

import (
	"testing"
)

func TestGeneratePuzzle(t *testing.T) {
	testCases := []struct {
		name          string
		size          int
		difficulty    int
		words         []string
		expectedError bool
	}{
		{
			name:          "Valid puzzle configuration",
			size:          10,
			difficulty:    1,
			words:         []string{"APPLE", "BANANA", "CHERRY"},
			expectedError: false,
		},
		{
			name:          "Invalid grid size",
			size:          3,
			difficulty:    1,
			words:         []string{"APPLE", "BANANA", "CHERRY"},
			expectedError: true,
		},
		{
			name:          "Invalid difficulty",
			size:          10,
			difficulty:    0,
			words:         []string{"APPLE", "BANANA", "CHERRY"},
			expectedError: true,
		},
		{
			name:          "Invalid words",
			size:          10,
			difficulty:    1,
			words:         []string{"APPLE", "BANANA", "1234"},
			expectedError: true,
		},
		{
			name:          "Invalid words",
			size:          10,
			difficulty:    1,
			words:         []string{"APPLE", "BANANA", "HOWNOWBROWNCOW"},
			expectedError: true,
		},
		{
			name:          "Valid words",
			size:          7,
			difficulty:    1,
			words:         []string{"APPLE", "BANANA", "ORANGE", "SHIH TZU"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GeneratePuzzle(tc.size, tc.words, 0, tc.difficulty, "../words.txt", false)
			if tc.expectedError && err == nil {
				t.Errorf("Expected an error but didn't get one")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Didn't expect an error but got one: %v", err)
			}
		})
	}
}

func TestProcessWords(t *testing.T) {
	testCases := []struct {
		name          string
		words         []string
		gridSize      int
		expectedError bool
		expectedWords []string
	}{
		{
			name:          "Valid words",
			words:         []string{"apple", "banana", "cherry"},
			gridSize:      10,
			expectedError: false,
			expectedWords: []string{"APPLE", "BANANA", "CHERRY"},
		},
		{
			name:          "Invalid word length",
			words:         []string{"apple", "banana", "watermelon"},
			gridSize:      9,
			expectedError: true,
		},
		{
			name:          "Empty words list",
			words:         []string{},
			gridSize:      10,
			expectedError: true,
		},
		{
			name:          "Words with whitespace",
			words:         []string{" apple ", " banana ", " shih tzu "},
			gridSize:      10,
			expectedError: false,
			expectedWords: []string{"APPLE", "BANANA", "SHIHTZU"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			processedWords, err := processWords(tc.words, tc.gridSize)
			if tc.expectedError && err == nil {
				t.Errorf("Expected an error but didn't get one")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Didn't expect an error but got one: %v", err)
			}

			if !tc.expectedError && !compareStringSlices(processedWords, tc.expectedWords) {
				t.Errorf("Expected words %v, but got %v", tc.expectedWords, processedWords)
			}
		})
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
