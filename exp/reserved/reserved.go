package reserved

import (
	"bufio"
	"bytes"
	_ "embed" // used to embed reserved words and patterns at compile time
	"regexp"
)

//go:embed reserved.txt
var data []byte
var reserved = map[string]struct{}{}
var predicates = make([]func(string) bool, 0, 20)

func init() {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		tx := scanner.Text()
		if tx[0] == '#' {
			continue
		}
		if tx[0] == '^' {
			predicate := regexp.MustCompile(tx).MatchString
			predicates = append(predicates, predicate)

			continue
		}
		reserved[tx] = struct{}{}
	}
}

// IsReserved checks weather the provided value could be considered a reserved
// subdomain matching words and patterns such as; auth, help, smtp, www3, etc.
func IsReserved(value string) bool {
	if _, ok := reserved[value]; ok {
		return true
	}
	for _, p := range predicates {
		if p(value) {
			return true
		}
	}
	return false
}
