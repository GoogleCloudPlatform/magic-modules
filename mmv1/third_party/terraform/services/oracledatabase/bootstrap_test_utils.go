package oracledatabase

import (
	"fmt"
	"log"
	"testing"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// BootstrapOdbDnsZone returns the name of a long-lived OdbDnsZone in the given
// OdbNetwork, creating it with the given domain if it doesn't exist.
// OdbDnsZone has no Terraform resource, so tests that reference one (such as
// the `odb_dns_zone` field of `google_oracle_database_cloud_vm_cluster`) use
// this shared zone instead of creating their own.
func BootstrapOdbDnsZone(t *testing.T, project, location, odbNetworkId, odbDnsZoneId, domain string) string {
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap config.")
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/odbNetworks/%s", project, location, odbNetworkId)
	name := fmt.Sprintf("%s/odbDnsZones/%s", parent, odbDnsZoneId)
	baseUrl := transport_tpg.BaseUrl(Product, config)

	log.Printf("[DEBUG] Getting OdbDnsZone %q", name)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    baseUrl + name,
		UserAgent: config.UserAgent,
	})
	if err == nil {
		return name
	}
	if !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		t.Fatalf("Error getting OdbDnsZone %q: %s", name, err)
	}

	log.Printf("[DEBUG] OdbDnsZone %q not found, bootstrapping", name)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    fmt.Sprintf("%s%s/odbDnsZones?odbDnsZoneId=%s", baseUrl, parent, odbDnsZoneId),
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"domain": domain,
		},
		Timeout: 20 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Error creating OdbDnsZone %q: %s", name, err)
	}
	if err := OracleDatabaseOperationWaitTime(config, res, project, "Creating OdbDnsZone", config.UserAgent, 20*time.Minute); err != nil {
		t.Fatalf("Error waiting to create OdbDnsZone %q: %s", name, err)
	}
	log.Printf("[DEBUG] Finished creating OdbDnsZone %q", name)

	return name
}
