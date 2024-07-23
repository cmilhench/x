package log

import (
	"fmt"
	"strings"
)

type Level int

const (
	Fatal Level = iota
	Error
	Info
	Debug
)

// String returns the string representation of the level.
func (l Level) String() string {
	var names = []string{
		"Fatal",
		"Error",
		"Info",
		"Debug",
	}
	if Fatal <= l && l <= Debug {
		return names[l]
	}
	return fmt.Sprintf("%%!Level(%d)", l)
}

// Parse a level string, case-insensitive.
func (l *Level) Parse(s string) {
	switch strings.ToLower(s) {
	case "fatal":
		*l = Fatal
	case "error":
		*l = Error
	case "info":
		*l = Info
	case "debug":
		*l = Debug
	default:
		*l = Error
	}
}
