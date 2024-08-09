package irc

import (
	"fmt"
	"testing"
)

func TestMessage(t *testing.T) {
	var tests = []struct {
		name string
		line string
		Message
	}{
		{"1", ":example.freenode.net NOTICE * :*** Looking up your hostname...\r\n", Message{"example.freenode.net", "NOTICE", "*", "*** Looking up your hostname...", ""}},
		{"2", "ERROR :Closing Link: 127.0.0.1 (Connection timed out)\r\n", Message{"", "ERROR", "", "Closing Link: 127.0.0.1 (Connection timed out)", ""}},
		{"3", ":user!~mail@example.net JOIN #channel\r\n", Message{"user!~mail@example.net", "JOIN", "#channel", "", ""}},
		{"4", ":user!~mail@example.com PRIVMSG user :Hello :)\r\n", Message{"user!~mail@example.com", "PRIVMSG", "user", "Hello :)", ""}},
		{"6", ":user!~mail@example.com PRIVMSG #channel :Hello :)\r\n", Message{"user!~mail@example.com", "PRIVMSG", "#channel", "Hello :)", ""}},
		{"6", ":NickServ!NickServ@services. NOTICE user :Some message.\r\n", Message{"NickServ!NickServ@services.", "NOTICE", "user", "Some message.", ""}},
		{"7", ":user PRIVMSG #chan :Hello!\r\n", Message{"user", "PRIVMSG", "#chan", "Hello!", ""}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("method%v", test.name), func(t *testing.T) {
			m := Message{}
			m.Parse(test.line)
			if m.Prefix != test.Prefix {
				t.Errorf("expected prefix '%s', got '%s'", test.Prefix, m.Prefix)
			}
			if m.Command != test.Command {
				t.Errorf("expected command '%s', got '%s'", test.Command, m.Command)
			}
			if m.Params != test.Params {
				t.Errorf("expected params '%s', got '%s'", test.Params, m.Params)
			}
			if m.Trailing != test.Trailing {
				t.Errorf("expected trailing '%s', got '%s'", test.Trailing, m.Trailing)
			}
		})
	}
}
