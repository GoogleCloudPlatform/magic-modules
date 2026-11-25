package bigquery

import (
	"net/http"
	"testing"
)

func TestSetBigQueryDatasetEtagHeader(t *testing.T) {
	for _, testCase := range []struct {
		name string
		etag string
		want string
	}{
		{name: "present", etag: "etag-value", want: "etag-value"},
		{name: "missing", etag: ""},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			headers := make(http.Header)
			setBigQueryDatasetEtagHeader(headers, testCase.etag)
			if got := headers.Get("If-Match"); got != testCase.want {
				t.Errorf("If-Match = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestAccessToPolicyForIamMemberPreservesEtag(t *testing.T) {
	policy, err := accessToPolicyForIamMember([]interface{}{}, "etag-value")
	if err != nil {
		t.Fatalf("accessToPolicyForIamMember() error = %v", err)
	}
	if got, want := policy.Etag, "etag-value"; got != want {
		t.Errorf("policy ETag = %q, want %q", got, want)
	}
}
