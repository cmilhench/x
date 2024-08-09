package kata

import "testing"

func Test_IsPalindrome(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		wantOut bool
	}{
		{"0 is a palindrome", 0, true},
		{"1 is a palindrome", 1, true},
		{"121 is a palindrome", 121, true},
		{"123 is not a palindrome", 123, false},
		{"-121 is not a palindrome", -121, false},
		{"10 is not a palindrome", 10, false},
		{"-101 is not a palindrome", -101, false},
		{"123454321 is a palindrome", 123454321, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := IsPalindromeNumber(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}

func Test_IsPalindromeText(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantOut bool
	}{
		{"A is a palindrome", "A", true},
		{"AA is a palindrome", "AA", true},
		{"AB is not a palindrome", "AB", false},
		{"A A is not a palindrome", "A A", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := IsPalindromeText(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}

// ----------------------------------------------------------------------------

func IsPalindromeText(text string) bool {
	for i, j := 0, len(text)-1; i < j; i, j = i+1, j-1 {
		if text[i] != text[j] {
			return false
		}
	}
	return true
}

func IsPalindromeNumber(num int) bool {
	var r, d int
	// 1234 % 10 = 4, 1234 / 10 = 123
	//  123 % 10 = 3,  123 / 10 = 12
	//   12 % 10 = 2,  12  / 10 = 1
	//    1 % 10 = 1,   1  / 10 = 0
	for i := num; i > 0; {
		d, i = i%10, i/10
		r = r*10 + d
	}

	return num == r
}
