package main

import (
	"blackjack"
	"deck"
	"fmt"
)

type BasicAI struct {
	score int
	seen  int
	decks int
}

func (ai *BasicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}
		if (score == 10 || score == 11) && !blackjack.SoftHand(hand...) {
			return blackjack.MoveDouble
		}
	}
	dScore := blackjack.Score(dealer)
	if dScore >= 5 && dScore <= 6 {
		return blackjack.MoveStand
	}
	if score < 13 {
		return blackjack.MoveHit
	}
	return blackjack.MoveStand
}

func (ai *BasicAI) Results(hands [][]deck.Card, dealersHand []deck.Card) {

	for _, card := range dealersHand {
		ai.count(card)
	}
	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}
func (ai *BasicAI) count(card deck.Card) {
	score := blackjack.Score(card)
	if score >= 10 {
		ai.score--
	} else if score <= 6 {
		ai.score++
	}
	ai.seen++
}
func (ai *BasicAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	trueScore := ai.score / ((ai.decks*52 - ai.seen) / 52)
	switch {
	case trueScore >= 10:
		fmt.Printf("trueScore>14\n")
		return 10000
	case trueScore >= 8:
		return 1000
	default:
		return 100
	}

}

func main() {
	options := blackjack.Options{
		Decks:           4,
		Hands:           9999999,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(options)

	winnings := game.Play(&BasicAI{
		seen:  0,
		score: 0,
		decks: 4,
	})

	fmt.Println(winnings)
}
