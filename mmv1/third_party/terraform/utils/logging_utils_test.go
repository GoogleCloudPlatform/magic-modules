package google

import "testing"

func TestParseLoggingSinkParentId(t *testing.T) {
	tests := []struct {
		val         string
		out         string
		errExpected bool
	}{
		{"projects/my-project/sinks/my-sink", "my-project", false},
		{"folders/foofolder/sinks/woo", "foofolder", false},
		{"kitchens/the-big-one/sinks/second-from-the-left", "", true},
	}

	for _, test := range tests {
		out, err := parseLoggingSinkParentId(test.val)
		if err != nil {
			if !test.errExpected {
				t.Errorf("Got error with val %#v: error = %#v", test.val, err)
			}
		} else {
			if out != test.out {
				t.Errorf("Mismatch on val %#v: expected %#v but got %#v", test.val, test.out, out)
			}
		}
	}
}
