package generated

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/generated/converters/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

// Options is a struct for options so that adding new options does not
// require updating function signatures all along the pipe.
type ConvertOptions struct {
	ErrorLogger *zap.Logger
}

// Convert converts Asset into HCL.
func Convert(assets []*caiasset.Asset, converterNames map[string]string, converterMap map[string]common.Converter, options *ConvertOptions) ([]byte, error) {
	if options == nil || options.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

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
		for _, item := range items {
			block := rootBody.AppendNewBlock("resource", item.Labels)
			if err := hclWriteBlock(item.Value, block.Body()); err != nil {
				return nil, err
			}
		}
	}

	t, err := printer.Format(f.Bytes())
	options.ErrorLogger.Debug(string(t))
	return t, err
}

func hclWriteBlock(val cty.Value, body *hclwrite.Body) error {
	if val.IsNull() {
		return nil
	}
	if !val.Type().IsObjectType() {
		return fmt.Errorf("expect object type only, but type = %s", val.Type().FriendlyName())
	}
	it := val.ElementIterator()
	for it.Next() {
		objKey, objVal := it.Element()
		if objVal.IsNull() {
			continue
		}
		objValType := objVal.Type()
		switch {
		case objValType.IsObjectType():
			newBlock := body.AppendNewBlock(objKey.AsString(), nil)
			if err := hclWriteBlock(objVal, newBlock.Body()); err != nil {
				return err
			}
		case objValType.IsCollectionType():
			if objVal.LengthInt() == 0 {
				continue
			}
			// Presumes map should not contain object type.
			if !objValType.IsMapType() && objValType.ElementType().IsObjectType() {
				listIterator := objVal.ElementIterator()
				for listIterator.Next() {
					_, listVal := listIterator.Element()
					subBlock := body.AppendNewBlock(objKey.AsString(), nil)
					if err := hclWriteBlock(listVal, subBlock.Body()); err != nil {
						return err
					}
				}
				continue
			}
			fallthrough
		default:
			if objValType.FriendlyName() == "string" && objVal.AsString() == "" {
				continue
			}
			body.SetAttributeValue(objKey.AsString(), objVal)
		}
	}
	return nil
}
