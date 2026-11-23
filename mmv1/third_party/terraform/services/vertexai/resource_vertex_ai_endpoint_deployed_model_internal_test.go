package vertexai

import (
	"io"
	"net/http"
	"strings"
	"testing"

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
