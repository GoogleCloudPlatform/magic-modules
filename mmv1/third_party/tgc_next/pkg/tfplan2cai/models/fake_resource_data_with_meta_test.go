// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package models

import (
	"testing"

	provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
	"github.com/stretchr/testify/assert"
)

func TestFakeResourceDataWithMeta_kind(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	assert.Equal(t, "google_compute_disk", d.Kind())
}

func TestFakeResourceDataWithMeta_id(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	assert.Equal(t, d.Id(), "")
}

func TestFakeResourceDataWithMeta_get(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	assert.Equal(t, d.Get("name"), "test-disk")
}

func TestFakeResourceDataWithMeta_getOkOk(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	res, ok := d.GetOk("name")
	assert.Equal(t, "test-disk", res)
	assert.True(t, ok)
}

func TestFakeResourceDataWithMeta_getOkNonexistentField(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	res, ok := d.GetOk("incorrect")
	assert.Nil(t, res)
	assert.False(t, ok)
}

func TestFakeResourceDataWithMeta_getOkEmptyString(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		false,
		"google_compute_disk.test-disk",
	)
	res, ok := d.GetOk("image")
	assert.Equal(t, "", res)
	assert.False(t, ok)
}

func TestFakeResourceDataWithMeta_getOkUnsetString(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":     "my-node-pool",
		"location": "us-central1",
		"cluster":  "projects/my-project-id/global/clusters/my-gke-cluster",
		"config": map[string]interface{}{
			"machineType": "n1-standard-1",
			"metadata": map[string]string{
				"disable-legacy-endpoints": "true",
			},
			"oauthScopes": []string{
				"https://www.googleapis.com/auth/cloud-platform",
			},
			"preemptible": true,
		},
	}
	d := NewFakeResourceDataWithMeta(
		"google_container_cluster",
		p.ResourcesMap["google_container_cluster"].Schema,
		values,
		false,
		"google_container_cluster.my-node-pool",
	)
	res, ok := d.GetOk("subnetwork")
	assert.Equal(t, "", res)
	assert.False(t, ok)
}

func TestFakeResourceDataWithMeta_getOkTypeObject(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"advanced_machine_features": []interface{}{},
		"allow_stopping_for_update": nil,
		"attached_disk": []interface{}{
			map[string]interface{}{
				"device_name":                     "test-device_name",
				"disk_encryption_key_raw":         nil,
				"disk_encryption_key_rsa":         nil,
				"kms_key_self_link":               "test-kms_key_self_link",
				"disk_encryption_service_account": nil,
				"mode":                            "READ_ONLY",
				"source":                          "test-source",
			},
			map[string]interface{}{
				"disk_encryption_key_raw": nil,
				"mode":                    "READ_WRITE",
				"source":                  "test-source2",
			},
		},
		"boot_disk": []interface{}{
			map[string]interface{}{
				"auto_delete":             true,
				"disk_encryption_key_raw": nil,
				"initialize_params": []interface{}{
					map[string]interface{}{
						"image": "debian-cloud/debian-9",
					},
				},
				"mode": "READ_WRITE",
			},
		},
		"can_ip_forward":          false,
		"deletion_protection":     false,
		"description":             nil,
		"desired_status":          nil,
		"enable_display":          nil,
		"hostname":                nil,
		"labels":                  nil,
		"machine_type":            "n1-standard-1",
		"metadata":                nil,
		"metadata_startup_script": nil,
		"name":                    "test",
		"network_interface": []interface{}{
			map[string]interface{}{
				"access_config": []interface{}{
					map[string]interface{}{
						"public_ptr_domain_name": nil,
					},
				},
				"alias_ip_range":     []interface{}{},
				"ipv6_access_config": []interface{}{},
				"network":            "default",
				"nic_type":           nil,
			},
		},
		"resource_policies": nil,
		"scratch_disk": []interface{}{
			map[string]interface{}{
				"interface": "SCSI",
			},
		},
		"service_account":          []interface{}{},
		"shielded_instance_config": []interface{}{},
		"tags": []interface{}{
			"bar",
			"foo",
		},
		"timeouts": nil,
		"zone":     "us-central1-a",
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_instance",
		p.ResourcesMap["google_compute_instance"].Schema,
		values,
		false,
		"google_compute_instance.test",
	)
	res, ok := d.GetOk("attached_disk.0")
	assert.Equal(t, map[string]interface{}{
		"device_name":                     "test-device_name",
		"disk_encryption_key_raw":         "",
		"disk_encryption_key_sha256":      "",
		"disk_encryption_key_rsa":         "",
		"disk_encryption_service_account": "",
		"kms_key_self_link":               "test-kms_key_self_link",
		"mode":                            "READ_ONLY",
		"source":                          "test-source",
	}, res)
	assert.True(t, ok)
}

func TestFakeResourceDataWithMeta_getOknsetTypeObject(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"advanced_machine_features": []interface{}{},
		"allow_stopping_for_update": nil,
		"attached_disk":             []interface{}{},
		"boot_disk": []interface{}{
			map[string]interface{}{
				"auto_delete":             true,
				"disk_encryption_key_raw": nil,
				"initialize_params": []interface{}{
					map[string]interface{}{
						"image": "debian-cloud/debian-9",
					},
				},
				"mode": "READ_WRITE",
			},
		},
		"can_ip_forward":          false,
		"deletion_protection":     false,
		"description":             nil,
		"desired_status":          nil,
		"enable_display":          nil,
		"hostname":                nil,
		"labels":                  nil,
		"machine_type":            "n1-standard-1",
		"metadata":                nil,
		"metadata_startup_script": nil,
		"name":                    "test",
		"network_interface": []interface{}{
			map[string]interface{}{
				"access_config": []interface{}{
					map[string]interface{}{
						"public_ptr_domain_name": nil,
					},
				},
				"alias_ip_range":     []interface{}{},
				"ipv6_access_config": []interface{}{},
				"network":            "default",
				"nic_type":           nil,
			},
		},
		"resource_policies": nil,
		"scratch_disk": []interface{}{
			map[string]interface{}{
				"interface": "SCSI",
			},
		},
		"service_account":          []interface{}{},
		"shielded_instance_config": []interface{}{},
		"tags": []interface{}{
			"bar",
			"foo",
		},
		"timeouts": nil,
		"zone":     "us-central1-a",
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_instance",
		p.ResourcesMap["google_compute_instance"].Schema,
		values,
		false,
		"google_compute_instance.test",
	)
	res, ok := d.GetOk("attached_disk.0")
	assert.Equal(t, map[string]interface{}{
		"device_name":                     "",
		"disk_encryption_key_raw":         "",
		"disk_encryption_key_sha256":      "",
		"disk_encryption_key_rsa":         "",
		"disk_encryption_service_account": "",
		"kms_key_self_link":               "",
		"mode":                            "",
		"source":                          "",
	}, res)
	assert.False(t, ok)
}

func TestFakeResourceDataWithMeta_isDelelted(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		true,
		"google_compute_disk.test-disk",
	)
	assert.Equal(t, true, d.IsDeleted())
}

func TestFakeResourceDataWithMeta_address(t *testing.T) {
	p := provider.Provider()

	values := map[string]interface{}{
		"name":                      "test-disk",
		"type":                      "pd-ssd",
		"zone":                      "us-central1-a",
		"image":                     "projects/debian-cloud/global/images/debian-8-jessie-v20170523",
		"physical_block_size_bytes": 4096,
	}
	d := NewFakeResourceDataWithMeta(
		"google_compute_disk",
		p.ResourcesMap["google_compute_disk"].Schema,
		values,
		true,
		"google_compute_disk.test-disk",
	)
	assert.Equal(t, "google_compute_disk.test-disk", d.Address())
}
