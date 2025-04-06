package main

import (
	"fmt"

	"github.com/Kacperbrozyna/learning_go_repo/blackjack_ai/blackjack"
	deck "github.com/kacperbrozyna/learning_go_repo/deck_of_cards"
)

type basicAI struct {
	score int
	seen  int
	decks int
}

func (ai *basicAI) Bet(shuffled bool) int {

	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	true_score := ai.score / ((ai.decks*52 - ai.seen) / 52)

	switch {
	case true_score > 14:
		return 10000
	case true_score > 8:
		return 500
	default:
		return 100
	}
}

// Needs a better strat, loses bad
func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)

	if len(hand) == 2 {

		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}

		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}

	dealer_score := blackjack.Score(dealer)
	if dealer_score >= 5 && dealer_score <= 6 {
		return blackjack.MoveStand
	}

	if dealer_score < 13 {
		return blackjack.MoveHit
	}

	return blackjack.MoveStand
}

func (ai *basicAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}

	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *basicAI) count(card deck.Card) {
	score := blackjack.Score(card)

	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	default:
	}

	ai.seen++
}

func main() {
	options := blackjack.Options{
		Decks:           4,
		Hands:           50000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(options)
	winnings := game.Play(&basicAI{
		seen:  0,
		score: 0,
		decks: 4,
	}) //Replace with blackjack.HumanAI() if you would like to play
	fmt.Println(winnings)

}
