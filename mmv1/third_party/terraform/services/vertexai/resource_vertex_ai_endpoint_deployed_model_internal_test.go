package vertexai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		for _, action := range []string{"create", "update", "delete"} {
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
				if action == "update" {
					fields := func(min int) map[string]interface{} {
						return map[string]interface{}{
							"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
							"region":            "us-central1",
							"deployed_model_id": "456",
							"model":             "projects/test-project/locations/us-central1/models/example",
							"automatic_resources": []interface{}{map[string]interface{}{
								"min_replica_count": min,
							}},
						}
					}
					d, _ = endpointDeployedModelUpdateData(t, fields(1), fields(2), config)
				}

				var err error
				switch action {
				case "create":
					err = resourceVertexAIEndpointDeployedModelCreate(d, config)
				case "update":
					err = resourceVertexAIEndpointDeployedModelUpdate(d, config)
				case "delete":
					err = resourceVertexAIEndpointDeployedModelDelete(d, config)
				}
				if err != nil {
					t.Fatalf("%s failed: %v", action, err)
				}
				wantRequests := 2
				if action != "delete" {
					wantRequests++
				}
				if requests != wantRequests {
					t.Errorf("got %d requests, want %d", requests, wantRequests)
				}
			})
		}
	}
}

func TestVertexAIEndpointDeployedModelEmptyResourceBlocks(t *testing.T) {
	tests := []struct {
		name       string
		fields     map[string]interface{}
		want       map[string]interface{}
		response   string
		stateKey   string
		stateValue interface{}
	}{
		{
			name: "empty automatic resources",
			fields: map[string]interface{}{
				"automatic_resources": []interface{}{map[string]interface{}{}},
			},
			want:       map[string]interface{}{"automaticResources": map[string]interface{}{}},
			response:   `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","automaticResources":{}}]}`,
			stateKey:   "automatic_resources.#",
			stateValue: 1,
		},
		{
			name: "empty machine spec",
			fields: map[string]interface{}{
				"dedicated_resources": []interface{}{map[string]interface{}{
					"machine_spec":      []interface{}{map[string]interface{}{}},
					"min_replica_count": 1,
				}},
			},
			want: map[string]interface{}{"dedicatedResources": map[string]interface{}{
				"machineSpec": map[string]interface{}{}, "minReplicaCount": float64(1),
			}},
			response:   `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","dedicatedResources":{"machineSpec":{"machineType":"n1-standard-2"},"minReplicaCount":1}}]}`,
			stateKey:   "dedicated_resources.0.machine_spec.0.machine_type",
			stateValue: "n1-standard-2",
		},
		{
			name:       "shared resources",
			fields:     map[string]interface{}{"shared_resources": "projects/test-project/locations/us-central1/deploymentResourcePools/pool"},
			want:       map[string]interface{}{"sharedResources": "projects/test-project/locations/us-central1/deploymentResourcePools/pool"},
			response:   `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","sharedResources":"projects/test-project/locations/us-central1/deploymentResourcePools/pool"}]}`,
			stateKey:   "shared_resources",
			stateValue: "projects/test-project/locations/us-central1/deploymentResourcePools/pool",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields["endpoint"] = "projects/test-project/locations/us-central1/endpoints/123"
			tt.fields["deployed_model_id"] = "456"
			tt.fields["model"] = "projects/test-project/locations/us-central1/models/example"
			d := schema.TestResourceDataRaw(t, ResourceVertexAIEndpointDeployedModel().Schema, tt.fields)
			var body map[string]interface{}
			client := &http.Client{Transport: endpointDeployedModelRoundTripper(func(req *http.Request) (*http.Response, error) {
				response := tt.response
				if req.Method == http.MethodPost {
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Fatalf("decode request: %v", err)
					}
					response = `{}`
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(response)),
					Request:    req,
				}, nil
			})}
			if err := resourceVertexAIEndpointDeployedModelCreate(d, &transport_tpg.Config{Client: client}); err != nil {
				t.Fatalf("create failed: %v", err)
			}
			model, ok := body["deployedModel"].(map[string]interface{})
			if !ok {
				t.Fatalf("missing deployedModel in request: %#v", body)
			}
			for _, key := range []string{"automaticResources", "dedicatedResources", "sharedResources"} {
				got, gotOK := model[key]
				want, wantOK := tt.want[key]
				if gotOK != wantOK || !reflect.DeepEqual(got, want) {
					t.Errorf("%s = %#v (present %t), want %#v (present %t)", key, got, gotOK, want, wantOK)
				}
			}
			if got := d.Get(tt.stateKey); !reflect.DeepEqual(got, tt.stateValue) {
				t.Errorf("state %s = %#v, want %#v", tt.stateKey, got, tt.stateValue)
			}
			if tt.name == "empty machine spec" {
				diff, err := schema.InternalMap(ResourceVertexAIEndpointDeployedModel().Schema).Diff(
					context.Background(), d.State(), terraform.NewResourceConfigRaw(tt.fields), nil, nil, false,
				)
				if err != nil {
					t.Fatalf("plan after refresh failed: %v", err)
				}
				if diff != nil && diff.Attributes[tt.stateKey] != nil {
					t.Errorf("default machine type caused a plan diff: %#v", diff.Attributes[tt.stateKey])
				}
			}
		})
	}
}

func TestVertexAIEndpointDeployedModelAutoscalingTargetDefault(t *testing.T) {
	const metricName = "aiplatform.googleapis.com/prediction/online/cpu/utilization"
	for _, tt := range []struct {
		name       string
		target     int
		wantTarget float64
	}{
		{name: "default target", wantTarget: 60},
		{name: "explicit target", target: 75, wantTarget: 75},
	} {
		t.Run(tt.name, func(t *testing.T) {
			metric := map[string]interface{}{"metric_name": metricName}
			if tt.target != 0 {
				metric["target"] = tt.target
			}
			raw := map[string]interface{}{
				"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
				"deployed_model_id": "456",
				"model":             "projects/test-project/locations/us-central1/models/example",
				"dedicated_resources": []interface{}{map[string]interface{}{
					"machine_spec":             []interface{}{map[string]interface{}{"machine_type": "n1-standard-2"}},
					"min_replica_count":        1,
					"autoscaling_metric_specs": []interface{}{metric},
				}},
			}
			sm := schema.InternalMap(ResourceVertexAIEndpointDeployedModel().Schema)
			if diags := sm.Validate(terraform.NewResourceConfigRaw(raw)); diags.HasError() {
				t.Fatalf("configuration validation failed: %v", diags)
			}
			d := schema.TestResourceDataRaw(t, ResourceVertexAIEndpointDeployedModel().Schema, raw)
			var body map[string]interface{}
			client := &http.Client{Transport: endpointDeployedModelRoundTripper(func(req *http.Request) (*http.Response, error) {
				response := `{}`
				if req.Method == http.MethodPost {
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						t.Fatalf("decode request: %v", err)
					}
				} else {
					responseBytes, err := json.Marshal(map[string]interface{}{
						"deployedModels": []interface{}{map[string]interface{}{
							"id": "456", "model": raw["model"],
							"dedicatedResources": map[string]interface{}{
								"machineSpec":     map[string]interface{}{"machineType": "n1-standard-2"},
								"minReplicaCount": 1,
								"autoscalingMetricSpecs": []interface{}{map[string]interface{}{
									"metricName": metricName, "target": tt.wantTarget,
								}},
							},
						}},
					})
					if err != nil {
						t.Fatal(err)
					}
					response = string(responseBytes)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: req}, nil
			})}
			if err := resourceVertexAIEndpointDeployedModelCreate(d, &transport_tpg.Config{Client: client}); err != nil {
				t.Fatalf("create failed: %v", err)
			}
			model := body["deployedModel"].(map[string]interface{})
			dedicated := model["dedicatedResources"].(map[string]interface{})
			sentMetric := dedicated["autoscalingMetricSpecs"].([]interface{})[0].(map[string]interface{})
			if got, present := sentMetric["target"]; present != (tt.target != 0) || (present && got != float64(tt.target)) {
				t.Errorf("sent target = %#v (present %t), configured %d", got, present, tt.target)
			}
			const targetKey = "dedicated_resources.0.autoscaling_metric_specs.0.target"
			if got := d.Get(targetKey); got != int(tt.wantTarget) {
				t.Errorf("state target = %v, want %v", got, tt.wantTarget)
			}
			diff, err := sm.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(raw), nil, nil, false)
			if err != nil {
				t.Fatalf("plan after refresh failed: %v", err)
			}
			if diff != nil && diff.Attributes[targetKey] != nil {
				t.Errorf("target caused a plan diff: %#v", diff.Attributes[targetKey])
			}
		})
	}
}

func endpointDeployedModelUpdateData(t *testing.T, before, after map[string]interface{}, config *transport_tpg.Config) (*schema.ResourceData, *terraform.InstanceDiff) {
	t.Helper()
	r := ResourceVertexAIEndpointDeployedModel()
	state := schema.TestResourceDataRaw(t, r.Schema, before)
	state.SetId("projects/test-project/locations/us-central1/endpoints/123/deployedModels/456")
	instanceState := state.State()
	instanceState.RawConfig = cty.ObjectVal(map[string]cty.Value{
		"deletion_policy": cty.NullVal(cty.String),
	})
	diff, err := schema.InternalMap(r.Schema).Diff(
		context.Background(), instanceState, terraform.NewResourceConfigRaw(after),
		r.CustomizeDiff, config, true,
	)
	if err != nil {
		t.Fatalf("plan failed: %v", err)
	}
	data, err := schema.InternalMap(r.Schema).Data(instanceState, diff)
	if err != nil {
		t.Fatalf("resource data from plan failed: %v", err)
	}
	return data, diff
}

func TestVertexAIEndpointDeployedModelOmittedTargetRetainsState(t *testing.T) {
	config := func(target interface{}) map[string]interface{} {
		metric := map[string]interface{}{"metric_name": "aiplatform.googleapis.com/prediction/online/cpu/utilization"}
		if target != nil {
			metric["target"] = target
		}
		return map[string]interface{}{
			"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
			"deployed_model_id": "456",
			"model":             "projects/test-project/locations/us-central1/models/example",
			"dedicated_resources": []interface{}{map[string]interface{}{
				"machine_spec":             []interface{}{map[string]interface{}{"machine_type": "n1-standard-2"}},
				"min_replica_count":        1,
				"autoscaling_metric_specs": []interface{}{metric},
			}},
		}
	}
	d, _ := endpointDeployedModelUpdateData(t, config(75), config(nil), &transport_tpg.Config{})
	if d.HasChange("dedicated_resources.0.autoscaling_metric_specs") {
		t.Fatal("omitting the API-computed target should retain its existing value")
	}
	if got := d.Get("dedicated_resources.0.autoscaling_metric_specs.0.target"); got != 75 {
		t.Fatalf("target after omission = %v, want 75", got)
	}
}

func TestVertexAIEndpointDeployedModelUpdatesReplicasInPlace(t *testing.T) {
	base := func(min, max int) map[string]interface{} {
		return map[string]interface{}{
			"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
			"region":            "us-central1",
			"deployed_model_id": "456",
			"model":             "projects/test-project/locations/us-central1/models/example",
			"automatic_resources": []interface{}{map[string]interface{}{
				"min_replica_count": min, "max_replica_count": max,
			}},
		}
	}
	var body map[string]interface{}
	var requests int
	client := &http.Client{Transport: endpointDeployedModelRoundTripper(func(req *http.Request) (*http.Response, error) {
		requests++
		if got := req.Header.Get("X-Goog-User-Project"); got != "billing-project" {
			t.Errorf("%s %s quota project = %q, want billing-project", req.Method, req.URL.Path, got)
		}
		response := `{"name":"projects/test-project/locations/us-central1/endpoints/123/operations/op1"}`
		switch {
		case req.Method == http.MethodPost:
			if !strings.HasSuffix(req.URL.Path, "/endpoints/123:mutateDeployedModel") || req.URL.Query().Get("updateMask") != "" {
				t.Errorf("unexpected update URL: %s", req.URL)
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatalf("decode update request: %v", err)
			}
		case strings.Contains(req.URL.Path, "/operations/"):
			response = `{"name":"projects/test-project/locations/us-central1/endpoints/123/operations/op1","done":true}`
		default:
			response = `{"deployedModels":[{"id":"456","model":"projects/test-project/locations/us-central1/models/example","automaticResources":{"minReplicaCount":2,"maxReplicaCount":4}}]}`
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response)), Request: req}, nil
	})}
	config := &transport_tpg.Config{
		Project: "provider-project", BillingProject: "billing-project", UserProjectOverride: true,
		Client: client, PollInterval: time.Millisecond,
	}
	d, diff := endpointDeployedModelUpdateData(t, base(1, 2), base(2, 4), config)
	if diff == nil || diff.RequiresNew() {
		t.Fatalf("replica update should be in-place, got diff %#v", diff)
	}
	if err := resourceVertexAIEndpointDeployedModelUpdate(d, config); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	want := map[string]interface{}{
		"deployedModel": map[string]interface{}{
			"id":                 "456",
			"automaticResources": map[string]interface{}{"minReplicaCount": float64(2), "maxReplicaCount": float64(4)},
		},
		"updateMask": "automaticResources.minReplicaCount,automaticResources.maxReplicaCount",
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("update body = %#v, want %#v", body, want)
	}
	if requests != 3 {
		t.Errorf("got %d requests, want update, poll, and read", requests)
	}
}

func TestVertexAIEndpointDeployedModelMutationPayloads(t *testing.T) {
	const metricName = "aiplatform.googleapis.com/prediction/online/cpu/utilization"
	dedicated := func(min, max, required int, target int) map[string]interface{} {
		fields := map[string]interface{}{
			"machine_spec":           []interface{}{map[string]interface{}{"machine_type": "n1-standard-2"}},
			"min_replica_count":      min,
			"max_replica_count":      max,
			"required_replica_count": required,
		}
		if target != 0 {
			fields["autoscaling_metric_specs"] = []interface{}{map[string]interface{}{"metric_name": metricName, "target": target}}
		}
		return map[string]interface{}{
			"endpoint":            "projects/test-project/locations/us-central1/endpoints/123",
			"region":              "us-central1",
			"deployed_model_id":   "456",
			"model":               "projects/test-project/locations/us-central1/models/example",
			"dedicated_resources": []interface{}{fields},
		}
	}
	tests := []struct {
		name      string
		before    map[string]interface{}
		after     map[string]interface{}
		wantModel map[string]interface{}
		wantMask  string
	}{
		{
			name:   "dedicated replicas and metrics",
			before: dedicated(1, 2, 1, 60),
			after:  dedicated(2, 4, 2, 75),
			wantModel: map[string]interface{}{
				"id": "456", "dedicatedResources": map[string]interface{}{
					"minReplicaCount": 2, "maxReplicaCount": 4, "requiredReplicaCount": 2,
					"autoscalingMetricSpecs": []interface{}{map[string]interface{}{"metricName": metricName, "target": 75}},
				},
			},
			wantMask: "dedicatedResources.minReplicaCount,dedicatedResources.maxReplicaCount,dedicatedResources.requiredReplicaCount,dedicatedResources.autoscalingMetricSpecs",
		},
		{
			name:   "clear metrics",
			before: dedicated(1, 2, 1, 60),
			after:  dedicated(1, 2, 1, 0),
			wantModel: map[string]interface{}{
				"id": "456", "dedicatedResources": map[string]interface{}{"autoscalingMetricSpecs": []interface{}{}},
			},
			wantMask: "dedicatedResources.autoscalingMetricSpecs",
		},
	}
	loggingField, loggingAPI := "disable_container_logging", "disableContainerLogging"
	if _, beta := ResourceVertexAIEndpointDeployedModel().Schema["enable_container_logging"]; beta {
		loggingField, loggingAPI = "enable_container_logging", "enableContainerLogging"
	}
	for _, value := range []bool{false, true} {
		before, after := dedicated(1, 2, 1, 60), dedicated(1, 2, 1, 60)
		before[loggingField], after[loggingField] = !value, value
		tests = append(tests, struct {
			name      string
			before    map[string]interface{}
			after     map[string]interface{}
			wantModel map[string]interface{}
			wantMask  string
		}{
			name: fmt.Sprintf("%s to %t", loggingField, value), before: before, after: after,
			wantModel: map[string]interface{}{"id": "456", loggingAPI: value}, wantMask: loggingAPI,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &transport_tpg.Config{Project: "test-project", Region: "us-central1"}
			d, diff := endpointDeployedModelUpdateData(t, tt.before, tt.after, config)
			if diff == nil || diff.RequiresNew() {
				t.Fatalf("change should be in-place, got diff %#v", diff)
			}
			body, err := resourceVertexAIEndpointDeployedModelUpdateEncoder(d, config, nil)
			if err != nil {
				t.Fatalf("encode update: %v", err)
			}
			want := map[string]interface{}{"deployedModel": tt.wantModel, "updateMask": tt.wantMask}
			if !reflect.DeepEqual(body, want) {
				t.Errorf("update body = %#v, want %#v", body, want)
			}
		})
	}
}

func TestVertexAIEndpointDeployedModelImmutableChangesReplace(t *testing.T) {
	automatic := func() map[string]interface{} {
		return map[string]interface{}{
			"endpoint":          "projects/test-project/locations/us-central1/endpoints/123",
			"region":            "us-central1",
			"deployed_model_id": "456",
			"model":             "projects/test-project/locations/us-central1/models/example",
			"automatic_resources": []interface{}{map[string]interface{}{
				"min_replica_count": 1, "max_replica_count": 2,
			}},
		}
	}
	dedicated := func() map[string]interface{} {
		fields := automatic()
		delete(fields, "automatic_resources")
		fields["dedicated_resources"] = []interface{}{map[string]interface{}{
			"machine_spec":      []interface{}{map[string]interface{}{"machine_type": "n1-standard-2"}},
			"min_replica_count": 1, "max_replica_count": 2,
		}}
		return fields
	}
	shared := func() map[string]interface{} {
		fields := automatic()
		delete(fields, "automatic_resources")
		fields["shared_resources"] = "projects/test-project/locations/us-central1/deploymentResourcePools/pool"
		return fields
	}
	type testCase struct {
		name          string
		before, after map[string]interface{}
	}
	tests := []testCase{
		{name: "automatic to dedicated", before: automatic(), after: dedicated()},
		{name: "dedicated to automatic", before: dedicated(), after: automatic()},
		{name: "dedicated to shared", before: dedicated(), after: shared()},
		{name: "shared to automatic", before: shared(), after: automatic()},
	}
	changedModel := automatic()
	changedModel["model"] = "projects/test-project/locations/us-central1/models/other"
	tests = append(tests, testCase{name: "model", before: automatic(), after: changedModel})
	changedMachine := dedicated()
	changedMachine["dedicated_resources"].([]interface{})[0].(map[string]interface{})["machine_spec"] =
		[]interface{}{map[string]interface{}{"machine_type": "n1-standard-4"}}
	tests = append(tests, testCase{name: "machine", before: dedicated(), after: changedMachine})
	changedSpot := dedicated()
	changedSpot["dedicated_resources"].([]interface{})[0].(map[string]interface{})["spot"] = true
	tests = append(tests, testCase{name: "spot", before: dedicated(), after: changedSpot})
	if _, beta := ResourceVertexAIEndpointDeployedModel().Schema["enable_container_logging"]; beta {
		before, after := automatic(), automatic()
		before["disable_container_logging"], after["disable_container_logging"] = false, true
		tests = append(tests, testCase{name: "beta disable container logging", before: before, after: after})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, diff := endpointDeployedModelUpdateData(t, tt.before, tt.after, &transport_tpg.Config{})
			if diff == nil || !diff.RequiresNew() {
				t.Errorf("change must replace the deployment, got diff %#v", diff)
			}
		})
	}
}
