package cloudrunv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// The update PATCH replaces the whole revision template, so an action that the
// configuration omits must be resent from refreshed state rather than dropped.
func TestCloudRunV2Service_encryptionKeyRevocationActionUpdate(t *testing.T) {
	t.Parallel()

	const (
		project  = "my-project"
		location = "us-central1"
		name     = "my-service"
		image    = "us-docker.pkg.dev/cloudrun/container/hello"
		key      = "projects/my-project/locations/us-central1/keyRings/my-ring/cryptoKeys/my-key"
	)
	id := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, location, name)

	cases := []struct {
		name       string
		refreshed  string // action returned by the API before the update
		configured string // action in configuration; empty means omitted
		want       string // action sent in the PATCH request
	}{
		{name: "omitted keeps PREVENT_NEW", refreshed: "PREVENT_NEW", want: "PREVENT_NEW"},
		{name: "omitted keeps SHUTDOWN", refreshed: "SHUTDOWN", want: "SHUTDOWN"},
		{name: "configured overrides refreshed", refreshed: "SHUTDOWN", configured: "PREVENT_NEW", want: "PREVENT_NEW"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var mu sync.Mutex
			var patched map[string]interface{}
			service := map[string]interface{}{
				"name": id,
				"template": map[string]interface{}{
					"timeout":                       "300s",
					"encryptionKey":                 key,
					"encryptionKeyRevocationAction": tc.refreshed,
					"containers":                    []interface{}{map[string]interface{}{"image": image}},
				},
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v2/"+id {
					http.NotFound(w, r)
					return
				}
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodGet:
					_ = json.NewEncoder(w).Encode(service)
				case http.MethodPatch:
					if err := json.NewDecoder(r.Body).Decode(&patched); err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					service["template"] = patched["template"]
					_ = json.NewEncoder(w).Encode(map[string]interface{}{"name": "operations/update", "done": true})
				default:
					http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
				}
			}))
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			if err != nil {
				t.Fatalf("parsing test server URL: %v", err)
			}
			// Fail locally instead of reaching Google APIs if a request misses the
			// custom endpoint.
			client := server.Client()
			loopback := client.Transport
			client.Transport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Host != serverURL.Host {
					return nil, fmt.Errorf("request to %s is outside the test server", r.URL)
				}
				return loopback.RoundTrip(r)
			})
			config := &transport_tpg.Config{
				Context:         context.Background(),
				Client:          client,
				Project:         project,
				CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/v2/"},
			}

			res := ResourceCloudRunV2Service()
			prior, diags := res.RefreshWithoutUpgrade(context.Background(), &terraform.InstanceState{
				ID:         id,
				Attributes: map[string]string{"id": id, "project": project, "location": location, "name": name},
			}, config)
			if diags.HasError() {
				t.Fatalf("refresh: %v", diags)
			}
			if got := prior.Attributes["template.0.encryption_key_revocation_action"]; got != tc.refreshed {
				t.Fatalf("refreshed action = %q, want %q", got, tc.refreshed)
			}

			template := map[string]interface{}{
				"timeout":        "301s",
				"encryption_key": key,
				"containers":     []interface{}{map[string]interface{}{"image": image}},
			}
			if tc.configured != "" {
				template["encryption_key_revocation_action"] = tc.configured
			}
			// The resource's CustomizeDiff functions read the raw plan, which only
			// Terraform core supplies, and none of them touch the template.
			planner := &schema.Resource{Schema: res.Schema}
			diff, err := planner.SimpleDiff(context.Background(), prior, terraform.NewResourceConfigRaw(map[string]interface{}{
				"project":  project,
				"location": location,
				"name":     name,
				"template": []interface{}{template},
			}), config)
			if err != nil {
				t.Fatalf("plan: %v", err)
			}
			if diff.RequiresNew() {
				t.Fatalf("plan replaces the service: %v", diff.Attributes)
			}

			if _, diags := res.Apply(context.Background(), prior, diff, config); diags.HasError() {
				t.Fatalf("apply: %v", diags)
			}

			mu.Lock()
			defer mu.Unlock()
			if patched == nil {
				t.Fatal("apply sent no PATCH request")
			}
			sent, _ := patched["template"].(map[string]interface{})
			for field, want := range map[string]string{
				"encryptionKeyRevocationAction": tc.want,
				"timeout":                       "301s",
				"encryptionKey":                 key,
			} {
				if got := sent[field]; got != want {
					t.Errorf("PATCH template.%s = %v, want %q", field, got, want)
				}
			}
		})
	}
}
