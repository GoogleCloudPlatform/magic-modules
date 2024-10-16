package resourcemanager

import (
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func ResourceConverterFolder() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: "cloudresourcemanager.googleapis.com/Folder",
		Convert:   GetFolderCaiObject,
	}
}

func GetFolderCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//cloudresourcemanager.googleapis.com/folders/{{folder_id}}")

	if err != nil {
		return []cai.Asset{}, nil
	}

	if obj, err := GetFolderApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Folder",
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Folder",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetFolderApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {

	folder := &cai.Folder{
		Name:        d.Get("name").(string),
		Parent:      d.Get("parent").(string),
		DisplayName: d.Get("display_name").(string),
		State:       d.Get("lifecycle_state").(string),
	}

	if v, ok := d.GetOkExists("create_time"); ok {
		folder.CreateTime = constructTime(v.(string))
	}

	return cai.JsonMap(folder)
}

func constructTime(create_time string) *cai.Timestamp {
	if create_time == "" {
		return &cai.Timestamp{}
	}
	t, _ := time.Parse(time.RFC3339, create_time)
	return &cai.Timestamp{
		Seconds: t.Unix(),
		Nanos:   t.UnixNano(),
	}
}
