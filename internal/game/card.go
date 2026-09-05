package main

import (
	"errors"
	"math/rand"
	"strconv"
)

type Card struct {
	Value     string
	Type      string
	RealValue int
	Point     int
}

func (c Card) Beats(other Card) bool {
	return c.RealValue > other.RealValue
}

func NewDeck() []Card {
	deck := make([]Card, 54)
	names := [4]string{"Clubs", "Diamonds", "Hearts", "Spades"}

	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			index := j + 13*i
			deck[index].Type = names[i]
			deck[index].RealValue = j + 1

			switch j {
			case 0:
				deck[index].RealValue = 14
				deck[index].Value = "Ace"
				deck[index].Point = 11
			case 10:
				deck[index].Value = "Jack"
				deck[index].Point = 10
			case 11:
				deck[index].Value = "Queen"
				deck[index].Point = 10
			case 12:
				deck[index].Value = "King"
				deck[index].Point = 10
			default:
				deck[index].Value = strconv.Itoa(j + 1)
				deck[index].Point = j + 1
			}
		}
	}
	deck[52] = Card{"Joker", "Joker", 15, 20}
	deck[53] = Card{"Joker", "Joker", 15, 20}

	return deck
}

func Shuffle(deck []Card) {

	rand.Shuffle(len(deck), func(a, b int) { deck[a], deck[b] = deck[b], deck[a] })
}

type Player struct {
	Hand     []Card
	Captured []Card
}

func (p *Player) Score() int {
	score := 0
	for i := 0; i < len(p.Captured); i++ {
		score += p.Captured[i].Point
	}

	return score
}

type GameState struct {
	Players   [2]*Player
	Pile      []Card
	Remaining []Card
	Turn      bool
	Over      bool
	LastGuy   bool
}

func (g *GameState) PlayCard(playerIdx int, cardIdx int) (captured bool, err error) {

	if (playerIdx == 0 && g.Turn) || (playerIdx == 1 && !g.Turn) {
		return false, errors.New("Not your turn!")
	}

	if len(g.Pile) == 0 {
		g.Pile = append(g.Pile, g.Players[playerIdx].Hand[cardIdx])
	} else if (g.Players[playerIdx].Hand[cardIdx]).Beats(g.Pile[len(g.Pile)-1]) {
		g.LastGuy = g.Turn
		captured = true
		for i := 0; i < len(g.Pile); i++ {
			g.Players[playerIdx].Captured = append(g.Players[playerIdx].Captured, g.Pile[i])
		}
		g.Players[playerIdx].Captured = append(g.Players[playerIdx].Captured, g.Players[playerIdx].Hand[cardIdx])
		g.Pile = g.Pile[:0]
	} else {
		g.Pile = append(g.Pile, g.Players[playerIdx].Hand[cardIdx])
	}

	for i := cardIdx; i < len(g.Players[playerIdx].Hand)-1; i++ {
		g.Players[playerIdx].Hand[i] = g.Players[playerIdx].Hand[i+1]
	}
	g.Players[playerIdx].Hand = g.Players[playerIdx].Hand[:len(g.Players[playerIdx].Hand)-1]

	g.redealIfNeeded()

	if (len(g.Players[0].Hand) == 0) && (len(g.Players[1].Hand) == 0) {
		x := 0
		if g.LastGuy {
			x = 1
		}

		for i := 0; i < len(g.Pile); i++ {
			g.Players[x].Captured = append(g.Players[x].Captured, g.Pile[i])
		}

		g.Pile = g.Pile[:0]

		g.Over = true
	}

	g.Turn = !(g.Turn)

	return captured, err
}

func (g *GameState) Winner() int {
	if g.Players[0].Score() > g.Players[1].Score() {
		return 0
	} else if g.Players[0].Score() < g.Players[1].Score() {
		return 1
	} else {
		return -1
	}
}

func NewGame(deck []Card) *GameState {
	Shuffle(deck)

	player1_hand := make([]Card, 5)
	copy(player1_hand, deck[45:50])
	player2_hand := make([]Card, 5)
	copy(player2_hand, deck[40:45])
	players := [2]*Player{&Player{player1_hand, nil}, &Player{player2_hand, nil}}
	remaining := make([]Card, 40)
	copy(remaining, deck[:40])
	pile := make([]Card, 4)
	copy(pile, deck[50:54])
	turn := false
	over := false
	lastGuy := false
	game := &GameState{players, pile, remaining, turn, over, lastGuy}

	return game

}

func (g *GameState) redealIfNeeded() {
	if (len(g.Players[0].Hand) == 0) && (len(g.Players[1].Hand) == 0) && (len(g.Remaining) > 0) {
		player1s := make([]Card, 5)
		copy(player1s, g.Remaining[(len(g.Remaining)-5):])
		g.Players[0].Hand = player1s
		g.Remaining = g.Remaining[:(len(g.Remaining) - 5)]
		player2s := make([]Card, 5)
		copy(player2s, g.Remaining[(len(g.Remaining)-5):])
		g.Players[1].Hand = player2s
		g.Remaining = g.Remaining[:(len(g.Remaining) - 5)]

	}
}
