package skysim

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
	"github.com/erikbryant/skysim/tableau"
	"github.com/erikbryant/skysim/util"
	"github.com/fatih/color"
)

type SkySim struct {
	cards    cards.Cards
	tableaus []tableau.Tableau
	player   int
	firstOut int
}

// New returns a new game, ready to play
func New(players int) SkySim {
	s := SkySim{}

	s.cards = cards.New()

	for i := 0; i < players; i++ {
		s.tableaus = append(s.tableaus, tableau.Deal(&s.cards))
	}

	s.player = 0
	s.firstOut = -1

	return s
}

// tableau returns a poiter to the current player's tableau
func (s SkySim) tableau() *tableau.Tableau {
	return &s.tableaus[s.player]
}

// robot returns true if this player is a robot
func (s SkySim) robot() bool {
	return s.player != 0
}

// Reveal reveals a single card
func (s *SkySim) Reveal() {
	if s.robot() {
		s.tableau().Reveal(0, 0, &s.cards)
		s.tableau().Reveal(1, 0, &s.cards)
		return
	}

	var vRow int
	var hRow int

	for {
		fmt.Print("Choose a card to reveal (vRow hRow): ")
		fmt.Scanf("%d %d", &vRow, &hRow)
		if s.tableau().Reveal(vRow, hRow, &s.cards) == nil {
			break
		}
	}
}

// Replace replaces a card in the tableau with the given card
func (s *SkySim) Replace(rank int) {
	var vRow int
	var hRow int

	fmt.Print("Choose a card to replace (vRow hRow): ")
	fmt.Scanf("%d %d", &vRow, &hRow)
	s.tableau().Replace(vRow, hRow, rank, &s.cards)
}

// Draw draws (and plays) a card
func (s *SkySim) Draw() {
	rank := s.cards.Draw()
	fmt.Print("Drew: ")
	mask := cards.MaskForRank(rank)
	mask.Printf("%2d", rank)

	fmt.Println()
	fmt.Println("(r)eplace a tableau card")
	fmt.Println("(d)iscard it and reveal one")

	switch util.Choose("rd") {
	case 'r':
		s.Replace(rank)
	case 'd':
		s.cards.Discard(rank)
		s.Reveal()
	}
}

// TakeAnotherTurn processes a player's turn and returns whether they have gone out (or quit)
func (s *SkySim) TakeTurn() bool {
	fmt.Println()
	s.Print()

	fmt.Println("(d)raw a new card")
	fmt.Println("(r)eplace a tableau card with the discard")
	fmt.Println("(p)rint debug")
	fmt.Println("(q)uit")

	switch util.Choose("drpq") {
	case 'd':
		s.Draw()
	case 'r':
		rank := s.cards.DrawDiscard()
		s.Replace(rank)
	case 'p':
		fmt.Println()
		s.PrintDebug()
		s.TakeTurn()
	case 'q':
		return false
	}

	return !s.tableau().Out()
}

// gameOver returns true when a player has gone out
func (s SkySim) gameOver() bool {
	return s.firstOut >= 0 && s.player == s.firstOut
}

// Play plays the game
func (s *SkySim) Play() {
	// Players each reveal two cards
	for s.player = range s.tableaus {
		s.Print()
		s.Reveal()
		s.Print()
		s.Reveal()
	}
	s.Print()

	// Players alternate turns until someone goes out
	// then each other player gets one more turn
	for !s.gameOver() {
		for s.player = range s.tableaus {
			if s.gameOver() {
				break
			}
			if !s.TakeTurn() && s.firstOut < 0 {
				// Record which player went out first
				s.firstOut = s.player
			}
		}
	}

	fmt.Println()

	// Players reveal and score
	for s.player = range s.tableaus {
		s.tableau().RevealAll(&s.cards)
	}
	s.Print()
}

// Print prints the current game state
func (s SkySim) Print() {
	fmt.Printf("\n\n")
	s.cards.Print()
	for i, t := range s.tableaus {
		fmt.Println()
		header := fmt.Sprintf("** Player %d **", i)
		if i == s.player {
			mask := color.New(color.FgGreen, color.Bold)
			mask.Printf(header)
		} else if i == s.firstOut {
			header += " <-- First out!"
			mask := color.New(color.FgRed, color.Bold)
			mask.Printf(header)
		} else {
			fmt.Printf(header)
		}
		fmt.Println()
		t.Print(s.cards)
	}
}

// PrintDebug prints the current game state, revealing any hidden information
func (s SkySim) PrintDebug() {
	s.tableau().PrintDebug(s.cards)
}
