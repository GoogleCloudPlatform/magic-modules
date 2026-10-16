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
		})
	}
}
