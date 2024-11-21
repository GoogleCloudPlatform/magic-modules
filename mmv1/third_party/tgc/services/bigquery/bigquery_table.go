package bigquery

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const BigQueryTableAssetType string = "bigquery.googleapis.com/Table"

func ResourceConverterBigQueryTable() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: BigQueryTableAssetType,
		Convert:   GetBigQueryTableCaiObject,
	}
}

func GetBigQueryTableCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")

	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetBigQueryTableApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: BigQueryTableAssetType,
			Resource: &cai.AssetResource{
				Version:              "v2",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/bigquery/v2/rest",
				DiscoveryName:        "Table",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetBigQueryTableApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	tableReferenceProp, err := expandBigQueryTableTableReference(nil, d, config)

	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("table_reference"); !tpgresource.IsEmptyValue(reflect.ValueOf(tableReferenceProp)) && (ok || !reflect.DeepEqual(v, tableReferenceProp)) {
		obj["tableReference"] = tableReferenceProp
	}

	descriptionProp, err := expandBigQueryTableDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}

	friendlyNameProp, err := expandBigQueryTableFriendlyName(d.Get("friendly_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("friendly_name"); !tpgresource.IsEmptyValue(reflect.ValueOf(friendlyNameProp)) && (ok || !reflect.DeepEqual(v, friendlyNameProp)) {
		obj["friendlyName"] = friendlyNameProp
	}

	labelsProp, err := expandBigQueryTableLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !tpgresource.IsEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	locationProp, err := expandBigQueryTableLocation(d.Get("location"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("location"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationProp)) && (ok || !reflect.DeepEqual(v, locationProp)) {
		obj["location"] = locationProp
	}

	expirationTimeProp, err := expandBigQueryTableExpirationTime(nil, d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("expiration_time"); !tpgresource.IsEmptyValue(reflect.ValueOf(expirationTimeProp)) && (ok || !reflect.DeepEqual(v, expirationTimeProp)) {
		obj["expirationTime"] = expirationTimeProp
	}

	encryptionConfigurationProp, err := expandBigQueryTableEncyptionConfiguration(d.Get("encryption_configuration"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("encryption_configuration"); !tpgresource.IsEmptyValue(reflect.ValueOf(encryptionConfigurationProp)) && (ok || !reflect.DeepEqual(v, encryptionConfigurationProp)) {
		obj["encryptionConfiguration"] = encryptionConfigurationProp
	}

	timePartitionProp, err := expandBigQueryTableTimePartition(d.Get("time_partitioning"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("time_partitioning"); !tpgresource.IsEmptyValue(reflect.ValueOf(timePartitionProp)) && (ok || !reflect.DeepEqual(v, timePartitionProp)) {
		obj["timePartitioning"] = timePartitionProp
	}

	return obj, nil
}

func expandBigQueryTableTableReference(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	transformed := make(map[string]interface{})

	transformedProjectId, err := expandBigQueryTableTableReferenceProjectId(d.Get("project"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProjectId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["project"] = transformedProjectId
	}

	transformedDatasetId, err := expandBigQueryTableTableReferenceDatasetId(d.Get("dataset_id"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDatasetId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["datasetId"] = transformedDatasetId
	}

	transformedTableId, err := expandBigQueryTableTableReferenceTableId(d.Get("table_id"), d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTableId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["tableId"] = transformedTableId
	}

	return transformed, nil
}

func expandBigQueryTableExpirationTime(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	transformed := make(map[string]interface{})

	transformedExpirationTimeValue, err := expandBigQueryTableExpirationTimeValue(d.Get("expiration_time"), d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpirationTimeValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["value"] = transformedExpirationTimeValue
	}

	return transformed, nil
}

func expandBigQueryTableEncyptionConfiguration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
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
	} else if val := reflect.ValueOf(transformedKmsKeyName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["kmsKeyName"] = transformedKmsKeyName
	}
	return transformed, nil
}

func expandBigQueryTableTableReferenceProjectId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTableReferenceDatasetId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTableReferenceTableId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableExpirationTimeValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableEncyptionConfigurationKmsKeyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimePartition(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
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
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	transformedExpirationMs, err := expandBigQueryTableTimePartitionExpirationMs(original["expiration_ms"], d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedExpirationMs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["expirationMs"] = transformedExpirationMs
	}

	transformedField, err := expandBigQueryTableTimeField(original["field"], d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedField); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["field"] = transformedField
	}

	transformedRequirePartitionFilter, err := expandBigQueryTableRequirePartitionFilter(original["require_partition_filter"], d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRequirePartitionFilter); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["requirePartitionFilter"] = transformedRequirePartitionFilter
	}
	return transformed, nil
}

func expandBigQueryTableTimePartitionType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimePartitionExpirationMs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableTimeField(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	transformed := make(map[string]interface{})

	transformedValue, err := expandBigQueryTableTimeFieldValue(v, d, config)

	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedValue); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["value"] = transformedValue
	}

	return transformed, nil
}

func expandBigQueryTableTimeFieldValue(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableRequirePartitionFilter(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableFriendlyName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandBigQueryTableLabels(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandBigQueryTableLocation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}
