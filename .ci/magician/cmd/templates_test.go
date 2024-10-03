package cmd

import (
	"testing"
)

func TestColor(t *testing.T) {
	cases := []struct {
		name  string
		color string
		text  string
		want  string
	}{
		{
			name:  "red",
			color: "red",
			text:  "Test text",
			want:  "🔴 Test text",
		},
		{
			name:  "yellow",
			color: "yellow",
			text:  "Test text",
			want:  "🟡 Test text",
		},
		{
			name: "green",
			color: "green",
			text: "Test text",
			want: "🟢 Test text",
		},
		{
			name: "unsupported color",
			color: "mauve",
			text: "Test text",
			want: "Test text",
		},
		{
			name: "empty color",
			text: "Test text",
			want: "Test text",
		},
		{
			name:  "empty text",
			color: "green",
			want:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := color(tc.color, tc.text)
			if got != tc.want {
				t.Errorf("color(%s, %s) got %s; want %s", tc.color, tc.text, got, tc.want)
			}
		})
	}
}
