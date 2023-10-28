package resourcemanager

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
)

var UberConverter = common.UberConverter{
	ConverterByAssetType: map[string]string{
		ProjectAssetType:        "google_project",
		ProjectBillingAssetType: "google_project",
	},
	Converters: common.CreateConverterMap(map[string]common.ConverterFactory{
		"google_project": NewProjectConverter,
	}),
}
