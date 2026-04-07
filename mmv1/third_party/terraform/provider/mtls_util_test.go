package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestUnitMtls_urlSwitching(t *testing.T) {
	t.Parallel()
	for _, p := range registry.ListProducts() {
		url = getMtlsEndpoint(p.BaseUrl)
		if !strings.Contains(url, ".mtls.") {
			t.Errorf("%s: mtls conversion unsuccessful preconv - %s postconv - %s", key, bp, url)
		}
	}
	for key, bp := range transport_tpg.DefaultBasePaths {
		url := getMtlsEndpoint(bp)
		if !strings.Contains(url, ".mtls.") {
			t.Errorf("%s: mtls conversion unsuccessful preconv - %s postconv - %s", key, bp, url)
		}
	}
}
