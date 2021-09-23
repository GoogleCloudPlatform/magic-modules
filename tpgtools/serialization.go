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
	cloudbuildBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudbuild/beta"
	cloudresourcemanager "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager"
	cloudresourcemanagerBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager/beta"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	computeBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute/beta"
	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"
	dataprocBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc/beta"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	eventarcBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc/beta"
	gkehubBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/gkehub/beta"
	monitoring "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring"
	monitoringBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring/beta"
	orgpolicy "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy"
	orgpolicyBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/orgpolicy/beta"
	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"
	privatecaBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca/beta"
	fmtcmd "github.com/hashicorp/hcl/hcl/fmtcmd"
)

// DCLToTerraformReference converts a DCL resource name to the final tpgtools name
// after overrides are applied
func DCLToTerraformReference(resourceType, version string) (string, error) {
	if version == "beta" {
		switch resourceType {
		case "AssuredWorkloadsWorkload":
			return "google_assured_workloads_workload", nil
		case "CloudbuildWorkerPool":
			return "google_cloudbuild_worker_pool", nil
		case "CloudResourceManagerFolder":
			return "google_folder", nil
		case "CloudResourceManagerProject":
			return "google_project", nil
		case "ComputeFirewallPolicy":
			return "google_compute_firewall_policy", nil
		case "ComputeFirewallPolicyAssociation":
			return "google_compute_firewall_policy_association", nil
		case "ComputeFirewallPolicyRule":
			return "google_compute_firewall_policy_rule", nil
		case "ComputeForwardingRule":
			return "google_compute_forwarding_rule", nil
		case "ComputeGlobalForwardingRule":
			return "google_compute_global_forwarding_rule", nil
		case "DataprocWorkflowTemplate":
			return "google_dataproc_workflow_template", nil
		case "EventarcTrigger":
			return "google_eventarc_trigger", nil
		case "GkeHubFeature":
			return "google_gke_hub_feature", nil
		case "GkeHubFeatureMembership":
			return "google_gke_hub_feature_membership", nil
		case "MonitoringMetricsScope":
			return "google_monitoring_metrics_scope", nil
		case "MonitoringMonitoredProject":
			return "google_monitoring_monitored_project", nil
		case "OrgPolicyPolicy":
			return "google_org_policy_policy", nil
		case "PrivatecaCertificateTemplate":
			return "google_privateca_certificate_template", nil
		}
	}
	// If not found in sample version, fallthrough to GA
	switch resourceType {
	case "AssuredWorkloadsWorkload":
		return "google_assured_workloads_workload", nil
	case "CloudResourceManagerFolder":
		return "google_folder", nil
	case "CloudResourceManagerProject":
		return "google_project", nil
	case "ComputeFirewallPolicy":
		return "google_compute_firewall_policy", nil
	case "ComputeFirewallPolicyAssociation":
		return "google_compute_firewall_policy_association", nil
	case "ComputeFirewallPolicyRule":
		return "google_compute_firewall_policy_rule", nil
	case "ComputeForwardingRule":
		return "google_compute_forwarding_rule", nil
	case "ComputeGlobalForwardingRule":
		return "google_compute_global_forwarding_rule", nil
	case "DataprocWorkflowTemplate":
		return "google_dataproc_workflow_template", nil
	case "EventarcTrigger":
		return "google_eventarc_trigger", nil
	case "MonitoringMetricsScope":
		return "google_monitoring_metrics_scope", nil
	case "MonitoringMonitoredProject":
		return "google_monitoring_monitored_project", nil
	case "OrgPolicyPolicy":
		return "google_org_policy_policy", nil
	case "PrivatecaCertificateTemplate":
		return "google_privateca_certificate_template", nil
	default:
		return "", fmt.Errorf("Error retrieving Terraform name from DCL resource type: %s not found", resourceType)
	}

}

// DCLToTerraformSampleName converts a DCL resource name to the final tpgtools name
// after overrides are applied.
// e.g. cloudresourcemanager.project -> CloudResourceManagerProject
func DCLToTerraformSampleName(service, resource string) (string, string, error) {
	switch service + resource {
	case "assuredworkloadsworkload":
		return "AssuredWorkloads", "Workload", nil
	case "cloudresourcemanagerfolder":
		return "CloudResourceManager", "Folder", nil
	case "cloudresourcemanagerproject":
		return "CloudResourceManager", "Project", nil
	case "computefirewallpolicy":
		return "Compute", "FirewallPolicy", nil
	case "computefirewallpolicyassociation":
		return "Compute", "FirewallPolicyAssociation", nil
	case "computefirewallpolicyrule":
		return "Compute", "FirewallPolicyRule", nil
	case "computeforwardingrule":
		return "Compute", "ForwardingRule", nil
	case "dataprocworkflowtemplate":
		return "Dataproc", "WorkflowTemplate", nil
	case "eventarctrigger":
		return "Eventarc", "Trigger", nil
	case "monitoringmetricsscope":
		return "Monitoring", "MetricsScope", nil
	case "monitoringmonitoredproject":
		return "Monitoring", "MonitoredProject", nil
	case "orgpolicypolicy":
		return "OrgPolicy", "Policy", nil
	case "privatecacertificatetemplate":
		return "Privateca", "CertificateTemplate", nil
	default:
		return "", "", fmt.Errorf("Error retrieving Terraform sample name from DCL resource type: %s.%s not found", service, resource)
	}

}

// ConvertSampleJSONToHCL unmarshals json to an HCL string.
func ConvertSampleJSONToHCL(resourceType string, version string, b []byte) (string, error) {
	if version == "beta" {
		switch resourceType {
		case "AssuredWorkloadsWorkload":
			r := &assuredworkloadsBeta.Workload{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return AssuredWorkloadsWorkloadBetaAsHCL(*r)
		case "CloudbuildWorkerPool":
			r := &cloudbuildBeta.WorkerPool{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return CloudbuildWorkerPoolBetaAsHCL(*r)
		case "CloudResourceManagerFolder":
			r := &cloudresourcemanagerBeta.Folder{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return CloudResourceManagerFolderBetaAsHCL(*r)
		case "CloudResourceManagerProject":
			r := &cloudresourcemanagerBeta.Project{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return serializeBetaProjectToHCL(*r)
		case "ComputeFirewallPolicy":
			r := &computeBeta.FirewallPolicy{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyBetaAsHCL(*r)
		case "ComputeFirewallPolicyAssociation":
			r := &computeBeta.FirewallPolicyAssociation{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyAssociationBetaAsHCL(*r)
		case "ComputeFirewallPolicyRule":
			r := &computeBeta.FirewallPolicyRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeFirewallPolicyRuleBetaAsHCL(*r)
		case "ComputeForwardingRule":
			r := &computeBeta.ForwardingRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeForwardingRuleBetaAsHCL(*r)
		case "ComputeGlobalForwardingRule":
			r := &computeBeta.ForwardingRule{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return ComputeGlobalForwardingRuleBetaAsHCL(*r)
		case "DataprocWorkflowTemplate":
			r := &dataprocBeta.WorkflowTemplate{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return DataprocWorkflowTemplateBetaAsHCL(*r)
		case "EventarcTrigger":
			r := &eventarcBeta.Trigger{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return EventarcTriggerBetaAsHCL(*r)
		case "GkeHubFeature":
			r := &gkehubBeta.Feature{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return GkeHubFeatureBetaAsHCL(*r)
		case "GkeHubFeatureMembership":
			r := &gkehubBeta.FeatureMembership{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return GkeHubFeatureMembershipBetaAsHCL(*r)
		case "MonitoringMetricsScope":
			r := &monitoringBeta.MetricsScope{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return MonitoringMetricsScopeBetaAsHCL(*r)
		case "MonitoringMonitoredProject":
			r := &monitoringBeta.MonitoredProject{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return MonitoringMonitoredProjectBetaAsHCL(*r)
		case "OrgPolicyPolicy":
			r := &orgpolicyBeta.Policy{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return OrgPolicyPolicyBetaAsHCL(*r)
		case "PrivatecaCertificateTemplate":
			r := &privatecaBeta.CertificateTemplate{}
			if err := json.Unmarshal(b, r); err != nil {
				return "", err
			}
			return PrivatecaCertificateTemplateBetaAsHCL(*r)
		}
	}
	// If not found in sample version, fallthrough to GA
	switch resourceType {
	case "AssuredWorkloadsWorkload":
		r := &assuredworkloads.Workload{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return AssuredWorkloadsWorkloadAsHCL(*r)
	case "CloudResourceManagerFolder":
		r := &cloudresourcemanager.Folder{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return CloudResourceManagerFolderAsHCL(*r)
	case "CloudResourceManagerProject":
		r := &cloudresourcemanager.Project{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return serializeGAProjectToHCL(*r)
	case "ComputeFirewallPolicy":
		r := &compute.FirewallPolicy{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyAsHCL(*r)
	case "ComputeFirewallPolicyAssociation":
		r := &compute.FirewallPolicyAssociation{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyAssociationAsHCL(*r)
	case "ComputeFirewallPolicyRule":
		r := &compute.FirewallPolicyRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeFirewallPolicyRuleAsHCL(*r)
	case "ComputeForwardingRule":
		r := &compute.ForwardingRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeForwardingRuleAsHCL(*r)
	case "ComputeGlobalForwardingRule":
		r := &compute.ForwardingRule{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return ComputeGlobalForwardingRuleAsHCL(*r)
	case "DataprocWorkflowTemplate":
		r := &dataproc.WorkflowTemplate{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return DataprocWorkflowTemplateAsHCL(*r)
	case "EventarcTrigger":
		r := &eventarc.Trigger{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return EventarcTriggerAsHCL(*r)
	case "MonitoringMetricsScope":
		r := &monitoring.MetricsScope{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return MonitoringMetricsScopeAsHCL(*r)
	case "MonitoringMonitoredProject":
		r := &monitoring.MonitoredProject{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return MonitoringMonitoredProjectAsHCL(*r)
	case "OrgPolicyPolicy":
		r := &orgpolicy.Policy{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return OrgPolicyPolicyAsHCL(*r)
	case "PrivatecaCertificateTemplate":
		r := &privateca.CertificateTemplate{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return PrivatecaCertificateTemplateAsHCL(*r)
	default:
		//return fmt.Sprintf("%s resource not supported in tpgtools", resourceType), nil
		return "", fmt.Errorf("Error converting sample JSON to HCL: %s not found", resourceType)
	}

}

// AssuredWorkloadsWorkloadBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func AssuredWorkloadsWorkloadBetaAsHCL(r assuredworkloadsBeta.Workload) (string, error) {
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
	if r.ProvisionedResourcesParent != nil {
		outputConfig += fmt.Sprintf("\tprovisioned_resources_parent = %#v\n", *r.ProvisionedResourcesParent)
	}
	if r.ResourceSettings != nil {
		for _, v := range r.ResourceSettings {
			outputConfig += fmt.Sprintf("\tresource_settings %s\n", convertAssuredWorkloadsWorkloadBetaResourceSettingsToHCL(&v))
		}
	}
	return formatHCL(outputConfig + "}")
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
func CloudbuildWorkerPoolBetaAsHCL(r cloudbuildBeta.WorkerPool) (string, error) {
	outputConfig := "resource \"google_cloudbuild_worker_pool\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
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
	return formatHCL(outputConfig + "}")
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
func CloudResourceManagerFolderBetaAsHCL(r cloudresourcemanagerBeta.Folder) (string, error) {
	outputConfig := "resource \"google_folder\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	return formatHCL(outputConfig + "}")
}

// CloudResourceManagerProjectBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerProjectBetaAsHCL(r cloudresourcemanagerBeta.Project) (string, error) {
	outputConfig := "resource \"google_project\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplayname = %#v\n", *r.DisplayName)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyBetaAsHCL(r computeBeta.FirewallPolicy) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyAssociationBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAssociationBetaAsHCL(r computeBeta.FirewallPolicyAssociation) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyRuleBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyRuleBetaAsHCL(r computeBeta.FirewallPolicyRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
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
func ComputeForwardingRuleBetaAsHCL(r computeBeta.ForwardingRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeGlobalForwardingRuleBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeGlobalForwardingRuleBetaAsHCL(r computeBeta.ForwardingRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
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

// DataprocWorkflowTemplateBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func DataprocWorkflowTemplateBetaAsHCL(r dataprocBeta.WorkflowTemplate) (string, error) {
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
	return formatHCL(outputConfig + "}")
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsHadoopJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsHadoopJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsHiveJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsPigJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsPigJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsPigJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkRJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkRJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateBetaJobsSparkSqlJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateBetaJobsSparkSqlJobLoggingConfigToHCL(r *dataprocBeta.WorkflowTemplateJobsSparkSqlJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

// EventarcTriggerBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func EventarcTriggerBetaAsHCL(r eventarcBeta.Trigger) (string, error) {
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
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if v := convertEventarcTriggerBetaTransportToHCL(r.Transport); v != "" {
		outputConfig += fmt.Sprintf("\ttransport %s\n", v)
	}
	return formatHCL(outputConfig + "}")
}

func convertEventarcTriggerBetaDestinationToHCL(r *eventarcBeta.TriggerDestination) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
func GkeHubFeatureBetaAsHCL(r gkehubBeta.Feature) (string, error) {
	outputConfig := "resource \"google_gke_hub_feature\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertGkeHubFeatureBetaSpecToHCL(r.Spec); v != "" {
		outputConfig += fmt.Sprintf("\tspec %s\n", v)
	}
	return formatHCL(outputConfig + "}")
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

// GkeHubFeatureMembershipBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func GkeHubFeatureMembershipBetaAsHCL(r gkehubBeta.FeatureMembership) (string, error) {
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
	return formatHCL(outputConfig + "}")
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

// MonitoringMetricsScopeBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func MonitoringMetricsScopeBetaAsHCL(r monitoringBeta.MetricsScope) (string, error) {
	outputConfig := "resource \"google_monitoring_metrics_scope\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return formatHCL(outputConfig + "}")
}

func convertMonitoringMetricsScopeBetaMonitoredProjectsToHCL(r *monitoringBeta.MetricsScopeMonitoredProjects) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// MonitoringMonitoredProjectBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func MonitoringMonitoredProjectBetaAsHCL(r monitoringBeta.MonitoredProject) (string, error) {
	outputConfig := "resource \"google_monitoring_monitored_project\" \"output\" {\n"
	if r.MetricsScope != nil {
		outputConfig += fmt.Sprintf("\tmetrics_scope = %#v\n", *r.MetricsScope)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return formatHCL(outputConfig + "}")
}

// OrgPolicyPolicyBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OrgPolicyPolicyBetaAsHCL(r orgpolicyBeta.Policy) (string, error) {
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
	return formatHCL(outputConfig + "}")
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

// PrivatecaCertificateTemplateBetaAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func PrivatecaCertificateTemplateBetaAsHCL(r privatecaBeta.CertificateTemplate) (string, error) {
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
	if v := convertPrivatecaCertificateTemplateBetaPassthroughExtensionsToHCL(r.PassthroughExtensions); v != "" {
		outputConfig += fmt.Sprintf("\tpassthrough_extensions %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplateBetaPredefinedValuesToHCL(r.PredefinedValues); v != "" {
		outputConfig += fmt.Sprintf("\tpredefined_values %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return formatHCL(outputConfig + "}")
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

// AssuredWorkloadsWorkloadAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func AssuredWorkloadsWorkloadAsHCL(r assuredworkloads.Workload) (string, error) {
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
	if r.ProvisionedResourcesParent != nil {
		outputConfig += fmt.Sprintf("\tprovisioned_resources_parent = %#v\n", *r.ProvisionedResourcesParent)
	}
	if r.ResourceSettings != nil {
		for _, v := range r.ResourceSettings {
			outputConfig += fmt.Sprintf("\tresource_settings %s\n", convertAssuredWorkloadsWorkloadResourceSettingsToHCL(&v))
		}
	}
	return formatHCL(outputConfig + "}")
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

// CloudResourceManagerFolderAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerFolderAsHCL(r cloudresourcemanager.Folder) (string, error) {
	outputConfig := "resource \"google_folder\" \"output\" {\n"
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplay_name = %#v\n", *r.DisplayName)
	}
	return formatHCL(outputConfig + "}")
}

// CloudResourceManagerProjectAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func CloudResourceManagerProjectAsHCL(r cloudresourcemanager.Project) (string, error) {
	outputConfig := "resource \"google_project\" \"output\" {\n"
	if r.DisplayName != nil {
		outputConfig += fmt.Sprintf("\tdisplayname = %#v\n", *r.DisplayName)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Parent != nil {
		outputConfig += fmt.Sprintf("\tparent = %#v\n", *r.Parent)
	}
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAsHCL(r compute.FirewallPolicy) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyAssociationAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyAssociationAsHCL(r compute.FirewallPolicyAssociation) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeFirewallPolicyRuleAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeFirewallPolicyRuleAsHCL(r compute.FirewallPolicyRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
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
func ComputeForwardingRuleAsHCL(r compute.ForwardingRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
}

// ComputeGlobalForwardingRuleAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func ComputeGlobalForwardingRuleAsHCL(r compute.ForwardingRule) (string, error) {
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
	return formatHCL(outputConfig + "}")
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

// DataprocWorkflowTemplateAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func DataprocWorkflowTemplateAsHCL(r dataproc.WorkflowTemplate) (string, error) {
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
	return formatHCL(outputConfig + "}")
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsHadoopJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsHadoopJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsHiveJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsPigJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsPigJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsPigJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkRJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkRJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	if r.QueryFileUri != nil {
		outputConfig += fmt.Sprintf("\tquery_file_uri = %#v\n", *r.QueryFileUri)
	}
	if v := convertDataprocWorkflowTemplateJobsSparkSqlJobQueryListToHCL(r.QueryList); v != "" {
		outputConfig += fmt.Sprintf("\tquery_list %s\n", v)
	}
	return outputConfig + "}"
}

func convertDataprocWorkflowTemplateJobsSparkSqlJobLoggingConfigToHCL(r *dataproc.WorkflowTemplateJobsSparkSqlJobLoggingConfig) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
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
	return outputConfig + "}"
}

// EventarcTriggerAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func EventarcTriggerAsHCL(r eventarc.Trigger) (string, error) {
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
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if r.ServiceAccount != nil {
		outputConfig += fmt.Sprintf("\tservice_account = %#v\n", *r.ServiceAccount)
	}
	if v := convertEventarcTriggerTransportToHCL(r.Transport); v != "" {
		outputConfig += fmt.Sprintf("\ttransport %s\n", v)
	}
	return formatHCL(outputConfig + "}")
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

// MonitoringMetricsScopeAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func MonitoringMetricsScopeAsHCL(r monitoring.MetricsScope) (string, error) {
	outputConfig := "resource \"google_monitoring_metrics_scope\" \"output\" {\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return formatHCL(outputConfig + "}")
}

func convertMonitoringMetricsScopeMonitoredProjectsToHCL(r *monitoring.MetricsScopeMonitoredProjects) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

// MonitoringMonitoredProjectAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func MonitoringMonitoredProjectAsHCL(r monitoring.MonitoredProject) (string, error) {
	outputConfig := "resource \"google_monitoring_monitored_project\" \"output\" {\n"
	if r.MetricsScope != nil {
		outputConfig += fmt.Sprintf("\tmetrics_scope = %#v\n", *r.MetricsScope)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return formatHCL(outputConfig + "}")
}

// OrgPolicyPolicyAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func OrgPolicyPolicyAsHCL(r orgpolicy.Policy) (string, error) {
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
	return formatHCL(outputConfig + "}")
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

// PrivatecaCertificateTemplateAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func PrivatecaCertificateTemplateAsHCL(r privateca.CertificateTemplate) (string, error) {
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
	if v := convertPrivatecaCertificateTemplatePassthroughExtensionsToHCL(r.PassthroughExtensions); v != "" {
		outputConfig += fmt.Sprintf("\tpassthrough_extensions %s\n", v)
	}
	if v := convertPrivatecaCertificateTemplatePredefinedValuesToHCL(r.PredefinedValues); v != "" {
		outputConfig += fmt.Sprintf("\tpredefined_values %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	return formatHCL(outputConfig + "}")
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
		"imageVersion": in["image_version"],
		"properties":   in["properties"],
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
		"httpsProxy":   in["https_proxy"],
		"policyDir":    in["policy_dir"],
		"secretType":   in["secret_type"],
		"syncBranch":   in["sync_branch"],
		"syncRepo":     in["sync_repo"],
		"syncRev":      in["sync_rev"],
		"syncWaitSecs": in["sync_wait_secs"],
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

func convertMonitoringMetricsScopeBetaMonitoredProjects(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"createTime": in["create_time"],
		"name":       in["name"],
	}
}

func convertMonitoringMetricsScopeBetaMonitoredProjectsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertMonitoringMetricsScopeBetaMonitoredProjects(v))
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
		"imageVersion": in["image_version"],
		"properties":   in["properties"],
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

func convertMonitoringMetricsScopeMonitoredProjects(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"createTime": in["create_time"],
		"name":       in["name"],
	}
}

func convertMonitoringMetricsScopeMonitoredProjectsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertMonitoringMetricsScopeMonitoredProjects(v))
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

func formatHCL(hcl string) (string, error) {
	b := bytes.Buffer{}
	r := strings.NewReader(hcl)
	if err := fmtcmd.Run(nil, nil, r, &b, fmtcmd.Options{}); err != nil {
		return "", err
	}
	return b.String(), nil
}
