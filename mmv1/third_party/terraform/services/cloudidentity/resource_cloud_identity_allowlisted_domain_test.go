package cloudidentity

import (
	"errors"
	"testing"
)

func TestCloudIdentityAllowlistedDomainPollCheckExistence(t *testing.T) {
	t.Parallel()

	pollErr := errors.New("list failed")

	cases := []struct {
		name          string
		domain        string
		resp          map[string]interface{}
		respErr       error
		wantSuccess   bool
		wantRetryable bool
		wantErr       error
		wantName      string
	}{
		{
			name:          "non-retryable error on poll failure",
			domain:        "example.com",
			respErr:       pollErr,
			wantRetryable: false,
			wantErr:       pollErr,
		},
		{
			name:          "empty response map is pending",
			domain:        "example.com",
			resp:          map[string]interface{}{},
			wantRetryable: true,
		},
		{
			name:   "empty allowlistedDomains slice is pending",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{},
			},
			wantRetryable: true,
		},
		{
			name:   "other domain only is pending",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/other",
						"domain": "other.example.com",
					},
				},
			},
			wantRetryable: true,
		},
		{
			name:   "matching domain without name is pending",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"domain": "example.com",
					},
				},
			},
			wantRetryable: true,
		},
		{
			name:   "exact lowercase match succeeds and captures name",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/12345",
						"domain": "example.com",
					},
				},
			},
			wantSuccess: true,
			wantName:    "allowlistedDomains/12345",
		},
		{
			name:   "case-insensitive domain match succeeds and captures name",
			domain: "Example.COM",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/67890",
						"domain": "example.com",
					},
				},
			},
			wantSuccess: true,
			wantName:    "allowlistedDomains/67890",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotName string
			res := cloudIdentityAllowlistedDomainPollCheckExistence(tc.domain, &gotName)(tc.resp, tc.respErr)
			if tc.wantSuccess {
				if res != nil {
					t.Fatalf("expected nil PollResult (success), got retryable=%v err=%v", res.Retryable, res.Err)
				}
				if gotName != tc.wantName {
					t.Fatalf("got name %q, want %q", gotName, tc.wantName)
				}
				return
			}

			if res == nil {
				t.Fatalf("expected non-nil PollResult, got nil")
			}
			if res.Retryable != tc.wantRetryable {
				t.Fatalf("got Retryable=%v, want %v (err=%v)", res.Retryable, tc.wantRetryable, res.Err)
			}
			if tc.wantErr != nil && !errors.Is(res.Err, tc.wantErr) {
				t.Fatalf("got err %v, want %v", res.Err, tc.wantErr)
			}
			if gotName != "" {
				t.Fatalf("expected empty name on non-success, got %q", gotName)
			}
		})
	}
}

func TestCloudIdentityAllowlistedDomainPollCheckAbsence(t *testing.T) {
	t.Parallel()

	pollErr := errors.New("list failed")

	cases := []struct {
		name          string
		domain        string
		resp          map[string]interface{}
		respErr       error
		wantSuccess   bool
		wantRetryable bool
		wantErr       error
	}{
		{
			name:          "non-retryable error on poll failure",
			domain:        "example.com",
			respErr:       pollErr,
			wantRetryable: false,
			wantErr:       pollErr,
		},
		{
			name:        "empty response map succeeds",
			domain:      "example.com",
			resp:        map[string]interface{}{},
			wantSuccess: true,
		},
		{
			name:   "empty allowlistedDomains slice succeeds",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{},
			},
			wantSuccess: true,
		},
		{
			name:   "other domain only succeeds",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/other",
						"domain": "other.example.com",
					},
				},
			},
			wantSuccess: true,
		},
		{
			name:   "matching domain still listed is pending",
			domain: "example.com",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/12345",
						"domain": "example.com",
					},
				},
			},
			wantRetryable: true,
		},
		{
			name:   "case-insensitive matching domain still listed is pending",
			domain: "Example.COM",
			resp: map[string]interface{}{
				"allowlistedDomains": []interface{}{
					map[string]interface{}{
						"name":   "allowlistedDomains/12345",
						"domain": "example.com",
					},
				},
			},
			wantRetryable: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := cloudIdentityAllowlistedDomainPollCheckAbsence(tc.domain)(tc.resp, tc.respErr)
			if tc.wantSuccess {
				if res != nil {
					t.Fatalf("expected nil PollResult (success), got retryable=%v err=%v", res.Retryable, res.Err)
				}
				return
			}

			if res == nil {
				t.Fatalf("expected non-nil PollResult, got nil")
			}
			if res.Retryable != tc.wantRetryable {
				t.Fatalf("got Retryable=%v, want %v (err=%v)", res.Retryable, tc.wantRetryable, res.Err)
			}
			if tc.wantErr != nil && !errors.Is(res.Err, tc.wantErr) {
				t.Fatalf("got err %v, want %v", res.Err, tc.wantErr)
			}
		})
	}
}
