package google

import "time"

func resourceConverterFolder() ResourceConverter {
	return ResourceConverter{
		AssetType:         "cloudresourcemanager.googleapis.com/Folder",
		Convert:           GetFolderCaiObject,
	}
}

func GetFolderCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//cloudresourcemanager.googleapis.com/folders/{{folder_id}}")
	
	if err != nil {
		return []Asset{}, nil
	}

	if obj, err := GetFolderApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: "cloudresourcemanager.googleapis.com/Folder",
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/compute/v1/rest",
				DiscoveryName:        "Folder",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetFolderApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {

	folder := &Folder{
		Name:        d.Get("name").(string),
		Parent:	     d.Get("parent").(string),
		DisplayName: d.Get("display_name").(string),
		State:       d.Get("lifecycle_state").(string), 
	}

	if v, ok := d.GetOkExists("create_time"); ok {
		folder.CreateTime = constructTime(v.(string))
	}

	return jsonMap(folder)
}

func constructTime(create_time string) *Timestamp{
	if create_time == "" {
		return &Timestamp{}
	}
	t,_:= time.Parse(time.RFC3339, create_time)
	return &Timestamp{
		Seconds: t.Unix(),
		Nanos: t.UnixNano(),
	}
}