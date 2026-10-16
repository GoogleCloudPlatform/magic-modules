package iambeta

import (
	"fmt"
	"net/url"
)

// validateWorkloadIdentityPoolURL requires raw to be an https URL with no
// embedded userinfo.
//
// The openid_config and jwks data sources fetch the issuer's OIDC discovery
// document and its public keys, which are Workload Identity Federation trust
// material, and both the OIDC Discovery spec and these fields' own contracts
// require the https scheme. resource_name is commonly wired from the Computed
// jwks_uri of the openid_config data source, an unknown value at plan time that
// skips the schema-level ValidateFunc, so the scheme is enforced here at fetch
// time as well. Rejecting userinfo avoids a "https://trusted@attacker/" host
// confusion where the authority that gets contacted differs from the one an
// operator reading the config would expect.
func validateWorkloadIdentityPoolURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("resource_name %q is not a valid URL: %w", raw, err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("resource_name must use the https scheme, got %q", u.Scheme)
	}
	if u.User != nil {
		return fmt.Errorf("resource_name must not embed userinfo credentials")
	}
	return nil
}
