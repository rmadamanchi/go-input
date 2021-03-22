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

	input           string
	cursorIndex     int
	selectedIndex   int
	matchingChoices []Choice
	viewStartIndex  int

	closed bool
}

type Choice struct {
	Data  interface{}
	Value string
}

func (s *Selection) Hide() {
	fmt.Print("\x1b[1000D\x1b[0J") // move cursor to beginning of line; clear till end of screen
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

	s.input = ""
	s.cursorIndex = 0
	s.selectedIndex = 0
	if s.DefaultSelection != "" {
		for i, choice := range s.Choices {
			if s.getValue(&choice) == s.DefaultSelection {
				s.selectedIndex = i
			}
		}
	}

	s.matchingChoices = s.match(s.input)
	s.viewStartIndex = 0
	if s.selectedIndex >= s.PageSize {
		if s.selectedIndex > len(s.matchingChoices)-s.PageSize {
			s.viewStartIndex = len(s.matchingChoices) - s.PageSize
		} else {
			s.viewStartIndex = s.selectedIndex
		}
	}

	s.inputLoop()
}

func (s *Selection) Show() {
	s.closed = false
	s.inputLoop()
}

func (s *Selection) inputLoop() {
	s.render()

	for {
		if s.closed {
			return
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if char >= 32 && char <= 126 {
			s.input = s.input[:s.cursorIndex] + string(char) + s.input[s.cursorIndex:]
			s.cursorIndex += 1
			s.matchingChoices = s.match(s.input)
			s.viewStartIndex = 0
			s.selectedIndex = 0
		} else if key == keyboard.KeySpace {
			s.input = s.input[:s.cursorIndex] + string(" ") + s.input[s.cursorIndex:]
			s.cursorIndex += 1
			s.matchingChoices = s.match(s.input)
			s.viewStartIndex = 0
			s.selectedIndex = 0
		} else if key == keyboard.KeyArrowLeft {
			s.cursorIndex = max(0, s.cursorIndex-1)
		} else if key == keyboard.KeyArrowRight {
			s.cursorIndex = min(len(s.input), s.cursorIndex+1)
		} else if key == keyboard.KeyArrowDown {
			s.selectedIndex = min(len(s.matchingChoices)-1, s.selectedIndex+1)
			if s.selectedIndex >= s.viewStartIndex+s.PageSize {
				s.viewStartIndex++
			}
		} else if key == keyboard.KeyArrowUp {
			s.selectedIndex = max(0, s.selectedIndex-1)
			if s.selectedIndex < s.viewStartIndex {
				s.viewStartIndex--
			}
		} else if key == keyboard.KeyBackspace {
			if s.cursorIndex > 0 {
				s.input = s.input[:s.cursorIndex-1] + s.input[s.cursorIndex:]
				s.cursorIndex -= 1
				s.matchingChoices = s.match(s.input)
				s.viewStartIndex = 0
				s.selectedIndex = 0
			}
		} else if keyBindingFn, ok := s.KeyBindings[key]; ok {
			keyBindingFn(s, s.matchingChoices[s.selectedIndex])
		}

		if !s.closed {
			s.render()
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

func (s *Selection) render() {
	fmt.Print("\x1b[s")                                      // save cursor
	fmt.Print("\x1b[1000D")                                  // move cursor to left
	fmt.Print("\x1b[K")                                      // clear line
	fmt.Print("\x1b[1;34m" + s.Prompt + "\x1b[0m" + s.input) // print input
	fmt.Println()
	for i := 0; i < len(s.matchingChoices); i++ {
		choice := s.matchingChoices[i]
		if i >= s.viewStartIndex && i < s.viewStartIndex+s.PageSize {
			value := s.getValue(&choice)
			fmt.Print("\x1b[K") // clear current line
			matchIndex := strings.Index(strings.ToLower(value), strings.ToLower(s.input))

			preMatchPart := value[0:matchIndex]
			matchPart := value[matchIndex : matchIndex+len(s.input)]
			postMatchPart := value[matchIndex+len(s.input):]
			if i == s.selectedIndex {
				fmt.Println("\x1b[30;47m" + preMatchPart + "\x1b[36m" + matchPart + "\x1b[30;47m" + postMatchPart + "\x1b[0m")
			} else {
				fmt.Println(preMatchPart + "\x1b[36m" + matchPart + "\x1b[0m" + postMatchPart)
			}
		}
	}
	//debug
	//fmt.Println(strconv.Itoa(selectedIndex) + "-" + strconv.Itoa(viewStartIndex) + "\x1b[0J")

	fmt.Print("\x1b[0J")                                                 // clear till end of screen
	fmt.Print("\x1b[u")                                                  // restore cursor
	fmt.Print("\x1b[1000D")                                              // move cursor back to left
	fmt.Print("\x1b[" + strconv.Itoa(s.cursorIndex+len(s.Prompt)) + "C") // move cursor right
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
