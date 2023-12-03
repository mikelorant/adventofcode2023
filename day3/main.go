package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Parts map[string][]Number

type Number struct {
	Value int
	X, Y  int
}

type Gear struct {
	X, Y int
}

func main() {
	pn := SumPartNumbers("input1.txt")
	log.Println("Part 1: Sum of part numbers:", pn)

	gr := SumGearRatio("input1.txt")
	log.Println("Part 2: Sum of gear ratio:", gr)
}

func SumPartNumbers(filename string) int {
	var sum int

	schem, num := schematic(filename)
	parts := scan(schem, num)

	for _, v := range parts {
		for _, vv := range v {
			sum += vv.Value
		}
	}

	return sum
}

func SumGearRatio(filename string) int {
	var gs []Gear

	schem, _ := schematic(filename)

	for idx, line := range schem {
		gs = append(gs, gears(line, idx)...)
	}

	return scanGear(schem, gs)
}

func schematic(filename string) ([]string, []Number) {
	var num []Number

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to open file: %w", err)
	}

	var schematic []string
	var i int

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()

		num = append(num, numbers(line, i)...)

		schematic = append(schematic, line)

		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}

	return schematic, num
}

func numbers(txt string, y int) []Number {
	re := regexp.MustCompile(`\d+`)

	matches := re.FindAllString(txt, -1)
	location := re.FindAllStringIndex(txt, -1)

	num := make([]Number, len(matches))

	for idx, v := range matches {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("unable to convert to string: %w", err)
		}

		num[idx].Value = val
		num[idx].X = location[idx][0]
		num[idx].Y = y
	}

	return num
}

func scan(schem []string, nums []Number) Parts {
	parts := make(Parts)

	for _, v := range nums {
		startX := v.X - 1
		startY := v.Y - 1

		endX := v.X + len(strconv.Itoa(v.Value)) + 1
		endY := v.Y + 1

		if startX < 0 {
			startX = 0
		}

		if startY < 0 {
			startY = 0
		}

		if endX > len(schem[v.X]) {
			endX = len(schem[v.X])
		}

		if endY > len(schem[v.Y])-1 {
			endY = len(schem[v.Y]) - 1
		}

		before := string(schem[startY][startX:endX])
		line := string(schem[v.Y][startX:endX])
		after := string(schem[endY][startX:endX])

		sym := symbols(before + line + after)

		for _, vv := range strings.Split(sym, "") {
			parts[vv] = append(parts[vv], Number{
				Value: v.Value,
				X:     v.X,
				Y:     v.Y,
			})
		}
	}

	return parts
}

func scanGear(schem []string, gs []Gear) int {
	var sum int

	for _, v := range gs {
		startX := v.X - 1
		startY := v.Y - 1

		endX := v.X + 1 + 1
		endY := v.Y + 1

		if startX < 0 {
			startX = 0
		}

		if startY < 0 {
			startY = 0
		}

		if endX > len(schem[0]) {
			endX = len(schem[0])
		}

		if endY > len(schem[0])-1 {
			endY = len(schem[0]) - 1
		}

		before := expand(schem[startY], startX, endX)
		line := expand(schem[v.Y], startX, endX)
		after := expand(schem[endY], startX, endX)

		sum += gearRatio(fmt.Sprintf("%s:%s:%s", before, line, after))
	}

	return sum
}

func symbols(txt string) string {
	const PartSymbolExclusion = `0|1|2|3|4|5|6|7|8|9|\.`

	re := regexp.MustCompile(PartSymbolExclusion)
	sym := re.ReplaceAllString(txt, "")

	return sym
}

func gears(txt string, id int) []Gear {
	var gs []Gear

	re := regexp.MustCompile(`\*`)
	loc := re.FindAllStringIndex(txt, -1)

	for _, v := range loc {
		gs = append(gs, Gear{
			X: v[0],
			Y: id,
		})
	}

	return gs
}

func gearRatio(txt string) int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(txt, -1)

	if len(matches) != 2 {
		return 0
	}

	sum := 1

	for _, v := range matches {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("unable to convert to string: %w", err)
		}

		sum *= val
	}

	return sum
}

func expand(txt string, begin, end int) string {
	const num = "01234567890"

	for begin > 0 {
		if !strings.ContainsAny(string(txt[begin]), num) ||
			!strings.ContainsAny(string(txt[begin-1]), num) {
			break
		}

		begin--
	}

	for end < len(txt) {
		if !strings.ContainsAny(string(txt[end-1]), num) ||
			!strings.ContainsAny(string(txt[end]), num) {
			break
		}

		end++
	}

	return string(txt[begin:end])
}
