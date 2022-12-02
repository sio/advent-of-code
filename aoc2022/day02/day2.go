package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

type GameMove int

const (
	Rock     GameMove = 0
	Paper    GameMove = 1
	Scissors GameMove = 2
)

type GameRound struct {
	Them GameMove
	Us   GameMove
}

func (g *GameRound) Outcome() int {
	result := (g.Us - g.Them + 3) % 3 // golang modulo operator is special
	switch result {
	case 0:
		return 3 // draw
	case 1:
		return 6 // we won
	case 2:
		return 0 // we lost
	default:
		log.Fatalf("unexpected result: %d (round %v)", result, g)
		return 0
	}
}

func (g *GameRound) Score() int {
	return 1 + int(g.Us) + g.Outcome()
}

func (g *GameRound) Valid() bool {
	return g.Us >= 0 && g.Us < 3 && g.Them >= 0 && g.Them < 3
}

func main() {
	part := flag.Int("part", 0, "puzzle part")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	switch *part {
	case 1:
		part1(flag.Args()[0])
	case 2:
		part2(flag.Args()[0])
	default:
		part1(flag.Args()[0])
		part2(flag.Args()[0])
	}
}

func part1(filename string) {
	opponentMoves := map[string]GameMove{
		"A": Rock,
		"B": Paper,
		"C": Scissors,
	}
	ourMoves := map[string]GameMove{
		"X": Rock,
		"Y": Paper,
		"Z": Scissors,
	}
	var moves []string
	var score int
	for line := range ReadLines(filename) {
		moves = strings.Split(line, " ")
		if len(moves) != 2 {
			log.Fatalf("unexpected input line: %s", line)
		}
		round := GameRound{
			Us:   ourMoves[moves[1]],
			Them: opponentMoves[moves[0]],
		}
		if !round.Valid() {
			log.Fatalf("invalid moves parsed: %v", round)
		}
		score += round.Score()
		//log.Printf("Round %v, score %d", round, round.Score())
	}
	log.Printf("Part 1 score: %d", score)
}

func part2(filename string) {
	opponentMoves := map[string]GameMove{
		"A": Rock,
		"B": Paper,
		"C": Scissors,
	}
	outcomes := map[string]int{
		"X": -1, // lose
		"Y": 0,  // draw
		"Z": 1,  // win
	}
	var moves []string
	var score int
	for line := range ReadLines(filename) {
		moves = strings.Split(line, " ")
		if len(moves) != 2 {
			log.Fatalf("unexpected input line: %s", line)
		}
		round := GameRound{
			Them: opponentMoves[moves[0]],
		}
		delta, ok := outcomes[moves[1]]
		if !ok {
			log.Fatalf("invalid input: %s", line)
		}
		round.Us = GameMove((int(round.Them) + delta + 3) % 3)
		if !round.Valid() {
			log.Fatalf("invalid round: %v (from line %q)", round, line)
		}
		score += round.Score()
		//log.Printf("Round %v, score %d", round, round.Score())
	}
	log.Printf("Part 2 score: %d", score)
}
