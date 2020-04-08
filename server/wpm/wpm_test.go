package wpm

import "testing"

func TestWildcardPatternMatch(t *testing.T) {

	type args struct {
		str     string
		pattern string
	}
	tests := []struct {
		name string
		in   args
		out  bool
	}{
		{
			name: "with * passing",
			in: args{str: "nerico",
				pattern: "*eric*",
			},
			out: true,
		},
		{
			name: "without * passing",
			in: args{str: "nerico",
				pattern: "nerico",
			},
			out: true,
		},
		{
			name: "with * failing",
			in: args{str: "nerico",
				pattern: "*ricos",
			},
			out: false,
		},
		{
			name: "without * failing",
			in: args{str: "nerico",
				pattern: "nerim",
			},
			out: false,
		},
		{
			name: "empty pattern match with empty string",
			in: args{str: "",
				pattern: "",
			},
			out: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWildcardPattern(tt.in.str, tt.in.pattern).Match(); got != tt.out {
				t.Errorf("WildcardPattern().Match = %v, want %v", got, tt.out)
			}
		})
	}
}
