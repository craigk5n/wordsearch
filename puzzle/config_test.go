package puzzle

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	// Create a temporary YAML config file for testing purposes
	configContent := []byte(`
title: "Sample Puzzle"
size: 15
difficulty: 3
words:
  - "apple"
  - "banana"
  - "orange"
output_basename: "sample_output"
`)
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary test config file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file after test

	if _, err := tmpFile.Write(configContent); err != nil {
		t.Fatalf("Failed to write temporary test config file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary test config file: %v", err)
	}

	config, err := ParseConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseConfig returned error: %v", err)
	}

	if config.Title != "Sample Puzzle" {
		t.Errorf("Expected title to be %q, got %q", "Sample Puzzle", config.Title)
	}

	if config.Size != 15 {
		t.Errorf("Expected size to be %d, got %d", 15, config.Size)
	}

	if config.Difficulty != 3 {
		t.Errorf("Expected difficulty to be %d, got %d", 3, config.Difficulty)
	}

	expectedWords := []string{"apple", "banana", "orange"}
	if len(config.Words) != len(expectedWords) {
		t.Fatalf("Expected %d words, got %d", len(expectedWords), len(config.Words))
	}
	for i, word := range config.Words {
		if word != expectedWords[i] {
			t.Errorf("Expected word %d to be %q, got %q", i, expectedWords[i], word)
		}
	}

	if config.OutputBasename != "sample_output" {
		t.Errorf("Expected output_basename to be %q, got %q", "sample_output", config.OutputBasename)
	}
}
