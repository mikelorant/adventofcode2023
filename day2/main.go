package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Games struct {
	ID    int
	Games []Game
}

type Game struct {
	Legal bool
	Red   int
	Green int
	Blue  int
}

var GameMax = Game{
	Red:   12,
	Green: 13,
	Blue:  14,
}

const (
	FILENAME = "input1.txt"
)

func main() {
	fh, err := os.Open(FILENAME)
	if err != nil {
		log.Fatal("unable to open file: %w", err)
	}

	var all []Games

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		all = append(all, parse(scanner.Text()))
	}

	sum := sumLegal(all)
	power := sumPower(all)

	log.Println("Sum of legal games:", sum)
	log.Println("Power of all games:", power)

	if err := scanner.Err(); err != nil {
		log.Fatal("scanner error: %w", err)
	}
}

func sumLegal(all []Games) int {
	var sum int

	for _, v := range all {
		legal := true

		for _, game := range v.Games {
			if !game.Legal {
				legal = false
			}
		}

		if legal {
			sum += v.ID
		}
	}

	return sum
}

func sumPower(all []Games) int {
	var power int

	for _, v := range all {
		var gameMin Game

		for _, game := range v.Games {
			if game.Red > gameMin.Red {
				gameMin.Red = game.Red
			}

			if game.Green > gameMin.Green {
				gameMin.Green = game.Green
			}

			if game.Blue > gameMin.Blue {
				gameMin.Blue = game.Blue
			}
		}

		power += gameMin.Red * gameMin.Green * gameMin.Blue
	}

	return power
}

func parse(txt string) Games {
	res := strings.SplitN(txt, ":", 2)

	return Games{
		ID:    parseID(res[0]),
		Games: parseGames(res[1]),
	}
}

func parseID(txt string) int {
	re := regexp.MustCompile(`Game\ (.*)`)

	matches := re.FindStringSubmatch(txt)

	i, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal("unable to convert to string: %w", err)
	}

	return i
}

func parseGames(txt string) []Game {
	var games []Game

	gs := strings.Split(txt, ";")

	for k, v := range gs {
		gs[k] = strings.TrimSpace(v)
	}

	for _, g := range gs {
		games = append(games, parseGame(g))
	}

	return games
}

func parseGame(txt string) Game {
	var game Game

	res := strings.Split(txt, ",")

	for _, v := range res {
		re := regexp.MustCompile(`^(\d*)\ (\w*)`)
		matches := re.FindStringSubmatch(strings.TrimSpace(v))

		num, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal("unable to convert to string: %w", err)
		}

		colour := strings.ToLower(matches[2])

		switch colour {
		case "red":
			game.Red += num
		case "green":
			game.Green += num
		case "blue":
			game.Blue += num
		}
	}

	game.Legal = legal(game, GameMax)

	return game
}

func legal(game, max Game) bool {
	if game.Red > max.Red {
		return false
	}

	if game.Green > max.Green {
		return false
	}

	if game.Blue > max.Blue {
		return false
	}

	return true
}
