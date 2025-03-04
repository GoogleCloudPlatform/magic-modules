/*
* Copyright 2025 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertTestNameToResource(t *testing.T) {
	cases := map[string]struct {
		testName string
		want     string
	}{
		"Resource Test: standard format": {
			testName: "TestAccCloudRunV2Job_cloudrunv2JobBasicExample",
			want:     "google_cloud_run_v2_job",
		},
		"Resource Test: standard format with ALL CAPS (GKE)": {
			testName: "TestAccGKEHubFeatureMembership_gkehubFeaturePolicyController",
			want:     "google_gke_hub_feature_membership",
		},
		"Resource Test: non-standard format": {
			testName: "TestAccSecurityPosturePostureDeployment_securityposturePostureDeployment_update",
			want:     "google_securityposture_posture_deployment",
		},
		"Data Source Test: standard format": {
			testName: "TestAccDataSourceGoogleCloudRunV2Job_basic",
			want:     "google_cloud_run_v2_job",
		},
		"Data Source Test: non-standard format": {
			testName: "TestAccDataSourceGoogleCloudBackupDRDataSource_basic",
			want:     "google_backup_dr_data_source",
		},
		"IAM Resource Test: standard format": {
			testName: "TestAccBeyondcorpSecurityGatewayIamMemberGenerated",
			want:     "google_beyondcorp_security_gateway_iam_member",
		},
		"IAM Resource Test: standard format with underscore": {
			testName: "TestAccBeyondcorpSecurityGatewayIamBindingGenerated_withCondition",
			want:     "google_beyondcorp_security_gateway_iam_binding",
		},
		"IAM Resource Test: non-standard format": {
			testName: "TestAccIAM3OrganizationsPolicyBinding_iam3OrganizationsPolicyBindingExample_update",
			want:     "google_iam_organizations_policy_binding",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := convertTestNameToResource(tc.testName)
			assert.Equal(t, tc.want, got)
		})
	}
}
