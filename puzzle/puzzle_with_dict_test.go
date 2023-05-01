//go:build !no_dict
// +build !no_dict

package puzzle

import (
	"testing"
)

func TestGeneratePuzzle_RequiresDict(t *testing.T) {
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
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GeneratePuzzle(tc.size, tc.words, 0, tc.difficulty, "../dict-en.txt", false)
			if tc.expectedError && err == nil {
				t.Errorf("Expected an error but didn't get one")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Didn't expect an error but got one: %v", err)
			}
		})
	}
}
