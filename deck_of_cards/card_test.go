package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Rank: King, Suit: Heart})
	fmt.Println(Card{Suit: Joker})

	//Output:
	//Ace of Hearts
	//Two of Spades
	//Nine of Diamonds
	//Jack of Clubs
	//King of Hearts
	//Joker
}

func TestNew(t *testing.T) {
	cards := New()

	if len(cards) != 52 {
		t.Error("Wrong number of cards in a new deck")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Error("Expected Ace of Spades as first card recevied: ", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Error("Expected Ace of Spades as first card recevied: ", cards[0])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0

	for _, card := range cards {
		if card.Suit == Joker {
			count++
		}
	}

	if count != 3 {
		t.Error("Expected 3 Jokers, received: ", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}

	cards := New(Filter(filter))
	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Error("Expected to filter out all Twos and Threes: ", card.Rank.String())
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))

	if len(cards) != 156 {
		t.Errorf("Expected %d cards, received %d cards: ", 156, len(cards))
	}
}

func TestShuffle(t *testing.T) {
	shuffleRand = rand.New(rand.NewSource(0))

	original := New()
	first := original[40]
	second := original[35]
	cards := New(Shuffle)

	if cards[0] != first {
		t.Errorf("Expected the first card to be %s, received %s", first, cards[0])
	}

	if cards[1] != second {
		t.Errorf("Expected the second card to be %s, received %s", second, cards[1])
	}
}
