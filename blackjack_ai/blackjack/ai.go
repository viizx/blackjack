package blackjack

import (
	"deck"
	"fmt"
)

type AI interface {
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hands [][]deck.Card, dealersHand []deck.Card)
	Bet(shuffled bool) int
}

type dealerAI struct{}

type humanAI struct{}

func HumanAI() AI {
	return humanAI{}
}

func (AI dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && SoftHand(hand...)) {
		return MoveHit
	}
	return MoveStand
}
func (AI dealerAI) Results(hands [][]deck.Card, dealersHand []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("==FINAL HANDS==")
	for _, h := range hands {
		fmt.Println(" ", h)
	}
	fmt.Println("Player: ", hands)
	fmt.Println("Dealer: ", dealersHand)

}
func (AI dealerAI) Bet(shuffled bool) int {
	//no op
	return 1
}

func (AI humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	var input string
	fmt.Println("Player: ", hand)
	fmt.Println("Dealer: ", dealer)
	fmt.Println("(H)it, (S)tand, (D)ouble or s(P)lit?")

	_, err := fmt.Scanf("%s\n", &input)

	for err != nil {
		fmt.Println("Invalid input. Use 'h' to hit or 's' for stand.")
		fmt.Println("(H)it, (S)tand, (D)ouble or s(P)lit?")
		_, err = fmt.Scanf("%s\n", &input)
	}

	switch input {
	case "h":
		return MoveHit
	case "s":
		return MoveStand
	case "d":
		return MoveDouble
	case "p":
		return MoveSplit
	default:
		panic("Invalid option")
	}
}

func (AI humanAI) Results(hands [][]deck.Card, dealersHand []deck.Card) {

	fmt.Println("** FINAL RESULTS **")
	fmt.Println("Player: ")
	for _, h := range hands {
		fmt.Println(h)
	}
	fmt.Println("Dealer: ", dealersHand)

}

func (AI humanAI) Bet(shuffled bool) int {
	var bet int

	fmt.Print("\nHow much do you want to bet?\n")

	_, err := fmt.Scanf("%d\n", &bet)

	for err != nil {
		fmt.Println("Invalid bet. Please enter a number.")
		fmt.Print("\nHow much do you want to bet?\n")
		_, err = fmt.Scanf("%d\n", &bet)
	}

	return bet
}
