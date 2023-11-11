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

func Deal(c cards.Cards) Tableau {
	var t Tableau

	for i := 0; i < vRows; i++ {
		t = append(t, VRow{})
	}

	for row := range t {
		for col := range t[row] {
			t[row][col].rank = c.Draw()
			t[row][col].visible = false
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

func (t Tableau) Get(row, col int) (tableauCard, error) {
	if row < 0 || col < 0 {
		return tableauCard{}, fmt.Errorf("Out of bounds")
	}
	if row > t.vRows() || col > t.vRowLen() {
		return tableauCard{}, fmt.Errorf("Out of bounds")
	}
	return t[row][col], nil
}

// Reveal reveals the given card, or returns error if that is invalid
func (t *Tableau) Reveal(vRow, col int, c *cards.Cards) error {
	tc, err := t.Get(vRow, col)
	if err != nil || tc.visible {
		return fmt.Errorf("Reveal failed")
	}

	(*t)[vRow][col].visible = true
	t.removeCompletedVRows(c)

	return nil
}

// Replace replaces a tableau card with the given card
func (t *Tableau) Replace(vRow, col, rank int, c *cards.Cards) {
	c.Discard((*t)[vRow][col].rank)
	(*t)[vRow][col].rank = rank
	(*t)[vRow][col].visible = true
	t.removeCompletedVRows(c)
}

func matchingVRow(vRow VRow) bool {
	val := vRow[0].rank

	for _, card := range vRow {
		if !card.visible || card.rank != val {
			return false
		}
	}

	return true
}

func (t *Tableau) removeCompletedVRows(c *cards.Cards) {
	for vRow := len(*t) - 1; vRow >= 0; vRow-- {
		if matchingVRow((*t)[vRow]) {
			for i := range (*t)[vRow] {
				c.Discard((*t)[vRow][i].rank)
			}
			*t = slices.DeleteFunc(*t, matchingVRow)
		}
	}
}

// Out returns true if the player is out (all cards exposed)
func (t Tableau) Out() bool {
	for row := range t {
		for col := range t[row] {
			if !t[row][col].visible {
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

	for row := range t {
		for col := range t[row] {
			if t[row][col].visible {
				vScore += t[row][col].rank
				eScore += t[row][col].rank
			} else {
				eScore += cards.AvgRank()
			}
			aScore += t[row][col].rank
		}
	}

	return vScore, eScore, aScore
}

func (t Tableau) Print(c cards.Cards) {
	visible, estimated, _ := t.Score()
	fmt.Printf("%2d %2d ", visible, estimated)
	c.Print()

	if t.vRows() == 0 {
		fmt.Println("All vertical rows empty! :)")
		return
	}

	for col := range t[0] {
		for vRow := range t {
			if t[vRow][col].visible {
				rank := t[vRow][col].rank
				mask := cards.MaskForRank(rank)
				mask.Printf("%2d", rank)
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

	for col := range t[0] {
		for vRow := range t {
			rank := t[vRow][col].rank
			mask := cards.MaskForRank(rank)
			mask.Printf("%2d", rank)
			fmt.Print(" ")
		}
		fmt.Println()
	}

	c.PrintDebug()

	fmt.Println(t.Score())
}
