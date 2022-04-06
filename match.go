package input

import (
	"strings"
)

func formatMatches(value string, terms []string, matchStart string, matchEnd string) string {
	//fmt.Println(terms)
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
