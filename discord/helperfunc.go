package discord

import (
	"strings"
)

// MultilineString represents a multiline string.
type MultilineString struct {
	Lines []string
}

// NewMultilineString creates a new MultilineString.
func NewMultilineString(lines ...string) *MultilineString {
	return &MultilineString{Lines: lines}
}

// Format returns the formatted multiline string.
func (ms *MultilineString) Format() string {
	return strings.Join(ms.Lines, "\n")
}

// Append appends lines to the MultilineString.
func (ms *MultilineString) Append(lines ...string) {
	ms.Lines = append(ms.Lines, lines...)
}
