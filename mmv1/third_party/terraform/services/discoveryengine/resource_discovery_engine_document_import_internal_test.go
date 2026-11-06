package discoveryengine

import (
	"encoding/json"
	"testing"
)

func TestDiscoveryEngineDocumentImportParseCount(t *testing.T) {
	cases := map[string]struct {
		in      interface{}
		wantN   int64
		wantStr string
		wantErr bool
	}{
		"nil (proto3 omits zero values)": {
			in:      nil,
			wantN:   0,
			wantStr: "0",
		},
		"empty string": {
			in:      "",
			wantN:   0,
			wantStr: "0",
		},
		"int64 serialized as string": {
			in:      "42",
			wantN:   42,
			wantStr: "42",
		},
		"string exceeding int32": {
			in:      "9000000000",
			wantN:   9000000000,
			wantStr: "9000000000",
		},
		"json.Number": {
			in:      json.Number("7"),
			wantN:   7,
			wantStr: "7",
		},
		"empty json.Number": {
			in:      json.Number(""),
			wantN:   0,
			wantStr: "0",
		},
		"malformed string": {
			in:      "abc",
			wantErr: true,
		},
		"non-integer json.Number": {
			in:      json.Number("1.5"),
			wantErr: true,
		},
		"unexpected type": {
			in:      true,
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			n, s, err := discoveryEngineDocumentImportParseCount(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %#v, got n=%d s=%q", tc.in, n, s)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %#v: %s", tc.in, err)
			}
			if n != tc.wantN || s != tc.wantStr {
				t.Errorf("got (%d, %q), want (%d, %q)", n, s, tc.wantN, tc.wantStr)
			}
		})
	}
}

func TestDiscoveryEngineDocumentImportFormatErrorSamples(t *testing.T) {
	cases := map[string]struct {
		samples []interface{}
		limit   int
		want    string
	}{
		"nil samples": {
			samples: nil,
			limit:   3,
			want:    "",
		},
		"empty samples": {
			samples: []interface{}{},
			limit:   3,
			want:    "",
		},
		"zero limit": {
			samples: []interface{}{map[string]interface{}{"code": json.Number("3"), "message": "bad"}},
			limit:   0,
			want:    "",
		},
		"negative limit": {
			samples: []interface{}{map[string]interface{}{"code": json.Number("3"), "message": "bad"}},
			limit:   -1,
			want:    "",
		},
		"code and message": {
			samples: []interface{}{map[string]interface{}{"code": json.Number("3"), "message": "invalid JSON"}},
			limit:   3,
			want:    "[code=3] invalid JSON",
		},
		"message only": {
			samples: []interface{}{map[string]interface{}{"message": "invalid JSON"}},
			limit:   3,
			want:    "invalid JSON",
		},
		"map without message": {
			samples: []interface{}{map[string]interface{}{"code": json.Number("3")}},
			limit:   3,
			want:    "map[code:3]",
		},
		"non-map sample": {
			samples: []interface{}{"raw error"},
			limit:   3,
			want:    "raw error",
		},
		"exactly limit samples": {
			samples: []interface{}{
				map[string]interface{}{"message": "a"},
				map[string]interface{}{"message": "b"},
			},
			limit: 2,
			want:  "a; b",
		},
		"truncated beyond limit": {
			samples: []interface{}{
				map[string]interface{}{"code": json.Number("3"), "message": "a"},
				map[string]interface{}{"code": json.Number("3"), "message": "b"},
				map[string]interface{}{"code": json.Number("3"), "message": "c"},
				map[string]interface{}{"code": json.Number("3"), "message": "d"},
				map[string]interface{}{"code": json.Number("3"), "message": "e"},
			},
			limit: 3,
			want:  "[code=3] a; [code=3] b; [code=3] c; ... (2 more error sample(s))",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := discoveryEngineDocumentImportFormatErrorSamples(tc.samples, tc.limit); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
