package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestIpv6AddressDiffSuppressForDNS(t *testing.T) {
	cases := map[string]struct {
		Old, New       string
		ShouldSuppress bool
	}{
		"compact form should suppress diff": {
			Old:            "2a03:b0c0:1:e0::29b:8001",
			New:            "2a03:b0c0:0001:00e0:0000:0000:029b:8001",
			ShouldSuppress: true,
		},
		"different address should not suppress diff": {
			Old:            "2a03:b0c0:1:e00::29b:8001",
			New:            "2a03:b0c0:0001:00e0:0000:0000:029b:8001",
			ShouldSuppress: false,
		},
	}

	for tn, tc := range cases {
		shouldSuppress := AddressDiffSuppress("", tc.Old, tc.New, nil)
		if shouldSuppress != tc.ShouldSuppress {
			t.Errorf("%s: expected %t", tn, tc.ShouldSuppress)
		}
	}
}

func TestAccDNSResourceRecordSet_basic(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s.hashicorptest.com./A", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check both import formats
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSResourceRecordSet_Update(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				Config: testAccDnsResourceRecordSet_basic(zoneName, "127.0.0.11", 300),
			},
			{
				Config: testAccDnsResourceRecordSet_basic(zoneName, "127.0.0.11", 600),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/%s.hashicorptest.com./A", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSResourceRecordSet_changeType(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_basic(zoneName, "127.0.0.10", 300),
			},
			{
				Config: testAccDnsResourceRecordSet_bigChange(zoneName, 600),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s/sub-domain.%s.hashicorptest.com./CNAME", getTestProjectFromEnv(), zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSResourceRecordSet_nestedNS(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-ns-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_nestedNS(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/nested.%s.hashicorptest.com./NS", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSResourceRecordSet_quotedTXT(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_quotedTXT(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s.hashicorptest.com./TXT", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDNSResourceRecordSet_MX(t *testing.T) {
	t.Parallel()

	zoneName := fmt.Sprintf("dnszone-test-txt-%s", randString(t, 10))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsResourceRecordSetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDnsResourceRecordSet_MX(zoneName, 300),
			},
			{
				ResourceName:      "google_dns_resource_dns_record_set.foobar",
				ImportStateId:     fmt.Sprintf("%s/%s.hashicorptest.com./MX", zoneName, zoneName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDnsResourceRecordSetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_dns_resource_dns_record_set" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{DNSBasePath}}projects/{{project}}/managedZones/{{managed_zone}}/rrsets/{{name}}/{{type}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = sendRequest(config, "GET", billingProject, url, config.userAgent, nil)
			if err == nil {
				return fmt.Errorf("DNSResourceDnsRecordSet still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccDnsResourceRecordSet_basic(zoneName string, addr2 string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_resource_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "%s.hashicorptest.com."
  type         = "A"
  rrdatas      = ["127.0.0.1", "%s"]
  ttl          = %d
}
`, zoneName, zoneName, zoneName, addr2, ttl)
}

func testAccDnsResourceRecordSet_nestedNS(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_resource_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "nested.%s.hashicorptest.com."
  type         = "NS"
  rrdatas      = ["ns.hashicorp.services.", "ns2.hashicorp.services."]
  ttl          = %d
}
`, name, name, name, ttl)
}

func testAccDnsResourceRecordSet_bigChange(zoneName string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_resource_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "sub-domain.%s.hashicorptest.com."
  type         = "CNAME"
  rrdatas      = ["www.terraform.io."]
  ttl          = %d
}
`, zoneName, zoneName, zoneName, ttl)
}

func testAccDnsResourceRecordSet_quotedTXT(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_resource_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "%s.hashicorptest.com."
  type         = "TXT"
  rrdatas      = ["\"quoted text\"", "\"test space\""]
  ttl          = %d
}
`, name, name, name, ttl)
}

func testAccDnsResourceRecordSet_MX(name string, ttl int) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "parent-zone" {
  name        = "%s"
  dns_name    = "%s.hashicorptest.com."
  description = "Test Description"
}

resource "google_dns_resource_dns_record_set" "foobar" {
  managed_zone = google_dns_managed_zone.parent-zone.name
  name         = "%s.hashicorptest.com."
  type         = "MX"
  rrdatas = [
	"1 aspmx.l.google.com.",
    "5 alt1.aspmx.l.google.com.",
    "5 alt2.aspmx.l.google.com.",
    "10 aspmx2.googlemail.com.",
    "10 aspmx3.googlemail.com.",
  ]
  ttl = %d
}
`, name, name, name, ttl)
}
