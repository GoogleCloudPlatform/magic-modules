package google

import (
	"fmt"
	"math/rand"
	"regexp"
)

type ConvertFunc func(d TerraformResourceData, config *Config) ([]Asset, error)
type GetApiObjectFunc func(d TerraformResourceData, config *Config) (map[string]interface{}, error)

// FetchFullResourceFunc allows initial data for a resource to be fetched from the API and merged
// with the planned changes. This is useful for resources that are only partially managed
// by Terraform, like IAM policies managed with member/binding resources.
type FetchFullResourceFunc func(d TerraformResourceData, config *Config) (Asset, error)

// MergeFunc combines multiple terraform resources into a single CAI asset.
// The incoming asset will either be an asset that was created/updated or deleted.
type MergeFunc func(existing, incoming Asset) Asset

type ResourceConverter struct {
	AssetType         string
	Convert           ConvertFunc
	FetchFullResource FetchFullResourceFunc
	MergeCreateUpdate MergeFunc
	MergeDelete       MergeFunc
}

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string `json:"name"`
	// The type name in `google.<api>.<resourcename>` format.
	Type      string         `json:"asset_type"`
	Resource  *AssetResource `json:"resource,omitempty"`
	IAMPolicy *IAMPolicy     `json:"iam_policy,omitempty"`
	OrgPolicy []*OrgPolicy   `json:"org_policy,omitempty"`
}

// AssetResource is the Asset's Resource field.
type AssetResource struct {
	// Api version
	Version string `json:"version"`
	// URI including scheme for the discovery doc - assembled from
	// product name and version.
	DiscoveryDocumentURI string `json:"discovery_document_uri"`
	// Resource name.
	DiscoveryName string `json:"discovery_name"`
	// Actual resource state as per Terraform.  Note that this does
	// not necessarily correspond perfectly with the CAI representation
	// as there are occasional deviations between CAI and API responses.
	// This returns the API response values instead.
	Data map[string]interface{} `json:"data,omitempty"`
}

type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
}

type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

type OrgPolicy struct {
	Constraint     string          `json:"constraint,omitempty"`
	ListPolicy     *ListPolicy     `json:"listPolicy"`
	BooleanPolicy  *BooleanPolicy  `json:"booleanPolicy"`
	RestoreDefault *RestoreDefault `json:"restoreDefault"`
	UpdateTime     *Timestamp      `json:"update_time,omitempty"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nanos   int64 `json:"nanos,omitempty"`
}

type ListPolicyAllValues int32

type ListPolicy struct {
	AllowedValues     []string            `json:"allowed_values,omitempty"`
	DeniedValues      []string            `json:"denied_values,omitempty"`
	AllValues         ListPolicyAllValues `json:"all_values,omitempty"`
	SuggestedValue    string              `json:"suggested_value,omitempty"`
	InheritFromParent bool                `json:"inherit_from_parent,omitempty"`
}

type BooleanPolicy struct {
	Enforced bool `json:"enforced,omitempty"`
}

type RestoreDefault struct {
}

// assetName templates an asset.name by looking up and replacing all instances
// of {{field}}. In the case where a field would resolve to an empty string, a
// generated unique string will be used: "placeholder-" + randomString().
// This is done to preserve uniqueness of asset.name for a given asset.asset_type.
func assetName(d TerraformResourceData, config *Config, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{{([[:word:]]+)}}")

	// workaround for empty project
	placeholderSet := false
	if config.Project == "" {
		config.Project = fmt.Sprintf("placeholder-%s", randString(8))
		placeholderSet = true
	}

	f, err := buildReplacementFunc(re, d, config, linkTmpl, false)
	if err != nil {
		return "", err
	}
	if placeholderSet {
		config.Project = ""
	}

	fWithPlaceholder := func(key string) string {
		val := f(key)
		if val == "" {
			val = fmt.Sprintf("placeholder-%s", randString(8))
		}
		return val
	}

	return re.ReplaceAllStringFunc(linkTmpl, fWithPlaceholder), nil
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
