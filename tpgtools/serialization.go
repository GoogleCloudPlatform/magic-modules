// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	assuredworkloadsBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads/beta"
	cloudbuild "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild"
	cloudbuildBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild/beta"
	cloudresourcemanager "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager"
	cloudresourcemanagerBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager/beta"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	computeBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute/beta"
	containeraws "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws"
	containerawsBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containeraws/beta"
	containerazure "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure"
	containerazureBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/containerazure/beta"
	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"
	dataprocBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc/beta"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	eventarcBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc/beta"
	gkehubBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/gkehub/beta"
	monitoringBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring/beta"
	networkconnectivity "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity"
	networkconnectivityBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/networkconnectivity/beta"
	orgpolicy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy"
	orgpolicyBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy/beta"
	osconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig"
	osconfigBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/osconfig/beta"
	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"
	privatecaBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca/beta"
	recaptchaenterprise "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/recaptchaenterprise"
	recaptchaenterpriseBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/recaptchaenterprise/beta"
	fmtcmd "github.com/hashicorp/hcl/hcl/fmtcmd"
)

// DCLToTerraformReference converts a DCL resource name to the final tpgtools name
// after overrides are applied
func DCLToTerraformReference(product DCLPackageName, resource miscellaneousNameSnakeCase, version string) (string, error) {
	if version == "alpha" {
		switch fmt.Sprintf("%s/%s", product, resource) {
		}
	}
	if version == "beta" {
		switch fmt.Sprintf("%s/%s", product, resource) {
		case "assuredworkloads/workload":
			return "google_assured_workloads_workload", nil
		case "cloudbuild/worker_pool":
			return "google_cloudbuild_worker_pool", nil
		case "cloudresourcemanager/folder":
			return "google_folder", nil
		case "cloudresourcemanager/project":
			return "google_project", nil
		case "compute/firewall_policy":
			return "google_compute_firewall_policy", nil
		case "compute/firewall_policy_association":
			return "google_compute_firewall_policy_association", nil
		case "compute/firewall_policy_rule":
			return "google_compute_firewall_policy_rule", nil
		case "compute/forwarding_rule":
			return "google_compute_forwarding_rule", nil
		case "compute/global_forwarding_rule":
			return "google_compute_global_forwarding_rule", nil
		case "containeraws/cluster":
			return "google_container_aws_cluster", nil
		case "containeraws/node_pool":
			return "google_container_aws_node_pool", nil
		case "containerazure/client":
			return "google_container_azure_client", nil
		case "containerazure/cluster":
			return "google_container_azure_cluster", nil
		case "containerazure/node_pool":
			return "google_container_azure_node_pool", nil
		case "dataproc/workflow_template":
			return "google_dataproc_workflow_template", nil
		case "eventarc/trigger":
			return "google_eventarc_trigger", nil
		case "gkehub/feature":
			return "google_gke_hub_feature", nil
		case "gkehub/feature_membership":
			return "google_gke_hub_feature_membership", nil
		case "monitoring/monitored_project":
			return "google_monitoring_monitored_project", nil
		case "networkconnectivity/hub":
			return "google_network_connectivity_hub", nil
		case "networkconnectivity/spoke":
			return "google_network_connectivity_spoke", nil
		case "orgpolicy/policy":
			return "google_org_policy_policy", nil
		case "osconfig/os_policy_assignment":
			return "google_os_config_os_policy_assignment", nil
		case "privateca/certificate_template":
			return "google_privateca_certificate_template", nil
		case "recaptchaenterprise/key":
			return "google_recaptcha_enterprise_key", nil
		}
	}
	// If not found in sample version, fallthrough to GA
	switch fmt.Sprintf("%s/%s", product, resource) {
	case "assuredworkloads/workload":
		return "google_assured_workloads_workload", nil
	case "cloudbuild/worker_pool":
		return "google_cloudbuild_worker_pool", nil
	case "cloudresourcemanager/folder":
		return "google_folder", nil
	case "cloudresourcemanager/project":
		return "google_project", nil
	case "compute/firewall_policy":
		return "google_compute_firewall_policy", nil
	case "compute/firewall_policy_association":
		return "google_compute_firewall_policy_association", nil
	case "compute/firewall_policy_rule":
		return "google_compute_firewall_policy_rule", nil
	case "compute/forwarding_rule":
		return "google_compute_forwarding_rule", nil
	case "compute/global_forwarding_rule":
		return "google_compute_global_forwarding_rule", nil
	case "containeraws/cluster":
		return "google_container_aws_cluster", nil
	case "containeraws/node_pool":
		return "google_container_aws_node_pool", nil
	case "containerazure/client":
		return "google_container_azure_client", nil
	case "containerazure/cluster":
		return "google_container_azure_cluster", nil
	case "containerazure/node_pool":
		return "google_container_azure_node_pool", nil
	case "dataproc/workflow_template":
		return "google_dataproc_workflow_template", nil
	case "eventarc/trigger":
		return "google_eventarc_trigger", nil
	case "networkconnectivity/hub":
		return "google_network_connectivity_hub", nil
	case "networkconnectivity/spoke":
		return "google_network_connectivity_spoke", nil
	case "orgpolicy/policy":
		return "google_org_policy_policy", nil
	case "osconfig/os_policy_assignment":
		return "google_os_config_os_policy_assignment", nil
	case "privateca/certificate_template":
		return "google_privateca_certificate_template", nil
	case "recaptchaenterprise/key":
		return "google_recaptcha_enterprise_key", nil
	default:
		return "", fmt.Errorf("Error retrieving Terraform name from DCL resource type: %s/%s not found", product, resource)
	}

}

// ConvertSampleJSONToHCL unmarshals json to an HCL string.
func ConvertSampleJSONToHCL(product DCLPackageName, resource miscellaneousNameSnakeCase, version string, hasGAEquivalent bool, b []byte) (string, error) {
	if version == "alpha" {
		switch fmt.Sprintf("%s/%s", product, resource) {
		}
	}
	if version == "beta" {
		switch fmt.Sprintf("%s/%s", product, resource) {
		case "assuredworkloads/workload":
			r := &assuredworkloadsBeta.Workload{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return AssuredWorkloadsWorkloadBetaAsHCL(*r, hasGAEquivalent)
		case "cloudbuild/worker_pool":
			r := &cloudbuildBeta.WorkerPool{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return CloudbuildWorkerPoolBetaAsHCL(*r, hasGAEquivalent)
		case "cloudresourcemanager/folder":
			r := &cloudresourcemanagerBeta.Folder{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return CloudResourceManagerFolderBetaAsHCL(*r, hasGAEquivalent)
		case "cloudresourcemanager/project":
			r := &cloudresourcemanagerBeta.Project{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return serializeBetaProjectToHCL(*r, hasGAEquivalent)
		case "compute/firewall_policy":
			r := &computeBeta.FirewallPolicy{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyBetaAsHCL(*r, hasGAEquivalent)
		case "compute/firewall_policy_association":
			r := &computeBeta.FirewallPolicyAssociation{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyAssociationBetaAsHCL(*r, hasGAEquivalent)
		case "compute/firewall_policy_rule":
			r := &computeBeta.FirewallPolicyRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyRuleBetaAsHCL(*r, hasGAEquivalent)
		case "compute/forwarding_rule":
			r := &computeBeta.ForwardingRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeForwardingRuleBetaAsHCL(*r, hasGAEquivalent)
		case "compute/global_forwarding_rule":
			r := &computeBeta.ForwardingRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeGlobalForwardingRuleBetaAsHCL(*r, hasGAEquivalent)
		case "containeraws/cluster":
			r := &containerawsBeta.Cluster{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ContainerAwsClusterBetaAsHCL(*r, hasGAEquivalent)
		case "containeraws/node_pool":
			r := &containerawsBeta.NodePool{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ContainerAwsNodePoolBetaAsHCL(*r, hasGAEquivalent)
		case "containerazure/client":
			r := &containerazureBeta.AzureClient{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ContainerAzureClientBetaAsHCL(*r, hasGAEquivalent)
		case "containerazure/cluster":
			r := &containerazureBeta.Cluster{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ContainerAzureClusterBetaAsHCL(*r, hasGAEquivalent)
		case "containerazure/node_pool":
			r := &containerazureBeta.NodePool{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ContainerAzureNodePoolBetaAsHCL(*r, hasGAEquivalent)
		case "dataproc/workflow_template":
			r := &dataprocBeta.WorkflowTemplate{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return DataprocWorkflowTemplateBetaAsHCL(*r, hasGAEquivalent)
		case "eventarc/trigger":
			r := &eventarcBeta.Trigger{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return EventarcTriggerBetaAsHCL(*r, hasGAEquivalent)
		case "gkehub/feature":
			r := &gkehubBeta.Feature{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return GkeHubFeatureBetaAsHCL(*r, hasGAEquivalent)
		case "gkehub/feature_membership":
			r := &gkehubBeta.FeatureMembership{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return GkeHubFeatureMembershipBetaAsHCL(*r, hasGAEquivalent)
		case "monitoring/monitored_project":
			r := &monitoringBeta.MonitoredProject{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return MonitoringMonitoredProjectBetaAsHCL(*r, hasGAEquivalent)
		case "networkconnectivity/hub":
			r := &networkconnectivityBeta.Hub{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return NetworkConnectivityHubBetaAsHCL(*r, hasGAEquivalent)
		case "networkconnectivity/spoke":
			r := &networkconnectivityBeta.Spoke{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return NetworkConnectivitySpokeBetaAsHCL(*r, hasGAEquivalent)
		case "orgpolicy/policy":
			r := &orgpolicyBeta.Policy{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return OrgPolicyPolicyBetaAsHCL(*r, hasGAEquivalent)
		case "osconfig/os_policy_assignment":
			r := &osconfigBeta.OSPolicyAssignment{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return OsConfigOsPolicyAssignmentBetaAsHCL(*r, hasGAEquivalent)
		case "privateca/certificate_template":
			r := &privatecaBeta.CertificateTemplate{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return PrivatecaCertificateTemplateBetaAsHCL(*r, hasGAEquivalent)
		case "recaptchaenterprise/key":
			r := &recaptchaenterpriseBeta.Key{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return RecaptchaEnterpriseKeyBetaAsHCL(*r, hasGAEquivalent)
		}
	}
	// If not found in sample version, fallthrough to GA
	switch fmt.Sprintf("%s/%s", product, resource) {
	case "assuredworkloads/workload":
		r := &assuredworkloads.Workload{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return AssuredWorkloadsWorkloadAsHCL(*r, hasGAEquivalent)
	case "cloudbuild/worker_pool":
		r := &cloudbuild.WorkerPool{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return CloudbuildWorkerPoolAsHCL(*r, hasGAEquivalent)
	case "cloudresourcemanager/folder":
		r := &cloudresourcemanager.Folder{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return CloudResourceManagerFolderAsHCL(*r, hasGAEquivalent)
	case "cloudresourcemanager/project":
		r := &cloudresourcemanager.Project{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return serializeGAProjectToHCL(*r, hasGAEquivalent)
	case "compute/firewall_policy":
		r := &compute.FirewallPolicy{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyAsHCL(*r, hasGAEquivalent)
	case "compute/firewall_policy_association":
		r := &compute.FirewallPolicyAssociation{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyAssociationAsHCL(*r, hasGAEquivalent)
	case "compute/firewall_policy_rule":
		r := &compute.FirewallPolicyRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyRuleAsHCL(*r, hasGAEquivalent)
	case "compute/forwarding_rule":
		r := &compute.ForwardingRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeForwardingRuleAsHCL(*r, hasGAEquivalent)
	case "compute/global_forwarding_rule":
		r := &compute.ForwardingRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeGlobalForwardingRuleAsHCL(*r, hasGAEquivalent)
	case "containeraws/cluster":
		r := &containeraws.Cluster{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ContainerAwsClusterAsHCL(*r, hasGAEquivalent)
	case "containeraws/node_pool":
		r := &containeraws.NodePool{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ContainerAwsNodePoolAsHCL(*r, hasGAEquivalent)
	case "containerazure/client":
		r := &containerazure.AzureClient{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ContainerAzureClientAsHCL(*r, hasGAEquivalent)
	case "containerazure/cluster":
		r := &containerazure.Cluster{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ContainerAzureClusterAsHCL(*r, hasGAEquivalent)
	case "containerazure/node_pool":
		r := &containerazure.NodePool{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ContainerAzureNodePoolAsHCL(*r, hasGAEquivalent)
	case "dataproc/workflow_template":
		r := &dataproc.WorkflowTemplate{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return DataprocWorkflowTemplateAsHCL(*r, hasGAEquivalent)
	case "eventarc/trigger":
		r := &eventarc.Trigger{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return EventarcTriggerAsHCL(*r, hasGAEquivalent)
	case "networkconnectivity/hub":
		r := &networkconnectivity.Hub{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return NetworkConnectivityHubAsHCL(*r, hasGAEquivalent)
	case "networkconnectivity/spoke":
		r := &networkconnectivity.Spoke{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return NetworkConnectivitySpokeAsHCL(*r, hasGAEquivalent)
	case "orgpolicy/policy":
		r := &orgpolicy.Policy{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return OrgPolicyPolicyAsHCL(*r, hasGAEquivalent)
	case "osconfig/os_policy_assignment":
		r := &osconfig.OSPolicyAssignment{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return OsConfigOsPolicyAssignmentAsHCL(*r, hasGAEquivalent)
	case "privateca/certificate_template":
		r := &privateca.CertificateTemplate{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return PrivatecaCertificateTemplateAsHCL(*r, hasGAEquivalent)
	case "recaptchaenterprise/key":
		r := &recaptchaenterprise.Key{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return RecaptchaEnterpriseKeyAsHCL(*r, hasGAEquivalent)
	default:
		return "", fmt.Errorf("Error converting sample JSON to HCL: %s not found", product, resource)
	}

}

// AssuredWorkloadsWorkloadBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func AssuredWorkloadsWorkloadBetaAsHCL(r assuredworkloadsBeta.Workload, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_assured_workloads_workload\" \"output\" {\n"
	if r.BillingAccount != nil {
		outputConfig += fmt.Sprintf("\tbilling_account = %#v\n", *r.BillingAccount)
	}
	if r.ComplianceRegime != nil {
		outputConfig += fmt.Sprintf("\tcompliance_regime = %#v\n", *r.ComplianceRegime)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Organization != nil {
		outputConfig += fmt.Sprintf("\torganization = %#v\n", *r.Organization)
	}
	if v := convertAssuredWorkloadsWorkloadBetaKmsSettingsToHCL(r.KmsSettings); v != "" {
		outputConfig += fmt.Sprintf("\tkms_settings %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.ProvisionedResourcesParent != nil {
		outputConfig += fmt.Sprintf("\tprovisioned_resources_parent = %#v\n", *r.ProvisionedResourcesParent)
	}
	if r.ResourceSettings != nil {
		for _, v := range r.ResourceSettings {
			outputConfig += fmt.Sprintf("\tresource_settings %s\n", convertAssuredWorkloadsWorkloadBetaResourceSettingsToHCL(&v))
		}
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertAssuredWorkloadsWorkloadBetaKmsSettingsToHCL(r *assuredworkloadsBeta.WorkloadKmsSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.NextRotationTime != nil {
		outputConfig += fmt.Sprintf("\tnext_rotation_time = %#v\n", *r.NextRotationTime)
	}
	if r.RotationPeriod != nil {
		outputConfig += fmt.Sprintf("\trotation_period = %#v\n", *r.RotationPeriod)
	}
	return outputConfig + "}"
}

func convertAssuredWorkloadsWorkloadBetaResourceSettingsToHCL(r *assuredworkloadsBeta.WorkloadResourceSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ResourceId != nil {
		outputConfig += fmt.Sprintf("\tresource_id = %#v\n", *r.ResourceId)
	}
	if r.ResourceType != nil {
		outputConfig += fmt.Sprintf("\tresource_type = %#v\n", *r.ResourceType)
	}
	return outputConfig + "}"
}

func convertAssuredWorkloadsWorkloadBetaResourcesToHCL(r *assuredworkloadsBeta.WorkloadResources) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// CloudbuildWorkerPoolBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudbuildWorkerPoolBetaAsHCL(r cloudbuildBeta.WorkerPool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_cloudbuild_worker_pool\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if v := convertCloudbuildWorkerPoolBetaNetworkConfigToHCL(r.NetworkConfig); v != "" {
		outputConfig += fmt.Sprintf("\tnetwork_config %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertCloudbuildWorkerPoolBetaWorkerConfigToHCL(r.WorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tworker_config %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertCloudbuildWorkerPoolBetaNetworkConfigToHCL(r *cloudbuildBeta.WorkerPoolNetworkConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PeeredNetwork != nil {
		outputConfig += fmt.Sprintf("\tpeered_network = %#v\n", *r.PeeredNetwork)
	}
	return outputConfig + "}"
}

func convertCloudbuildWorkerPoolBetaWorkerConfigToHCL(r *cloudbuildBeta.WorkerPoolWorkerConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DiskSizeGb != nil {
		outputConfig += fmt.Sprintf("\tdisk_size_gb = %#v\n", *r.DiskSizeGb)
	}
	if r.MachineType != nil {
		outputConfig += fmt.Sprintf("\tmachine_type = %#v\n", *r.MachineType)
	}
	if r.NoExternalIP != nil {
		outputConfig += fmt.Sprintf("\tno_external_ip = %#v\n", *r.NoExternalIP)
	}
	return outputConfig + "}"
}

// CloudResourceManagerFolderBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerFolderBetaAsHCL(r cloudresourcemanagerBeta.Folder, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_folder\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// CloudResourceManagerProjectBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerProjectBetaAsHCL(r cloudresourcemanagerBeta.Project, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_project\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplayname = %#v\n", *r.DisplayName)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyBetaAsHCL(r computeBeta.FirewallPolicy, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.ShortName != nil {
		outputConfig += fmt.Sprintf("\tshort_name = %#v\n", *r.ShortName)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyAssociationBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAssociationBetaAsHCL(r computeBeta.FirewallPolicyAssociation, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy_association\" \"output\" {\n"
	if r.AttachmentTarget != nil {
		outputConfig += fmt.Sprintf("\tattachment_target = %#v\n", *r.AttachmentTarget)
	}
	if r.FirewallPolicy != nil {
		outputConfig += fmt.Sprintf("\tfirewall_policy = %#v\n", *r.FirewallPolicy)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyRuleBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyRuleBetaAsHCL(r computeBeta.FirewallPolicyRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy_rule\" \"output\" {\n"
	if r.Action != nil {
		outputConfig += fmt.Sprintf("\taction = %#v\n", *r.Action)
	}
	if r.Direction != nil {
		outputConfig += fmt.Sprintf("\tdirection = %#v\n", *r.Direction)
	}
	if r.FirewallPolicy != nil {
		outputConfig += fmt.Sprintf("\tfirewall_policy = %#v\n", *r.FirewallPolicy)
	}
	if v := convertComputeFirewallPolicyRuleBetaMatchToHCL(r.Match); v != "" {
		outputConfig += fmt.Sprintf("\tmatch %s\n", v)
	}
	if r.Priority != nil {
		outputConfig += fmt.Sprintf("\tpriority = %#v\n", *r.Priority)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Disabled != nil {
		outputConfig += fmt.Sprintf("\tdisabled = %#v\n", *r.Disabled)
	}
	if r.EnableLogging != nil {
		outputConfig += fmt.Sprintf("\tenable_logging = %#v\n", *r.EnableLogging)
	}
	if r.TargetResources != nil {
		outputConfig += "\ttarget_resources = ["
		for _, v := range r.TargetResources {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.TargetServiceAccounts != nil {
		outputConfig += "\ttarget_service_accounts = ["
		for _, v := range r.TargetServiceAccounts {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertComputeFirewallPolicyRuleBetaMatchToHCL(r *computeBeta.FirewallPolicyRuleMatch) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Layer4Configs != nil {
		for _, v := range r.Layer4Configs {
			outputConfig += fmt.Sprintf("\tlayer4_configs %s\n", convertComputeFirewallPolicyRuleBetaMatchLayer4ConfigsToHCL(&v))
		}
	}
	if r.DestIPRanges != nil {
		outputConfig += "\tdest_ip_ranges = ["
		for _, v := range r.DestIPRanges {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.SrcIPRanges != nil {
		outputConfig += "\tsrc_ip_ranges = ["
		for _, v := range r.SrcIPRanges {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertComputeFirewallPolicyRuleBetaMatchLayer4ConfigsToHCL(r *computeBeta.FirewallPolicyRuleMatchLayer4Configs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.Ports != nil {
		outputConfig += "\tports = ["
		for _, v := range r.Ports {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// ComputeForwardingRuleBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeForwardingRuleBetaAsHCL(r computeBeta.ForwardingRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_forwarding_rule\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.AllPorts != nil {
		outputConfig += fmt.Sprintf("\tall_ports = %#v\n", *r.AllPorts)
	}
	if r.AllowGlobalAccess != nil {
		outputConfig += fmt.Sprintf("\tallow_global_access = %#v\n", *r.AllowGlobalAccess)
	}
	if r.BackendService != nil {
		outputConfig += fmt.Sprintf("\tbackend_service = %#v\n", *r.BackendService)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.IsMirroringCollector != nil {
		outputConfig += fmt.Sprintf("\tis_mirroring_collector = %#v\n", *r.IsMirroringCollector)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.LoadBalancingScheme != nil {
		outputConfig += fmt.Sprintf("\tload_balancing_scheme = %#v\n", *r.LoadBalancingScheme)
	}
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if r.NetworkTier != nil {
		outputConfig += fmt.Sprintf("\tnetwork_tier = %#v\n", *r.NetworkTier)
	}
	if r.PortRange != nil {
		outputConfig += fmt.Sprintf("\tport_range = %#v\n", *r.PortRange)
	}
	if r.Ports != nil {
		outputConfig += "\tports = ["
		for _, v := range r.Ports {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tregion = %#v\n", *r.Location)
	}
	if r.ServiceLabel != nil {
		outputConfig += fmt.Sprintf("\tservice_label = %#v\n", *r.ServiceLabel)
	}
	if r.Subnetwork != nil {
		outputConfig += fmt.Sprintf("\tsubnetwork = %#v\n", *r.Subnetwork)
	}
	if r.Target != nil {
		outputConfig += fmt.Sprintf("\ttarget = %#v\n", *r.Target)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeGlobalForwardingRuleBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeGlobalForwardingRuleBetaAsHCL(r computeBeta.ForwardingRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_global_forwarding_rule\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Target != nil {
		outputConfig += fmt.Sprintf("\ttarget = %#v\n", *r.Target)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.IPVersion != nil {
		outputConfig += fmt.Sprintf("\tip_version = %#v\n", *r.IPVersion)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.LoadBalancingScheme != nil {
		outputConfig += fmt.Sprintf("\tload_balancing_scheme = %#v\n", *r.LoadBalancingScheme)
	}
	if r.MetadataFilter != nil {
		for _, v := range r.MetadataFilter {
			outputConfig += fmt.Sprintf("\tmetadata_filters %s\n", convertComputeGlobalForwardingRuleBetaMetadataFilterToHCL(&v))
		}
	}
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if r.PortRange != nil {
		outputConfig += fmt.Sprintf("\tport_range = %#v\n", *r.PortRange)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertComputeGlobalForwardingRuleBetaMetadataFilterToHCL(r *computeBeta.ForwardingRuleMetadataFilter) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.FilterLabel != nil {
		for _, v := range r.FilterLabel {
			outputConfig += fmt.Sprintf("\tfilter_labels %s\n", convertComputeGlobalForwardingRuleBetaMetadataFilterFilterLabelToHCL(&v))
		}
	}
	if r.FilterMatchCriteria != nil {
		outputConfig += fmt.Sprintf("\tfilter_match_criteria = %#v\n", *r.FilterMatchCriteria)
	}
	return outputConfig + "}"
}

func convertComputeGlobalForwardingRuleBetaMetadataFilterFilterLabelToHCL(r *computeBeta.ForwardingRuleMetadataFilterFilterLabel) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

// ContainerAwsClusterBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAwsClusterBetaAsHCL(r containerawsBeta.Cluster, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_aws_cluster\" \"output\" {\n"
	if v := convertContainerAwsClusterBetaAuthorizationToHCL(r.Authorization); v != "" {
		outputConfig += fmt.Sprintf("\tauthorization %s\n", v)
	}
	if r.AwsRegion != nil {
		outputConfig += fmt.Sprintf("\taws_region = %#v\n", *r.AwsRegion)
	}
	if v := convertContainerAwsClusterBetaControlPlaneToHCL(r.ControlPlane); v != "" {
		outputConfig += fmt.Sprintf("\tcontrol_plane %s\n", v)
	}
	if v := convertContainerAwsClusterBetaFleetToHCL(r.Fleet); v != "" {
		outputConfig += fmt.Sprintf("\tfleet %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertContainerAwsClusterBetaNetworkingToHCL(r.Networking); v != "" {
		outputConfig += fmt.Sprintf("\tnetworking %s\n", v)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAwsClusterBetaAuthorizationToHCL(r *containerawsBeta.ClusterAuthorization) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdminUsers != nil {
		for _, v := range r.AdminUsers {
			outputConfig += fmt.Sprintf("\tadmin_users %s\n", convertContainerAwsClusterBetaAuthorizationAdminUsersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaAuthorizationAdminUsersToHCL(r *containerawsBeta.ClusterAuthorizationAdminUsers) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Username != nil {
		outputConfig += fmt.Sprintf("\tusername = %#v\n", *r.Username)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneToHCL(r *containerawsBeta.ClusterControlPlane) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAwsClusterBetaControlPlaneAwsServicesAuthenticationToHCL(r.AwsServicesAuthentication); v != "" {
		outputConfig += fmt.Sprintf("\taws_services_authentication %s\n", v)
	}
	if v := convertContainerAwsClusterBetaControlPlaneConfigEncryptionToHCL(r.ConfigEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_encryption %s\n", v)
	}
	if v := convertContainerAwsClusterBetaControlPlaneDatabaseEncryptionToHCL(r.DatabaseEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tdatabase_encryption %s\n", v)
	}
	if r.IamInstanceProfile != nil {
		outputConfig += fmt.Sprintf("\tiam_instance_profile = %#v\n", *r.IamInstanceProfile)
	}
	if r.SubnetIds != nil {
		outputConfig += "\tsubnet_ids = ["
		for _, v := range r.SubnetIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	if r.InstanceType != nil {
		outputConfig += fmt.Sprintf("\tinstance_type = %#v\n", *r.InstanceType)
	}
	if v := convertContainerAwsClusterBetaControlPlaneMainVolumeToHCL(r.MainVolume); v != "" {
		outputConfig += fmt.Sprintf("\tmain_volume %s\n", v)
	}
	if v := convertContainerAwsClusterBetaControlPlaneProxyConfigToHCL(r.ProxyConfig); v != "" {
		outputConfig += fmt.Sprintf("\tproxy_config %s\n", v)
	}
	if v := convertContainerAwsClusterBetaControlPlaneRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	if r.SecurityGroupIds != nil {
		outputConfig += "\tsecurity_group_ids = ["
		for _, v := range r.SecurityGroupIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertContainerAwsClusterBetaControlPlaneSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneAwsServicesAuthenticationToHCL(r *containerawsBeta.ClusterControlPlaneAwsServicesAuthentication) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.RoleArn != nil {
		outputConfig += fmt.Sprintf("\trole_arn = %#v\n", *r.RoleArn)
	}
	if r.RoleSessionName != nil {
		outputConfig += fmt.Sprintf("\trole_session_name = %#v\n", *r.RoleSessionName)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneConfigEncryptionToHCL(r *containerawsBeta.ClusterControlPlaneConfigEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneDatabaseEncryptionToHCL(r *containerawsBeta.ClusterControlPlaneDatabaseEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneMainVolumeToHCL(r *containerawsBeta.ClusterControlPlaneMainVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneProxyConfigToHCL(r *containerawsBeta.ClusterControlPlaneProxyConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SecretArn != nil {
		outputConfig += fmt.Sprintf("\tsecret_arn = %#v\n", *r.SecretArn)
	}
	if r.SecretVersion != nil {
		outputConfig += fmt.Sprintf("\tsecret_version = %#v\n", *r.SecretVersion)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneRootVolumeToHCL(r *containerawsBeta.ClusterControlPlaneRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaControlPlaneSshConfigToHCL(r *containerawsBeta.ClusterControlPlaneSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Ec2KeyPair != nil {
		outputConfig += fmt.Sprintf("\tec2_key_pair = %#v\n", *r.Ec2KeyPair)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaFleetToHCL(r *containerawsBeta.ClusterFleet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaNetworkingToHCL(r *containerawsBeta.ClusterNetworking) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PodAddressCidrBlocks != nil {
		outputConfig += "\tpod_address_cidr_blocks = ["
		for _, v := range r.PodAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ServiceAddressCidrBlocks != nil {
		outputConfig += "\tservice_address_cidr_blocks = ["
		for _, v := range r.ServiceAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.VPCId != nil {
		outputConfig += fmt.Sprintf("\tvpc_id = %#v\n", *r.VPCId)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterBetaWorkloadIdentityConfigToHCL(r *containerawsBeta.ClusterWorkloadIdentityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// ContainerAwsNodePoolBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAwsNodePoolBetaAsHCL(r containerawsBeta.NodePool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_aws_node_pool\" \"output\" {\n"
	if v := convertContainerAwsNodePoolBetaAutoscalingToHCL(r.Autoscaling); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling %s\n", v)
	}
	if r.Cluster != nil {
		outputConfig += fmt.Sprintf("\tcluster = %#v\n", *r.Cluster)
	}
	if v := convertContainerAwsNodePoolBetaConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if v := convertContainerAwsNodePoolBetaMaxPodsConstraintToHCL(r.MaxPodsConstraint); v != "" {
		outputConfig += fmt.Sprintf("\tmax_pods_constraint %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAwsNodePoolBetaAutoscalingToHCL(r *containerawsBeta.NodePoolAutoscaling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmax_node_count = %#v\n", *r.MaxNodeCount)
	}
	if r.MinNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmin_node_count = %#v\n", *r.MinNodeCount)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaConfigToHCL(r *containerawsBeta.NodePoolConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAwsNodePoolBetaConfigConfigEncryptionToHCL(r.ConfigEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_encryption %s\n", v)
	}
	if r.IamInstanceProfile != nil {
		outputConfig += fmt.Sprintf("\tiam_instance_profile = %#v\n", *r.IamInstanceProfile)
	}
	if r.InstanceType != nil {
		outputConfig += fmt.Sprintf("\tinstance_type = %#v\n", *r.InstanceType)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertContainerAwsNodePoolBetaConfigRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	if r.SecurityGroupIds != nil {
		outputConfig += "\tsecurity_group_ids = ["
		for _, v := range r.SecurityGroupIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertContainerAwsNodePoolBetaConfigSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Taints != nil {
		for _, v := range r.Taints {
			outputConfig += fmt.Sprintf("\ttaints %s\n", convertContainerAwsNodePoolBetaConfigTaintsToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaConfigConfigEncryptionToHCL(r *containerawsBeta.NodePoolConfigConfigEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaConfigRootVolumeToHCL(r *containerawsBeta.NodePoolConfigRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaConfigSshConfigToHCL(r *containerawsBeta.NodePoolConfigSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Ec2KeyPair != nil {
		outputConfig += fmt.Sprintf("\tec2_key_pair = %#v\n", *r.Ec2KeyPair)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaConfigTaintsToHCL(r *containerawsBeta.NodePoolConfigTaints) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Effect != nil {
		outputConfig += fmt.Sprintf("\teffect = %#v\n", *r.Effect)
	}
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolBetaMaxPodsConstraintToHCL(r *containerawsBeta.NodePoolMaxPodsConstraint) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxPodsPerNode != nil {
		outputConfig += fmt.Sprintf("\tmax_pods_per_node = %#v\n", *r.MaxPodsPerNode)
	}
	return outputConfig + "}"
}

// ContainerAzureClientBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureClientBetaAsHCL(r containerazureBeta.AzureClient, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_client\" \"output\" {\n"
	if r.ApplicationId != nil {
		outputConfig += fmt.Sprintf("\tapplication_id = %#v\n", *r.ApplicationId)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.TenantId != nil {
		outputConfig += fmt.Sprintf("\ttenant_id = %#v\n", *r.TenantId)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ContainerAzureClusterBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureClusterBetaAsHCL(r containerazureBeta.Cluster, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_cluster\" \"output\" {\n"
	if v := convertContainerAzureClusterBetaAuthorizationToHCL(r.Authorization); v != "" {
		outputConfig += fmt.Sprintf("\tauthorization %s\n", v)
	}
	if r.AzureRegion != nil {
		outputConfig += fmt.Sprintf("\tazure_region = %#v\n", *r.AzureRegion)
	}
	if r.Client != nil {
		outputConfig += fmt.Sprintf("\tclient = %#v\n", *r.Client)
	}
	if v := convertContainerAzureClusterBetaControlPlaneToHCL(r.ControlPlane); v != "" {
		outputConfig += fmt.Sprintf("\tcontrol_plane %s\n", v)
	}
	if v := convertContainerAzureClusterBetaFleetToHCL(r.Fleet); v != "" {
		outputConfig += fmt.Sprintf("\tfleet %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertContainerAzureClusterBetaNetworkingToHCL(r.Networking); v != "" {
		outputConfig += fmt.Sprintf("\tnetworking %s\n", v)
	}
	if r.ResourceGroupId != nil {
		outputConfig += fmt.Sprintf("\tresource_group_id = %#v\n", *r.ResourceGroupId)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAzureClusterBetaAuthorizationToHCL(r *containerazureBeta.ClusterAuthorization) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdminUsers != nil {
		for _, v := range r.AdminUsers {
			outputConfig += fmt.Sprintf("\tadmin_users %s\n", convertContainerAzureClusterBetaAuthorizationAdminUsersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaAuthorizationAdminUsersToHCL(r *containerazureBeta.ClusterAuthorizationAdminUsers) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Username != nil {
		outputConfig += fmt.Sprintf("\tusername = %#v\n", *r.Username)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneToHCL(r *containerazureBeta.ClusterControlPlane) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAzureClusterBetaControlPlaneSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	if v := convertContainerAzureClusterBetaControlPlaneDatabaseEncryptionToHCL(r.DatabaseEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tdatabase_encryption %s\n", v)
	}
	if v := convertContainerAzureClusterBetaControlPlaneMainVolumeToHCL(r.MainVolume); v != "" {
		outputConfig += fmt.Sprintf("\tmain_volume %s\n", v)
	}
	if v := convertContainerAzureClusterBetaControlPlaneProxyConfigToHCL(r.ProxyConfig); v != "" {
		outputConfig += fmt.Sprintf("\tproxy_config %s\n", v)
	}
	if r.ReplicaPlacements != nil {
		for _, v := range r.ReplicaPlacements {
			outputConfig += fmt.Sprintf("\treplica_placements %s\n", convertContainerAzureClusterBetaControlPlaneReplicaPlacementsToHCL(&v))
		}
	}
	if v := convertContainerAzureClusterBetaControlPlaneRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.VmSize != nil {
		outputConfig += fmt.Sprintf("\tvm_size = %#v\n", *r.VmSize)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneSshConfigToHCL(r *containerazureBeta.ClusterControlPlaneSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AuthorizedKey != nil {
		outputConfig += fmt.Sprintf("\tauthorized_key = %#v\n", *r.AuthorizedKey)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneDatabaseEncryptionToHCL(r *containerazureBeta.ClusterControlPlaneDatabaseEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KeyId != nil {
		outputConfig += fmt.Sprintf("\tkey_id = %#v\n", *r.KeyId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneMainVolumeToHCL(r *containerazureBeta.ClusterControlPlaneMainVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneProxyConfigToHCL(r *containerazureBeta.ClusterControlPlaneProxyConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ResourceGroupId != nil {
		outputConfig += fmt.Sprintf("\tresource_group_id = %#v\n", *r.ResourceGroupId)
	}
	if r.SecretId != nil {
		outputConfig += fmt.Sprintf("\tsecret_id = %#v\n", *r.SecretId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneReplicaPlacementsToHCL(r *containerazureBeta.ClusterControlPlaneReplicaPlacements) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AzureAvailabilityZone != nil {
		outputConfig += fmt.Sprintf("\tazure_availability_zone = %#v\n", *r.AzureAvailabilityZone)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaControlPlaneRootVolumeToHCL(r *containerazureBeta.ClusterControlPlaneRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaFleetToHCL(r *containerazureBeta.ClusterFleet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaNetworkingToHCL(r *containerazureBeta.ClusterNetworking) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PodAddressCidrBlocks != nil {
		outputConfig += "\tpod_address_cidr_blocks = ["
		for _, v := range r.PodAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ServiceAddressCidrBlocks != nil {
		outputConfig += "\tservice_address_cidr_blocks = ["
		for _, v := range r.ServiceAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.VirtualNetworkId != nil {
		outputConfig += fmt.Sprintf("\tvirtual_network_id = %#v\n", *r.VirtualNetworkId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterBetaWorkloadIdentityConfigToHCL(r *containerazureBeta.ClusterWorkloadIdentityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// ContainerAzureNodePoolBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureNodePoolBetaAsHCL(r containerazureBeta.NodePool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_node_pool\" \"output\" {\n"
	if v := convertContainerAzureNodePoolBetaAutoscalingToHCL(r.Autoscaling); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling %s\n", v)
	}
	if r.Cluster != nil {
		outputConfig += fmt.Sprintf("\tcluster = %#v\n", *r.Cluster)
	}
	if v := convertContainerAzureNodePoolBetaConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if v := convertContainerAzureNodePoolBetaMaxPodsConstraintToHCL(r.MaxPodsConstraint); v != "" {
		outputConfig += fmt.Sprintf("\tmax_pods_constraint %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.AzureAvailabilityZone != nil {
		outputConfig += fmt.Sprintf("\tazure_availability_zone = %#v\n", *r.AzureAvailabilityZone)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAzureNodePoolBetaAutoscalingToHCL(r *containerazureBeta.NodePoolAutoscaling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmax_node_count = %#v\n", *r.MaxNodeCount)
	}
	if r.MinNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmin_node_count = %#v\n", *r.MinNodeCount)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolBetaConfigToHCL(r *containerazureBeta.NodePoolConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAzureNodePoolBetaConfigSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	if v := convertContainerAzureNodePoolBetaConfigRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.VmSize != nil {
		outputConfig += fmt.Sprintf("\tvm_size = %#v\n", *r.VmSize)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolBetaConfigSshConfigToHCL(r *containerazureBeta.NodePoolConfigSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AuthorizedKey != nil {
		outputConfig += fmt.Sprintf("\tauthorized_key = %#v\n", *r.AuthorizedKey)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolBetaConfigRootVolumeToHCL(r *containerazureBeta.NodePoolConfigRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolBetaMaxPodsConstraintToHCL(r *containerazureBeta.NodePoolMaxPodsConstraint) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxPodsPerNode != nil {
		outputConfig += fmt.Sprintf("\tmax_pods_per_node = %#v\n", *r.MaxPodsPerNode)
	}
	return outputConfig + "}"
}

// DataprocWorkflowTemplateBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func DataprocWorkflowTemplateBetaAsHCL(r dataprocBeta.WorkflowTemplate, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_dataproc_workflow_template\" \"output\" {\n"
	if r.Jobs != nil {
		for _, v := range r.Jobs {
			outputConfig += fmt.Sprintf("\tjobs %s\n", convertDataprocWorkflowTemplateBetaJobsToHCL(&v))
		}
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertDataprocWorkflowTemplateBetaPlacementToHCL(r.Placement); v != "" {
		outputConfig += fmt.Sprintf("\tplacement %s\n", v)
	}
	if r.DagTimeout != nil {
		outputConfig += fmt.Sprintf("\tdag_timeout = %#v\n", *r.DagTimeout)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Parameters != nil {
		for _, v := range r.Parameters {
			outputConfig += fmt.Sprintf("\tparameters %s\n", convertDataprocWorkflowTemplateBetaParametersToHCL(&v))
		}
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertDataprocWorkflowTemplateBetaJobsToHCL(r *dataprocBeta.WorkflowTemplateJobs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.StepId != nil {
		outputConfig += fmt.Sprintf("\tstep_id = %#v\n", *r.StepId)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsHadoopJobToHCL(r.HadoopJob); v != "" {
		outputConfig += fmt.Sprintf("\thadoop_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsHiveJobToHCL(r.HiveJob); v != "" {
		outputConfig += fmt.Sprintf("\thive_job %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertDataprocWorkflowTemplateBetaJobsPigJobToHCL(r.PigJob); v != "" {
		outputConfig += fmt.Sprintf("\tpig_job %s\n", v)
	}
	if r.PrerequisiteStepIds != nil {
		outputConfig += "\tprerequisite_step_ids = ["
		for _, v := range r.PrerequisiteStepIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPrestoJobToHCL(r.PrestoJob); v != "" {
		outputConfig += fmt.Sprintf("\tpresto_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPysparkJobToHCL(r.PysparkJob); v != "" {
		outputConfig += fmt.Sprintf("\tpyspark_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSchedulingToHCL(r.Scheduling); v != "" {
		outputConfig += fmt.Sprintf("\tscheduling %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkJobToHCL(r.SparkJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkRJobToHCL(r.SparkRJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_r_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkSqlJobToHCL(r.SparkSqlJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_sql_job %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobToHCL(r *dataprocBeta.WorkflowTemplateJobsHadoopJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.MainClass != nil {
		outputConfig += fmt.Sprintf("\tmain_class = %#v\n", *r.MainClass)
	}
	if r.MainJarFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_jar_file_uri = %#v\n", *r.MainJarFileUri)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsHadoopJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsHiveJobToHCL(r *dataprocBeta.WorkflowTemplateJobsHiveJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsHiveJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsHiveJobQueryListToHCL(r *dataprocBeta.WorkflowTemplateJobsHiveJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPigJobToHCL(r *dataprocBeta.WorkflowTemplateJobsPigJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPigJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsPigJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPigJobQueryListToHCL(r *dataprocBeta.WorkflowTemplateJobsPigJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobToHCL(r *dataprocBeta.WorkflowTemplateJobsPrestoJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClientTags != nil {
		outputConfig += "\tclient_tags = ["
		for _, v := range r.ClientTags {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.OutputFormat != nil {
		outputConfig += fmt.Sprintf("\toutput_format = %#v\n", *r.OutputFormat)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsPrestoJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryListToHCL(r *dataprocBeta.WorkflowTemplateJobsPrestoJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJobToHCL(r *dataprocBeta.WorkflowTemplateJobsPysparkJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MainPythonFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_python_file_uri = %#v\n", *r.MainPythonFileUri)
	}
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.PythonFileUris != nil {
		outputConfig += "\tpython_file_uris = ["
		for _, v := range r.PythonFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsPysparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSchedulingToHCL(r *dataprocBeta.WorkflowTemplateJobsScheduling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxFailuresPerHour != nil {
		outputConfig += fmt.Sprintf("\tmax_failures_per_hour = %#v\n", *r.MaxFailuresPerHour)
	}
	if r.MaxFailuresTotal != nil {
		outputConfig += fmt.Sprintf("\tmax_failures_total = %#v\n", *r.MaxFailuresTotal)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.MainClass != nil {
		outputConfig += fmt.Sprintf("\tmain_class = %#v\n", *r.MainClass)
	}
	if r.MainJarFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_jar_file_uri = %#v\n", *r.MainJarFileUri)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkRJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MainRFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_r_file_uri = %#v\n", *r.MainRFileUri)
	}
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkRJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkSqlJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkSqlJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryListToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkSqlJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaPlacementToHCL(r *dataprocBeta.WorkflowTemplatePlacement) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateBetaPlacementClusterSelectorToHCL(r.ClusterSelector); v != "" {
		outputConfig += fmt.Sprintf("\tcluster_selector %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaPlacementManagedClusterToHCL(r.ManagedCluster); v != "" {
		outputConfig += fmt.Sprintf("\tmanaged_cluster %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaPlacementClusterSelectorToHCL(r *dataprocBeta.WorkflowTemplatePlacementClusterSelector) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tcluster_labels = {"
	for k, v := range r.ClusterLabels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Zone != nil {
		outputConfig += fmt.Sprintf("\tzone = %#v\n", *r.Zone)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaPlacementManagedClusterToHCL(r *dataprocBeta.WorkflowTemplatePlacementManagedCluster) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClusterName != nil {
		outputConfig += fmt.Sprintf("\tcluster_name = %#v\n", *r.ClusterName)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaParametersToHCL(r *dataprocBeta.WorkflowTemplateParameters) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Fields != nil {
		outputConfig += "\tfields = ["
		for _, v := range r.Fields {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if v := convertDataprocWorkflowTemplateBetaParametersValidationToHCL(r.Validation); v != "" {
		outputConfig += fmt.Sprintf("\tvalidation %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaParametersValidationToHCL(r *dataprocBeta.WorkflowTemplateParametersValidation) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateBetaParametersValidationRegexToHCL(r.Regex); v != "" {
		outputConfig += fmt.Sprintf("\tregex %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaParametersValidationValuesToHCL(r.Values); v != "" {
		outputConfig += fmt.Sprintf("\tvalues %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaParametersValidationRegexToHCL(r *dataprocBeta.WorkflowTemplateParametersValidationRegex) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Regexes != nil {
		outputConfig += "\tregexes = ["
		for _, v := range r.Regexes {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaParametersValidationValuesToHCL(r *dataprocBeta.WorkflowTemplateParametersValidationValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Values != nil {
		outputConfig += "\tvalues = ["
		for _, v := range r.Values {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigToHCL(r *dataprocBeta.ClusterInstanceGroupConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Accelerators != nil {
		for _, v := range r.Accelerators {
			outputConfig += fmt.Sprintf("\taccelerators %s\n", convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigAcceleratorsToHCL(&v))
		}
	}
	if v := convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfigToHCL(r.DiskConfig); v != "" {
		outputConfig += fmt.Sprintf("\tdisk_config %s\n", v)
	}
	if r.Image != nil {
		outputConfig += fmt.Sprintf("\timage = %#v\n", *r.Image)
	}
	if r.MachineType != nil {
		outputConfig += fmt.Sprintf("\tmachine_type = %#v\n", *r.MachineType)
	}
	if r.MinCpuPlatform != nil {
		outputConfig += fmt.Sprintf("\tmin_cpu_platform = %#v\n", *r.MinCpuPlatform)
	}
	if r.NumInstances != nil {
		outputConfig += fmt.Sprintf("\tnum_instances = %#v\n", *r.NumInstances)
	}
	if r.Preemptibility != nil {
		outputConfig += fmt.Sprintf("\tpreemptibility = %#v\n", *r.Preemptibility)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigAcceleratorsToHCL(r *dataprocBeta.ClusterInstanceGroupConfigAccelerators) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AcceleratorCount != nil {
		outputConfig += fmt.Sprintf("\taccelerator_count = %#v\n", *r.AcceleratorCount)
	}
	if r.AcceleratorType != nil {
		outputConfig += fmt.Sprintf("\taccelerator_type = %#v\n", *r.AcceleratorType)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfigToHCL(r *dataprocBeta.ClusterInstanceGroupConfigDiskConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BootDiskSizeGb != nil {
		outputConfig += fmt.Sprintf("\tboot_disk_size_gb = %#v\n", *r.BootDiskSizeGb)
	}
	if r.BootDiskType != nil {
		outputConfig += fmt.Sprintf("\tboot_disk_type = %#v\n", *r.BootDiskType)
	}
	if r.NumLocalSsds != nil {
		outputConfig += fmt.Sprintf("\tnum_local_ssds = %#v\n", *r.NumLocalSsds)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigManagedGroupConfigToHCL(r *dataprocBeta.ClusterInstanceGroupConfigManagedGroupConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigToHCL(r *dataprocBeta.ClusterClusterConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfigToHCL(r.AutoscalingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfigToHCL(r.EncryptionConfig); v != "" {
		outputConfig += fmt.Sprintf("\tencryption_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfigToHCL(r.EndpointConfig); v != "" {
		outputConfig += fmt.Sprintf("\tendpoint_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigToHCL(r.GceClusterConfig); v != "" {
		outputConfig += fmt.Sprintf("\tgce_cluster_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigToHCL(r.GkeClusterConfig); v != "" {
		outputConfig += fmt.Sprintf("\tgke_cluster_config %s\n", v)
	}
	if r.InitializationActions != nil {
		for _, v := range r.InitializationActions {
			outputConfig += fmt.Sprintf("\tinitialization_actions %s\n", convertDataprocWorkflowTemplateBetaClusterClusterConfigInitializationActionsToHCL(&v))
		}
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfigToHCL(r.LifecycleConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlifecycle_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigToHCL(r.MasterConfig); v != "" {
		outputConfig += fmt.Sprintf("\tmaster_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfigToHCL(r.MetastoreConfig); v != "" {
		outputConfig += fmt.Sprintf("\tmetastore_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigToHCL(r.SecondaryWorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsecondary_worker_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigToHCL(r.SecurityConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsecurity_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfigToHCL(r.SoftwareConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsoftware_config %s\n", v)
	}
	if r.StagingBucket != nil {
		outputConfig += fmt.Sprintf("\tstaging_bucket = %#v\n", *r.StagingBucket)
	}
	if r.TempBucket != nil {
		outputConfig += fmt.Sprintf("\ttemp_bucket = %#v\n", *r.TempBucket)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigToHCL(r.WorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tworker_config %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfigToHCL(r *dataprocBeta.ClusterClusterConfigAutoscalingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Policy != nil {
		outputConfig += fmt.Sprintf("\tpolicy = %#v\n", *r.Policy)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfigToHCL(r *dataprocBeta.ClusterClusterConfigEncryptionConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.GcePdKmsKeyName != nil {
		outputConfig += fmt.Sprintf("\tgce_pd_kms_key_name = %#v\n", *r.GcePdKmsKeyName)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfigToHCL(r *dataprocBeta.ClusterClusterConfigEndpointConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.EnableHttpPortAccess != nil {
		outputConfig += fmt.Sprintf("\tenable_http_port_access = %#v\n", *r.EnableHttpPortAccess)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigToHCL(r *dataprocBeta.ClusterClusterConfigGceClusterConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.InternalIPOnly != nil {
		outputConfig += fmt.Sprintf("\tinternal_ip_only = %#v\n", *r.InternalIPOnly)
	}
	outputConfig += "\tmetadata = {"
	for k, v := range r.Metadata {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinityToHCL(r.NodeGroupAffinity); v != "" {
		outputConfig += fmt.Sprintf("\tnode_group_affinity %s\n", v)
	}
	if r.PrivateIPv6GoogleAccess != nil {
		outputConfig += fmt.Sprintf("\tprivate_ipv6_google_access = %#v\n", *r.PrivateIPv6GoogleAccess)
	}
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinityToHCL(r.ReservationAffinity); v != "" {
		outputConfig += fmt.Sprintf("\treservation_affinity %s\n", v)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if r.ServiceAccountScopes != nil {
		outputConfig += "\tservice_account_scopes = ["
		for _, v := range r.ServiceAccountScopes {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Subnetwork != nil {
		outputConfig += fmt.Sprintf("\tsubnetwork = %#v\n", *r.Subnetwork)
	}
	if r.Tags != nil {
		outputConfig += "\ttags = ["
		for _, v := range r.Tags {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Zone != nil {
		outputConfig += fmt.Sprintf("\tzone = %#v\n", *r.Zone)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinityToHCL(r *dataprocBeta.ClusterClusterConfigGceClusterConfigNodeGroupAffinity) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.NodeGroup != nil {
		outputConfig += fmt.Sprintf("\tnode_group = %#v\n", *r.NodeGroup)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinityToHCL(r *dataprocBeta.ClusterClusterConfigGceClusterConfigReservationAffinity) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ConsumeReservationType != nil {
		outputConfig += fmt.Sprintf("\tconsume_reservation_type = %#v\n", *r.ConsumeReservationType)
	}
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Values != nil {
		outputConfig += "\tvalues = ["
		for _, v := range r.Values {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigToHCL(r *dataprocBeta.ClusterClusterConfigGkeClusterConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTargetToHCL(r.NamespacedGkeDeploymentTarget); v != "" {
		outputConfig += fmt.Sprintf("\tnamespaced_gke_deployment_target %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTargetToHCL(r *dataprocBeta.ClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTarget) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClusterNamespace != nil {
		outputConfig += fmt.Sprintf("\tcluster_namespace = %#v\n", *r.ClusterNamespace)
	}
	if r.TargetGkeCluster != nil {
		outputConfig += fmt.Sprintf("\ttarget_gke_cluster = %#v\n", *r.TargetGkeCluster)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigInitializationActionsToHCL(r *dataprocBeta.ClusterClusterConfigInitializationActions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ExecutableFile != nil {
		outputConfig += fmt.Sprintf("\texecutable_file = %#v\n", *r.ExecutableFile)
	}
	if r.ExecutionTimeout != nil {
		outputConfig += fmt.Sprintf("\texecution_timeout = %#v\n", *r.ExecutionTimeout)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfigToHCL(r *dataprocBeta.ClusterClusterConfigLifecycleConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AutoDeleteTime != nil {
		outputConfig += fmt.Sprintf("\tauto_delete_time = %#v\n", *r.AutoDeleteTime)
	}
	if r.AutoDeleteTtl != nil {
		outputConfig += fmt.Sprintf("\tauto_delete_ttl = %#v\n", *r.AutoDeleteTtl)
	}
	if r.IdleDeleteTtl != nil {
		outputConfig += fmt.Sprintf("\tidle_delete_ttl = %#v\n", *r.IdleDeleteTtl)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfigToHCL(r *dataprocBeta.ClusterClusterConfigMetastoreConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DataprocMetastoreService != nil {
		outputConfig += fmt.Sprintf("\tdataproc_metastore_service = %#v\n", *r.DataprocMetastoreService)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigToHCL(r *dataprocBeta.ClusterClusterConfigSecurityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfigToHCL(r.KerberosConfig); v != "" {
		outputConfig += fmt.Sprintf("\tkerberos_config %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfigToHCL(r *dataprocBeta.ClusterClusterConfigSecurityConfigKerberosConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CrossRealmTrustAdminServer != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_admin_server = %#v\n", *r.CrossRealmTrustAdminServer)
	}
	if r.CrossRealmTrustKdc != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_kdc = %#v\n", *r.CrossRealmTrustKdc)
	}
	if r.CrossRealmTrustRealm != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_realm = %#v\n", *r.CrossRealmTrustRealm)
	}
	if r.CrossRealmTrustSharedPassword != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_shared_password = %#v\n", *r.CrossRealmTrustSharedPassword)
	}
	if r.EnableKerberos != nil {
		outputConfig += fmt.Sprintf("\tenable_kerberos = %#v\n", *r.EnableKerberos)
	}
	if r.KdcDbKey != nil {
		outputConfig += fmt.Sprintf("\tkdc_db_key = %#v\n", *r.KdcDbKey)
	}
	if r.KeyPassword != nil {
		outputConfig += fmt.Sprintf("\tkey_password = %#v\n", *r.KeyPassword)
	}
	if r.Keystore != nil {
		outputConfig += fmt.Sprintf("\tkeystore = %#v\n", *r.Keystore)
	}
	if r.KeystorePassword != nil {
		outputConfig += fmt.Sprintf("\tkeystore_password = %#v\n", *r.KeystorePassword)
	}
	if r.KmsKey != nil {
		outputConfig += fmt.Sprintf("\tkms_key = %#v\n", *r.KmsKey)
	}
	if r.Realm != nil {
		outputConfig += fmt.Sprintf("\trealm = %#v\n", *r.Realm)
	}
	if r.RootPrincipalPassword != nil {
		outputConfig += fmt.Sprintf("\troot_principal_password = %#v\n", *r.RootPrincipalPassword)
	}
	if r.TgtLifetimeHours != nil {
		outputConfig += fmt.Sprintf("\ttgt_lifetime_hours = %#v\n", *r.TgtLifetimeHours)
	}
	if r.Truststore != nil {
		outputConfig += fmt.Sprintf("\ttruststore = %#v\n", *r.Truststore)
	}
	if r.TruststorePassword != nil {
		outputConfig += fmt.Sprintf("\ttruststore_password = %#v\n", *r.TruststorePassword)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfigToHCL(r *dataprocBeta.ClusterClusterConfigSoftwareConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ImageVersion != nil {
		outputConfig += fmt.Sprintf("\timage_version = %#v\n", *r.ImageVersion)
	}
	if r.OptionalComponents != nil {
		outputConfig += "\toptional_components = ["
		for _, v := range r.OptionalComponents {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

// EventarcTriggerBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func EventarcTriggerBetaAsHCL(r eventarcBeta.Trigger, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_eventarc_trigger\" \"output\" {\n"
	if v := convertEventarcTriggerBetaDestinationToHCL(r.Destination); v != "" {
		outputConfig += fmt.Sprintf("\tdestination %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.MatchingCriteria != nil {
		for _, v := range r.MatchingCriteria {
			outputConfig += fmt.Sprintf("\tmatching_criteria %s\n", convertEventarcTriggerBetaMatchingCriteriaToHCL(&v))
		}
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if v := convertEventarcTriggerBetaTransportToHCL(r.Transport); v != "" {
		outputConfig += fmt.Sprintf("\ttransport %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertEventarcTriggerBetaDestinationToHCL(r *eventarcBeta.TriggerDestination) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CloudFunction != nil {
		outputConfig += fmt.Sprintf("\tcloud_function = %#v\n", *r.CloudFunction)
	}
	if v := convertEventarcTriggerBetaDestinationCloudRunServiceToHCL(r.CloudRunService); v != "" {
		outputConfig += fmt.Sprintf("\tcloud_run_service %s\n", v)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerBetaDestinationCloudRunServiceToHCL(r *eventarcBeta.TriggerDestinationCloudRunService) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Service != nil {
		outputConfig += fmt.Sprintf("\tservice = %#v\n", *r.Service)
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.Region != nil {
		outputConfig += fmt.Sprintf("\tregion = %#v\n", *r.Region)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerBetaMatchingCriteriaToHCL(r *eventarcBeta.TriggerMatchingCriteria) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Attribute != nil {
		outputConfig += fmt.Sprintf("\tattribute = %#v\n", *r.Attribute)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerBetaTransportToHCL(r *eventarcBeta.TriggerTransport) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertEventarcTriggerBetaTransportPubsubToHCL(r.Pubsub); v != "" {
		outputConfig += fmt.Sprintf("\tpubsub %s\n", v)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerBetaTransportPubsubToHCL(r *eventarcBeta.TriggerTransportPubsub) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Topic != nil {
		outputConfig += fmt.Sprintf("\ttopic = %#v\n", *r.Topic)
	}
	return outputConfig + "}"
}

// GkeHubFeatureBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func GkeHubFeatureBetaAsHCL(r gkehubBeta.Feature, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_gke_hub_feature\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertGkeHubFeatureBetaSpecToHCL(r.Spec); v != "" {
		outputConfig += fmt.Sprintf("\tspec %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertGkeHubFeatureBetaSpecToHCL(r *gkehubBeta.FeatureSpec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertGkeHubFeatureBetaSpecMulticlusteringressToHCL(r.Multiclusteringress); v != "" {
		outputConfig += fmt.Sprintf("\tmulticlusteringress %s\n", v)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureBetaSpecMulticlusteringressToHCL(r *gkehubBeta.FeatureSpecMulticlusteringress) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ConfigMembership != nil {
		outputConfig += fmt.Sprintf("\tconfig_membership = %#v\n", *r.ConfigMembership)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureBetaResourceStateToHCL(r *gkehubBeta.FeatureResourceState) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertGkeHubFeatureBetaStateToHCL(r *gkehubBeta.FeatureState) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertGkeHubFeatureBetaStateStateToHCL(r *gkehubBeta.FeatureStateState) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// GkeHubFeatureMembershipBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func GkeHubFeatureMembershipBetaAsHCL(r gkehubBeta.FeatureMembership, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_gke_hub_feature_membership\" \"output\" {\n"
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementToHCL(r.Configmanagement); v != "" {
		outputConfig += fmt.Sprintf("\tconfigmanagement %s\n", v)
	}
	if r.Feature != nil {
		outputConfig += fmt.Sprintf("\tfeature = %#v\n", *r.Feature)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Membership != nil {
		outputConfig += fmt.Sprintf("\tmembership = %#v\n", *r.Membership)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertGkeHubFeatureMembershipBetaConfigmanagementToHCL(r *gkehubBeta.FeatureMembershipConfigmanagement) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementBinauthzToHCL(r.Binauthz); v != "" {
		outputConfig += fmt.Sprintf("\tbinauthz %s\n", v)
	}
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncToHCL(r.ConfigSync); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_sync %s\n", v)
	}
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyControllerToHCL(r.HierarchyController); v != "" {
		outputConfig += fmt.Sprintf("\thierarchy_controller %s\n", v)
	}
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementPolicyControllerToHCL(r.PolicyController); v != "" {
		outputConfig += fmt.Sprintf("\tpolicy_controller %s\n", v)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureMembershipBetaConfigmanagementBinauthzToHCL(r *gkehubBeta.FeatureMembershipConfigmanagementBinauthz) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Enabled != nil {
		outputConfig += fmt.Sprintf("\tenabled = %#v\n", *r.Enabled)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncToHCL(r *gkehubBeta.FeatureMembershipConfigmanagementConfigSync) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGitToHCL(r.Git); v != "" {
		outputConfig += fmt.Sprintf("\tgit %s\n", v)
	}
	if r.SourceFormat != nil {
		outputConfig += fmt.Sprintf("\tsource_format = %#v\n", *r.SourceFormat)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGitToHCL(r *gkehubBeta.FeatureMembershipConfigmanagementConfigSyncGit) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.GcpServiceAccountEmail != nil {
		outputConfig += fmt.Sprintf("\tgcp_service_account_email = %#v\n", *r.GcpServiceAccountEmail)
	}
	if r.HttpsProxy != nil {
		outputConfig += fmt.Sprintf("\thttps_proxy = %#v\n", *r.HttpsProxy)
	}
	if r.PolicyDir != nil {
		outputConfig += fmt.Sprintf("\tpolicy_dir = %#v\n", *r.PolicyDir)
	}
	if r.SecretType != nil {
		outputConfig += fmt.Sprintf("\tsecret_type = %#v\n", *r.SecretType)
	}
	if r.SyncBranch != nil {
		outputConfig += fmt.Sprintf("\tsync_branch = %#v\n", *r.SyncBranch)
	}
	if r.SyncRepo != nil {
		outputConfig += fmt.Sprintf("\tsync_repo = %#v\n", *r.SyncRepo)
	}
	if r.SyncRev != nil {
		outputConfig += fmt.Sprintf("\tsync_rev = %#v\n", *r.SyncRev)
	}
	if r.SyncWaitSecs != nil {
		outputConfig += fmt.Sprintf("\tsync_wait_secs = %#v\n", *r.SyncWaitSecs)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyControllerToHCL(r *gkehubBeta.FeatureMembershipConfigmanagementHierarchyController) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.EnableHierarchicalResourceQuota != nil {
		outputConfig += fmt.Sprintf("\tenable_hierarchical_resource_quota = %#v\n", *r.EnableHierarchicalResourceQuota)
	}
	if r.EnablePodTreeLabels != nil {
		outputConfig += fmt.Sprintf("\tenable_pod_tree_labels = %#v\n", *r.EnablePodTreeLabels)
	}
	if r.Enabled != nil {
		outputConfig += fmt.Sprintf("\tenabled = %#v\n", *r.Enabled)
	}
	return outputConfig + "}"
}

func convertGkeHubFeatureMembershipBetaConfigmanagementPolicyControllerToHCL(r *gkehubBeta.FeatureMembershipConfigmanagementPolicyController) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AuditIntervalSeconds != nil {
		outputConfig += fmt.Sprintf("\taudit_interval_seconds = %#v\n", *r.AuditIntervalSeconds)
	}
	if r.Enabled != nil {
		outputConfig += fmt.Sprintf("\tenabled = %#v\n", *r.Enabled)
	}
	if r.ExemptableNamespaces != nil {
		outputConfig += "\texemptable_namespaces = ["
		for _, v := range r.ExemptableNamespaces {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.LogDeniesEnabled != nil {
		outputConfig += fmt.Sprintf("\tlog_denies_enabled = %#v\n", *r.LogDeniesEnabled)
	}
	if r.ReferentialRulesEnabled != nil {
		outputConfig += fmt.Sprintf("\treferential_rules_enabled = %#v\n", *r.ReferentialRulesEnabled)
	}
	if r.TemplateLibraryInstalled != nil {
		outputConfig += fmt.Sprintf("\ttemplate_library_installed = %#v\n", *r.TemplateLibraryInstalled)
	}
	return outputConfig + "}"
}

// MonitoringMonitoredProjectBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func MonitoringMonitoredProjectBetaAsHCL(r monitoringBeta.MonitoredProject, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_monitoring_monitored_project\" \"output\" {\n"
	if r.MetricsScope != nil {
		outputConfig += fmt.Sprintf("\tmetrics_scope = %#v\n", *r.MetricsScope)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// NetworkConnectivityHubBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func NetworkConnectivityHubBetaAsHCL(r networkconnectivityBeta.Hub, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_network_connectivity_hub\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertNetworkConnectivityHubBetaRoutingVpcsToHCL(r *networkconnectivityBeta.HubRoutingVpcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// NetworkConnectivitySpokeBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func NetworkConnectivitySpokeBetaAsHCL(r networkconnectivityBeta.Spoke, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_network_connectivity_spoke\" \"output\" {\n"
	if r.Hub != nil {
		outputConfig += fmt.Sprintf("\thub = %#v\n", *r.Hub)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertNetworkConnectivitySpokeBetaLinkedInterconnectAttachmentsToHCL(r.LinkedInterconnectAttachments); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_interconnect_attachments %s\n", v)
	}
	if v := convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesToHCL(r.LinkedRouterApplianceInstances); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_router_appliance_instances %s\n", v)
	}
	if v := convertNetworkConnectivitySpokeBetaLinkedVpnTunnelsToHCL(r.LinkedVpnTunnels); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_vpn_tunnels %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertNetworkConnectivitySpokeBetaLinkedInterconnectAttachmentsToHCL(r *networkconnectivityBeta.SpokeLinkedInterconnectAttachments) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	if r.Uris != nil {
		outputConfig += "\turis = ["
		for _, v := range r.Uris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesToHCL(r *networkconnectivityBeta.SpokeLinkedRouterApplianceInstances) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Instances != nil {
		for _, v := range r.Instances {
			outputConfig += fmt.Sprintf("\tinstances %s\n", convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesInstancesToHCL(&v))
		}
	}
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesInstancesToHCL(r *networkconnectivityBeta.SpokeLinkedRouterApplianceInstancesInstances) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.VirtualMachine != nil {
		outputConfig += fmt.Sprintf("\tvirtual_machine = %#v\n", *r.VirtualMachine)
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeBetaLinkedVpnTunnelsToHCL(r *networkconnectivityBeta.SpokeLinkedVpnTunnels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	if r.Uris != nil {
		outputConfig += "\turis = ["
		for _, v := range r.Uris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// OrgPolicyPolicyBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OrgPolicyPolicyBetaAsHCL(r orgpolicyBeta.Policy, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_org_policy_policy\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if v := convertOrgPolicyPolicyBetaSpecToHCL(r.Spec); v != "" {
		outputConfig += fmt.Sprintf("\tspec %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertOrgPolicyPolicyBetaSpecToHCL(r *orgpolicyBeta.PolicySpec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.InheritFromParent != nil {
		outputConfig += fmt.Sprintf("\tinherit_from_parent = %#v\n", *r.InheritFromParent)
	}
	if r.Reset != nil {
		outputConfig += fmt.Sprintf("\treset = %#v\n", *r.Reset)
	}
	if r.Rules != nil {
		for _, v := range r.Rules {
			outputConfig += fmt.Sprintf("\trules %s\n", convertOrgPolicyPolicyBetaSpecRulesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicyBetaSpecRulesToHCL(r *orgpolicyBeta.PolicySpecRules) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAll != nil {
		outputConfig += fmt.Sprintf("\tallow_all = %q\n", serializeEnumBool(r.AllowAll))
	}
	if v := convertOrgPolicyPolicyBetaSpecRulesConditionToHCL(r.Condition); v != "" {
		outputConfig += fmt.Sprintf("\tcondition %s\n", v)
	}
	if r.DenyAll != nil {
		outputConfig += fmt.Sprintf("\tdeny_all = %q\n", serializeEnumBool(r.DenyAll))
	}
	if r.Enforce != nil {
		outputConfig += fmt.Sprintf("\tenforce = %q\n", serializeEnumBool(r.Enforce))
	}
	if v := convertOrgPolicyPolicyBetaSpecRulesValuesToHCL(r.Values); v != "" {
		outputConfig += fmt.Sprintf("\tvalues %s\n", v)
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicyBetaSpecRulesConditionToHCL(r *orgpolicyBeta.PolicySpecRulesCondition) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Expression != nil {
		outputConfig += fmt.Sprintf("\texpression = %#v\n", *r.Expression)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Title != nil {
		outputConfig += fmt.Sprintf("\ttitle = %#v\n", *r.Title)
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicyBetaSpecRulesValuesToHCL(r *orgpolicyBeta.PolicySpecRulesValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowedValues != nil {
		outputConfig += "\tallowed_values = ["
		for _, v := range r.AllowedValues {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.DeniedValues != nil {
		outputConfig += "\tdenied_values = ["
		for _, v := range r.DeniedValues {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// OsConfigOsPolicyAssignmentBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OsConfigOsPolicyAssignmentBetaAsHCL(r osconfigBeta.OSPolicyAssignment, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_os_config_os_policy_assignment\" \"output\" {\n"
	if v := convertOsConfigOsPolicyAssignmentBetaInstanceFilterToHCL(r.InstanceFilter); v != "" {
		outputConfig += fmt.Sprintf("\tinstance_filter %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.OSPolicies != nil {
		for _, v := range r.OSPolicies {
			outputConfig += fmt.Sprintf("\tos_policies %s\n", convertOsConfigOsPolicyAssignmentBetaOSPoliciesToHCL(&v))
		}
	}
	if v := convertOsConfigOsPolicyAssignmentBetaRolloutToHCL(r.Rollout); v != "" {
		outputConfig += fmt.Sprintf("\trollout %s\n", v)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterToHCL(r *osconfigBeta.OSPolicyAssignmentInstanceFilter) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.All != nil {
		outputConfig += fmt.Sprintf("\tall = %#v\n", *r.All)
	}
	if r.ExclusionLabels != nil {
		for _, v := range r.ExclusionLabels {
			outputConfig += fmt.Sprintf("\texclusion_labels %s\n", convertOsConfigOsPolicyAssignmentBetaInstanceFilterExclusionLabelsToHCL(&v))
		}
	}
	if r.InclusionLabels != nil {
		for _, v := range r.InclusionLabels {
			outputConfig += fmt.Sprintf("\tinclusion_labels %s\n", convertOsConfigOsPolicyAssignmentBetaInstanceFilterInclusionLabelsToHCL(&v))
		}
	}
	if r.Inventories != nil {
		for _, v := range r.Inventories {
			outputConfig += fmt.Sprintf("\tinventories %s\n", convertOsConfigOsPolicyAssignmentBetaInstanceFilterInventoriesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterExclusionLabelsToHCL(r *osconfigBeta.OSPolicyAssignmentInstanceFilterExclusionLabels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInclusionLabelsToHCL(r *osconfigBeta.OSPolicyAssignmentInstanceFilterInclusionLabels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInventoriesToHCL(r *osconfigBeta.OSPolicyAssignmentInstanceFilterInventories) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.OSShortName != nil {
		outputConfig += fmt.Sprintf("\tos_short_name = %#v\n", *r.OSShortName)
	}
	if r.OSVersion != nil {
		outputConfig += fmt.Sprintf("\tos_version = %#v\n", *r.OSVersion)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesToHCL(r *osconfigBeta.OSPolicyAssignmentOSPolicies) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.Mode != nil {
		outputConfig += fmt.Sprintf("\tmode = %#v\n", *r.Mode)
	}
	if r.ResourceGroups != nil {
		for _, v := range r.ResourceGroups {
			outputConfig += fmt.Sprintf("\tresource_groups %s\n", convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsToHCL(&v))
		}
	}
	if r.AllowNoResourceGroupMatch != nil {
		outputConfig += fmt.Sprintf("\tallow_no_resource_group_match = %#v\n", *r.AllowNoResourceGroupMatch)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroups) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Resources != nil {
		for _, v := range r.Resources {
			outputConfig += fmt.Sprintf("\tresources %s\n", convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesToHCL(&v))
		}
	}
	if r.InventoryFilters != nil {
		for _, v := range r.InventoryFilters {
			outputConfig += fmt.Sprintf("\tinventory_filters %s\n", convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsInventoryFiltersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResources) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecToHCL(r.Exec); v != "" {
		outputConfig += fmt.Sprintf("\texec %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgToHCL(r.Pkg); v != "" {
		outputConfig += fmt.Sprintf("\tpkg %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryToHCL(r.Repository); v != "" {
		outputConfig += fmt.Sprintf("\trepository %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateToHCL(r.Validate); v != "" {
		outputConfig += fmt.Sprintf("\tvalidate %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceToHCL(r.Enforce); v != "" {
		outputConfig += fmt.Sprintf("\tenforce %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Interpreter != nil {
		outputConfig += fmt.Sprintf("\tinterpreter = %#v\n", *r.Interpreter)
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if r.OutputFilePath != nil {
		outputConfig += fmt.Sprintf("\toutput_file_path = %#v\n", *r.OutputFilePath)
	}
	if r.Script != nil {
		outputConfig += fmt.Sprintf("\tscript = %#v\n", *r.Script)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Interpreter != nil {
		outputConfig += fmt.Sprintf("\tinterpreter = %#v\n", *r.Interpreter)
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if r.OutputFilePath != nil {
		outputConfig += fmt.Sprintf("\toutput_file_path = %#v\n", *r.OutputFilePath)
	}
	if r.Script != nil {
		outputConfig += fmt.Sprintf("\tscript = %#v\n", *r.Script)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.State != nil {
		outputConfig += fmt.Sprintf("\tstate = %#v\n", *r.State)
	}
	if r.Content != nil {
		outputConfig += fmt.Sprintf("\tcontent = %#v\n", *r.Content)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DesiredState != nil {
		outputConfig += fmt.Sprintf("\tdesired_state = %#v\n", *r.DesiredState)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgAptToHCL(r.Apt); v != "" {
		outputConfig += fmt.Sprintf("\tapt %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebToHCL(r.Deb); v != "" {
		outputConfig += fmt.Sprintf("\tdeb %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGoogetToHCL(r.Googet); v != "" {
		outputConfig += fmt.Sprintf("\tgooget %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiToHCL(r.Msi); v != "" {
		outputConfig += fmt.Sprintf("\tmsi %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmToHCL(r.Rpm); v != "" {
		outputConfig += fmt.Sprintf("\trpm %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYumToHCL(r.Yum); v != "" {
		outputConfig += fmt.Sprintf("\tyum %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypperToHCL(r.Zypper); v != "" {
		outputConfig += fmt.Sprintf("\tzypper %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgAptToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.PullDeps != nil {
		outputConfig += fmt.Sprintf("\tpull_deps = %#v\n", *r.PullDeps)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGoogetToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.Properties != nil {
		outputConfig += "\tproperties = ["
		for _, v := range r.Properties {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.PullDeps != nil {
		outputConfig += fmt.Sprintf("\tpull_deps = %#v\n", *r.PullDeps)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYumToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypperToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryAptToHCL(r.Apt); v != "" {
		outputConfig += fmt.Sprintf("\tapt %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGooToHCL(r.Goo); v != "" {
		outputConfig += fmt.Sprintf("\tgoo %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYumToHCL(r.Yum); v != "" {
		outputConfig += fmt.Sprintf("\tyum %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypperToHCL(r.Zypper); v != "" {
		outputConfig += fmt.Sprintf("\tzypper %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryAptToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveType != nil {
		outputConfig += fmt.Sprintf("\tarchive_type = %#v\n", *r.ArchiveType)
	}
	if r.Components != nil {
		outputConfig += "\tcomponents = ["
		for _, v := range r.Components {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Distribution != nil {
		outputConfig += fmt.Sprintf("\tdistribution = %#v\n", *r.Distribution)
	}
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.GpgKey != nil {
		outputConfig += fmt.Sprintf("\tgpg_key = %#v\n", *r.GpgKey)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGooToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Url != nil {
		outputConfig += fmt.Sprintf("\turl = %#v\n", *r.Url)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYumToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BaseUrl != nil {
		outputConfig += fmt.Sprintf("\tbase_url = %#v\n", *r.BaseUrl)
	}
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.GpgKeys != nil {
		outputConfig += "\tgpg_keys = ["
		for _, v := range r.GpgKeys {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypperToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BaseUrl != nil {
		outputConfig += fmt.Sprintf("\tbase_url = %#v\n", *r.BaseUrl)
	}
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.GpgKeys != nil {
		outputConfig += "\tgpg_keys = ["
		for _, v := range r.GpgKeys {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsInventoryFiltersToHCL(r *osconfigBeta.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.OSShortName != nil {
		outputConfig += fmt.Sprintf("\tos_short_name = %#v\n", *r.OSShortName)
	}
	if r.OSVersion != nil {
		outputConfig += fmt.Sprintf("\tos_version = %#v\n", *r.OSVersion)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaRolloutToHCL(r *osconfigBeta.OSPolicyAssignmentRollout) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudgetToHCL(r.DisruptionBudget); v != "" {
		outputConfig += fmt.Sprintf("\tdisruption_budget %s\n", v)
	}
	if r.MinWaitDuration != nil {
		outputConfig += fmt.Sprintf("\tmin_wait_duration = %#v\n", *r.MinWaitDuration)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudgetToHCL(r *osconfigBeta.OSPolicyAssignmentRolloutDisruptionBudget) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Fixed != nil {
		outputConfig += fmt.Sprintf("\tfixed = %#v\n", *r.Fixed)
	}
	if r.Percent != nil {
		outputConfig += fmt.Sprintf("\tpercent = %#v\n", *r.Percent)
	}
	return outputConfig + "}"
}

// PrivatecaCertificateTemplateBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func PrivatecaCertificateTemplateBetaAsHCL(r privatecaBeta.CertificateTemplate, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_privateca_certificate_template\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if v := convertPrivatecaCertificateTemplateBetaIdentityConstraintsToHCL(r.IdentityConstraints); v != "" {
		outputConfig += fmt.Sprintf("\tidentity_constraints %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertPrivatecaCertificateTemplateBetaPassthroughExtensionsToHCL(r.PassthroughExtensions); v != "" {
		outputConfig += fmt.Sprintf("\tpassthrough_extensions %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesToHCL(r.PredefinedValues); v != "" {
		outputConfig += fmt.Sprintf("\tpredefined_values %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraintsToHCL(r *privatecaBeta.CertificateTemplateIdentityConstraints) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowSubjectAltNamesPassthrough != nil {
		outputConfig += fmt.Sprintf("\tallow_subject_alt_names_passthrough = %#v\n", *r.AllowSubjectAltNamesPassthrough)
	}
	if r.AllowSubjectPassthrough != nil {
		outputConfig += fmt.Sprintf("\tallow_subject_passthrough = %#v\n", *r.AllowSubjectPassthrough)
	}
	if v := convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpressionToHCL(r.CelExpression); v != "" {
		outputConfig += fmt.Sprintf("\tcel_expression %s\n", v)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpressionToHCL(r *privatecaBeta.CertificateTemplateIdentityConstraintsCelExpression) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Expression != nil {
		outputConfig += fmt.Sprintf("\texpression = %#v\n", *r.Expression)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Title != nil {
		outputConfig += fmt.Sprintf("\ttitle = %#v\n", *r.Title)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensionsToHCL(r *privatecaBeta.CertificateTemplatePassthroughExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdditionalExtensions != nil {
		for _, v := range r.AdditionalExtensions {
			outputConfig += fmt.Sprintf("\tadditional_extensions %s\n", convertPrivatecaCertificateTemplateBetaPassthroughExtensionsAdditionalExtensionsToHCL(&v))
		}
	}
	if r.KnownExtensions != nil {
		outputConfig += "\tknown_extensions = ["
		for _, v := range r.KnownExtensions {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensionsAdditionalExtensionsToHCL(r *privatecaBeta.CertificateTemplatePassthroughExtensionsAdditionalExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesToHCL(r *privatecaBeta.CertificateTemplatePredefinedValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdditionalExtensions != nil {
		for _, v := range r.AdditionalExtensions {
			outputConfig += fmt.Sprintf("\tadditional_extensions %s\n", convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsToHCL(&v))
		}
	}
	if r.AiaOcspServers != nil {
		outputConfig += "\taia_ocsp_servers = ["
		for _, v := range r.AiaOcspServers {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptionsToHCL(r.CaOptions); v != "" {
		outputConfig += fmt.Sprintf("\tca_options %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageToHCL(r.KeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\tkey_usage %s\n", v)
	}
	if r.PolicyIds != nil {
		for _, v := range r.PolicyIds {
			outputConfig += fmt.Sprintf("\tpolicy_ids %s\n", convertPrivatecaCertificateTemplateBetaPredefinedValuesPolicyIdsToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesAdditionalExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectIdToHCL(r.ObjectId); v != "" {
		outputConfig += fmt.Sprintf("\tobject_id %s\n", v)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	if r.Critical != nil {
		outputConfig += fmt.Sprintf("\tcritical = %#v\n", *r.Critical)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectIdToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptionsToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesCaOptions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IsCa != nil {
		outputConfig += fmt.Sprintf("\tis_ca = %#v\n", *r.IsCa)
	}
	if r.MaxIssuerPathLength != nil {
		outputConfig += fmt.Sprintf("\tmax_issuer_path_length = %#v\n", *r.MaxIssuerPathLength)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsageToHCL(r.BaseKeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\tbase_key_usage %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsageToHCL(r.ExtendedKeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\textended_key_usage %s\n", v)
	}
	if r.UnknownExtendedKeyUsages != nil {
		for _, v := range r.UnknownExtendedKeyUsages {
			outputConfig += fmt.Sprintf("\tunknown_extended_key_usages %s\n", convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageUnknownExtendedKeyUsagesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsageToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CertSign != nil {
		outputConfig += fmt.Sprintf("\tcert_sign = %#v\n", *r.CertSign)
	}
	if r.ContentCommitment != nil {
		outputConfig += fmt.Sprintf("\tcontent_commitment = %#v\n", *r.ContentCommitment)
	}
	if r.CrlSign != nil {
		outputConfig += fmt.Sprintf("\tcrl_sign = %#v\n", *r.CrlSign)
	}
	if r.DataEncipherment != nil {
		outputConfig += fmt.Sprintf("\tdata_encipherment = %#v\n", *r.DataEncipherment)
	}
	if r.DecipherOnly != nil {
		outputConfig += fmt.Sprintf("\tdecipher_only = %#v\n", *r.DecipherOnly)
	}
	if r.DigitalSignature != nil {
		outputConfig += fmt.Sprintf("\tdigital_signature = %#v\n", *r.DigitalSignature)
	}
	if r.EncipherOnly != nil {
		outputConfig += fmt.Sprintf("\tencipher_only = %#v\n", *r.EncipherOnly)
	}
	if r.KeyAgreement != nil {
		outputConfig += fmt.Sprintf("\tkey_agreement = %#v\n", *r.KeyAgreement)
	}
	if r.KeyEncipherment != nil {
		outputConfig += fmt.Sprintf("\tkey_encipherment = %#v\n", *r.KeyEncipherment)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsageToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClientAuth != nil {
		outputConfig += fmt.Sprintf("\tclient_auth = %#v\n", *r.ClientAuth)
	}
	if r.CodeSigning != nil {
		outputConfig += fmt.Sprintf("\tcode_signing = %#v\n", *r.CodeSigning)
	}
	if r.EmailProtection != nil {
		outputConfig += fmt.Sprintf("\temail_protection = %#v\n", *r.EmailProtection)
	}
	if r.OcspSigning != nil {
		outputConfig += fmt.Sprintf("\tocsp_signing = %#v\n", *r.OcspSigning)
	}
	if r.ServerAuth != nil {
		outputConfig += fmt.Sprintf("\tserver_auth = %#v\n", *r.ServerAuth)
	}
	if r.TimeStamping != nil {
		outputConfig += fmt.Sprintf("\ttime_stamping = %#v\n", *r.TimeStamping)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageUnknownExtendedKeyUsagesToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesPolicyIdsToHCL(r *privatecaBeta.CertificateTemplatePredefinedValuesPolicyIds) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// RecaptchaEnterpriseKeyBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func RecaptchaEnterpriseKeyBetaAsHCL(r recaptchaenterpriseBeta.Key, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_recaptcha_enterprise_key\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if v := convertRecaptchaEnterpriseKeyBetaAndroidSettingsToHCL(r.AndroidSettings); v != "" {
		outputConfig += fmt.Sprintf("\tandroid_settings %s\n", v)
	}
	if v := convertRecaptchaEnterpriseKeyBetaIosSettingsToHCL(r.IosSettings); v != "" {
		outputConfig += fmt.Sprintf("\tios_settings %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertRecaptchaEnterpriseKeyBetaTestingOptionsToHCL(r.TestingOptions); v != "" {
		outputConfig += fmt.Sprintf("\ttesting_options %s\n", v)
	}
	if v := convertRecaptchaEnterpriseKeyBetaWebSettingsToHCL(r.WebSettings); v != "" {
		outputConfig += fmt.Sprintf("\tweb_settings %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertRecaptchaEnterpriseKeyBetaAndroidSettingsToHCL(r *recaptchaenterpriseBeta.KeyAndroidSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAllPackageNames != nil {
		outputConfig += fmt.Sprintf("\tallow_all_package_names = %#v\n", *r.AllowAllPackageNames)
	}
	if r.AllowedPackageNames != nil {
		outputConfig += "\tallowed_package_names = ["
		for _, v := range r.AllowedPackageNames {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyBetaIosSettingsToHCL(r *recaptchaenterpriseBeta.KeyIosSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAllBundleIds != nil {
		outputConfig += fmt.Sprintf("\tallow_all_bundle_ids = %#v\n", *r.AllowAllBundleIds)
	}
	if r.AllowedBundleIds != nil {
		outputConfig += "\tallowed_bundle_ids = ["
		for _, v := range r.AllowedBundleIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyBetaTestingOptionsToHCL(r *recaptchaenterpriseBeta.KeyTestingOptions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.TestingChallenge != nil {
		outputConfig += fmt.Sprintf("\ttesting_challenge = %#v\n", *r.TestingChallenge)
	}
	if r.TestingScore != nil {
		outputConfig += fmt.Sprintf("\ttesting_score = %#v\n", *r.TestingScore)
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyBetaWebSettingsToHCL(r *recaptchaenterpriseBeta.KeyWebSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IntegrationType != nil {
		outputConfig += fmt.Sprintf("\tintegration_type = %#v\n", *r.IntegrationType)
	}
	if r.AllowAllDomains != nil {
		outputConfig += fmt.Sprintf("\tallow_all_domains = %#v\n", *r.AllowAllDomains)
	}
	if r.AllowAmpTraffic != nil {
		outputConfig += fmt.Sprintf("\tallow_amp_traffic = %#v\n", *r.AllowAmpTraffic)
	}
	if r.AllowedDomains != nil {
		outputConfig += "\tallowed_domains = ["
		for _, v := range r.AllowedDomains {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ChallengeSecurityPreference != nil {
		outputConfig += fmt.Sprintf("\tchallenge_security_preference = %#v\n", *r.ChallengeSecurityPreference)
	}
	return outputConfig + "}"
}

// AssuredWorkloadsWorkloadAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func AssuredWorkloadsWorkloadAsHCL(r assuredworkloads.Workload, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_assured_workloads_workload\" \"output\" {\n"
	if r.BillingAccount != nil {
		outputConfig += fmt.Sprintf("\tbilling_account = %#v\n", *r.BillingAccount)
	}
	if r.ComplianceRegime != nil {
		outputConfig += fmt.Sprintf("\tcompliance_regime = %#v\n", *r.ComplianceRegime)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Organization != nil {
		outputConfig += fmt.Sprintf("\torganization = %#v\n", *r.Organization)
	}
	if v := convertAssuredWorkloadsWorkloadKmsSettingsToHCL(r.KmsSettings); v != "" {
		outputConfig += fmt.Sprintf("\tkms_settings %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.ProvisionedResourcesParent != nil {
		outputConfig += fmt.Sprintf("\tprovisioned_resources_parent = %#v\n", *r.ProvisionedResourcesParent)
	}
	if r.ResourceSettings != nil {
		for _, v := range r.ResourceSettings {
			outputConfig += fmt.Sprintf("\tresource_settings %s\n", convertAssuredWorkloadsWorkloadResourceSettingsToHCL(&v))
		}
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertAssuredWorkloadsWorkloadKmsSettingsToHCL(r *assuredworkloads.WorkloadKmsSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.NextRotationTime != nil {
		outputConfig += fmt.Sprintf("\tnext_rotation_time = %#v\n", *r.NextRotationTime)
	}
	if r.RotationPeriod != nil {
		outputConfig += fmt.Sprintf("\trotation_period = %#v\n", *r.RotationPeriod)
	}
	return outputConfig + "}"
}

func convertAssuredWorkloadsWorkloadResourceSettingsToHCL(r *assuredworkloads.WorkloadResourceSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ResourceId != nil {
		outputConfig += fmt.Sprintf("\tresource_id = %#v\n", *r.ResourceId)
	}
	if r.ResourceType != nil {
		outputConfig += fmt.Sprintf("\tresource_type = %#v\n", *r.ResourceType)
	}
	return outputConfig + "}"
}

func convertAssuredWorkloadsWorkloadResourcesToHCL(r *assuredworkloads.WorkloadResources) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// CloudbuildWorkerPoolAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudbuildWorkerPoolAsHCL(r cloudbuild.WorkerPool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_cloudbuild_worker_pool\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if v := convertCloudbuildWorkerPoolNetworkConfigToHCL(r.NetworkConfig); v != "" {
		outputConfig += fmt.Sprintf("\tnetwork_config %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertCloudbuildWorkerPoolWorkerConfigToHCL(r.WorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tworker_config %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertCloudbuildWorkerPoolNetworkConfigToHCL(r *cloudbuild.WorkerPoolNetworkConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PeeredNetwork != nil {
		outputConfig += fmt.Sprintf("\tpeered_network = %#v\n", *r.PeeredNetwork)
	}
	return outputConfig + "}"
}

func convertCloudbuildWorkerPoolWorkerConfigToHCL(r *cloudbuild.WorkerPoolWorkerConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DiskSizeGb != nil {
		outputConfig += fmt.Sprintf("\tdisk_size_gb = %#v\n", *r.DiskSizeGb)
	}
	if r.MachineType != nil {
		outputConfig += fmt.Sprintf("\tmachine_type = %#v\n", *r.MachineType)
	}
	if r.NoExternalIP != nil {
		outputConfig += fmt.Sprintf("\tno_external_ip = %#v\n", *r.NoExternalIP)
	}
	return outputConfig + "}"
}

// CloudResourceManagerFolderAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerFolderAsHCL(r cloudresourcemanager.Folder, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_folder\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// CloudResourceManagerProjectAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerProjectAsHCL(r cloudresourcemanager.Project, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_project\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplayname = %#v\n", *r.DisplayName)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAsHCL(r compute.FirewallPolicy, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.ShortName != nil {
		outputConfig += fmt.Sprintf("\tshort_name = %#v\n", *r.ShortName)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyAssociationAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAssociationAsHCL(r compute.FirewallPolicyAssociation, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy_association\" \"output\" {\n"
	if r.AttachmentTarget != nil {
		outputConfig += fmt.Sprintf("\tattachment_target = %#v\n", *r.AttachmentTarget)
	}
	if r.FirewallPolicy != nil {
		outputConfig += fmt.Sprintf("\tfirewall_policy = %#v\n", *r.FirewallPolicy)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeFirewallPolicyRuleAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyRuleAsHCL(r compute.FirewallPolicyRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_firewall_policy_rule\" \"output\" {\n"
	if r.Action != nil {
		outputConfig += fmt.Sprintf("\taction = %#v\n", *r.Action)
	}
	if r.Direction != nil {
		outputConfig += fmt.Sprintf("\tdirection = %#v\n", *r.Direction)
	}
	if r.FirewallPolicy != nil {
		outputConfig += fmt.Sprintf("\tfirewall_policy = %#v\n", *r.FirewallPolicy)
	}
	if v := convertComputeFirewallPolicyRuleMatchToHCL(r.Match); v != "" {
		outputConfig += fmt.Sprintf("\tmatch %s\n", v)
	}
	if r.Priority != nil {
		outputConfig += fmt.Sprintf("\tpriority = %#v\n", *r.Priority)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Disabled != nil {
		outputConfig += fmt.Sprintf("\tdisabled = %#v\n", *r.Disabled)
	}
	if r.EnableLogging != nil {
		outputConfig += fmt.Sprintf("\tenable_logging = %#v\n", *r.EnableLogging)
	}
	if r.TargetResources != nil {
		outputConfig += "\ttarget_resources = ["
		for _, v := range r.TargetResources {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.TargetServiceAccounts != nil {
		outputConfig += "\ttarget_service_accounts = ["
		for _, v := range r.TargetServiceAccounts {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertComputeFirewallPolicyRuleMatchToHCL(r *compute.FirewallPolicyRuleMatch) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Layer4Configs != nil {
		for _, v := range r.Layer4Configs {
			outputConfig += fmt.Sprintf("\tlayer4_configs %s\n", convertComputeFirewallPolicyRuleMatchLayer4ConfigsToHCL(&v))
		}
	}
	if r.DestIPRanges != nil {
		outputConfig += "\tdest_ip_ranges = ["
		for _, v := range r.DestIPRanges {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.SrcIPRanges != nil {
		outputConfig += "\tsrc_ip_ranges = ["
		for _, v := range r.SrcIPRanges {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertComputeFirewallPolicyRuleMatchLayer4ConfigsToHCL(r *compute.FirewallPolicyRuleMatchLayer4Configs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.Ports != nil {
		outputConfig += "\tports = ["
		for _, v := range r.Ports {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// ComputeForwardingRuleAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeForwardingRuleAsHCL(r compute.ForwardingRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_forwarding_rule\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.AllPorts != nil {
		outputConfig += fmt.Sprintf("\tall_ports = %#v\n", *r.AllPorts)
	}
	if r.AllowGlobalAccess != nil {
		outputConfig += fmt.Sprintf("\tallow_global_access = %#v\n", *r.AllowGlobalAccess)
	}
	if r.BackendService != nil {
		outputConfig += fmt.Sprintf("\tbackend_service = %#v\n", *r.BackendService)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.IsMirroringCollector != nil {
		outputConfig += fmt.Sprintf("\tis_mirroring_collector = %#v\n", *r.IsMirroringCollector)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.LoadBalancingScheme != nil {
		outputConfig += fmt.Sprintf("\tload_balancing_scheme = %#v\n", *r.LoadBalancingScheme)
	}
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if r.NetworkTier != nil {
		outputConfig += fmt.Sprintf("\tnetwork_tier = %#v\n", *r.NetworkTier)
	}
	if r.PortRange != nil {
		outputConfig += fmt.Sprintf("\tport_range = %#v\n", *r.PortRange)
	}
	if r.Ports != nil {
		outputConfig += "\tports = ["
		for _, v := range r.Ports {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tregion = %#v\n", *r.Location)
	}
	if r.ServiceLabel != nil {
		outputConfig += fmt.Sprintf("\tservice_label = %#v\n", *r.ServiceLabel)
	}
	if r.Subnetwork != nil {
		outputConfig += fmt.Sprintf("\tsubnetwork = %#v\n", *r.Subnetwork)
	}
	if r.Target != nil {
		outputConfig += fmt.Sprintf("\ttarget = %#v\n", *r.Target)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ComputeGlobalForwardingRuleAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeGlobalForwardingRuleAsHCL(r compute.ForwardingRule, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_compute_global_forwarding_rule\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Target != nil {
		outputConfig += fmt.Sprintf("\ttarget = %#v\n", *r.Target)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.IPProtocol != nil {
		outputConfig += fmt.Sprintf("\tip_protocol = %#v\n", *r.IPProtocol)
	}
	if r.IPVersion != nil {
		outputConfig += fmt.Sprintf("\tip_version = %#v\n", *r.IPVersion)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.LoadBalancingScheme != nil {
		outputConfig += fmt.Sprintf("\tload_balancing_scheme = %#v\n", *r.LoadBalancingScheme)
	}
	if r.MetadataFilter != nil {
		for _, v := range r.MetadataFilter {
			outputConfig += fmt.Sprintf("\tmetadata_filters %s\n", convertComputeGlobalForwardingRuleMetadataFilterToHCL(&v))
		}
	}
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if r.PortRange != nil {
		outputConfig += fmt.Sprintf("\tport_range = %#v\n", *r.PortRange)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertComputeGlobalForwardingRuleMetadataFilterToHCL(r *compute.ForwardingRuleMetadataFilter) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.FilterLabel != nil {
		for _, v := range r.FilterLabel {
			outputConfig += fmt.Sprintf("\tfilter_labels %s\n", convertComputeGlobalForwardingRuleMetadataFilterFilterLabelToHCL(&v))
		}
	}
	if r.FilterMatchCriteria != nil {
		outputConfig += fmt.Sprintf("\tfilter_match_criteria = %#v\n", *r.FilterMatchCriteria)
	}
	return outputConfig + "}"
}

func convertComputeGlobalForwardingRuleMetadataFilterFilterLabelToHCL(r *compute.ForwardingRuleMetadataFilterFilterLabel) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

// ContainerAwsClusterAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAwsClusterAsHCL(r containeraws.Cluster, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_aws_cluster\" \"output\" {\n"
	if v := convertContainerAwsClusterAuthorizationToHCL(r.Authorization); v != "" {
		outputConfig += fmt.Sprintf("\tauthorization %s\n", v)
	}
	if r.AwsRegion != nil {
		outputConfig += fmt.Sprintf("\taws_region = %#v\n", *r.AwsRegion)
	}
	if v := convertContainerAwsClusterControlPlaneToHCL(r.ControlPlane); v != "" {
		outputConfig += fmt.Sprintf("\tcontrol_plane %s\n", v)
	}
	if v := convertContainerAwsClusterFleetToHCL(r.Fleet); v != "" {
		outputConfig += fmt.Sprintf("\tfleet %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertContainerAwsClusterNetworkingToHCL(r.Networking); v != "" {
		outputConfig += fmt.Sprintf("\tnetworking %s\n", v)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAwsClusterAuthorizationToHCL(r *containeraws.ClusterAuthorization) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdminUsers != nil {
		for _, v := range r.AdminUsers {
			outputConfig += fmt.Sprintf("\tadmin_users %s\n", convertContainerAwsClusterAuthorizationAdminUsersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterAuthorizationAdminUsersToHCL(r *containeraws.ClusterAuthorizationAdminUsers) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Username != nil {
		outputConfig += fmt.Sprintf("\tusername = %#v\n", *r.Username)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneToHCL(r *containeraws.ClusterControlPlane) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAwsClusterControlPlaneAwsServicesAuthenticationToHCL(r.AwsServicesAuthentication); v != "" {
		outputConfig += fmt.Sprintf("\taws_services_authentication %s\n", v)
	}
	if v := convertContainerAwsClusterControlPlaneConfigEncryptionToHCL(r.ConfigEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_encryption %s\n", v)
	}
	if v := convertContainerAwsClusterControlPlaneDatabaseEncryptionToHCL(r.DatabaseEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tdatabase_encryption %s\n", v)
	}
	if r.IamInstanceProfile != nil {
		outputConfig += fmt.Sprintf("\tiam_instance_profile = %#v\n", *r.IamInstanceProfile)
	}
	if r.SubnetIds != nil {
		outputConfig += "\tsubnet_ids = ["
		for _, v := range r.SubnetIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	if r.InstanceType != nil {
		outputConfig += fmt.Sprintf("\tinstance_type = %#v\n", *r.InstanceType)
	}
	if v := convertContainerAwsClusterControlPlaneMainVolumeToHCL(r.MainVolume); v != "" {
		outputConfig += fmt.Sprintf("\tmain_volume %s\n", v)
	}
	if v := convertContainerAwsClusterControlPlaneProxyConfigToHCL(r.ProxyConfig); v != "" {
		outputConfig += fmt.Sprintf("\tproxy_config %s\n", v)
	}
	if v := convertContainerAwsClusterControlPlaneRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	if r.SecurityGroupIds != nil {
		outputConfig += "\tsecurity_group_ids = ["
		for _, v := range r.SecurityGroupIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertContainerAwsClusterControlPlaneSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneAwsServicesAuthenticationToHCL(r *containeraws.ClusterControlPlaneAwsServicesAuthentication) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.RoleArn != nil {
		outputConfig += fmt.Sprintf("\trole_arn = %#v\n", *r.RoleArn)
	}
	if r.RoleSessionName != nil {
		outputConfig += fmt.Sprintf("\trole_session_name = %#v\n", *r.RoleSessionName)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneConfigEncryptionToHCL(r *containeraws.ClusterControlPlaneConfigEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneDatabaseEncryptionToHCL(r *containeraws.ClusterControlPlaneDatabaseEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneMainVolumeToHCL(r *containeraws.ClusterControlPlaneMainVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneProxyConfigToHCL(r *containeraws.ClusterControlPlaneProxyConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SecretArn != nil {
		outputConfig += fmt.Sprintf("\tsecret_arn = %#v\n", *r.SecretArn)
	}
	if r.SecretVersion != nil {
		outputConfig += fmt.Sprintf("\tsecret_version = %#v\n", *r.SecretVersion)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneRootVolumeToHCL(r *containeraws.ClusterControlPlaneRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterControlPlaneSshConfigToHCL(r *containeraws.ClusterControlPlaneSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Ec2KeyPair != nil {
		outputConfig += fmt.Sprintf("\tec2_key_pair = %#v\n", *r.Ec2KeyPair)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterFleetToHCL(r *containeraws.ClusterFleet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterNetworkingToHCL(r *containeraws.ClusterNetworking) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PodAddressCidrBlocks != nil {
		outputConfig += "\tpod_address_cidr_blocks = ["
		for _, v := range r.PodAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ServiceAddressCidrBlocks != nil {
		outputConfig += "\tservice_address_cidr_blocks = ["
		for _, v := range r.ServiceAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.VPCId != nil {
		outputConfig += fmt.Sprintf("\tvpc_id = %#v\n", *r.VPCId)
	}
	return outputConfig + "}"
}

func convertContainerAwsClusterWorkloadIdentityConfigToHCL(r *containeraws.ClusterWorkloadIdentityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// ContainerAwsNodePoolAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAwsNodePoolAsHCL(r containeraws.NodePool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_aws_node_pool\" \"output\" {\n"
	if v := convertContainerAwsNodePoolAutoscalingToHCL(r.Autoscaling); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling %s\n", v)
	}
	if r.Cluster != nil {
		outputConfig += fmt.Sprintf("\tcluster = %#v\n", *r.Cluster)
	}
	if v := convertContainerAwsNodePoolConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if v := convertContainerAwsNodePoolMaxPodsConstraintToHCL(r.MaxPodsConstraint); v != "" {
		outputConfig += fmt.Sprintf("\tmax_pods_constraint %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAwsNodePoolAutoscalingToHCL(r *containeraws.NodePoolAutoscaling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmax_node_count = %#v\n", *r.MaxNodeCount)
	}
	if r.MinNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmin_node_count = %#v\n", *r.MinNodeCount)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolConfigToHCL(r *containeraws.NodePoolConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAwsNodePoolConfigConfigEncryptionToHCL(r.ConfigEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_encryption %s\n", v)
	}
	if r.IamInstanceProfile != nil {
		outputConfig += fmt.Sprintf("\tiam_instance_profile = %#v\n", *r.IamInstanceProfile)
	}
	if r.InstanceType != nil {
		outputConfig += fmt.Sprintf("\tinstance_type = %#v\n", *r.InstanceType)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertContainerAwsNodePoolConfigRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	if r.SecurityGroupIds != nil {
		outputConfig += "\tsecurity_group_ids = ["
		for _, v := range r.SecurityGroupIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertContainerAwsNodePoolConfigSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Taints != nil {
		for _, v := range r.Taints {
			outputConfig += fmt.Sprintf("\ttaints %s\n", convertContainerAwsNodePoolConfigTaintsToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolConfigConfigEncryptionToHCL(r *containeraws.NodePoolConfigConfigEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolConfigRootVolumeToHCL(r *containeraws.NodePoolConfigRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Iops != nil {
		outputConfig += fmt.Sprintf("\tiops = %#v\n", *r.Iops)
	}
	if r.KmsKeyArn != nil {
		outputConfig += fmt.Sprintf("\tkms_key_arn = %#v\n", *r.KmsKeyArn)
	}
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	if r.VolumeType != nil {
		outputConfig += fmt.Sprintf("\tvolume_type = %#v\n", *r.VolumeType)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolConfigSshConfigToHCL(r *containeraws.NodePoolConfigSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Ec2KeyPair != nil {
		outputConfig += fmt.Sprintf("\tec2_key_pair = %#v\n", *r.Ec2KeyPair)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolConfigTaintsToHCL(r *containeraws.NodePoolConfigTaints) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Effect != nil {
		outputConfig += fmt.Sprintf("\teffect = %#v\n", *r.Effect)
	}
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

func convertContainerAwsNodePoolMaxPodsConstraintToHCL(r *containeraws.NodePoolMaxPodsConstraint) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxPodsPerNode != nil {
		outputConfig += fmt.Sprintf("\tmax_pods_per_node = %#v\n", *r.MaxPodsPerNode)
	}
	return outputConfig + "}"
}

// ContainerAzureClientAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureClientAsHCL(r containerazure.AzureClient, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_client\" \"output\" {\n"
	if r.ApplicationId != nil {
		outputConfig += fmt.Sprintf("\tapplication_id = %#v\n", *r.ApplicationId)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.TenantId != nil {
		outputConfig += fmt.Sprintf("\ttenant_id = %#v\n", *r.TenantId)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// ContainerAzureClusterAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureClusterAsHCL(r containerazure.Cluster, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_cluster\" \"output\" {\n"
	if v := convertContainerAzureClusterAuthorizationToHCL(r.Authorization); v != "" {
		outputConfig += fmt.Sprintf("\tauthorization %s\n", v)
	}
	if r.AzureRegion != nil {
		outputConfig += fmt.Sprintf("\tazure_region = %#v\n", *r.AzureRegion)
	}
	if r.Client != nil {
		outputConfig += fmt.Sprintf("\tclient = %#v\n", *r.Client)
	}
	if v := convertContainerAzureClusterControlPlaneToHCL(r.ControlPlane); v != "" {
		outputConfig += fmt.Sprintf("\tcontrol_plane %s\n", v)
	}
	if v := convertContainerAzureClusterFleetToHCL(r.Fleet); v != "" {
		outputConfig += fmt.Sprintf("\tfleet %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertContainerAzureClusterNetworkingToHCL(r.Networking); v != "" {
		outputConfig += fmt.Sprintf("\tnetworking %s\n", v)
	}
	if r.ResourceGroupId != nil {
		outputConfig += fmt.Sprintf("\tresource_group_id = %#v\n", *r.ResourceGroupId)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAzureClusterAuthorizationToHCL(r *containerazure.ClusterAuthorization) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdminUsers != nil {
		for _, v := range r.AdminUsers {
			outputConfig += fmt.Sprintf("\tadmin_users %s\n", convertContainerAzureClusterAuthorizationAdminUsersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterAuthorizationAdminUsersToHCL(r *containerazure.ClusterAuthorizationAdminUsers) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Username != nil {
		outputConfig += fmt.Sprintf("\tusername = %#v\n", *r.Username)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneToHCL(r *containerazure.ClusterControlPlane) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAzureClusterControlPlaneSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	if v := convertContainerAzureClusterControlPlaneDatabaseEncryptionToHCL(r.DatabaseEncryption); v != "" {
		outputConfig += fmt.Sprintf("\tdatabase_encryption %s\n", v)
	}
	if v := convertContainerAzureClusterControlPlaneMainVolumeToHCL(r.MainVolume); v != "" {
		outputConfig += fmt.Sprintf("\tmain_volume %s\n", v)
	}
	if v := convertContainerAzureClusterControlPlaneProxyConfigToHCL(r.ProxyConfig); v != "" {
		outputConfig += fmt.Sprintf("\tproxy_config %s\n", v)
	}
	if r.ReplicaPlacements != nil {
		for _, v := range r.ReplicaPlacements {
			outputConfig += fmt.Sprintf("\treplica_placements %s\n", convertContainerAzureClusterControlPlaneReplicaPlacementsToHCL(&v))
		}
	}
	if v := convertContainerAzureClusterControlPlaneRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.VmSize != nil {
		outputConfig += fmt.Sprintf("\tvm_size = %#v\n", *r.VmSize)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneSshConfigToHCL(r *containerazure.ClusterControlPlaneSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AuthorizedKey != nil {
		outputConfig += fmt.Sprintf("\tauthorized_key = %#v\n", *r.AuthorizedKey)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneDatabaseEncryptionToHCL(r *containerazure.ClusterControlPlaneDatabaseEncryption) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.KeyId != nil {
		outputConfig += fmt.Sprintf("\tkey_id = %#v\n", *r.KeyId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneMainVolumeToHCL(r *containerazure.ClusterControlPlaneMainVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneProxyConfigToHCL(r *containerazure.ClusterControlPlaneProxyConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ResourceGroupId != nil {
		outputConfig += fmt.Sprintf("\tresource_group_id = %#v\n", *r.ResourceGroupId)
	}
	if r.SecretId != nil {
		outputConfig += fmt.Sprintf("\tsecret_id = %#v\n", *r.SecretId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneReplicaPlacementsToHCL(r *containerazure.ClusterControlPlaneReplicaPlacements) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AzureAvailabilityZone != nil {
		outputConfig += fmt.Sprintf("\tazure_availability_zone = %#v\n", *r.AzureAvailabilityZone)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterControlPlaneRootVolumeToHCL(r *containerazure.ClusterControlPlaneRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterFleetToHCL(r *containerazure.ClusterFleet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterNetworkingToHCL(r *containerazure.ClusterNetworking) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.PodAddressCidrBlocks != nil {
		outputConfig += "\tpod_address_cidr_blocks = ["
		for _, v := range r.PodAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ServiceAddressCidrBlocks != nil {
		outputConfig += "\tservice_address_cidr_blocks = ["
		for _, v := range r.ServiceAddressCidrBlocks {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.VirtualNetworkId != nil {
		outputConfig += fmt.Sprintf("\tvirtual_network_id = %#v\n", *r.VirtualNetworkId)
	}
	return outputConfig + "}"
}

func convertContainerAzureClusterWorkloadIdentityConfigToHCL(r *containerazure.ClusterWorkloadIdentityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// ContainerAzureNodePoolAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ContainerAzureNodePoolAsHCL(r containerazure.NodePool, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_container_azure_node_pool\" \"output\" {\n"
	if v := convertContainerAzureNodePoolAutoscalingToHCL(r.Autoscaling); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling %s\n", v)
	}
	if r.Cluster != nil {
		outputConfig += fmt.Sprintf("\tcluster = %#v\n", *r.Cluster)
	}
	if v := convertContainerAzureNodePoolConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if v := convertContainerAzureNodePoolMaxPodsConstraintToHCL(r.MaxPodsConstraint); v != "" {
		outputConfig += fmt.Sprintf("\tmax_pods_constraint %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.SubnetId != nil {
		outputConfig += fmt.Sprintf("\tsubnet_id = %#v\n", *r.SubnetId)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	outputConfig += "\tannotations = {"
	for k, v := range r.Annotations {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.AzureAvailabilityZone != nil {
		outputConfig += fmt.Sprintf("\tazure_availability_zone = %#v\n", *r.AzureAvailabilityZone)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertContainerAzureNodePoolAutoscalingToHCL(r *containerazure.NodePoolAutoscaling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmax_node_count = %#v\n", *r.MaxNodeCount)
	}
	if r.MinNodeCount != nil {
		outputConfig += fmt.Sprintf("\tmin_node_count = %#v\n", *r.MinNodeCount)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolConfigToHCL(r *containerazure.NodePoolConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertContainerAzureNodePoolConfigSshConfigToHCL(r.SshConfig); v != "" {
		outputConfig += fmt.Sprintf("\tssh_config %s\n", v)
	}
	if v := convertContainerAzureNodePoolConfigRootVolumeToHCL(r.RootVolume); v != "" {
		outputConfig += fmt.Sprintf("\troot_volume %s\n", v)
	}
	outputConfig += "\ttags = {"
	for k, v := range r.Tags {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.VmSize != nil {
		outputConfig += fmt.Sprintf("\tvm_size = %#v\n", *r.VmSize)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolConfigSshConfigToHCL(r *containerazure.NodePoolConfigSshConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AuthorizedKey != nil {
		outputConfig += fmt.Sprintf("\tauthorized_key = %#v\n", *r.AuthorizedKey)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolConfigRootVolumeToHCL(r *containerazure.NodePoolConfigRootVolume) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SizeGib != nil {
		outputConfig += fmt.Sprintf("\tsize_gib = %#v\n", *r.SizeGib)
	}
	return outputConfig + "}"
}

func convertContainerAzureNodePoolMaxPodsConstraintToHCL(r *containerazure.NodePoolMaxPodsConstraint) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxPodsPerNode != nil {
		outputConfig += fmt.Sprintf("\tmax_pods_per_node = %#v\n", *r.MaxPodsPerNode)
	}
	return outputConfig + "}"
}

// DataprocWorkflowTemplateAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func DataprocWorkflowTemplateAsHCL(r dataproc.WorkflowTemplate, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_dataproc_workflow_template\" \"output\" {\n"
	if r.Jobs != nil {
		for _, v := range r.Jobs {
			outputConfig += fmt.Sprintf("\tjobs %s\n", convertDataprocWorkflowTemplateJobsToHCL(&v))
		}
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertDataprocWorkflowTemplatePlacementToHCL(r.Placement); v != "" {
		outputConfig += fmt.Sprintf("\tplacement %s\n", v)
	}
	if r.DagTimeout != nil {
		outputConfig += fmt.Sprintf("\tdag_timeout = %#v\n", *r.DagTimeout)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Parameters != nil {
		for _, v := range r.Parameters {
			outputConfig += fmt.Sprintf("\tparameters %s\n", convertDataprocWorkflowTemplateParametersToHCL(&v))
		}
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.Version != nil {
		outputConfig += fmt.Sprintf("\tversion = %#v\n", *r.Version)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertDataprocWorkflowTemplateJobsToHCL(r *dataproc.WorkflowTemplateJobs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.StepId != nil {
		outputConfig += fmt.Sprintf("\tstep_id = %#v\n", *r.StepId)
	}
	if v := convertDataprocWorkflowTemplateJobsHadoopJobToHCL(r.HadoopJob); v != "" {
		outputConfig += fmt.Sprintf("\thadoop_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsHiveJobToHCL(r.HiveJob); v != "" {
		outputConfig += fmt.Sprintf("\thive_job %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertDataprocWorkflowTemplateJobsPigJobToHCL(r.PigJob); v != "" {
		outputConfig += fmt.Sprintf("\tpig_job %s\n", v)
	}
	if r.PrerequisiteStepIds != nil {
		outputConfig += "\tprerequisite_step_ids = ["
		for _, v := range r.PrerequisiteStepIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsPrestoJobToHCL(r.PrestoJob); v != "" {
		outputConfig += fmt.Sprintf("\tpresto_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsPysparkJobToHCL(r.PysparkJob); v != "" {
		outputConfig += fmt.Sprintf("\tpyspark_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsSchedulingToHCL(r.Scheduling); v != "" {
		outputConfig += fmt.Sprintf("\tscheduling %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsSparkJobToHCL(r.SparkJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsSparkRJobToHCL(r.SparkRJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_r_job %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateJobsSparkSqlJobToHCL(r.SparkSqlJob); v != "" {
		outputConfig += fmt.Sprintf("\tspark_sql_job %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsHadoopJobToHCL(r *dataproc.WorkflowTemplateJobsHadoopJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.MainClass != nil {
		outputConfig += fmt.Sprintf("\tmain_class = %#v\n", *r.MainClass)
	}
	if r.MainJarFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_jar_file_uri = %#v\n", *r.MainJarFileUri)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsHadoopJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsHiveJobToHCL(r *dataproc.WorkflowTemplateJobsHiveJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsHiveJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsHiveJobQueryListToHCL(r *dataproc.WorkflowTemplateJobsHiveJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPigJobToHCL(r *dataproc.WorkflowTemplateJobsPigJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsPigJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsPigJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPigJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsPigJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPigJobQueryListToHCL(r *dataproc.WorkflowTemplateJobsPigJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPrestoJobToHCL(r *dataproc.WorkflowTemplateJobsPrestoJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClientTags != nil {
		outputConfig += "\tclient_tags = ["
		for _, v := range r.ClientTags {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ContinueOnFailure != nil {
		outputConfig += fmt.Sprintf("\tcontinue_on_failure = %#v\n", *r.ContinueOnFailure)
	}
	if v := convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.OutputFormat != nil {
		outputConfig += fmt.Sprintf("\toutput_format = %#v\n", *r.OutputFormat)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsPrestoJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsPrestoJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPrestoJobQueryListToHCL(r *dataproc.WorkflowTemplateJobsPrestoJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPysparkJobToHCL(r *dataproc.WorkflowTemplateJobsPysparkJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MainPythonFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_python_file_uri = %#v\n", *r.MainPythonFileUri)
	}
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.PythonFileUris != nil {
		outputConfig += "\tpython_file_uris = ["
		for _, v := range r.PythonFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsPysparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSchedulingToHCL(r *dataproc.WorkflowTemplateJobsScheduling) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MaxFailuresPerHour != nil {
		outputConfig += fmt.Sprintf("\tmax_failures_per_hour = %#v\n", *r.MaxFailuresPerHour)
	}
	if r.MaxFailuresTotal != nil {
		outputConfig += fmt.Sprintf("\tmax_failures_total = %#v\n", *r.MaxFailuresTotal)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkJobToHCL(r *dataproc.WorkflowTemplateJobsSparkJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsSparkJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	if r.MainClass != nil {
		outputConfig += fmt.Sprintf("\tmain_class = %#v\n", *r.MainClass)
	}
	if r.MainJarFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_jar_file_uri = %#v\n", *r.MainJarFileUri)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkRJobToHCL(r *dataproc.WorkflowTemplateJobsSparkRJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MainRFileUri != nil {
		outputConfig += fmt.Sprintf("\tmain_r_file_uri = %#v\n", *r.MainRFileUri)
	}
	if r.ArchiveUris != nil {
		outputConfig += "\tarchive_uris = ["
		for _, v := range r.ArchiveUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.FileUris != nil {
		outputConfig += "\tfile_uris = ["
		for _, v := range r.FileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkRJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobToHCL(r *dataproc.WorkflowTemplateJobsSparkSqlJob) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.JarFileUris != nil {
		outputConfig += "\tjar_file_uris = ["
		for _, v := range r.JarFileUris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigToHCL(r.LoggingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlogging_config %s\n", v)
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsSparkSqlJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	outputConfig += "\tscript_variables = {"
	for k, v := range r.ScriptVariables {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkSqlJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tdriver_log_levels = {"
	for k, v := range r.DriverLogLevels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobQueryListToHCL(r *dataproc.WorkflowTemplateJobsSparkSqlJobQueryList) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Queries != nil {
		outputConfig += "\tqueries = ["
		for _, v := range r.Queries {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplatePlacementToHCL(r *dataproc.WorkflowTemplatePlacement) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplatePlacementClusterSelectorToHCL(r.ClusterSelector); v != "" {
		outputConfig += fmt.Sprintf("\tcluster_selector %s\n", v)
	}
	if v := convertDataprocWorkflowTemplatePlacementManagedClusterToHCL(r.ManagedCluster); v != "" {
		outputConfig += fmt.Sprintf("\tmanaged_cluster %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplatePlacementClusterSelectorToHCL(r *dataproc.WorkflowTemplatePlacementClusterSelector) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tcluster_labels = {"
	for k, v := range r.ClusterLabels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Zone != nil {
		outputConfig += fmt.Sprintf("\tzone = %#v\n", *r.Zone)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplatePlacementManagedClusterToHCL(r *dataproc.WorkflowTemplatePlacementManagedCluster) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClusterName != nil {
		outputConfig += fmt.Sprintf("\tcluster_name = %#v\n", *r.ClusterName)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigToHCL(r.Config); v != "" {
		outputConfig += fmt.Sprintf("\tconfig %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateParametersToHCL(r *dataproc.WorkflowTemplateParameters) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Fields != nil {
		outputConfig += "\tfields = ["
		for _, v := range r.Fields {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if v := convertDataprocWorkflowTemplateParametersValidationToHCL(r.Validation); v != "" {
		outputConfig += fmt.Sprintf("\tvalidation %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateParametersValidationToHCL(r *dataproc.WorkflowTemplateParametersValidation) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateParametersValidationRegexToHCL(r.Regex); v != "" {
		outputConfig += fmt.Sprintf("\tregex %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateParametersValidationValuesToHCL(r.Values); v != "" {
		outputConfig += fmt.Sprintf("\tvalues %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateParametersValidationRegexToHCL(r *dataproc.WorkflowTemplateParametersValidationRegex) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Regexes != nil {
		outputConfig += "\tregexes = ["
		for _, v := range r.Regexes {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateParametersValidationValuesToHCL(r *dataproc.WorkflowTemplateParametersValidationValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Values != nil {
		outputConfig += "\tvalues = ["
		for _, v := range r.Values {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigToHCL(r *dataproc.ClusterInstanceGroupConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Accelerators != nil {
		for _, v := range r.Accelerators {
			outputConfig += fmt.Sprintf("\taccelerators %s\n", convertDataprocWorkflowTemplateClusterInstanceGroupConfigAcceleratorsToHCL(&v))
		}
	}
	if v := convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfigToHCL(r.DiskConfig); v != "" {
		outputConfig += fmt.Sprintf("\tdisk_config %s\n", v)
	}
	if r.Image != nil {
		outputConfig += fmt.Sprintf("\timage = %#v\n", *r.Image)
	}
	if r.MachineType != nil {
		outputConfig += fmt.Sprintf("\tmachine_type = %#v\n", *r.MachineType)
	}
	if r.MinCpuPlatform != nil {
		outputConfig += fmt.Sprintf("\tmin_cpu_platform = %#v\n", *r.MinCpuPlatform)
	}
	if r.NumInstances != nil {
		outputConfig += fmt.Sprintf("\tnum_instances = %#v\n", *r.NumInstances)
	}
	if r.Preemptibility != nil {
		outputConfig += fmt.Sprintf("\tpreemptibility = %#v\n", *r.Preemptibility)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigAcceleratorsToHCL(r *dataproc.ClusterInstanceGroupConfigAccelerators) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AcceleratorCount != nil {
		outputConfig += fmt.Sprintf("\taccelerator_count = %#v\n", *r.AcceleratorCount)
	}
	if r.AcceleratorType != nil {
		outputConfig += fmt.Sprintf("\taccelerator_type = %#v\n", *r.AcceleratorType)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfigToHCL(r *dataproc.ClusterInstanceGroupConfigDiskConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BootDiskSizeGb != nil {
		outputConfig += fmt.Sprintf("\tboot_disk_size_gb = %#v\n", *r.BootDiskSizeGb)
	}
	if r.BootDiskType != nil {
		outputConfig += fmt.Sprintf("\tboot_disk_type = %#v\n", *r.BootDiskType)
	}
	if r.NumLocalSsds != nil {
		outputConfig += fmt.Sprintf("\tnum_local_ssds = %#v\n", *r.NumLocalSsds)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigManagedGroupConfigToHCL(r *dataproc.ClusterInstanceGroupConfigManagedGroupConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigToHCL(r *dataproc.ClusterClusterConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfigToHCL(r.AutoscalingConfig); v != "" {
		outputConfig += fmt.Sprintf("\tautoscaling_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfigToHCL(r.EncryptionConfig); v != "" {
		outputConfig += fmt.Sprintf("\tencryption_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfigToHCL(r.EndpointConfig); v != "" {
		outputConfig += fmt.Sprintf("\tendpoint_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigToHCL(r.GceClusterConfig); v != "" {
		outputConfig += fmt.Sprintf("\tgce_cluster_config %s\n", v)
	}
	if r.InitializationActions != nil {
		for _, v := range r.InitializationActions {
			outputConfig += fmt.Sprintf("\tinitialization_actions %s\n", convertDataprocWorkflowTemplateClusterClusterConfigInitializationActionsToHCL(&v))
		}
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfigToHCL(r.LifecycleConfig); v != "" {
		outputConfig += fmt.Sprintf("\tlifecycle_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterInstanceGroupConfigToHCL(r.MasterConfig); v != "" {
		outputConfig += fmt.Sprintf("\tmaster_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterInstanceGroupConfigToHCL(r.SecondaryWorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsecondary_worker_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigToHCL(r.SecurityConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsecurity_config %s\n", v)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfigToHCL(r.SoftwareConfig); v != "" {
		outputConfig += fmt.Sprintf("\tsoftware_config %s\n", v)
	}
	if r.StagingBucket != nil {
		outputConfig += fmt.Sprintf("\tstaging_bucket = %#v\n", *r.StagingBucket)
	}
	if r.TempBucket != nil {
		outputConfig += fmt.Sprintf("\ttemp_bucket = %#v\n", *r.TempBucket)
	}
	if v := convertDataprocWorkflowTemplateClusterInstanceGroupConfigToHCL(r.WorkerConfig); v != "" {
		outputConfig += fmt.Sprintf("\tworker_config %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfigToHCL(r *dataproc.ClusterClusterConfigAutoscalingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Policy != nil {
		outputConfig += fmt.Sprintf("\tpolicy = %#v\n", *r.Policy)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfigToHCL(r *dataproc.ClusterClusterConfigEncryptionConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.GcePdKmsKeyName != nil {
		outputConfig += fmt.Sprintf("\tgce_pd_kms_key_name = %#v\n", *r.GcePdKmsKeyName)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfigToHCL(r *dataproc.ClusterClusterConfigEndpointConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.EnableHttpPortAccess != nil {
		outputConfig += fmt.Sprintf("\tenable_http_port_access = %#v\n", *r.EnableHttpPortAccess)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigToHCL(r *dataproc.ClusterClusterConfigGceClusterConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.InternalIPOnly != nil {
		outputConfig += fmt.Sprintf("\tinternal_ip_only = %#v\n", *r.InternalIPOnly)
	}
	outputConfig += "\tmetadata = {"
	for k, v := range r.Metadata {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Network != nil {
		outputConfig += fmt.Sprintf("\tnetwork = %#v\n", *r.Network)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinityToHCL(r.NodeGroupAffinity); v != "" {
		outputConfig += fmt.Sprintf("\tnode_group_affinity %s\n", v)
	}
	if r.PrivateIPv6GoogleAccess != nil {
		outputConfig += fmt.Sprintf("\tprivate_ipv6_google_access = %#v\n", *r.PrivateIPv6GoogleAccess)
	}
	if v := convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinityToHCL(r.ReservationAffinity); v != "" {
		outputConfig += fmt.Sprintf("\treservation_affinity %s\n", v)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if r.ServiceAccountScopes != nil {
		outputConfig += "\tservice_account_scopes = ["
		for _, v := range r.ServiceAccountScopes {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Subnetwork != nil {
		outputConfig += fmt.Sprintf("\tsubnetwork = %#v\n", *r.Subnetwork)
	}
	if r.Tags != nil {
		outputConfig += "\ttags = ["
		for _, v := range r.Tags {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Zone != nil {
		outputConfig += fmt.Sprintf("\tzone = %#v\n", *r.Zone)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinityToHCL(r *dataproc.ClusterClusterConfigGceClusterConfigNodeGroupAffinity) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.NodeGroup != nil {
		outputConfig += fmt.Sprintf("\tnode_group = %#v\n", *r.NodeGroup)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinityToHCL(r *dataproc.ClusterClusterConfigGceClusterConfigReservationAffinity) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ConsumeReservationType != nil {
		outputConfig += fmt.Sprintf("\tconsume_reservation_type = %#v\n", *r.ConsumeReservationType)
	}
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Values != nil {
		outputConfig += "\tvalues = ["
		for _, v := range r.Values {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigInitializationActionsToHCL(r *dataproc.ClusterClusterConfigInitializationActions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ExecutableFile != nil {
		outputConfig += fmt.Sprintf("\texecutable_file = %#v\n", *r.ExecutableFile)
	}
	if r.ExecutionTimeout != nil {
		outputConfig += fmt.Sprintf("\texecution_timeout = %#v\n", *r.ExecutionTimeout)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfigToHCL(r *dataproc.ClusterClusterConfigLifecycleConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AutoDeleteTime != nil {
		outputConfig += fmt.Sprintf("\tauto_delete_time = %#v\n", *r.AutoDeleteTime)
	}
	if r.AutoDeleteTtl != nil {
		outputConfig += fmt.Sprintf("\tauto_delete_ttl = %#v\n", *r.AutoDeleteTtl)
	}
	if r.IdleDeleteTtl != nil {
		outputConfig += fmt.Sprintf("\tidle_delete_ttl = %#v\n", *r.IdleDeleteTtl)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigToHCL(r *dataproc.ClusterClusterConfigSecurityConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfigToHCL(r.KerberosConfig); v != "" {
		outputConfig += fmt.Sprintf("\tkerberos_config %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfigToHCL(r *dataproc.ClusterClusterConfigSecurityConfigKerberosConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CrossRealmTrustAdminServer != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_admin_server = %#v\n", *r.CrossRealmTrustAdminServer)
	}
	if r.CrossRealmTrustKdc != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_kdc = %#v\n", *r.CrossRealmTrustKdc)
	}
	if r.CrossRealmTrustRealm != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_realm = %#v\n", *r.CrossRealmTrustRealm)
	}
	if r.CrossRealmTrustSharedPassword != nil {
		outputConfig += fmt.Sprintf("\tcross_realm_trust_shared_password = %#v\n", *r.CrossRealmTrustSharedPassword)
	}
	if r.EnableKerberos != nil {
		outputConfig += fmt.Sprintf("\tenable_kerberos = %#v\n", *r.EnableKerberos)
	}
	if r.KdcDbKey != nil {
		outputConfig += fmt.Sprintf("\tkdc_db_key = %#v\n", *r.KdcDbKey)
	}
	if r.KeyPassword != nil {
		outputConfig += fmt.Sprintf("\tkey_password = %#v\n", *r.KeyPassword)
	}
	if r.Keystore != nil {
		outputConfig += fmt.Sprintf("\tkeystore = %#v\n", *r.Keystore)
	}
	if r.KeystorePassword != nil {
		outputConfig += fmt.Sprintf("\tkeystore_password = %#v\n", *r.KeystorePassword)
	}
	if r.KmsKey != nil {
		outputConfig += fmt.Sprintf("\tkms_key = %#v\n", *r.KmsKey)
	}
	if r.Realm != nil {
		outputConfig += fmt.Sprintf("\trealm = %#v\n", *r.Realm)
	}
	if r.RootPrincipalPassword != nil {
		outputConfig += fmt.Sprintf("\troot_principal_password = %#v\n", *r.RootPrincipalPassword)
	}
	if r.TgtLifetimeHours != nil {
		outputConfig += fmt.Sprintf("\ttgt_lifetime_hours = %#v\n", *r.TgtLifetimeHours)
	}
	if r.Truststore != nil {
		outputConfig += fmt.Sprintf("\ttruststore = %#v\n", *r.Truststore)
	}
	if r.TruststorePassword != nil {
		outputConfig += fmt.Sprintf("\ttruststore_password = %#v\n", *r.TruststorePassword)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfigToHCL(r *dataproc.ClusterClusterConfigSoftwareConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ImageVersion != nil {
		outputConfig += fmt.Sprintf("\timage_version = %#v\n", *r.ImageVersion)
	}
	if r.OptionalComponents != nil {
		outputConfig += "\toptional_components = ["
		for _, v := range r.OptionalComponents {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	outputConfig += "\tproperties = {"
	for k, v := range r.Properties {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

// EventarcTriggerAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func EventarcTriggerAsHCL(r eventarc.Trigger, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_eventarc_trigger\" \"output\" {\n"
	if v := convertEventarcTriggerDestinationToHCL(r.Destination); v != "" {
		outputConfig += fmt.Sprintf("\tdestination %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.MatchingCriteria != nil {
		for _, v := range r.MatchingCriteria {
			outputConfig += fmt.Sprintf("\tmatching_criteria %s\n", convertEventarcTriggerMatchingCriteriaToHCL(&v))
		}
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if v := convertEventarcTriggerTransportToHCL(r.Transport); v != "" {
		outputConfig += fmt.Sprintf("\ttransport %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertEventarcTriggerDestinationToHCL(r *eventarc.TriggerDestination) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CloudFunction != nil {
		outputConfig += fmt.Sprintf("\tcloud_function = %#v\n", *r.CloudFunction)
	}
	if v := convertEventarcTriggerDestinationCloudRunServiceToHCL(r.CloudRunService); v != "" {
		outputConfig += fmt.Sprintf("\tcloud_run_service %s\n", v)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerDestinationCloudRunServiceToHCL(r *eventarc.TriggerDestinationCloudRunService) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Service != nil {
		outputConfig += fmt.Sprintf("\tservice = %#v\n", *r.Service)
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.Region != nil {
		outputConfig += fmt.Sprintf("\tregion = %#v\n", *r.Region)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerMatchingCriteriaToHCL(r *eventarc.TriggerMatchingCriteria) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Attribute != nil {
		outputConfig += fmt.Sprintf("\tattribute = %#v\n", *r.Attribute)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerTransportToHCL(r *eventarc.TriggerTransport) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertEventarcTriggerTransportPubsubToHCL(r.Pubsub); v != "" {
		outputConfig += fmt.Sprintf("\tpubsub %s\n", v)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerTransportPubsubToHCL(r *eventarc.TriggerTransportPubsub) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Topic != nil {
		outputConfig += fmt.Sprintf("\ttopic = %#v\n", *r.Topic)
	}
	return outputConfig + "}"
}

// NetworkConnectivityHubAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func NetworkConnectivityHubAsHCL(r networkconnectivity.Hub, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_network_connectivity_hub\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertNetworkConnectivityHubRoutingVpcsToHCL(r *networkconnectivity.HubRoutingVpcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// NetworkConnectivitySpokeAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func NetworkConnectivitySpokeAsHCL(r networkconnectivity.Spoke, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_network_connectivity_spoke\" \"output\" {\n"
	if r.Hub != nil {
		outputConfig += fmt.Sprintf("\thub = %#v\n", *r.Hub)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertNetworkConnectivitySpokeLinkedInterconnectAttachmentsToHCL(r.LinkedInterconnectAttachments); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_interconnect_attachments %s\n", v)
	}
	if v := convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesToHCL(r.LinkedRouterApplianceInstances); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_router_appliance_instances %s\n", v)
	}
	if v := convertNetworkConnectivitySpokeLinkedVpnTunnelsToHCL(r.LinkedVpnTunnels); v != "" {
		outputConfig += fmt.Sprintf("\tlinked_vpn_tunnels %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertNetworkConnectivitySpokeLinkedInterconnectAttachmentsToHCL(r *networkconnectivity.SpokeLinkedInterconnectAttachments) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	if r.Uris != nil {
		outputConfig += "\turis = ["
		for _, v := range r.Uris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesToHCL(r *networkconnectivity.SpokeLinkedRouterApplianceInstances) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Instances != nil {
		for _, v := range r.Instances {
			outputConfig += fmt.Sprintf("\tinstances %s\n", convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesToHCL(&v))
		}
	}
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesToHCL(r *networkconnectivity.SpokeLinkedRouterApplianceInstancesInstances) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IPAddress != nil {
		outputConfig += fmt.Sprintf("\tip_address = %#v\n", *r.IPAddress)
	}
	if r.VirtualMachine != nil {
		outputConfig += fmt.Sprintf("\tvirtual_machine = %#v\n", *r.VirtualMachine)
	}
	return outputConfig + "}"
}

func convertNetworkConnectivitySpokeLinkedVpnTunnelsToHCL(r *networkconnectivity.SpokeLinkedVpnTunnels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.SiteToSiteDataTransfer != nil {
		outputConfig += fmt.Sprintf("\tsite_to_site_data_transfer = %#v\n", *r.SiteToSiteDataTransfer)
	}
	if r.Uris != nil {
		outputConfig += "\turis = ["
		for _, v := range r.Uris {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// OrgPolicyPolicyAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OrgPolicyPolicyAsHCL(r orgpolicy.Policy, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_org_policy_policy\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if v := convertOrgPolicyPolicySpecToHCL(r.Spec); v != "" {
		outputConfig += fmt.Sprintf("\tspec %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertOrgPolicyPolicySpecToHCL(r *orgpolicy.PolicySpec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.InheritFromParent != nil {
		outputConfig += fmt.Sprintf("\tinherit_from_parent = %#v\n", *r.InheritFromParent)
	}
	if r.Reset != nil {
		outputConfig += fmt.Sprintf("\treset = %#v\n", *r.Reset)
	}
	if r.Rules != nil {
		for _, v := range r.Rules {
			outputConfig += fmt.Sprintf("\trules %s\n", convertOrgPolicyPolicySpecRulesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicySpecRulesToHCL(r *orgpolicy.PolicySpecRules) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAll != nil {
		outputConfig += fmt.Sprintf("\tallow_all = %q\n", serializeEnumBool(r.AllowAll))
	}
	if v := convertOrgPolicyPolicySpecRulesConditionToHCL(r.Condition); v != "" {
		outputConfig += fmt.Sprintf("\tcondition %s\n", v)
	}
	if r.DenyAll != nil {
		outputConfig += fmt.Sprintf("\tdeny_all = %q\n", serializeEnumBool(r.DenyAll))
	}
	if r.Enforce != nil {
		outputConfig += fmt.Sprintf("\tenforce = %q\n", serializeEnumBool(r.Enforce))
	}
	if v := convertOrgPolicyPolicySpecRulesValuesToHCL(r.Values); v != "" {
		outputConfig += fmt.Sprintf("\tvalues %s\n", v)
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicySpecRulesConditionToHCL(r *orgpolicy.PolicySpecRulesCondition) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Expression != nil {
		outputConfig += fmt.Sprintf("\texpression = %#v\n", *r.Expression)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Title != nil {
		outputConfig += fmt.Sprintf("\ttitle = %#v\n", *r.Title)
	}
	return outputConfig + "}"
}

func convertOrgPolicyPolicySpecRulesValuesToHCL(r *orgpolicy.PolicySpecRulesValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowedValues != nil {
		outputConfig += "\tallowed_values = ["
		for _, v := range r.AllowedValues {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.DeniedValues != nil {
		outputConfig += "\tdenied_values = ["
		for _, v := range r.DeniedValues {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// OsConfigOsPolicyAssignmentAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OsConfigOsPolicyAssignmentAsHCL(r osconfig.OSPolicyAssignment, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_os_config_os_policy_assignment\" \"output\" {\n"
	if v := convertOsConfigOsPolicyAssignmentInstanceFilterToHCL(r.InstanceFilter); v != "" {
		outputConfig += fmt.Sprintf("\tinstance_filter %s\n", v)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.OSPolicies != nil {
		for _, v := range r.OSPolicies {
			outputConfig += fmt.Sprintf("\tos_policies %s\n", convertOsConfigOsPolicyAssignmentOSPoliciesToHCL(&v))
		}
	}
	if v := convertOsConfigOsPolicyAssignmentRolloutToHCL(r.Rollout); v != "" {
		outputConfig += fmt.Sprintf("\trollout %s\n", v)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertOsConfigOsPolicyAssignmentInstanceFilterToHCL(r *osconfig.OSPolicyAssignmentInstanceFilter) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.All != nil {
		outputConfig += fmt.Sprintf("\tall = %#v\n", *r.All)
	}
	if r.ExclusionLabels != nil {
		for _, v := range r.ExclusionLabels {
			outputConfig += fmt.Sprintf("\texclusion_labels %s\n", convertOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsToHCL(&v))
		}
	}
	if r.InclusionLabels != nil {
		for _, v := range r.InclusionLabels {
			outputConfig += fmt.Sprintf("\tinclusion_labels %s\n", convertOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsToHCL(&v))
		}
	}
	if r.Inventories != nil {
		for _, v := range r.Inventories {
			outputConfig += fmt.Sprintf("\tinventories %s\n", convertOsConfigOsPolicyAssignmentInstanceFilterInventoriesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsToHCL(r *osconfig.OSPolicyAssignmentInstanceFilterExclusionLabels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsToHCL(r *osconfig.OSPolicyAssignmentInstanceFilterInclusionLabels) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInventoriesToHCL(r *osconfig.OSPolicyAssignmentInstanceFilterInventories) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.OSShortName != nil {
		outputConfig += fmt.Sprintf("\tos_short_name = %#v\n", *r.OSShortName)
	}
	if r.OSVersion != nil {
		outputConfig += fmt.Sprintf("\tos_version = %#v\n", *r.OSVersion)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesToHCL(r *osconfig.OSPolicyAssignmentOSPolicies) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.Mode != nil {
		outputConfig += fmt.Sprintf("\tmode = %#v\n", *r.Mode)
	}
	if r.ResourceGroups != nil {
		for _, v := range r.ResourceGroups {
			outputConfig += fmt.Sprintf("\tresource_groups %s\n", convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsToHCL(&v))
		}
	}
	if r.AllowNoResourceGroupMatch != nil {
		outputConfig += fmt.Sprintf("\tallow_no_resource_group_match = %#v\n", *r.AllowNoResourceGroupMatch)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroups) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Resources != nil {
		for _, v := range r.Resources {
			outputConfig += fmt.Sprintf("\tresources %s\n", convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesToHCL(&v))
		}
	}
	if r.InventoryFilters != nil {
		for _, v := range r.InventoryFilters {
			outputConfig += fmt.Sprintf("\tinventory_filters %s\n", convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResources) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecToHCL(r.Exec); v != "" {
		outputConfig += fmt.Sprintf("\texec %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgToHCL(r.Pkg); v != "" {
		outputConfig += fmt.Sprintf("\tpkg %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryToHCL(r.Repository); v != "" {
		outputConfig += fmt.Sprintf("\trepository %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateToHCL(r.Validate); v != "" {
		outputConfig += fmt.Sprintf("\tvalidate %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceToHCL(r.Enforce); v != "" {
		outputConfig += fmt.Sprintf("\tenforce %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Interpreter != nil {
		outputConfig += fmt.Sprintf("\tinterpreter = %#v\n", *r.Interpreter)
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if r.OutputFilePath != nil {
		outputConfig += fmt.Sprintf("\toutput_file_path = %#v\n", *r.OutputFilePath)
	}
	if r.Script != nil {
		outputConfig += fmt.Sprintf("\tscript = %#v\n", *r.Script)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Interpreter != nil {
		outputConfig += fmt.Sprintf("\tinterpreter = %#v\n", *r.Interpreter)
	}
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	if r.OutputFilePath != nil {
		outputConfig += fmt.Sprintf("\toutput_file_path = %#v\n", *r.OutputFilePath)
	}
	if r.Script != nil {
		outputConfig += fmt.Sprintf("\tscript = %#v\n", *r.Script)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.State != nil {
		outputConfig += fmt.Sprintf("\tstate = %#v\n", *r.State)
	}
	if r.Content != nil {
		outputConfig += fmt.Sprintf("\tcontent = %#v\n", *r.Content)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileToHCL(r.File); v != "" {
		outputConfig += fmt.Sprintf("\tfile %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DesiredState != nil {
		outputConfig += fmt.Sprintf("\tdesired_state = %#v\n", *r.DesiredState)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptToHCL(r.Apt); v != "" {
		outputConfig += fmt.Sprintf("\tapt %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebToHCL(r.Deb); v != "" {
		outputConfig += fmt.Sprintf("\tdeb %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetToHCL(r.Googet); v != "" {
		outputConfig += fmt.Sprintf("\tgooget %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiToHCL(r.Msi); v != "" {
		outputConfig += fmt.Sprintf("\tmsi %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmToHCL(r.Rpm); v != "" {
		outputConfig += fmt.Sprintf("\trpm %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumToHCL(r.Yum); v != "" {
		outputConfig += fmt.Sprintf("\tyum %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperToHCL(r.Zypper); v != "" {
		outputConfig += fmt.Sprintf("\tzypper %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.PullDeps != nil {
		outputConfig += fmt.Sprintf("\tpull_deps = %#v\n", *r.PullDeps)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.Properties != nil {
		outputConfig += "\tproperties = ["
		for _, v := range r.Properties {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceToHCL(r.Source); v != "" {
		outputConfig += fmt.Sprintf("\tsource %s\n", v)
	}
	if r.PullDeps != nil {
		outputConfig += fmt.Sprintf("\tpull_deps = %#v\n", *r.PullDeps)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowInsecure != nil {
		outputConfig += fmt.Sprintf("\tallow_insecure = %#v\n", *r.AllowInsecure)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsToHCL(r.Gcs); v != "" {
		outputConfig += fmt.Sprintf("\tgcs %s\n", v)
	}
	if r.LocalPath != nil {
		outputConfig += fmt.Sprintf("\tlocal_path = %#v\n", *r.LocalPath)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteToHCL(r.Remote); v != "" {
		outputConfig += fmt.Sprintf("\tremote %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Bucket != nil {
		outputConfig += fmt.Sprintf("\tbucket = %#v\n", *r.Bucket)
	}
	if r.Object != nil {
		outputConfig += fmt.Sprintf("\tobject = %#v\n", *r.Object)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.Sha256Checksum != nil {
		outputConfig += fmt.Sprintf("\tsha256_checksum = %#v\n", *r.Sha256Checksum)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptToHCL(r.Apt); v != "" {
		outputConfig += fmt.Sprintf("\tapt %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooToHCL(r.Goo); v != "" {
		outputConfig += fmt.Sprintf("\tgoo %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumToHCL(r.Yum); v != "" {
		outputConfig += fmt.Sprintf("\tyum %s\n", v)
	}
	if v := convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperToHCL(r.Zypper); v != "" {
		outputConfig += fmt.Sprintf("\tzypper %s\n", v)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ArchiveType != nil {
		outputConfig += fmt.Sprintf("\tarchive_type = %#v\n", *r.ArchiveType)
	}
	if r.Components != nil {
		outputConfig += "\tcomponents = ["
		for _, v := range r.Components {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Distribution != nil {
		outputConfig += fmt.Sprintf("\tdistribution = %#v\n", *r.Distribution)
	}
	if r.Uri != nil {
		outputConfig += fmt.Sprintf("\turi = %#v\n", *r.Uri)
	}
	if r.GpgKey != nil {
		outputConfig += fmt.Sprintf("\tgpg_key = %#v\n", *r.GpgKey)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Url != nil {
		outputConfig += fmt.Sprintf("\turl = %#v\n", *r.Url)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BaseUrl != nil {
		outputConfig += fmt.Sprintf("\tbase_url = %#v\n", *r.BaseUrl)
	}
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.GpgKeys != nil {
		outputConfig += "\tgpg_keys = ["
		for _, v := range r.GpgKeys {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.BaseUrl != nil {
		outputConfig += fmt.Sprintf("\tbase_url = %#v\n", *r.BaseUrl)
	}
	if r.Id != nil {
		outputConfig += fmt.Sprintf("\tid = %#v\n", *r.Id)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if r.GpgKeys != nil {
		outputConfig += "\tgpg_keys = ["
		for _, v := range r.GpgKeys {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersToHCL(r *osconfig.OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.OSShortName != nil {
		outputConfig += fmt.Sprintf("\tos_short_name = %#v\n", *r.OSShortName)
	}
	if r.OSVersion != nil {
		outputConfig += fmt.Sprintf("\tos_version = %#v\n", *r.OSVersion)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentRolloutToHCL(r *osconfig.OSPolicyAssignmentRollout) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertOsConfigOsPolicyAssignmentRolloutDisruptionBudgetToHCL(r.DisruptionBudget); v != "" {
		outputConfig += fmt.Sprintf("\tdisruption_budget %s\n", v)
	}
	if r.MinWaitDuration != nil {
		outputConfig += fmt.Sprintf("\tmin_wait_duration = %#v\n", *r.MinWaitDuration)
	}
	return outputConfig + "}"
}

func convertOsConfigOsPolicyAssignmentRolloutDisruptionBudgetToHCL(r *osconfig.OSPolicyAssignmentRolloutDisruptionBudget) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Fixed != nil {
		outputConfig += fmt.Sprintf("\tfixed = %#v\n", *r.Fixed)
	}
	if r.Percent != nil {
		outputConfig += fmt.Sprintf("\tpercent = %#v\n", *r.Percent)
	}
	return outputConfig + "}"
}

// PrivatecaCertificateTemplateAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func PrivatecaCertificateTemplateAsHCL(r privateca.CertificateTemplate, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_privateca_certificate_template\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if v := convertPrivatecaCertificateTemplateIdentityConstraintsToHCL(r.IdentityConstraints); v != "" {
		outputConfig += fmt.Sprintf("\tidentity_constraints %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if v := convertPrivatecaCertificateTemplatePassthroughExtensionsToHCL(r.PassthroughExtensions); v != "" {
		outputConfig += fmt.Sprintf("\tpassthrough_extensions %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplatePredefinedValuesToHCL(r.PredefinedValues); v != "" {
		outputConfig += fmt.Sprintf("\tpredefined_values %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertPrivatecaCertificateTemplateIdentityConstraintsToHCL(r *privateca.CertificateTemplateIdentityConstraints) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowSubjectAltNamesPassthrough != nil {
		outputConfig += fmt.Sprintf("\tallow_subject_alt_names_passthrough = %#v\n", *r.AllowSubjectAltNamesPassthrough)
	}
	if r.AllowSubjectPassthrough != nil {
		outputConfig += fmt.Sprintf("\tallow_subject_passthrough = %#v\n", *r.AllowSubjectPassthrough)
	}
	if v := convertPrivatecaCertificateTemplateIdentityConstraintsCelExpressionToHCL(r.CelExpression); v != "" {
		outputConfig += fmt.Sprintf("\tcel_expression %s\n", v)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplateIdentityConstraintsCelExpressionToHCL(r *privateca.CertificateTemplateIdentityConstraintsCelExpression) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Description != nil {
		outputConfig += fmt.Sprintf("\tdescription = %#v\n", *r.Description)
	}
	if r.Expression != nil {
		outputConfig += fmt.Sprintf("\texpression = %#v\n", *r.Expression)
	}
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Title != nil {
		outputConfig += fmt.Sprintf("\ttitle = %#v\n", *r.Title)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePassthroughExtensionsToHCL(r *privateca.CertificateTemplatePassthroughExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdditionalExtensions != nil {
		for _, v := range r.AdditionalExtensions {
			outputConfig += fmt.Sprintf("\tadditional_extensions %s\n", convertPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsToHCL(&v))
		}
	}
	if r.KnownExtensions != nil {
		outputConfig += "\tknown_extensions = ["
		for _, v := range r.KnownExtensions {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsToHCL(r *privateca.CertificateTemplatePassthroughExtensionsAdditionalExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesToHCL(r *privateca.CertificateTemplatePredefinedValues) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AdditionalExtensions != nil {
		for _, v := range r.AdditionalExtensions {
			outputConfig += fmt.Sprintf("\tadditional_extensions %s\n", convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsToHCL(&v))
		}
	}
	if r.AiaOcspServers != nil {
		outputConfig += "\taia_ocsp_servers = ["
		for _, v := range r.AiaOcspServers {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if v := convertPrivatecaCertificateTemplatePredefinedValuesCaOptionsToHCL(r.CaOptions); v != "" {
		outputConfig += fmt.Sprintf("\tca_options %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageToHCL(r.KeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\tkey_usage %s\n", v)
	}
	if r.PolicyIds != nil {
		for _, v := range r.PolicyIds {
			outputConfig += fmt.Sprintf("\tpolicy_ids %s\n", convertPrivatecaCertificateTemplatePredefinedValuesPolicyIdsToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsToHCL(r *privateca.CertificateTemplatePredefinedValuesAdditionalExtensions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdToHCL(r.ObjectId); v != "" {
		outputConfig += fmt.Sprintf("\tobject_id %s\n", v)
	}
	if r.Value != nil {
		outputConfig += fmt.Sprintf("\tvalue = %#v\n", *r.Value)
	}
	if r.Critical != nil {
		outputConfig += fmt.Sprintf("\tcritical = %#v\n", *r.Critical)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdToHCL(r *privateca.CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesCaOptionsToHCL(r *privateca.CertificateTemplatePredefinedValuesCaOptions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IsCa != nil {
		outputConfig += fmt.Sprintf("\tis_ca = %#v\n", *r.IsCa)
	}
	if r.MaxIssuerPathLength != nil {
		outputConfig += fmt.Sprintf("\tmax_issuer_path_length = %#v\n", *r.MaxIssuerPathLength)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageToHCL(r *privateca.CertificateTemplatePredefinedValuesKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageToHCL(r.BaseKeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\tbase_key_usage %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageToHCL(r.ExtendedKeyUsage); v != "" {
		outputConfig += fmt.Sprintf("\textended_key_usage %s\n", v)
	}
	if r.UnknownExtendedKeyUsages != nil {
		for _, v := range r.UnknownExtendedKeyUsages {
			outputConfig += fmt.Sprintf("\tunknown_extended_key_usages %s\n", convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageToHCL(r *privateca.CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.CertSign != nil {
		outputConfig += fmt.Sprintf("\tcert_sign = %#v\n", *r.CertSign)
	}
	if r.ContentCommitment != nil {
		outputConfig += fmt.Sprintf("\tcontent_commitment = %#v\n", *r.ContentCommitment)
	}
	if r.CrlSign != nil {
		outputConfig += fmt.Sprintf("\tcrl_sign = %#v\n", *r.CrlSign)
	}
	if r.DataEncipherment != nil {
		outputConfig += fmt.Sprintf("\tdata_encipherment = %#v\n", *r.DataEncipherment)
	}
	if r.DecipherOnly != nil {
		outputConfig += fmt.Sprintf("\tdecipher_only = %#v\n", *r.DecipherOnly)
	}
	if r.DigitalSignature != nil {
		outputConfig += fmt.Sprintf("\tdigital_signature = %#v\n", *r.DigitalSignature)
	}
	if r.EncipherOnly != nil {
		outputConfig += fmt.Sprintf("\tencipher_only = %#v\n", *r.EncipherOnly)
	}
	if r.KeyAgreement != nil {
		outputConfig += fmt.Sprintf("\tkey_agreement = %#v\n", *r.KeyAgreement)
	}
	if r.KeyEncipherment != nil {
		outputConfig += fmt.Sprintf("\tkey_encipherment = %#v\n", *r.KeyEncipherment)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageToHCL(r *privateca.CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClientAuth != nil {
		outputConfig += fmt.Sprintf("\tclient_auth = %#v\n", *r.ClientAuth)
	}
	if r.CodeSigning != nil {
		outputConfig += fmt.Sprintf("\tcode_signing = %#v\n", *r.CodeSigning)
	}
	if r.EmailProtection != nil {
		outputConfig += fmt.Sprintf("\temail_protection = %#v\n", *r.EmailProtection)
	}
	if r.OcspSigning != nil {
		outputConfig += fmt.Sprintf("\tocsp_signing = %#v\n", *r.OcspSigning)
	}
	if r.ServerAuth != nil {
		outputConfig += fmt.Sprintf("\tserver_auth = %#v\n", *r.ServerAuth)
	}
	if r.TimeStamping != nil {
		outputConfig += fmt.Sprintf("\ttime_stamping = %#v\n", *r.TimeStamping)
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesToHCL(r *privateca.CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertPrivatecaCertificateTemplatePredefinedValuesPolicyIdsToHCL(r *privateca.CertificateTemplatePredefinedValuesPolicyIds) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ObjectIdPath != nil {
		outputConfig += "\tobject_id_path = ["
		for _, v := range r.ObjectIdPath {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

// RecaptchaEnterpriseKeyAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func RecaptchaEnterpriseKeyAsHCL(r recaptchaenterprise.Key, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_recaptcha_enterprise_key\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	if v := convertRecaptchaEnterpriseKeyAndroidSettingsToHCL(r.AndroidSettings); v != "" {
		outputConfig += fmt.Sprintf("\tandroid_settings %s\n", v)
	}
	if v := convertRecaptchaEnterpriseKeyIosSettingsToHCL(r.IosSettings); v != "" {
		outputConfig += fmt.Sprintf("\tios_settings %s\n", v)
	}
	outputConfig += "\tlabels = {"
	for k, v := range r.Labels {
		outputConfig += fmt.Sprintf("%v = %q, ", k, v)
	}
	outputConfig += "}\n"
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertRecaptchaEnterpriseKeyTestingOptionsToHCL(r.TestingOptions); v != "" {
		outputConfig += fmt.Sprintf("\ttesting_options %s\n", v)
	}
	if v := convertRecaptchaEnterpriseKeyWebSettingsToHCL(r.WebSettings); v != "" {
		outputConfig += fmt.Sprintf("\tweb_settings %s\n", v)
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		return "", err
	}
	if !hasGAEquivalent {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

func convertRecaptchaEnterpriseKeyAndroidSettingsToHCL(r *recaptchaenterprise.KeyAndroidSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAllPackageNames != nil {
		outputConfig += fmt.Sprintf("\tallow_all_package_names = %#v\n", *r.AllowAllPackageNames)
	}
	if r.AllowedPackageNames != nil {
		outputConfig += "\tallowed_package_names = ["
		for _, v := range r.AllowedPackageNames {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyIosSettingsToHCL(r *recaptchaenterprise.KeyIosSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.AllowAllBundleIds != nil {
		outputConfig += fmt.Sprintf("\tallow_all_bundle_ids = %#v\n", *r.AllowAllBundleIds)
	}
	if r.AllowedBundleIds != nil {
		outputConfig += "\tallowed_bundle_ids = ["
		for _, v := range r.AllowedBundleIds {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyTestingOptionsToHCL(r *recaptchaenterprise.KeyTestingOptions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.TestingChallenge != nil {
		outputConfig += fmt.Sprintf("\ttesting_challenge = %#v\n", *r.TestingChallenge)
	}
	if r.TestingScore != nil {
		outputConfig += fmt.Sprintf("\ttesting_score = %#v\n", *r.TestingScore)
	}
	return outputConfig + "}"
}

func convertRecaptchaEnterpriseKeyWebSettingsToHCL(r *recaptchaenterprise.KeyWebSettings) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.IntegrationType != nil {
		outputConfig += fmt.Sprintf("\tintegration_type = %#v\n", *r.IntegrationType)
	}
	if r.AllowAllDomains != nil {
		outputConfig += fmt.Sprintf("\tallow_all_domains = %#v\n", *r.AllowAllDomains)
	}
	if r.AllowAmpTraffic != nil {
		outputConfig += fmt.Sprintf("\tallow_amp_traffic = %#v\n", *r.AllowAmpTraffic)
	}
	if r.AllowedDomains != nil {
		outputConfig += "\tallowed_domains = ["
		for _, v := range r.AllowedDomains {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.ChallengeSecurityPreference != nil {
		outputConfig += fmt.Sprintf("\tchallenge_security_preference = %#v\n", *r.ChallengeSecurityPreference)
	}
	return outputConfig + "}"
}

func convertAssuredWorkloadsWorkloadBetaKmsSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nextRotationTime": in["next_rotation_time"],
		"rotationPeriod":   in["rotation_period"],
	}
}

func convertAssuredWorkloadsWorkloadBetaKmsSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadBetaKmsSettings(v))
	}
	return out
}

func convertAssuredWorkloadsWorkloadBetaResourceSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceId":   in["resource_id"],
		"resourceType": in["resource_type"],
	}
}

func convertAssuredWorkloadsWorkloadBetaResourceSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadBetaResourceSettings(v))
	}
	return out
}

func convertAssuredWorkloadsWorkloadBetaResources(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceId":   in["resource_id"],
		"resourceType": in["resource_type"],
	}
}

func convertAssuredWorkloadsWorkloadBetaResourcesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadBetaResources(v))
	}
	return out
}

func convertCloudbuildWorkerPoolBetaNetworkConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"peeredNetwork": in["peered_network"],
	}
}

func convertCloudbuildWorkerPoolBetaNetworkConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertCloudbuildWorkerPoolBetaNetworkConfig(v))
	}
	return out
}

func convertCloudbuildWorkerPoolBetaWorkerConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"diskSizeGb":   in["disk_size_gb"],
		"machineType":  in["machine_type"],
		"noExternalIP": in["no_external_ip"],
	}
}

func convertCloudbuildWorkerPoolBetaWorkerConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertCloudbuildWorkerPoolBetaWorkerConfig(v))
	}
	return out
}

func convertComputeFirewallPolicyRuleBetaMatch(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"layer4Configs": in["layer4_configs"],
		"destIPRanges":  in["dest_ip_ranges"],
		"srcIPRanges":   in["src_ip_ranges"],
	}
}

func convertComputeFirewallPolicyRuleBetaMatchList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeFirewallPolicyRuleBetaMatch(v))
	}
	return out
}

func convertComputeFirewallPolicyRuleBetaMatchLayer4Configs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iPProtocol": in["ip_protocol"],
		"ports":      in["ports"],
	}
}

func convertComputeFirewallPolicyRuleBetaMatchLayer4ConfigsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeFirewallPolicyRuleBetaMatchLayer4Configs(v))
	}
	return out
}

func convertComputeGlobalForwardingRuleBetaMetadataFilter(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"filterLabel":         in["filter_labels"],
		"filterMatchCriteria": in["filter_match_criteria"],
	}
}

func convertComputeGlobalForwardingRuleBetaMetadataFilterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeGlobalForwardingRuleBetaMetadataFilter(v))
	}
	return out
}

func convertComputeGlobalForwardingRuleBetaMetadataFilterFilterLabel(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name":  in["name"],
		"value": in["value"],
	}
}

func convertComputeGlobalForwardingRuleBetaMetadataFilterFilterLabelList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeGlobalForwardingRuleBetaMetadataFilterFilterLabel(v))
	}
	return out
}

func convertContainerAwsClusterBetaAuthorization(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"adminUsers": in["admin_users"],
	}
}

func convertContainerAwsClusterBetaAuthorizationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaAuthorization(v))
	}
	return out
}

func convertContainerAwsClusterBetaAuthorizationAdminUsers(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"username": in["username"],
	}
}

func convertContainerAwsClusterBetaAuthorizationAdminUsersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaAuthorizationAdminUsers(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlane(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"awsServicesAuthentication": convertContainerAwsClusterBetaControlPlaneAwsServicesAuthentication(in["aws_services_authentication"]),
		"configEncryption":          convertContainerAwsClusterBetaControlPlaneConfigEncryption(in["config_encryption"]),
		"databaseEncryption":        convertContainerAwsClusterBetaControlPlaneDatabaseEncryption(in["database_encryption"]),
		"iamInstanceProfile":        in["iam_instance_profile"],
		"subnetIds":                 in["subnet_ids"],
		"version":                   in["version"],
		"instanceType":              in["instance_type"],
		"mainVolume":                convertContainerAwsClusterBetaControlPlaneMainVolume(in["main_volume"]),
		"proxyConfig":               convertContainerAwsClusterBetaControlPlaneProxyConfig(in["proxy_config"]),
		"rootVolume":                convertContainerAwsClusterBetaControlPlaneRootVolume(in["root_volume"]),
		"securityGroupIds":          in["security_group_ids"],
		"sshConfig":                 convertContainerAwsClusterBetaControlPlaneSshConfig(in["ssh_config"]),
		"tags":                      in["tags"],
	}
}

func convertContainerAwsClusterBetaControlPlaneList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlane(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneAwsServicesAuthentication(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"roleArn":         in["role_arn"],
		"roleSessionName": in["role_session_name"],
	}
}

func convertContainerAwsClusterBetaControlPlaneAwsServicesAuthenticationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneAwsServicesAuthentication(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneConfigEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsClusterBetaControlPlaneConfigEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneConfigEncryption(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneDatabaseEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsClusterBetaControlPlaneDatabaseEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneDatabaseEncryption(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneMainVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsClusterBetaControlPlaneMainVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneMainVolume(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneProxyConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"secretArn":     in["secret_arn"],
		"secretVersion": in["secret_version"],
	}
}

func convertContainerAwsClusterBetaControlPlaneProxyConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneProxyConfig(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsClusterBetaControlPlaneRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneRootVolume(v))
	}
	return out
}

func convertContainerAwsClusterBetaControlPlaneSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"ec2KeyPair": in["ec2_key_pair"],
	}
}

func convertContainerAwsClusterBetaControlPlaneSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaControlPlaneSshConfig(v))
	}
	return out
}

func convertContainerAwsClusterBetaFleet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"project":    in["project"],
		"membership": in["membership"],
	}
}

func convertContainerAwsClusterBetaFleetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaFleet(v))
	}
	return out
}

func convertContainerAwsClusterBetaNetworking(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"podAddressCidrBlocks":     in["pod_address_cidr_blocks"],
		"serviceAddressCidrBlocks": in["service_address_cidr_blocks"],
		"vPCId":                    in["vpc_id"],
	}
}

func convertContainerAwsClusterBetaNetworkingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaNetworking(v))
	}
	return out
}

func convertContainerAwsClusterBetaWorkloadIdentityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"identityProvider": in["identity_provider"],
		"issuerUri":        in["issuer_uri"],
		"workloadPool":     in["workload_pool"],
	}
}

func convertContainerAwsClusterBetaWorkloadIdentityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterBetaWorkloadIdentityConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaAutoscaling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxNodeCount": in["max_node_count"],
		"minNodeCount": in["min_node_count"],
	}
}

func convertContainerAwsNodePoolBetaAutoscalingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaAutoscaling(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configEncryption":   convertContainerAwsNodePoolBetaConfigConfigEncryption(in["config_encryption"]),
		"iamInstanceProfile": in["iam_instance_profile"],
		"instanceType":       in["instance_type"],
		"labels":             in["labels"],
		"rootVolume":         convertContainerAwsNodePoolBetaConfigRootVolume(in["root_volume"]),
		"securityGroupIds":   in["security_group_ids"],
		"sshConfig":          convertContainerAwsNodePoolBetaConfigSshConfig(in["ssh_config"]),
		"tags":               in["tags"],
		"taints":             in["taints"],
	}
}

func convertContainerAwsNodePoolBetaConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaConfigConfigEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsNodePoolBetaConfigConfigEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaConfigConfigEncryption(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaConfigRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsNodePoolBetaConfigRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaConfigRootVolume(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaConfigSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"ec2KeyPair": in["ec2_key_pair"],
	}
}

func convertContainerAwsNodePoolBetaConfigSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaConfigSshConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaConfigTaints(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"effect": in["effect"],
		"key":    in["key"],
		"value":  in["value"],
	}
}

func convertContainerAwsNodePoolBetaConfigTaintsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaConfigTaints(v))
	}
	return out
}

func convertContainerAwsNodePoolBetaMaxPodsConstraint(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxPodsPerNode": in["max_pods_per_node"],
	}
}

func convertContainerAwsNodePoolBetaMaxPodsConstraintList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolBetaMaxPodsConstraint(v))
	}
	return out
}

func convertContainerAzureClusterBetaAuthorization(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"adminUsers": in["admin_users"],
	}
}

func convertContainerAzureClusterBetaAuthorizationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaAuthorization(v))
	}
	return out
}

func convertContainerAzureClusterBetaAuthorizationAdminUsers(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"username": in["username"],
	}
}

func convertContainerAzureClusterBetaAuthorizationAdminUsersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaAuthorizationAdminUsers(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlane(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sshConfig":          convertContainerAzureClusterBetaControlPlaneSshConfig(in["ssh_config"]),
		"subnetId":           in["subnet_id"],
		"version":            in["version"],
		"databaseEncryption": convertContainerAzureClusterBetaControlPlaneDatabaseEncryption(in["database_encryption"]),
		"mainVolume":         convertContainerAzureClusterBetaControlPlaneMainVolume(in["main_volume"]),
		"proxyConfig":        convertContainerAzureClusterBetaControlPlaneProxyConfig(in["proxy_config"]),
		"replicaPlacements":  in["replica_placements"],
		"rootVolume":         convertContainerAzureClusterBetaControlPlaneRootVolume(in["root_volume"]),
		"tags":               in["tags"],
		"vmSize":             in["vm_size"],
	}
}

func convertContainerAzureClusterBetaControlPlaneList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlane(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"authorizedKey": in["authorized_key"],
	}
}

func convertContainerAzureClusterBetaControlPlaneSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneSshConfig(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneDatabaseEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"keyId": in["key_id"],
	}
}

func convertContainerAzureClusterBetaControlPlaneDatabaseEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneDatabaseEncryption(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneMainVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureClusterBetaControlPlaneMainVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneMainVolume(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneProxyConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceGroupId": in["resource_group_id"],
		"secretId":        in["secret_id"],
	}
}

func convertContainerAzureClusterBetaControlPlaneProxyConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneProxyConfig(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneReplicaPlacements(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"azureAvailabilityZone": in["azure_availability_zone"],
		"subnetId":              in["subnet_id"],
	}
}

func convertContainerAzureClusterBetaControlPlaneReplicaPlacementsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneReplicaPlacements(v))
	}
	return out
}

func convertContainerAzureClusterBetaControlPlaneRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureClusterBetaControlPlaneRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaControlPlaneRootVolume(v))
	}
	return out
}

func convertContainerAzureClusterBetaFleet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"project":    in["project"],
		"membership": in["membership"],
	}
}

func convertContainerAzureClusterBetaFleetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaFleet(v))
	}
	return out
}

func convertContainerAzureClusterBetaNetworking(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"podAddressCidrBlocks":     in["pod_address_cidr_blocks"],
		"serviceAddressCidrBlocks": in["service_address_cidr_blocks"],
		"virtualNetworkId":         in["virtual_network_id"],
	}
}

func convertContainerAzureClusterBetaNetworkingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaNetworking(v))
	}
	return out
}

func convertContainerAzureClusterBetaWorkloadIdentityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"identityProvider": in["identity_provider"],
		"issuerUri":        in["issuer_uri"],
		"workloadPool":     in["workload_pool"],
	}
}

func convertContainerAzureClusterBetaWorkloadIdentityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterBetaWorkloadIdentityConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolBetaAutoscaling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxNodeCount": in["max_node_count"],
		"minNodeCount": in["min_node_count"],
	}
}

func convertContainerAzureNodePoolBetaAutoscalingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolBetaAutoscaling(v))
	}
	return out
}

func convertContainerAzureNodePoolBetaConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sshConfig":  convertContainerAzureNodePoolBetaConfigSshConfig(in["ssh_config"]),
		"rootVolume": convertContainerAzureNodePoolBetaConfigRootVolume(in["root_volume"]),
		"tags":       in["tags"],
		"vmSize":     in["vm_size"],
	}
}

func convertContainerAzureNodePoolBetaConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolBetaConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolBetaConfigSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"authorizedKey": in["authorized_key"],
	}
}

func convertContainerAzureNodePoolBetaConfigSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolBetaConfigSshConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolBetaConfigRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureNodePoolBetaConfigRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolBetaConfigRootVolume(v))
	}
	return out
}

func convertContainerAzureNodePoolBetaMaxPodsConstraint(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxPodsPerNode": in["max_pods_per_node"],
	}
}

func convertContainerAzureNodePoolBetaMaxPodsConstraintList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolBetaMaxPodsConstraint(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"stepId":              in["step_id"],
		"hadoopJob":           convertDataprocWorkflowTemplateBetaJobsHadoopJob(in["hadoop_job"]),
		"hiveJob":             convertDataprocWorkflowTemplateBetaJobsHiveJob(in["hive_job"]),
		"labels":              in["labels"],
		"pigJob":              convertDataprocWorkflowTemplateBetaJobsPigJob(in["pig_job"]),
		"prerequisiteStepIds": in["prerequisite_step_ids"],
		"prestoJob":           convertDataprocWorkflowTemplateBetaJobsPrestoJob(in["presto_job"]),
		"pysparkJob":          convertDataprocWorkflowTemplateBetaJobsPysparkJob(in["pyspark_job"]),
		"scheduling":          convertDataprocWorkflowTemplateBetaJobsScheduling(in["scheduling"]),
		"sparkJob":            convertDataprocWorkflowTemplateBetaJobsSparkJob(in["spark_job"]),
		"sparkRJob":           convertDataprocWorkflowTemplateBetaJobsSparkRJob(in["spark_r_job"]),
		"sparkSqlJob":         convertDataprocWorkflowTemplateBetaJobsSparkSqlJob(in["spark_sql_job"]),
	}
}

func convertDataprocWorkflowTemplateBetaJobsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobs(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveUris":    in["archive_uris"],
		"args":           in["args"],
		"fileUris":       in["file_uris"],
		"jarFileUris":    in["jar_file_uris"],
		"loggingConfig":  convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfig(in["logging_config"]),
		"mainClass":      in["main_class"],
		"mainJarFileUri": in["main_jar_file_uri"],
		"properties":     in["properties"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsHadoopJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsHiveJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"continueOnFailure": in["continue_on_failure"],
		"jarFileUris":       in["jar_file_uris"],
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateBetaJobsHiveJobQueryList(in["query_list"]),
		"scriptVariables":   in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsHiveJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsHiveJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsHiveJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsHiveJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsHiveJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPigJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"continueOnFailure": in["continue_on_failure"],
		"jarFileUris":       in["jar_file_uris"],
		"loggingConfig":     convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfig(in["logging_config"]),
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateBetaJobsPigJobQueryList(in["query_list"]),
		"scriptVariables":   in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPigJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPigJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPigJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPigJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPigJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clientTags":        in["client_tags"],
		"continueOnFailure": in["continue_on_failure"],
		"loggingConfig":     convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfig(in["logging_config"]),
		"outputFormat":      in["output_format"],
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryList(in["query_list"]),
	}
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPrestoJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPrestoJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPrestoJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"mainPythonFileUri": in["main_python_file_uri"],
		"archiveUris":       in["archive_uris"],
		"args":              in["args"],
		"fileUris":          in["file_uris"],
		"jarFileUris":       in["jar_file_uris"],
		"loggingConfig":     convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfig(in["logging_config"]),
		"properties":        in["properties"],
		"pythonFileUris":    in["python_file_uris"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPysparkJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsPysparkJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsScheduling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxFailuresPerHour": in["max_failures_per_hour"],
		"maxFailuresTotal":   in["max_failures_total"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSchedulingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsScheduling(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveUris":    in["archive_uris"],
		"args":           in["args"],
		"fileUris":       in["file_uris"],
		"jarFileUris":    in["jar_file_uris"],
		"loggingConfig":  convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfig(in["logging_config"]),
		"mainClass":      in["main_class"],
		"mainJarFileUri": in["main_jar_file_uri"],
		"properties":     in["properties"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"mainRFileUri":  in["main_r_file_uri"],
		"archiveUris":   in["archive_uris"],
		"args":          in["args"],
		"fileUris":      in["file_uris"],
		"loggingConfig": convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfig(in["logging_config"]),
		"properties":    in["properties"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkRJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"jarFileUris":     in["jar_file_uris"],
		"loggingConfig":   convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfig(in["logging_config"]),
		"properties":      in["properties"],
		"queryFileUri":    in["query_file_uri"],
		"queryList":       convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryList(in["query_list"]),
		"scriptVariables": in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkSqlJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaPlacement(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterSelector": convertDataprocWorkflowTemplateBetaPlacementClusterSelector(in["cluster_selector"]),
		"managedCluster":  convertDataprocWorkflowTemplateBetaPlacementManagedCluster(in["managed_cluster"]),
	}
}

func convertDataprocWorkflowTemplateBetaPlacementList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaPlacement(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaPlacementClusterSelector(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterLabels": in["cluster_labels"],
		"zone":          in["zone"],
	}
}

func convertDataprocWorkflowTemplateBetaPlacementClusterSelectorList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaPlacementClusterSelector(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaPlacementManagedCluster(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterName": in["cluster_name"],
		"config":      convertDataprocWorkflowTemplateBetaClusterClusterConfig(in["config"]),
		"labels":      in["labels"],
	}
}

func convertDataprocWorkflowTemplateBetaPlacementManagedClusterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaPlacementManagedCluster(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaParameters(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"fields":      in["fields"],
		"name":        in["name"],
		"description": in["description"],
		"validation":  convertDataprocWorkflowTemplateBetaParametersValidation(in["validation"]),
	}
}

func convertDataprocWorkflowTemplateBetaParametersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaParameters(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaParametersValidation(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"regex":  convertDataprocWorkflowTemplateBetaParametersValidationRegex(in["regex"]),
		"values": convertDataprocWorkflowTemplateBetaParametersValidationValues(in["values"]),
	}
}

func convertDataprocWorkflowTemplateBetaParametersValidationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaParametersValidation(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaParametersValidationRegex(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"regexes": in["regexes"],
	}
}

func convertDataprocWorkflowTemplateBetaParametersValidationRegexList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaParametersValidationRegex(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaParametersValidationValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"values": in["values"],
	}
}

func convertDataprocWorkflowTemplateBetaParametersValidationValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaParametersValidationValues(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"accelerators":       in["accelerators"],
		"diskConfig":         convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfig(in["disk_config"]),
		"image":              in["image"],
		"machineType":        in["machine_type"],
		"minCpuPlatform":     in["min_cpu_platform"],
		"numInstances":       in["num_instances"],
		"preemptibility":     in["preemptibility"],
		"instanceNames":      in["instance_names"],
		"isPreemptible":      in["is_preemptible"],
		"managedGroupConfig": convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigManagedGroupConfig(in["managed_group_config"]),
	}
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigAccelerators(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"acceleratorCount": in["accelerator_count"],
		"acceleratorType":  in["accelerator_type"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigAcceleratorsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigAccelerators(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bootDiskSizeGb": in["boot_disk_size_gb"],
		"bootDiskType":   in["boot_disk_type"],
		"numLocalSsds":   in["num_local_ssds"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigDiskConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigManagedGroupConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"instanceGroupManagerName": in["instance_group_manager_name"],
		"instanceTemplateName":     in["instance_template_name"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigManagedGroupConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfigManagedGroupConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"autoscalingConfig":     convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfig(in["autoscaling_config"]),
		"encryptionConfig":      convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfig(in["encryption_config"]),
		"endpointConfig":        convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfig(in["endpoint_config"]),
		"gceClusterConfig":      convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfig(in["gce_cluster_config"]),
		"gkeClusterConfig":      convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfig(in["gke_cluster_config"]),
		"initializationActions": in["initialization_actions"],
		"lifecycleConfig":       convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfig(in["lifecycle_config"]),
		"masterConfig":          convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfig(in["master_config"]),
		"metastoreConfig":       convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfig(in["metastore_config"]),
		"secondaryWorkerConfig": convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfig(in["secondary_worker_config"]),
		"securityConfig":        convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfig(in["security_config"]),
		"softwareConfig":        convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfig(in["software_config"]),
		"stagingBucket":         in["staging_bucket"],
		"tempBucket":            in["temp_bucket"],
		"workerConfig":          convertDataprocWorkflowTemplateBetaClusterInstanceGroupConfig(in["worker_config"]),
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"policy": in["policy"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigAutoscalingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"gcePdKmsKeyName": in["gce_pd_kms_key_name"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigEncryptionConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"enableHttpPortAccess": in["enable_http_port_access"],
		"httpPorts":            in["http_ports"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigEndpointConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"internalIPOnly":          in["internal_ip_only"],
		"metadata":                in["metadata"],
		"network":                 in["network"],
		"nodeGroupAffinity":       convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinity(in["node_group_affinity"]),
		"privateIPv6GoogleAccess": in["private_ipv6_google_access"],
		"reservationAffinity":     convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinity(in["reservation_affinity"]),
		"serviceAccount":          in["service_account"],
		"serviceAccountScopes":    in["service_account_scopes"],
		"subnetwork":              in["subnetwork"],
		"tags":                    in["tags"],
		"zone":                    in["zone"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinity(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nodeGroup": in["node_group"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinityList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigNodeGroupAffinity(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinity(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"consumeReservationType": in["consume_reservation_type"],
		"key":                    in["key"],
		"values":                 in["values"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinityList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigGceClusterConfigReservationAffinity(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"namespacedGkeDeploymentTarget": convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTarget(in["namespaced_gke_deployment_target"]),
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTarget(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterNamespace": in["cluster_namespace"],
		"targetGkeCluster": in["target_gke_cluster"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTargetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigGkeClusterConfigNamespacedGkeDeploymentTarget(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigInitializationActions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"executableFile":   in["executable_file"],
		"executionTimeout": in["execution_timeout"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigInitializationActionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigInitializationActions(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"autoDeleteTime": in["auto_delete_time"],
		"autoDeleteTtl":  in["auto_delete_ttl"],
		"idleDeleteTtl":  in["idle_delete_ttl"],
		"idleStartTime":  in["idle_start_time"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigLifecycleConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"dataprocMetastoreService": in["dataproc_metastore_service"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigMetastoreConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kerberosConfig": convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfig(in["kerberos_config"]),
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"crossRealmTrustAdminServer":    in["cross_realm_trust_admin_server"],
		"crossRealmTrustKdc":            in["cross_realm_trust_kdc"],
		"crossRealmTrustRealm":          in["cross_realm_trust_realm"],
		"crossRealmTrustSharedPassword": in["cross_realm_trust_shared_password"],
		"enableKerberos":                in["enable_kerberos"],
		"kdcDbKey":                      in["kdc_db_key"],
		"keyPassword":                   in["key_password"],
		"keystore":                      in["keystore"],
		"keystorePassword":              in["keystore_password"],
		"kmsKey":                        in["kms_key"],
		"realm":                         in["realm"],
		"rootPrincipalPassword":         in["root_principal_password"],
		"tgtLifetimeHours":              in["tgt_lifetime_hours"],
		"truststore":                    in["truststore"],
		"truststorePassword":            in["truststore_password"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigSecurityConfigKerberosConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"imageVersion":       in["image_version"],
		"optionalComponents": in["optional_components"],
		"properties":         in["properties"],
	}
}

func convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateBetaClusterClusterConfigSoftwareConfig(v))
	}
	return out
}

func convertEventarcTriggerBetaDestination(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"cloudFunction":   in["cloud_function"],
		"cloudRunService": convertEventarcTriggerBetaDestinationCloudRunService(in["cloud_run_service"]),
	}
}

func convertEventarcTriggerBetaDestinationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerBetaDestination(v))
	}
	return out
}

func convertEventarcTriggerBetaDestinationCloudRunService(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"service": in["service"],
		"path":    in["path"],
		"region":  in["region"],
	}
}

func convertEventarcTriggerBetaDestinationCloudRunServiceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerBetaDestinationCloudRunService(v))
	}
	return out
}

func convertEventarcTriggerBetaMatchingCriteria(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"attribute": in["attribute"],
		"value":     in["value"],
	}
}

func convertEventarcTriggerBetaMatchingCriteriaList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerBetaMatchingCriteria(v))
	}
	return out
}

func convertEventarcTriggerBetaTransport(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"pubsub": convertEventarcTriggerBetaTransportPubsub(in["pubsub"]),
	}
}

func convertEventarcTriggerBetaTransportList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerBetaTransport(v))
	}
	return out
}

func convertEventarcTriggerBetaTransportPubsub(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"topic":        in["topic"],
		"subscription": in["subscription"],
	}
}

func convertEventarcTriggerBetaTransportPubsubList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerBetaTransportPubsub(v))
	}
	return out
}

func convertGkeHubFeatureBetaSpec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"multiclusteringress": convertGkeHubFeatureBetaSpecMulticlusteringress(in["multiclusteringress"]),
	}
}

func convertGkeHubFeatureBetaSpecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureBetaSpec(v))
	}
	return out
}

func convertGkeHubFeatureBetaSpecMulticlusteringress(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configMembership": in["config_membership"],
	}
}

func convertGkeHubFeatureBetaSpecMulticlusteringressList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureBetaSpecMulticlusteringress(v))
	}
	return out
}

func convertGkeHubFeatureBetaResourceState(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"hasResources": in["has_resources"],
		"state":        in["state"],
	}
}

func convertGkeHubFeatureBetaResourceStateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureBetaResourceState(v))
	}
	return out
}

func convertGkeHubFeatureBetaState(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"state": convertGkeHubFeatureBetaStateState(in["state"]),
	}
}

func convertGkeHubFeatureBetaStateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureBetaState(v))
	}
	return out
}

func convertGkeHubFeatureBetaStateState(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"code":        in["code"],
		"description": in["description"],
		"updateTime":  in["update_time"],
	}
}

func convertGkeHubFeatureBetaStateStateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureBetaStateState(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagement(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"binauthz":            convertGkeHubFeatureMembershipBetaConfigmanagementBinauthz(in["binauthz"]),
		"configSync":          convertGkeHubFeatureMembershipBetaConfigmanagementConfigSync(in["config_sync"]),
		"hierarchyController": convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyController(in["hierarchy_controller"]),
		"policyController":    convertGkeHubFeatureMembershipBetaConfigmanagementPolicyController(in["policy_controller"]),
		"version":             in["version"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagement(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagementBinauthz(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"enabled": in["enabled"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementBinauthzList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagementBinauthz(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSync(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"git":          convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGit(in["git"]),
		"sourceFormat": in["source_format"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagementConfigSync(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGit(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"gcpServiceAccountEmail": in["gcp_service_account_email"],
		"httpsProxy":             in["https_proxy"],
		"policyDir":              in["policy_dir"],
		"secretType":             in["secret_type"],
		"syncBranch":             in["sync_branch"],
		"syncRepo":               in["sync_repo"],
		"syncRev":                in["sync_rev"],
		"syncWaitSecs":           in["sync_wait_secs"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGitList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagementConfigSyncGit(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyController(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"enableHierarchicalResourceQuota": in["enable_hierarchical_resource_quota"],
		"enablePodTreeLabels":             in["enable_pod_tree_labels"],
		"enabled":                         in["enabled"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyControllerList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagementHierarchyController(v))
	}
	return out
}

func convertGkeHubFeatureMembershipBetaConfigmanagementPolicyController(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"auditIntervalSeconds":     in["audit_interval_seconds"],
		"enabled":                  in["enabled"],
		"exemptableNamespaces":     in["exemptable_namespaces"],
		"logDeniesEnabled":         in["log_denies_enabled"],
		"referentialRulesEnabled":  in["referential_rules_enabled"],
		"templateLibraryInstalled": in["template_library_installed"],
	}
}

func convertGkeHubFeatureMembershipBetaConfigmanagementPolicyControllerList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertGkeHubFeatureMembershipBetaConfigmanagementPolicyController(v))
	}
	return out
}

func convertNetworkConnectivityHubBetaRoutingVpcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri": in["uri"],
	}
}

func convertNetworkConnectivityHubBetaRoutingVpcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivityHubBetaRoutingVpcs(v))
	}
	return out
}

func convertNetworkConnectivitySpokeBetaLinkedInterconnectAttachments(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
		"uris":                   in["uris"],
	}
}

func convertNetworkConnectivitySpokeBetaLinkedInterconnectAttachmentsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeBetaLinkedInterconnectAttachments(v))
	}
	return out
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstances(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"instances":              in["instances"],
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
	}
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstances(v))
	}
	return out
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesInstances(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iPAddress":      in["ip_address"],
		"virtualMachine": in["virtual_machine"],
	}
}

func convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesInstancesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeBetaLinkedRouterApplianceInstancesInstances(v))
	}
	return out
}

func convertNetworkConnectivitySpokeBetaLinkedVpnTunnels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
		"uris":                   in["uris"],
	}
}

func convertNetworkConnectivitySpokeBetaLinkedVpnTunnelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeBetaLinkedVpnTunnels(v))
	}
	return out
}

func convertOrgPolicyPolicyBetaSpec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"inheritFromParent": in["inherit_from_parent"],
		"reset":             in["reset"],
		"rules":             in["rules"],
		"etag":              in["etag"],
		"updateTime":        in["update_time"],
	}
}

func convertOrgPolicyPolicyBetaSpecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicyBetaSpec(v))
	}
	return out
}

func convertOrgPolicyPolicyBetaSpecRules(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAll":  in["allow_all"],
		"condition": convertOrgPolicyPolicyBetaSpecRulesCondition(in["condition"]),
		"denyAll":   in["deny_all"],
		"enforce":   in["enforce"],
		"values":    convertOrgPolicyPolicyBetaSpecRulesValues(in["values"]),
	}
}

func convertOrgPolicyPolicyBetaSpecRulesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicyBetaSpecRules(v))
	}
	return out
}

func convertOrgPolicyPolicyBetaSpecRulesCondition(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"description": in["description"],
		"expression":  in["expression"],
		"location":    in["location"],
		"title":       in["title"],
	}
}

func convertOrgPolicyPolicyBetaSpecRulesConditionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicyBetaSpecRulesCondition(v))
	}
	return out
}

func convertOrgPolicyPolicyBetaSpecRulesValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowedValues": in["allowed_values"],
		"deniedValues":  in["denied_values"],
	}
}

func convertOrgPolicyPolicyBetaSpecRulesValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicyBetaSpecRulesValues(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilter(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"all":             in["all"],
		"exclusionLabels": in["exclusion_labels"],
		"inclusionLabels": in["inclusion_labels"],
		"inventories":     in["inventories"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaInstanceFilter(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterExclusionLabels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"labels": in["labels"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterExclusionLabelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaInstanceFilterExclusionLabels(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInclusionLabels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"labels": in["labels"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInclusionLabelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaInstanceFilterInclusionLabels(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInventories(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"oSShortName": in["os_short_name"],
		"oSVersion":   in["os_version"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaInstanceFilterInventoriesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaInstanceFilterInventories(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPolicies(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"id":                        in["id"],
		"mode":                      in["mode"],
		"resourceGroups":            in["resource_groups"],
		"allowNoResourceGroupMatch": in["allow_no_resource_group_match"],
		"description":               in["description"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPolicies(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroups(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resources":        in["resources"],
		"inventoryFilters": in["inventory_filters"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroups(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResources(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"id":         in["id"],
		"exec":       convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExec(in["exec"]),
		"file":       convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFile(in["file"]),
		"pkg":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkg(in["pkg"]),
		"repository": convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepository(in["repository"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResources(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"validate": convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidate(in["validate"]),
		"enforce":  convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforce(in["enforce"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExec(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidate(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"interpreter":    in["interpreter"],
		"args":           in["args"],
		"file":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFile(in["file"]),
		"outputFilePath": in["output_file_path"],
		"script":         in["script"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidate(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecValidateFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforce(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"interpreter":    in["interpreter"],
		"args":           in["args"],
		"file":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFile(in["file"]),
		"outputFilePath": in["output_file_path"],
		"script":         in["script"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforce(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"path":        in["path"],
		"state":       in["state"],
		"content":     in["content"],
		"file":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFile(in["file"]),
		"permissions": in["permissions"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesFileFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkg(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"desiredState": in["desired_state"],
		"apt":          convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgApt(in["apt"]),
		"deb":          convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDeb(in["deb"]),
		"googet":       convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGooget(in["googet"]),
		"msi":          convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsi(in["msi"]),
		"rpm":          convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpm(in["rpm"]),
		"yum":          convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYum(in["yum"]),
		"zypper":       convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypper(in["zypper"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkg(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgApt(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgAptList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgApt(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDeb(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":   convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSource(in["source"]),
		"pullDeps": in["pull_deps"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDeb(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGooget(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGoogetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgGooget(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsi(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":     convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSource(in["source"]),
		"properties": in["properties"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsi(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpm(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":   convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSource(in["source"]),
		"pullDeps": in["pull_deps"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpm(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYum(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYumList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgYum(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypper(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypperList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesPkgZypper(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepository(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"apt":    convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryApt(in["apt"]),
		"goo":    convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGoo(in["goo"]),
		"yum":    convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYum(in["yum"]),
		"zypper": convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypper(in["zypper"]),
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepository(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryApt(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveType":  in["archive_type"],
		"components":   in["components"],
		"distribution": in["distribution"],
		"uri":          in["uri"],
		"gpgKey":       in["gpg_key"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryAptList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryApt(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGoo(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
		"url":  in["url"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGooList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryGoo(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYum(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseUrl":     in["base_url"],
		"id":          in["id"],
		"displayName": in["display_name"],
		"gpgKeys":     in["gpg_keys"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYumList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryYum(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypper(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseUrl":     in["base_url"],
		"id":          in["id"],
		"displayName": in["display_name"],
		"gpgKeys":     in["gpg_keys"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypperList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsResourcesRepositoryZypper(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsInventoryFilters(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"oSShortName": in["os_short_name"],
		"oSVersion":   in["os_version"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsInventoryFiltersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaOSPoliciesResourceGroupsInventoryFilters(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaRollout(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"disruptionBudget": convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudget(in["disruption_budget"]),
		"minWaitDuration":  in["min_wait_duration"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaRolloutList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaRollout(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudget(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"fixed":   in["fixed"],
		"percent": in["percent"],
	}
}

func convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudgetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentBetaRolloutDisruptionBudget(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraints(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowSubjectAltNamesPassthrough": in["allow_subject_alt_names_passthrough"],
		"allowSubjectPassthrough":         in["allow_subject_passthrough"],
		"celExpression":                   convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpression(in["cel_expression"]),
	}
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraintsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaIdentityConstraints(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpression(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"description": in["description"],
		"expression":  in["expression"],
		"location":    in["location"],
		"title":       in["title"],
	}
}

func convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpressionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaIdentityConstraintsCelExpression(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"additionalExtensions": in["additional_extensions"],
		"knownExtensions":      in["known_extensions"],
	}
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPassthroughExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensionsAdditionalExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplateBetaPassthroughExtensionsAdditionalExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPassthroughExtensionsAdditionalExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"additionalExtensions": in["additional_extensions"],
		"aiaOcspServers":       in["aia_ocsp_servers"],
		"caOptions":            convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptions(in["ca_options"]),
		"keyUsage":             convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsage(in["key_usage"]),
		"policyIds":            in["policy_ids"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValues(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectId": convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectId(in["object_id"]),
		"value":    in["value"],
		"critical": in["critical"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectId(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectIdList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesAdditionalExtensionsObjectId(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"isCa":                in["is_ca"],
		"maxIssuerPathLength": in["max_issuer_path_length"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesCaOptions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseKeyUsage":             convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsage(in["base_key_usage"]),
		"extendedKeyUsage":         convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsage(in["extended_key_usage"]),
		"unknownExtendedKeyUsages": in["unknown_extended_key_usages"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"certSign":          in["cert_sign"],
		"contentCommitment": in["content_commitment"],
		"crlSign":           in["crl_sign"],
		"dataEncipherment":  in["data_encipherment"],
		"decipherOnly":      in["decipher_only"],
		"digitalSignature":  in["digital_signature"],
		"encipherOnly":      in["encipher_only"],
		"keyAgreement":      in["key_agreement"],
		"keyEncipherment":   in["key_encipherment"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageBaseKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clientAuth":      in["client_auth"],
		"codeSigning":     in["code_signing"],
		"emailProtection": in["email_protection"],
		"ocspSigning":     in["ocsp_signing"],
		"serverAuth":      in["server_auth"],
		"timeStamping":    in["time_stamping"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageExtendedKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageUnknownExtendedKeyUsages(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageUnknownExtendedKeyUsagesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesKeyUsageUnknownExtendedKeyUsages(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesPolicyIds(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplateBetaPredefinedValuesPolicyIdsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateBetaPredefinedValuesPolicyIds(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyBetaAndroidSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAllPackageNames": in["allow_all_package_names"],
		"allowedPackageNames":  in["allowed_package_names"],
	}
}

func convertRecaptchaEnterpriseKeyBetaAndroidSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyBetaAndroidSettings(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyBetaIosSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAllBundleIds": in["allow_all_bundle_ids"],
		"allowedBundleIds":  in["allowed_bundle_ids"],
	}
}

func convertRecaptchaEnterpriseKeyBetaIosSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyBetaIosSettings(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyBetaTestingOptions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"testingChallenge": in["testing_challenge"],
		"testingScore":     in["testing_score"],
	}
}

func convertRecaptchaEnterpriseKeyBetaTestingOptionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyBetaTestingOptions(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyBetaWebSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"integrationType":             in["integration_type"],
		"allowAllDomains":             in["allow_all_domains"],
		"allowAmpTraffic":             in["allow_amp_traffic"],
		"allowedDomains":              in["allowed_domains"],
		"challengeSecurityPreference": in["challenge_security_preference"],
	}
}

func convertRecaptchaEnterpriseKeyBetaWebSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyBetaWebSettings(v))
	}
	return out
}

func convertAssuredWorkloadsWorkloadKmsSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nextRotationTime": in["next_rotation_time"],
		"rotationPeriod":   in["rotation_period"],
	}
}

func convertAssuredWorkloadsWorkloadKmsSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadKmsSettings(v))
	}
	return out
}

func convertAssuredWorkloadsWorkloadResourceSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceId":   in["resource_id"],
		"resourceType": in["resource_type"],
	}
}

func convertAssuredWorkloadsWorkloadResourceSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadResourceSettings(v))
	}
	return out
}

func convertAssuredWorkloadsWorkloadResources(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceId":   in["resource_id"],
		"resourceType": in["resource_type"],
	}
}

func convertAssuredWorkloadsWorkloadResourcesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertAssuredWorkloadsWorkloadResources(v))
	}
	return out
}

func convertCloudbuildWorkerPoolNetworkConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"peeredNetwork": in["peered_network"],
	}
}

func convertCloudbuildWorkerPoolNetworkConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertCloudbuildWorkerPoolNetworkConfig(v))
	}
	return out
}

func convertCloudbuildWorkerPoolWorkerConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"diskSizeGb":   in["disk_size_gb"],
		"machineType":  in["machine_type"],
		"noExternalIP": in["no_external_ip"],
	}
}

func convertCloudbuildWorkerPoolWorkerConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertCloudbuildWorkerPoolWorkerConfig(v))
	}
	return out
}

func convertComputeFirewallPolicyRuleMatch(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"layer4Configs": in["layer4_configs"],
		"destIPRanges":  in["dest_ip_ranges"],
		"srcIPRanges":   in["src_ip_ranges"],
	}
}

func convertComputeFirewallPolicyRuleMatchList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeFirewallPolicyRuleMatch(v))
	}
	return out
}

func convertComputeFirewallPolicyRuleMatchLayer4Configs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iPProtocol": in["ip_protocol"],
		"ports":      in["ports"],
	}
}

func convertComputeFirewallPolicyRuleMatchLayer4ConfigsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeFirewallPolicyRuleMatchLayer4Configs(v))
	}
	return out
}

func convertComputeGlobalForwardingRuleMetadataFilter(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"filterLabel":         in["filter_labels"],
		"filterMatchCriteria": in["filter_match_criteria"],
	}
}

func convertComputeGlobalForwardingRuleMetadataFilterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeGlobalForwardingRuleMetadataFilter(v))
	}
	return out
}

func convertComputeGlobalForwardingRuleMetadataFilterFilterLabel(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name":  in["name"],
		"value": in["value"],
	}
}

func convertComputeGlobalForwardingRuleMetadataFilterFilterLabelList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertComputeGlobalForwardingRuleMetadataFilterFilterLabel(v))
	}
	return out
}

func convertContainerAwsClusterAuthorization(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"adminUsers": in["admin_users"],
	}
}

func convertContainerAwsClusterAuthorizationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterAuthorization(v))
	}
	return out
}

func convertContainerAwsClusterAuthorizationAdminUsers(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"username": in["username"],
	}
}

func convertContainerAwsClusterAuthorizationAdminUsersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterAuthorizationAdminUsers(v))
	}
	return out
}

func convertContainerAwsClusterControlPlane(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"awsServicesAuthentication": convertContainerAwsClusterControlPlaneAwsServicesAuthentication(in["aws_services_authentication"]),
		"configEncryption":          convertContainerAwsClusterControlPlaneConfigEncryption(in["config_encryption"]),
		"databaseEncryption":        convertContainerAwsClusterControlPlaneDatabaseEncryption(in["database_encryption"]),
		"iamInstanceProfile":        in["iam_instance_profile"],
		"subnetIds":                 in["subnet_ids"],
		"version":                   in["version"],
		"instanceType":              in["instance_type"],
		"mainVolume":                convertContainerAwsClusterControlPlaneMainVolume(in["main_volume"]),
		"proxyConfig":               convertContainerAwsClusterControlPlaneProxyConfig(in["proxy_config"]),
		"rootVolume":                convertContainerAwsClusterControlPlaneRootVolume(in["root_volume"]),
		"securityGroupIds":          in["security_group_ids"],
		"sshConfig":                 convertContainerAwsClusterControlPlaneSshConfig(in["ssh_config"]),
		"tags":                      in["tags"],
	}
}

func convertContainerAwsClusterControlPlaneList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlane(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneAwsServicesAuthentication(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"roleArn":         in["role_arn"],
		"roleSessionName": in["role_session_name"],
	}
}

func convertContainerAwsClusterControlPlaneAwsServicesAuthenticationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneAwsServicesAuthentication(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneConfigEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsClusterControlPlaneConfigEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneConfigEncryption(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneDatabaseEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsClusterControlPlaneDatabaseEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneDatabaseEncryption(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneMainVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsClusterControlPlaneMainVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneMainVolume(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneProxyConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"secretArn":     in["secret_arn"],
		"secretVersion": in["secret_version"],
	}
}

func convertContainerAwsClusterControlPlaneProxyConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneProxyConfig(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsClusterControlPlaneRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneRootVolume(v))
	}
	return out
}

func convertContainerAwsClusterControlPlaneSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"ec2KeyPair": in["ec2_key_pair"],
	}
}

func convertContainerAwsClusterControlPlaneSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterControlPlaneSshConfig(v))
	}
	return out
}

func convertContainerAwsClusterFleet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"project":    in["project"],
		"membership": in["membership"],
	}
}

func convertContainerAwsClusterFleetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterFleet(v))
	}
	return out
}

func convertContainerAwsClusterNetworking(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"podAddressCidrBlocks":     in["pod_address_cidr_blocks"],
		"serviceAddressCidrBlocks": in["service_address_cidr_blocks"],
		"vPCId":                    in["vpc_id"],
	}
}

func convertContainerAwsClusterNetworkingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterNetworking(v))
	}
	return out
}

func convertContainerAwsClusterWorkloadIdentityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"identityProvider": in["identity_provider"],
		"issuerUri":        in["issuer_uri"],
		"workloadPool":     in["workload_pool"],
	}
}

func convertContainerAwsClusterWorkloadIdentityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsClusterWorkloadIdentityConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolAutoscaling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxNodeCount": in["max_node_count"],
		"minNodeCount": in["min_node_count"],
	}
}

func convertContainerAwsNodePoolAutoscalingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolAutoscaling(v))
	}
	return out
}

func convertContainerAwsNodePoolConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configEncryption":   convertContainerAwsNodePoolConfigConfigEncryption(in["config_encryption"]),
		"iamInstanceProfile": in["iam_instance_profile"],
		"instanceType":       in["instance_type"],
		"labels":             in["labels"],
		"rootVolume":         convertContainerAwsNodePoolConfigRootVolume(in["root_volume"]),
		"securityGroupIds":   in["security_group_ids"],
		"sshConfig":          convertContainerAwsNodePoolConfigSshConfig(in["ssh_config"]),
		"tags":               in["tags"],
		"taints":             in["taints"],
	}
}

func convertContainerAwsNodePoolConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolConfigConfigEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kmsKeyArn": in["kms_key_arn"],
	}
}

func convertContainerAwsNodePoolConfigConfigEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolConfigConfigEncryption(v))
	}
	return out
}

func convertContainerAwsNodePoolConfigRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iops":       in["iops"],
		"kmsKeyArn":  in["kms_key_arn"],
		"sizeGib":    in["size_gib"],
		"volumeType": in["volume_type"],
	}
}

func convertContainerAwsNodePoolConfigRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolConfigRootVolume(v))
	}
	return out
}

func convertContainerAwsNodePoolConfigSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"ec2KeyPair": in["ec2_key_pair"],
	}
}

func convertContainerAwsNodePoolConfigSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolConfigSshConfig(v))
	}
	return out
}

func convertContainerAwsNodePoolConfigTaints(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"effect": in["effect"],
		"key":    in["key"],
		"value":  in["value"],
	}
}

func convertContainerAwsNodePoolConfigTaintsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolConfigTaints(v))
	}
	return out
}

func convertContainerAwsNodePoolMaxPodsConstraint(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxPodsPerNode": in["max_pods_per_node"],
	}
}

func convertContainerAwsNodePoolMaxPodsConstraintList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAwsNodePoolMaxPodsConstraint(v))
	}
	return out
}

func convertContainerAzureClusterAuthorization(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"adminUsers": in["admin_users"],
	}
}

func convertContainerAzureClusterAuthorizationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterAuthorization(v))
	}
	return out
}

func convertContainerAzureClusterAuthorizationAdminUsers(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"username": in["username"],
	}
}

func convertContainerAzureClusterAuthorizationAdminUsersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterAuthorizationAdminUsers(v))
	}
	return out
}

func convertContainerAzureClusterControlPlane(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sshConfig":          convertContainerAzureClusterControlPlaneSshConfig(in["ssh_config"]),
		"subnetId":           in["subnet_id"],
		"version":            in["version"],
		"databaseEncryption": convertContainerAzureClusterControlPlaneDatabaseEncryption(in["database_encryption"]),
		"mainVolume":         convertContainerAzureClusterControlPlaneMainVolume(in["main_volume"]),
		"proxyConfig":        convertContainerAzureClusterControlPlaneProxyConfig(in["proxy_config"]),
		"replicaPlacements":  in["replica_placements"],
		"rootVolume":         convertContainerAzureClusterControlPlaneRootVolume(in["root_volume"]),
		"tags":               in["tags"],
		"vmSize":             in["vm_size"],
	}
}

func convertContainerAzureClusterControlPlaneList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlane(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"authorizedKey": in["authorized_key"],
	}
}

func convertContainerAzureClusterControlPlaneSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneSshConfig(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneDatabaseEncryption(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"keyId": in["key_id"],
	}
}

func convertContainerAzureClusterControlPlaneDatabaseEncryptionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneDatabaseEncryption(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneMainVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureClusterControlPlaneMainVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneMainVolume(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneProxyConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resourceGroupId": in["resource_group_id"],
		"secretId":        in["secret_id"],
	}
}

func convertContainerAzureClusterControlPlaneProxyConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneProxyConfig(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneReplicaPlacements(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"azureAvailabilityZone": in["azure_availability_zone"],
		"subnetId":              in["subnet_id"],
	}
}

func convertContainerAzureClusterControlPlaneReplicaPlacementsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneReplicaPlacements(v))
	}
	return out
}

func convertContainerAzureClusterControlPlaneRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureClusterControlPlaneRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterControlPlaneRootVolume(v))
	}
	return out
}

func convertContainerAzureClusterFleet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"project":    in["project"],
		"membership": in["membership"],
	}
}

func convertContainerAzureClusterFleetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterFleet(v))
	}
	return out
}

func convertContainerAzureClusterNetworking(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"podAddressCidrBlocks":     in["pod_address_cidr_blocks"],
		"serviceAddressCidrBlocks": in["service_address_cidr_blocks"],
		"virtualNetworkId":         in["virtual_network_id"],
	}
}

func convertContainerAzureClusterNetworkingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterNetworking(v))
	}
	return out
}

func convertContainerAzureClusterWorkloadIdentityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"identityProvider": in["identity_provider"],
		"issuerUri":        in["issuer_uri"],
		"workloadPool":     in["workload_pool"],
	}
}

func convertContainerAzureClusterWorkloadIdentityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureClusterWorkloadIdentityConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolAutoscaling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxNodeCount": in["max_node_count"],
		"minNodeCount": in["min_node_count"],
	}
}

func convertContainerAzureNodePoolAutoscalingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolAutoscaling(v))
	}
	return out
}

func convertContainerAzureNodePoolConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sshConfig":  convertContainerAzureNodePoolConfigSshConfig(in["ssh_config"]),
		"rootVolume": convertContainerAzureNodePoolConfigRootVolume(in["root_volume"]),
		"tags":       in["tags"],
		"vmSize":     in["vm_size"],
	}
}

func convertContainerAzureNodePoolConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolConfigSshConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"authorizedKey": in["authorized_key"],
	}
}

func convertContainerAzureNodePoolConfigSshConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolConfigSshConfig(v))
	}
	return out
}

func convertContainerAzureNodePoolConfigRootVolume(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"sizeGib": in["size_gib"],
	}
}

func convertContainerAzureNodePoolConfigRootVolumeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolConfigRootVolume(v))
	}
	return out
}

func convertContainerAzureNodePoolMaxPodsConstraint(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxPodsPerNode": in["max_pods_per_node"],
	}
}

func convertContainerAzureNodePoolMaxPodsConstraintList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertContainerAzureNodePoolMaxPodsConstraint(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"stepId":              in["step_id"],
		"hadoopJob":           convertDataprocWorkflowTemplateJobsHadoopJob(in["hadoop_job"]),
		"hiveJob":             convertDataprocWorkflowTemplateJobsHiveJob(in["hive_job"]),
		"labels":              in["labels"],
		"pigJob":              convertDataprocWorkflowTemplateJobsPigJob(in["pig_job"]),
		"prerequisiteStepIds": in["prerequisite_step_ids"],
		"prestoJob":           convertDataprocWorkflowTemplateJobsPrestoJob(in["presto_job"]),
		"pysparkJob":          convertDataprocWorkflowTemplateJobsPysparkJob(in["pyspark_job"]),
		"scheduling":          convertDataprocWorkflowTemplateJobsScheduling(in["scheduling"]),
		"sparkJob":            convertDataprocWorkflowTemplateJobsSparkJob(in["spark_job"]),
		"sparkRJob":           convertDataprocWorkflowTemplateJobsSparkRJob(in["spark_r_job"]),
		"sparkSqlJob":         convertDataprocWorkflowTemplateJobsSparkSqlJob(in["spark_sql_job"]),
	}
}

func convertDataprocWorkflowTemplateJobsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobs(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsHadoopJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveUris":    in["archive_uris"],
		"args":           in["args"],
		"fileUris":       in["file_uris"],
		"jarFileUris":    in["jar_file_uris"],
		"loggingConfig":  convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(in["logging_config"]),
		"mainClass":      in["main_class"],
		"mainJarFileUri": in["main_jar_file_uri"],
		"properties":     in["properties"],
	}
}

func convertDataprocWorkflowTemplateJobsHadoopJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsHadoopJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsHiveJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"continueOnFailure": in["continue_on_failure"],
		"jarFileUris":       in["jar_file_uris"],
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateJobsHiveJobQueryList(in["query_list"]),
		"scriptVariables":   in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateJobsHiveJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsHiveJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsHiveJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateJobsHiveJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsHiveJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPigJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"continueOnFailure": in["continue_on_failure"],
		"jarFileUris":       in["jar_file_uris"],
		"loggingConfig":     convertDataprocWorkflowTemplateJobsPigJobLoggingConfig(in["logging_config"]),
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateJobsPigJobQueryList(in["query_list"]),
		"scriptVariables":   in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateJobsPigJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPigJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPigJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsPigJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPigJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPigJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateJobsPigJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPigJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPrestoJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clientTags":        in["client_tags"],
		"continueOnFailure": in["continue_on_failure"],
		"loggingConfig":     convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(in["logging_config"]),
		"outputFormat":      in["output_format"],
		"properties":        in["properties"],
		"queryFileUri":      in["query_file_uri"],
		"queryList":         convertDataprocWorkflowTemplateJobsPrestoJobQueryList(in["query_list"]),
	}
}

func convertDataprocWorkflowTemplateJobsPrestoJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPrestoJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPrestoJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPrestoJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateJobsPrestoJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPrestoJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPysparkJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"mainPythonFileUri": in["main_python_file_uri"],
		"archiveUris":       in["archive_uris"],
		"args":              in["args"],
		"fileUris":          in["file_uris"],
		"jarFileUris":       in["jar_file_uris"],
		"loggingConfig":     convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(in["logging_config"]),
		"properties":        in["properties"],
		"pythonFileUris":    in["python_file_uris"],
	}
}

func convertDataprocWorkflowTemplateJobsPysparkJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPysparkJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsPysparkJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsScheduling(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"maxFailuresPerHour": in["max_failures_per_hour"],
		"maxFailuresTotal":   in["max_failures_total"],
	}
}

func convertDataprocWorkflowTemplateJobsSchedulingList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsScheduling(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveUris":    in["archive_uris"],
		"args":           in["args"],
		"fileUris":       in["file_uris"],
		"jarFileUris":    in["jar_file_uris"],
		"loggingConfig":  convertDataprocWorkflowTemplateJobsSparkJobLoggingConfig(in["logging_config"]),
		"mainClass":      in["main_class"],
		"mainJarFileUri": in["main_jar_file_uri"],
		"properties":     in["properties"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkRJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"mainRFileUri":  in["main_r_file_uri"],
		"archiveUris":   in["archive_uris"],
		"args":          in["args"],
		"fileUris":      in["file_uris"],
		"loggingConfig": convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(in["logging_config"]),
		"properties":    in["properties"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkRJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkRJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkSqlJob(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"jarFileUris":     in["jar_file_uris"],
		"loggingConfig":   convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(in["logging_config"]),
		"properties":      in["properties"],
		"queryFileUri":    in["query_file_uri"],
		"queryList":       convertDataprocWorkflowTemplateJobsSparkSqlJobQueryList(in["query_list"]),
		"scriptVariables": in["script_variables"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkSqlJob(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"driverLogLevels": in["driver_log_levels"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobQueryList(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"queries": in["queries"],
	}
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobQueryListList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateJobsSparkSqlJobQueryList(v))
	}
	return out
}

func convertDataprocWorkflowTemplatePlacement(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterSelector": convertDataprocWorkflowTemplatePlacementClusterSelector(in["cluster_selector"]),
		"managedCluster":  convertDataprocWorkflowTemplatePlacementManagedCluster(in["managed_cluster"]),
	}
}

func convertDataprocWorkflowTemplatePlacementList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplatePlacement(v))
	}
	return out
}

func convertDataprocWorkflowTemplatePlacementClusterSelector(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterLabels": in["cluster_labels"],
		"zone":          in["zone"],
	}
}

func convertDataprocWorkflowTemplatePlacementClusterSelectorList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplatePlacementClusterSelector(v))
	}
	return out
}

func convertDataprocWorkflowTemplatePlacementManagedCluster(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clusterName": in["cluster_name"],
		"config":      convertDataprocWorkflowTemplateClusterClusterConfig(in["config"]),
		"labels":      in["labels"],
	}
}

func convertDataprocWorkflowTemplatePlacementManagedClusterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplatePlacementManagedCluster(v))
	}
	return out
}

func convertDataprocWorkflowTemplateParameters(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"fields":      in["fields"],
		"name":        in["name"],
		"description": in["description"],
		"validation":  convertDataprocWorkflowTemplateParametersValidation(in["validation"]),
	}
}

func convertDataprocWorkflowTemplateParametersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateParameters(v))
	}
	return out
}

func convertDataprocWorkflowTemplateParametersValidation(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"regex":  convertDataprocWorkflowTemplateParametersValidationRegex(in["regex"]),
		"values": convertDataprocWorkflowTemplateParametersValidationValues(in["values"]),
	}
}

func convertDataprocWorkflowTemplateParametersValidationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateParametersValidation(v))
	}
	return out
}

func convertDataprocWorkflowTemplateParametersValidationRegex(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"regexes": in["regexes"],
	}
}

func convertDataprocWorkflowTemplateParametersValidationRegexList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateParametersValidationRegex(v))
	}
	return out
}

func convertDataprocWorkflowTemplateParametersValidationValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"values": in["values"],
	}
}

func convertDataprocWorkflowTemplateParametersValidationValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateParametersValidationValues(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"accelerators":       in["accelerators"],
		"diskConfig":         convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfig(in["disk_config"]),
		"image":              in["image"],
		"machineType":        in["machine_type"],
		"minCpuPlatform":     in["min_cpu_platform"],
		"numInstances":       in["num_instances"],
		"preemptibility":     in["preemptibility"],
		"instanceNames":      in["instance_names"],
		"isPreemptible":      in["is_preemptible"],
		"managedGroupConfig": convertDataprocWorkflowTemplateClusterInstanceGroupConfigManagedGroupConfig(in["managed_group_config"]),
	}
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterInstanceGroupConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigAccelerators(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"acceleratorCount": in["accelerator_count"],
		"acceleratorType":  in["accelerator_type"],
	}
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigAcceleratorsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterInstanceGroupConfigAccelerators(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bootDiskSizeGb": in["boot_disk_size_gb"],
		"bootDiskType":   in["boot_disk_type"],
		"numLocalSsds":   in["num_local_ssds"],
	}
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterInstanceGroupConfigDiskConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigManagedGroupConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"instanceGroupManagerName": in["instance_group_manager_name"],
		"instanceTemplateName":     in["instance_template_name"],
	}
}

func convertDataprocWorkflowTemplateClusterInstanceGroupConfigManagedGroupConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterInstanceGroupConfigManagedGroupConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"autoscalingConfig":     convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfig(in["autoscaling_config"]),
		"encryptionConfig":      convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfig(in["encryption_config"]),
		"endpointConfig":        convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfig(in["endpoint_config"]),
		"gceClusterConfig":      convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfig(in["gce_cluster_config"]),
		"initializationActions": in["initialization_actions"],
		"lifecycleConfig":       convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfig(in["lifecycle_config"]),
		"masterConfig":          convertDataprocWorkflowTemplateClusterInstanceGroupConfig(in["master_config"]),
		"secondaryWorkerConfig": convertDataprocWorkflowTemplateClusterInstanceGroupConfig(in["secondary_worker_config"]),
		"securityConfig":        convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfig(in["security_config"]),
		"softwareConfig":        convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfig(in["software_config"]),
		"stagingBucket":         in["staging_bucket"],
		"tempBucket":            in["temp_bucket"],
		"workerConfig":          convertDataprocWorkflowTemplateClusterInstanceGroupConfig(in["worker_config"]),
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"policy": in["policy"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigAutoscalingConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"gcePdKmsKeyName": in["gce_pd_kms_key_name"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigEncryptionConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"enableHttpPortAccess": in["enable_http_port_access"],
		"httpPorts":            in["http_ports"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigEndpointConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"internalIPOnly":          in["internal_ip_only"],
		"metadata":                in["metadata"],
		"network":                 in["network"],
		"nodeGroupAffinity":       convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinity(in["node_group_affinity"]),
		"privateIPv6GoogleAccess": in["private_ipv6_google_access"],
		"reservationAffinity":     convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinity(in["reservation_affinity"]),
		"serviceAccount":          in["service_account"],
		"serviceAccountScopes":    in["service_account_scopes"],
		"subnetwork":              in["subnetwork"],
		"tags":                    in["tags"],
		"zone":                    in["zone"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinity(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nodeGroup": in["node_group"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinityList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigNodeGroupAffinity(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinity(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"consumeReservationType": in["consume_reservation_type"],
		"key":                    in["key"],
		"values":                 in["values"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinityList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigGceClusterConfigReservationAffinity(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigInitializationActions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"executableFile":   in["executable_file"],
		"executionTimeout": in["execution_timeout"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigInitializationActionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigInitializationActions(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"autoDeleteTime": in["auto_delete_time"],
		"autoDeleteTtl":  in["auto_delete_ttl"],
		"idleDeleteTtl":  in["idle_delete_ttl"],
		"idleStartTime":  in["idle_start_time"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigLifecycleConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"kerberosConfig": convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfig(in["kerberos_config"]),
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"crossRealmTrustAdminServer":    in["cross_realm_trust_admin_server"],
		"crossRealmTrustKdc":            in["cross_realm_trust_kdc"],
		"crossRealmTrustRealm":          in["cross_realm_trust_realm"],
		"crossRealmTrustSharedPassword": in["cross_realm_trust_shared_password"],
		"enableKerberos":                in["enable_kerberos"],
		"kdcDbKey":                      in["kdc_db_key"],
		"keyPassword":                   in["key_password"],
		"keystore":                      in["keystore"],
		"keystorePassword":              in["keystore_password"],
		"kmsKey":                        in["kms_key"],
		"realm":                         in["realm"],
		"rootPrincipalPassword":         in["root_principal_password"],
		"tgtLifetimeHours":              in["tgt_lifetime_hours"],
		"truststore":                    in["truststore"],
		"truststorePassword":            in["truststore_password"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigSecurityConfigKerberosConfig(v))
	}
	return out
}

func convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfig(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"imageVersion":       in["image_version"],
		"optionalComponents": in["optional_components"],
		"properties":         in["properties"],
	}
}

func convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfigList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertDataprocWorkflowTemplateClusterClusterConfigSoftwareConfig(v))
	}
	return out
}

func convertEventarcTriggerDestination(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"cloudFunction":   in["cloud_function"],
		"cloudRunService": convertEventarcTriggerDestinationCloudRunService(in["cloud_run_service"]),
	}
}

func convertEventarcTriggerDestinationList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerDestination(v))
	}
	return out
}

func convertEventarcTriggerDestinationCloudRunService(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"service": in["service"],
		"path":    in["path"],
		"region":  in["region"],
	}
}

func convertEventarcTriggerDestinationCloudRunServiceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerDestinationCloudRunService(v))
	}
	return out
}

func convertEventarcTriggerMatchingCriteria(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"attribute": in["attribute"],
		"value":     in["value"],
	}
}

func convertEventarcTriggerMatchingCriteriaList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerMatchingCriteria(v))
	}
	return out
}

func convertEventarcTriggerTransport(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"pubsub": convertEventarcTriggerTransportPubsub(in["pubsub"]),
	}
}

func convertEventarcTriggerTransportList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerTransport(v))
	}
	return out
}

func convertEventarcTriggerTransportPubsub(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"topic":        in["topic"],
		"subscription": in["subscription"],
	}
}

func convertEventarcTriggerTransportPubsubList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerTransportPubsub(v))
	}
	return out
}

func convertNetworkConnectivityHubRoutingVpcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri": in["uri"],
	}
}

func convertNetworkConnectivityHubRoutingVpcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivityHubRoutingVpcs(v))
	}
	return out
}

func convertNetworkConnectivitySpokeLinkedInterconnectAttachments(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
		"uris":                   in["uris"],
	}
}

func convertNetworkConnectivitySpokeLinkedInterconnectAttachmentsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeLinkedInterconnectAttachments(v))
	}
	return out
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstances(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"instances":              in["instances"],
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
	}
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeLinkedRouterApplianceInstances(v))
	}
	return out
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"iPAddress":      in["ip_address"],
		"virtualMachine": in["virtual_machine"],
	}
}

func convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstancesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeLinkedRouterApplianceInstancesInstances(v))
	}
	return out
}

func convertNetworkConnectivitySpokeLinkedVpnTunnels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"siteToSiteDataTransfer": in["site_to_site_data_transfer"],
		"uris":                   in["uris"],
	}
}

func convertNetworkConnectivitySpokeLinkedVpnTunnelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertNetworkConnectivitySpokeLinkedVpnTunnels(v))
	}
	return out
}

func convertOrgPolicyPolicySpec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"inheritFromParent": in["inherit_from_parent"],
		"reset":             in["reset"],
		"rules":             in["rules"],
		"etag":              in["etag"],
		"updateTime":        in["update_time"],
	}
}

func convertOrgPolicyPolicySpecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicySpec(v))
	}
	return out
}

func convertOrgPolicyPolicySpecRules(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAll":  in["allow_all"],
		"condition": convertOrgPolicyPolicySpecRulesCondition(in["condition"]),
		"denyAll":   in["deny_all"],
		"enforce":   in["enforce"],
		"values":    convertOrgPolicyPolicySpecRulesValues(in["values"]),
	}
}

func convertOrgPolicyPolicySpecRulesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicySpecRules(v))
	}
	return out
}

func convertOrgPolicyPolicySpecRulesCondition(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"description": in["description"],
		"expression":  in["expression"],
		"location":    in["location"],
		"title":       in["title"],
	}
}

func convertOrgPolicyPolicySpecRulesConditionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicySpecRulesCondition(v))
	}
	return out
}

func convertOrgPolicyPolicySpecRulesValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowedValues": in["allowed_values"],
		"deniedValues":  in["denied_values"],
	}
}

func convertOrgPolicyPolicySpecRulesValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOrgPolicyPolicySpecRulesValues(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentInstanceFilter(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"all":             in["all"],
		"exclusionLabels": in["exclusion_labels"],
		"inclusionLabels": in["inclusion_labels"],
		"inventories":     in["inventories"],
	}
}

func convertOsConfigOsPolicyAssignmentInstanceFilterList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentInstanceFilter(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"labels": in["labels"],
	}
}

func convertOsConfigOsPolicyAssignmentInstanceFilterExclusionLabelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentInstanceFilterExclusionLabels(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"labels": in["labels"],
	}
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInclusionLabelsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentInstanceFilterInclusionLabels(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInventories(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"oSShortName": in["os_short_name"],
		"oSVersion":   in["os_version"],
	}
}

func convertOsConfigOsPolicyAssignmentInstanceFilterInventoriesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentInstanceFilterInventories(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPolicies(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"id":                        in["id"],
		"mode":                      in["mode"],
		"resourceGroups":            in["resource_groups"],
		"allowNoResourceGroupMatch": in["allow_no_resource_group_match"],
		"description":               in["description"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPolicies(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"resources":        in["resources"],
		"inventoryFilters": in["inventory_filters"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroups(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"id":         in["id"],
		"exec":       convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(in["exec"]),
		"file":       convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(in["file"]),
		"pkg":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(in["pkg"]),
		"repository": convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(in["repository"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResources(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"validate": convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(in["validate"]),
		"enforce":  convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(in["enforce"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"interpreter":    in["interpreter"],
		"args":           in["args"],
		"file":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(in["file"]),
		"outputFilePath": in["output_file_path"],
		"script":         in["script"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"interpreter":    in["interpreter"],
		"args":           in["args"],
		"file":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(in["file"]),
		"outputFilePath": in["output_file_path"],
		"script":         in["script"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"path":        in["path"],
		"state":       in["state"],
		"content":     in["content"],
		"file":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(in["file"]),
		"permissions": in["permissions"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"desiredState": in["desired_state"],
		"apt":          convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(in["apt"]),
		"deb":          convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(in["deb"]),
		"googet":       convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(in["googet"]),
		"msi":          convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(in["msi"]),
		"rpm":          convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(in["rpm"]),
		"yum":          convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(in["yum"]),
		"zypper":       convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(in["zypper"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":   convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(in["source"]),
		"pullDeps": in["pull_deps"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":     convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(in["source"]),
		"properties": in["properties"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"source":   convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(in["source"]),
		"pullDeps": in["pull_deps"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowInsecure": in["allow_insecure"],
		"gcs":           convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(in["gcs"]),
		"localPath":     in["local_path"],
		"remote":        convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(in["remote"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"bucket":     in["bucket"],
		"object":     in["object"],
		"generation": in["generation"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"uri":            in["uri"],
		"sha256Checksum": in["sha256_checksum"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"apt":    convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(in["apt"]),
		"goo":    convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(in["goo"]),
		"yum":    convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(in["yum"]),
		"zypper": convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(in["zypper"]),
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"archiveType":  in["archive_type"],
		"components":   in["components"],
		"distribution": in["distribution"],
		"uri":          in["uri"],
		"gpgKey":       in["gpg_key"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
		"url":  in["url"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseUrl":     in["base_url"],
		"id":          in["id"],
		"displayName": in["display_name"],
		"gpgKeys":     in["gpg_keys"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseUrl":     in["base_url"],
		"id":          in["id"],
		"displayName": in["display_name"],
		"gpgKeys":     in["gpg_keys"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"oSShortName": in["os_short_name"],
		"oSVersion":   in["os_version"],
	}
}

func convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentRollout(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"disruptionBudget": convertOsConfigOsPolicyAssignmentRolloutDisruptionBudget(in["disruption_budget"]),
		"minWaitDuration":  in["min_wait_duration"],
	}
}

func convertOsConfigOsPolicyAssignmentRolloutList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentRollout(v))
	}
	return out
}

func convertOsConfigOsPolicyAssignmentRolloutDisruptionBudget(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"fixed":   in["fixed"],
		"percent": in["percent"],
	}
}

func convertOsConfigOsPolicyAssignmentRolloutDisruptionBudgetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertOsConfigOsPolicyAssignmentRolloutDisruptionBudget(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateIdentityConstraints(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowSubjectAltNamesPassthrough": in["allow_subject_alt_names_passthrough"],
		"allowSubjectPassthrough":         in["allow_subject_passthrough"],
		"celExpression":                   convertPrivatecaCertificateTemplateIdentityConstraintsCelExpression(in["cel_expression"]),
	}
}

func convertPrivatecaCertificateTemplateIdentityConstraintsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateIdentityConstraints(v))
	}
	return out
}

func convertPrivatecaCertificateTemplateIdentityConstraintsCelExpression(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"description": in["description"],
		"expression":  in["expression"],
		"location":    in["location"],
		"title":       in["title"],
	}
}

func convertPrivatecaCertificateTemplateIdentityConstraintsCelExpressionList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplateIdentityConstraintsCelExpression(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePassthroughExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"additionalExtensions": in["additional_extensions"],
		"knownExtensions":      in["known_extensions"],
	}
}

func convertPrivatecaCertificateTemplatePassthroughExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePassthroughExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePassthroughExtensionsAdditionalExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValues(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"additionalExtensions": in["additional_extensions"],
		"aiaOcspServers":       in["aia_ocsp_servers"],
		"caOptions":            convertPrivatecaCertificateTemplatePredefinedValuesCaOptions(in["ca_options"]),
		"keyUsage":             convertPrivatecaCertificateTemplatePredefinedValuesKeyUsage(in["key_usage"]),
		"policyIds":            in["policy_ids"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValues(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectId": convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(in["object_id"]),
		"value":    in["value"],
		"critical": in["critical"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectIdList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesCaOptions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"isCa":                in["is_ca"],
		"maxIssuerPathLength": in["max_issuer_path_length"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesCaOptionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesCaOptions(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"baseKeyUsage":             convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(in["base_key_usage"]),
		"extendedKeyUsage":         convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(in["extended_key_usage"]),
		"unknownExtendedKeyUsages": in["unknown_extended_key_usages"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"certSign":          in["cert_sign"],
		"contentCommitment": in["content_commitment"],
		"crlSign":           in["crl_sign"],
		"dataEncipherment":  in["data_encipherment"],
		"decipherOnly":      in["decipher_only"],
		"digitalSignature":  in["digital_signature"],
		"encipherOnly":      in["encipher_only"],
		"keyAgreement":      in["key_agreement"],
		"keyEncipherment":   in["key_encipherment"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"clientAuth":      in["client_auth"],
		"codeSigning":     in["code_signing"],
		"emailProtection": in["email_protection"],
		"ocspSigning":     in["ocsp_signing"],
		"serverAuth":      in["server_auth"],
		"timeStamping":    in["time_stamping"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsageList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsagesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages(v))
	}
	return out
}

func convertPrivatecaCertificateTemplatePredefinedValuesPolicyIds(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"objectIdPath": in["object_id_path"],
	}
}

func convertPrivatecaCertificateTemplatePredefinedValuesPolicyIdsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertPrivatecaCertificateTemplatePredefinedValuesPolicyIds(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyAndroidSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAllPackageNames": in["allow_all_package_names"],
		"allowedPackageNames":  in["allowed_package_names"],
	}
}

func convertRecaptchaEnterpriseKeyAndroidSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyAndroidSettings(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyIosSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"allowAllBundleIds": in["allow_all_bundle_ids"],
		"allowedBundleIds":  in["allowed_bundle_ids"],
	}
}

func convertRecaptchaEnterpriseKeyIosSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyIosSettings(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyTestingOptions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"testingChallenge": in["testing_challenge"],
		"testingScore":     in["testing_score"],
	}
}

func convertRecaptchaEnterpriseKeyTestingOptionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyTestingOptions(v))
	}
	return out
}

func convertRecaptchaEnterpriseKeyWebSettings(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"integrationType":             in["integration_type"],
		"allowAllDomains":             in["allow_all_domains"],
		"allowAmpTraffic":             in["allow_amp_traffic"],
		"allowedDomains":              in["allowed_domains"],
		"challengeSecurityPreference": in["challenge_security_preference"],
	}
}

func convertRecaptchaEnterpriseKeyWebSettingsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRecaptchaEnterpriseKeyWebSettings(v))
	}
	return out
}

func formatHCL(hcl string) (string, error) {
	b := bytes.Buffer{}
	r := strings.NewReader(hcl)
	if err := fmtcmd.Run(nil, nil, r, &b, fmtcmd.Options{}); err != nil {
		return "", err
	}
	return b.String(), nil
}
