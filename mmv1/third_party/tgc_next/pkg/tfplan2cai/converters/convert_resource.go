package converters

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/ancestrymanager"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters/cai"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/models"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/transport"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Converts the single resource into CAI assets
func ConvertResource(rdList []*models.FakeResourceDataWithMeta, cfg *transport_tpg.Config, am ancestrymanager.AncestryManager, errLogger *zap.Logger) ([]caiasset.Asset, error) {
	if rdList == nil || len(rdList) == 0 {
		return nil, nil
	}

	var assets []caiasset.Asset
	for _, rd := range rdList {
		// Skip unsupported resources
		converter, ok := ConverterMap[rd.Kind()]
		if !ok {
			errLogger.Debug(fmt.Sprintf("%s: resource type cannot be converted for CAI-based policies: %s. For details, see https://cloud.google.com/docs/terraform/policy-validation/create-cai-constraints#supported_resources", rd.Address(), rd.Kind()))
			continue
		} else {
			convertedAssets, err := converter.Convert(rd, cfg)
			if err != nil {
				if errors.Cause(err) == cai.ErrNoConversion {
					continue
				}
				return assets, err
			}

			// TODO: combine assets and fetch full policy for IAM bindings/members
			// TODO: combine tfplan address

			for _, asset := range convertedAssets {
				asset.TfplanAddress = []string{rd.Address()}
				err := am.SetAncestors(rd, cfg, &asset)
				if err != nil {
					return nil, err
				}
				assets = append(assets, asset)
			}
		}
	}

	return assets, nil
}
