package log

import (
	"fmt"
	"log"
	"os"
)

var (
	stdOut   = log.New(os.Stdout, "", 0)
	stdErr   = log.New(os.Stderr, "[ERROR] => ", log.LstdFlags|log.Lshortfile)
	debugOut = log.New(os.Stdout, "[DEBUG] => ", log.LstdFlags|log.Lshortfile)
)

// Errorf writes the output for a logging event.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	stdErr.Printf(format, args...)
}

// Infof writes the output for a logging event.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	var level Level
	level.Parse(os.Getenv("LOG_LEVEL"))
	if level >= Info {
		stdOut.Printf(format, args...)
	}
}

// Debugf writes the output for a logging event.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	var level Level
	level.Parse(os.Getenv("LOG_LEVEL"))
	if level >= Debug {
		debugOut.Printf(format, args...)
	}
}

// Panicf writes the output for a logging event.
// Arguments are handled in the manner of fmt.Printf.
func Panicf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	stdErr.Printf(s)
	panic(s)
}
