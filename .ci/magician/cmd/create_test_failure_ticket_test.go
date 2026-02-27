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
	"testing"

	"magician/provider"

	"github.com/stretchr/testify/assert"
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

func TestIsTerraformTeamOwned(t *testing.T) {
	cases := map[string]struct {
		tf   testFailure
		want bool
	}{
		"GA Quota Error": {
			tf: testFailure{
				ErrorTypes: map[provider.Version]string{
					provider.GA: "Quota",
				},
			},
			want: true,
		},
		"Beta API Enablement Error": {
			tf: testFailure{
				ErrorTypes: map[provider.Version]string{
					provider.Beta: "API enablement (Test environment)",
				},
			},
			want: true,
		},
		"GA Other Error": {
			tf: testFailure{
				ErrorTypes: map[provider.Version]string{
					provider.GA: "Some other error",
				},
			},
			want: false,
		},
		"Beta Other Error": {
			tf: testFailure{
				ErrorTypes: map[provider.Version]string{
					provider.Beta: "Another different error",
				},
			},
			want: false,
		},
		"Mixed Errors - One TF Owned": {
			tf: testFailure{
				ErrorTypes: map[provider.Version]string{
					provider.GA:   "Some other error",
					provider.Beta: "Quota",
				},
			},
			want: true,
		},
		"No Error Types": {
			tf: testFailure{
				ErrorTypes: make(map[provider.Version]string),
			},
			want: false,
		},
		"Nil Error Types": {
			tf:   testFailure{},
			want: false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := IsTerraformTeamOwned(&tc.tf)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestShouldCreateTicket(t *testing.T) {
	cases := map[string]struct {
		tf                   testFailure
		existTestNames       []string
		todayClosedTestNames []string
		want                 bool
	}{
		"No failures": {
			tf: testFailure{
				TestName: "TestAccSomething",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA:   testFailureNone,
					provider.Beta: testFailureNone,
				},
			},
			want: false,
		},
		"Already exists": {
			tf: testFailure{
				TestName: "TestAccExisting",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure100,
				},
			},
			existTestNames: []string{"TestAccExisting"},
			want:           false,
		},
		"Closed today": {
			tf: testFailure{
				TestName: "TestAccClosed",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure100,
				},
			},
			todayClosedTestNames: []string{"TestAccClosed"},
			want:                 false,
		},
		"Team owned error - create": {
			tf: testFailure{
				TestName:   "TestAccTeamOwned",
				ErrorTypes: map[provider.Version]string{provider.GA: "Quota"},
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA:   testFailure10, // Normally wouldn't create
					provider.Beta: testFailureNone,
				},
			},
			want: true,
		},
		"GA 50% failure - create": {
			tf: testFailure{
				TestName: "TestAccGA50",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure50,
				},
			},
			want: true,
		},
		"Beta 100% failure - create": {
			tf: testFailure{
				TestName: "TestAccBeta100",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.Beta: testFailure100,
				},
			},
			want: true,
		},
		"Low failure - don't create": {
			tf: testFailure{
				TestName: "TestAccLowFailure",
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA:   testFailure10,
					provider.Beta: testFailure10,
				},
			},
			want: false,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got := shouldCreateTicket(&tc.tf, tc.existTestNames, tc.todayClosedTestNames)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetTicketLabels(t *testing.T) {
	cases := map[string]struct {
		tf           testFailure
		expectLabels []string
	}{
		"Team owned error": {
			tf: testFailure{
				TestName:         "TestAccTeamOwnedQuota",
				AffectedResource: "google_compute_instance",
				ErrorTypes:       map[provider.Version]string{provider.GA: "Quota"},
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure10,
				},
			},
			expectLabels: []string{"size/xs", "test-failure", "test-failure-10", "service/terraform"},
		},
		"Non-Team owned - has service label": {
			tf: testFailure{
				TestName:         "TestAccComputeInstance",
				AffectedResource: "google_compute_instance",
				ErrorTypes:       map[provider.Version]string{provider.GA: "API Error"},
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure50,
				},
			},
			expectLabels: []string{"size/xs", "test-failure", "test-failure-50", "service/compute-instances"},
		},
		"Non-Team owned - fallback label": {
			tf: testFailure{
				TestName:         "TestAccUnknownResource",
				AffectedResource: "google_unknown_resource",
				ErrorTypes:       map[provider.Version]string{provider.GA: "API Error"},
				FailureRateLabels: map[provider.Version]testFailureRateLabel{
					provider.GA: testFailure100,
				},
			},
			expectLabels: []string{"size/xs", "test-failure", "test-failure-100", "service/terraform"},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			labels, err := computeTicketLabels(&tc.tf)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expectLabels, labels)
		})
	}
}
