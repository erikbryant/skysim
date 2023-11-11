package cards

import (
	"fmt"
	"github.com/fatih/color"
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

// maskForRank returns the console text mask for the given rank
func MaskForRank(rank int) *color.Color {
	mask := color.New(color.FgBlack, color.Bold)

	if rank <= -1 {
		return mask.Add(color.BgBlue)
	}
	if rank == 0 {
		return mask.Add(color.BgCyan)
	}
	if rank <= 4 {
		return mask.Add(color.BgGreen)
	}
	if rank <= 8 {
		return mask.Add(color.BgYellow)
	}

	return mask.Add(color.BgRed)
}

func MaskForBack() *color.Color {
	mask := color.New(color.FgBlack, color.Bold)
	return mask.Add(color.BgWhite)
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
		c.drawPile = c.discardPile[:len(c.discardPile)-1]
		c.discardPile = []int{c.discardPile[len(c.discardPile)-1]}
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
	mask := MaskForBack()
	mask.Printf("SJ")
	fmt.Print(" ")
	rank := c.LookDiscard()
	mask = MaskForRank(rank)
	mask.Printf("%2d", rank)
	fmt.Println()
}

// shortPrint prints a string with up to 20 elements from the start/end of the slice
func shortPrint(s []int) {
	l := len(s)

	if l <= 20 {
		for _, rank := range s {
			mask := MaskForRank(rank)
			mask.Printf("%2d", rank)
			fmt.Print(" ")
		}
		return
	}

	for i := 0; i <= 5; i++ {
		mask := MaskForRank(s[i])
		mask.Printf("%2d", s[i])
		fmt.Print(" ")
	}
	fmt.Print("... ")
	for i := l - 6; i < l; i++ {
		mask := MaskForRank(s[i])
		mask.Printf("%2d", s[i])
		fmt.Print(" ")
	}
}

// Print prints the given cards in order
func (c Cards) PrintDebug() {
	fmt.Print("Draw: -> ")
	shortPrint(c.drawPile)
	fmt.Println()
	fmt.Print("Discard: ")
	shortPrint(c.discardPile)
	fmt.Println("<-")
}
