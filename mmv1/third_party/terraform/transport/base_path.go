package transport

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/registry"
)

// Returns the base path for a resource taking into account the following rules:
// Overridden path takes precedence over everything
// Regional endpoint should be returned if preferred
// Global endpoint should be returned if preferred
// If no preferences, return the product default based on DefaultRepStatus map
func BasePath(product registry.Product, config *Config, location string) (string, error) {
	var path string
	if v := config.CustomEndpoints[product.CustomEndpointField]; v != "" {
		path = v
	} else if config.PreferRegionalEndpoints || (DefaultRepStatus[product.Name] && !config.PreferGlobalEndpoints) {
		path = product.RepUrl
	} else {
		if config.IsMtls {
			path = GetMtlsEndpoint(product.BaseUrl)
		}
		if config.UniverseDomain != "" && config.UniverseDomain != "googleapis.com" {
			path = strings.ReplaceAll(path, "googleapis.com", config.UniverseDomain)
		}
	}

	if strings.Contains(path, "{{location}}") && location == "" {
		log.Printf("[WARN] Found base path with location but no location provided: %s", path)
		return path, fmt.Errorf("failed to find location for a resource with a regionalized endpoint %s", path)
	}
	// Still try to replace location even if it may not exist, this allows
	// for products that only support REP on their base_url
	return strings.ReplaceAll(path, "{{location}}", location), nil
}
