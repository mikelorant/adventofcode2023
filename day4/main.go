package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	ID      int
	Wins    int
	Score   int
	Numbers []int
	Winners []int
}

type Winners struct {
	ID   int
	Card []int
}

func main() {
	wn := SumWinningNumbers("input1.txt")
	log.Println("Part 1: Sum of winning numbers:", wn)

	tc := TotalWinningCards("input1.txt")
	log.Println("Part 2: Total cards:", tc)
}

func SumWinningNumbers(filename string) int {
	var sum int

	for _, card := range parseCards(filename) {
		sum += card.Score
	}

	return sum
}

func TotalWinningCards(filename string) int {
	var sum int

	cards := parseCards(filename)
	winners := make([]Winners, len(cards))

	for idx := range winners {
		winners[idx].ID = idx
		winners[idx].Card = append(winners[idx].Card, idx)
	}

	for idx, card := range cards {
		n := len(winners[idx].Card) - 1
		i := 0

		for i <= n {
			winners = updateWinners(winners, idx, card.Wins)

			i++
		}
	}

	for _, v := range winners {
		sum += len(v.Card)
	}

	return sum
}

func updateWinners(ws []Winners, idx, wins int) []Winners {
	i := 1

	for i <= wins {
		newIdx := idx + i

		ws[newIdx].Card = append(ws[newIdx].Card, idx)

		i++
	}

	return ws
}

func parseCards(filename string) []Card {
	var cards []Card

	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal("unable to open file: %w", err)
	}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		cards = append(cards, card(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}

	return cards
}

func card(txt string) Card {
	re := regexp.MustCompile(`Card\s+(\d+):(.*)`)
	matches := re.FindStringSubmatch(txt)
	id := toInt(matches[1])[0]

	arr := strings.SplitN(matches[2], "|", 2)
	wins := toInt(arr[0])
	nums := toInt(arr[1])

	c := Card{
		ID:      id,
		Numbers: nums,
		Winners: wins,
	}

	c.Wins = winners(c)
	c.Score = score(c.Wins)

	return c
}

func score(wins int) int {
	switch wins {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return power(2, wins-1)
	}
}

func toInt(txt string) []int {
	var nums []int

	for _, val := range strings.Fields(txt) {
		num, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			log.Fatal("unable to convert to int: %w", err)
		}

		nums = append(nums, num)
	}

	return nums
}

func winners(c Card) int {
	var wins int

	for _, win := range c.Winners {
		var res bool

		for _, num := range c.Numbers {
			if num == win {
				res = true

				break
			}
		}

		if res {
			wins++
		}
	}

	return wins
}

func power(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
