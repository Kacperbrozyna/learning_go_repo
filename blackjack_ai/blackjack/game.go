package blackjack

import (
	"errors"
	"fmt"

	deck "github.com/kacperbrozyna/learning_go_repo/deck_of_cards"
)

type state uint8

const (
	stateBetting state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}
type Game struct {
	numDecks        int
	numHands        int
	deck            []deck.Card
	state           state
	player          []hand
	playerBet       int
	handIdx         int
	dealer          []deck.Card
	dealerAI        AI
	balance         int
	blackjackPayout float64
}

func New(options Options) Game {
	game := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}

	if options.Decks == 0 {
		game.numDecks = 3
	}

	if options.Hands == 0 {
		game.numHands = 100
	}

	if options.BlackjackPayout == 0 {
		game.blackjackPayout = 1.5
	}

	game.numDecks = options.Decks
	game.numHands = options.Hands
	game.blackjackPayout = options.BlackjackPayout

	return game
}

func (game *Game) currentHand() *[]deck.Card {
	switch game.state {
	case statePlayerTurn:
		return &game.player[game.handIdx].cards
	case stateDealerTurn:
		return &game.dealer
	default:
		panic("Not any players turn")
	}
}

type hand struct {
	cards []deck.Card
	bet   int
}

func bet(game *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)

	if bet < 100 {
		panic("Not allowing bets below 100")
	}

	game.playerBet = bet
}

func deal(game *Game) {
	game.handIdx = 0
	playerHand := make([]deck.Card, 0, 5)
	game.dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, game.deck = draw(game.deck)
		playerHand = append(playerHand, card)

		card, game.deck = draw(game.deck)
		game.dealer = append(game.dealer, card)
	}

	game.player = []hand{
		{
			cards: playerHand,
			bet:   game.playerBet,
		},
	}
	game.state = statePlayerTurn
}

func (game *Game) Play(ai AI) int {
	game.deck = nil
	min := 52 * game.numDecks / 3

	for i := 0; i < game.numHands; i++ {
		shuffled := false
		if len(game.deck) < min {
			game.deck = deck.New(deck.Deck(game.numDecks), deck.Shuffle)
			shuffled = true
		}

		bet(game, ai, shuffled)

		deal(game)

		if Blackjack(game.dealer...) {
			endRound(game, ai)
			continue
		}

		for game.state == statePlayerTurn {
			hand := make([]deck.Card, len(*game.currentHand()))
			copy(hand, *game.currentHand())

			move := ai.Play(hand, game.dealer[0])
			err := move(game)

			switch err {
			case errBust:
				MoveStand(game)
			case errNotRecognisedState:
				panic(err)
			case errMorethanTwoSplit:
			case errMoreThanTwoDouble:
			case errNotTwoOfTheSame:
			case nil:
			default:
				panic(err)
			}
		}

		for game.state == stateDealerTurn {
			hand := make([]deck.Card, len(game.player))
			copy(hand, game.dealer)

			move := game.dealerAI.Play(hand, game.dealer[0])
			move(game)
		}

		endRound(game, ai)
	}

	return game.balance
}

var (
	errBust               = errors.New("Hand score exceeded 21")
	errNotRecognisedState = errors.New("Unexpected State")
	errMoreThanTwoDouble  = errors.New("Can only double on a hand with 2 cards")
	errMorethanTwoSplit   = errors.New("Can only split on a hand with 2 cards")
	errNotTwoOfTheSame    = errors.New("Can only split with two of the same rank")
)

type Move func(*Game) error

func MoveHit(game *Game) error {
	hand := game.currentHand()

	var card deck.Card
	card, game.deck = draw(game.deck)
	*hand = append(*hand, card)

	if Score(*hand...) > 21 {
		MoveStand(game)
		return errBust
	}

	return nil
}

func MoveSplit(game *Game) error {

	if game.state != statePlayerTurn {
		return errNotRecognisedState
	}

	cards := game.currentHand()
	if (len(*cards)) != 2 {
		return errMorethanTwoSplit
	}

	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errNotTwoOfTheSame
	}

	game.player = append(game.player, hand{
		cards: []deck.Card{(*cards)[1]},
		bet:   game.player[game.handIdx].bet,
	})

	game.player[game.handIdx].cards = (*cards)[:1]

	return nil
}

func MoveDouble(game *Game) error {
	if len(*game.currentHand()) != 2 {
		return errMoreThanTwoDouble
	}

	game.playerBet *= 2

	MoveHit(game)
	return MoveStand(game)
}

func MoveStand(game *Game) error {
	if game.state == stateDealerTurn {
		game.state++
		return nil
	}

	if game.state == statePlayerTurn {
		game.handIdx++
		if game.handIdx >= len(game.player) {
			game.state++
		}

		return nil
	}

	return errNotRecognisedState
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, card := range hand {
		score += min(int(card.Rank), 10)
	}

	return score
}

func endRound(game *Game, ai AI) {
	dealer_score := Score(game.dealer...)
	dealerBlackjack := Blackjack(game.dealer...)
	allHands := make([][]deck.Card, len(game.player))

	for i, hand := range game.player {
		cards := hand.cards
		allHands[i] = cards
		winnings := hand.bet

		player_score := Score(cards...)
		playerBlackJack := Blackjack(cards...)

		switch {
		case dealerBlackjack && playerBlackJack:
			fmt.Println("Double Blackjack!")
			winnings = 0
		case dealerBlackjack:
			fmt.Println("Git gud, Dealer Blackjack!")
			winnings *= -1
		case playerBlackJack:
			fmt.Println("Lucky you, Blackjack!")
			winnings = int(float64(winnings) * game.blackjackPayout)
		case player_score > 21:
			fmt.Println("You're Bust!")
			winnings *= -1
		case dealer_score > 21:
			fmt.Println("Dealer is Bust!")
		case player_score > dealer_score:
			fmt.Println("You Won!")
		case dealer_score >= player_score:
			fmt.Println("You Lost!")
			winnings *= -1
		}
		game.balance += winnings
	}

	fmt.Println()
	ai.Results(allHands, game.dealer)
	game.player = nil
	game.dealer = nil
	fmt.Println()
}

func Score(hand ...deck.Card) int {
	min_score := minScore(hand...)
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

func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}
