package tfplan2cai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/ancestrymanager"
	convert "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google"
	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/tfplan2cai/tfplan"
	"go.uber.org/zap"
)

// Options struct to avoid updating function signatures all along the pipe.
type Options struct {
	ConvertUnchanged bool
	ErrorLogger      *zap.Logger
	Offline          bool
	DefaultProject   string
	DefaultRegion    string
	DefaultZone      string
	// UserAgent for all requests (if online)
	UserAgent string
	// HTTPClient for all requests (if online)
	HTTPClient *http.Client
	// Map hierarchy resource (like projects/<number> or folders/<number>)
	// to an ancestry path (like organizations/123/folders/456/projects/789)
	AncestryCache map[string]string
}

// Convert converts terraform json plan to CAI Assets.
func Convert(ctx context.Context, jsonPlan []byte, o *Options) ([]caiasset.Asset, error) {
	// Creates ancestry manager and converter internally; they are
	// implementation details private to this package.

	errorLogger := o.ErrorLogger
	if errorLogger == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("building logger: %w", err)
		}
		errorLogger = logger
	}

	// contents of newConverter
	// Set up config and ancestry manager using the same http client and user agent
	cfg, err := resources.NewConfig(ctx, o.DefaultProject, o.DefaultZone, o.DefaultRegion, o.Offline, o.UserAgent, o.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	ancestryManager, err := ancestrymanager.New(cfg, o.Offline, o.AncestryCache, errorLogger)
	if err != nil {
		return nil, fmt.Errorf("building ancestry manager: %w", err)
	}

	converter := convert.NewConverter(cfg, ancestryManager, o.Offline, o.ConvertUnchanged, errorLogger)

	// ReadResourceChanges
	changes, err := tfplan.ReadResourceChanges(jsonPlan)
	if err != nil {
		return nil, err
	}

	err = converter.AddResourceChanges(changes)
	if err != nil {
		return nil, err
	}

	return converter.Assets(), nil
}
