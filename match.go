package input

import (
	"strings"
)

// this is not a super robust matching algorithm but it works for simple purposes
func formatMatches(value string, terms []string, matchStart string, matchEnd string) string {
	formattedValue := value
	for _, term := range terms {
		matchIndex := strings.Index(strings.ToLower(formattedValue), strings.ToLower(term))

		if matchIndex != -1 {
			preMatchPart := formattedValue[0:matchIndex]
			matchPart := formattedValue[matchIndex : matchIndex+len(term)]
			postMatchPart := formattedValue[matchIndex+len(term):]
			formattedValue = preMatchPart + matchStart + matchPart + matchEnd + postMatchPart
		}
	}
	return formattedValue
}
