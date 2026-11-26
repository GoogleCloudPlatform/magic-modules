package vertexai_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/vertexai"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestVertexAIEndpointDeployedModelImport(t *testing.T) {
	r := vertexai.ResourceVertexAIEndpointDeployedModel()
	d := r.TestResourceData()
	d.SetId("projects/my-project/locations/us-central1/endpoints/endpoint-123/deployedModels/456")

	out, err := r.Importer.State(d, &transport_tpg.Config{})
	if err != nil {
		t.Fatalf("Importer.State returned error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected one imported resource, got %d", len(out))
	}

	got := out[0]
	if got.Get("endpoint").(string) != "projects/my-project/locations/us-central1/endpoints/endpoint-123" {
		t.Errorf("endpoint = %q", got.Get("endpoint"))
	}
	if got.Get("region").(string) != "us-central1" {
		t.Errorf("region = %q", got.Get("region"))
	}
	if got.Get("deployed_model_id").(string) != "456" {
		t.Errorf("deployed_model_id = %q", got.Get("deployed_model_id"))
	}
	if got.Id() != d.Id() {
		t.Errorf("id = %q, want %q", got.Id(), d.Id())
	}
}
