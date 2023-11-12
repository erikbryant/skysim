package util

import (
	"golang.org/x/term"
	"os"
	"strings"
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
func Choose(choices string) byte {
	for {
		choice := readChar()
		if strings.Index(choices, string(choice)) != -1 {
			return choice
		}
	}
}
