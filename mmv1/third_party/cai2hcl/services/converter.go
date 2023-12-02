package services

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/compute"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/services/resourcemanager"
)

var uberConverters = []common.UberConverter{
	compute.UberConverter,
	resourcemanager.UberConverter,
}

var UberConverter common.UberConverter

func init() {
	var converterByAssetType = make(map[string]string)
	var converters = make(map[string]common.Converter)

	for _, uberConverter := range uberConverters {
		appendMap(converterByAssetType, uberConverter.ConverterByAssetType)
		appendMap(converters, uberConverter.Converters)
	}

	UberConverter = common.UberConverter{
		ConverterByAssetType: converterByAssetType,
		Converters:           converters,
	}
}

func appendMap[V interface{}](to map[string]V, from map[string]V) {
	for key, val := range from {
		if _, hasKey := to[key]; hasKey {
			panic(fmt.Sprintf("Map keys are not unique. Duplicate: %s", key))
		}

		to[key] = val
	}
}
