package caiasset

import (
	"fmt"
	"strings"
	"time"
)

// Asset is the CAI representation of a resource.
type Asset struct {
	// The name, in a peculiar format: `\\<api>.googleapis.com/<self_link>`
	Name string `json:"name"`
	// The type name in `google.<api>.<resourcename>` format.
	Type          string           `json:"asset_type"`
	Resource      *AssetResource   `json:"resource,omitempty"`
	IAMPolicy     *IAMPolicy       `json:"iam_policy,omitempty"`
	OrgPolicy     []*OrgPolicy     `json:"org_policy,omitempty"`
	V2OrgPolicies []*V2OrgPolicies `json:"v2_org_policies,omitempty"`
	Ancestors     []string         `json:"ancestors"`
}

// IAMPolicy is the representation of a Cloud IAM policy set on a cloud resource.
type IAMPolicy struct {
	Bindings []IAMBinding `json:"bindings"`
}

// IAMBinding binds a role to a set of members.
type IAMBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

// AssetResource is nested within the Asset type.
type AssetResource struct {
	Version              string                 `json:"version"`
	DiscoveryDocumentURI string                 `json:"discovery_document_uri"`
	DiscoveryName        string                 `json:"discovery_name"`
	Parent               string                 `json:"parent"`
	Data                 map[string]interface{} `json:"data"`
	Location             string                 `json:"location,omitempty"`
}

// OrgPolicy is for managing organization policies.
type OrgPolicy struct {
	Constraint     string          `json:"constraint,omitempty"`
	ListPolicy     *ListPolicy     `json:"list_policy,omitempty"`
	BooleanPolicy  *BooleanPolicy  `json:"boolean_policy,omitempty"`
	RestoreDefault *RestoreDefault `json:"restore_default,omitempty"`
	UpdateTime     *Timestamp      `json:"update_time,omitempty"`
}

// V2OrgPolicies is the represtation of V2OrgPolicies
type V2OrgPolicies struct {
	Name       string      `json:"name"`
	PolicySpec *PolicySpec `json:"spec,omitempty"`
}

// Spec is the representation of Spec for Custom Org Policy
type PolicySpec struct {
	Etag              string        `json:"etag,omitempty"`
	UpdateTime        *Timestamp    `json:"update_time,omitempty"`
	PolicyRules       []*PolicyRule `json:"rules,omitempty"`
	InheritFromParent bool          `json:"inherit_from_parent,omitempty"`
	Reset             bool          `json:"reset,omitempty"`
}

type PolicyRule struct {
	Values    *StringValues `json:"values,omitempty"`
	AllowAll  bool          `json:"allow_all,omitempty"`
	DenyAll   bool          `json:"deny_all,omitempty"`
	Enforce   bool          `json:"enforce,omitempty"`
	Condition *Expr         `json:"condition,omitempty"`
}

type StringValues struct {
	AllowedValues []string `json:"allowed_values,omitempty"`
	DeniedValues  []string `json:"denied_values,omitempty"`
}

type Expr struct {
	Expression  string `json:"expression,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Timestamp struct {
	Seconds int64 `json:"seconds,omitempty"`
	Nanos   int64 `json:"nanos,omitempty"`
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Unix(0, t.Nanos).UTC().Format(time.RFC3339Nano) + `"`), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	p, err := time.Parse(time.RFC3339Nano, strings.Trim(string(b), `"`))
	if err != nil {
		return fmt.Errorf("bad Timestamp: %v", err)
	}
	t.Seconds = p.Unix()
	t.Nanos = p.UnixNano()
	return nil
}

// ListPolicyAllValues is used to set `Policies` that apply to all possible
// configuration values rather than specific values in `allowed_values` or
// `denied_values`.
type ListPolicyAllValues int32

// ListPolicy can define specific values and subtrees of Cloud Resource
// Manager resource hierarchy (`Organizations`, `Folders`, `Projects`) that
// are allowed or denied by setting the `allowed_values` and `denied_values`
// fields.
type ListPolicy struct {
	AllowedValues     []string            `json:"allowed_values,omitempty"`
	DeniedValues      []string            `json:"denied_values,omitempty"`
	AllValues         ListPolicyAllValues `json:"all_values,omitempty"`
	SuggestedValue    string              `json:"suggested_value,omitempty"`
	InheritFromParent bool                `json:"inherit_from_parent,omitempty"`
}

// BooleanPolicy If `true`, then the `Policy` is enforced. If `false`,
// then any configuration is acceptable.
type BooleanPolicy struct {
	Enforced bool `json:"enforced,omitempty"`
}

// RestoreDefault determines if the default values of the `Constraints` are active for the
// resources.
type RestoreDefault struct {
}
