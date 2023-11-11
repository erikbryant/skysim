package main

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
	"github.com/erikbryant/skysim/tableau"
	"os"
	"strings"

	"golang.org/x/term"
)

// readChar returns the key that was pressed
func readChar() byte {
	// switch stdin into 'raw' mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		panic(err)
	}

	return b[0]
}

// choose returns which of the choices the user selects
func choose(choices string) byte {
	for {
		choice := readChar()
		if strings.Index(choices, string(choice)) != -1 {
			return choice
		}
	}
}

// Reveal reveals a single card
func Reveal(t *tableau.Tableau, c *cards.Cards) {
	var vRow int
	var hRow int

	for {
		t.Print(*c)
		fmt.Print("\nChoose a card to reveal (vRow hRow): ")
		fmt.Scanf("%d %d", &vRow, &hRow)
		if t.Reveal(vRow, hRow, c) == nil {
			break
		}
	}
}

// Draw draws (and plays) a card
func Draw(t *tableau.Tableau, c *cards.Cards) {
	rank := c.Draw()
	fmt.Print("Drew: ")
	mask := cards.MaskForRank(rank)
	mask.Printf("%2d", rank)

	fmt.Println()
	fmt.Println("(r)eplace a tableau card")
	fmt.Println("(d)iscard it and reveal one")

	switch choose("rd") {
	case 'r':
		Replace(t, c, rank)
	case 'd':
		c.Discard(rank)
		Reveal(t, c)
	}
}

// Replace replaces a card in the tableau with the given card
func Replace(t *tableau.Tableau, c *cards.Cards, rank int) {
	var vRow int
	var hRow int

	fmt.Print("\nChoose a card to replace (vRow hRow): ")
	fmt.Scanf("%d %d", &vRow, &hRow)
	t.Replace(vRow, hRow, rank, c)
}

// TakeAnotherTurn processes a player's turn and returns whether they have gone out (or quit)
func TakeAnotherTurn(t *tableau.Tableau, c *cards.Cards) bool {
	fmt.Println()
	t.Print(*c)

	fmt.Println("(d)raw a new card")
	fmt.Println("(r)eplace a tableau card with the discard")
	fmt.Println("(p)rint debug")
	fmt.Println("(q)uit")

	switch choose("drpq") {
	case 'd':
		Draw(t, c)
	case 'r':
		rank := c.DrawDiscard()
		Replace(t, c, rank)
	case 'p':
		fmt.Println()
		t.PrintDebug(*c)
	case 'q':
		return false
	}

	return !t.Out()
}

func main() {
	c := cards.New()
	p1 := tableau.Deal(c)

	Reveal(&p1, &c)
	Reveal(&p1, &c)

	for TakeAnotherTurn(&p1, &c) {
	}

	fmt.Println()
	p1.PrintDebug(c)
}
