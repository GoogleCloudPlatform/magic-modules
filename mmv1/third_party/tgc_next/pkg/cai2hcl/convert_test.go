package cai2hcl

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
)

func TestConvertWithResourceName(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Parent:               "//cloudresourcemanager.googleapis.com/folders/456",
				Data: map[string]interface{}{
					"name":      "My Project",
					"projectId": "example-project",
				},
			},
		},
	}

	got, err := converters.ConvertResource(assets, &models.ResourceConverterOptions{
		ResourceName: "custom_project_name",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := `resource "google_project" "custom_project_name" {
  folder_id  = "456"
  name       = "My Project"
  project_id = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertWithoutResourceName(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Parent:               "//cloudresourcemanager.googleapis.com/folders/456",
				Data: map[string]interface{}{
					"name":      "My Project",
					"projectId": "example-project",
				},
			},
		},
	}

	got, err := converters.ConvertResource(assets, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := `resource "google_project" "example-project" {
  folder_id  = "456"
  name       = "My Project"
  project_id = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertProjectWithLabelsAreNewResourcesTrue(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Parent:               "//cloudresourcemanager.googleapis.com/folders/456",
				Data: map[string]interface{}{
					"name":      "My Project",
					"projectId": "example-project",
					"labels": map[string]interface{}{
						"key1": "value1",
					},
				},
			},
		},
	}

	got, err := converters.ConvertResource(assets, &models.ResourceConverterOptions{
		AreNewResources: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := `resource "google_project" "example-project" {
  folder_id = "456"
  labels = {
    key1 = "value1"
  }
  name       = "My Project"
  project_id = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertProjectWithLabelsAreNewResourcesFalse(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//cloudresourcemanager.googleapis.com/projects/example-project",
			Type: "cloudresourcemanager.googleapis.com/Project",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Project",
				Parent:               "//cloudresourcemanager.googleapis.com/folders/456",
				Data: map[string]interface{}{
					"name":      "My Project",
					"projectId": "example-project",
					"labels": map[string]interface{}{
						"key1": "value1",
					},
				},
			},
		},
	}

	got, err := converters.ConvertResource(assets, &models.ResourceConverterOptions{
		AreNewResources: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := `resource "google_project" "example-project" {
  folder_id  = "456"
  name       = "My Project"
  project_id = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}
