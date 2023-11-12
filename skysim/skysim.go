package skysim

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
	"github.com/erikbryant/skysim/tableau"
	"github.com/erikbryant/skysim/util"
	"github.com/fatih/color"
)

type Player struct {
	tableau tableau.Tableau
	robot   bool
}

type SkySim struct {
	cards    cards.Cards
	players  []Player
	player   int
	firstOut int
}

// New returns a new game, ready to play
func New(humans, robots int) SkySim {
	s := SkySim{}

	s.cards = cards.New()

	for i := 0; i < humans+robots; i++ {
		p := Player{
			tableau.Deal(&s.cards),
			i >= humans,
		}
		s.players = append(s.players, p)
	}

	s.player = 0
	s.firstOut = -1

	return s
}

// tableau returns a pointer to the current player's tableau
func (s SkySim) tableau() *tableau.Tableau {
	return &s.players[s.player].tableau
}

// robot returns true if this player is a robot
func (s SkySim) robot() bool {
	return s.players[s.player].robot
}

// revealRobot has the robot reveal a single card
func (s *SkySim) revealRobot() {
	hRow, vRow := s.tableau().FirstHidden()
	s.tableau().Reveal(vRow, hRow, &s.cards)
}

// reveal reveals a single card
func (s *SkySim) reveal() {
	if s.robot() {
		s.revealRobot()
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

// replace replaces a card in the tableau with the given card
func (s *SkySim) replace(rank int) {
	var vRow int
	var hRow int

	fmt.Print("Choose a card to replace (vRow hRow): ")
	fmt.Scanf("%d %d", &vRow, &hRow)
	s.tableau().Replace(vRow, hRow, rank, &s.cards)
}

// draw draws (and plays) a card
func (s *SkySim) draw() {
	rank := s.cards.Draw()
	fmt.Print("Drew: ")
	mask := cards.MaskForRank(rank)
	mask.Printf("%2d", rank)

	fmt.Println()
	fmt.Println("(r)eplace a tableau card")
	fmt.Println("(d)iscard it and reveal one")

	switch util.Choose("rd") {
	case 'r':
		s.replace(rank)
	case 'd':
		s.cards.Discard(rank)
		s.reveal()
	}
}

// takeTurnRobot processes a robot's turn and returns whether they have gone out (or quit)
func (s *SkySim) takeTurnRobot() bool {
	fmt.Println("Robot #", s.player)
	// Should we replace a tableau card with the discard?
	discardRank := s.cards.LookDiscard()
	if discardRank < cards.AvgRank() {
		// This is a low card; we'll use it
		s.cards.DrawDiscard()
		visibleRank, vRow, hRow := s.tableau().HighestVisible()
		if visibleRank < discardRank {
			// Our visible cards are lower; replace first hidden card instead
			vRow, hRow = s.tableau().FirstHidden()
		}
		fmt.Printf("  Replacing %d, %d with discard (%d)\n", vRow, hRow, discardRank)
		s.tableau().Replace(vRow, hRow, discardRank, &s.cards)
		return !s.tableau().Out()
	}

	// Draw a card
	fmt.Printf("  Ignoring discard (%d); drawing instead...\n", discardRank)
	drawRank := s.cards.Draw()
	fmt.Println("  Drew:", drawRank)

	// Can we use it?
	if drawRank < cards.AvgRank() {
		// This is a low card; we'll use it
		visibleRank, vRow, hRow := s.tableau().HighestVisible()
		if visibleRank < discardRank {
			// Our visible cards are lower; replace first hidden card instead
			vRow, hRow = s.tableau().FirstHidden()
		}
		fmt.Printf("  Replacing %d, %d with draw (%d)\n", vRow, hRow, drawRank)
		s.tableau().Replace(vRow, hRow, drawRank, &s.cards)
		return !s.tableau().Out()
	}

	// We can't use it. Discard it and reveal a card.
	s.cards.Discard(drawRank)
	vRow, hRow := s.tableau().FirstHidden()
	fmt.Printf("  Ignoring draw; revealing %d, %d\n", vRow, hRow)
	s.tableau().Reveal(vRow, hRow, &s.cards)

	return !s.tableau().Out()
}

// takeTurn processes a player's turn and returns whether they have gone out (or quit)
func (s *SkySim) takeTurn() bool {
	if s.robot() {
		return s.takeTurnRobot()
	}

	fmt.Println()
	s.print()

	fmt.Println("(d)raw a new card")
	fmt.Println("(r)eplace a tableau card with the discard")
	fmt.Println("(p)rint debug")
	fmt.Println("(q)uit")

	switch util.Choose("drpq") {
	case 'd':
		s.draw()
	case 'r':
		rank := s.cards.DrawDiscard()
		s.replace(rank)
	case 'p':
		fmt.Println()
		s.printDebug()
		s.takeTurn()
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
	for s.player = range s.players {
		s.print()
		s.reveal()
		s.print()
		s.reveal()
	}
	s.print()

	// Players alternate turns until someone goes out
	// then each other player gets one more turn
	playing := true
	for playing {
		for s.player = range s.players {
			if s.gameOver() {
				playing = false
				break
			}
			if !s.takeTurn() && s.firstOut < 0 {
				// Record which player went out first
				s.firstOut = s.player
			}
		}
	}

	fmt.Println()

	// Players reveal and score
	for s.player = range s.players {
		s.tableau().RevealAll(&s.cards)
	}
	s.print()
}

// print prints the current game state
func (s SkySim) print() {
	fmt.Printf("\n\n")
	s.cards.Print()
	for i, p := range s.players {
		fmt.Println()
		pType := "human"
		if p.robot {
			pType = "robot"
		}
		header := fmt.Sprintf("** Player %d (%s) **", i, pType)
		if i == s.firstOut {
			header += " <-- First out!"
			mask := color.New(color.FgRed, color.Bold)
			mask.Printf(header)
		} else if i == s.player {
			mask := color.New(color.FgGreen, color.Bold)
			mask.Printf(header)
		} else {
			fmt.Printf(header)
		}
		fmt.Println()
		p.tableau.Print(s.cards)
	}
}

// printDebug prints the current game state, revealing any hidden information
func (s SkySim) printDebug() {
	s.tableau().PrintDebug(s.cards)
}
