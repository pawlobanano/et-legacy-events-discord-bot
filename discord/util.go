package discord

import (
	"strconv"
	"strings"
)

// multilineString represents a multiline string.
type multilineString struct {
	Lines []string
}

// newMultilineString creates a new MultilineString.
func newMultilineString(lines ...string) *multilineString {
	return &multilineString{Lines: lines}
}

// format returns the formatted multiline string.
func (ms *multilineString) format() string {
	return strings.Join(ms.Lines, "\n")
}

// append appends lines to the MultilineString.
func (ms *multilineString) append(lines ...string) {
	ms.Lines = append(ms.Lines, lines...)
}

// extractSuffixNumber extracts suffix number in a string and return it as integer.
func extractSuffixNumber(input string) (int, error) {
	// Split the input string by spaces.
	parts := strings.Fields(input)

	if len(parts) >= 3 {
		number := parts[2]

		num, err := strconv.Atoi(number)
		if err != nil {
			num = 0
			return num, err
		}

		return num, nil
	}

	return 0, nil
}
