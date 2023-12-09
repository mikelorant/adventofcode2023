package main

import (
	"log"
	"os"
	"slices"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Report struct {
	Histories []History `@@+`
}

type History struct {
	Values []int `@Int* EOL`
}

func main() {
	valuesEnd := SumExtrapolatedValues("input1.txt", false)
	log.Println("Part 1: Sum of ending extrapolated values:", valuesEnd)

	valuesBegin := SumExtrapolatedValues("input1.txt", true)
	log.Println("Part 2: Sum of beginning extrapolated values:", valuesBegin)
}

func SumExtrapolatedValues(filename string, begin bool) int {
	var sum int

	report := parse(filename)

	if begin == true {
		for _, history := range report.Histories {
			slices.Reverse(history.Values)
		}
	}

	for _, history := range report.Histories {
		sum += history.Predict()
	}

	return sum
}

func (h *History) Predict() int {
	var next int

	seqs := sequences(h.Values)
	slices.Reverse(seqs)

	for idx, seq := range seqs {
		if idx == 0 {
			seqs[idx] = append(seqs[idx], 0)
			continue
		}

		prev := idx - 1
		diff := last(seqs[prev])
		next = last(seq) + diff

		seqs[idx] = append(seqs[idx], next)
	}

	return next
}

func sequences(values []int) [][]int {
	var seqs [][]int

	seqs = append(seqs, values)
	seq := seqs[0]

	for !isZero(seq) {
		seq = difference(seq)
		seqs = append(seqs, seq)
	}

	return seqs
}

func difference(values []int) []int {
	if len(values) <= 1 {
		log.Fatal("unable to diff values")
	}

	d := make([]int, len(values)-1)
	for idx := range d {
		d[idx] = values[idx+1] - values[idx]
	}

	return d
}

func isZero(values []int) bool {
	if len(values) == 0 {
		return false
	}

	return slices.Max(values) == 0 && slices.Min(values) == 0
}

func parse(filename string) Report {
	reportLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Int", `[\d-]+`},
		{"EOL", `\n`},
		{"whitespace", `\s+`},
	})

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	parser := participle.MustBuild[Report](
		participle.Lexer(reportLexer),
	)
	report, err := parser.Parse(filename, fh)
	if err != nil {
		log.Fatal("unable to parse paper: %w", err)
	}

	return *report
}

func last(ints []int) int {
	return ints[len(ints)-1]
}
