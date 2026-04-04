package transport

import (
	"context"
	"os"
	"strings"

	"cloud.google.com/go/bigtable"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type BigtableClientFactory struct {
	BasePath            string
	AdminBasePath       string
	UniverseDomain      string
	gRPCLoggingOptions  []option.ClientOption
	UserAgent           string
	TokenSource         oauth2.TokenSource
	BillingProject      string
	UserProjectOverride bool
}

func (s BigtableClientFactory) getClientOptions(isDataClient bool) []option.ClientOption {
	var opts []option.ClientOption
	if requestReason := os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"); requestReason != "" {
		opts = append(opts, option.WithRequestReason(requestReason))
	}

	if s.UserProjectOverride && s.BillingProject != "" {
		opts = append(opts, option.WithQuotaProject(s.BillingProject))
	}

	if s.UniverseDomain != "" {
		opts = append(opts, option.WithUniverseDomain(s.UniverseDomain))
	}

	if isDataClient && s.BasePath != "" {
		endpoint := strings.TrimPrefix(s.BasePath, "https://")
		endpoint = strings.TrimSuffix(endpoint, "/")

		if s.BasePath == s.AdminBasePath && strings.HasPrefix(endpoint, "bigtableadmin.") {
			endpoint = strings.Replace(endpoint, "bigtableadmin.", "bigtable.", 1)
		}

		opts = append(opts, option.WithEndpoint(endpoint))
	} else if !isDataClient && s.AdminBasePath != "" {
		endpoint := strings.TrimPrefix(s.AdminBasePath, "https://")
		endpoint = strings.TrimSuffix(endpoint, "/")
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	opts = append(opts, option.WithTokenSource(s.TokenSource), option.WithUserAgent(s.UserAgent))
	opts = append(opts, s.gRPCLoggingOptions...)

	return opts
}

func (s BigtableClientFactory) NewInstanceAdminClient(project string) (*bigtable.InstanceAdminClient, error) {
	opts := s.getClientOptions(false)
	return bigtable.NewInstanceAdminClient(context.Background(), project, opts...)
}

func (s BigtableClientFactory) NewAdminClient(project, instance string) (*bigtable.AdminClient, error) {
	opts := s.getClientOptions(false)
	return bigtable.NewAdminClient(context.Background(), project, instance, opts...)
}

func (s BigtableClientFactory) NewClient(project, instance string) (*bigtable.Client, error) {
	opts := s.getClientOptions(true)
	return bigtable.NewClient(context.Background(), project, instance, opts...)
}
