package apihub

import (
	"context"
	"testing"
)

func TestUnitApihubAttribute_isApihubAttributeCardinalityDecrease(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		old  interface{}
		new  interface{}
		want bool
	}{
		{name: "decrease", old: 3, new: 2, want: true},
		{name: "decrease to the minimum", old: 20, new: 1, want: true},
		{name: "increase", old: 1, new: 3, want: false},
		{name: "unchanged", old: 2, new: 2, want: false},
		{name: "create with no prior value", old: 0, new: 2, want: false},
		{name: "new value unset or unknown", old: 3, new: 0, want: false},
		{name: "nil values", old: nil, new: nil, want: false},
		{name: "unexpected types", old: "3", new: "2", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := isApihubAttributeCardinalityDecrease(context.Background(), tc.old, tc.new, nil); got != tc.want {
				t.Errorf("isApihubAttributeCardinalityDecrease(%v, %v) = %v, want %v", tc.old, tc.new, got, tc.want)
			}
		})
	}
}
