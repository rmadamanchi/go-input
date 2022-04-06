package input

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	formatted := formatMatches("my name is bob", []string{"is", "name", "bob", "na", "foo", "b"}, "<<", ">>")
	fmt.Println(formatted)
}
