package transport

import "testing"

func TestNextListPageToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		res       map[string]interface{}
		wantToken string
		wantMore  bool
	}{
		"omitted": {
			res: map[string]interface{}{},
		},
		"null": {
			res: map[string]interface{}{"nextPageToken": nil},
		},
		"empty": {
			res: map[string]interface{}{"nextPageToken": ""},
		},
		"non-string": {
			res: map[string]interface{}{"nextPageToken": 1},
		},
		"present": {
			res:       map[string]interface{}{"nextPageToken": "abc"},
			wantToken: "abc",
			wantMore:  true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gotToken, gotMore := nextListPageToken(tc.res)
			if gotToken != tc.wantToken || gotMore != tc.wantMore {
				t.Fatalf("nextListPageToken() = (%q, %v), want (%q, %v)", gotToken, gotMore, tc.wantToken, tc.wantMore)
			}
		})
	}
}
