package google

import (
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/tpgresource"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func resourceConverterFolder() tpgresource.ResourceConverter {
	return tpgresource.ResourceConverter{
		AssetType: "cloudresourcemanager.googleapis.com/Folder",
		Convert:   GetFolderCaiObject,
	}
}

func GetFolderCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]tpgresource.Asset, error) {
	name, err := tpgresource.AssetName(d, config, "//cloudresourcemanager.googleapis.com/folders/{{folder_id}}")

	if err != nil {
		return []tpgresource.Asset{}, nil
	}

	if obj, err := GetFolderApiObject(d, config); err == nil {
		return []tpgresource.Asset{{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Folder",
			Resource: &tpgresource.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Folder",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []tpgresource.Asset{}, err
	}
}

func GetFolderApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {

	folder := &tpgresource.Folder{
		Name:        d.Get("name").(string),
		Parent:      d.Get("parent").(string),
		DisplayName: d.Get("display_name").(string),
		State:       d.Get("lifecycle_state").(string),
	}

	if v, ok := d.GetOkExists("create_time"); ok {
		folder.CreateTime = constructTime(v.(string))
	}

	return tpgresource.JsonMap(folder)
}

func constructTime(create_time string) *tpgresource.Timestamp {
	if create_time == "" {
		return &tpgresource.Timestamp{}
	}
	t, _ := time.Parse(time.RFC3339, create_time)
	return &tpgresource.Timestamp{
		Seconds: t.Unix(),
		Nanos:   t.UnixNano(),
	}
}
