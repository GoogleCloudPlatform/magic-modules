package resourcemanager

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

func TestRetryProjectDefaultNetworkDeletion_RetriesApiNotEnabled(t *testing.T) {
	tests := map[string]error{
		"typed":      &googleapi.Error{Code: 403, Errors: []googleapi.ErrorItem{{Reason: "accessNotConfigured"}}},
		"plain text": fmt.Errorf("googleapi: Error 403: Compute Engine API has not been used in project 123456789 before or it is disabled."),
		"operation polling": fmt.Errorf("Error waiting for Deleting Firewall: error while retrieving operation: %s", &googleapi.Error{
			Code:    403,
			Message: "Compute Engine API has not been used in project 123456789 before or it is disabled.",
		}),
	}
	for name, apiErr := range tests {
		t.Run(name, func(t *testing.T) {
			attempts := 0
			err := retryProjectDefaultNetworkDeletion(func() error {
				attempts++
				if attempts < 3 {
					return apiErr
				}
				return nil
			}, 5*time.Second)

			if err != nil {
				t.Fatalf("retryProjectDefaultNetworkDeletion() returned error: %v", err)
			}
			if attempts != 3 {
				t.Fatalf("expected 3 deletion attempts, got %d", attempts)
			}
		})
	}
}

func TestRetryProjectDefaultNetworkDeletion_DoesNotRetryOther403(t *testing.T) {
	tests := map[string]error{
		"typed":      &googleapi.Error{Code: 403},
		"plain text": fmt.Errorf("googleapi: Error 403: Required 'compute.networks.delete' permission"),
		"operation polling": fmt.Errorf("Error waiting for Deleting Firewall: error while retrieving operation: %s", &googleapi.Error{
			Code:    403,
			Message: "Required 'compute.globalOperations.get' permission",
		}),
	}
	for name, apiErr := range tests {
		t.Run(name, func(t *testing.T) {
			attempts := 0
			err := retryProjectDefaultNetworkDeletion(func() error {
				attempts++
				return apiErr
			}, time.Second)

			if err == nil {
				t.Fatal("retryProjectDefaultNetworkDeletion() returned nil, want error")
			}
			if attempts != 1 {
				t.Fatalf("expected 1 deletion attempt, got %d", attempts)
			}
		})
	}
}

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
