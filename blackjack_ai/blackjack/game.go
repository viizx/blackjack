package blackjack

import (
	"deck"
	"errors"
	"fmt"
)

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

var errBust = errors.New("hand score exceeded 21")

type Move func(*Game) error
type state int8
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

type Game struct {
	nDecks          int
	nHands          int
	blackJackPayout float64

	state state
	deck  []deck.Card

	player    []hand
	handIdx   int
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI dealerAI
}

type hand struct {
	cards []deck.Card
	bet   int
}

func deal(g *Game) {
	playerHand := make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)

	var card deck.Card

	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{
		{
			cards: playerHand,
			bet:   g.playerBet,
		},
	}
	g.state = statePlayerTurn

}

func New(options Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}

	if options.Decks == 0 {
		options.Decks = 3
	}

	if options.Hands == 0 {
		options.Hands = 10
	}

	if options.BlackjackPayout == 0 {
		options.BlackjackPayout = 1.5
	}

	g.nDecks = options.Decks
	g.nHands = options.Hands
	g.blackJackPayout = options.BlackjackPayout

	return g

}

func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3

	for i := 0; i < g.nDecks; i++ {
		deckIsShuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			deckIsShuffled = true
		}
		bet(g, ai, deckIsShuffled)
		deal(g)

		if Blackjack(g.dealer...) {
			endRound(g, ai)
			continue
		}

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				// noop
			default:
				panic(err)
			}
		}

		endRound(g, ai)
	}
	return g.balance
}

func MoveHit(g *Game) error {
	hand := g.currentHand()

	var card deck.Card

	card, g.deck = draw(g.deck)

	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

func MoveDouble(g *Game) error {
	if len(g.player) != 2 {
		return errors.New("you can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}

	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("invalid state")
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)

	if minScore > 11 {
		return minScore
	}

	for _, card := range hand {
		if card.Rank == deck.Ace {
			return minScore + 10
		}
	}

	return minScore
}

func SoftHand(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, card := range hand {
		score += min(int(card.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func endRound(g *Game, ai AI) {
	dScore := Score(g.dealer...)
	dBlackjack := Blackjack(g.dealer...)
	allHands := make([][]deck.Card, len(g.player))
	for idx, hand := range g.player {
		cards := hand.cards
		allHands[idx] = cards
		pScore, pBlackjack := Score(cards...), Blackjack(cards...)
		winnings := hand.bet

		switch {
		case pBlackjack && dBlackjack:
			winnings = 0
		case dBlackjack:
			winnings *= -1
		case pBlackjack:
			winnings = int(float64(winnings) * g.blackJackPayout)
		case pScore > 21:
			g.balance--
			winnings *= -1
		case dScore > 21:
		case pScore > dScore:
			//win
		case pScore < dScore:
			//win
			g.balance--
			winnings *= -1
		case pScore == dScore:
			fmt.Println("Draw")
			winnings = 0
		}
		g.balance += winnings
	}

	fmt.Println("\nbalance", g.balance)
	ai.Results(allHands, g.dealer)
	g.player = nil
	g.dealer = nil

}

func (g *Game) currentHand() *[]deck.Card {

	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("No players turn?")
	}
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}
