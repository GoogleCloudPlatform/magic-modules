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

func TestConvertComputeInstanceWithLabelsAreNewResourcesTrue(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//compute.googleapis.com/projects/example-project/zones/us-central1-a/instances/example-instance",
			Type: "compute.googleapis.com/Instance",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data: map[string]interface{}{
					"name":        "example-instance",
					"machineType": "zones/us-central1-a/machineTypes/n1-standard-1",
					"disks": []interface{}{
						map[string]interface{}{
							"boot":   true,
							"source": "projects/example-project/zones/us-central1-a/disks/example-instance",
						},
					},
					"networkInterfaces": []interface{}{
						map[string]interface{}{
							"network": "projects/example-project/global/networks/default",
						},
					},
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

	expected := `resource "google_compute_instance" "example-instance" {
  boot_disk {
    source = "projects/example-project/zones/us-central1-a/disks/example-instance"
  }
  labels = {
    key1 = "value1"
  }
  machine_type = "n1-standard-1"
  name         = "example-instance"
  network_interface {
    network = "projects/example-project/global/networks/default"
  }
  project = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertComputeInstanceWithLabelsAreNewResourcesFalse(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//compute.googleapis.com/projects/example-project/zones/us-central1-a/instances/example-instance",
			Type: "compute.googleapis.com/Instance",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Instance",
				Data: map[string]interface{}{
					"name":        "example-instance",
					"machineType": "zones/us-central1-a/machineTypes/n1-standard-1",
					"disks": []interface{}{
						map[string]interface{}{
							"boot":   true,
							"source": "projects/example-project/zones/us-central1-a/disks/example-instance",
						},
					},
					"networkInterfaces": []interface{}{
						map[string]interface{}{
							"network": "projects/example-project/global/networks/default",
						},
					},
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

	expected := `resource "google_compute_instance" "example-instance" {
  boot_disk {
    source = "projects/example-project/zones/us-central1-a/disks/example-instance"
  }
  machine_type = "n1-standard-1"
  name         = "example-instance"
  network_interface {
    network = "projects/example-project/global/networks/default"
  }
  project = "example-project"
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertContainerClusterWithLabelsAreNewResourcesTrue(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//container.googleapis.com/projects/example-project/locations/us-central1-a/clusters/example-cluster",
			Type: "container.googleapis.com/Cluster",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
				DiscoveryName:        "Cluster",
				Data: map[string]interface{}{
					"name":     "example-cluster",
					"location": "us-central1-a",
					"resourceLabels": map[string]interface{}{
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

	expected := `resource "google_container_cluster" "example-cluster" {
  location        = "us-central1-a"
  name            = "example-cluster"
  networking_mode = "ROUTES"
  project         = "example-project"
  release_channel {
    channel = "UNSPECIFIED"
  }
  resource_labels = {
    key1 = "value1"
  }
  secret_manager_config {
    enabled = false
  }
  secret_sync_config {
    enabled = false
  }
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}

func TestConvertContainerClusterWithLabelsAreNewResourcesFalse(t *testing.T) {
	assets := []caiasset.Asset{
		{
			Name: "//container.googleapis.com/projects/example-project/locations/us-central1-a/clusters/example-cluster",
			Type: "container.googleapis.com/Cluster",
			Resource: &caiasset.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/container/v1/rest",
				DiscoveryName:        "Cluster",
				Data: map[string]interface{}{
					"name":     "example-cluster",
					"location": "us-central1-a",
					"resourceLabels": map[string]interface{}{
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

	expected := `resource "google_container_cluster" "example-cluster" {
  location        = "us-central1-a"
  name            = "example-cluster"
  networking_mode = "ROUTES"
  project         = "example-project"
  release_channel {
    channel = "UNSPECIFIED"
  }
  secret_manager_config {
    enabled = false
  }
  secret_sync_config {
    enabled = false
  }
}
`
	if string(got) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(got))
	}
}
