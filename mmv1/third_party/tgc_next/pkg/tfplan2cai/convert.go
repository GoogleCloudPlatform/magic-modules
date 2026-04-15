package tfplan2cai

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/ancestrymanager"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/resolvers"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/transport"
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
	// HTTPClient for all requests (if online)
	HTTPClient *http.Client
	// Map hierarchy resource (like projects/<number> or folders/<number>)
	// to an ancestry path (like organizations/123/folders/456/projects/789)
	AncestryCache map[string]string

	// If true, the ancestry manager will be a no-op.
	NoOpAncestryManager bool
}

// Convert converts terraform json plan to CAI Assets.
func Convert(ctx context.Context, jsonPlan []byte, o *Options) ([]caiasset.Asset, error) {
	if o == nil || o.ErrorLogger == nil {
		return nil, fmt.Errorf("logger is not initialized")
	}

	// IAM resource resolver, do not run until IAM resources included
	resolvers.NewIamAdvancedResolver(o.ErrorLogger).Resolve(jsonPlan)
	resourceDataMap := resolvers.NewDefaultPreResolver(o.ErrorLogger).Resolve(jsonPlan)

	// TODO: add remaining advanced resolvers for resources
	ParentResolver := resolvers.NewParentResourceResolver(o.ErrorLogger)
	dependencyMap := ParentResolver.Resolve(jsonPlan)

	ParentChildMap := make(map[string][]string)
	for child, attrs := range dependencyMap {
		seenParents := make(map[string]bool)
		for _, parent := range attrs {
			if !seenParents[parent] {
				ParentChildMap[parent] = append(ParentChildMap[parent], child)
				seenParents[parent] = true
			}
		}
	}

	orderMap, err := resolvers.SortTraversalOrder(ParentChildMap)
	if err != nil {
		return nil, fmt.Errorf("sorting traversal order: %w", err)
	}

	// Set up config and ancestry manager using the same user agent.
	// Config and ancestry manager are shared among resources.

	cfg, err := transport.NewConfig(ctx, o.DefaultProject, o.DefaultZone, o.DefaultRegion, o.Offline, o.UserAgent, o.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	var ancestryManager ancestrymanager.AncestryManager
	if !o.NoOpAncestryManager {
		ancestryManager, err = ancestrymanager.New(cfg, o.Offline, o.AncestryCache, o.ErrorLogger)
		if err != nil {
			return nil, fmt.Errorf("building ancestry manager: %w", err)
		}
	} else {
		ancestryManager = &ancestrymanager.NoOpAncestryManager{}
	}

	var assets []caiasset.Asset
	convertedAssetsByAddress := make(map[string][]caiasset.Asset)
	convertedAddresses := make(map[string]bool)

	if orderMap != nil && len(orderMap) > 0 {

		var levels []int
		for level := range orderMap {
			levels = append(levels, level)
		}
		sort.Ints(levels)

		for _, level := range levels {
			addresses := orderMap[level]
			for _, address := range addresses {
				convertedAddresses[address] = true
				resourceDataList := resourceDataMap[address]
				if resourceDataList == nil {
					continue
				}

				deps := dependencyMap[address]
				for _, rd := range resourceDataList {
					if deps != nil {
						for attrName, parentAddr := range deps {
							parentRds := resourceDataMap[parentAddr]
							if len(parentRds) > 0 {
								if parentId := parentRds[0].Id(); parentId != "" {
									rd.Set(attrName, parentId)
								}
							}
						}
					}
				}

				convertedAssets, err := converters.ConvertResource(resourceDataList, cfg, ancestryManager, o.ErrorLogger)
				if err != nil {
					return nil, fmt.Errorf("tfplan2cai converting: %w", err)
				}
				if len(convertedAssets) > 0 {
					parts := strings.SplitN(convertedAssets[0].Name, "/", 4)
					if len(parts) == 4 {
						for _, rd := range resourceDataList {
							rd.SetId(parts[3])
						}
					}
				}
				convertedAssetsByAddress[address] = convertedAssets
				assets = append(assets, convertedAssets...)
			}
		}
	}

	for address, resourceDataList := range resourceDataMap {
		if !convertedAddresses[address] {
			convertedAssets, err := converters.ConvertResource(resourceDataList, cfg, ancestryManager, o.ErrorLogger)
			if err != nil {
				return nil, fmt.Errorf("tfplan2cai converting: %w", err)
			}
			if len(convertedAssets) > 0 {
				parts := strings.SplitN(convertedAssets[0].Name, "/", 4)
				if len(parts) == 4 {
					for _, rd := range resourceDataList {
						rd.SetId(parts[3])
					}
				}
			}
			convertedAssetsByAddress[address] = convertedAssets
			assets = append(assets, convertedAssets...)
		}
	}

	return assets, nil
}
