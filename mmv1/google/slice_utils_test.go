package google

import (
	"reflect"
	"testing"
)

func TestSliceSelect(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		S           []int
		testFun     func(int) bool
		expected    int
	}{
		{
			description: "interger slice selects even numbers",
			S:           []int{0, 1, 2},
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 2,
		},
		{
			description: "empty slice",
			S:           make([]int, 0),
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 0,
		},
		{
			description: "nil slice",
			S:           nil,
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := len(Select(tc.S, tc.testFun)), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestSliceReject(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		S           []int
		testFun     func(int) bool
		expected    int
	}{
		{
			description: "interger slice rejects even numbers",
			S:           []int{0, 1, 2},
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 1,
		},
		{
			description: "empty slice",
			S:           make([]int, 0),
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 0,
		},
		{
			description: "nil slice",
			S:           nil,
			testFun: func(n int) bool {
				return n%2 == 0
			},
			expected: 0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := len(Reject(tc.S, tc.testFun)), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestSliceConcat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		S1          []int
		S2          []int
		expected    int
	}{
		{
			description: "interger slice rejects even numbers",
			S1:          []int{0, 1, 2},
			S2:          []int{3, 4},
			expected:    5,
		},
		{
			description: "empty slice",
			S1:          nil,
			S2:          make([]int, 0),
			expected:    0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := len(Concat(tc.S1, tc.S2)), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}
