package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type GameMove uint8

const (
	Rock     GameMove = 1
	Paper    GameMove = 2
	Scissors GameMove = 3
)

type GameRound struct {
	Them GameMove
	Us   GameMove
}

func (g *GameRound) Outcome() uint8 {
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
	return int(uint8(g.Us) + g.Outcome())
}

func (g *GameRound) Valid() bool {
	return g.Us > 0 && g.Us < 4 && g.Them > 0 && g.Them < 4
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s INPUT\n", os.Args[0])
		os.Exit(1)
	}
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
	for line := range ReadLines(os.Args[1]) {
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
	log.Printf("Total score: %d", score)
}
