package input

import (
	"fmt"
	"strconv"
)

type Printer struct {
}

var (
	ansiReset = "\x1b[0m"
)

func (p Printer) CursorStartOfLine() Printer {
	fmt.Print("\x1b[1000D")
	return p
}

func (p Printer) ClearScreen() Printer {
	fmt.Print("\x1b[0J")
	return p
}

func (p Printer) SaveCursor() Printer {
	fmt.Print("\x1b[s")
	return p
}

func (p Printer) RestoreCursor() Printer {
	fmt.Print("\x1b[u")
	return p
}

func (p Printer) CursorRight(chars int) Printer {
	fmt.Print("\x1b[" + strconv.Itoa(chars) + "C")
	return p
}

func (p Printer) HideCursor() Printer {
	fmt.Print("\x1B[?25l")
	return p
}

func (p Printer) ShowCursor() Printer {
	fmt.Print("\x1B[?25h")
	return p
}

func (p Printer) ClearLine() Printer {
	fmt.Print("\x1b[K")
	return p
}

func (p Printer) Blue(s string) Printer {
	fmt.Print("\x1b[1;34m" + s + ansiReset)
	return p
}

func (p Printer) Yellow(s string) Printer {
	fmt.Print("\x1b[1;33m" + s + ansiReset)
	return p
}

func (p Printer) Green(s string) Printer {
	fmt.Print("\x1b[1;32m" + s + ansiReset)
	return p
}

func (p Printer) Gray(s string) Printer {
	fmt.Print("\x1b[1;30;1m" + s + ansiReset)
	return p
}

func (p Printer) BgWhite(s string) Printer {
	fmt.Print("\x1b[30;47m" + s + ansiReset)
	return p
}

func (p Printer) Print(s string) Printer {
	fmt.Print(s)
	return p
}

func (p Printer) NewLine() Printer {
	fmt.Println()
	return p
}
