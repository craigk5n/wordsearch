package puzzle

import (
	"io/ioutil"

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

	return &config, nil
}
