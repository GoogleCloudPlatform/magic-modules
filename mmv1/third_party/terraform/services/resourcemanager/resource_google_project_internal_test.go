package resourcemanager

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestPopulateGoogleProjectResourceData_DoesNotCopyAllLabelsToTerraformLabelsWhenUnset(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceGoogleProject().Schema, map[string]interface{}{
		"project_id": "example-project",
		"name":       "Example Project",
	})

	project := &cloudresourcemanager.Project{
		ProjectNumber: 123456789,
		Name:          "Example Project",
		Labels: map[string]string{
			"firebase":     "enabled",
			"earth-engine": "",
		},
	}

	if err := populateGoogleProjectResourceData(d, project, "example-project", &transport_tpg.Config{}); err != nil {
		t.Fatalf("populateGoogleProjectResourceData() returned error: %v", err)
	}

	if got := d.Get("labels").(map[string]interface{}); len(got) != 0 {
		t.Fatalf("expected labels to remain empty when unset in config, got %#v", got)
	}

	if got := d.Get("terraform_labels").(map[string]interface{}); len(got) != 0 {
		t.Fatalf("expected terraform_labels to remain empty when unset in config, got %#v", got)
	}

	if got := d.Get("effective_labels").(map[string]interface{}); !reflect.DeepEqual(got, map[string]interface{}{
		"firebase":     "enabled",
		"earth-engine": "",
	}) {
		t.Fatalf("expected effective_labels to contain all project labels, got %#v", got)
	}
}
