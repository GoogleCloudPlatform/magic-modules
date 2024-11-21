package dataflow

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const DataflowJobAssetType string = "dataflow.googleapis.com/Job"

func ResourceDataflowJob() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: DataflowJobAssetType,
		Convert:   GetDataflowJobCaiObject,
	}
}

func GetDataflowJobCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//dataflow.googleapis.com/projects/{{project}}/locations/{{location}}/jobs")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetDataflowApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: DataflowJobAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1beta3",
				DiscoveryDocumentURI: "https://dataflow.googleapis.com/$discovery/rest?version=v1b3",
				DiscoveryName:        "Job",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetDataflowApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	nameProp, err := expandDataflowName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	projectIdProp, err := expandDataflowProjectId(d.Get("project_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("project_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(projectIdProp)) && (ok || !reflect.DeepEqual(v, projectIdProp)) {
		obj["projectId"] = projectIdProp
	}
	typeProp, err := expandDataflowType(d.Get("type"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("type"); !tpgresource.IsEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
		obj["type"] = typeProp
	}
	environmentProp, err := expandDataflowEnvironment(d.Get("environment"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("environment"); !tpgresource.IsEmptyValue(reflect.ValueOf(environmentProp)) && (ok || !reflect.DeepEqual(v, environmentProp)) {
		obj["environment"] = environmentProp
	}
	stepsProp, err := expandDataflowSteps(d.Get("steps"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("steps"); !tpgresource.IsEmptyValue(reflect.ValueOf(stepsProp)) && (ok || !reflect.DeepEqual(v, stepsProp)) {
		obj["steps"] = stepsProp
	}
	stepsLocationProp, err := expandDataflowStepsLocation(d.Get("steps_location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("steps_location"); !tpgresource.IsEmptyValue(reflect.ValueOf(stepsLocationProp)) && (ok || !reflect.DeepEqual(v, stepsLocationProp)) {
		obj["stepsLocation"] = stepsLocationProp
	}
	currentStateProp, err := expandDataflowCurrentState(d.Get("current_state"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("current_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(currentStateProp)) && (ok || !reflect.DeepEqual(v, currentStateProp)) {
		obj["currentState"] = currentStateProp
	}
	currentStateTimeProp, err := expandDataflowCurrentStateTime(d.Get("current_state_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("current_state_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(currentStateTimeProp)) && (ok || !reflect.DeepEqual(v, currentStateTimeProp)) {
		obj["currentStateTime"] = currentStateTimeProp
	}
	requestedStateProp, err := expandDataflowRequestedState(d.Get("requested_state"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("requested_state"); !tpgresource.IsEmptyValue(reflect.ValueOf(requestedStateProp)) && (ok || !reflect.DeepEqual(v, requestedStateProp)) {
		obj["requestedState"] = requestedStateProp
	}
	executionInfoProp, err := expandDataflowExecutionInfo(d.Get("execution_info"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("execution_info"); !tpgresource.IsEmptyValue(reflect.ValueOf(executionInfoProp)) && (ok || !reflect.DeepEqual(v, executionInfoProp)) {
		obj["executionInfo"] = executionInfoProp
	}
	createTimeProp, err := expandDataflowCreateTime(d.Get("create_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("create_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(createTimeProp)) && (ok || !reflect.DeepEqual(v, createTimeProp)) {
		obj["createTime"] = createTimeProp
	}
	replaceJobIdProp, err := expandDataflowReplaceJobId(d.Get("replace_job_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("replace_job_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(replaceJobIdProp)) && (ok || !reflect.DeepEqual(v, replaceJobIdProp)) {
		obj["replaceJobId"] = replaceJobIdProp
	}
	transformNameMappingProp, err := expandDataflowTransformNameMapping(d.Get("transform_name_mapping"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("transform_name_mapping"); !tpgresource.IsEmptyValue(reflect.ValueOf(transformNameMappingProp)) && (ok || !reflect.DeepEqual(v, transformNameMappingProp)) {
		obj["transformNameMapping"] = transformNameMappingProp
	}
	clientRequestIdProp, err := expandDataflowClientRequestId(d.Get("client_request_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("client_request_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(clientRequestIdProp)) && (ok || !reflect.DeepEqual(v, clientRequestIdProp)) {
		obj["clientRequestId"] = clientRequestIdProp
	}
	replacedByJobIdProp, err := expandDataflowReplacedByJobId(d.Get("replaced_by_job_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("replaced_by_job_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(replacedByJobIdProp)) && (ok || !reflect.DeepEqual(v, replacedByJobIdProp)) {
		obj["replacedByJobId"] = replacedByJobIdProp
	}
	tempFilesProp, err := expandDataflowTempFiles(d.Get("temp_files"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("temp_files"); !tpgresource.IsEmptyValue(reflect.ValueOf(tempFilesProp)) && (ok || !reflect.DeepEqual(v, tempFilesProp)) {
		obj["tempFiles"] = tempFilesProp
	}
	labelsProp, err := expandDataflowLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	locationProp, err := expandDataflowLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}
	pipelineDescriptionProp, err := expandDataflowPipelineDescription(d.Get("pipeline_description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("pipeline_description"); !tpgresource.IsEmptyValue(reflect.ValueOf(pipelineDescriptionProp)) && (ok || !reflect.DeepEqual(v, pipelineDescriptionProp)) {
		obj["pipelineDescription"] = pipelineDescriptionProp
	}
	stageStatesProp, err := expandDataflowStageStates(d.Get("stage_states"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("stage_states"); !tpgresource.IsEmptyValue(reflect.ValueOf(stageStatesProp)) && (ok || !reflect.DeepEqual(v, stageStatesProp)) {
		obj["stageStates"] = stageStatesProp
	}
	jobMetadataProp, err := expandDataflowJobMetadata(d.Get("job_metadata"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("job_metadata"); !tpgresource.IsEmptyValue(reflect.ValueOf(jobMetadataProp)) && (ok || !reflect.DeepEqual(v, jobMetadataProp)) {
		obj["jobMetadata"] = jobMetadataProp
	}
	startTimeProp, err := expandDataflowStartTime(d.Get("start_time"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("start_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(startTimeProp)) && (ok || !reflect.DeepEqual(v, startTimeProp)) {
		obj["startTime"] = startTimeProp
	}
	createdFromSnapshotIdProp, err := expandDataflowCreatedFromSnapshotId(d.Get("created_from_snapshot_id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("created_from_snapshot_id"); !tpgresource.IsEmptyValue(reflect.ValueOf(createdFromSnapshotIdProp)) && (ok || !reflect.DeepEqual(v, createdFromSnapshotIdProp)) {
		obj["createdFromSnapshotId"] = createdFromSnapshotIdProp
	}
	satisfiesPzsProp, err := expandDataflowSatisfiesPzs(d.Get("satisfies_pzs"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("satisfies_pzs"); !tpgresource.IsEmptyValue(reflect.ValueOf(satisfiesPzsProp)) && (ok || !reflect.DeepEqual(v, satisfiesPzsProp)) {
		obj["satisfiesPzs"] = satisfiesPzsProp
	}
	runtimeUpdatableParamsProp, err := expandDataflowRuntimeUpdatableParams(d.Get("runtime_updatable_params"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("runtime_updatable_params"); !tpgresource.IsEmptyValue(reflect.ValueOf(runtimeUpdatableParamsProp)) && (ok || !reflect.DeepEqual(v, runtimeUpdatableParamsProp)) {
		obj["runtimeUpdatableParams"] = runtimeUpdatableParamsProp
	}
	satisfiesPziProp, err := expandDataflowSatisfiesPzi(d.Get("satisfies_pzi"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("satisfies_pzi"); !tpgresource.IsEmptyValue(reflect.ValueOf(satisfiesPziProp)) && (ok || !reflect.DeepEqual(v, satisfiesPziProp)) {
		obj["satisfiesPzi"] = satisfiesPziProp
	}

	return obj, nil
}

func expandDataflowName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowEnvironment(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowSteps(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowStepsLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowCurrentState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowCurrentStateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowRequestedState(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowExecutionInfo(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowCreateTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowReplaceJobId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowTransformNameMapping(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowClientRequestId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowReplacedByJobId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowTempFiles(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowPipelineDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowStageStates(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowJobMetadata(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowStartTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowCreatedFromSnapshotId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowSatisfiesPzs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowRuntimeUpdatableParams(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataflowSatisfiesPzi(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
