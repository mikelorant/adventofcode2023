package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Hands []Hand

type Hand struct {
	RawCards string
	Cards    map[string]int
	Bid      int
	Type     Type
	Strength []int
	Joker    bool
}

type Cards map[string]int

type Type int

func (t Type) String() string {
	switch t {
	case nothing:
		return "nothing"
	case highCard:
		return "highcard"
	case onePair:
		return "onepair"
	case twoPair:
		return "twopair"
	case threeKind:
		return "threekind"
	case fullHouse:
		return "fullhouse"
	case fourKind:
		return "fourkind"
	case fiveKind:
		return "fivekind"
	default:
		return "nothing"
	}
}

const (
	nothing   Type = iota // 0
	highCard              // 1
	onePair               // 2
	twoPair               // 3
	threeKind             // 4
	fullHouse             // 5
	fourKind              // 6
	fiveKind              // 7
)

func main() {
	total1 := TotalWinnings("input1.txt", false)
	log.Println("Part 1: Total winnings:", total1)

	total2 := TotalWinnings("input1.txt", true)
	log.Println("Part 2: Total winnings:", total2)
}

func TotalWinnings(filename string, joker bool) int {
	hands := parse(filename, joker)
	hands.Sort()

	return hands.Winnings()
}

func (h Hands) Len() int {
	return len(h)
}

func (h Hands) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h Hands) Less(i, j int) bool {
	if h[i].Type < h[j].Type {
		return true
	}

	if h[i].Type > h[j].Type {
		return false
	}

	for idx := range h[i].Strength {
		if h[i].Strength[idx] < h[j].Strength[idx] {
			return true
		}

		if h[i].Strength[idx] > h[j].Strength[idx] {
			return false
		}
	}

	return false
}

func (h *Hands) Sort() {
	sort.Sort(h)
}

func (h Hands) Winnings() int {
	var wins int

	for i := 0; i < len(h); i++ {
		wins += h[i].Bid * (i + 1)
	}

	return wins
}

func (h *Hand) Result() {
	j := h.Jokers()

	switch c := h.Cards; {
	case isFiveKind(c, j):
		h.Type = fiveKind
	case isFourKind(c, j):
		h.Type = fourKind
	case isFullHouse(c, j):
		h.Type = fullHouse
	case isThreeKind(c, j):
		h.Type = threeKind
	case isTwoPair(c, j):
		h.Type = twoPair
	case isOnePair(c, j):
		h.Type = onePair
	default:
		h.Type = highCard
	}
}

func (h *Hand) Score() {
	var val int

	cards := strings.Split(h.RawCards, "")

	for _, card := range cards {
		switch c := string(card); {
		case c == "A":
			val = 14
		case c == "K":
			val = 13
		case c == "Q":
			val = 12
		case c == "J":
			switch h.Joker {
			case true:
				val = 1
			default:
				val = 11
			}
		case c == "T":
			val = 10
		default:
			val = toInt(c)
		}

		h.Strength = append(h.Strength, val)
	}
}

func (h *Hand) Jokers() int {
	if h.Joker {
		_, ok := h.Cards["J"]
		if ok {
			return h.Cards["J"]
		}
	}

	return 0
}

func isOnePair(c Cards, j int) bool {
	if j > 0 {
		return true
	}

	for _, v := range c {
		if v == 2 {
			return true
		}
	}

	return false
}

func isTwoPair(c Cards, j int) bool {
	var res int

	switch {
	case j == 1:
		return isOnePair(deleteJoker(c), j-1)
	case j >= 2:
		return true
	}

	for _, v := range c {
		if v == 2 {
			res++
		}
	}

	return res == 2
}

func isThreeKind(c Cards, j int) bool {
	switch {
	case j == 1:
		return isOnePair(deleteJoker(c), j-1)
	case j >= 2:
		return true
	}

	for _, v := range c {
		if v == 3 {
			return true
		}
	}

	return false
}

func isFullHouse(c Cards, j int) bool {
	var three, pair bool

	switch {
	case j == 1:
		if isThreeKind(deleteJoker(c), j-1) {
			return true
		}

		if isTwoPair(deleteJoker(c), j-1) {
			return true
		}
	case j == 2:
		if isThreeKind(deleteJoker(c), j-2) {
			return true
		}

		if isOnePair(deleteJoker(c), j-2) {
			return true
		}
	case j >= 3:
		return true
	}

	for _, v := range c {
		if v == 3 {
			three = true
		}

		if v == 2 {
			pair = true
		}
	}

	return three && pair
}

func isFourKind(c Cards, j int) bool {
	switch {
	case j == 1:
		if isThreeKind(deleteJoker(c), j-1) {
			return true
		}
	case j == 2:
		if isOnePair(deleteJoker(c), j-2) {
			return true
		}
	case j >= 3:
		return true
	}

	for _, v := range c {
		if v == 4 {
			return true
		}
	}

	return false
}

func isFiveKind(c Cards, j int) bool {
	switch {
	case j == 1:
		if isFourKind(deleteJoker(c), j-1) {
			return true
		}
	case j == 2:
		if isThreeKind(deleteJoker(c), j-2) {
			return true
		}
	case j == 3:
		if isOnePair(deleteJoker(c), j-3) {
			return true
		}
	case j >= 4:
		return true
	}

	return len(c) == 1
}

func parse(filename string, joker bool) Hands {
	var hands Hands

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to read file: %w", err)
	}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		arr := strings.Fields(scanner.Text())
		hand := Hand{
			Cards:    make(map[string]int),
			RawCards: arr[0],
			Bid:      toInt(arr[1]),
			Joker:    joker,
		}
		hand.Check()
		hand.Result()
		hand.Score()

		hands = append(hands, hand)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}

	return hands
}

func toInt(txt string) int {
	num, err := strconv.Atoi(txt)
	if err != nil {
		log.Fatal("unable to convert to int: %w", err)
	}

	return num
}

func (h *Hand) Check() {
	cards := strings.Split(h.RawCards, "")
	for _, card := range cards {
		h.Cards[card]++
	}
}

func deleteJoker(c Cards) Cards {
	delete(c, "J")

	return c
}
