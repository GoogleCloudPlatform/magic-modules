package logging

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

// TestDataSourceGoogleLoggingProjectSinksRead checks paging and the field mapping against a fake API.
func TestDataSourceGoogleLoggingProjectSinksRead(t *testing.T) {
	t.Parallel()

	// The first sink sets every field, so each mapping is checked, and a second sink shares its page.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"sinks": [{
			"name": "sink-1",
			"resourceName": "projects/test-project/sinks/sink-1",
			"destination": "bigquery.googleapis.com/projects/test-project/datasets/logs",
			"filter": "severity>=ERROR",
			"description": "first sink",
			"disabled": true,
			"writerIdentity": "serviceAccount:service-123@gcp-sa-logging.iam.gserviceaccount.com",
			"exclusions": [{"name": "exclusion-1", "description": "first exclusion", "filter": "resource.type=\"gce_instance\"", "disabled": true}],
			"bigqueryOptions": {"usePartitionedTables": true, "usesTimestampColumnPartitioning": true}
		}, {
			"name": "sink-2",
			"destination": "storage.googleapis.com/bucket-2"
		}]
	}`
	// The last sink has only the fields the API returns for the built-in sinks.
	const page2JSON = `{
		"sinks": [{
			"name": "_Default",
			"resourceName": "projects/test-project/sinks/_Default",
			"destination": "logging.googleapis.com/projects/test-project/locations/global/buckets/_Default",
			"filter": "NOT LOG_ID(\"cloudaudit.googleapis.com/activity\")"
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/projects/test-project/sinks":
			pageToken := r.URL.Query().Get("pageToken")
			mu.Lock()
			pageTokens = append(pageTokens, pageToken)
			mu.Unlock()
			switch pageToken {
			case "":
				_, _ = w.Write([]byte(page1JSON))
			case "page-2":
				_, _ = w.Write([]byte(page2JSON))
			default:
				http.Error(w, "unexpected page token", http.StatusBadRequest)
			}
		case "/v2/projects/empty-project/sinks":
			_, _ = w.Write([]byte(`{}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleLoggingProjectSinks()
	cfg := &transport_tpg.Config{
		Context:         context.Background(),
		Client:          server.Client(),
		Project:         "test-project",
		CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/v2/"},
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/sinks"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	wantFirst := map[string]interface{}{
		"name":            "sink-1",
		"destination":     "bigquery.googleapis.com/projects/test-project/datasets/logs",
		"filter":          "severity>=ERROR",
		"description":     "first sink",
		"disabled":        true,
		"writer_identity": "serviceAccount:service-123@gcp-sa-logging.iam.gserviceaccount.com",
		"exclusions": []interface{}{map[string]interface{}{
			"name":        "exclusion-1",
			"description": "first exclusion",
			"filter":      `resource.type="gce_instance"`,
			"disabled":    true,
		}},
		"bigquery_options": []interface{}{map[string]interface{}{"use_partitioned_tables": true}},
	}
	sinks := d.Get("sinks").([]interface{})
	if len(sinks) != 3 {
		t.Fatalf("got %d sinks, want 3 (two from the first page, one from the second): %#v", len(sinks), sinks)
	}
	if !reflect.DeepEqual(sinks[0], wantFirst) {
		t.Errorf("first sink = %#v, want %#v", sinks[0], wantFirst)
	}
	if got, want := sinks[1].(map[string]interface{})["name"], "sink-2"; got != want {
		t.Errorf("second sink name = %q, want %q", got, want)
	}
	wantLast := map[string]interface{}{
		"name":             "_Default",
		"destination":      "logging.googleapis.com/projects/test-project/locations/global/buckets/_Default",
		"filter":           `NOT LOG_ID("cloudaudit.googleapis.com/activity")`,
		"description":      "",
		"disabled":         false,
		"writer_identity":  "",
		"exclusions":       []interface{}{},
		"bigquery_options": []interface{}{},
	}
	if !reflect.DeepEqual(sinks[2], wantLast) {
		t.Errorf("last sink = %#v, want %#v", sinks[2], wantLast)
	}

	mu.Lock()
	gotTokens := pageTokens
	mu.Unlock()
	if !reflect.DeepEqual(gotTokens, []string{"", "page-2"}) {
		t.Errorf("list requests sent page tokens %q, want the first page and then page-2", gotTokens)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project"})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read of a project with no sinks returned an error: %s", err)
	}
	if got := d.Get("sinks").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no sinks", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "missing-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read of a project the API does not know returned no error")
	}
}
