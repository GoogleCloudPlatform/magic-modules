package google

import (
	"fmt"
	"log"
	"reflect"
)

const ComputeAttachedDiskAssetType string = "compute.googleapis.com/Disk"

func resourceConverterComputeAttachedDisk() ResourceConverter {
	return ResourceConverter{
		AssetType: ComputeAttachedDiskAssetType,
		Convert:   GetComputeAttachedDiskCaiObject,
	}
}

func GetComputeAttachedDiskCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//compute.googleapis.com/projects/{{project}}/zones/{{zone}}/disks/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetComputeAttachedDiskApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ComputeAttachedDiskAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Disk",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetComputeAttachedDiskApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	descriptionProp, err := expandComputeAttachedDiskDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	nameProp, err := expandComputeAttachedDiskName(d.Get("name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	zoneProp, err := expandComputeAttachedDiskZone(d.Get("zone"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("zone"); !isEmptyValue(reflect.ValueOf(zoneProp)) && (ok || !reflect.DeepEqual(v, zoneProp)) {
		obj["zone"] = zoneProp
	}
	return resourceComputeAttachedDiskEncoder(d, config, obj)
}

func resourceComputeAttachedDiskEncoder(d TerraformResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk("type"); ok {
		log.Printf("[DEBUG] Loading disk type: %s", v.(string))
		diskType, err := readDiskType(config, d, v.(string))
		if err != nil {
			return nil, fmt.Errorf(
				"Error loading disk type '%s': %s",
				v.(string), err)
		}

		obj["type"] = diskType.RelativeLink()
	}

	if v, ok := d.GetOk("image"); ok {
		log.Printf("[DEBUG] Resolving image name: %s", v.(string))
		imageUrl, err := resolveImage(config, project, v.(string), userAgent)
		if err != nil {
			return nil, fmt.Errorf(
				"Error resolving image name '%s': %s",
				v.(string), err)
		}

		obj["sourceImage"] = imageUrl
		log.Printf("[DEBUG] Image name resolved to: %s", imageUrl)
	}

	return obj, nil
}

func expandComputeAttachedDiskDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAttachedDiskName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandComputeAttachedDiskZone(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	f, err := parseGlobalFieldValue("zones", v.(string), "project", d, config, true)
	if err != nil {
		return nil, fmt.Errorf("Invalid value for zone: %s", err)
	}
	return f.RelativeLink(), nil
}
