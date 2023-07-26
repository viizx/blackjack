package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Diamond})
	fmt.Println(Card{Rank: Three, Suit: Club})
	fmt.Println(Card{Rank: Jack, Suit: Spade})
	fmt.Println(Card{Suit: Joker})

	//Output:
	// Ace of Hearts
	// Two of Diamonds
	// Three of Clubs
	// Jack of Spades
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("Wrong Number of Cards in Deck")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	firstCard := Card{Rank: Ace, Suit: Spade}
	if cards[0] != firstCard {
		t.Error("Expected Ace of Spades to be the first card. Recieved:", firstCard)
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	firstCard := Card{Rank: Ace, Suit: Spade}
	if cards[0] != firstCard {
		t.Error("Expected Ace of Spades to be the first card. Recieved:", firstCard)
	}
}

func TestShuffle(t *testing.T) {
	cards := New(Sort(Less))
	shuffledCards := Shuffle(cards)
	isSame := true
	for i, card := range cards {
		if card != shuffledCards[i] {
			isSame = false
			break
		}
	}
	if isSame {
		t.Error("Deck is not shuffled")
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(4))
	count := 0
	for _, card := range cards {
		if card.Suit == Joker {
			count++
		}
	}

	if count != 4 {
		t.Error("Expected 4 Jokers. Recieved:", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}

	cards := New(Filter(filter))

	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Error("Expected all Twos and Threes to be filtered out")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))

	if len(cards) != 13*4*3 {
		t.Errorf("Expected %d cards, recieved %d cards", 13*4*3, len(cards))
	}
}
