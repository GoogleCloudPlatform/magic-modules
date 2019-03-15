package google

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string
	// The type name in `google.<api>.<resourcename>` format.
	Type      string
	Resource  *AssetResource
	IAMPolicy *IAMPolicy
}

// AssetResource is the Asset's Resource field.
type AssetResource struct {
	// Api version
	Version string
	// URI including scheme for the discovery doc - assembled from
	// product name and version.
	DiscoveryDocumentURI string
	// Resource name.
	DiscoveryName string
	// Actual resource state as per Terraform.  Note that this does
	// not necessarily correspond perfectly with the CAI representation
	// as there are occasional deviations between CAI and API responses.
	// This returns the API response values instead.
	Data map[string]interface{}
}

type IAMPolicy struct {
	Bindings []IAMBinding
}

type IAMBinding struct {
	Role    string
	Members []string
}
