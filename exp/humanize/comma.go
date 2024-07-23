package humanize

import (
	"fmt"
	"strings"
)

// Comma formats an integer with commas separating every three digits.
func Comma(x int64) string {
	sign := ""
	value := []string{}
	if x < 0 {
		sign = "-"
		x = -x
	}
	for x > 999 {
		value = append([]string{fmt.Sprintf("%03d", (x % 1000))}, value...)
		x = x / 1000
	}
	value = append([]string{fmt.Sprintf("%d", (x))}, value...)
	return sign + strings.Join(value, ",")
}
