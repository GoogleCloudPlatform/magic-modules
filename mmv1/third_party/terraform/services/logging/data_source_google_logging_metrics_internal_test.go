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

// TestDataSourceGoogleLoggingMetricsRead checks paging and the field mapping against a fake API.
func TestDataSourceGoogleLoggingMetricsRead(t *testing.T) {
	t.Parallel()

	// The first metric sets every field, so each mapping is checked, and a second metric shares its page.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"metrics": [{
			"name": "metric-1",
			"resourceName": "projects/test-project/metrics/metric-1",
			"description": "first metric",
			"filter": "resource.type=gae_app AND severity>=ERROR",
			"bucketName": "projects/test-project/locations/global/buckets/_Default",
			"disabled": true,
			"metricDescriptor": {
				"name": "projects/test-project/metricDescriptors/logging.googleapis.com/user/metric-1",
				"type": "logging.googleapis.com/user/metric-1",
				"metricKind": "DELTA",
				"valueType": "DISTRIBUTION",
				"unit": "1",
				"displayName": "First metric",
				"labels": [
					{"key": "mass", "valueType": "STRING", "description": "amount of matter"},
					{"key": "sku", "valueType": "INT64"}
				]
			},
			"labelExtractors": {"mass": "EXTRACT(jsonPayload.request)", "sku": "EXTRACT(jsonPayload.id)"},
			"valueExtractor": "EXTRACT(jsonPayload.request)",
			"bucketOptions": {"linearBuckets": {"numFiniteBuckets": 3, "width": 1.5, "offset": 2.5}},
			"createTime": "2026-11-10T00:00:00.000Z",
			"updateTime": "2026-11-10T00:00:00.000Z"
		}, {
			"name": "metric-2",
			"filter": "severity>=WARNING",
			"metricDescriptor": {"metricKind": "DELTA", "valueType": "INT64"}
		}]
	}`
	// The last metric has a slash in its name, which the API allows, and only the fields of a counter.
	const page2JSON = `{
		"metrics": [{
			"name": "nginx/requests",
			"filter": "resource.type=gce_instance",
			"metricDescriptor": {"metricKind": "DELTA", "valueType": "INT64"}
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/projects/test-project/metrics":
			if got, want := r.Header.Get("X-Goog-User-Project"), "billing-project"; got != want {
				t.Errorf("quota project header = %q, want the billing project %q", got, want)
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
		case "/v2/projects/empty-project/metrics":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v2/projects/bad-project/metrics":
			_, _ = w.Write([]byte(`{"metrics": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleLoggingMetrics()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v2/"},
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/metrics"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	metrics := d.Get("metrics").([]interface{})
	if len(metrics) != 3 {
		t.Fatalf("got %d metrics, want 3 (two from the first page, one from the second): %#v", len(metrics), metrics)
	}

	first := metrics[0].(map[string]interface{})
	descriptor := first["metric_descriptor"].([]interface{})[0].(map[string]interface{})
	// The labels are a set, so they are compared as a list of their own and taken out of the descriptor.
	labels := descriptor["labels"].(*schema.Set)
	delete(descriptor, "labels")
	// This comparison is exact on purpose: a field added to Metric.yaml must be added here and in flattenLoggingMetricsItem.
	wantFirst := map[string]interface{}{
		"name":        "metric-1",
		"description": "first metric",
		"filter":      "resource.type=gae_app AND severity>=ERROR",
		"bucket_name": "projects/test-project/locations/global/buckets/_Default",
		"disabled":    true,
		"metric_descriptor": []interface{}{map[string]interface{}{
			"metric_kind":  "DELTA",
			"value_type":   "DISTRIBUTION",
			"unit":         "1",
			"display_name": "First metric",
		}},
		"label_extractors": map[string]interface{}{"mass": "EXTRACT(jsonPayload.request)", "sku": "EXTRACT(jsonPayload.id)"},
		"value_extractor":  "EXTRACT(jsonPayload.request)",
		"bucket_options": []interface{}{map[string]interface{}{
			"linear_buckets":      []interface{}{map[string]interface{}{"num_finite_buckets": 3, "width": 1.5, "offset": 2.5}},
			"exponential_buckets": []interface{}{},
			"explicit_buckets":    []interface{}{},
		}},
	}
	if !reflect.DeepEqual(first, wantFirst) {
		t.Errorf("first metric = %#v, want %#v", first, wantFirst)
	}
	wantLabels := map[string]map[string]interface{}{
		"mass": {"key": "mass", "value_type": "STRING", "description": "amount of matter"},
		"sku":  {"key": "sku", "value_type": "INT64", "description": ""},
	}
	if labels.Len() != len(wantLabels) {
		t.Errorf("first metric has %d descriptor labels, want %d: %#v", labels.Len(), len(wantLabels), labels.List())
	}
	for _, raw := range labels.List() {
		label := raw.(map[string]interface{})
		if want := wantLabels[label["key"].(string)]; !reflect.DeepEqual(label, want) {
			t.Errorf("descriptor label = %#v, want %#v", label, want)
		}
	}

	if got, want := metrics[1].(map[string]interface{})["name"], "metric-2"; got != want {
		t.Errorf("second metric name = %q, want %q", got, want)
	}
	last := metrics[2].(map[string]interface{})
	if got, want := last["name"], "nginx/requests"; got != want {
		t.Errorf("last metric name = %q, want %q", got, want)
	}
	if got, want := last["metric_descriptor"].([]interface{})[0].(map[string]interface{})["metric_kind"], "DELTA"; got != want {
		t.Errorf("last metric kind = %q, want %q", got, want)
	}
	if got := last["bucket_options"].([]interface{}); len(got) != 0 {
		t.Errorf("last metric bucket options = %#v, want none", got)
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
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v2/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project"})
	if err := ds.Read(d, cfgNoBilling); err != nil {
		t.Fatalf("Read of a project with no metrics returned an error: %s", err)
	}
	if got := d.Get("metrics").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no metrics", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "bad-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read accepted a list item that is not an object, want an error")
	}

	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "missing-project"})
	if err := ds.Read(d, cfg); err == nil {
		t.Errorf("Read of a project the API does not know returned no error")
	}
}
