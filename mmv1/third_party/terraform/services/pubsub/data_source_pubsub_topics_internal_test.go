package pubsub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// TestDataSourceGooglePubsubTopicsRead checks paging and the fields of each topic against a fake API.
func TestDataSourceGooglePubsubTopicsRead(t *testing.T) {
	t.Parallel()

	// The first topic sets every top-level field, and a second topic with only a name shares its page.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"topics": [{
			"name": "projects/test-project/topics/topic-1",
			"kmsKeyName": "projects/test-project/locations/global/keyRings/ring/cryptoKeys/key",
			"labels": {"env": "test", "team": "data"},
			"messageStoragePolicy": {"allowedPersistenceRegions": ["europe-west3", "europe-west1"], "enforceInTransit": true},
			"schemaSettings": {
				"schema": "projects/test-project/schemas/schema-1",
				"encoding": "JSON",
				"firstRevisionId": "rev-1",
				"lastRevisionId": "rev-2"
			},
			"messageRetentionDuration": "86600s",
			"ingestionDataSourceSettings": {
				"cloudStorage": {
					"bucket": "bucket-1",
					"avroFormat": {},
					"minimumObjectCreateTime": "2026-11-10T00:00:00Z",
					"matchGlob": "**.avro"
				},
				"platformLogsSettings": {"severity": "WARNING"}
			},
			"messageTransforms": [{
				"javascriptUdf": {"functionName": "redact", "code": "function redact(message, metadata) { return message; }"},
				"disabled": true
			}],
			"state": "ACTIVE"
		}, {
			"name": "projects/test-project/topics/topic-2"
		}]
	}`
	const page2JSON = `{
		"topics": [{
			"name": "projects/test-project/topics/topic-3",
			"labels": {"env": "prod"}
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/projects/test-project/topics":
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
		case "/v1/projects/empty-project/topics":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v1/projects/bad-project/topics":
			_, _ = w.Write([]byte(`{"topics": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGooglePubsubTopics()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v1/"},
	}

	// These fields of the resource are not part of a topic: the API does not return them, or labels already holds them.
	topicSchema := ds.Schema["topics"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"project", "deletion_policy", "tags", "terraform_labels", "effective_labels"} {
		if _, ok := topicSchema[field]; ok {
			t.Errorf("the schema of a topic has the field %q, want it left out", field)
		}
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/topics"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	topics := d.Get("topics").([]interface{})
	if len(topics) != 3 {
		t.Fatalf("got %d topics, want 3 (two from the first page, one from the second): %#v", len(topics), topics)
	}

	// Only the attributes the fake API sets are compared, in the flat form of the state, so a field added to Topic.yaml needs no change here.
	wantFirst := map[string]string{
		"topics.0.id":                         "projects/test-project/topics/topic-1",
		"topics.0.name":                       "topic-1",
		"topics.0.kms_key_name":               "projects/test-project/locations/global/keyRings/ring/cryptoKeys/key",
		"topics.0.labels.%":                   "2",
		"topics.0.labels.env":                 "test",
		"topics.0.labels.team":                "data",
		"topics.0.message_retention_duration": "86600s",
		"topics.0.message_storage_policy.0.enforce_in_transit":             "true",
		"topics.0.message_storage_policy.0.allowed_persistence_regions.#":  "2",
		"topics.0.schema_settings.0.schema":                                "projects/test-project/schemas/schema-1",
		"topics.0.schema_settings.0.encoding":                              "JSON",
		"topics.0.schema_settings.0.first_revision_id":                     "rev-1",
		"topics.0.schema_settings.0.last_revision_id":                      "rev-2",
		"topics.0.message_transforms.#":                                    "1",
		"topics.0.message_transforms.0.disabled":                           "true",
		"topics.0.message_transforms.0.javascript_udf.0.function_name":     "redact",
		"topics.0.message_transforms.0.javascript_udf.0.code":              "function redact(message, metadata) { return message; }",
		"topics.0.ingestion_data_source_settings.0.cloud_storage.0.bucket": "bucket-1",
		// A block with no fields has nothing but its presence, so it must still be counted.
		"topics.0.ingestion_data_source_settings.0.cloud_storage.0.avro_format.#":              "1",
		"topics.0.ingestion_data_source_settings.0.cloud_storage.0.minimum_object_create_time": "2026-11-10T00:00:00Z",
		"topics.0.ingestion_data_source_settings.0.cloud_storage.0.match_glob":                 "**.avro",
		"topics.0.ingestion_data_source_settings.0.platform_logs_settings.0.severity":          "WARNING",
	}
	attributes := d.State().Attributes
	for key, want := range wantFirst {
		if got, ok := attributes[key]; !ok || got != want {
			t.Errorf("%s = %q (present: %t), want %q", key, got, ok, want)
		}
	}
	// The state keys of a set hold a hash of each value, so the regions are read as a list and sorted.
	storagePolicy := topics[0].(map[string]interface{})["message_storage_policy"].([]interface{})[0].(map[string]interface{})
	regions := storagePolicy["allowed_persistence_regions"].(*schema.Set).List()
	sort.Slice(regions, func(i, j int) bool { return regions[i].(string) < regions[j].(string) })
	if want := []interface{}{"europe-west1", "europe-west3"}; !reflect.DeepEqual(regions, want) {
		t.Errorf("first topic allowed persistence regions = %#v, want %#v", regions, want)
	}

	second := topics[1].(map[string]interface{})
	if got, want := second["id"], "projects/test-project/topics/topic-2"; got != want {
		t.Errorf("second topic id = %q, want %q", got, want)
	}
	if got := second["labels"].(map[string]interface{}); len(got) != 0 {
		t.Errorf("second topic labels = %#v, want none", got)
	}
	if got := second["message_storage_policy"].([]interface{}); len(got) != 0 {
		t.Errorf("second topic message storage policy = %#v, want none", got)
	}
	last := topics[2].(map[string]interface{})
	if got, want := last["name"], "topic-3"; got != want {
		t.Errorf("last topic name = %q, want %q", got, want)
	}
	if got, want := last["labels"], map[string]interface{}{"env": "prod"}; !reflect.DeepEqual(got, want) {
		t.Errorf("last topic labels = %#v, want %#v", got, want)
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
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v1/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project"})
	if err := ds.Read(d, cfgNoBilling); err != nil {
		t.Fatalf("Read of a project with no topics returned an error: %s", err)
	}
	if got := d.Get("topics").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no topics", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}
	if got, want := d.Id(), "projects/empty-project/topics"; got != want {
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
		CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/v1/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfgNoProject); err == nil {
		t.Errorf("Read with no project argument and no provider project returned no error")
	}
}
