package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestDataSourceGoogleComputeInterconnectAttachmentsRead checks both list calls, paging and the field mapping against a fake API.
func TestDataSourceGoogleComputeInterconnectAttachmentsRead(t *testing.T) {
	t.Parallel()

	// The first attachment sets every field and uses beta links, so each mapping and each link conversion is checked.
	const regionPage1JSON = `{
		"nextPageToken": "page-2",
		"items": [{
			"name": "attachment-1",
			"description": "first attachment",
			"selfLink": "https://www.googleapis.com/compute/beta/projects/test-project/regions/europe-west1/interconnectAttachments/attachment-1",
			"region": "https://www.googleapis.com/compute/beta/projects/test-project/regions/europe-west1",
			"router": "https://www.googleapis.com/compute/beta/projects/test-project/regions/europe-west1/routers/router-1",
			"interconnect": "https://www.googleapis.com/compute/beta/projects/test-project/global/interconnects/interconnect-1",
			"attachmentGroup": "https://www.googleapis.com/compute/beta/projects/test-project/global/interconnectAttachmentGroups/group-1",
			"type": "DEDICATED",
			"state": "ACTIVE",
			"adminEnabled": true,
			"bandwidth": "BPS_10G",
			"edgeAvailabilityDomain": "AVAILABILITY_DOMAIN_1",
			"mtu": 1500,
			"vlanTag8021q": 1234,
			"stackType": "IPV4_IPV6",
			"encryption": "IPSEC",
			"ipsecInternalAddresses": ["https://www.googleapis.com/compute/beta/projects/test-project/regions/europe-west1/addresses/range-1"],
			"cloudRouterIpAddress": "169.254.10.1/29",
			"customerRouterIpAddress": "169.254.10.2/29",
			"cloudRouterIpv6Address": "2600:2d00:0:1::1/125",
			"customerRouterIpv6Address": "2600:2d00:0:1::2/125",
			"candidateCloudRouterIpAddress": "169.254.20.1/29",
			"candidateCustomerRouterIpAddress": "169.254.20.2/29",
			"candidateCloudRouterIpv6Address": "2600:2d00:0:2::1/125",
			"candidateCustomerRouterIpv6Address": "2600:2d00:0:2::2/125",
			"pairingKey": "pairing-key-1",
			"partnerAsn": "65001",
			"privateInterconnectInfo": {"tag8021q": 1234},
			"googleReferenceId": "reference-1",
			"creationTimestamp": "2026-11-10T00:00:00.000-08:00",
			"labels": {"env": "test", "team": "network"},
			"labelFingerprint": "fingerprint-1",
			"l2Forwarding": {
				"network": "https://www.googleapis.com/compute/beta/projects/test-project/global/networks/network-1",
				"tunnelEndpointIpAddress": "10.0.0.1",
				"defaultApplianceIpAddress": "10.0.0.2",
				"geneveHeader": {"vni": 42},
				"applianceMappings": {
					"1000": {"name": "appliance-1000", "applianceIpAddress": "10.0.0.4"},
					"30": {"name": "appliance-30"},
					"4-5": {"name": "appliance-4-5"},
					"200": {
						"name": "appliance-200",
						"applianceIpAddress": "10.0.0.3",
						"innerVlanToApplianceMappings": [{"innerApplianceIpAddress": "10.0.0.5", "innerVlanTags": ["10", "20-30"]}]
					}
				}
			}
		}, {
			"name": "attachment-2",
			"region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1",
			"type": "PARTNER"
		}]
	}`
	const regionPage2JSON = `{
		"items": [{
			"name": "attachment-3",
			"region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1",
			"type": "PARTNER"
		}]
	}`
	// A region with no attachments carries a warning in place of the list, and the regions are split over two pages.
	const aggregatedPage1JSON = `{
		"nextPageToken": "page-2",
		"items": {
			"regions/us-central1": {"interconnectAttachments": [
				{"name": "us-b", "region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1"},
				{"name": "us-a", "region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1"}
			]},
			"regions/asia-east1": {"warning": {"code": "NO_RESULTS_ON_PAGE", "message": "No results for the scope 'regions/asia-east1'."}},
			"regions/europe-west1": {"interconnectAttachments": [
				{"name": "eu-a", "region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1"}
			]}
		}
	}`
	const aggregatedPage2JSON = `{
		"items": {
			"regions/asia-south1": {"interconnectAttachments": [
				{"name": "asia-a", "region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/asia-south1"}
			]}
		}
	}`

	var mu sync.Mutex
	pageTokens := map[string][]string{}
	// servePages answers the first page, then the second, and fails a repeated first page, which would mean an endless read.
	servePages := func(w http.ResponseWriter, r *http.Request, page1, page2 string) {
		pageToken := r.URL.Query().Get("pageToken")
		mu.Lock()
		pageTokens[r.URL.Path] = append(pageTokens[r.URL.Path], pageToken)
		requests := len(pageTokens[r.URL.Path])
		mu.Unlock()
		switch {
		case pageToken == "" && requests > 1:
			http.Error(w, "first page requested twice", http.StatusBadRequest)
		case pageToken == "":
			_, _ = w.Write([]byte(page1))
		case pageToken == "page-2":
			_, _ = w.Write([]byte(page2))
		default:
			http.Error(w, "unexpected page token", http.StatusBadRequest)
		}
	}
	const regionPath = "/compute/v1/projects/test-project/regions/europe-west1/interconnectAttachments"
	const aggregatedPath = "/compute/v1/projects/test-project/aggregated/interconnectAttachments"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Has("returnPartialSuccess") {
			t.Errorf("request asked for partial results, want a failing region to fail the read: %s", r.URL)
		}
		switch r.URL.Path {
		case regionPath:
			if got, want := r.Header.Get("X-Goog-User-Project"), "billing-project"; got != want {
				t.Errorf("quota project header = %q, want the billing project %q", got, want)
			}
			servePages(w, r, regionPage1JSON, regionPage2JSON)
		case aggregatedPath:
			servePages(w, r, aggregatedPage1JSON, aggregatedPage2JSON)
		case "/compute/v1/projects/empty-project/aggregated/interconnectAttachments":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{"items": {"regions/europe-west1": {"warning": {"code": "NO_RESULTS_ON_PAGE"}}}}`))
		case "/compute/v1/projects/bad-item/regions/europe-west1/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": ["not-an-object"]}`))
		case "/compute/v1/projects/bad-scope/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": {"regions/europe-west1": "not-an-object"}}`))
		case "/compute/v1/projects/empty-project/regions/europe-west1/interconnectAttachments":
			_, _ = w.Write([]byte(`{"kind": "compute#interconnectAttachmentList"}`))
		case "/compute/v1/projects/bad-items/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": "not-a-list"}`))
		case "/compute/v1/projects/bad-scope-list/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": {"regions/europe-west1": {"interconnectAttachments": "not-a-list"}}}`))
		case "/compute/v1/projects/unreachable-scope/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": {
				"regions/europe-west1": {"interconnectAttachments": [{"name": "eu-a"}]},
				"regions/us-central1": {"warning": {"code": "UNREACHABLE", "message": "Scope regions/us-central1 is unreachable."}}
			}}`))
		case "/compute/v1/projects/unreachable-list/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": {"regions/europe-west1": {"interconnectAttachments": [{"name": "eu-a"}]}}, "unreachables": ["regions/us-central1"]}`))
		case "/compute/v1/projects/bad-scope-item/aggregated/interconnectAttachments":
			_, _ = w.Write([]byte(`{"items": {"regions/europe-west1": {"interconnectAttachments": ["not-an-object"]}}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleComputeInterconnectAttachments()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		Region:              "us-east1",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/compute/v1/"},
	}
	names := func(d *schema.ResourceData) []string {
		got := []string{}
		for _, raw := range d.Get("interconnect_attachments").([]interface{}) {
			got = append(got, raw.(map[string]interface{})["name"].(string))
		}
		return got
	}

	// The region is given as a self link, the form google_compute_interconnect_attachment keeps in its region attribute.
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"region": "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1",
	})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read of one region returned an error: %s", err)
	}
	if got, want := d.Id(), "projects/test-project/regions/europe-west1/interconnectAttachments"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}
	if got, want := names(d), []string{"attachment-1", "attachment-2", "attachment-3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("attachments of one region = %q, want %q (two from the first page, one from the second)", got, want)
	}

	first := d.Get("interconnect_attachments").([]interface{})[0].(map[string]interface{})
	// This comparison is exact on purpose: a field added to InterconnectAttachment.yaml must be added here and in flattenComputeInterconnectAttachmentsItem.
	wantFirst := map[string]interface{}{
		"name":                                   "attachment-1",
		"description":                            "first attachment",
		"self_link":                              "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1/interconnectAttachments/attachment-1",
		"region":                                 "europe-west1",
		"router":                                 "https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1/routers/router-1",
		"interconnect":                           "https://www.googleapis.com/compute/v1/projects/test-project/global/interconnects/interconnect-1",
		"attachment_group":                       "https://www.googleapis.com/compute/v1/projects/test-project/global/interconnectAttachmentGroups/group-1",
		"type":                                   "DEDICATED",
		"state":                                  "ACTIVE",
		"admin_enabled":                          true,
		"bandwidth":                              "BPS_10G",
		"edge_availability_domain":               "AVAILABILITY_DOMAIN_1",
		"mtu":                                    "1500",
		"vlan_tag8021q":                          1234,
		"stack_type":                             "IPV4_IPV6",
		"encryption":                             "IPSEC",
		"ipsec_internal_addresses":               []interface{}{"https://www.googleapis.com/compute/v1/projects/test-project/regions/europe-west1/addresses/range-1"},
		"cloud_router_ip_address":                "169.254.10.1/29",
		"customer_router_ip_address":             "169.254.10.2/29",
		"cloud_router_ipv6_address":              "2600:2d00:0:1::1/125",
		"customer_router_ipv6_address":           "2600:2d00:0:1::2/125",
		"candidate_cloud_router_ip_address":      "169.254.20.1/29",
		"candidate_customer_router_ip_address":   "169.254.20.2/29",
		"candidate_cloud_router_ipv6_address":    "2600:2d00:0:2::1/125",
		"candidate_customer_router_ipv6_address": "2600:2d00:0:2::2/125",
		"pairing_key":                            "pairing-key-1",
		"partner_asn":                            "65001",
		"private_interconnect_info":              []interface{}{map[string]interface{}{"tag8021q": 1234}},
		"google_reference_id":                    "reference-1",
		"creation_timestamp":                     "2026-11-10T00:00:00.000-08:00",
		"labels":                                 map[string]interface{}{"env": "test", "team": "network"},
		"label_fingerprint":                      "fingerprint-1",
		"l2_forwarding": []interface{}{map[string]interface{}{
			"network":                      "https://www.googleapis.com/compute/v1/projects/test-project/global/networks/network-1",
			"tunnel_endpoint_ip_address":   "10.0.0.1",
			"default_appliance_ip_address": "10.0.0.2",
			"geneve_header":                []interface{}{map[string]interface{}{"vni": 42}},
			// The keys are ordered by their first VLAN tag as a number: a key can be a range, and as text 1000 would come first.
			"appliance_mappings": []interface{}{
				map[string]interface{}{
					"vlan_id":                          "4-5",
					"name":                             "appliance-4-5",
					"appliance_ip_address":             "",
					"inner_vlan_to_appliance_mappings": []interface{}{},
				},
				map[string]interface{}{
					"vlan_id":                          "30",
					"name":                             "appliance-30",
					"appliance_ip_address":             "",
					"inner_vlan_to_appliance_mappings": []interface{}{},
				},
				map[string]interface{}{
					"vlan_id":              "200",
					"name":                 "appliance-200",
					"appliance_ip_address": "10.0.0.3",
					"inner_vlan_to_appliance_mappings": []interface{}{map[string]interface{}{
						"inner_appliance_ip_address": "10.0.0.5",
						"inner_vlan_tags":            []interface{}{"10", "20-30"},
					}},
				},
				map[string]interface{}{
					"vlan_id":                          "1000",
					"name":                             "appliance-1000",
					"appliance_ip_address":             "10.0.0.4",
					"inner_vlan_to_appliance_mappings": []interface{}{},
				},
			},
		}},
	}
	if !reflect.DeepEqual(first, wantFirst) {
		t.Errorf("first attachment = %#v, want %#v", first, wantFirst)
	}
	// This names the field when the exact comparison above fails because the schema gained one.
	mapped := flattenComputeInterconnectAttachmentsItem(map[string]interface{}{}, d, cfg)
	for field := range ds.Schema["interconnect_attachments"].Elem.(*schema.Resource).Schema {
		if _, ok := mapped[field]; !ok {
			t.Errorf("flattenComputeInterconnectAttachmentsItem does not set %q: map it there or remove it from the schema", field)
		}
	}

	second := d.Get("interconnect_attachments").([]interface{})[1].(map[string]interface{})
	if got, want := second["encryption"], "NONE"; got != want {
		t.Errorf("encryption of an attachment without one = %q, want the default %q", got, want)
	}
	if got := second["l2_forwarding"].([]interface{}); len(got) != 0 {
		t.Errorf("l2_forwarding of an attachment without one = %#v, want none", got)
	}

	// With no region argument every region is listed, and the provider region us-east1 is not used.
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read of all regions returned an error: %s", err)
	}
	if got, want := d.Id(), "projects/test-project/aggregated/interconnectAttachments"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	// The result is sorted by region across pages, and attachments keep the order the API gave inside a region.
	if got, want := names(d), []string{"asia-a", "eu-a", "us-b", "us-a"}; !reflect.DeepEqual(got, want) {
		t.Errorf("attachments of all regions = %q, want %q", got, want)
	}
	if got, want := d.Get("interconnect_attachments").([]interface{})[0].(map[string]interface{})["region"], "asia-south1"; got != want {
		t.Errorf("region of the first attachment = %q, want %q", got, want)
	}

	mu.Lock()
	gotTokens := map[string][]string{regionPath: pageTokens[regionPath], aggregatedPath: pageTokens[aggregatedPath]}
	mu.Unlock()
	for path, tokens := range gotTokens {
		if !reflect.DeepEqual(tokens, []string{"", "page-2"}) {
			t.Errorf("list requests to %s sent page tokens %q, want the first page and then page-2", path, tokens)
		}
	}

	// Without a billing project, the project being read is the quota project.
	cfgNoBilling := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/compute/v1/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project"})
	if err := ds.Read(d, cfgNoBilling); err != nil {
		t.Fatalf("Read of a project with no attachments returned an error: %s", err)
	}
	if got := d.Get("interconnect_attachments").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no attachments", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}

	// The list of one region has no items key at all when the region is empty.
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project", "region": "europe-west1"})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read of a region with no attachments returned an error: %s", err)
	}
	if got := d.Get("interconnect_attachments").([]interface{}); len(got) != 0 {
		t.Errorf("empty region returned %#v, want no attachments", got)
	}

	for _, tc := range []struct {
		name   string
		config map[string]interface{}
	}{
		{"a list item that is not an object", map[string]interface{}{"project": "bad-item", "region": "europe-west1"}},
		{"a region entry that is not an object", map[string]interface{}{"project": "bad-scope"}},
		{"an item of a region that is not an object", map[string]interface{}{"project": "bad-scope-item"}},
		{"items that are neither a list nor an object", map[string]interface{}{"project": "bad-items"}},
		{"attachments of a region that are not a list", map[string]interface{}{"project": "bad-scope-list"}},
		{"a region the API marks as unreachable", map[string]interface{}{"project": "unreachable-scope"}},
		{"a list of unreachable regions", map[string]interface{}{"project": "unreachable-list"}},
		{"a project the API does not know", map[string]interface{}{"project": "missing-project"}},
	} {
		d = schema.TestResourceDataRaw(t, ds.Schema, tc.config)
		if err := ds.Read(d, cfg); err == nil {
			t.Errorf("Read accepted %s, want an error", tc.name)
		}
	}
}
