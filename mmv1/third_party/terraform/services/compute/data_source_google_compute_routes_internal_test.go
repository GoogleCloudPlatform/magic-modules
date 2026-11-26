package compute

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestDataSourceGoogleComputeRoutesRead checks paging, the filter parameter and the field mapping against a fake API.
func TestDataSourceGoogleComputeRoutesRead(t *testing.T) {
	t.Parallel()

	// The first route sets every field at once, which no real route does, so each mapping is checked.
	const page1JSON = `{
		"kind": "compute#routeList",
		"nextPageToken": "page-2",
		"items": [{
			"kind": "compute#route",
			"id": "1001",
			"name": "route-1",
			"description": "first route",
			"network": "https://www.googleapis.com/compute/v1/projects/test-project/global/networks/vpc-1",
			"destRange": "10.10.0.0/16",
			"priority": 900,
			"tags": ["web", "db"],
			"nextHopGateway": "https://www.googleapis.com/compute/v1/projects/test-project/global/gateways/default-internet-gateway",
			"nextHopInstance": "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/instances/vm-1",
			"nextHopIp": "10.0.0.5",
			"nextHopVpnTunnel": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/vpnTunnels/tunnel-1",
			"nextHopIlb": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/forwardingRules/ilb-1",
			"nextHopNetwork": "https://www.googleapis.com/compute/v1/projects/test-project/global/networks/vpc-2",
			"nextHopPeering": "peering-1",
			"nextHopHub": "projects/test-project/locations/global/hubs/hub-1",
			"nextHopInterconnectAttachment": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/interconnectAttachments/attachment-1",
			"nextHopOrigin": "IGP",
			"nextHopMed": 50,
			"nextHopInterRegionCost": 200,
			"routeType": "STATIC",
			"routeStatus": "ACTIVE",
			"selfLink": "https://www.googleapis.com/compute/v1/projects/test-project/global/routes/route-1",
			"creationTimestamp": "2026-11-10T00:00:00.000-08:00"
		}]
	}`
	// The second route sends one metric as a string and the other as the largest 32-bit number.
	const page2JSON = `{
		"kind": "compute#routeList",
		"items": [{"name": "route-2", "destRange": "10.20.0.0/16", "nextHopMed": "70", "nextHopInterRegionCost": 4294967295}]
	}`

	var mu sync.Mutex
	var queries []url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got, want := r.Header.Get("X-Goog-User-Project"), "billing-project"; got != want {
			t.Errorf("quota project header = %q, want %q", got, want)
		}
		switch r.URL.Path {
		case "/compute/v1/projects/test-project/global/routes":
			mu.Lock()
			queries = append(queries, r.URL.Query())
			mu.Unlock()
			switch r.URL.Query().Get("pageToken") {
			case "":
				_, _ = w.Write([]byte(page1JSON))
			case "page-2":
				_, _ = w.Write([]byte(page2JSON))
			default:
				http.Error(w, "unexpected page token", http.StatusBadRequest)
			}
		case "/compute/v1/projects/bad-project/global/routes":
			_, _ = w.Write([]byte(`{"items": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleComputeRoutes()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{"compute_custom_endpoint": server.URL + "/compute/v1/"},
	}

	const filter = "network eq .*vpc-1"
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"filter": filter})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read with a filter returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/global/routes"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	wantFirst := map[string]interface{}{
		"name":                             "route-1",
		"description":                      "first route",
		"network":                          "https://www.googleapis.com/compute/v1/projects/test-project/global/networks/vpc-1",
		"dest_range":                       "10.10.0.0/16",
		"priority":                         900,
		"tags":                             []interface{}{"web", "db"},
		"next_hop_gateway":                 "https://www.googleapis.com/compute/v1/projects/test-project/global/gateways/default-internet-gateway",
		"next_hop_instance":                "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/instances/vm-1",
		"next_hop_ip":                      "10.0.0.5",
		"next_hop_vpn_tunnel":              "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/vpnTunnels/tunnel-1",
		"next_hop_ilb":                     "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/forwardingRules/ilb-1",
		"next_hop_network":                 "https://www.googleapis.com/compute/v1/projects/test-project/global/networks/vpc-2",
		"next_hop_peering":                 "peering-1",
		"next_hop_hub":                     "projects/test-project/locations/global/hubs/hub-1",
		"next_hop_interconnect_attachment": "https://www.googleapis.com/compute/v1/projects/test-project/regions/us-central1/interconnectAttachments/attachment-1",
		"next_hop_origin":                  "IGP",
		"next_hop_med":                     "50",
		"next_hop_inter_region_cost":       "200",
		"route_type":                       "STATIC",
		"route_status":                     "ACTIVE",
		"self_link":                        "https://www.googleapis.com/compute/v1/projects/test-project/global/routes/route-1",
		"creation_timestamp":               "2026-11-10T00:00:00.000-08:00",
	}
	routes := d.Get("routes").([]interface{})
	if len(routes) != 2 {
		t.Fatalf("got %d routes, want 2 (one from each page): %#v", len(routes), routes)
	}
	if !reflect.DeepEqual(routes[0], wantFirst) {
		t.Errorf("first route = %#v, want %#v", routes[0], wantFirst)
	}
	second := routes[1].(map[string]interface{})
	for field, want := range map[string]string{"name": "route-2", "next_hop_med": "70", "next_hop_inter_region_cost": "4294967295"} {
		if got := second[field]; got != want {
			t.Errorf("second route %s = %q, want %q", field, got, want)
		}
	}

	mu.Lock()
	filtered := queries
	queries = nil
	mu.Unlock()
	if len(filtered) != 2 || filtered[0].Get("filter") != filter || filtered[1].Get("filter") != filter || filtered[1].Get("pageToken") != "page-2" {
		t.Errorf("filtered read sent queries %v, want the filter on both pages and the page token on the second", filtered)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read without a filter returned an error: %s", err)
	}
	mu.Lock()
	unfiltered := queries
	mu.Unlock()
	if len(unfiltered) != 2 || unfiltered[0].Has("filter") || unfiltered[1].Has("filter") {
		t.Errorf("unfiltered read sent queries %v, want no filter parameter", unfiltered)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "bad-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read accepted a list item that is not an object, want an error")
	}
}
