package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
)

type Paper struct {
	Time     []int `"Time" ":" @Int+`
	Distance []int `"Distance" ":" @Int+`
}

func main() {
	win1 := WaysToWin("input1.txt", false)
	log.Println("Part 1: Ways to win:", win1)

	win2 := WaysToWin("input1.txt", true)
	log.Println("Part 2: Ways to win:", win2)
}

func WaysToWin(filename string, combine bool) int {
	paper := parse(filename)

	if combine {
		paper = merge(paper)
	}

	var rs []int
	var ways int
	var idx int

	for idx < len(paper.Time) {
		r := result(paper.Time[idx], paper.Distance[idx])
		rs = append(rs, len(r))
		idx++
	}

	for _, v := range rs {
		if ways == 0 {
			ways = v
			continue
		}
		ways *= v
	}

	return ways
}

func result(time, distance int) []int {
	var res []int
	var i int

	for i < time {
		d := i * (time - i)
		if d > distance {
			res = append(res, i)
		}
		i++
	}

	return res
}

func parse(filename string) Paper {
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	parser := participle.MustBuild[Paper]()
	paper, err := parser.Parse(filename, fh)
	if err != nil {
		log.Fatal("unable to parse paper: %w", err)
	}

	return *paper
}

func merge(paper Paper) Paper {
	return Paper{
		Time:     []int{intsToInt(paper.Time)},
		Distance: []int{intsToInt(paper.Distance)},
	}
}

func intsToInt(ints []int) int {
	txt := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ints)), ""), "[]")

	val, err := strconv.Atoi(txt)
	if err != nil {
		log.Fatal("unable to convert to int: %w", err)
	}

	return val
}
