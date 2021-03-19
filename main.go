package main

import (
	"fmt"
	"os"

	"github.com/rmadamanchi/go-input/input"
)

func main() {
	choices := make([]string, 0)
	choices = append(choices, "this is cool")
	choices = append(choices, "this is a choice")
	choices = append(choices, "this is another choice")

	i := &input.Input{
		Prompt:  "> ",
		Choices: choices,
		SelectedFn: func(choice string) {
			fmt.Println("\nselected: " + choice)
			fmt.Print("\x1b[0J")
			os.Exit(0)
		},
	}

	i.Ask()

}
