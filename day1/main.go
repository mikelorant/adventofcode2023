package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Index struct {
	Number int
	Index  int
}

type Math string

const (
	FILENAME = "input2.txt"

	MIN Math = "<"
	MAX Math = ">"
)

func main() {
	fh, err := os.Open(FILENAME)
	if err != nil {
		log.Fatal("unable to open file: %w", err)
	}

	var sum int

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		idx := indexNumbersWords(scanner.Text())

		min := value(idx, MIN)
		max := value(idx, MAX)

		cali := caliValue(min, max)
		sum += cali

		log.Println(scanner.Text(), value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}

	log.Println("Sum of calibration values:", sum)
}

func caliValue(first, last int) int {
	val, err := strconv.Atoi(fmt.Sprintf("%d%d", first, last))
	if err != nil {
		log.Fatal("unable to convert to string: %w", err)
	}

	return val
}

func indexNumbersWords(txt string) []Index {
	var idx []Index

	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	words := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	idx = append(idx, index(txt, numbers)...)
	idx = append(idx, index(txt, words)...)

	return idx
}

func index(txt string, arr []string) []Index {
	var idx []Index

	for k, v := range arr {
		re := regexp.MustCompile(v)
		for _, vv := range re.FindAllStringIndex(txt, -1) {
			idx = append(idx, Index{
				Number: k + 1,
				Index:  vv[0],
			})
		}
	}

	return idx
}

func value(idx []Index, m Math) int {
	var res int

	i := -1

	for _, v := range idx {
		if i == -1 || update(v.Index, i, m) {
			res = v.Number
			i = v.Index
		}
	}

	return res
}

func update(x, y int, m Math) bool {
	if m == MIN {
		return x < y
	}

	return x > y
}
