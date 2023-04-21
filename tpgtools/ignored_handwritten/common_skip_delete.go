package google

import (
	"context"
	"log"

	dns "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dns"
)

// Skip delete for DNS record set if the record is for the primary NS record
func rrefSkipDelete(c *Config, recordSet *dns.ResourceRecordSet) (bool, error) {
	if *recordSet.DnsType != "NS" {
		// Only skip for NS records in some circumstances
		return false, nil
	}

	mz := &dns.ManagedZone{
		Name:    recordSet.ManagedZone,
		Project: recordSet.Project,
	}

	res, err := c.clientDnsDCL.GetManagedZone(context.Background(), mz)
	if err != nil {
		return false, err
	}

	// Subdomains can be deleted, so check if this is one
	if *res.DnsName == *recordSet.DnsName {
		log.Println("[DEBUG] NS records can't be deleted due to API restrictions, so they're being left in place. See https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/dns_record_set for more information.")
		return true, nil
	}
	return false, nil
}
