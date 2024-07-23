package humanize

import "testing"

func Test_Ordinal(t *testing.T) {
	tests := []struct {
		name    string
		in      int
		wantOut string
	}{
		{"1st", 1, "1st"},
		{"2nd", 2, "2nd"},
		{"3rd", 3, "3rd"},
		{"4th", 4, "4th"},
		{"11th", 11, "11th"},
		{"12th", 12, "12th"},
		{"13th", 13, "13th"},
		{"14th", 14, "14th"},
		{"111th", 111, "111th"},
		{"112th", 112, "112th"},
		{"113th", 113, "113th"},
		{"121th", 121, "121st"},
		{"122th", 122, "122nd"},
		{"123th", 123, "123rd"},
		{"1111th", 1111, "1111th"},
		{"1112th", 1112, "1112th"},
		{"1113th", 1113, "1113th"},
		{"1121th", 1121, "1121st"},
		{"1122th", 1122, "1122nd"},
		{"1123th", 1123, "1123rd"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOut := Ordinal(tc.in)
			if gotOut != tc.wantOut {
				t.Errorf("Ordinal(%v) = %v, want %v", tc.in, gotOut, tc.wantOut)
			}
		})
	}
}
