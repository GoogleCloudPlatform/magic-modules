package google

import (
	"reflect"
)

func GetBigQueryTableCaiObject(d TerraformResourceData, config *Config) (Asset, error) {
	name, err := assetName(d, config, "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")

	if err != nil {
		return Asset{}, err
	}
	if obj, err := GetBigQueryTableApiObject(d, config); err == nil {
		return Asset{
			Name: name,
			Type: "bigquery.googleapis.com/Table",
			Resource: &AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/bigquery/v2/rest",
				DiscoveryName:        "Table",
				Data:                 obj,
			},
		}, nil
	} else {
		return Asset{}, err
	}
}

func GetBigQueryTableApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	tableReferenceProp, err := expandBigQueryTableTableReference(nil, d, config)

	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("table_reference"); !isEmptyValue(reflect.ValueOf(tableReferenceProp)) && (ok || !reflect.DeepEqual(v, tableReferenceProp)) {
		obj["tableReference"] = tableReferenceProp
	}

    descriptionProp, err := expandBigQueryTableDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	friendlyNameProp, err := expandBigQueryTableFriendlyName(d.Get("friendly_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("friendly_name"); !isEmptyValue(reflect.ValueOf(friendlyNameProp)) && (ok || !reflect.DeepEqual(v, friendlyNameProp)) {
		obj["friendlyName"] = friendlyNameProp
	}
	
	labelsProp, err := expandBigQueryTableLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	locationProp, err := expandBigQueryTableLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !isEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}
	
	expirationTimeProp, err := expandBigQueryTableExpirationTime(nil, d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("expiration_time"); !isEmptyValue(reflect.ValueOf(expirationTimeProp)) && (ok || !reflect.DeepEqual(v, expirationTimeProp)) {
		obj["expirationTime"] = expirationTimeProp
	}

	encryptionConfigurationProp, err := expandBigQueryTableEncyptionConfiguration(d.Get("encryption_configuration"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("encryption_configuration"); !isEmptyValue(reflect.ValueOf(encryptionConfigurationProp)) && (ok || !reflect.DeepEqual(v, encryptionConfigurationProp)) {
		obj["encryptionConfiguration"] = encryptionConfigurationProp
	}

	timePartitionProp, err := expandBigQueryTableTimePartition(d.Get("time_partitioning"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("time_partitioning"); !isEmptyValue(reflect.ValueOf(timePartitionProp)) && (ok || !reflect.DeepEqual(v, timePartitionProp)) {
		obj["timePartitioning"] = timePartitionProp
	}

	return obj, nil
}

func expandBigQueryTableTableReference(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	transformed := make(map[string]interface{})
	
	transformedProjectId, err := expandBigQueryTableTableReferenceProjectId(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProjectId); val.IsValid() && !isEmptyValue(val) {
		transformed["project"] = transformedProjectId
	}

	transformedDatasetId, err := expandBigQueryTableTableReferenceDatasetId(d.Get("dataset_id"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDatasetId); val.IsValid() && !isEmptyValue(val) {
		transformed["datasetId"] = transformedDatasetId
	}

	transformedTableId, err := expandBigQueryTableTableReferenceTableId(d.Get("table_id"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTableId); val.IsValid() && !isEmptyValue(val) {
		transformed["tableId"] = transformedDatasetId
	}

	return transformed, nil
}

func expandBigQueryTableExpirationTime(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	transformed := make(map[string]interface{})
	
	transformedExpirationTimeValue, err := expandBigQueryTableExpirationTimeValue(d.Get("expiration_time"), d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpirationTimeValue); val.IsValid() && !isEmptyValue(val) {
		transformed["value"] = transformedExpirationTimeValue
	}

	return transformed, nil
}

func expandBigQueryTableEncyptionConfiguration(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedKmsKeyName, err := expandBigQueryTableEncyptionConfigurationKmsKeyName(original["kms_key_name"], d, config)
	
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedKmsKeyName); val.IsValid() && !isEmptyValue(val) {
		transformed["kmsKeyName"] = transformedKmsKeyName
	}
	return transformed, nil
}

func expandBigQueryTableTableReferenceProjectId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTableReferenceDatasetId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTableReferenceTableId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableExpirationTimeValue(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableEncyptionConfigurationKmsKeyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimePartition(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedType, err := expandBigQueryTableTimePartitionType(original["type"], d, config)
	
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !isEmptyValue(val) {
		transformed["type"] = transformedType
	}

	transformedExpirationMs, err := expandBigQueryTableTimePartitionExpirationMs(original["expiration_ms"], d, config)
	
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpirationMs); val.IsValid() && !isEmptyValue(val) {
		transformed["expirationMs"] = transformedExpirationMs
	}

	transformedField, err := expandBigQueryTableTimeField(original["field"], d, config)
	
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedField); val.IsValid() && !isEmptyValue(val) {
		transformed["field"] = transformedField
	}
	
	transformedRequirePartitionFilter, err := expandBigQueryTableRequirePartitionFilter(original["require_partition_filter"], d, config)
	
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRequirePartitionFilter); val.IsValid() && !isEmptyValue(val) {
		transformed["requirePartitionFilter"] = transformedRequirePartitionFilter
	}
	return transformed, nil
}

func expandBigQueryTableTimePartitionType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimePartitionExpirationMs(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimeField(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	transformed := make(map[string]interface{})
	
	transformedValue, err := expandBigQueryTableTimeFieldValue(v, d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !isEmptyValue(val) {
		transformed["value"] = transformedValue
	}

	return transformed, nil
}

func expandBigQueryTableTimeFieldValue(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableRequirePartitionFilter(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableFriendlyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandBigQueryTableLocation(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
