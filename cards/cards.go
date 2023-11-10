package cards

import (
	"sort"
)

const distribution = [][]int{
	[]int{5, -2},
	[]int{10, -1},
	[]int{15, 0},
	[]int{10, 1},
	[]int{10, 2},
	[]int{10, 3},
	[]int{10, 4},
	[]int{10, 5},
	[]int{10, 6},
	[]int{10, 7},
	[]int{10, 8},
	[]int{10, 9},
	[]int{10, 10},
	[]int{10, 11},
	[]int{10, 12},
}

type Cards []int

func Deck() Cards {
	deck := []int{}

	for _, val := range distribution {
		for i := 1; i <= val[0]; i++ {
			deck = append(deck, val[1])
		}
	}

	return deck
}

func (c *Cards) Shuffle() Cards {
	sort.Shuffle(rand.New(rand.NewSource(time.Now().UnixNano())), c)
}

func (c *Cards) Print() {
	for _, val := range c {
		fmt.Print(c)
		fmt.Print(" ")
	}
	fmt.Println()
}
