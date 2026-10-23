package detector

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDetectUnreadableTests(t *testing.T) {
	cases := map[string]struct {
		testNames []string
		readErrs  map[string]error
		known     []string
		want      UnreadableTestsReport
	}{
		"no tests": {
			want: UnreadableTestsReport{
				NewUnreadableTests:   map[string]string{},
				StaleUnreadableTests: []string{},
			},
		},
		"known unreadable tests are not reported": {
			testNames: []string{"TestAccReadable", "TestAccKnown"},
			readErrs:  map[string]error{"TestAccKnown": errors.New("failed")},
			known:     []string{"TestAccKnown"},
			want: UnreadableTestsReport{
				NewUnreadableTests:   map[string]string{},
				StaleUnreadableTests: []string{},
			},
		},
		"unknown unreadable tests are reported with their errors": {
			readErrs: map[string]error{
				"TestAccKnown": errors.New("failed"),
				"TestAccNew":   errors.New("failed to find variable declaration config"),
			},
			known: []string{"TestAccKnown"},
			want: UnreadableTestsReport{
				NewUnreadableTests:   map[string]string{"TestAccNew": "failed to find variable declaration config"},
				StaleUnreadableTests: []string{},
			},
		},
		"known tests that were read are reported as stale, sorted": {
			testNames: []string{"TestAccKnown", "TestAccReadableB", "TestAccOther", "TestAccReadableA"},
			readErrs:  map[string]error{"TestAccKnown": errors.New("failed")},
			known:     []string{"TestAccKnown", "TestAccReadableA", "TestAccReadableB", "TestAccElsewhere"},
			want: UnreadableTestsReport{
				NewUnreadableTests:   map[string]string{},
				StaleUnreadableTests: []string{"TestAccReadableA", "TestAccReadableB"},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := DetectUnreadableTests(tc.testNames, tc.readErrs, tc.known)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("DetectUnreadableTests() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDetectUnreadableTestsTruncatesErrors(t *testing.T) {
	long := "errors reading test steps:\n\t" + strings.Repeat("x", 2*maxUnreadableTestErrorLength)
	got := DetectUnreadableTests(nil, map[string]error{"TestAccLong": errors.New(long)}, nil).NewUnreadableTests["TestAccLong"]
	if !strings.HasPrefix(got, "errors reading test steps: xxx") {
		t.Errorf("error whitespace was not collapsed: %q", got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("truncated error should end with an ellipsis: %q", got)
	}
	if n := len([]rune(got)); n != maxUnreadableTestErrorLength+1 {
		t.Errorf("truncated error has %d runes, want %d", n, maxUnreadableTestErrorLength+1)
	}
}
