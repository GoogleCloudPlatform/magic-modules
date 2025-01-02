package converters

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/ancestrymanager"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters/cai"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/models"

	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	"github.com/pkg/errors"
)

// Converts the single resource into CAI assets
func ConvertResource(rdList []*models.FakeResourceData, cfg *transport_tpg.Config, am ancestrymanager.AncestryManager) ([]caiasset.Asset, error) {
	if rdList == nil || len(rdList) == 0 {
		return nil, nil
	}

	var assets []caiasset.Asset
	for _, rd := range rdList {
		converter := ConverterMap[rd.Kind()]
		convertedAssets, err := converter.Convert(rd, cfg)
		if err != nil {
			if errors.Cause(err) == cai.ErrNoConversion {
				continue
			}
		}

		// TODO: combine assets and fetch full policy for IAM bindings/members
		// TODO: combine tfplan address

		for _, asset := range convertedAssets {
			asset.TfplanAddress = []string{rd.TfplanAddr()}
			err := am.SetAncestors(rd, cfg, &asset)
			if err != nil {
				return nil, err
			}
			assets = append(assets, asset)
		}
	}

	return assets, nil
}
