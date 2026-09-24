package service_test

import (
	"testing"
)

// A table-driven unit test that carries the TestAcc prefix. It ranges over a
// map just like a serial test does, but the map holds test cases rather than
// test functions, so there is no config to read.
func TestAccTableDriven(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Expected string
	}{
		"first case": {
			Input:    "value-one",
			Expected: "value-one",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.Input != tc.Expected {
				t.Errorf("got %s, expected %s", tc.Input, tc.Expected)
			}
		})
	}
}
