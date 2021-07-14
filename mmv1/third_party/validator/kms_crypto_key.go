package google

import (
	"fmt"
	"reflect"
	"time"
)

func GetKMSCryptoKeyCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudkms.googleapis.com/{{id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetKMSCryptoKeyApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "cloudkms.googleapis.com/CryptoKey",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/kms/v1/rest",
				DiscoveryName:        "CryptoKey",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetKMSCryptoKeyApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	labelsProp, err := expandKMSCryptoKeyLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	purposeProp, err := expandKMSCryptoKeyPurpose(d.Get("purpose"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("purpose"); !isEmptyValue(reflect.ValueOf(purposeProp)) && (ok || !reflect.DeepEqual(v, purposeProp)) {
		obj["purpose"] = purposeProp
	}
	rotationPeriodProp, err := expandKMSCryptoKeyRotationPeriod(d.Get("rotation_period"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("rotation_period"); !isEmptyValue(reflect.ValueOf(rotationPeriodProp)) && (ok || !reflect.DeepEqual(v, rotationPeriodProp)) {
		obj["rotationPeriod"] = rotationPeriodProp
	}
	versionTemplateProp, err := expandKMSCryptoKeyVersionTemplate(d.Get("version_template"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("version_template"); !isEmptyValue(reflect.ValueOf(versionTemplateProp)) && (ok || !reflect.DeepEqual(v, versionTemplateProp)) {
		obj["versionTemplate"] = versionTemplateProp
	}

	return resourceKMSCryptoKeyEncoder(d, config, obj)
}

func resourceKMSCryptoKeyEncoder(d TerraformResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	// if rotationPeriod is set, nextRotationTime must also be set.
	if d.Get("rotation_period") != "" {
		rotationPeriod := d.Get("rotation_period").(string)
		nextRotation, err := kmsCryptoKeyNextRotation(time.Now(), rotationPeriod)

		if err != nil {
			return nil, fmt.Errorf("Error setting CryptoKey rotation period: %s", err.Error())
		}

		obj["nextRotationTime"] = nextRotation
	}

	// set to false if it is not true explicitly
	if !(d.Get("skip_initial_version_creation").(bool)) {
		if err := d.Set("skip_initial_version_creation", false); err != nil {
			return nil, fmt.Errorf("Error setting skip_initial_version_creation: %s", err)
		}
	}

	return obj, nil
}

func expandKMSCryptoKeyLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandKMSCryptoKeyPurpose(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyRotationPeriod(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyVersionTemplate(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAlgorithm, err := expandKMSCryptoKeyVersionTemplateAlgorithm(original["algorithm"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAlgorithm); val.IsValid() && !isEmptyValue(val) {
		transformed["algorithm"] = transformedAlgorithm
	}

	transformedProtectionLevel, err := expandKMSCryptoKeyVersionTemplateProtectionLevel(original["protection_level"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedProtectionLevel); val.IsValid() && !isEmptyValue(val) {
		transformed["protectionLevel"] = transformedProtectionLevel
	}

	return transformed, nil
}

func expandKMSCryptoKeyVersionTemplateAlgorithm(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandKMSCryptoKeyVersionTemplateProtectionLevel(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
