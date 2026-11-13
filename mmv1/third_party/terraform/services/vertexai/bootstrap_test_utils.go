package vertexai

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const sharedVertexAIModelID = "tf-bootstrap-vertex-ai-model"

// BootstrapVertexAIModel returns a shared Model in the test project, uploading
// it if it does not exist. Tests use it as a copy source, so it must not be deleted.
func BootstrapVertexAIModel(t *testing.T, region string) string {
	config := transport_tpg.BootstrapConfig(t)
	if config == nil {
		return ""
	}

	project := envvar.GetTestProjectFromEnv()
	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)
	apiURL := strings.ReplaceAll(transport_tpg.BaseUrl(Product, config), "{{region}}", region)
	modelName := fmt.Sprintf("%s/models/%s", parent, sharedVertexAIModelID)

	exists, err := vertexAIModelExists(config, project, apiURL, modelName)
	if err != nil {
		t.Fatalf("Error checking for shared Vertex AI model %q: %s", modelName, err)
	}
	if exists {
		return modelName
	}

	t.Logf("Uploading shared Vertex AI model %q", modelName)
	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    apiURL + parent + "/models:upload",
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"modelId": sharedVertexAIModelID,
			"model": map[string]interface{}{
				"displayName": sharedVertexAIModelID,
				"artifactUri": "gs://cloud-samples-data/vertex-ai/model-deployment/models/resnet",
				"containerSpec": map[string]interface{}{
					"imageUri": "us-docker.pkg.dev/vertex-ai/prediction/tf2-cpu.2-12:latest",
				},
			},
		},
	})
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
			t.Fatalf("Error uploading shared Vertex AI model %q: %s", modelName, err)
		}
		return waitForVertexAIModel(t, config, project, apiURL, modelName)
	}

	opName, ok := op["name"].(string)
	if !ok || opName == "" {
		t.Fatalf("Vertex AI model upload did not return an operation: %#v", op)
	}
	for deadline := time.Now().Add(20 * time.Minute); ; time.Sleep(10 * time.Second) {
		if time.Now().After(deadline) {
			t.Fatalf("Timed out waiting for Vertex AI model upload operation %q", opName)
		}
		op, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    apiURL + opName,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			t.Fatalf("Error polling Vertex AI model upload operation %q: %s", opName, err)
		}
		if done, _ := op["done"].(bool); !done {
			continue
		}
		if opErr, ok := op["error"]; ok {
			t.Fatalf("Vertex AI model upload failed: %#v", opErr)
		}
		if response, ok := op["response"].(map[string]interface{}); ok {
			if name, ok := response["model"].(string); ok && name != "" {
				return name
			}
		}
		t.Fatalf("Vertex AI model upload response did not include a model name: %#v", op)
	}
}

func waitForVertexAIModel(t *testing.T, config *transport_tpg.Config, project, apiURL, modelName string) string {
	for deadline := time.Now().Add(5 * time.Minute); ; time.Sleep(5 * time.Second) {
		if time.Now().After(deadline) {
			t.Fatalf("Timed out waiting for shared Vertex AI model %q to become available", modelName)
		}
		exists, err := vertexAIModelExists(config, project, apiURL, modelName)
		if err != nil {
			t.Fatalf("Error checking for shared Vertex AI model %q: %s", modelName, err)
		}
		if exists {
			return modelName
		}
	}
}

func vertexAIModelExists(config *transport_tpg.Config, project, apiURL, modelName string) (bool, error) {
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    apiURL + modelName,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
