package input

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"golang.org/x/term"
	"os"
)

func InitTerminal() (*term.State, error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("error setting stdin to raw: %v", err)
	}

	if err := keyboard.Open(); err != nil {
		return nil, fmt.Errorf("error openning keyboard: %v", err)
	}

	return state, nil
}

func RestoreTerminal(state *term.State) error {
	if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
		return fmt.Errorf("error restoring terminal state: %v", err)
	}

	if err := keyboard.Close(); err != nil {
		return fmt.Errorf("error closing keyboard: %v", err)
	}

	return nil
}
