package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	dataproc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc"
	dataprocBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataproc/beta"
	eventarc "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc"
	eventarcBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/eventarc/beta"
	run "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/run"
	fmtcmd "github.com/hashicorp/hcl/hcl/fmtcmd"
)

// DCLToTerraformReferencce converts a DCL resource name to the final tpgtools name
// after overrides are applied
func DCLToTerraformReference(resourceType, version string) (string, error) {
	if version == "beta" {
		switch resourceType {
		case "DataprocWorkflowTemplate":
			return "google_dataproc_workflow_template", nil
		case "EventarcTrigger":
			return "google_eventarc_trigger", nil
		}
	}
	// If not found in sample version, fallthrough to GA
	switch resourceType {
	case "DataprocWorkflowTemplate":
		return "google_dataproc_workflow_template", nil
	case "EventarcTrigger":
		return "google_eventarc_trigger", nil
	case "RunService":
		return "google_cloud_run_service", nil
	default:
		return "", fmt.Errorf("Error retrieving Terraform name from DCL resource type: %s not found", resourceType)
	}

}

// ConvertSampleJSONToDCLResource unmarshals json to a DCL resource specified by the resource type
func ConvertSampleJSONToHCL(resourceType string, version string, b []byte) (string, error) {
	if version == "beta" {
		switch resourceType {
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
		}
	}
	// If not found in sample version, fallthrough to GA
	switch resourceType {
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
	case "RunService":
		r := &run.Service{}
		if err := json.Unmarshal(b, r); err != nil {
			return "", err
		}
		return RunServiceAsHCL(*r)
	default:
		//return fmt.Sprintf("%s resource not supported in tpgtools", resourceType), nil
		return "", fmt.Errorf("Error converting sample JSON to HCL: %s not found", resourceType)
	}

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
	if r.EventFilters != nil {
		for _, v := range r.EventFilters {
			outputConfig += fmt.Sprintf("\tmatching_criteria %s\n", convertEventarcTriggerEventFiltersToHCL(&v))
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
	if v := convertEventarcTriggerDestinationCloudRunToHCL(r.CloudRun); v != "" {
		outputConfig += fmt.Sprintf("\tcloud_run_service %s\n", v)
	}
	return outputConfig + "}"
}

func convertEventarcTriggerDestinationCloudRunToHCL(r *eventarc.TriggerDestinationCloudRun) string {
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

func convertEventarcTriggerEventFiltersToHCL(r *eventarc.TriggerEventFilters) string {
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

// RunServiceAsHCL returns a string representation of the specified resource in HCL.
// The generated HCL will include every settable field as a literal - that is, no
// variables, no references.  This may not be the best possible representation, but
// the crucial point is that `terraform import; terraform apply` will not produce
// any changes.  We do not validate that the resource specified will pass terraform
// validation unless is an object returned from the API after an Apply.
func RunServiceAsHCL(r run.Service) (string, error) {
	outputConfig := "resource \"google_cloud_run_service\" \"output\" {\n"
	if r.Location != nil {
		outputConfig += fmt.Sprintf("\tlocation = %#v\n", *r.Location)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.ApiVersion != nil {
		outputConfig += fmt.Sprintf("\tapi_version = %#v\n", *r.ApiVersion)
	}
	if r.Kind != nil {
		outputConfig += fmt.Sprintf("\tkind = %#v\n", *r.Kind)
	}
	if v := convertRunServiceMetadataToHCL(r.Metadata); v != "" {
		outputConfig += fmt.Sprintf("\tmetadata %s\n", v)
	}
	if r.Project != nil {
		outputConfig += fmt.Sprintf("\tproject = %#v\n", *r.Project)
	}
	if v := convertRunServiceSpecToHCL(r.Spec); v != "" {
		outputConfig += v
	}
	return formatHCL(outputConfig + "}")
}

func convertRunServiceMetadataToHCL(r *run.ServiceMetadata) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClusterName != nil {
		outputConfig += fmt.Sprintf("\tcluster_name = %#v\n", *r.ClusterName)
	}
	if v := convertRunServiceMetadataCreateTimeToHCL(r.CreateTime); v != "" {
		outputConfig += fmt.Sprintf("\tcreate_time %s\n", v)
	}
	if v := convertRunServiceMetadataDeleteTimeToHCL(r.DeleteTime); v != "" {
		outputConfig += fmt.Sprintf("\tdelete_time %s\n", v)
	}
	if r.DeletionGracePeriodSeconds != nil {
		outputConfig += fmt.Sprintf("\tdeletion_grace_period_seconds = %#v\n", *r.DeletionGracePeriodSeconds)
	}
	if r.Finalizers != nil {
		outputConfig += "\tfinalizers = ["
		for _, v := range r.Finalizers {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.GenerateName != nil {
		outputConfig += fmt.Sprintf("\tgenerate_name = %#v\n", *r.GenerateName)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Namespace != nil {
		outputConfig += fmt.Sprintf("\tnamespace = %#v\n", *r.Namespace)
	}
	if r.OwnerReferences != nil {
		for _, v := range r.OwnerReferences {
			outputConfig += fmt.Sprintf("\towner_references %s\n", convertRunServiceMetadataOwnerReferencesToHCL(&v))
		}
	}
	if r.ResourceVersion != nil {
		outputConfig += fmt.Sprintf("\tresource_version = %#v\n", *r.ResourceVersion)
	}
	if r.SelfLink != nil {
		outputConfig += fmt.Sprintf("\tself_link = %#v\n", *r.SelfLink)
	}
	if r.Uid != nil {
		outputConfig += fmt.Sprintf("\tuid = %#v\n", *r.Uid)
	}
	return outputConfig + "}"
}

func convertRunServiceMetadataCreateTimeToHCL(r *run.ServiceMetadataCreateTime) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Nanos != nil {
		outputConfig += fmt.Sprintf("\tnanos = %#v\n", *r.Nanos)
	}
	if r.Seconds != nil {
		outputConfig += fmt.Sprintf("\tseconds = %#v\n", *r.Seconds)
	}
	return outputConfig + "}"
}

func convertRunServiceMetadataDeleteTimeToHCL(r *run.ServiceMetadataDeleteTime) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Nanos != nil {
		outputConfig += fmt.Sprintf("\tnanos = %#v\n", *r.Nanos)
	}
	if r.Seconds != nil {
		outputConfig += fmt.Sprintf("\tseconds = %#v\n", *r.Seconds)
	}
	return outputConfig + "}"
}

func convertRunServiceMetadataOwnerReferencesToHCL(r *run.ServiceMetadataOwnerReferences) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ApiVersion != nil {
		outputConfig += fmt.Sprintf("\tapi_version = %#v\n", *r.ApiVersion)
	}
	if r.BlockOwnerDeletion != nil {
		outputConfig += fmt.Sprintf("\tblock_owner_deletion = %#v\n", *r.BlockOwnerDeletion)
	}
	if r.Controller != nil {
		outputConfig += fmt.Sprintf("\tcontroller = %#v\n", *r.Controller)
	}
	if r.Kind != nil {
		outputConfig += fmt.Sprintf("\tkind = %#v\n", *r.Kind)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Uid != nil {
		outputConfig += fmt.Sprintf("\tuid = %#v\n", *r.Uid)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecToHCL(r *run.ServiceSpec) string {
	if r == nil {
		return ""
	}
	outputConfig := ""
	if v := convertRunServiceSpecTemplateToHCL(r.Template); v != "" {
		outputConfig += fmt.Sprintf("template %s\n", v)
	}
	if r.Traffic != nil {
		for _, v := range r.Traffic {
			outputConfig += fmt.Sprintf("traffic %s\n", convertRunServiceSpecTrafficToHCL(&v))
		}
	}
	return outputConfig
}

func convertRunServiceSpecTemplateToHCL(r *run.ServiceSpecTemplate) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateMetadataToHCL(r.Metadata); v != "" {
		outputConfig += fmt.Sprintf("\tmetadata %s\n", v)
	}
	if v := convertRunServiceSpecTemplateSpecToHCL(r.Spec); v != "" {
		outputConfig += fmt.Sprintf("\tspec %s\n", v)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateMetadataToHCL(r *run.ServiceSpecTemplateMetadata) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ClusterName != nil {
		outputConfig += fmt.Sprintf("\tcluster_name = %#v\n", *r.ClusterName)
	}
	if v := convertRunServiceSpecTemplateMetadataCreateTimeToHCL(r.CreateTime); v != "" {
		outputConfig += fmt.Sprintf("\tcreate_time %s\n", v)
	}
	if v := convertRunServiceSpecTemplateMetadataDeleteTimeToHCL(r.DeleteTime); v != "" {
		outputConfig += fmt.Sprintf("\tdelete_time %s\n", v)
	}
	if r.DeletionGracePeriodSeconds != nil {
		outputConfig += fmt.Sprintf("\tdeletion_grace_period_seconds = %#v\n", *r.DeletionGracePeriodSeconds)
	}
	if r.Finalizers != nil {
		outputConfig += "\tfinalizers = ["
		for _, v := range r.Finalizers {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.GenerateName != nil {
		outputConfig += fmt.Sprintf("\tgenerate_name = %#v\n", *r.GenerateName)
	}
	if r.Generation != nil {
		outputConfig += fmt.Sprintf("\tgeneration = %#v\n", *r.Generation)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Namespace != nil {
		outputConfig += fmt.Sprintf("\tnamespace = %#v\n", *r.Namespace)
	}
	if r.OwnerReferences != nil {
		for _, v := range r.OwnerReferences {
			outputConfig += fmt.Sprintf("\towner_references %s\n", convertRunServiceSpecTemplateMetadataOwnerReferencesToHCL(&v))
		}
	}
	if r.ResourceVersion != nil {
		outputConfig += fmt.Sprintf("\tresource_version = %#v\n", *r.ResourceVersion)
	}
	if r.SelfLink != nil {
		outputConfig += fmt.Sprintf("\tself_link = %#v\n", *r.SelfLink)
	}
	if r.Uid != nil {
		outputConfig += fmt.Sprintf("\tuid = %#v\n", *r.Uid)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateMetadataCreateTimeToHCL(r *run.ServiceSpecTemplateMetadataCreateTime) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Nanos != nil {
		outputConfig += fmt.Sprintf("\tnanos = %#v\n", *r.Nanos)
	}
	if r.Seconds != nil {
		outputConfig += fmt.Sprintf("\tseconds = %#v\n", *r.Seconds)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateMetadataDeleteTimeToHCL(r *run.ServiceSpecTemplateMetadataDeleteTime) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Nanos != nil {
		outputConfig += fmt.Sprintf("\tnanos = %#v\n", *r.Nanos)
	}
	if r.Seconds != nil {
		outputConfig += fmt.Sprintf("\tseconds = %#v\n", *r.Seconds)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateMetadataOwnerReferencesToHCL(r *run.ServiceSpecTemplateMetadataOwnerReferences) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ApiVersion != nil {
		outputConfig += fmt.Sprintf("\tapi_version = %#v\n", *r.ApiVersion)
	}
	if r.BlockOwnerDeletion != nil {
		outputConfig += fmt.Sprintf("\tblock_owner_deletion = %#v\n", *r.BlockOwnerDeletion)
	}
	if r.Controller != nil {
		outputConfig += fmt.Sprintf("\tcontroller = %#v\n", *r.Controller)
	}
	if r.Kind != nil {
		outputConfig += fmt.Sprintf("\tkind = %#v\n", *r.Kind)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Uid != nil {
		outputConfig += fmt.Sprintf("\tuid = %#v\n", *r.Uid)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecToHCL(r *run.ServiceSpecTemplateSpec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContainerConcurrency != nil {
		outputConfig += fmt.Sprintf("\tcontainer_concurrency = %#v\n", *r.ContainerConcurrency)
	}
	if r.Containers != nil {
		for _, v := range r.Containers {
			outputConfig += fmt.Sprintf("\tcontainers %s\n", convertRunServiceSpecTemplateSpecContainersToHCL(&v))
		}
	}
	if r.ServiceAccountName != nil {
		outputConfig += fmt.Sprintf("\tservice_account_name = %#v\n", *r.ServiceAccountName)
	}
	if r.TimeoutSeconds != nil {
		outputConfig += fmt.Sprintf("\ttimeout_seconds = %#v\n", *r.TimeoutSeconds)
	}
	if r.Volumes != nil {
		for _, v := range r.Volumes {
			outputConfig += fmt.Sprintf("\tvolumes %s\n", convertRunServiceSpecTemplateSpecVolumesToHCL(&v))
		}
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersToHCL(r *run.ServiceSpecTemplateSpecContainers) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Args != nil {
		outputConfig += "\targs = ["
		for _, v := range r.Args {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Command != nil {
		outputConfig += "\tcommand = ["
		for _, v := range r.Command {
			outputConfig += fmt.Sprintf("%#v, ", v)
		}
		outputConfig += "]\n"
	}
	if r.Env != nil {
		for _, v := range r.Env {
			outputConfig += fmt.Sprintf("\tenv %s\n", convertRunServiceSpecTemplateSpecContainersEnvToHCL(&v))
		}
	}
	if r.EnvFrom != nil {
		for _, v := range r.EnvFrom {
			outputConfig += fmt.Sprintf("\tenv_from %s\n", convertRunServiceSpecTemplateSpecContainersEnvFromToHCL(&v))
		}
	}
	if r.Image != nil {
		outputConfig += fmt.Sprintf("\timage = %#v\n", *r.Image)
	}
	if r.ImagePullPolicy != nil {
		outputConfig += fmt.Sprintf("\timage_pull_policy = %#v\n", *r.ImagePullPolicy)
	}
	if v := convertRunServiceSpecTemplateSpecContainersLivenessProbeToHCL(r.LivenessProbe); v != "" {
		outputConfig += fmt.Sprintf("\tliveness_probe %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Ports != nil {
		for _, v := range r.Ports {
			outputConfig += fmt.Sprintf("\tports %s\n", convertRunServiceSpecTemplateSpecContainersPortsToHCL(&v))
		}
	}
	if v := convertRunServiceSpecTemplateSpecContainersReadinessProbeToHCL(r.ReadinessProbe); v != "" {
		outputConfig += fmt.Sprintf("\treadiness_probe %s\n", v)
	}
	if v := convertRunServiceSpecTemplateSpecContainersResourcesToHCL(r.Resources); v != "" {
		outputConfig += fmt.Sprintf("\tresources %s\n", v)
	}
	if v := convertRunServiceSpecTemplateSpecContainersSecurityContextToHCL(r.SecurityContext); v != "" {
		outputConfig += fmt.Sprintf("\tsecurity_context %s\n", v)
	}
	if r.TerminationMessagePath != nil {
		outputConfig += fmt.Sprintf("\ttermination_message_path = %#v\n", *r.TerminationMessagePath)
	}
	if r.TerminationMessagePolicy != nil {
		outputConfig += fmt.Sprintf("\ttermination_message_policy = %#v\n", *r.TerminationMessagePolicy)
	}
	if r.VolumeMounts != nil {
		for _, v := range r.VolumeMounts {
			outputConfig += fmt.Sprintf("\tvolume_mounts %s\n", convertRunServiceSpecTemplateSpecContainersVolumeMountsToHCL(&v))
		}
	}
	if r.WorkingDir != nil {
		outputConfig += fmt.Sprintf("\tworking_dir = %#v\n", *r.WorkingDir)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvToHCL(r *run.ServiceSpecTemplateSpecContainersEnv) string {
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
	if v := convertRunServiceSpecTemplateSpecContainersEnvValueFromToHCL(r.ValueFrom); v != "" {
		outputConfig += fmt.Sprintf("\tvalue_from %s\n", v)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromToHCL(r *run.ServiceSpecTemplateSpecContainersEnvValueFrom) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefToHCL(r.ConfigMapKeyRef); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_map_key_ref %s\n", v)
	}
	if v := convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefToHCL(r.SecretKeyRef); v != "" {
		outputConfig += fmt.Sprintf("\tsecret_key_ref %s\n", v)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefToHCL(r *run.ServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRef) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if v := convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReferenceToHCL(r.LocalObjectReference); v != "" {
		outputConfig += fmt.Sprintf("\tlocal_object_reference %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReferenceToHCL(r *run.ServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReference) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefToHCL(r *run.ServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRef) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if v := convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReferenceToHCL(r.LocalObjectReference); v != "" {
		outputConfig += fmt.Sprintf("\tlocal_object_reference %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReferenceToHCL(r *run.ServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReference) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvFromToHCL(r *run.ServiceSpecTemplateSpecContainersEnvFrom) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefToHCL(r.ConfigMapRef); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_map_ref %s\n", v)
	}
	if r.Prefix != nil {
		outputConfig += fmt.Sprintf("\tprefix = %#v\n", *r.Prefix)
	}
	if v := convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefToHCL(r.SecretRef); v != "" {
		outputConfig += fmt.Sprintf("\tsecret_ref %s\n", v)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefToHCL(r *run.ServiceSpecTemplateSpecContainersEnvFromConfigMapRef) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReferenceToHCL(r.LocalObjectReference); v != "" {
		outputConfig += fmt.Sprintf("\tlocal_object_reference %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReferenceToHCL(r *run.ServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReference) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefToHCL(r *run.ServiceSpecTemplateSpecContainersEnvFromSecretRef) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReferenceToHCL(r.LocalObjectReference); v != "" {
		outputConfig += fmt.Sprintf("\tlocal_object_reference %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReferenceToHCL(r *run.ServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReference) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeToHCL(r *run.ServiceSpecTemplateSpecContainersLivenessProbe) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersLivenessProbeExecToHCL(r.Exec); v != "" {
		outputConfig += fmt.Sprintf("\texec %s\n", v)
	}
	if r.FailureThreshold != nil {
		outputConfig += fmt.Sprintf("\tfailure_threshold = %#v\n", *r.FailureThreshold)
	}
	if v := convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetToHCL(r.HttpGet); v != "" {
		outputConfig += fmt.Sprintf("\thttp_get %s\n", v)
	}
	if r.InitialDelaySeconds != nil {
		outputConfig += fmt.Sprintf("\tinitial_delay_seconds = %#v\n", *r.InitialDelaySeconds)
	}
	if r.PeriodSeconds != nil {
		outputConfig += fmt.Sprintf("\tperiod_seconds = %#v\n", *r.PeriodSeconds)
	}
	if r.SuccessThreshold != nil {
		outputConfig += fmt.Sprintf("\tsuccess_threshold = %#v\n", *r.SuccessThreshold)
	}
	if v := convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocketToHCL(r.TcpSocket); v != "" {
		outputConfig += fmt.Sprintf("\ttcp_socket %s\n", v)
	}
	if r.TimeoutSeconds != nil {
		outputConfig += fmt.Sprintf("\ttimeout_seconds = %#v\n", *r.TimeoutSeconds)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeExecToHCL(r *run.ServiceSpecTemplateSpecContainersLivenessProbeExec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Command != nil {
		outputConfig += fmt.Sprintf("\tcommand = %#v\n", *r.Command)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetToHCL(r *run.ServiceSpecTemplateSpecContainersLivenessProbeHttpGet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Host != nil {
		outputConfig += fmt.Sprintf("\thost = %#v\n", *r.Host)
	}
	if r.HttpHeaders != nil {
		for _, v := range r.HttpHeaders {
			outputConfig += fmt.Sprintf("\thttp_headers %s\n", convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeadersToHCL(&v))
		}
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.Scheme != nil {
		outputConfig += fmt.Sprintf("\tscheme = %#v\n", *r.Scheme)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeadersToHCL(r *run.ServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeaders) string {
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

func convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocketToHCL(r *run.ServiceSpecTemplateSpecContainersLivenessProbeTcpSocket) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Host != nil {
		outputConfig += fmt.Sprintf("\thost = %#v\n", *r.Host)
	}
	if r.Port != nil {
		outputConfig += fmt.Sprintf("\tport = %#v\n", *r.Port)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersPortsToHCL(r *run.ServiceSpecTemplateSpecContainersPorts) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ContainerPort != nil {
		outputConfig += fmt.Sprintf("\tcontainer_port = %#v\n", *r.ContainerPort)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Protocol != nil {
		outputConfig += fmt.Sprintf("\tprotocol = %#v\n", *r.Protocol)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeToHCL(r *run.ServiceSpecTemplateSpecContainersReadinessProbe) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecContainersReadinessProbeExecToHCL(r.Exec); v != "" {
		outputConfig += fmt.Sprintf("\texec %s\n", v)
	}
	if r.FailureThreshold != nil {
		outputConfig += fmt.Sprintf("\tfailure_threshold = %#v\n", *r.FailureThreshold)
	}
	if v := convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetToHCL(r.HttpGet); v != "" {
		outputConfig += fmt.Sprintf("\thttp_get %s\n", v)
	}
	if r.InitialDelaySeconds != nil {
		outputConfig += fmt.Sprintf("\tinitial_delay_seconds = %#v\n", *r.InitialDelaySeconds)
	}
	if r.PeriodSeconds != nil {
		outputConfig += fmt.Sprintf("\tperiod_seconds = %#v\n", *r.PeriodSeconds)
	}
	if r.SuccessThreshold != nil {
		outputConfig += fmt.Sprintf("\tsuccess_threshold = %#v\n", *r.SuccessThreshold)
	}
	if v := convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocketToHCL(r.TcpSocket); v != "" {
		outputConfig += fmt.Sprintf("\ttcp_socket %s\n", v)
	}
	if r.TimeoutSeconds != nil {
		outputConfig += fmt.Sprintf("\ttimeout_seconds = %#v\n", *r.TimeoutSeconds)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeExecToHCL(r *run.ServiceSpecTemplateSpecContainersReadinessProbeExec) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Command != nil {
		outputConfig += fmt.Sprintf("\tcommand = %#v\n", *r.Command)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetToHCL(r *run.ServiceSpecTemplateSpecContainersReadinessProbeHttpGet) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Host != nil {
		outputConfig += fmt.Sprintf("\thost = %#v\n", *r.Host)
	}
	if r.HttpHeaders != nil {
		for _, v := range r.HttpHeaders {
			outputConfig += fmt.Sprintf("\thttp_headers %s\n", convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeadersToHCL(&v))
		}
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	if r.Scheme != nil {
		outputConfig += fmt.Sprintf("\tscheme = %#v\n", *r.Scheme)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeadersToHCL(r *run.ServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeaders) string {
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

func convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocketToHCL(r *run.ServiceSpecTemplateSpecContainersReadinessProbeTcpSocket) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Host != nil {
		outputConfig += fmt.Sprintf("\thost = %#v\n", *r.Host)
	}
	if r.Port != nil {
		outputConfig += fmt.Sprintf("\tport = %#v\n", *r.Port)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersResourcesToHCL(r *run.ServiceSpecTemplateSpecContainersResources) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersSecurityContextToHCL(r *run.ServiceSpecTemplateSpecContainersSecurityContext) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.RunAsUser != nil {
		outputConfig += fmt.Sprintf("\trun_as_user = %#v\n", *r.RunAsUser)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecContainersVolumeMountsToHCL(r *run.ServiceSpecTemplateSpecContainersVolumeMounts) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.MountPath != nil {
		outputConfig += fmt.Sprintf("\tmount_path = %#v\n", *r.MountPath)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.ReadOnly != nil {
		outputConfig += fmt.Sprintf("\tread_only = %#v\n", *r.ReadOnly)
	}
	if r.SubPath != nil {
		outputConfig += fmt.Sprintf("\tsub_path = %#v\n", *r.SubPath)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecVolumesToHCL(r *run.ServiceSpecTemplateSpecVolumes) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if v := convertRunServiceSpecTemplateSpecVolumesConfigMapToHCL(r.ConfigMap); v != "" {
		outputConfig += fmt.Sprintf("\tconfig_map %s\n", v)
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if v := convertRunServiceSpecTemplateSpecVolumesSecretToHCL(r.Secret); v != "" {
		outputConfig += fmt.Sprintf("\tsecret %s\n", v)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecVolumesConfigMapToHCL(r *run.ServiceSpecTemplateSpecVolumesConfigMap) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DefaultMode != nil {
		outputConfig += fmt.Sprintf("\tdefault_mode = %#v\n", *r.DefaultMode)
	}
	if r.Items != nil {
		for _, v := range r.Items {
			outputConfig += fmt.Sprintf("\titems %s\n", convertRunServiceSpecTemplateSpecVolumesConfigMapItemsToHCL(&v))
		}
	}
	if r.Name != nil {
		outputConfig += fmt.Sprintf("\tname = %#v\n", *r.Name)
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecVolumesConfigMapItemsToHCL(r *run.ServiceSpecTemplateSpecVolumesConfigMapItems) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Mode != nil {
		outputConfig += fmt.Sprintf("\tmode = %#v\n", *r.Mode)
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecVolumesSecretToHCL(r *run.ServiceSpecTemplateSpecVolumesSecret) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.DefaultMode != nil {
		outputConfig += fmt.Sprintf("\tdefault_mode = %#v\n", *r.DefaultMode)
	}
	if r.Items != nil {
		for _, v := range r.Items {
			outputConfig += fmt.Sprintf("\titems %s\n", convertRunServiceSpecTemplateSpecVolumesSecretItemsToHCL(&v))
		}
	}
	if r.Optional != nil {
		outputConfig += fmt.Sprintf("\toptional = %#v\n", *r.Optional)
	}
	if r.SecretName != nil {
		outputConfig += fmt.Sprintf("\tsecret_name = %#v\n", *r.SecretName)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTemplateSpecVolumesSecretItemsToHCL(r *run.ServiceSpecTemplateSpecVolumesSecretItems) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.Key != nil {
		outputConfig += fmt.Sprintf("\tkey = %#v\n", *r.Key)
	}
	if r.Mode != nil {
		outputConfig += fmt.Sprintf("\tmode = %#v\n", *r.Mode)
	}
	if r.Path != nil {
		outputConfig += fmt.Sprintf("\tpath = %#v\n", *r.Path)
	}
	return outputConfig + "}"
}

func convertRunServiceSpecTrafficToHCL(r *run.ServiceSpecTraffic) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	if r.ConfigurationName != nil {
		outputConfig += fmt.Sprintf("\tconfiguration_name = %#v\n", *r.ConfigurationName)
	}
	if r.LatestRevision != nil {
		outputConfig += fmt.Sprintf("\tlatest_revision = %#v\n", *r.LatestRevision)
	}
	if r.Percent != nil {
		outputConfig += fmt.Sprintf("\tpercent = %#v\n", *r.Percent)
	}
	if r.RevisionName != nil {
		outputConfig += fmt.Sprintf("\trevision_name = %#v\n", *r.RevisionName)
	}
	if r.Tag != nil {
		outputConfig += fmt.Sprintf("\ttag = %#v\n", *r.Tag)
	}
	return outputConfig + "}"
}

func convertRunServiceStatusToHCL(r *run.ServiceStatus) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertRunServiceStatusAddressToHCL(r *run.ServiceStatusAddress) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertRunServiceStatusConditionsToHCL(r *run.ServiceStatusConditions) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertRunServiceStatusConditionsLastTransitionTimeToHCL(r *run.ServiceStatusConditionsLastTransitionTime) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
}

func convertRunServiceStatusTrafficToHCL(r *run.ServiceStatusTraffic) string {
	if r == nil {
		return ""
	}
	outputConfig := "{\n"
	return outputConfig + "}"
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
		"cloudFunction": in["cloud_function"],
		"cloudRun":      convertEventarcTriggerDestinationCloudRun(in["cloud_run_service"]),
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

func convertEventarcTriggerDestinationCloudRun(i interface{}) map[string]interface{} {
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

func convertEventarcTriggerDestinationCloudRunList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerDestinationCloudRun(v))
	}
	return out
}

func convertEventarcTriggerEventFilters(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"attribute": in["attribute"],
		"value":     in["value"],
	}
}

func convertEventarcTriggerEventFiltersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertEventarcTriggerEventFilters(v))
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

func convertRunServiceMetadata(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"annotations":                in["annotations"],
		"clusterName":                in["cluster_name"],
		"createTime":                 convertRunServiceMetadataCreateTime(in["create_time"]),
		"deleteTime":                 convertRunServiceMetadataDeleteTime(in["delete_time"]),
		"deletionGracePeriodSeconds": in["deletion_grace_period_seconds"],
		"finalizers":                 in["finalizers"],
		"generateName":               in["generate_name"],
		"generation":                 in["generation"],
		"labels":                     in["labels"],
		"name":                       in["name"],
		"namespace":                  in["namespace"],
		"ownerReferences":            in["owner_references"],
		"resourceVersion":            in["resource_version"],
		"selfLink":                   in["self_link"],
		"uid":                        in["uid"],
	}
}

func convertRunServiceMetadataList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceMetadata(v))
	}
	return out
}

func convertRunServiceMetadataCreateTime(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nanos":   in["nanos"],
		"seconds": in["seconds"],
	}
}

func convertRunServiceMetadataCreateTimeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceMetadataCreateTime(v))
	}
	return out
}

func convertRunServiceMetadataDeleteTime(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nanos":   in["nanos"],
		"seconds": in["seconds"],
	}
}

func convertRunServiceMetadataDeleteTimeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceMetadataDeleteTime(v))
	}
	return out
}

func convertRunServiceMetadataOwnerReferences(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"apiVersion":         in["api_version"],
		"blockOwnerDeletion": in["block_owner_deletion"],
		"controller":         in["controller"],
		"kind":               in["kind"],
		"name":               in["name"],
		"uid":                in["uid"],
	}
}

func convertRunServiceMetadataOwnerReferencesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceMetadataOwnerReferences(v))
	}
	return out
}

func convertRunServiceSpec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"template": convertRunServiceSpecTemplate(in["template"]),
		"traffic":  in["traffic"],
	}
}

func convertRunServiceSpecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpec(v))
	}
	return out
}

func convertRunServiceSpecTemplate(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"metadata": convertRunServiceSpecTemplateMetadata(in["metadata"]),
		"spec":     convertRunServiceSpecTemplateSpec(in["spec"]),
	}
}

func convertRunServiceSpecTemplateList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplate(v))
	}
	return out
}

func convertRunServiceSpecTemplateMetadata(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"annotations":                in["annotations"],
		"clusterName":                in["cluster_name"],
		"createTime":                 convertRunServiceSpecTemplateMetadataCreateTime(in["create_time"]),
		"deleteTime":                 convertRunServiceSpecTemplateMetadataDeleteTime(in["delete_time"]),
		"deletionGracePeriodSeconds": in["deletion_grace_period_seconds"],
		"finalizers":                 in["finalizers"],
		"generateName":               in["generate_name"],
		"generation":                 in["generation"],
		"labels":                     in["labels"],
		"name":                       in["name"],
		"namespace":                  in["namespace"],
		"ownerReferences":            in["owner_references"],
		"resourceVersion":            in["resource_version"],
		"selfLink":                   in["self_link"],
		"uid":                        in["uid"],
	}
}

func convertRunServiceSpecTemplateMetadataList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateMetadata(v))
	}
	return out
}

func convertRunServiceSpecTemplateMetadataCreateTime(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nanos":   in["nanos"],
		"seconds": in["seconds"],
	}
}

func convertRunServiceSpecTemplateMetadataCreateTimeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateMetadataCreateTime(v))
	}
	return out
}

func convertRunServiceSpecTemplateMetadataDeleteTime(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nanos":   in["nanos"],
		"seconds": in["seconds"],
	}
}

func convertRunServiceSpecTemplateMetadataDeleteTimeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateMetadataDeleteTime(v))
	}
	return out
}

func convertRunServiceSpecTemplateMetadataOwnerReferences(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"apiVersion":         in["api_version"],
		"blockOwnerDeletion": in["block_owner_deletion"],
		"controller":         in["controller"],
		"kind":               in["kind"],
		"name":               in["name"],
		"uid":                in["uid"],
	}
}

func convertRunServiceSpecTemplateMetadataOwnerReferencesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateMetadataOwnerReferences(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"containerConcurrency": in["container_concurrency"],
		"containers":           in["containers"],
		"serviceAccountName":   in["service_account_name"],
		"timeoutSeconds":       in["timeout_seconds"],
		"volumes":              in["volumes"],
	}
}

func convertRunServiceSpecTemplateSpecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpec(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainers(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"args":                     in["args"],
		"command":                  in["command"],
		"env":                      in["env"],
		"envFrom":                  in["env_from"],
		"image":                    in["image"],
		"imagePullPolicy":          in["image_pull_policy"],
		"livenessProbe":            convertRunServiceSpecTemplateSpecContainersLivenessProbe(in["liveness_probe"]),
		"name":                     in["name"],
		"ports":                    in["ports"],
		"readinessProbe":           convertRunServiceSpecTemplateSpecContainersReadinessProbe(in["readiness_probe"]),
		"resources":                convertRunServiceSpecTemplateSpecContainersResources(in["resources"]),
		"securityContext":          convertRunServiceSpecTemplateSpecContainersSecurityContext(in["security_context"]),
		"terminationMessagePath":   in["termination_message_path"],
		"terminationMessagePolicy": in["termination_message_policy"],
		"volumeMounts":             in["volume_mounts"],
		"workingDir":               in["working_dir"],
	}
}

func convertRunServiceSpecTemplateSpecContainersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainers(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnv(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name":      in["name"],
		"value":     in["value"],
		"valueFrom": convertRunServiceSpecTemplateSpecContainersEnvValueFrom(in["value_from"]),
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnv(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFrom(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configMapKeyRef": convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRef(in["config_map_key_ref"]),
		"secretKeyRef":    convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRef(in["secret_key_ref"]),
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvValueFrom(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRef(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"key":                  in["key"],
		"localObjectReference": convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReference(in["local_object_reference"]),
		"name":                 in["name"],
		"optional":             in["optional"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRef(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReference(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReferenceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvValueFromConfigMapKeyRefLocalObjectReference(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRef(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"key":                  in["key"],
		"localObjectReference": convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReference(in["local_object_reference"]),
		"name":                 in["name"],
		"optional":             in["optional"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRef(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReference(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReferenceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvValueFromSecretKeyRefLocalObjectReference(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvFrom(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configMapRef": convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRef(in["config_map_ref"]),
		"prefix":       in["prefix"],
		"secretRef":    convertRunServiceSpecTemplateSpecContainersEnvFromSecretRef(in["secret_ref"]),
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvFromList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvFrom(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRef(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"localObjectReference": convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReference(in["local_object_reference"]),
		"name":                 in["name"],
		"optional":             in["optional"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRef(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReference(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReferenceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvFromConfigMapRefLocalObjectReference(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRef(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"localObjectReference": convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReference(in["local_object_reference"]),
		"name":                 in["name"],
		"optional":             in["optional"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvFromSecretRef(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReference(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name": in["name"],
	}
}

func convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReferenceList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersEnvFromSecretRefLocalObjectReference(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbe(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"exec":                convertRunServiceSpecTemplateSpecContainersLivenessProbeExec(in["exec"]),
		"failureThreshold":    in["failure_threshold"],
		"httpGet":             convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGet(in["http_get"]),
		"initialDelaySeconds": in["initial_delay_seconds"],
		"periodSeconds":       in["period_seconds"],
		"successThreshold":    in["success_threshold"],
		"tcpSocket":           convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocket(in["tcp_socket"]),
		"timeoutSeconds":      in["timeout_seconds"],
	}
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersLivenessProbe(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeExec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"command": in["command"],
	}
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeExecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersLivenessProbeExec(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"host":        in["host"],
		"httpHeaders": in["http_headers"],
		"path":        in["path"],
		"scheme":      in["scheme"],
	}
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGet(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeaders(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name":  in["name"],
		"value": in["value"],
	}
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeadersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersLivenessProbeHttpGetHttpHeaders(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocket(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"host": in["host"],
		"port": in["port"],
	}
}

func convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocketList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersLivenessProbeTcpSocket(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersPorts(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"containerPort": in["container_port"],
		"name":          in["name"],
		"protocol":      in["protocol"],
	}
}

func convertRunServiceSpecTemplateSpecContainersPortsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersPorts(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbe(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"exec":                convertRunServiceSpecTemplateSpecContainersReadinessProbeExec(in["exec"]),
		"failureThreshold":    in["failure_threshold"],
		"httpGet":             convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGet(in["http_get"]),
		"initialDelaySeconds": in["initial_delay_seconds"],
		"periodSeconds":       in["period_seconds"],
		"successThreshold":    in["success_threshold"],
		"tcpSocket":           convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocket(in["tcp_socket"]),
		"timeoutSeconds":      in["timeout_seconds"],
	}
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersReadinessProbe(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeExec(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"command": in["command"],
	}
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeExecList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersReadinessProbeExec(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGet(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"host":        in["host"],
		"httpHeaders": in["http_headers"],
		"path":        in["path"],
		"scheme":      in["scheme"],
	}
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGet(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeaders(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"name":  in["name"],
		"value": in["value"],
	}
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeadersList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersReadinessProbeHttpGetHttpHeaders(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocket(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"host": in["host"],
		"port": in["port"],
	}
}

func convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocketList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersReadinessProbeTcpSocket(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersResources(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"limits":   in["limits"],
		"requests": in["requests"],
	}
}

func convertRunServiceSpecTemplateSpecContainersResourcesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersResources(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersSecurityContext(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"runAsUser": in["run_as_user"],
	}
}

func convertRunServiceSpecTemplateSpecContainersSecurityContextList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersSecurityContext(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecContainersVolumeMounts(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"mountPath": in["mount_path"],
		"name":      in["name"],
		"readOnly":  in["read_only"],
		"subPath":   in["sub_path"],
	}
}

func convertRunServiceSpecTemplateSpecContainersVolumeMountsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecContainersVolumeMounts(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecVolumes(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configMap": convertRunServiceSpecTemplateSpecVolumesConfigMap(in["config_map"]),
		"name":      in["name"],
		"secret":    convertRunServiceSpecTemplateSpecVolumesSecret(in["secret"]),
	}
}

func convertRunServiceSpecTemplateSpecVolumesList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecVolumes(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecVolumesConfigMap(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"defaultMode": in["default_mode"],
		"items":       in["items"],
		"name":        in["name"],
		"optional":    in["optional"],
	}
}

func convertRunServiceSpecTemplateSpecVolumesConfigMapList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecVolumesConfigMap(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecVolumesConfigMapItems(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"key":  in["key"],
		"mode": in["mode"],
		"path": in["path"],
	}
}

func convertRunServiceSpecTemplateSpecVolumesConfigMapItemsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecVolumesConfigMapItems(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecVolumesSecret(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"defaultMode": in["default_mode"],
		"items":       in["items"],
		"optional":    in["optional"],
		"secretName":  in["secret_name"],
	}
}

func convertRunServiceSpecTemplateSpecVolumesSecretList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecVolumesSecret(v))
	}
	return out
}

func convertRunServiceSpecTemplateSpecVolumesSecretItems(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"key":  in["key"],
		"mode": in["mode"],
		"path": in["path"],
	}
}

func convertRunServiceSpecTemplateSpecVolumesSecretItemsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTemplateSpecVolumesSecretItems(v))
	}
	return out
}

func convertRunServiceSpecTraffic(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configurationName": in["configuration_name"],
		"latestRevision":    in["latest_revision"],
		"percent":           in["percent"],
		"revisionName":      in["revision_name"],
		"tag":               in["tag"],
		"url":               in["url"],
	}
}

func convertRunServiceSpecTrafficList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceSpecTraffic(v))
	}
	return out
}

func convertRunServiceStatus(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"address":                   convertRunServiceStatusAddress(in["address"]),
		"conditions":                in["conditions"],
		"latestCreatedRevisionName": in["latest_created_revision_name"],
		"latestReadyRevisionName":   in["latest_ready_revision_name"],
		"observedGeneration":        in["observed_generation"],
		"traffic":                   in["traffic"],
		"url":                       in["url"],
	}
}

func convertRunServiceStatusList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceStatus(v))
	}
	return out
}

func convertRunServiceStatusAddress(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"url": in["url"],
	}
}

func convertRunServiceStatusAddressList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceStatusAddress(v))
	}
	return out
}

func convertRunServiceStatusConditions(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"lastTransitionTime": convertRunServiceStatusConditionsLastTransitionTime(in["last_transition_time"]),
		"message":            in["message"],
		"reason":             in["reason"],
		"severity":           in["severity"],
		"status":             in["status"],
		"type":               in["type"],
	}
}

func convertRunServiceStatusConditionsList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceStatusConditions(v))
	}
	return out
}

func convertRunServiceStatusConditionsLastTransitionTime(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"nanos":   in["nanos"],
		"seconds": in["seconds"],
	}
}

func convertRunServiceStatusConditionsLastTransitionTimeList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceStatusConditionsLastTransitionTime(v))
	}
	return out
}

func convertRunServiceStatusTraffic(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}
	in := i.(map[string]interface{})
	return map[string]interface{}{
		"configurationName": in["configuration_name"],
		"latestRevision":    in["latest_revision"],
		"percent":           in["percent"],
		"revisionName":      in["revision_name"],
		"tag":               in["tag"],
		"url":               in["url"],
	}
}

func convertRunServiceStatusTrafficList(i interface{}) (out []map[string]interface{}) {
	if i == nil {
		return nil
	}

	for _, v := range i.([]interface{}) {
		out = append(out, convertRunServiceStatusTraffic(v))
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
