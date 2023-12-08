package main

import (
	"log"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Map struct {
	Instructions []Instruction
	Nodes        map[string]Node
	Steps        int
	Raw          Raw `@@`
}

type Raw struct {
	Instructions string `@Value`
	Nodes        []Node `@@+`
}

type Instruction int

type Node struct {
	Name  string `@Value Symbol Symbol`
	Left  string `@Value Symbol`
	Right string `@Value Symbol`
}

const (
	Left Instruction = iota
	Right
)

func main() {
	steps := Steps("input1.txt", "AAA", "ZZZ")
	log.Println("Part 1: Steps:", steps)

	stepsAZ := Steps("input2.txt", "A", "Z")
	log.Println("Part 2: Steps:", stepsAZ)
}

func Steps(filename string, start, end string) int {
	m := parse(filename)

	var steps []int
	for _, node := range m.SuffixNodes(start) {
		steps = append(steps, m.Traverse(node, end))
	}

	switch len(steps) {
	case 0:
	case 1:
		return steps[0]
	case 2:
		return lcm(steps[0], steps[1])
	}

	return lcm(steps[0], steps[1], steps[2:]...)
}

func (m *Map) Traverse(start, suffix string) int {
	var steps int

	at := start

	for !isSuffix(at, suffix) || steps == 0 {
		for _, instruct := range m.Instructions {
			at = m.step(at, instruct)
			steps++
		}

		if isSuffix(at, suffix) {
			break
		}
	}

	return steps
}

func (m *Map) SuffixNodes(suffix string) []string {
	var nodes []string

	for node := range m.Nodes {
		if strings.HasSuffix(node, suffix) {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (m *Map) ConvertInstructions() {
	for _, val := range strings.Split(m.Raw.Instructions, "") {
		switch val {
		case "L":
			m.Instructions = append(m.Instructions, Left)
		case "R":
			m.Instructions = append(m.Instructions, Right)
		}
	}
}

func (m *Map) ConvertNodes() {
	m.Nodes = make(map[string]Node)

	for _, val := range m.Raw.Nodes {
		m.Nodes[val.Name] = Node{
			Left:  val.Left,
			Right: val.Right,
		}
	}
}

func (m *Map) step(node string, instruct Instruction) string {
	switch instruct {
	case Left:
		return m.Nodes[node].Left
	default:
		return m.Nodes[node].Right
	}
}

func parse(filename string) Map {
	mapLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Value", `\w+`},
		{"Symbol", `[=(),]`},
		{"whitespace", `\s+`},
	})

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	parser := participle.MustBuild[Map](
		participle.Lexer(mapLexer),
	)
	m, err := parser.Parse(filename, fh)
	if err != nil {
		log.Fatal("unable to parse paper: %w", err)
	}

	m.ConvertInstructions()
	m.ConvertNodes()

	return *m
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}

func isSuffix(node string, suffix string) bool {
	return strings.HasSuffix(node, suffix)
}
