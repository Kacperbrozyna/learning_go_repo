//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)
const (
	minRank = Ace
	maxRank = King
)

type Card struct {
	Suit
	Rank
}

func (card Card) String() string {
	if card.Suit == Joker {
		return card.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", card.Rank.String(), card.Suit.String())
}

func New(options ...func([]Card) []Card) []Card {
	var cards []Card

	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range options {
		cards = opt(cards)
	}

	return cards
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absoulteRank(cards[i]) < absoulteRank(cards[j])
	}
}

func absoulteRank(card Card) int {
	return int(card.Suit)*int(maxRank) + int(card.Rank)
}

type Permer interface {
	Perm(numb int) []int16
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	permutation := shuffleRand.Perm(len(cards))

	for i, j := range permutation {
		ret[i] = cards[j]
	}

	return ret
}

func Jokers(numb int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < numb; i++ {
			cards = append(cards, Card{Rank: Rank(i), Suit: Joker})
		}

		return cards
	}
}

func Filter(filter func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card

		for _, card := range cards {
			if !filter(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

func Deck(numb int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card

		for i := 0; i < numb; i++ {
			ret = append(ret, cards...)
		}

		return ret
	}
}
