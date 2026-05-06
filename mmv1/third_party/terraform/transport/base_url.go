package transport

import (
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/registry"
)

// Returns the base URL for a product taking into account the following rules:
// 1. If there is a custom endpoint set, return that immediately.
// 2. Otherwise, determine whether to use the REP url or standard url.
// 3. Make adjustments for mTLS / universe domain.
// 4. Return final URL.
func BaseUrl(product registry.Product, config *Config) string {
	if v := config.CustomEndpoints[product.CustomEndpointField]; v != "" {
		return v
	}

	u := product.BaseUrl
	if config.PreferRegionalEndpoints && product.RepUrl != "" {
		u = product.RepUrl
	} else if config.PreferGlobalEndpoints {
		u = product.BaseUrl
	} else if product.RepByDefault && product.RepUrl != "" {
		u = product.RepUrl
	}

	if config.IsMtls {
		u = GetMtlsEndpoint(u)
	}
	if config.UniverseDomain != "" && config.UniverseDomain != "googleapis.com" {
		u = strings.ReplaceAll(u, "googleapis.com", config.UniverseDomain)
	}

	return u
}
