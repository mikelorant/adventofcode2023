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

type Map struct {
	Source      string
	Destination string
	Ranges      []Range
}

type Range struct {
	Destination int
	Source      int
	Length      int
}

const (
	soilToFertilizer      = "soil-to-fertilizer"
	fertilizerToWater     = "fertilizer-to-water"
	waterToLight          = "water-to-light"
	lightToTemperature    = "light-to-temperature"
	temperatureToHumidity = "temperature-to-humidity"
	humidityToLocation    = "humidity-to-location"
)

func main() {
	low := LowestLocationNumber("input1.txt", false)
	log.Println("Part 1: Lowest location number:", low)

	lowPair := LowestLocationNumber("input2.txt", true)
	log.Println("Part 2: Lowest location number:", lowPair)
}

func LowestLocationNumber(filename string, pairs bool) int {
	low := -1
	locs := parse(filename, pairs)
	for _, v := range locs {
		if low == -1 || v < low {
			low = v
		}
	}

	return low
}

func parse(filename string, pairs bool) []int {
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to open file: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}

	lines = mergeLines(lines)
	seeds, maps := parseLines(lines)
	print(seeds, maps)

	var locs []int
	min := -1

	switch pairs {
	case true:
		i := 0
		for i < len(seeds) / 2 {
			start := seeds[i*2]
			length := seeds[i*2+1]

			log.Println("Checking Pairs:", start, length)

			for j := 0; j < length; j++ {
				loc := location(maps, start + j)
				if loc < min || min == -1 {
					min = loc
				}

			}
			i++
		}
		locs = append(locs, min)
	default:
		for _, seed := range seeds {
			locs = append(locs, location(maps, seed))
		}
	}

	return locs
}

func lookup(m Map, num int) int {
	value := -1

	for _, v := range m.Ranges {
		delta := v.Destination - v.Source
		low := v.Source
		high := v.Source + v.Length

		if low > num || num >= high {
			continue
		}

		value = num + delta
	}

	if value == -1 {
		return num
	}

	return value
}

func location(maps map[string]Map, seed int) int {
	name := "seed"

	num := lookup(maps[name], seed)
	name = maps[name].Destination

	for {
		num = lookup(maps[name], num)
		name = maps[name].Destination

		if name == "location" {
			break
		}
	}

	return num
}

func mergeLines(lines []string) []string {
	var section string
	var value string
	var sections []string

	for _, line := range lines {
		if line == "" {
			if value == "" {
				sections = append(sections, section)
				section = ""
				continue
			}
			sections = append(sections, fmt.Sprintf("%s: %s", section, value))
			section = ""
			value = ""
			continue
		}

		if strings.HasPrefix(line, "seeds:") {
			section = line
			continue
		}

		re := regexp.MustCompile(`[a-zA-Z-]+`)
		res := re.FindString(line)
		if res != "" {
			section = res
			continue
		}

		if value == "" {
			value = line
			continue
		}

		value = fmt.Sprintf("%s, %s", value, line)
	}

	sections = append(sections, fmt.Sprintf("%s: %s", section, value))

	return sections
}

func parseLines(lines []string) ([]int, map[string]Map) {
	seeds := parseSeeds(lines[0])
	maps := parseMaps(lines[1:])

	return seeds, maps
}

func parseSeeds(txt string) []int {
	matches := strings.SplitN(txt, ":", 2)

	return toInt(matches[1])
}

func parseMaps(lines []string) map[string]Map {
	maps := make(map[string]Map)

	for _, line := range lines {
		matches := strings.FieldsFunc(line, split)

		name := matches[0]
		values := matches[1:]

		srcDest := strings.Split(name, "-to-")

		var rs []Range
		for _, v := range values {
			arr := toInt(v)

			rs = append(rs, Range{
				Destination: arr[0],
				Source:      arr[1],
				Length:      arr[2],
			})
		}

		maps[srcDest[0]] = Map{
			Source:      srcDest[0],
			Destination: srcDest[1],
			Ranges:      rs,
		}
	}

	return maps
}

func toInt(txt string) []int {
	arr := strings.Fields(txt)
	s := make([]int, len(arr))

	for idx, v := range arr {
		num, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal("unable to convert to int: %w", err)
		}
		s[idx] = num
	}

	return s
}

func split(r rune) bool {
	return r == ':' || r == ','
}

func print(seeds []int, maps map[string]Map) {
	log.Println("Seeds:", seeds)
	log.Println()

	for _, v := range maps {
		log.Println("Source", v.Source)
		log.Println("Destination:", v.Destination)
		for _, vv := range v.Ranges {
			log.Println(fmt.Sprintf("Src: %d Dest: %d Len: %d", vv.Source, vv.Destination, vv.Length))
		}
		log.Println()
	}
}
