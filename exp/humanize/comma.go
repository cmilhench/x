package humanize

import (
	"fmt"
	"strconv"
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

	value = append([]string{strconv.FormatInt(x, 10)}, value...)
	return sign + strings.Join(value, ",")
}
