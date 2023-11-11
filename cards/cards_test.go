package cards

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New()

	expectedLen := 0
	for _, val := range distribution {
		expectedLen += val.count
	}
	expectedLen-- // One card is moved to the discard pile

	answer := len(c.drawPile)
	if answer != expectedLen {
		t.Errorf("For drawPile expected len %d, got %d", expectedLen, answer)
	}

	answer = len(c.discardPile)
	if answer != 1 {
		t.Errorf("For discardPile expected len %d, got %d", 1, answer)
	}

}

func TestDraw(t *testing.T) {
	testCases := []struct {
		c        Cards
		expected int
	}{
		{Cards{[]int{0, 1, 2}, []int{}}, 0},
		{Cards{[]int{1}, []int{}}, 1},
		{Cards{[]int{}, []int{1, 0}}, 1},
		{Cards{[]int{}, []int{2, 2, 2, 2, 0}}, 2},
	}

	for _, testCase := range testCases {
		answer := testCase.c.Draw()
		if answer != testCase.expected {
			t.Errorf("For %v expected %d, got %d", testCase.c, testCase.expected, answer)
		}
	}
}

func TestDiscard(t *testing.T) {
	testCases := []struct {
		c        Cards
		expected int
	}{
		{Cards{[]int{}, []int{}}, 0},
		{Cards{[]int{}, []int{3}}, 1},
		{Cards{[]int{}, []int{9, 0}}, 2},
		{Cards{[]int{}, []int{-1, 3, 5}}, 3},
	}

	for _, testCase := range testCases {
		testCase.c.Discard(testCase.expected)
		answer := testCase.c.LookDiscard()
		if answer != testCase.expected {
			t.Errorf("For %v expected %d, got %d", testCase.c, testCase.expected, answer)
		}
	}
}

func TestLookDiscard(t *testing.T) {
	testCases := []struct {
		c        Cards
		expected int
	}{
		{Cards{[]int{}, []int{0}}, 0},
		{Cards{[]int{}, []int{1}}, 1},
		{Cards{[]int{}, []int{9, 0, 2}}, 2},
		{Cards{[]int{}, []int{-1, 3, 5, 3}}, 3},
	}

	for _, testCase := range testCases {
		answer := testCase.c.LookDiscard()
		if answer != testCase.expected {
			t.Errorf("For %v expected %d, got %d", testCase.c, testCase.expected, answer)
		}
	}
}

func TestDrawDiscard(t *testing.T) {
	testCases := []struct {
		c        Cards
		expected int
	}{
		{Cards{[]int{}, []int{0}}, 0},
		{Cards{[]int{}, []int{1}}, 1},
		{Cards{[]int{}, []int{9, 0, 2}}, 2},
		{Cards{[]int{}, []int{-1, 3, 5, 3}}, 3},
	}

	for _, testCase := range testCases {
		answer := testCase.c.DrawDiscard()
		if answer != testCase.expected {
			t.Errorf("For %v expected %d, got %d", testCase.c, testCase.expected, answer)
		}
	}
}
