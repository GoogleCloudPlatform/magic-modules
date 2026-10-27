package vertexai

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const SharedVertexAIModelDisplayName = "tf-bootstrap-vertex-ai-model"

// BootstrapVertexAIModel returns the resource name of a shared Model in the
// test project, uploading it first if it doesn't exist yet. Tests use it as a
// copy source, so it must not be deleted.
func BootstrapVertexAIModel(t *testing.T, region string) string {
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	project := envvar.GetTestProjectFromEnv()
	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)
	apiURL := strings.ReplaceAll(transport_tpg.BaseUrl(Product, config), "{{region}}", region)
	baseURL := apiURL + parent + "/models"

	filter := url.QueryEscape(fmt.Sprintf("display_name=%q", SharedVertexAIModelDisplayName))
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    baseURL + "?filter=" + filter,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		t.Fatalf("Error listing Vertex AI models: %s", err)
	}
	if models, ok := res["models"].([]interface{}); ok && len(models) > 0 {
		if model, ok := models[0].(map[string]interface{}); ok {
			if name, ok := model["name"].(string); ok && name != "" {
				return name
			}
		}
	}

	t.Logf("Uploading shared Vertex AI model %q in %s", SharedVertexAIModelDisplayName, parent)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    baseURL + ":upload",
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"model": map[string]interface{}{
				"displayName": SharedVertexAIModelDisplayName,
				"containerSpec": map[string]interface{}{
					"imageUri": "us-docker.pkg.dev/cloudrun/container/hello:latest",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Error uploading Vertex AI model: %s", err)
	}

	opName, ok := op["name"].(string)
	if !ok || opName == "" {
		t.Fatalf("Vertex AI model upload did not return an operation: %#v", op)
	}
	var opRes map[string]interface{}
	for deadline := time.Now().Add(20 * time.Minute); ; time.Sleep(10 * time.Second) {
		if time.Now().After(deadline) {
			t.Fatalf("Timed out waiting for Vertex AI model upload operation %s", opName)
		}
		op, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    apiURL + opName,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			t.Fatalf("Error polling Vertex AI model upload operation %s: %s", opName, err)
		}
		if done, _ := op["done"].(bool); !done {
			continue
		}
		if opErr, ok := op["error"]; ok {
			t.Fatalf("Vertex AI model upload failed: %#v", opErr)
		}
		opRes, _ = op["response"].(map[string]interface{})
		break
	}
	name, ok := opRes["model"].(string)
	if !ok || name == "" {
		t.Fatalf("Vertex AI model upload response did not include a model name: %#v", opRes)
	}
	return name
}
