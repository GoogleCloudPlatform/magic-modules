package common

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// HclWriteBlocks prints HCLResourceBlock objects as string.
func HclWriteBlocks(blocks []*HCLResourceBlock) ([]byte, error) {
	f := hclwrite.NewFile()
	rootBody := f.Body()

	for _, resourceBlock := range blocks {
		hclBlock := rootBody.AppendNewBlock("resource", resourceBlock.Labels)
		if err := hclWriteBlock(resourceBlock.Value, hclBlock.Body()); err != nil {
			return nil, err
		}
	}

	return printer.Format(f.Bytes())
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
