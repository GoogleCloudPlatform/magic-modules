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

// TestDataSourceMonitoringNotificationChannelsRead checks paging and the fields of each notification channel against a fake API.
func TestDataSourceMonitoringNotificationChannelsRead(t *testing.T) {
	t.Parallel()

	// The first channel sets every field of the resource, and the second holds a label whose value the API has obfuscated.
	const page1JSON = `{
		"nextPageToken": "page-2",
		"totalSize": 3,
		"notificationChannels": [{
			"name": "projects/test-project/notificationChannels/1111",
			"type": "email",
			"displayName": "on-call email",
			"description": "reaches the engineer on call",
			"enabled": true,
			"verificationStatus": "VERIFIED",
			"labels": {"email_address": "oncall@example.com"},
			"userLabels": {"env": "test", "team": "data"},
			"creationRecord": {"mutateTime": "2026-11-10T00:00:00.000000000Z", "mutatedBy": "creator@example.com"},
			"mutationRecords": [{"mutateTime": "2026-11-11T00:00:00.000000000Z", "mutatedBy": "editor@example.com"}]
		}, {
			"name": "projects/test-project/notificationChannels/2222",
			"type": "slack",
			"displayName": "alerts channel",
			"enabled": true,
			"labels": {"channel_name": "#alerts", "auth_token": "**********abcd"}
		}]
	}`
	const page2JSON = `{
		"totalSize": 3,
		"notificationChannels": [{
			"name": "projects/test-project/notificationChannels/3333",
			"type": "sms",
			"displayName": "on-call phone",
			"enabled": false,
			"verificationStatus": "UNVERIFIED",
			"labels": {"number": "+15555550100"}
		}]
	}`

	var mu sync.Mutex
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v3/projects/test-project/notificationChannels":
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
		case "/v3/projects/empty-project/notificationChannels":
			if got, want := r.Header.Get("X-Goog-User-Project"), "empty-project"; got != want {
				t.Errorf("quota project header = %q, want the project %q when no billing project is set", got, want)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v3/projects/bad-project/notificationChannels":
			_, _ = w.Write([]byte(`{"notificationChannels": ["not-an-object"]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceMonitoringNotificationChannels()
	cfg := &transport_tpg.Config{
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "test-project",
		BillingProject:      "billing-project",
		UserProjectOverride: true,
		CustomEndpoints:     map[string]string{Product.CustomEndpointField: server.URL + "/"},
	}

	// These fields of the resource are not part of a channel: project repeats the argument, and the others are settings the API never returns.
	channelSchema := ds.Schema["notification_channels"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"project", "deletion_policy", "force_delete", "sensitive_labels"} {
		if _, ok := channelSchema[field]; ok {
			t.Errorf("the schema of a channel has the field %q, want it left out", field)
		}
	}

	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read returned an error: %s", err)
	}

	if got, want := d.Id(), "projects/test-project/notificationChannels"; got != want {
		t.Errorf("id = %q, want %q", got, want)
	}
	if got, want := d.Get("project"), "test-project"; got != want {
		t.Errorf("project = %q, want %q", got, want)
	}

	// Only the attributes the fake API sets are compared, in the flat form of the state, so a field added to NotificationChannel.yaml needs no change here.
	want := map[string]string{
		"notification_channels.#":                      "3",
		"notification_channels.0.id":                   "projects/test-project/notificationChannels/1111",
		"notification_channels.0.name":                 "projects/test-project/notificationChannels/1111",
		"notification_channels.0.type":                 "email",
		"notification_channels.0.display_name":         "on-call email",
		"notification_channels.0.description":          "reaches the engineer on call",
		"notification_channels.0.enabled":              "true",
		"notification_channels.0.verification_status":  "VERIFIED",
		"notification_channels.0.labels.%":             "1",
		"notification_channels.0.labels.email_address": "oncall@example.com",
		"notification_channels.0.user_labels.%":        "2",
		"notification_channels.0.user_labels.env":      "test",
		"notification_channels.0.user_labels.team":     "data",

		"notification_channels.1.id":                  "projects/test-project/notificationChannels/2222",
		"notification_channels.1.type":                "slack",
		"notification_channels.1.display_name":        "alerts channel",
		"notification_channels.1.labels.%":            "2",
		"notification_channels.1.labels.channel_name": "#alerts",
		"notification_channels.1.labels.auth_token":   "**********abcd",
		"notification_channels.1.user_labels.%":       "0",

		"notification_channels.2.id":                  "projects/test-project/notificationChannels/3333",
		"notification_channels.2.name":                "projects/test-project/notificationChannels/3333",
		"notification_channels.2.type":                "sms",
		"notification_channels.2.enabled":             "false",
		"notification_channels.2.verification_status": "UNVERIFIED",
		"notification_channels.2.labels.number":       "+15555550100",
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
		t.Fatalf("Read of a project with no notification channels returned an error: %s", err)
	}
	if got := d.Get("notification_channels").([]interface{}); len(got) != 0 {
		t.Errorf("empty project returned %#v, want no notification channels", got)
	}
	if got, want := d.Get("project"), "empty-project"; got != want {
		t.Errorf("project = %q, want the project argument %q", got, want)
	}
	if got, want := d.Id(), "projects/empty-project/notificationChannels"; got != want {
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
