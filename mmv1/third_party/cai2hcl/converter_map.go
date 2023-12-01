package cai2hcl

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/services/resourcemanager"
)

var allConverterNames = []map[string]string{
	compute.ConverterNames,
	resourcemanager.ConverterNames,
}

var allConverterMaps = []map[string]common.Converter{
	compute.ConverterMap,
	resourcemanager.ConverterMap,
}

var ConverterNames = joinConverterNames(allConverterNames)
var ConverterMap = joinConverterMaps(allConverterMaps)

func joinConverterNames(arr []map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range arr {
		for key, value := range m {
			if _, hasKey := result[key]; hasKey {
				panic("Converters from different services are not unique")
			}

			result[key] = value
		}
	}

	return result
}

func joinConverterMaps(arr []map[string]common.Converter) map[string]common.Converter {
	result := make(map[string]common.Converter)

	for _, m := range arr {
		for key, value := range m {
			if _, hasKey := result[key]; hasKey {
				panic("Converters from different services are not unique")
			}

			result[key] = value
		}
	}

	return result
}
