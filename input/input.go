package input

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
)

type Input struct {
	Prompt     string
	Choices    []string
	SelectedFn func(string)
}

func (i *Input) Ask() {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("setting stdin to raw:", err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	input := ""
	index := 0
	selectedIndex := 0
	matchingChoices := match(i.Choices, input)
	viewStartIndex := 0
	viewEndIndex := min(len(matchingChoices), 10)

	printSelection(i.Prompt, input, index, matchingChoices, selectedIndex, viewStartIndex, viewEndIndex)

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			fmt.Print("\x1b[0J")
			break
		} else if char >= 32 && char <= 126 {
			input = input[:index] + string(char) + input[index:]
			index += 1
			matchingChoices = match(i.Choices, input)
			viewStartIndex = 0
			viewEndIndex = min(len(matchingChoices), 10)
			selectedIndex = 0
		} else if key == keyboard.KeySpace {
			input = input[:index] + string(" ") + input[index:]
			index += 1
			matchingChoices = match(i.Choices, input)
			viewStartIndex = 0
			viewEndIndex = min(len(matchingChoices), 10)
			selectedIndex = 0
		} else if key == keyboard.KeyArrowLeft {
			index = max(len(i.Prompt), index-1)
		} else if key == keyboard.KeyArrowRight {
			index = min(len(input), index+1)
		} else if key == keyboard.KeyArrowDown {
			selectedIndex = min(len(matchingChoices)-1, selectedIndex+1)
			if selectedIndex >= viewEndIndex {
				viewStartIndex++
				viewEndIndex++
			}
		} else if key == keyboard.KeyArrowUp {
			selectedIndex = max(0, selectedIndex-1)
			if selectedIndex < viewStartIndex {
				viewStartIndex--
				viewEndIndex--
			}
		} else if key == keyboard.KeyBackspace {
			if index > 0 {
				input = input[:index-1] + input[index:]
				index -= 1
				matchingChoices = match(i.Choices, input)
				viewStartIndex = 0
				viewEndIndex = min(len(matchingChoices), 10)
				selectedIndex = 0
			}
		} else if key == keyboard.KeyEnter && i.SelectedFn != nil {
			i.SelectedFn(matchingChoices[selectedIndex])
		}

		printSelection(i.Prompt, input, index, matchingChoices, selectedIndex, viewStartIndex, viewEndIndex)
	}
}

func match(choices []string, input string) []string {
	matchingChoices := make([]string, 0)
	for _, choice := range choices {
		if strings.Contains(strings.ToLower(choice), strings.ToLower(input)) {
			matchingChoices = append(matchingChoices, choice)
		}
	}
	return matchingChoices
}

func printSelection(prompt string, input string, index int, matchingChoices []string, selectedIndex int, viewStartIndex int, viewEndIndex int) {
	fmt.Print("\x1b[s")                                  // save cursor
	fmt.Print("\x1b[1000D")                              // move cursor to left
	fmt.Print("\x1b[K")                                  // clear line
	fmt.Print("\x1b[1;34m" + prompt + "\x1b[0m" + input) // print input
	fmt.Println()
	for i := 0; i < len(matchingChoices); i++ {
		choice := matchingChoices[i]
		if i >= viewStartIndex && i < viewEndIndex {
			fmt.Print("\x1b[K") // clear current line
			matchIndex := strings.Index(strings.ToLower(choice), strings.ToLower(input))

			preMatchPart := choice[0:matchIndex]
			matchPart := choice[matchIndex : matchIndex+len(input)]
			postMatchPart := choice[matchIndex+len(input):]
			if i == selectedIndex {
				fmt.Println("\x1b[30;47m" + preMatchPart + "\x1b[36m" + matchPart + "\x1b[30;47m" + postMatchPart + "\x1b[0m")
			} else {
				fmt.Println(preMatchPart + "\x1b[36m" + matchPart + "\x1b[0m" + postMatchPart)
			}
		}
	}

	fmt.Print("\x1b[0J")                                       // clear till end of screen
	fmt.Print("\x1b[u")                                        // restore cursor
	fmt.Print("\x1b[1000D")                                    // move cursor back to left
	fmt.Print("\x1b[" + strconv.Itoa(index+len(prompt)) + "C") // move cursor right
}

func max(i int, j int) int {
	if i > j {
		return i
	} else {
		return j
	}
}

func min(i int, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}
