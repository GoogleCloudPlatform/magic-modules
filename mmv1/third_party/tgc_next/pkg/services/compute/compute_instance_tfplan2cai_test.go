// Copyright 2026 Google LLC
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
package compute

import (
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/models"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseComputeInstancePlanValues() map[string]interface{} {
	return map[string]interface{}{
		"name":         "test-instance",
		"machine_type": "e2-medium",
		"zone":         "us-central1-a",
		"boot_disk": []interface{}{
			map[string]interface{}{
				"auto_delete": true,
				"initialize_params": []interface{}{
					map[string]interface{}{
						"image": "debian-cloud/debian-13",
					},
				},
			},
		},
		"network_interface": []interface{}{
			map[string]interface{}{
				"network": "default",
			},
		},
	}
}

func TestExpandComputeInstance_workloadIdentityConfig(t *testing.T) {
	t.Parallel()

	resSchema := ResourceComputeInstance().Schema
	cfg := &transport_tpg.Config{
		Project: "my-project",
		Zone:    "us-central1-a",
	}

	t.Run("explicit false identity_certificate_enabled preserved in ForceSendFields", func(t *testing.T) {
		t.Parallel()
		values := baseComputeInstancePlanValues()
		values["workload_identity_config"] = []interface{}{
			map[string]interface{}{
				"identity":                     "ns/my-ns/sa/my-id",
				"identity_certificate_enabled": false,
			},
		}
		d := models.NewFakeResourceDataWithMeta("google_compute_instance", resSchema, values, false, "google_compute_instance.foobar")
		inst, err := expandComputeInstance("my-project", d, cfg)
		require.NoError(t, err)
		require.NotNil(t, inst.WorkloadIdentityConfig)
		assert.Equal(t, "ns/my-ns/sa/my-id", inst.WorkloadIdentityConfig.Identity)
		assert.False(t, inst.WorkloadIdentityConfig.IdentityCertificateEnabled)
		assert.Contains(t, inst.WorkloadIdentityConfig.ForceSendFields, "Identity")
		assert.Contains(t, inst.WorkloadIdentityConfig.ForceSendFields, "IdentityCertificateEnabled")
	})

	t.Run("omitted nil identity_certificate_enabled omitted from ForceSendFields", func(t *testing.T) {
		t.Parallel()
		values := baseComputeInstancePlanValues()
		values["workload_identity_config"] = []interface{}{
			map[string]interface{}{
				"identity":                     "ns/my-ns/sa/my-id",
				"identity_certificate_enabled": nil,
			},
		}
		d := models.NewFakeResourceDataWithMeta("google_compute_instance", resSchema, values, false, "google_compute_instance.foobar")
		inst, err := expandComputeInstance("my-project", d, cfg)
		require.NoError(t, err)
		require.NotNil(t, inst.WorkloadIdentityConfig)
		assert.Equal(t, "ns/my-ns/sa/my-id", inst.WorkloadIdentityConfig.Identity)
		assert.Contains(t, inst.WorkloadIdentityConfig.ForceSendFields, "Identity")
		assert.NotContains(t, inst.WorkloadIdentityConfig.ForceSendFields, "IdentityCertificateEnabled")
	})

	t.Run("empty or beta-only workload_identity_config leaves WorkloadIdentityConfig nil", func(t *testing.T) {
		t.Parallel()
		values := baseComputeInstancePlanValues()
		values["workload_identity_config"] = []interface{}{
			map[string]interface{}{
				"identity":                     nil,
				"identity_certificate_enabled": nil,
			},
		}
		d := models.NewFakeResourceDataWithMeta("google_compute_instance", resSchema, values, false, "google_compute_instance.foobar")
		inst, err := expandComputeInstance("my-project", d, cfg)
		require.NoError(t, err)
		assert.Nil(t, inst.WorkloadIdentityConfig)
	})
}
