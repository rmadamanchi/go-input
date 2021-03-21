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

type Selection struct {
	Prompt           string
	Choices          []Choice
	PageSize         int
	DefaultSelection string
	ValueFn          func(*Choice) string
	KeyBindings      map[keyboard.Key]func(*Selection, Choice)
	closed           bool
}

type Choice struct {
	Data  interface{}
	Value string
}

func (s *Selection) Close() {
	fmt.Print("\x1b[0J") // clear till end of screen
	s.closed = true
}

func (s *Selection) Ask() {
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

	if s.PageSize == 0 {
		s.PageSize = 10
	}

	input := ""
	index := 0
	selectedIndex := 0
	if s.DefaultSelection != "" {
		for i, choice := range s.Choices {
			if s.getValue(&choice) == s.DefaultSelection {
				selectedIndex = i
			}
		}
	}

	matchingChoices := s.match(input)
	viewStartIndex := 0
	if selectedIndex >= s.PageSize {
		if selectedIndex > len(matchingChoices)-s.PageSize {
			viewStartIndex = len(matchingChoices) - s.PageSize
		} else {
			viewStartIndex = selectedIndex
		}
	}

	s.render(input, index, matchingChoices, selectedIndex, viewStartIndex)

	for {
		if s.closed {
			return
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if char >= 32 && char <= 126 {
			input = input[:index] + string(char) + input[index:]
			index += 1
			matchingChoices = s.match(input)
			viewStartIndex = 0
			selectedIndex = 0
		} else if key == keyboard.KeySpace {
			input = input[:index] + string(" ") + input[index:]
			index += 1
			matchingChoices = s.match(input)
			viewStartIndex = 0
			selectedIndex = 0
		} else if key == keyboard.KeyArrowLeft {
			index = max(0, index-1)
		} else if key == keyboard.KeyArrowRight {
			index = min(len(input), index+1)
		} else if key == keyboard.KeyArrowDown {
			selectedIndex = min(len(matchingChoices)-1, selectedIndex+1)
			if selectedIndex >= viewStartIndex+s.PageSize {
				viewStartIndex++
			}
		} else if key == keyboard.KeyArrowUp {
			selectedIndex = max(0, selectedIndex-1)
			if selectedIndex < viewStartIndex {
				viewStartIndex--
			}
		} else if key == keyboard.KeyBackspace {
			if index > 0 {
				input = input[:index-1] + input[index:]
				index -= 1
				matchingChoices = s.match(input)
				viewStartIndex = 0
				selectedIndex = 0
			}
		} else if keyBindingFn, ok := s.KeyBindings[key]; ok {
			keyBindingFn(s, matchingChoices[selectedIndex])
		}

		if !s.closed {
			s.render(input, index, matchingChoices, selectedIndex, viewStartIndex)
		}
	}
}

func (s *Selection) getValue(c *Choice) string {
	if s.ValueFn != nil {
		return s.ValueFn(c)
	} else {
		return c.Value
	}
}

func (s *Selection) match(input string) []Choice {
	matchingChoices := make([]Choice, 0)
	for _, choice := range s.Choices {
		if strings.Contains(strings.ToLower(s.getValue(&choice)), strings.ToLower(input)) {
			matchingChoices = append(matchingChoices, choice)
		}
	}
	return matchingChoices
}

func (s *Selection) render(input string, index int, matchingChoices []Choice, selectedIndex int, viewStartIndex int) {
	fmt.Print("\x1b[s")                                    // save cursor
	fmt.Print("\x1b[1000D")                                // move cursor to left
	fmt.Print("\x1b[K")                                    // clear line
	fmt.Print("\x1b[1;34m" + s.Prompt + "\x1b[0m" + input) // print input
	fmt.Println()
	for i := 0; i < len(matchingChoices); i++ {
		choice := matchingChoices[i]
		if i >= viewStartIndex && i < viewStartIndex+s.PageSize {
			value := s.getValue(&choice)
			fmt.Print("\x1b[K") // clear current line
			matchIndex := strings.Index(strings.ToLower(value), strings.ToLower(input))

			preMatchPart := value[0:matchIndex]
			matchPart := value[matchIndex : matchIndex+len(input)]
			postMatchPart := value[matchIndex+len(input):]
			if i == selectedIndex {
				fmt.Println("\x1b[30;47m" + preMatchPart + "\x1b[36m" + matchPart + "\x1b[30;47m" + postMatchPart + "\x1b[0m")
			} else {
				fmt.Println(preMatchPart + "\x1b[36m" + matchPart + "\x1b[0m" + postMatchPart)
			}
		}
	}
	//debug
	//fmt.Println(strconv.Itoa(selectedIndex) + "-" + strconv.Itoa(viewStartIndex) + "\x1b[0J")

	fmt.Print("\x1b[0J")                                         // clear till end of screen
	fmt.Print("\x1b[u")                                          // restore cursor
	fmt.Print("\x1b[1000D")                                      // move cursor back to left
	fmt.Print("\x1b[" + strconv.Itoa(index+len(s.Prompt)) + "C") // move cursor right
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
