package puzzle

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const leftMargin = 20
const rightMargin = 20
const topMargin = 20
const bottomMargin = 20

func GeneratePDF(puzzle Puzzle, title string, words []string, columns int, outputFile string,
	backgroundFile string) (Puzzle, error) {
	// Create a new PDF instance
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a new page
	pdf.AddPage()

	// Draw tiled background
	if backgroundFile != "" {
		drawBackground(pdf, backgroundFile)
		// Draw a rectangle in the center of the page for the puzzle and words
		pageW, pageH := pdf.GetPageSize()
		pdf.SetFillColor(255, 255, 255) // Set fill color to white
		pdf.Rect(leftMargin, topMargin, pageW-(leftMargin+rightMargin), pageH-(topMargin+bottomMargin), "F")
	}

	// Set margins (in mm)
	pdf.SetMargins(leftMargin, topMargin, rightMargin)

	// Set font, size, and styles
	pdf.SetFont("Arial", "", 14)

	// Add the title
	drawTitle(pdf, title, topMargin)
	pdf.Ln(10)

	// Draw the puzzle grid
	drawPuzzleGrid(pdf, puzzle, false)

	// Draw search words
	listSearchWords(pdf, words, columns)

	// Add solution page
	pdf.AddPage()

	// Draw tiled background
	if backgroundFile != "" {
		drawBackground(pdf, backgroundFile)
		// Draw a rectangle in the center of the page for the puzzle and words
		pageW, pageH := pdf.GetPageSize()
		pdf.SetFillColor(255, 255, 255) // Set fill color to white
		pdf.Rect(leftMargin, topMargin, pageW-(leftMargin+rightMargin), pageH-(topMargin+bottomMargin), "F")
	}

	// Set margins (in mm)
	pdf.SetMargins(leftMargin, topMargin, rightMargin)

	// Set font, size, and styles
	pdf.SetFont("Arial", "", 14)

	// Add the title
	drawTitle(pdf, title, topMargin)
	pdf.Ln(10)

	// Draw the puzzle grid
	drawPuzzleGrid(pdf, puzzle, true)

	// Save the PDF to a file
	err := pdf.OutputFileAndClose(outputFile)
	if err != nil {
		return puzzle, err
	}

	return puzzle, nil
}

func drawTitle(pdf *gofpdf.Fpdf, title string, marginY float64) {
	pdf.SetFont("Arial", "B", 24)
	titleWidth := pdf.GetStringWidth(title)
	pageWidth, _ := pdf.GetPageSize()
	titleX := (pageWidth - titleWidth) / 2

	pdf.SetXY(titleX, marginY)
	pdf.CellFormat(titleWidth, 10, title, "", 0, "C", false, 0, "")
}

// Draw a grid for the given puzzle, with the optimal font and cell size, and add
// a line for each placed word in the solution grid if applicable.
func drawPuzzleGrid(pdf *gofpdf.Fpdf, puzzle Puzzle, isSolution bool) {
	// Calculate the optimal font size and cell size based on the grid size and available page width
	gridSize := len(puzzle.grid)
	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()

	maxCellWidth := (pageWidth - (leftMargin + rightMargin)) / float64(gridSize)
	//fontSize := maxCellWidth * 0.8
	fontSize := maxCellWidth * 2.0
	cellSize := maxCellWidth

	pdf.SetFont("Courier", "B", fontSize)

	var gr [][]rune
	if isSolution {
		gr = puzzle.solution
		if isSolution {
			pageX, pageY := pdf.GetXY()
			pdf.SetLineWidth(1)             // Set the line width
			pdf.SetDrawColor(192, 192, 192) // Set the line color (black)
			for _, word := range puzzle.placedWords {
				startX := float64(word.x)*cellSize + (cellSize / 2)
				startY := float64(word.y)*cellSize + (cellSize / 2)
				endX := startX + float64(word.dx*(len(word.word)-1))*cellSize
				endY := startY + float64(word.dy*(len(word.word)-1))*cellSize
				pdf.Line(pageX+startX, pageY+startY, pageX+endX, pageY+endY)
			}
		}
	} else {
		gr = puzzle.grid
	}
	for y := 0; y < len(gr); y++ {
		for x := 0; x < len(puzzle.grid); x++ {
			pdf.CellFormat(cellSize, cellSize, string(gr[x][y]), "0", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
}

func drawBackground(pdf *gofpdf.Fpdf, imagePath string) error {
	info := pdf.RegisterImage(imagePath, "")
	if info == nil {
		return fmt.Errorf("failed to load background image %s", imagePath)
	}

	imgWidth := float64(info.Width()) / 10
	imgHeight := float64(info.Height()) / 10

	pageWidth, pageHeight := pdf.GetPageSize()
	x, y := 0.0, 0.0
	for x < pageWidth {
		pdf.ImageOptions(imagePath, x, y, imgWidth, imgHeight, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
		x += imgWidth
	}
	y = pageHeight - imgHeight
	x = 0.0
	for x < pageWidth {
		pdf.ImageOptions(imagePath, x, y, imgWidth, imgHeight, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
		x += imgWidth
	}

	// Draw left and right borders
	x, y = 0.0, imgHeight
	for y < pageHeight-imgHeight {
		pdf.ImageOptions(imagePath, x, y, imgWidth, imgHeight, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
		y += imgHeight
	}
	x = pageWidth - imgWidth
	y = imgHeight
	for y < pageHeight-imgHeight {
		pdf.ImageOptions(imagePath, x, y, imgWidth, imgHeight, false, gofpdf.ImageOptions{ReadDpi: true}, 0, "")
		y += imgHeight
	}
	return nil
}

func listSearchWords(pdf *gofpdf.Fpdf, words []string, columns int) {
	fontSize := float64(11)
	pdf.SetFont("Arial", "", fontSize)
	pdf.Ln(10)
	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	effectiveWidth := pageWidth - leftMargin - rightMargin

	if columns < 1 {
		columns = 5
	}
	wordCount := 0
	width := effectiveWidth / float64(columns)

	for _, word := range words {
		//pdf.CellFormat(40, 10, word, "0", 0, "L", false, 0, "")
		pdf.CellFormat(width, 10, word, "0", 0, "C", false, 0, "")
		wordCount++

		if wordCount%columns == 0 {
			//pdf.Ln(-1)
			pdf.Ln(fontSize / 2)
		}
	}
}
