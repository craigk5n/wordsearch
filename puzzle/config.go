package puzzle

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PuzzleConfig struct {
	Title          string   `yaml:"title"`
	Size           int      `yaml:"size"`
	Columns        int      `yaml:"columns"`
	Difficulty     int      `yaml:"difficulty"`
	Words          []string `yaml:"words"`
	OutputBasename string   `yaml:"output_basename"`
}

func basenameWithoutExt(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)

	return base[:len(base)-len(ext)]
}

func ParseConfig(filename string) (*PuzzleConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config PuzzleConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	// If no output_basename provided, use the basename of the input YAML file.
	if len(config.OutputBasename) == 0 {
		config.OutputBasename = basenameWithoutExt(filename)
	}

	return &config, nil
}
