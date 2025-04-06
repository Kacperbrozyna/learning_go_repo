package blackjack

import (
	"fmt"

	deck "github.com/kacperbrozyna/learning_go_repo/deck_of_cards"
)

type AI interface {
	Results(hand [][]deck.Card, dealer []deck.Card)
	Play(hand []deck.Card, dealer deck.Card) Move
	Bet(shuffled bool) int
}

type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dealer_score := Score(hand...)
	if dealer_score <= 16 || dealer_score == 17 && Soft(hand...) {
		return MoveHit
	}

	return MoveStand
}

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {}

type humanAI struct{}

func HumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled")
	}

	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		var input string
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("Hit, Stand, Double or Split?")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "Hit":
			return MoveHit
		case "Stand":
			return MoveStand
		case "Double":
			return MoveDouble
		case "Split":
			return MoveSplit
		default:
			fmt.Println("Not a valid option")
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("--FINAL HANDS--")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}
