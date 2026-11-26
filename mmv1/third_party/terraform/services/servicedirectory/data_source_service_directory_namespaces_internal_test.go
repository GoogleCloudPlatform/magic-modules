package servicedirectory

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

// TestDataSourceServiceDirectoryNamespacesRead checks paging, the filter parameter and the field mapping against a fake API.
func TestDataSourceServiceDirectoryNamespacesRead(t *testing.T) {
	t.Parallel()

	const page1JSON = `{
		"nextPageToken": "page-2",
		"namespaces": [{
			"name": "projects/test-project/locations/us-central1/namespaces/namespace-1",
			"labels": {"env": "prod", "team": "net"},
			"uid": "0e4a2c5e-7d0b-4c59-9f0e-1a2b3c4d5e6f"
		}]
	}`
	// The second namespace has no labels, as the API omits empty maps.
	const page2JSON = `{
		"namespaces": [{"name": "projects/test-project/locations/us-central1/namespaces/namespace-2"}]
	}`

	var mu sync.Mutex
	var queries []url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got, want := r.Header.Get("X-Goog-User-Project"), "billing-project"; got != want {
			t.Errorf("quota project header = %q, want %q", got, want)
		}
		switch r.URL.Path {
		case "/v1/projects/test-project/locations/us-central1/namespaces":
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
		case "/v1/projects/bad-project/locations/us-central1/namespaces":
			_, _ = w.Write([]byte(`{"namespaces": ["not-an-object"]}`))
		case "/v1/projects/empty-project/locations/us-central1/namespaces":
			_, _ = w.Write([]byte(`{}`))
		// This is projects.locations.get, which a location ending in an unescaped "#" would reach.
		case "/v1/projects/test-project/locations/us-central1":
			_, _ = w.Write([]byte(`{"name": "projects/test-project/locations/us-central1", "locationId": "us-central1"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceServiceDirectoryNamespaces()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v1/"},
	}

	const filter = "labels.env=prod"
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"location": "us-central1", "filter": filter})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read with a filter returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/locations/us-central1/namespaces"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	wantFirst := map[string]interface{}{
		"name":         "projects/test-project/locations/us-central1/namespaces/namespace-1",
		"namespace_id": "namespace-1",
		"labels":       map[string]interface{}{"env": "prod", "team": "net"},
	}
	namespaces := d.Get("namespaces").([]interface{})
	if len(namespaces) != 2 {
		t.Fatalf("got %d namespaces, want 2 (one from each page): %#v", len(namespaces), namespaces)
	}
	if !reflect.DeepEqual(namespaces[0], wantFirst) {
		t.Errorf("first namespace = %#v, want %#v", namespaces[0], wantFirst)
	}
	second := namespaces[1].(map[string]interface{})
	if got, want := second["namespace_id"], "namespace-2"; got != want {
		t.Errorf("second namespace_id = %q, want %q", got, want)
	}
	if got := second["labels"].(map[string]interface{}); len(got) != 0 {
		t.Errorf("second namespace labels = %#v, want none", got)
	}

	mu.Lock()
	filtered := queries
	queries = nil
	mu.Unlock()
	if len(filtered) != 2 || filtered[0].Get("filter") != filter || filtered[1].Get("filter") != filter || filtered[1].Get("pageToken") != "page-2" {
		t.Errorf("filtered read sent queries %v, want the filter on both pages and the page token on the second", filtered)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"location": "us-central1"})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read without a filter returned an error: %s", err)
	}
	mu.Lock()
	unfiltered := queries
	mu.Unlock()
	if len(unfiltered) != 2 || unfiltered[0].Has("filter") || unfiltered[1].Has("filter") {
		t.Errorf("unfiltered read sent queries %v, want no filter parameter", unfiltered)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"location": "us-central1", "project": "bad-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read accepted a list item that is not an object, want an error")
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"location": "us-central1", "project": "empty-project"})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read of a location with no namespaces returned an error: %s", err)
	}
	if got := d.Get("namespaces").([]interface{}); len(got) != 0 {
		t.Errorf("empty location returned %#v, want no namespaces", got)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"location": "us-central1#"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read of a location ending in # returned no error, want the request to fail")
	}
}
