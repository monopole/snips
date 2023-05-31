package query

import "testing"

func Test_upToFirstLfOrEnd(t *testing.T) {
	tests := map[string]struct {
		value string
		want  string
	}{
		"t1": {
			value: "the rain in spain",
			want:  "the rain in spain",
		},
		"t2": {
			value: `the rain in spain
falls mainly on the plain
`,
			want: "the rain in spain",
		},
	}
	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			if got := upToFirstLfOrEnd(tt.value); got != tt.want {
				t.Errorf("upToFirstLfOrEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}
