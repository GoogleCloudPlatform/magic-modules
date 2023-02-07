package google

// import (
// 	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"google.golang.org/api/cloudresourcemanager/v1"
// )

func resourceConverterFolder() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderCaiObject,
		MergeCreateUpdate: MergeFolder,
	}
}

func GetFolderCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/folders/{folder_id}") //check
	if err != nil {
		return []Asset{{
			Name: "name",
		}}, nil
	}
	// if obj, err := GetFolderApiObject(d, config); err == nil {
	// 	return []Asset{{
	// 		Name: "name",
	// 		Type: "cloudresourcemanager.googleapis.com/Folder",
	// 		// Resource: &AssetResource{
	// 		// 	Version:              "v1",
	// 		// 	DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
	// 		// 	DiscoveryName:        "Folder",
	// 		// 	Data:                 obj,
	// 		// },
	// 	}}, nil
	// } else {
	// 	return []Asset{}, err
	// }
}

// func GetFolderApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {

// 	folder := &Folder{
// 		Name:      d.Get("name").(string),
// 		Parent:	d.Get("parent").(string),
// 		DisplayName: d.Get("display_name").(string),
// 		State: d.Get("lifecycle_state").(string), //check
// 		//CreateTime: constructTime(d.Get("create_time").(string)),
// 	}

// 	return jsonMap(folder)
// }

// t,_:= time.Parse(time.RFC3339, "2021-08-02T06:07:23.051-07:00Z")
//     fmt.Println(t)
//     fmt.Println(t.Unix())
//     fmt.Println(t.UnixNano())

// func constructTime(create_time string) *time {
// 	t,_:= time.Parse(time.RFC3339, create_time)
// 	return &Time{
		
// 	}
// }




func MergeFolder(existing, incoming Asset) Asset {
	existing.Resource = incoming.Resource
	return existing
}
