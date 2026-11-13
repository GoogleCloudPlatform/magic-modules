package iambeta

import "testing"

func TestValidateWorkloadIdentityPoolURL(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "https sts endpoint",
			raw:  "https://sts.googleapis.com/v1/organizations/123/locations/global/workloadIdentityPools/pool/openid/jwks",
		},
		{
			name: "https scheme is case-insensitive",
			raw:  "HTTPS://sts.googleapis.com/openid/jwks",
		},
		{
			name:    "plaintext http is rejected",
			raw:     "http://sts.googleapis.com/openid/jwks",
			wantErr: true,
		},
		{
			name:    "non-web scheme is rejected",
			raw:     "file:///etc/passwd",
			wantErr: true,
		},
		{
			name:    "embedded userinfo is rejected",
			raw:     "https://sts.googleapis.com@attacker.example/openid/jwks",
			wantErr: true,
		},
		{
			name:    "missing scheme is rejected",
			raw:     "sts.googleapis.com/openid/jwks",
			wantErr: true,
		},
		{
			name:    "missing host is rejected",
			raw:     "https:///openid/jwks",
			wantErr: true,
		},
		{
			name:    "empty string is rejected",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "control characters make the url unparseable",
			raw:     "https://sts.googleapis.com/\x7f",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateWorkloadIdentityPoolURL(tc.raw)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for %q, got nil", tc.raw)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.raw, err)
			}

			// The schema ValidateFunc must reach the same verdict as the Read-time check.
			_, errs := validateWorkloadIdentityPoolURLSchema(tc.raw, "resource_name")
			if (len(errs) > 0) != (err != nil) {
				t.Fatalf("schema validator disagrees with read-time check for %q: schema errs=%v, read err=%v", tc.raw, errs, err)
			}
		})
	}
}

func TestValidateWorkloadIdentityPoolURLSchema_nonString(t *testing.T) {
	_, errs := validateWorkloadIdentityPoolURLSchema(42, "resource_name")
	if len(errs) != 1 {
		t.Fatalf("expected one error for non-string value, got %v", errs)
	}
}
