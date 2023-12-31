package tableau

import (
	"fmt"
	"github.com/erikbryant/skysim/cards"
	"slices"
)

// Cards in the player's tableau may or may not be visible
type tableauCard struct {
	rank    int
	visible bool
}

// Each player starts with the same size tableau
const vRows = 4
const vRowLen = 3

type VRow [vRowLen]tableauCard
type Tableau []VRow

func Deal(c *cards.Cards) Tableau {
	var t Tableau

	for i := 0; i < vRows; i++ {
		t = append(t, VRow{})
	}

	for vRow := range t {
		for hRow := range t[vRow] {
			t[vRow][hRow].rank = c.Draw()
			t[vRow][hRow].visible = false
		}
	}

	return t
}

func (t Tableau) vRows() int {
	return len(t)
}

func (t Tableau) vRowLen() int {
	return vRowLen
}

func (t Tableau) Expected() int {
	return cards.AvgRank() * t.vRows() * t.vRowLen()
}

// Get returns the card at the  given location or an error if that is out of bounds
func (t Tableau) Get(vRow, hRow int) (tableauCard, error) {
	if vRow < 0 || hRow < 0 {
		return tableauCard{}, fmt.Errorf("Out of bounds")
	}
	if vRow >= t.vRows() || hRow >= t.vRowLen() {
		return tableauCard{}, fmt.Errorf("Out of bounds")
	}
	return t[vRow][hRow], nil
}

// Set sets the rank at the  given location or returns an error if that is out of bounds
func (t *Tableau) SetRank(vRow, hRow, rank int) error {
	if vRow < 0 || hRow < 0 {
		return fmt.Errorf("Out of bounds")
	}
	if vRow > t.vRows() || hRow > t.vRowLen() {
		return fmt.Errorf("Out of bounds")
	}

	(*t)[vRow][hRow].rank = rank
	return nil
}

// Reveal reveals the given card, or returns error if that is invalid
func (t *Tableau) Reveal(vRow, hRow int, c *cards.Cards) error {
	tc, err := t.Get(vRow, hRow)
	if err != nil || tc.visible {
		return fmt.Errorf("Reveal failed")
	}

	(*t)[vRow][hRow].visible = true
	t.removeCompletedVRows(c)

	return nil
}

// Replace replaces a tableau card with the given card
func (t *Tableau) Replace(vRow, hRow, rank int, c *cards.Cards) error {
	tc, err := t.Get(vRow, hRow)
	if err != nil {
		return err
	}

	c.Discard(tc.rank)

	t.SetRank(vRow, hRow, rank)
	t.Reveal(vRow, hRow, c)

	t.removeCompletedVRows(c)

	return nil
}

// matchingVRow returns whether all cards in the vertical row match and are visible
func matchingVRow(vRow VRow) bool {
	val := vRow[0].rank

	for _, card := range vRow {
		if !card.visible || card.rank != val {
			return false
		}
	}

	return true
}

// removeCompletedVRows removed any vertical rows where the cards match and are visible
func (t *Tableau) removeCompletedVRows(c *cards.Cards) {
	for vRow := len(*t) - 1; vRow >= 0; vRow-- {
		if matchingVRow((*t)[vRow]) {
			for hRow := range (*t)[vRow] {
				tc, _ := t.Get(vRow, hRow)
				c.Discard(tc.rank)
			}
			*t = slices.DeleteFunc(*t, matchingVRow)
		}
	}
}

// FirstHidden returns the locatio of the first hidden card
func (t Tableau) FirstHidden() (int, int) {
	for vRow := range t {
		for hRow := range t[vRow] {
			c, _ := t.Get(vRow, hRow)
			if !c.visible {
				return vRow, hRow
			}
		}
	}

	return -1, -1
}

// HighestVisible returns the rank and location of the highest visible card
func (t Tableau) HighestVisible() (int, int, int) {
	rank := -999
	vr := -1
	hr := -1

	for vRow := range t {
		for hRow := range t[vRow] {
			c, _ := t.Get(vRow, hRow)
			if c.visible && c.rank > rank {
				rank = c.rank
				vr = vRow
				hr = hRow
			}
		}
	}

	return rank, vr, hr
}

func (t *Tableau) RevealAll(c *cards.Cards) {
	// Count backwards, as some of the vRows may disappear
	// when they are revealed if they are a matching row
	for vRow := t.vRows() - 1; vRow >= 0; vRow-- {
		for hRow := range (*t)[vRow] {
			t.Reveal(vRow, hRow, c)
		}
	}
}

// Out returns true if the player is out (all cards exposed)
func (t Tableau) Out() bool {
	for vRow := range t {
		for hRow := range t[vRow] {
			tc, _ := t.Get(vRow, hRow)
			if !tc.visible {
				return false
			}
		}
	}

	return true
}

// Score returns the visible, expected, and actual scores
func (t Tableau) Score() (int, int, int) {
	vScore := 0
	eScore := 0
	aScore := 0

	for vRow := range t {
		for hRow := range t[vRow] {
			tc, _ := t.Get(vRow, hRow)
			if tc.visible {
				vScore += tc.rank
				eScore += tc.rank
			} else {
				eScore += cards.AvgRank()
			}
			aScore += tc.rank
		}
	}

	return vScore, eScore, aScore
}

func (t Tableau) Print(c cards.Cards) {
	visible, estimated, _ := t.Score()
	fmt.Printf("%2d %2d\n", visible, estimated)

	if t.vRows() == 0 {
		fmt.Println("All vertical rows empty! :)")
		return
	}

	for hRow := range t[0] {
		for vRow := range t {
			tc, _ := t.Get(vRow, hRow)
			if tc.visible {
				mask := cards.MaskForRank(tc.rank)
				mask.Printf("%2d", tc.rank)
				fmt.Print(" ")
			} else {
				mask := cards.MaskForBack()
				mask.Printf("SJ")
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (t Tableau) PrintDebug(c cards.Cards) {
	if t.vRows() == 0 {
		fmt.Println("All vertical rows empty! :)")
		return
	}

	for hRow := range t[0] {
		for vRow := range t {
			tc, _ := t.Get(vRow, hRow)
			mask := cards.MaskForRank(tc.rank)
			mask.Printf("%2d", tc.rank)
			fmt.Print(" ")
		}
		fmt.Println()
	}

	c.PrintDebug()

	fmt.Println(t.Score())
}
