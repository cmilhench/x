package humanize

import "testing"

func Test_Comma(t *testing.T) {
	tests := []struct {
		name    string
		in      int64
		wantOut string
	}{
		{"digit", 5, "5"},
		{"hundred", 100, "100"},
		{"hundred", 999, "999"},
		{"thousand", 2345, "2,345"},
		{"negatives", -1234, "-1,234"},
		{"humongous", 12345678, "12,345,678"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := Comma(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("Comma(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}
