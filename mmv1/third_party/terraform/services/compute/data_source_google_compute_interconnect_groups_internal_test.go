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

// TestDataSourceGoogleComputeInterconnectGroupsRead checks paging and the fields of each group against a fake API.
func TestDataSourceGoogleComputeInterconnectGroupsRead(t *testing.T) {
	t.Parallel()

	// The first group sets every field, with one member link in the beta form, and a second group with only an intent shares its page.
	const page1JSON = `{
		"kind": "compute#InterconnectGroupsListResponse",
		"nextPageToken": "page-2",
		"items": [{
			"kind": "compute#InterconnectGroup",
			"id": "1111111111111111111",
			"name": "group-1",
			"description": "first group",
			"creationTimestamp": "2026-11-10T00:00:00.000-08:00",
			"selfLink": "https://www.googleapis.com/compute/v1/projects/test-project/global/interconnectGroups/group-1",
			"etag": "etag-1",
			"interconnects": {
				"interconnect-a": {"interconnect": "https://www.googleapis.com/compute/v1/projects/test-project/global/interconnects/interconnect-a"},
				"interconnect-b": {"interconnect": "https://www.googleapis.com/compute/beta/projects/test-project/global/interconnects/interconnect-b"}
			},
			"intent": {"topologyCapability": "PRODUCTION_CRITICAL"},
			"physicalStructure": {
				"metros": [{
					"metro": "iad",
					"facilities": [{
						"facility": "5467",
						"zones": [{
							"zone": "zone1",
							"interconnects": ["https://www.googleapis.com/compute/v1/projects/test-project/global/interconnects/interconnect-a"]
						}]
					}]
				}]
			},
			"configured": {
				"topologyCapability": {
					"supportedSla": "NO_SLA",
					"intendedCapabilityBlockers": [{
						"blockerType": "INCOMPATIBLE_METROS",
						"explanation": "the metros of this group cannot give the intended capability",
						"documentationLink": "https://cloud.google.com/network-connectivity/docs/interconnect",
						"metros": ["iad"],
						"facilities": ["5467"],
						"zones": ["zone1"],
						"interconnects": ["https://www.googleapis.com/compute/v1/projects/test-project/global/interconnects/interconnect-a"]
					}]
				}
			}
		}, {
			"name": "group-2",
			"intent": {"topologyCapability": "NO_SLA"}
		}]
	}`
	const page2JSON = `{
		"kind": "compute#InterconnectGroupsListResponse",
		"items": [{
			"name": "group-3",
			"description": "third group",
			"intent": {"topologyCapability": "PRODUCTION_NON_CRITICAL"}
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/compute/v1/projects/test-project/global/interconnectGroups":
			if got, want := r.Header.Get("X-Goog-User-Project"), "billing-project"; got != want {
				t.Errorf("quota project header = %q, want the billing project %q", got, want)
			}
			if _, ok := r.URL.Query()["returnPartialSuccess"]; ok {
				t.Errorf("the list request asks for partial results, want it left out so the API fails instead")
			}
			pageToken := r.URL.Query().Get("pageToken")
			mu.Lock()
			pageTokens = append(pageTokens, pageToken)
			requests := len(pageTokens)
			mu.Unlock()
			switch {
			// A second request for the first page means paging is broken, and the read would never end.
			case pageToken == "" && requests > 1:
				http.Error(w, "first page requested twice", http.StatusBadRequest)
			case pageToken == "":
				_, _ = w.Write([]byte(page1JSON))
			case pageToken == "page-2":
				_, _ = w.Write([]byte(page2JSON))
			default:
				http.Error(w, "unexpected page token", http.StatusBadRequest)
			}
		case "/compute/v1/projects/empty-project/global/interconnectGroups":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{"kind": "compute#InterconnectGroupsListResponse"}`))
		case "/compute/v1/projects/bad-project/global/interconnectGroups":
			_, _ = w.Write([]byte(`{"items": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleComputeInterconnectGroups()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/compute/v1/"},
	}

	// These fields of the resource are not part of a group: project repeats the argument and deletion_policy is a setting of the resource.
	groupSchema := ds.Schema["interconnect_groups"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"project", "deletion_policy"} {
		if _, ok := groupSchema[field]; ok {
			t.Errorf("the schema of a group has the field %q, want it left out", field)
		}
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/global/interconnectGroups"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	const linkA = "https://www.googleapis.com/compute/v1/projects/test-project/global/interconnects/interconnect-a"
	const blockers = "interconnect_groups.0.configured.0.topology_capability.0.intended_capability_blockers.0."
	// Only the attributes the fake API sets are compared, in the flat form of the state, so a field added to InterconnectGroup.yaml needs no change here.
	want := map[string]string{
		"interconnect_groups.#":                                                                    "3",
		"interconnect_groups.0.id":                                                                 "projects/test-project/global/interconnectGroups/group-1",
		"interconnect_groups.0.name":                                                               "group-1",
		"interconnect_groups.0.description":                                                        "first group",
		"interconnect_groups.0.creation_timestamp":                                                 "2026-11-10T00:00:00.000-08:00",
		"interconnect_groups.0.interconnects.#":                                                    "2",
		"interconnect_groups.0.intent.0.topology_capability":                                       "PRODUCTION_CRITICAL",
		"interconnect_groups.0.physical_structure.0.metros.0.metro":                                "iad",
		"interconnect_groups.0.physical_structure.0.metros.0.facilities.0.facility":                "5467",
		"interconnect_groups.0.physical_structure.0.metros.0.facilities.0.zones.0.zone":            "zone1",
		"interconnect_groups.0.physical_structure.0.metros.0.facilities.0.zones.0.interconnects.0": linkA,
		"interconnect_groups.0.configured.0.topology_capability.0.supported_sla":                   "NO_SLA",
		blockers + "blocker_type":                                                                  "INCOMPATIBLE_METROS",
		blockers + "explanation":                                                                   "the metros of this group cannot give the intended capability",
		blockers + "documentation_link":                                                            "https://cloud.google.com/network-connectivity/docs/interconnect",
		blockers + "metros.0":                                                                      "iad",
		blockers + "facilities.0":                                                                  "5467",
		blockers + "zones.0":                                                                       "zone1",
		blockers + "interconnects.0":                                                               linkA,

		"interconnect_groups.1.id":                           "projects/test-project/global/interconnectGroups/group-2",
		"interconnect_groups.1.name":                         "group-2",
		"interconnect_groups.1.description":                  "",
		"interconnect_groups.1.intent.0.topology_capability": "NO_SLA",
		"interconnect_groups.1.interconnects.#":              "0",
		"interconnect_groups.1.physical_structure.#":         "0",
		"interconnect_groups.1.configured.#":                 "0",

		"interconnect_groups.2.id":                           "projects/test-project/global/interconnectGroups/group-3",
		"interconnect_groups.2.name":                         "group-3",
		"interconnect_groups.2.description":                  "third group",
		"interconnect_groups.2.intent.0.topology_capability": "PRODUCTION_NON_CRITICAL",
	}
	attributes := d.State().Attributes
	for key, wantValue := range want {
		if got, ok := attributes[key]; !ok || got != wantValue {
			t.Errorf("%s = %q (present: %t), want %q", key, got, ok, wantValue)
		}
	}

	// The state keys of a set hold a hash of each member, so the member interconnects are read by name.
	members := map[string]interface{}{}
	first := d.Get("interconnect_groups").([]interface{})[0].(map[string]interface{})
	for _, raw := range first["interconnects"].(*schema.Set).List() {
		member := raw.(map[string]interface{})
		members[member["name"].(string)] = member["interconnect"]
	}
	// A link comes out as the API sent it, as in the google_compute_interconnect_group resource, so the beta form is kept.
	wantMembers := map[string]interface{}{
		"interconnect-a": linkA,
		"interconnect-b": "https://www.googleapis.com/compute/beta/projects/test-project/global/interconnects/interconnect-b",
	}
	if !reflect.DeepEqual(members, wantMembers) {
		t.Errorf("first group member interconnects = %#v, want %#v", members, wantMembers)
	}

	mu.Lock()
	gotTokens := pageTokens
	mu.Unlock()
	if !reflect.DeepEqual(gotTokens, []string{"", "page-2"}) {
		t.Errorf("list requests sent page tokens %q, want the first page and then page-2", gotTokens)
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
		t.Fatalf("Read of a project with no groups returned an error: %s", err)
	}
	if got := d.Get("interconnect_groups").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no groups", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}
	if got, want := d.Id(), "projects/empty-project/global/interconnectGroups"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "bad-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read accepted a list item that is not an object, want an error")
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "missing-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read of a project the API does not know returned no error")
	}

	cfgNoProject := &transport_tpg.Config{
		Context:         context.Background(),
		Client:          server.Client(),
		CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/compute/v1/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfgNoProject); err == nil {
		t.Errorf("Read with no project argument and no provider project returned no error")
	}
}
