package tfplan2cai

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/ancestrymanager"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/resolvers"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/transport"
)

// Options struct to avoid updating function signatures all along the pipe.
type Options struct {
	ErrorLogger    *zap.Logger
	Offline        bool
	DefaultProject string
	DefaultRegion  string
	DefaultZone    string
	// UserAgent for all requests (if online)
	UserAgent string
	// Map hierarchy resource (like projects/<number> or folders/<number>)
	// to an ancestry path (like organizations/123/folders/456/projects/789)
	AncestryCache map[string]string
}

// Convert converts terraform json plan to CAI Assets.
func Convert(ctx context.Context, jsonPlan []byte, o *Options) ([]caiasset.Asset, error) {
	if o == nil || o.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	resourceDataMap := resolvers.NewDefaultPreResolver(o.ErrorLogger).Resolve(jsonPlan)

	// TODO: add advanced resolvers for resources

	// Set up config and ancestry manager using the same user agent.
	// Config and ancestry manager are shared among resources.
	cfg, err := transport.NewConfig(ctx, o.DefaultProject, o.DefaultZone, o.DefaultRegion, o.Offline, o.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	ancestryManager, err := ancestrymanager.New(cfg, o.Offline, o.AncestryCache, o.ErrorLogger)
	if err != nil {
		return nil, fmt.Errorf("building ancestry manager: %w", err)
	}

	var assets []caiasset.Asset
	for _, resourceDataList := range resourceDataMap {
		convertedAssets, err := converters.ConvertResource(resourceDataList, cfg, ancestryManager, o.ErrorLogger)
		if err != nil {
			return nil, fmt.Errorf("tfplan2ai converting: %w", err)
		}
		assets = append(assets, convertedAssets...)
	}
	return assets, nil
}
