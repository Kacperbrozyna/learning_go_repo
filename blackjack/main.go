package main

import (
	"fmt"
	"strings"

	deck "github.com/kacperbrozyna/learning_go_repo/deck_of_cards"
)

type Hand []deck.Card

type State uint8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (game_state *GameState) CurrentPlayer() *Hand {
	switch game_state.State {
	case StatePlayerTurn:
		return &game_state.Player
	case StateDealerTurn:
		return &game_state.Dealer
	default:
		panic("Not any players turn")
	}
}

func clone(game_state GameState) GameState {
	ret := GameState{
		Deck:   make(Hand, len(game_state.Deck)),
		State:  game_state.State,
		Player: make(Hand, len(game_state.Player)),
		Dealer: make(Hand, len(game_state.Dealer)),
	}

	copy(ret.Deck, game_state.Deck)
	copy(ret.Player, game_state.Player)
	copy(ret.Dealer, game_state.Dealer)
	return ret
}

func (hand Hand) String() string {
	strs := make([]string, len(hand))
	for i := range hand {
		strs[i] = hand[i].String()
	}

	return strings.Join(strs, ", ")
}

func (hand Hand) DealerString() string {
	return hand[0].String() + ", **HIDDEN**"
}

func (hand Hand) MinScore() int {
	score := 0

	for _, card := range hand {
		score += min(int(card.Rank), 10)
	}

	return score
}

func (hand Hand) Score() int {
	min_score := hand.MinScore()
	if min_score > 11 {
		return min_score
	}

	for _, card := range hand {
		if card.Rank == deck.Ace {
			return min_score + 10
		}
	}

	return min_score
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Shuffle(game_state GameState) GameState {
	ret := clone(game_state)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(game_state GameState) GameState {
	ret := clone(game_state)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)

		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}

	ret.State = StatePlayerTurn
	return ret
}

func Stand(game_state GameState) GameState {
	ret := clone(game_state)
	ret.State++
	return ret
}

func Hit(game_state GameState) GameState {
	ret := clone(game_state)
	hand := ret.CurrentPlayer()

	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)

	if hand.Score() > 21 {
		return Stand(game_state)
	}

	return ret
}

func EndHand(game_state GameState) GameState {
	ret := clone(game_state)

	player_score, dealer_score := ret.Player.Score(), ret.Dealer.Score()

	fmt.Println("--FINAL HANDS--")
	fmt.Println("Player:", ret.Player, "\nScore:", player_score)
	fmt.Println("Dealer:", ret.Dealer, "\nScore:", dealer_score)

	switch {
	case player_score > 21:
		fmt.Println("You're Bust!")
	case dealer_score > 21:
		fmt.Println("Dealer is Bust!")
	case player_score > dealer_score:
		fmt.Println("You Won!")
	case dealer_score >= player_score:
		fmt.Println("You Lost!")
	}

	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil

	return ret
}

func main() {

	var game_state GameState
	game_state = Shuffle(game_state)
	game_state = Deal(game_state)

	var input string
	for game_state.State == StatePlayerTurn {
		fmt.Println("Player:", game_state.Player)
		fmt.Println("Dealer:", game_state.Dealer.DealerString())
		fmt.Println("Hit or Stand?")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "Hit":
			game_state = Hit(game_state)
		case "Stand":
			game_state = Stand(game_state)
		default:
			fmt.Println("Not a valid option")
		}
	}

	for game_state.State == StateDealerTurn {
		if game_state.Dealer.Score() <= 16 || game_state.Dealer.Score() == 17 && game_state.Dealer.MinScore() != 17 {
			game_state = Hit(game_state)
		} else {
			game_state = Stand(game_state)
		}
	}

	game_state = EndHand(game_state)
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
