package cards

import (
	"fmt"
	"math/rand"
	"time"
)

// How many cards are there of each rank
type dist struct {
	rank  int
	count int
}

var distribution = []dist{
	dist{-2, 5},
	dist{-1, 10},
	dist{0, 15},
	dist{1, 10},
	dist{2, 10},
	dist{3, 10},
	dist{4, 10},
	dist{5, 10},
	dist{6, 10},
	dist{7, 10},
	dist{8, 10},
	dist{9, 10},
	dist{10, 10},
	dist{11, 10},
	dist{12, 10},
}

type Cards struct {
	drawPile    []int
	discardPile []int
}

// AvgRank returns the average rank across all cards
func AvgRank() int {
	sum := 0
	count := 0
	for _, d := range distribution {
		count += d.count
		sum += d.count * d.rank
	}
	return sum / count
}

// New returns a shuffled draw pile and a discard pile with one card
func New() Cards {
	var c Cards

	for _, d := range distribution {
		for i := 1; i <= d.count; i++ {
			c.drawPile = append(c.drawPile, d.rank)
		}
	}

	Shuffle(c.drawPile)

	// Discard the top card
	c.Discard(c.Draw())

	return c
}

// Shuffle randomizes the order of the given cards
func Shuffle(c []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c), func(i, j int) { c[i], c[j] = c[j], c[i] })
}

// Draw returns the top card from the draw pile
func (c *Cards) Draw() int {
	if len(c.drawPile) == 0 {
		// Leave the top card of the discard pile in the discard pile
		// Shuffle the rest into the draw pile
		c.drawPile = c.discardPile[1:]
		c.discardPile = []int{c.discardPile[0]}
		Shuffle(c.drawPile)
	}

	rank := c.drawPile[0]
	c.drawPile = c.drawPile[1:]
	return rank
}

// Discard places the given card on top of the discard pile
func (c *Cards) Discard(rank int) {
	c.discardPile = append(c.discardPile, rank)
}

// LookDiscard returns the rank of the top card in the discard pile
func (c Cards) LookDiscard() int {
	return c.discardPile[len(c.discardPile)-1]
}

// DrawDiscard draws the top card from the discard pile
func (c *Cards) DrawDiscard() int {
	rank := c.LookDiscard()
	c.discardPile = c.discardPile[:len(c.discardPile)-1]
	return rank
}

func (c Cards) Print() {
	fmt.Println("X", c.LookDiscard())
}

// shorten returns a string with up to 20 elements from the start/end of the slice
func shorten(s []int) string {
	l := len(s)

	if l <= 20 {
		return fmt.Sprint(s)
	}

	str := ""
	str += fmt.Sprint(s[0:6])
	str += fmt.Sprint(" ... ")
	str += fmt.Sprint(s[l-6 : l])

	return str
}

// Print prints the given cards in order
func (c Cards) PrintDebug() {
	fmt.Println("Draw: ->", shorten(c.drawPile))
	fmt.Println("Discard:", shorten(c.discardPile), "<-")
}
