package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Layout struct {
	X, Y            int
	Rows            []Row `@@*`
	StartDirections []string
}

type Row struct {
	Tiles []Tile `@Tile+ EOL`
}

type Tile string

const (
	Vertical   Tile = "|"
	Horizontal Tile = "-"
	NorthEast  Tile = "L"
	NorthWest  Tile = "J"
	SouthWest  Tile = "7"
	SouthEast  Tile = "F"
	Ground     Tile = "."
	Start      Tile = "S"
)

type Pipe struct {
	X, Y           int
	EntryX, EntryY int
	Tile           Tile
}

func main() {
	steps := FurthestSteps("input1.txt")
	log.Println("Part 1: Furthest steps from loop starting position:", steps)

	tiles := EnclosedTiles("input1.txt")
	log.Println("Part 2: Tiles enclosed by loop:", tiles)
}

func FurthestSteps(filename string) int {
	var steps int

	layout := parse(filename)
	layout.SetStart()

	pipe := layout.FirstStep(layout.X, layout.Y)
	layout.X, layout.Y = pipe.X, pipe.Y
	steps++

	for layout.lookup(layout.X, layout.Y) != Start {
		pipe = layout.NextStep(pipe)
		layout.X, layout.Y = pipe.X, pipe.Y
		steps++
	}

	return steps / 2
}

func EnclosedTiles(filename string) int {
	layout := parse(filename)

	width := len(layout.Rows[0].Tiles)
	height := len(layout.Rows)

	output := make([][]Tile, height)
	i := 0
	for i < height {
		j := 0
		var row []Tile
		for j < width {
			row = append(row, Ground)
			j++
		}
		output[i] = row
		i++
	}

	layout.SetStart()
	startX, startY := layout.X, layout.Y

	pipe := layout.FirstStep(layout.X, layout.Y)
	layout.X, layout.Y = pipe.X, pipe.Y

	output[layout.Y][layout.X] = layout.lookup(layout.X, layout.Y)

	for layout.lookup(layout.X, layout.Y) != Start {
		pipe = layout.NextStep(pipe)
		layout.X, layout.Y = pipe.X, pipe.Y
		output[layout.Y][layout.X] = layout.lookup(layout.X, layout.Y)
	}

	printMap(output)

	output[startY][startX] = layout.convertStart()

	var sum int

	for _, line := range output {
		sum += sumEnclosedGround(line)
	}

	return sum
}

func (l *Layout) SetStart() {
	for idx, row := range l.Rows {
		s := slices.Index(row.Tiles, Start)
		if s == -1 {
			continue
		}

		l.Y = idx
		l.X = s

		break
	}
}

func (l *Layout) FirstStep(X, Y int) Pipe {
	var results []Pipe
	var relative []string

	x, y := l.above(X, Y)
	if x != -1 && y != -1 {
		tile := l.lookup(x, y)
		if tile == Start || tile == Vertical || tile == SouthWest || tile == SouthEast {
			results = append(results, Pipe{X: x, Y: y, EntryX: X, EntryY: Y, Tile: tile})
			relative = append(relative, "above")
		}
	}

	x, y = l.below(X, Y)
	if x != -1 && y != -1 {
		tile := l.lookup(x, y)
		if tile == Start || tile == Vertical || tile == NorthWest || tile == NorthEast {
			results = append(results, Pipe{X: x, Y: y, EntryX: X, EntryY: Y, Tile: tile})
			relative = append(relative, "below")
		}
	}

	x, y = l.left(X, Y)
	if x != -1 && y != -1 {
		tile := l.lookup(x, y)
		if tile == Start || tile == Horizontal || tile == NorthEast || tile == SouthEast {
			results = append(results, Pipe{X: x, Y: y, EntryX: X, EntryY: Y, Tile: tile})
			relative = append(relative, "left")
		}
	}

	x, y = l.right(X, Y)
	if x != -1 && y != -1 {
		tile := l.lookup(x, y)
		if tile == Start || tile == Horizontal || tile == NorthWest || tile == SouthWest {
			results = append(results, Pipe{X: x, Y: y, EntryX: X, EntryY: Y, Tile: tile})
			relative = append(relative, "right")
		}
	}

	l.StartDirections = relative

	return results[0]
}

func (l *Layout) NextStep(pipe Pipe) Pipe {
	tile := l.lookup(pipe.X, pipe.Y)
	switch tile {
	case Vertical:
		x, y := l.above(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.below(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	case Horizontal:
		x, y := l.left(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.right(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	case NorthEast:
		x, y := l.above(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.right(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	case NorthWest:
		x, y := l.above(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.left(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	case SouthWest:
		x, y := l.left(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.below(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	case SouthEast:
		x, y := l.right(pipe.X, pipe.Y)
		if (x != -1 && y != -1) && !(x == pipe.EntryX && y == pipe.EntryY) {
			return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
		}

		x, y = l.below(pipe.X, pipe.Y)

		return Pipe{X: x, Y: y, EntryX: pipe.X, EntryY: pipe.Y, Tile: l.lookup(x, y)}
	}

	return pipe
}

func (l *Layout) convertStart() Tile {
	switch {
	case slices.Contains(l.StartDirections, "above") && slices.Contains(l.StartDirections, "below"):
		return Vertical
	case slices.Contains(l.StartDirections, "left") && slices.Contains(l.StartDirections, "right"):
		return Horizontal
	case slices.Contains(l.StartDirections, "above") && slices.Contains(l.StartDirections, "left"):
		return NorthWest
	case slices.Contains(l.StartDirections, "above") && slices.Contains(l.StartDirections, "right"):
		return NorthEast
	case slices.Contains(l.StartDirections, "below") && slices.Contains(l.StartDirections, "left"):
		return SouthWest
	case slices.Contains(l.StartDirections, "below") && slices.Contains(l.StartDirections, "right"):
		return SouthEast
	default:
		return Start
	}
}

func (l *Layout) lookup(X, Y int) Tile {
	return l.Rows[Y].Tiles[X]
}

func (l *Layout) above(X, Y int) (int, int) {
	if Y-1 < 0 {
		return -1, -1
	}

	return X, Y - 1
}

func (l *Layout) below(X, Y int) (int, int) {
	if Y >= len(l.Rows[Y].Tiles) {
		return -1, -1
	}

	return X, Y + 1
}

func (l *Layout) left(X, Y int) (int, int) {
	if X-1 < 0 {
		return -1, -1
	}

	return X - 1, Y
}

func (l *Layout) right(X, Y int) (int, int) {
	if X >= len(l.Rows[Y].Tiles) {
		return -1, -1
	}

	return X + 1, Y
}

func parse(filename string) Layout {
	layoutLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Tile", `[SLJ7F\|\-\.]`},
		{"EOL", `\n`},
	})

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	parser := participle.MustBuild[Layout](
		participle.Lexer(layoutLexer),
	)
	layout, err := parser.Parse(filename, fh)
	if err != nil {
		log.Fatal("unable to parse paper: %w", err)
	}

	return *layout
}

func sumEnclosedGround(line []Tile) int {
	var sum int
	var inside bool

	for idx, tile := range line {
		if idx == 0 && tile != Ground {
			inside = !inside
			continue
		}

		if tile == Ground {
			if inside {
				sum++
			}
			continue
		}

		switch tile {
		case NorthWest, SouthWest:
			if isU(line, idx) {
				inside = !inside
			}
		case Horizontal:
		default:
			inside = !inside
		}
	}

	return sum
}

func isU(line []Tile, startIndex int) bool {
	if startIndex <= 0 {
		return false
	}

	tile := line[startIndex]

	i := startIndex - 1
	for i >= 0 {
		cur := line[i]
		switch {
		case tile == NorthWest && cur == NorthEast:
			return true
		case tile == SouthWest && cur == SouthEast:
			return true
		case cur == Horizontal:
		default:
			return false
		}

		i--
	}

	return false
}

func printMap(output [][]Tile) {
	for _, line := range output {
		var tiles []string
		for _, tile := range line {
			tiles = append(tiles, prettyTile(tile))
		}
		fmt.Println(strings.Join(tiles, ""))
	}
}

func prettyTile(tile Tile) string {
	switch tile {
	case Start:
		return ""
	case Vertical:
		return "│"
	case Horizontal:
		return "─"
	case NorthEast:
		return "└"
	case NorthWest:
		return "┘"
	case SouthEast:
		return "┌"
	case SouthWest:
		return "┐"
	default:
		return " "
	}
}
