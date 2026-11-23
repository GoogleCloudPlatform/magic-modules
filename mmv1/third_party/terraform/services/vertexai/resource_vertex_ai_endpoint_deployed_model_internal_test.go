package vertexai

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type endpointDeployedModelRoundTripper func(*http.Request) (*http.Response, error)

func (f endpointDeployedModelRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func endpointDeployedModelTestData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, ResourceVertexAIEndpointDeployedModel().Schema, map[string]interface{}{
		"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
		"deployed_model_id": "456",
		"model":             "projects/test-project/locations/us-central1/models/example",
		"automatic_resources": []interface{}{
			map[string]interface{}{"min_replica_count": 1},
		},
	})
}

func TestVertexAIEndpointDeployedModelCreateUsesEndpointRegion(t *testing.T) {
	for _, providerRegion := range []string{"", "europe-west1"} {
		t.Run("provider_region_"+providerRegion, func(t *testing.T) {
			d := endpointDeployedModelTestData(t)
			var hosts []string
			client := &http.Client{Transport: endpointDeployedModelRoundTripper(func(req *http.Request) (*http.Response, error) {
				hosts = append(hosts, req.URL.Host)
				body := `{}`
				if req.Method == http.MethodGet {
					body = `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","automaticResources":{"minReplicaCount":1}}]}`
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			})}
			config := &transport_tpg.Config{Project: "test-project", Region: providerRegion, Client: client}

			if err := resourceVertexAIEndpointDeployedModelCreate(d, config); err != nil {
				t.Fatalf("create failed: %v", err)
			}
			if got := d.Get("region").(string); got != "us-central1" {
				t.Errorf("region = %q, want us-central1", got)
			}
			if len(hosts) != 2 {
				t.Fatalf("got %d requests, want create and read", len(hosts))
			}
			for _, host := range hosts {
				if host != "us-central1-aiplatform.googleapis.com" {
					t.Errorf("request host = %q, want us-central1-aiplatform.googleapis.com", host)
				}
			}
		})
	}
}

func TestVertexAIEndpointDeployedModelCreateRejectsMissingEndpointRegion(t *testing.T) {
	d := endpointDeployedModelTestData(t)
	if err := d.Set("endpoint", "invalid-endpoint"); err != nil {
		t.Fatal(err)
	}
	config := &transport_tpg.Config{Project: "test-project"}

	err := resourceVertexAIEndpointDeployedModelCreate(d, config)
	if err == nil || !strings.Contains(err.Error(), "Cannot determine region from endpoint") {
		t.Fatalf("got error %v, want an endpoint region error", err)
	}
}

func TestVertexAIEndpointDeployedModelQuotaProject(t *testing.T) {
	tests := []struct {
		name           string
		override       bool
		billingProject string
		wantProject    string
	}{
		{name: "endpoint project", override: true, wantProject: "test-project"},
		{name: "billing project", override: true, billingProject: "billing-project", wantProject: "billing-project"},
		{name: "override disabled", billingProject: "billing-project"},
	}
	for _, tt := range tests {
		for _, action := range []string{"create", "delete"} {
			t.Run(tt.name+"/"+action, func(t *testing.T) {
				d := endpointDeployedModelTestData(t)
				if err := d.Set("region", "us-central1"); err != nil {
					t.Fatal(err)
				}
				d.SetId("projects/test-project/locations/us-central1/endpoints/123/deployedModels/456")
				var requests int
				client := &http.Client{Transport: endpointDeployedModelRoundTripper(func(req *http.Request) (*http.Response, error) {
					requests++
					if got := req.Header.Get("X-Goog-User-Project"); got != tt.wantProject {
						t.Errorf("%s %s quota project = %q, want %q", req.Method, req.URL.Path, got, tt.wantProject)
					}
					body := `{"name":"projects/test-project/locations/us-central1/endpoints/123/operations/op1"}`
					if strings.Contains(req.URL.Path, "/operations/") {
						body = `{"name":"projects/test-project/locations/us-central1/endpoints/123/operations/op1","done":true}`
					} else if req.Method == http.MethodGet {
						body = `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","automaticResources":{"minReplicaCount":1}}]}`
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})}
				config := &transport_tpg.Config{
					Project:             "provider-project",
					BillingProject:      tt.billingProject,
					UserProjectOverride: tt.override,
					Client:              client,
					PollInterval:        time.Millisecond,
				}

				var err error
				if action == "create" {
					err = resourceVertexAIEndpointDeployedModelCreate(d, config)
				} else {
					err = resourceVertexAIEndpointDeployedModelDelete(d, config)
				}
				if err != nil {
					t.Fatalf("%s failed: %v", action, err)
				}
				wantRequests := 2
				if action == "create" {
					wantRequests++
				}
				if requests != wantRequests {
					t.Errorf("got %d requests, want %d", requests, wantRequests)
				}
			})
		}
	}
}
