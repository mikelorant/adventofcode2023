package main

import (
	"log"
	"os"

	"github.com/alecthomas/participle/v2"
)

type Almanac struct {
	Seeds []int `"seeds" ":" @Int*`
	Maps  []Map `@@+`
}

type Map struct {
	Source      string  `@Ident "-" "to" "-"`
	Destination string  `@Ident "map" ":"`
	Ranges      []Range `@@+`
}

type Range struct {
	Destination int `@Int`
	Source      int `@Int`
	Length      int `@Int`
}

func main() {
	low := LowestLocationNumber("input1.txt", false)
	log.Println("Part 1: Lowest location number:", low)

	lowPair := LowestLocationNumber("input2.txt", true)
	log.Println("Part 2: Lowest location number:", lowPair)
}

func LowestLocationNumber(filename string, pairs bool) int {
	seeds, maps := parse(filename)
	locs := lowestLocation(seeds, maps, pairs)

	low := -1
	for _, v := range locs {
		if low == -1 || v < low {
			low = v
		}
	}

	return low
}

func lowestLocation(seeds []int, maps map[string]Map, pairs bool) []int {
	var locs []int
	min := -1

	switch pairs {
	case true:
		i := 0
		for i < len(seeds)/2 {
			start := seeds[i*2]
			length := seeds[i*2+1]

			log.Println("Checking Pairs:", start, length)

			for j := 0; j < length; j++ {
				loc := location(maps, start+j)
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

func parse(filemame string) ([]int, map[string]Map) {
	file, err := os.ReadFile(filemame)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	input := string(file)

	parser, err := participle.Build[Almanac]()
	if err != nil {
		log.Fatal("unable to build parser: %w", err)
	}

	almanac, err := parser.ParseString("", input)
	if err != nil {
		log.Fatal("unable to parse almanac: %w", err)
	}

	maps := make(map[string]Map, len(almanac.Maps))
	for _, v := range almanac.Maps {
		maps[v.Source] = v
	}

	return almanac.Seeds, maps
}
