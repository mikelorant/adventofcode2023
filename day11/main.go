package main

import (
	"log"
	"os"
	"slices"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Image struct {
	Rows         []Row `@@+`
	Galaxies     []Cell
	ExpandRow    []int
	ExpandColumn []int
}

type Row struct {
	Cells []Cell `@@+ EOL`
}

type Cell struct {
	ID    int
	Value string `@Cell`
	Type  CellType
	X, Y  int
}

type CellType int

const (
	Space CellType = iota
	Galaxy
)

func main() {
	sum1 := SumPathBetweenPairs("input1.txt", 2)
	log.Println("Part 1: Sum of path between pairs:", sum1)

	sum2 := SumPathBetweenPairs("input1.txt", 1000000)
	log.Println("Part 2: Sum of path between pairs:", sum2)
}

func SumPathBetweenPairs(filename string, emptySpaceSize int) int {
	var sum int

	image := parse(filename)
	image.Enhance(emptySpaceSize)

	for idx := range image.Galaxies {
		sum += image.SumGalaxies(idx + 1)
	}

	return sum
}

func parse(filename string) Image {
	imageLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Cell", `[#\.]`},
		{"EOL", `\n`},
	})

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	parser := participle.MustBuild[Image](
		participle.Lexer(imageLexer),
	)
	image, err := parser.Parse(filename, fh)
	if err != nil {
		log.Fatal("unable to parse paper: %w", err)
	}

	return *image
}

func (i *Image) Enhance(emptySpaceSize int) {
	i.expandSpace()
	i.addGalaxies(emptySpaceSize)
}

func (i *Image) addGalaxies(emptySpaceSize int) {
	var galaxies []Cell

	galaxy := 1

	for y, row := range i.Rows {
		for x, cell := range row.Cells {
			if cell.Value == "#" {
				var extraSpacesY int
				for _, row := range i.ExpandRow {
					if row <= y {
						extraSpacesY++
					}
				}

				var extraSpacesX int
				for _, column := range i.ExpandColumn {
					if column <= x {
						extraSpacesX++
					}
				}

				i.Rows[y].Cells[x].Type = Galaxy
				i.Rows[y].Cells[x].ID = galaxy
				i.Rows[y].Cells[x].X = x + (extraSpacesX * (emptySpaceSize - 1))
				i.Rows[y].Cells[x].Y = y + (extraSpacesY * (emptySpaceSize - 1))
				galaxies = append(galaxies, i.Rows[y].Cells[x])
				galaxy++
			}
		}
	}

	i.Galaxies = galaxies
}

func (i *Image) expandSpace() {
	i.expandColumns()
	i.expandRows()
}

func (i *Image) expandRows() {
	var expand []int

	for y, row := range i.Rows {
		var cells []string

		for _, cell := range row.Cells {
			cells = append(cells, cell.Value)
		}

		if slices.Contains(cells, "#") {
			continue
		}

		expand = append(expand, y)
	}

	i.ExpandRow = expand
}

func (i *Image) expandColumns() {
	var expand []int

	for j := 0; j < len(i.Rows[0].Cells); j++ {
		var cells []string

		for _, row := range i.Rows {
			cells = append(cells, row.Cells[j].Value)
		}

		if slices.Contains(cells, "#") {
			continue
		}

		expand = append(expand, j)
	}

	i.ExpandColumn = expand
}

func (i Image) String() string {
	var result []string

	for _, row := range i.Rows {
		var line []string
		for _, cell := range row.Cells {
			line = append(line, cell.Value)
		}
		result = append(result, strings.Join(line, ""))
	}

	return strings.Join(result, "\n")
}

func (i *Image) SumGalaxies(num int) int {
	var sum int

	for j := num; j < len(i.Galaxies); j++ {
		sum += distance(i.Galaxies[num-1], i.Galaxies[j])
	}

	return sum
}

func distance(a, b Cell) int {
	x := abs(a.X - b.X)
	y := abs(a.Y - b.Y)

	return x + y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
