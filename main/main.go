package main

import (
	"github.com/erikbryant/skysim/cards"
)

func main() {
	deck := cards.Deck()
	deck.Print()
	deck.Shuffle()
	deck.Print()
}
