package compute

import (
	"encoding/json"
	"errors"
	"sort"

	"google.golang.org/api/googleapi"

	compute "google.golang.org/api/compute/v0.beta"

	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
)

func expandComputeMetadata(m map[string]interface{}) []*compute.MetadataItems {
	metadata := make([]*compute.MetadataItems, 0, len(m))
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// Append new metadata to existing metadata
	for _, key := range keys {
		v := m[key].(string)
		metadata = append(metadata, &compute.MetadataItems{
			Key:   key,
			Value: &v,
		})
	}

	return metadata
}

func resourceInstanceMetadata(d tpgresource.TerraformResourceData) (*compute.Metadata, error) {
	m := &compute.Metadata{}
	mdMap := d.Get("metadata").(map[string]interface{})
	if v, ok := d.GetOk("metadata_startup_script"); ok && v.(string) != "" {
		if w, ok := mdMap["startup-script"]; ok {
			// metadata.startup-script could be from metadata_startup_script in the first place
			if v != w {
				return nil, errors.New("Cannot provide both metadata_startup_script and metadata.startup-script.")
			}
		}
		mdMap["startup-script"] = v
	}
	if len(mdMap) > 0 {
		m.Items = make([]*compute.MetadataItems, 0, len(mdMap))
		var keys []string
		for k := range mdMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := mdMap[k].(string)
			m.Items = append(m.Items, &compute.MetadataItems{
				Key:   k,
				Value: &v,
			})
		}

		// Set the fingerprint. If the metadata has never been set before
		// then this will just be blank.
		m.Fingerprint = d.Get("metadata_fingerprint").(string)
	}

	return m, nil
}

func resourceInstancePartnerMetadata(d tpgresource.TerraformResourceData) (map[string]compute.StructuredEntries, error) {
	partnerMetadata := make(map[string]compute.StructuredEntries)
	partnerMetadataMap := d.Get("partner_metadata").(map[string]interface{})
	if len(partnerMetadataMap) > 0 {
		for key, value := range partnerMetadataMap {
			var jsonMap map[string]interface{}
			err := json.Unmarshal([]byte(value.(string)), &jsonMap)
			if err != nil {
				return nil, err
			}
			structuredEntries := jsonMap["entries"].(map[string]interface{})
			structuredEntriesJson, err := json.Marshal(&structuredEntries)
			if err != nil {
				return nil, err
			}
			partnerMetadata[key] = compute.StructuredEntries{
				Entries: googleapi.RawMessage(structuredEntriesJson),
			}
		}
	}
	return partnerMetadata, nil
}

func resourceInstancePatchPartnerMetadata(d tpgresource.TerraformResourceData, currentPartnerMetadata map[string]compute.StructuredEntries) map[string]compute.StructuredEntries {
	partnerMetadata, _ := resourceInstancePartnerMetadata(d)
	for key := range currentPartnerMetadata {
		if _, ok := partnerMetadata[key]; !ok {
			partnerMetadata[key] = compute.StructuredEntries{}
		}
	}
	return partnerMetadata

}
