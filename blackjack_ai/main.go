package main

import (
	"blackjack"
	"fmt"
)

func main() {
	options := blackjack.Options{
		Decks:           3,
		Hands:           2,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(options)

	winnings := game.Play(blackjack.HumanAI())

	fmt.Println(winnings)
}
