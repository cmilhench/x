package irc

import (
	"fmt"
	"strings"
)

type Message struct {
	Prefix   string
	Command  string
	Params   string
	Trailing string
	// _raw     string
}

func ParseMessage(line string) *Message {
	c := &Message{}
	c.Parse(line)
	return c
}

func (c *Message) Parse(line string) {
	line = strings.TrimSuffix(line, "\r")
	line = strings.TrimSuffix(line, "\r\n")
	// orig := line
	// c._raw = orig
	// Prefix
	if line[0] == ':' {
		i := strings.Index(line, " ")
		c.Prefix = line[1:i]
		line = line[i+1:]
	}
	// Command
	i := strings.Index(line, " ")
	if i == -1 {
		i = len(line)
	}
	c.Command = line[0:i]
	line = line[i:]
	// Params
	i = strings.Index(line, " :")
	if i == -1 {
		i = len(line)
	}
	if i != 0 {
		c.Params = line[1:i]
	}
	// Trailing
	if len(line)-i > 2 {
		c.Trailing = line[i+2:]
	}
}

func (c *Message) String() string {
	var line string
	if len(c.Prefix) > 0 {
		line = fmt.Sprintf(":%s ", c.Prefix)
	}
	line += c.Command
	if len(c.Params) > 0 {
		line = fmt.Sprintf("%s %s", line, c.Params)
	}
	if len(c.Trailing) > 0 {
		line = fmt.Sprintf("%s :%s", line, c.Trailing)
	}
	return line
}

func (c *Message) Nick() string {
	return c.Prefix[0:strings.Index(c.Prefix, "!")]
}

func (c *Message) Username() string {
	return c.Prefix[strings.Index(c.Prefix, "!")+1 : strings.Index(c.Prefix, "@")]
}

func (c *Message) Hostname() string {
	return c.Prefix[strings.Index(c.Prefix, "@")+1:]
}
