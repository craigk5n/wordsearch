# Word Search Puzzle Generator

This application can generate a wordsearch puzzle from a provided YAML configuratio
file that specifies the title, puzzle size, search words, difficulty (1-9), and a dictionary file.
The dictionary file is used to insert random words into the puzzle when the
difficulty level is set higher.

## Building from Source
```
go build
```

## Example Configuration File
```yaml
title: Colors
difficulty: 5
words:
  - RED
  - GREEN
  - BLUE
  - YELLOW
  - BLACK
  - WHITE
  - PURPLE
  - PINK
  - ORANGE
  - GRAY
  - BROWN
  - TURQUOISE
  - GOLD
  - SILVER
  - NAVY
  - TEAL
  - MAROON
  - BEIGE
  - MAGENTA
  - OLIVE
```

## Usage

To create a puzzle and generate PDF and plain text output, use a command like the following:
```
./wordsearch -d dict-en.txt -i examples/colors.yaml
```
