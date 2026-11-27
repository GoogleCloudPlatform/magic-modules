package pubsub

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

// TestDataSourceGooglePubsubSubscriptionsRead checks paging and the fields of each subscription against a fake API.
func TestDataSourceGooglePubsubSubscriptionsRead(t *testing.T) {
	t.Parallel()

	// The first subscription is a push subscription with every common field, and the second writes to BigQuery and never expires.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"subscriptions": [{
			"name": "projects/test-project/subscriptions/sub-1",
			"topic": "projects/test-project/topics/topic-1",
			"labels": {"env": "test", "team": "data"},
			"pushConfig": {
				"pushEndpoint": "https://example.com/push",
				"attributes": {"x-goog-version": "v1"},
				"oidcToken": {"serviceAccountEmail": "push@test-project.iam.gserviceaccount.com", "audience": "https://example.com"},
				"noWrapper": {"writeMetadata": true}
			},
			"ackDeadlineSeconds": 20,
			"messageRetentionDuration": "1200s",
			"retainAckedMessages": true,
			"expirationPolicy": {"ttl": "300000.5s"},
			"filter": "attributes.foo = \"bar\"",
			"deadLetterPolicy": {"deadLetterTopic": "projects/test-project/topics/dead-letter", "maxDeliveryAttempts": 10},
			"retryPolicy": {"minimumBackoff": "10s", "maximumBackoff": "600s"},
			"enableMessageOrdering": true,
			"enableExactlyOnceDelivery": true,
			"messageTransforms": [{
				"javascriptUdf": {"functionName": "redact", "code": "function redact(message, metadata) { return message; }"},
				"disabled": true
			}],
			"state": "ACTIVE"
		}, {
			"name": "projects/test-project/subscriptions/sub-2",
			"topic": "projects/other-project/topics/topic-2",
			"bigqueryConfig": {
				"table": "test-project.dataset.table",
				"useTopicSchema": true,
				"writeMetadata": true,
				"dropUnknownFields": true,
				"serviceAccountEmail": "bigquery@test-project.iam.gserviceaccount.com"
			},
			"ackDeadlineSeconds": 10,
			"expirationPolicy": {}
		}]
	}`
	// The third subscription writes to Cloud Storage, with 64-bit numbers as text, and the last to Bigtable after its topic was deleted.
	const page2JSON = `{
		"subscriptions": [{
			"name": "projects/test-project/subscriptions/sub-3",
			"topic": "projects/test-project/topics/topic-1",
			"labels": {"env": "prod"},
			"cloudStorageConfig": {
				"bucket": "bucket-1",
				"filenamePrefix": "pre-",
				"filenameSuffix": ".avro",
				"filenameDatetimeFormat": "YYYY-MM-DD/hh_mm_ssZ",
				"maxDuration": "600s",
				"maxBytes": "1000",
				"maxMessages": "100",
				"state": "ACTIVE",
				"avroConfig": {"writeMetadata": true, "useTopicSchema": true},
				"serviceAccountEmail": "storage@test-project.iam.gserviceaccount.com"
			}
		}, {
			"name": "projects/test-project/subscriptions/sub-4",
			"topic": "_deleted-topic_",
			"bigtableConfig": {
				"table": "projects/test-project/instances/instance-1/tables/table-1",
				"columnFamilyMapping": {"rowKeySchema": {}},
				"state": "ACTIVE"
			}
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/projects/test-project/subscriptions":
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
		case "/v1/projects/empty-project/subscriptions":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v1/projects/bad-project/subscriptions":
			_, _ = w.Write([]byte(`{"subscriptions": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGooglePubsubSubscriptions()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/v1/"},
	}

	// These fields of the resource are not part of a subscription: the API does not return them, or labels already holds them.
	subscriptionSchema := ds.Schema["subscriptions"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"project", "deletion_policy", "tags", "terraform_labels", "effective_labels"} {
		if _, ok := subscriptionSchema[field]; ok {
			t.Errorf("the schema of a subscription has the field %q, want it left out", field)
		}
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/subscriptions"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	// Only the attributes the fake API sets are compared, in the flat form of the state, so a field added to Subscription.yaml needs no change here.
	want := map[string]string{
		"subscriptions.#":                                                     "4",
		"subscriptions.0.id":                                                  "projects/test-project/subscriptions/sub-1",
		"subscriptions.0.name":                                                "sub-1",
		"subscriptions.0.topic":                                               "projects/test-project/topics/topic-1",
		"subscriptions.0.labels.%":                                            "2",
		"subscriptions.0.labels.env":                                          "test",
		"subscriptions.0.labels.team":                                         "data",
		"subscriptions.0.push_config.0.push_endpoint":                         "https://example.com/push",
		"subscriptions.0.push_config.0.attributes.x-goog-version":             "v1",
		"subscriptions.0.push_config.0.oidc_token.0.service_account_email":    "push@test-project.iam.gserviceaccount.com",
		"subscriptions.0.push_config.0.oidc_token.0.audience":                 "https://example.com",
		"subscriptions.0.push_config.0.no_wrapper.0.write_metadata":           "true",
		"subscriptions.0.ack_deadline_seconds":                                "20",
		"subscriptions.0.message_retention_duration":                          "1200s",
		"subscriptions.0.retain_acked_messages":                               "true",
		"subscriptions.0.expiration_policy.0.ttl":                             "300000.5s",
		"subscriptions.0.filter":                                              `attributes.foo = "bar"`,
		"subscriptions.0.dead_letter_policy.0.dead_letter_topic":              "projects/test-project/topics/dead-letter",
		"subscriptions.0.dead_letter_policy.0.max_delivery_attempts":          "10",
		"subscriptions.0.retry_policy.0.minimum_backoff":                      "10s",
		"subscriptions.0.retry_policy.0.maximum_backoff":                      "600s",
		"subscriptions.0.enable_message_ordering":                             "true",
		"subscriptions.0.enable_exactly_once_delivery":                        "true",
		"subscriptions.0.message_transforms.#":                                "1",
		"subscriptions.0.message_transforms.0.disabled":                       "true",
		"subscriptions.0.message_transforms.0.javascript_udf.0.function_name": "redact",
		"subscriptions.0.message_transforms.0.javascript_udf.0.code":          "function redact(message, metadata) { return message; }",
		"subscriptions.0.bigquery_config.#":                                   "0",
		"subscriptions.0.cloud_storage_config.#":                              "0",

		"subscriptions.1.id":                                      "projects/test-project/subscriptions/sub-2",
		"subscriptions.1.name":                                    "sub-2",
		"subscriptions.1.topic":                                   "projects/other-project/topics/topic-2",
		"subscriptions.1.labels.%":                                "0",
		"subscriptions.1.ack_deadline_seconds":                    "10",
		"subscriptions.1.bigquery_config.0.table":                 "test-project.dataset.table",
		"subscriptions.1.bigquery_config.0.use_topic_schema":      "true",
		"subscriptions.1.bigquery_config.0.write_metadata":        "true",
		"subscriptions.1.bigquery_config.0.drop_unknown_fields":   "true",
		"subscriptions.1.bigquery_config.0.service_account_email": "bigquery@test-project.iam.gserviceaccount.com",
		// An expiration policy with no ttl means the subscription never expires, so the block must still be counted.
		"subscriptions.1.expiration_policy.#":     "1",
		"subscriptions.1.expiration_policy.0.ttl": "",
		"subscriptions.1.push_config.#":           "0",

		"subscriptions.2.id":                                                    "projects/test-project/subscriptions/sub-3",
		"subscriptions.2.name":                                                  "sub-3",
		"subscriptions.2.labels.%":                                              "1",
		"subscriptions.2.labels.env":                                            "prod",
		"subscriptions.2.cloud_storage_config.0.bucket":                         "bucket-1",
		"subscriptions.2.cloud_storage_config.0.filename_prefix":                "pre-",
		"subscriptions.2.cloud_storage_config.0.filename_suffix":                ".avro",
		"subscriptions.2.cloud_storage_config.0.filename_datetime_format":       "YYYY-MM-DD/hh_mm_ssZ",
		"subscriptions.2.cloud_storage_config.0.max_duration":                   "600s",
		"subscriptions.2.cloud_storage_config.0.max_bytes":                      "1000",
		"subscriptions.2.cloud_storage_config.0.max_messages":                   "100",
		"subscriptions.2.cloud_storage_config.0.state":                          "ACTIVE",
		"subscriptions.2.cloud_storage_config.0.avro_config.0.write_metadata":   "true",
		"subscriptions.2.cloud_storage_config.0.avro_config.0.use_topic_schema": "true",
		"subscriptions.2.cloud_storage_config.0.service_account_email":          "storage@test-project.iam.gserviceaccount.com",
		"subscriptions.2.expiration_policy.#":                                   "0",

		"subscriptions.3.id":                      "projects/test-project/subscriptions/sub-4",
		"subscriptions.3.topic":                   "_deleted-topic_",
		"subscriptions.3.bigtable_config.0.table": "projects/test-project/instances/instance-1/tables/table-1",
		"subscriptions.3.bigtable_config.0.state": "ACTIVE",
		// row_key_schema has no fields, so its presence is its whole value and the block must still be counted.
		"subscriptions.3.bigtable_config.0.column_family_mapping.0.row_key_schema.#": "1",
		"subscriptions.3.bigtable_config.0.column_family_mapping.0.delimited_key.#":  "0",
	}
	attributes := d.State().Attributes
	for key, wantValue := range want {
		if got, ok := attributes[key]; !ok || got != wantValue {
			t.Errorf("%s = %q (present: %t), want %q", key, got, ok, wantValue)
		}
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
		t.Fatalf("Read of a project with no subscriptions returned an error: %s", err)
	}
	if got := d.Get("subscriptions").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no subscriptions", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}
	if got, want := d.Id(), "projects/empty-project/subscriptions"; got != want {
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
