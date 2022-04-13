package input

import (
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type Selection struct {
	Prompt           string
	Choices          []Choice
	PageSize         int
	DefaultSelection string
	ValueFn          func(*Choice) string
	KeyBindings      map[keyboard.Key]func(*Selection, *Choice)
	Footer 			 string

	input           string
	cursorIndex     int
	selectedIndex   int
	matchingChoices []Choice
	viewStartIndex  int

	initialized bool
	hidden      bool
	message 	string
}

type Choice struct {
	Data  interface{}
	Value string
}

func (s *Selection) Hide() {
	Printer{}.CursorStartOfLine().ClearScreen()
	s.hidden = true
}

func (s *Selection) initialize() {

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

	s.initialized = true
}

func (s *Selection) Show() {
	if !s.initialized {
		s.initialize()
	}
	s.hidden = false

	s.render()

	for {
		if s.hidden {
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
		} else if key == keyboard.KeyHome {
			s.cursorIndex = 0
		} else if key == keyboard.KeyEnd {
			s.cursorIndex = len(s.input)
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
		} else if key == keyboard.KeyPgdn {
			s.selectedIndex = min(len(s.matchingChoices)-1, s.selectedIndex+s.PageSize)
			s.viewStartIndex = min(s.viewStartIndex + s.PageSize, len(s.matchingChoices)-s.PageSize)
		} else if key == keyboard.KeyPgup {
			s.selectedIndex = max(0, s.selectedIndex-s.PageSize)
			s.viewStartIndex = max(s.viewStartIndex - s.PageSize, 0)
		} else if key == keyboard.KeyBackspace {
			if s.cursorIndex > 0 {
				s.input = s.input[:s.cursorIndex-1] + s.input[s.cursorIndex:]
				s.cursorIndex -= 1
				s.matchingChoices = s.match(s.input)
				s.viewStartIndex = 0
				s.selectedIndex = 0
			}
		} else if keyBindingFn, ok := s.KeyBindings[key]; ok {
			if len(s.matchingChoices) > 0 {
				keyBindingFn(s, &(s.matchingChoices[s.selectedIndex]))
			} else {
				keyBindingFn(s, nil)
			}
		}

		if !s.hidden {
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
		fields := strings.Fields(input)
		matches := false
		if len(fields) == 0 {
			matches = true
		} else {
			allMatches := true
			for _, field := range fields {
				if !strings.Contains(strings.ToLower(s.getValue(&choice)), strings.ToLower(field)) {
					allMatches = false
				}
			}
			matches = allMatches
		}

		if matches {
			matchingChoices = append(matchingChoices, choice)
		}
	}
	return matchingChoices
}

func (s *Selection) FlashMessage(m string) {
	s.message = m
	s.render()

	go func() {
		time.Sleep(1 * time.Second)
		s.message = ""
		s.render()
	}()
}

func (s *Selection) render() {

	Printer{}.SaveCursor().HideCursor().CursorStartOfLine().ClearLine()
	Printer{}.Blue(s.Prompt).Print(s.input).NewLine()

	if len(s.matchingChoices) == 0 {
		Printer{}.ClearLine().Blue("no matches").NewLine()
	}

	for i := 0; i < len(s.matchingChoices); i++ {
		choice := s.matchingChoices[i]
		if i >= s.viewStartIndex && i < s.viewStartIndex+s.PageSize {
			value := s.getValue(&choice)
			Printer{}.ClearLine()
			if i == s.selectedIndex {
				Printer{}.BgWhite(formatMatches(value, strings.Fields(s.input), "\x1b[36m", "\x1b[30;47m")).NewLine()
			} else {
				Printer{}.Print(formatMatches(value, strings.Split(s.input, " "), "\x1b[36m", "\x1b[0m")).NewLine()
			}
		}
	}

	if s.Footer != "" {
		Printer{}.ClearLine().Yellow(s.Footer).NewLine()
	}

	if s.message != "" {
		Printer{}.ClearLine().Green(s.message).NewLine()
	}

	Printer{}.ClearScreen().RestoreCursor().CursorStartOfLine().CursorRight(s.cursorIndex+len(s.Prompt))
	Printer{}.ShowCursor()
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
