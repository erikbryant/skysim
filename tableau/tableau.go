package tableau

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
)

const width = 4
const height = 3

type Tableau [height][width]int

func Expected() int {
	return cards.Average() * width * height
}

func Deal(deck cards.Cards) Tableau {
	var t Tableau

	for row := range t {
		for col := range t[row] {
			t[row][col] = deck.Draw()
		}
	}

	return t
}

func (t Tableau) Score() int {
	score := 0

	for row := range t {
		for col := range t[row] {
			score += t[row][col]
		}
	}

	return score
}

func (t Tableau) Print() {
	for row := range t {
		for col := range t[row] {
			fmt.Printf("%2d ", t[row][col])
		}
		fmt.Println()
	}

}
