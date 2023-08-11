package cai2hcl

import (
	"context"
	"net/http"

	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

func NewConfig() (*transport_tpg.Config, error) {
	ctx := context.Background()
	defaultProject := ""
	defaultZone := ""
	defaultRegion := ""
	offline := true
	userAgent := ""
	var httpClient *http.Client = nil

	return resources.NewConfig(ctx, defaultProject, defaultZone, defaultRegion, offline, userAgent, httpClient)
}
