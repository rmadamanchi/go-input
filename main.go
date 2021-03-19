package main

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/rmadamanchi/go-input/input"
)

func main() {
	choices := make([]string, 0)
	choices = append(choices, "Alaska")
	choices = append(choices, "Arizona")
	choices = append(choices, "Arkansas")
	choices = append(choices, "California")
	choices = append(choices, "Colorado")
	choices = append(choices, "Connecticut")
	choices = append(choices, "Delaware")
	choices = append(choices, "Florida")
	choices = append(choices, "Georgia")
	choices = append(choices, "Hawaii")
	choices = append(choices, "Idaho")
	choices = append(choices, "Illinois")
	choices = append(choices, "Indiana")
	choices = append(choices, "Iowa")
	choices = append(choices, "Kansas")
	choices = append(choices, "Kentucky")
	choices = append(choices, "Louisiana")
	choices = append(choices, "Maine")
	choices = append(choices, "Maryland")
	choices = append(choices, "Massachusetts")
	choices = append(choices, "Michigan")
	choices = append(choices, "Minnesota")
	choices = append(choices, "Mississippi")
	choices = append(choices, "Missouri")
	choices = append(choices, "Montana")
	choices = append(choices, "Nebraska")
	choices = append(choices, "Nevada")
	choices = append(choices, "New Hampshire")
	choices = append(choices, "New Jersey")
	choices = append(choices, "New Mexico")
	choices = append(choices, "New York")
	choices = append(choices, "North Carolina")
	choices = append(choices, "North Dakota")
	choices = append(choices, "Ohio")
	choices = append(choices, "Oklahoma")
	choices = append(choices, "Oregon")
	choices = append(choices, "Pennsylvania")
	choices = append(choices, "Rhode Island")
	choices = append(choices, "South Carolina")
	choices = append(choices, "South Dakota")
	choices = append(choices, "Tennessee")
	choices = append(choices, "Texas")
	choices = append(choices, "Utah")
	choices = append(choices, "Vermont")
	choices = append(choices, "Virginia")
	choices = append(choices, "Washington")
	choices = append(choices, "West Virginia")
	choices = append(choices, "Wisconsin")
	choices = append(choices, "Wyoming")

	i := &input.Input{
		Prompt:  "> ",
		Choices: choices,
		KeyBindings: map[keyboard.Key]func(*input.Input, string){
			keyboard.KeyEnter: func(i *input.Input, choice string) {
				i.Clear()
				fmt.Println("\nselected: " + choice)
				os.Exit(0)
			},
			keyboard.KeyEsc: func(i *input.Input, choice string) {
				i.Clear()
				os.Exit(0)
			}},
	}

	i.Ask()

}
