package kata_test

import (
	"fmt"
	"strings"
	"testing"
)

func Test_RomanNumeral(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		wantOut string
	}{
		{"0001", 1, "I"},
		{"0004", 4, "IV"},
		{"0005", 5, "V"},
		{"0009", 9, "IX"},
		{"0010", 10, "X"},
		{"0039", 39, "XXXIX"},
		{"0246", 246, "CCXLVI"},
		{"0789", 789, "DCCLXXXIX"},
		{"2421", 2421, "MMCDXXI"},
		{"0160", 160, "CLX"},
		{"0027", 207, "CCVII"},
		{"1009", 1009, "MIX"},
		{"1066", 1066, "MLXVI"},
		{"3999", 3999, "MMMCMXCIX"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := RomanNumeral(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("Comma(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}

func Test_ParseRomanNumeral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		in      string
		wantOut int
	}{
		{"0001", "I", 1},
		{"0004", "IV", 4},
		{"0005", "V", 5},
		{"0009", "IX", 9},
		{"0010", "X", 10},
		{"0039", "XXXIX", 39},
		{"0246", "CCXLVI", 246},
		{"0789", "DCCLXXXIX", 789},
		{"2421", "MMCDXXI", 2421},
		{"0160", "CLX", 160},
		{"0027", "CCVII", 207},
		{"1009", "MIX", 1009},
		{"1066", "MLXVI", 1066},
		{"3999", "MMMCMXCIX", 3999},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := ParseRomanNumeral(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("Comma(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}

// ----------------------------------------------------------------------------

// 1 5 10 50 100 500 1000
// I V X  L  C   D   M
// I can be placed before V (5) and X (10) to make 4 and 9.
// X can be placed before L (50) and C (100) to make 40 and 90.
// C can be placed before D (500) and M (1000) to make 400 and 900.
// 3999 = MMMCMXCIX

func RomanNumeral1(num int) string {
	if 0 > num || num > 3999 {
		return ""
	}
	sym := "IVXLCDM"
	out := ""
	for idx := 1; num > 0; {
		r := num % 10
		switch r {
		case 1, 2, 3:
			I := sym[idx-1]
			out = fmt.Sprintf("%s%s", strings.Repeat(string(I), r), out)
		case 4:
			I, V := sym[idx-1], sym[idx]
			out = fmt.Sprintf("%s%s%s", string(I), string(V), out)
		case 5:
			V := sym[idx]
			out = fmt.Sprintf("%s%s", string(V), out)
		case 6, 7, 8:
			I, V := sym[idx-1], sym[idx]
			out = fmt.Sprintf("%s%s%s", string(V), strings.Repeat(string(I), r-5), out)
		case 9:
			I, X := sym[idx-1], sym[idx+1]
			out = fmt.Sprintf("%s%s%s", string(I), string(X), out)
		default:
		}
		num = num / 10
		idx = idx + 2
	}

	return out
}

func RomanNumeral(num int) string {
	val := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	sym := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	out := ""

	for i := 0; i < len(val); i++ {
		for num >= val[i] {
			out += sym[i]
			num -= val[i]
		}
	}
	return out
}

func ParseRomanNumeral(roman string) int {
	val := map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}
	out := 0
	prev := 0

	for i := len(roman) - 1; i >= 0; i-- {
		cur := val[roman[i]]
		if cur < prev {
			out -= cur
		} else {
			out += cur
		}
		prev = cur
	}
	return out
}
