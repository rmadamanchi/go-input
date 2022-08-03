package main

import (
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/rmadamanchi/go-input"
)

func main() {
	state, _ := input.InitTerminal()
	defer func() { _ = input.RestoreTerminal(state) }()

	choices := make([]input.Choice, 0)
	choices = append(choices, input.Choice{Value: "Alaska"})
	choices = append(choices, input.Choice{Value: "Arizona"})
	choices = append(choices, input.Choice{Value: "Arkansas"})
	choices = append(choices, input.Choice{Value: "California", Suffix: "(nice weather)"})
	choices = append(choices, input.Choice{Value: "Colorado"})
	choices = append(choices, input.Choice{Value: "Connecticut"})
	choices = append(choices, input.Choice{Value: "Delaware"})
	choices = append(choices, input.Choice{Value: "Florida"})
	choices = append(choices, input.Choice{Value: "Georgia"})
	choices = append(choices, input.Choice{Value: "Hawaii"})
	choices = append(choices, input.Choice{Value: "Idaho"})
	choices = append(choices, input.Choice{Value: "Illinois"})
	choices = append(choices, input.Choice{Value: "Indiana"})
	choices = append(choices, input.Choice{Value: "Iowa"})
	choices = append(choices, input.Choice{Value: "Kansas"})
	choices = append(choices, input.Choice{Value: "Kentucky"})
	choices = append(choices, input.Choice{Value: "Louisiana"})
	choices = append(choices, input.Choice{Value: "Maine"})
	choices = append(choices, input.Choice{Value: "Maryland"})
	choices = append(choices, input.Choice{Value: "Massachusetts"})
	choices = append(choices, input.Choice{Value: "Michigan"})
	choices = append(choices, input.Choice{Value: "Minnesota"})
	choices = append(choices, input.Choice{Value: "Mississippi", Suffix: "(too cold)"})
	choices = append(choices, input.Choice{Value: "Missouri"})
	choices = append(choices, input.Choice{Value: "Montana"})
	choices = append(choices, input.Choice{Value: "Nebraska"})
	choices = append(choices, input.Choice{Value: "Nevada"})
	choices = append(choices, input.Choice{Value: "New Hampshire"})
	choices = append(choices, input.Choice{Value: "New Jersey"})
	choices = append(choices, input.Choice{Value: "New Mexico"})
	choices = append(choices, input.Choice{Value: "New York"})
	choices = append(choices, input.Choice{Value: "North Carolina"})
	choices = append(choices, input.Choice{Value: "North Dakota"})
	choices = append(choices, input.Choice{Value: "Ohio"})
	choices = append(choices, input.Choice{Value: "Oklahoma"})
	choices = append(choices, input.Choice{Value: "Oregon"})
	choices = append(choices, input.Choice{Value: "Pennsylvania"})
	choices = append(choices, input.Choice{Value: "Rhode Island"})
	choices = append(choices, input.Choice{Value: "South Carolina"})
	choices = append(choices, input.Choice{Value: "South Dakota"})
	choices = append(choices, input.Choice{Value: "Tennessee"})
	choices = append(choices, input.Choice{Value: "Texas", Suffix: "(too hot)"})
	choices = append(choices, input.Choice{Value: "Utah"})
	choices = append(choices, input.Choice{Value: "Vermont"})
	choices = append(choices, input.Choice{Value: "Virginia"})
	choices = append(choices, input.Choice{Value: "Washington"})
	choices = append(choices, input.Choice{Value: "West Virginia"})
	choices = append(choices, input.Choice{Value: "Wisconsin"})
	choices = append(choices, input.Choice{Value: "Wyoming"})

	var s *input.Selection
	s = &input.Selection{
		Prompt:           "> ",
		Choices:          choices,
		PageSize:         10,
		DefaultSelection: "South Carolina",
		Footer:           "Enter: Select, Esc: Exit, Ctrl+C: Copy",
		ValueFn: func(c *input.Choice) string {
			return c.Value
		},
		KeyBindings: map[keyboard.Key]func(*input.Selection, *input.Choice){
			keyboard.KeyEnter: func(i *input.Selection, choice *input.Choice) {
				i.Hide()
				if choice != nil {
					fmt.Println("selected: " + choice.Value)
				}
			},
			keyboard.KeyEsc: func(i *input.Selection, choice *input.Choice) {
				i.Hide()
			},
			keyboard.KeyCtrlC: func(i *input.Selection, choice *input.Choice) {
				if choice != nil {
					s.FlashMessage("Copied")
				}
			}},
	}

	s.Show()

}
