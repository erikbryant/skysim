package main

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
	"github.com/erikbryant/skysim/tableau"
)

func Reveal(t *tableau.Tableau, c *cards.Cards) {
	var vRow int
	var hRow int

	t.Print(*c)
	fmt.Print("\nChoose a card to reveal (c r): ")
	fmt.Scanf("%d %d", &vRow, &hRow)
	t.Reveal(vRow, hRow, c)
}

func Draw(t *tableau.Tableau, c *cards.Cards) {
	var choice string

	rank := c.Draw()
	fmt.Print("Drew: ")
	mask := cards.MaskForRank(rank)
	mask.Printf("%2d", rank)
	fmt.Println()
	fmt.Println("(r)eplace a tableau card")
	fmt.Println("(d)iscard it and reveal one")
	fmt.Scanf("%s", &choice)
	switch choice {
	case "r":
		Replace(t, c, rank)
	case "d":
		c.Discard(rank)
		Reveal(t, c)
	default:
		fmt.Print("\nUnknown input. Please try again.\n\n")
	}
}

func Replace(t *tableau.Tableau, c *cards.Cards, rank int) {
	var vRow int
	var hRow int

	fmt.Print("\nChoose a card to replace (c r): ")
	fmt.Scanf("%d %d", &vRow, &hRow)
	t.Replace(vRow, hRow, rank, c)
}

func TakeAnotherTurn(t *tableau.Tableau, c *cards.Cards) bool {
	var choice string

	fmt.Println()
	t.Print(*c)
	fmt.Println("(d)raw a new card")
	fmt.Println("(r)eplace a tableau card with the discard")
	fmt.Println("(p)rint debug")
	fmt.Println("(q)uit")
	fmt.Print("Choose one (drpq): ")
	fmt.Scanf("%s", &choice)

	switch choice {
	case "d":
		Draw(t, c)
	case "r":
		rank := c.DrawDiscard()
		Replace(t, c, rank)
	case "p":
		fmt.Println()
		t.PrintDebug(*c)
	case "q":
		return false
	default:
		fmt.Print("\nUnknown input. Please try again.\n\n")
	}

	return true
}

func main() {
	c := cards.New()
	p1 := tableau.Deal(c)

	Reveal(&p1, &c)
	Reveal(&p1, &c)

	for {
		if !TakeAnotherTurn(&p1, &c) {
			break
		}
	}
}
