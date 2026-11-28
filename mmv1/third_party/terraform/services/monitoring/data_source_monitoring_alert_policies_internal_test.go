package monitoring

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

// TestDataSourceMonitoringAlertPoliciesRead checks paging and the fields of each alert policy against a fake API.
func TestDataSourceMonitoringAlertPoliciesRead(t *testing.T) {
	t.Parallel()

	// The first policy sets every top-level field with a threshold condition, and the second is a disabled log-based policy.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"totalSize": 3,
		"alertPolicies": [{
			"name": "projects/test-project/alertPolicies/1111",
			"displayName": "first policy",
			"combiner": "OR",
			"enabled": true,
			"severity": "WARNING",
			"userLabels": {"env": "test", "team": "data"},
			"notificationChannels": [
				"projects/test-project/notificationChannels/2222",
				"projects/test-project/notificationChannels/3333"
			],
			"documentation": {
				"content": "what to do when this fires",
				"mimeType": "text/markdown",
				"subject": "disk writes are high",
				"links": [{"displayName": "playbook", "url": "https://example.com/playbook"}]
			},
			"alertStrategy": {
				"autoClose": "1800s",
				"notificationPrompts": ["OPENED"],
				"notificationChannelStrategy": [{
					"notificationChannelNames": ["projects/test-project/notificationChannels/2222"],
					"renotifyInterval": "3600s"
				}]
			},
			"creationRecord": {"mutateTime": "2026-11-10T00:00:00.000000000Z", "mutatedBy": "creator@example.com"},
			"mutationRecord": {"mutateTime": "2026-11-11T00:00:00.000000000Z", "mutatedBy": "editor@example.com"},
			"conditions": [{
				"name": "projects/test-project/alertPolicies/1111/conditions/4444",
				"displayName": "disk writes",
				"conditionThreshold": {
					"filter": "metric.type=\"compute.googleapis.com/instance/disk/write_bytes_count\" AND resource.type=\"gce_instance\"",
					"comparison": "COMPARISON_GT",
					"duration": "60s",
					"thresholdValue": 0.5,
					"aggregations": [{
						"alignmentPeriod": "60s",
						"perSeriesAligner": "ALIGN_RATE",
						"crossSeriesReducer": "REDUCE_SUM",
						"groupByFields": ["resource.label.zone"]
					}],
					"trigger": {"count": 2}
				}
			}]
		}, {
			"name": "projects/test-project/alertPolicies/5555",
			"displayName": "second policy",
			"combiner": "OR",
			"enabled": false,
			"alertStrategy": {"notificationRateLimit": {"period": "300s"}},
			"conditions": [{
				"name": "projects/test-project/alertPolicies/5555/conditions/6666",
				"displayName": "log match",
				"conditionMatchedLog": {
					"filter": "severity>=ERROR",
					"labelExtractors": {"host": "EXTRACT(jsonPayload.host)"}
				}
			}]
		}]
	}`
	const page2JSON = `{
		"totalSize": 3,
		"alertPolicies": [{
			"name": "projects/test-project/alertPolicies/7777",
			"displayName": "third policy",
			"combiner": "AND",
			"enabled": true,
			"conditions": [{
				"name": "projects/test-project/alertPolicies/7777/conditions/8888",
				"displayName": "no data",
				"conditionAbsent": {
					"filter": "metric.type=\"compute.googleapis.com/instance/cpu/utilization\" AND resource.type=\"gce_instance\"",
					"duration": "300s"
				}
			}]
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v3/projects/test-project/alertPolicies":
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
		case "/v3/projects/empty-project/alertPolicies":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v3/projects/bad-project/alertPolicies":
			_, _ = w.Write([]byte(`{"alertPolicies": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceMonitoringAlertPolicies()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/"},
	}

	// These fields of the resource are not part of a policy: project repeats the argument and deletion_policy is a setting of the resource.
	policySchema := ds.Schema["alert_policies"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"project", "deletion_policy"} {
		if _, ok := policySchema[field]; ok {
			t.Errorf("the schema of a policy has the field %q, want it left out", field)
		}
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/alertPolicies"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	const threshold = "alert_policies.0.conditions.0.condition_threshold.0."
	// Only the attributes the fake API sets are compared, in the flat form of the state, so a field added to AlertPolicy.yaml needs no change here.
	want := map[string]string{
		"alert_policies.#":                                         "3",
		"alert_policies.0.id":                                      "projects/test-project/alertPolicies/1111",
		"alert_policies.0.name":                                    "projects/test-project/alertPolicies/1111",
		"alert_policies.0.display_name":                            "first policy",
		"alert_policies.0.combiner":                                "OR",
		"alert_policies.0.enabled":                                 "true",
		"alert_policies.0.severity":                                "WARNING",
		"alert_policies.0.user_labels.%":                           "2",
		"alert_policies.0.user_labels.env":                         "test",
		"alert_policies.0.user_labels.team":                        "data",
		"alert_policies.0.notification_channels.#":                 "2",
		"alert_policies.0.notification_channels.0":                 "projects/test-project/notificationChannels/2222",
		"alert_policies.0.notification_channels.1":                 "projects/test-project/notificationChannels/3333",
		"alert_policies.0.documentation.0.content":                 "what to do when this fires",
		"alert_policies.0.documentation.0.mime_type":               "text/markdown",
		"alert_policies.0.documentation.0.subject":                 "disk writes are high",
		"alert_policies.0.documentation.0.links.0.display_name":    "playbook",
		"alert_policies.0.documentation.0.links.0.url":             "https://example.com/playbook",
		"alert_policies.0.alert_strategy.0.auto_close":             "1800s",
		"alert_policies.0.alert_strategy.0.notification_prompts.0": "OPENED",
		"alert_policies.0.alert_strategy.0.notification_channel_strategy.0.notification_channel_names.0": "projects/test-project/notificationChannels/2222",
		"alert_policies.0.alert_strategy.0.notification_channel_strategy.0.renotify_interval":            "3600s",
		"alert_policies.0.creation_record.0.mutate_time":                                                 "2026-11-10T00:00:00.000000000Z",
		"alert_policies.0.creation_record.0.mutated_by":                                                  "creator@example.com",
		"alert_policies.0.conditions.#":                                                                  "1",
		"alert_policies.0.conditions.0.name":                                                             "projects/test-project/alertPolicies/1111/conditions/4444",
		"alert_policies.0.conditions.0.display_name":                                                     "disk writes",
		threshold + "filter":                              `metric.type="compute.googleapis.com/instance/disk/write_bytes_count" AND resource.type="gce_instance"`,
		threshold + "comparison":                          "COMPARISON_GT",
		threshold + "duration":                            "60s",
		threshold + "threshold_value":                     "0.5",
		threshold + "aggregations.0.alignment_period":     "60s",
		threshold + "aggregations.0.per_series_aligner":   "ALIGN_RATE",
		threshold + "aggregations.0.cross_series_reducer": "REDUCE_SUM",
		threshold + "aggregations.0.group_by_fields.0":    "resource.label.zone",
		threshold + "trigger.0.count":                     "2",

		"alert_policies.1.id":                                                         "projects/test-project/alertPolicies/5555",
		"alert_policies.1.display_name":                                               "second policy",
		"alert_policies.1.enabled":                                                    "false",
		"alert_policies.1.user_labels.%":                                              "0",
		"alert_policies.1.notification_channels.#":                                    "0",
		"alert_policies.1.alert_strategy.0.notification_rate_limit.0.period":          "300s",
		"alert_policies.1.conditions.0.condition_matched_log.0.filter":                "severity>=ERROR",
		"alert_policies.1.conditions.0.condition_matched_log.0.label_extractors.host": "EXTRACT(jsonPayload.host)",
		"alert_policies.1.conditions.0.condition_threshold.#":                         "0",

		"alert_policies.2.id":                                       "projects/test-project/alertPolicies/7777",
		"alert_policies.2.display_name":                             "third policy",
		"alert_policies.2.combiner":                                 "AND",
		"alert_policies.2.conditions.0.condition_absent.0.duration": "300s",
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
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"project": "empty-project"})
	if err := ds.Read(d, cfgNoBilling); err != nil {
		t.Fatalf("Read of a project with no alert policies returned an error: %s", err)
	}
	if got := d.Get("alert_policies").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no alert policies", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}
	if got, want := d.Id(), "projects/empty-project/alertPolicies"; got != want {
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
		CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/"},
	}
	d = schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfgNoProject); err == nil {
		t.Errorf("Read with no project argument and no provider project returned no error")
	}
}
