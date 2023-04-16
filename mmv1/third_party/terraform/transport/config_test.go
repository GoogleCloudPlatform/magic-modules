package transport

import "testing"

func TestGetRegionFromRegionSelfLink(t *testing.T) {
	cases := map[string]string{
		"https://www.googleapis.com/compute/v1/projects/test/regions/europe-west3": "europe-west3",
		"europe-west3": "europe-west3",
	}
	for input, expected := range cases {
		if result := GetRegionFromRegionSelfLink(input); result != expected {
			t.Errorf("expected to get %q from %q, got %q", expected, input, result)
		}
	}
}
