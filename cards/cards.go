package cards

import (
	"fmt"
	"math/rand"
	"time"
)

var distribution = [][2]int{
	[2]int{5, -2},
	[2]int{10, -1},
	[2]int{15, 0},
	[2]int{10, 1},
	[2]int{10, 2},
	[2]int{10, 3},
	[2]int{10, 4},
	[2]int{10, 5},
	[2]int{10, 6},
	[2]int{10, 7},
	[2]int{10, 8},
	[2]int{10, 9},
	[2]int{10, 10},
	[2]int{10, 11},
	[2]int{10, 12},
}

type Cards []int

func Average() int {
	deck := Deck()
	sum := 0
	count := 0
	for _, val := range deck {
		count++
		sum += val
	}
	return sum / count
}

func Deck() Cards {
	deck := []int{}

	for _, val := range distribution {
		for i := 1; i <= val[0]; i++ {
			deck = append(deck, val[1])
		}
	}

	return deck
}

func (c Cards) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c), func(i, j int) { c[i], c[j] = c[j], c[i] })
}

func (c *Cards) Draw() int {
	card := (*c)[0]

	*c = (*c)[1:]

	return card
}

func (c *Cards) Print() {
	for _, val := range *c {
		fmt.Print(val)
		fmt.Print(" ")
	}
	fmt.Println()
}
