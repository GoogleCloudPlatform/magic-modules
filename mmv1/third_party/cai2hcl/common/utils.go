package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	hashicorpcty "github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// Extracts named part from resource url.
func ParseFieldValue(url string, name string) string {
	fragments := strings.Split(url, "/")
	for ix, item := range fragments {
		if item == name && ix+1 < len(fragments) {
			return fragments[ix+1]
		}
	}
	return ""
}

// Decodes the map object into the target struct.
func DecodeJSON(data map[string]interface{}, v interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	return nil
}

// Converts resource from untyped map format to TF JSON.
func MapToCtyValWithSchema(m map[string]interface{}, s map[string]*schema.Schema) (cty.Value, error) {
	b, err := json.Marshal(&m)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error marshaling map as JSON: %v", err)
	}
	ty, err := hashicorpCtyTypeToZclconfCtyType(schema.InternalMap(s).CoreConfigSchema().ImpliedType())
	if err != nil {
		return cty.NilVal, fmt.Errorf("error casting type: %v", err)
	}
	ret, err := ctyjson.Unmarshal(b, ty)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error unmarshaling JSON as cty.Value: %v", err)
	}
	return ret, nil
}

func Convert(assets []*caiasset.Asset, converterNames map[string]string, converterMap map[string]Converter) ([]byte, error) {
	// Group resources from the same tf resource type for convert.
	// tf -> cai has 1:N mappings occasionally
	groups := make(map[string][]*caiasset.Asset)
	for _, asset := range assets {
		name, ok := converterNames[asset.Type]
		if !ok {
			continue
		}
		groups[name] = append(groups[name], asset)
	}

	f := hclwrite.NewFile()
	rootBody := f.Body()
	for name, v := range groups {
		converter, ok := converterMap[name]
		if !ok {
			continue
		}
		items, err := converter.Convert(v)
		if err != nil {
			return nil, err
		}

		for _, resourceBlock := range items {
			hclBlock := rootBody.AppendNewBlock("resource", resourceBlock.Labels)
			if err := hclWriteBlock(resourceBlock.Value, hclBlock.Body()); err != nil {
				return nil, err
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return printer.Format(f.Bytes())
}

func hashicorpCtyTypeToZclconfCtyType(t hashicorpcty.Type) (cty.Type, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return cty.NilType, err
	}
	var ret cty.Type
	if err := json.Unmarshal(b, &ret); err != nil {
		return cty.NilType, err
	}
	return ret, nil
}
