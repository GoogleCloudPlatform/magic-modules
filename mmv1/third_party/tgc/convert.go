// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	errorssyslib "errors"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/ancestrymanager"
	resources "github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/tfdata"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/tfplan"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var ErrDuplicateAsset = errors.New("duplicate asset")

// Asset contains the resource data and metadata in the same format as
// Google CAI (Cloud Asset Inventory).
type Asset struct {
	Name          string                    `json:"name"`
	Type          string                    `json:"asset_type"`
	Resource      *caiasset.AssetResource   `json:"resource,omitempty"`
	IAMPolicy     *caiasset.IAMPolicy       `json:"iam_policy,omitempty"`
	OrgPolicy     []*caiasset.OrgPolicy     `json:"org_policy,omitempty"`
	V2OrgPolicies []*caiasset.V2OrgPolicies `json:"v2_org_policies,omitempty"`

	// Store the converter's version of the asset to allow for merges which
	// operate on this type. When matching json tags land in the conversions
	// library, this could be nested to avoid the duplication of fields.
	converterAsset resources.Asset
	Ancestors      []string `json:"ancestors"`
}

// NewConverter is a factory function for Converter.
func NewConverter(cfg *transport_tpg.Config, ancestryManager ancestrymanager.AncestryManager, offline bool, convertUnchanged bool, errorLogger *zap.Logger) *Converter {
	return &Converter{
		schema:           provider.Provider(),
		converters:       resources.ResourceConverters(),
		offline:          offline,
		cfg:              cfg,
		ancestryManager:  ancestryManager,
		assets:           make(map[string]Asset),
		convertUnchanged: convertUnchanged,
		errorLogger:      errorLogger,
	}
}

// Converter knows how to convert terraform resources to their
// Google CAI (Cloud Asset Inventory) format (the Asset type).
type Converter struct {
	schema *schema.Provider

	// Map terraform resource kinds (i.e. "google_compute_instance")
	// to a ResourceConverter that can convert them to CAI assets.
	converters map[string][]resources.ResourceConverter

	offline bool
	cfg     *transport_tpg.Config

	// ancestryManager provides a manager to find the ancestry information for a project.
	ancestryManager ancestrymanager.AncestryManager

	// Map of converted assets (key = asset.Type + asset.Name)
	assets map[string]Asset

	// When set, Converter will convert ResourceChanges with no-op "actions".
	convertUnchanged bool

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

// AddResourceChange processes the resource changes in two stages:
// 1. Process deletions (fetching canonical resources from GCP as necessary)
// 2. Process creates, updates, and no-ops (fetching canonical resources from GCP as necessary)
// This will give us a deterministic end result even in cases where for example
// an IAM Binding and Member conflict with each other, but one is replacing the
// other.
func (c *Converter) AddResourceChanges(changes []*tfjson.ResourceChange) error {
	var createOrUpdateOrNoops []*tfjson.ResourceChange
	for _, rc := range changes {
		// Silently skip non-google resources
		if !strings.HasPrefix(rc.Type, "google_") {
			continue
		}

		// Warn about google-beta resources
		if rc.ProviderName == "registry.terraform.io/hashicorp/google-beta" {
			c.errorLogger.Debug(fmt.Sprintf("%s: resource uses the google-beta provider and may not be convertible", rc.Address))
		}

		// Skip resources not found in the google GA provider's schema
		if _, ok := c.schema.ResourcesMap[rc.Type]; !ok {
			c.errorLogger.Debug(fmt.Sprintf("%s: resource type not found in google GA provider: %s.", rc.Address, rc.Type))
			continue
		}

		// Skip unsupported resources
		if _, ok := c.converters[rc.Type]; !ok {
			c.errorLogger.Debug(fmt.Sprintf("%s: resource type cannot be converted for CAI-based policies: %s. For details, see https://cloud.google.com/docs/terraform/policy-validation/create-cai-constraints#supported_resources", rc.Address, rc.Type))
			continue
		}

		if tfplan.IsCreate(rc) || tfplan.IsUpdate(rc) || tfplan.IsDeleteCreate(rc) || (c.convertUnchanged && tfplan.IsNoOp(rc)) {
			createOrUpdateOrNoops = append(createOrUpdateOrNoops, rc)
		} else if tfplan.IsDelete(rc) {
			if err := c.addDelete(rc); err != nil {
				return fmt.Errorf("%s: converting deleted TF resource to CAI: %w", rc.Address, err)
			}
		}
	}

	for _, rc := range createOrUpdateOrNoops {
		if err := c.addCreateOrUpdateOrNoop(rc); err != nil {
			if errorssyslib.Is(err, ErrDuplicateAsset) {
				c.errorLogger.Warn(fmt.Sprintf("%s: converting TF resource to CAI: %v", rc.Address, err))
			} else {
				return fmt.Errorf("%s: converting TF resource to CAI: %w", rc.Address, err)
			}
		}
	}

	return nil
}

// For deletions, we only need to handle ResourceConverters that support
// both fetch and mergeDelete. Supporting just one doesn't
// make sense, and supporting neither means that the deletion
// can just happen without needing to be merged.
func (c *Converter) addDelete(rc *tfjson.ResourceChange) error {
	resource := c.schema.ResourcesMap[rc.Type]
	rd := tfdata.NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.Before.(map[string]interface{}),
	)
	for _, converter := range c.converters[rd.Kind()] {
		if converter.FetchFullResource == nil || converter.MergeDelete == nil {
			continue
		}
		convertedItems, err := convertWrapper(converter, rd, c.cfg)
		if err != nil {
			if errors.Cause(err) == resources.ErrNoConversion {
				continue
			}
			return err
		}

		for _, converted := range convertedItems {
			key := converted.Type + converted.Name
			var existingConverterAsset *resources.Asset
			if existing, exists := c.assets[key]; exists {
				existingConverterAsset = &existing.converterAsset
			} else if !c.offline {
				asset, err := converter.FetchFullResource(rd, c.cfg)
				if errors.Cause(err) == resources.ErrEmptyIdentityField {
					c.errorLogger.Debug(fmt.Sprintf("%s: Unable to fetch and merge remote %s asset due to unset or (known after apply) identity fields on the TF resource.", rc.Address, converted.Type))
					existingConverterAsset = nil
				} else if errors.Cause(err) == resources.ErrResourceInaccessible {
					c.errorLogger.Warn(fmt.Sprintf("%s: Fetching %s for merge failed due to not existing or insufficient permission.", rc.Address, key))
					existingConverterAsset = nil
				} else if err != nil {
					return fmt.Errorf("fetching remote asset %s: %w", key, err)
				} else {
					existingConverterAsset = &asset
				}
			}
			if existingConverterAsset != nil {
				converted = converter.MergeDelete(*existingConverterAsset, converted)
				augmented, err := c.augmentAsset(rd, c.cfg, converted)
				if err != nil {
					return err
				}
				c.assets[key] = augmented
			}
		}
	}

	return nil
}

// For create/update/no-op, we need to handle both the case of no merging,
// and the case of merging. If merging, we expect both fetch and mergeCreateUpdate
// to be present.
func (c *Converter) addCreateOrUpdateOrNoop(rc *tfjson.ResourceChange) error {
	resource := c.schema.ResourcesMap[rc.Type]
	rd := tfdata.NewFakeResourceData(
		rc.Type,
		resource.Schema,
		rc.Change.After.(map[string]interface{}),
	)

	for _, converter := range c.converters[rd.Kind()] {
		convertedAssets, err := convertWrapper(converter, rd, c.cfg)
		if err != nil {
			if errors.Cause(err) == resources.ErrNoConversion {
				continue
			}
			return err
		}

		for _, converted := range convertedAssets {
			key := converted.Type + converted.Name

			var existingConverterAsset *resources.Asset
			if existing, exists := c.assets[key]; exists {
				existingConverterAsset = &existing.converterAsset
			} else if converter.FetchFullResource != nil && !c.offline {
				asset, err := converter.FetchFullResource(rd, c.cfg)
				if errors.Cause(err) == resources.ErrEmptyIdentityField {
					c.errorLogger.Debug(fmt.Sprintf("%s: Unable to fetch and merge remote %s asset due to unset or (known after apply) identity fields on the TF resource.", rc.Address, converted.Type))
					existingConverterAsset = nil
				} else if errors.Cause(err) == resources.ErrResourceInaccessible {
					c.errorLogger.Warn(fmt.Sprintf("%s: Fetching %s for merge failed due to not existing or insufficient permission.", rc.Address, key))
					existingConverterAsset = nil
				} else if err != nil {
					return fmt.Errorf("fetching remote asset %s: %w", key, err)
				} else {
					existingConverterAsset = &asset
				}
			}

			if existingConverterAsset != nil {
				if converter.MergeCreateUpdate == nil {
					// If a merge function does not exist ignore the asset and return
					// a checkable error.
					return fmt.Errorf("%w: type %s: name %s", ErrDuplicateAsset, converted.Type, converted.Name)
				}
				converted = converter.MergeCreateUpdate(*existingConverterAsset, converted)
			}

			augmented, err := c.augmentAsset(rd, c.cfg, converted)
			if err != nil {
				return err
			}
			c.assets[key] = augmented
		}
	}

	return nil
}

type byName []caiasset.Asset

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Assets lists all converted assets previously added by calls to AddResource.
func (c *Converter) Assets() []caiasset.Asset {
	list := make([]caiasset.Asset, 0, len(c.assets))
	for _, a := range c.assets {
		list = append(list, caiasset.Asset{
			Name:          a.Name,
			Type:          a.Type,
			Resource:      a.Resource,
			IAMPolicy:     a.IAMPolicy,
			OrgPolicy:     a.OrgPolicy,
			V2OrgPolicies: a.V2OrgPolicies,
			Ancestors:     a.Ancestors,
		})
	}
	sort.Sort(byName(list))
	return list
}

// augmentAsset adds data to an asset that is not set by the conversion library.
func (c *Converter) augmentAsset(tfData tpgresource.TerraformResourceData, cfg *transport_tpg.Config, cai resources.Asset) (Asset, error) {
	ancestors, parent, err := c.ancestryManager.Ancestors(cfg, tfData, &cai)
	if err != nil {
		return Asset{}, fmt.Errorf("getting resource ancestry or parent failed: %w", err)
	}

	var resource *caiasset.AssetResource
	if cai.Resource != nil {
		resource = &caiasset.AssetResource{
			Version:              cai.Resource.Version,
			DiscoveryDocumentURI: cai.Resource.DiscoveryDocumentURI,
			DiscoveryName:        cai.Resource.DiscoveryName,
			Parent:               parent,
			Data:                 cai.Resource.Data,
		}
	}

	var policy *caiasset.IAMPolicy
	if cai.IAMPolicy != nil {
		policy = &caiasset.IAMPolicy{}
		for _, b := range cai.IAMPolicy.Bindings {
			policy.Bindings = append(policy.Bindings, caiasset.IAMBinding{
				Role:    b.Role,
				Members: b.Members,
			})
		}
	}

	var orgPolicy []*caiasset.OrgPolicy
	if cai.OrgPolicy != nil {
		for _, o := range cai.OrgPolicy {
			var listPolicy *caiasset.ListPolicy
			var booleanPolicy *caiasset.BooleanPolicy
			var restoreDefault *caiasset.RestoreDefault
			if o.ListPolicy != nil {
				listPolicy = &caiasset.ListPolicy{
					AllowedValues:     o.ListPolicy.AllowedValues,
					AllValues:         caiasset.ListPolicyAllValues(o.ListPolicy.AllValues),
					DeniedValues:      o.ListPolicy.DeniedValues,
					SuggestedValue:    o.ListPolicy.SuggestedValue,
					InheritFromParent: o.ListPolicy.InheritFromParent,
				}
			}
			if o.BooleanPolicy != nil {
				booleanPolicy = &caiasset.BooleanPolicy{
					Enforced: o.BooleanPolicy.Enforced,
				}
			}
			if o.RestoreDefault != nil {
				restoreDefault = &caiasset.RestoreDefault{}
			}
			//As time is not information in terraform resource data, time is fixed for testing purposes
			fixedTime := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
			orgPolicy = append(orgPolicy, &caiasset.OrgPolicy{
				Constraint:     o.Constraint,
				ListPolicy:     listPolicy,
				BooleanPolicy:  booleanPolicy,
				RestoreDefault: restoreDefault,
				UpdateTime: &caiasset.Timestamp{
					Seconds: int64(fixedTime.Unix()),
					Nanos:   int64(fixedTime.UnixNano()),
				},
			})
		}
	}

	var v2OrgPolicies []*caiasset.V2OrgPolicies
	if cai.V2OrgPolicies != nil {
		for _, o2 := range cai.V2OrgPolicies {
			var spec *caiasset.PolicySpec
			if o2.PolicySpec != nil {

				var rules []*caiasset.PolicyRule
				if o2.PolicySpec.PolicyRules != nil {
					for _, rule := range o2.PolicySpec.PolicyRules {
						var values *caiasset.StringValues
						if rule.Values != nil {
							values = &caiasset.StringValues{
								AllowedValues: rule.Values.AllowedValues,
								DeniedValues:  rule.Values.DeniedValues,
							}
						}

						var condition *caiasset.Expr
						if rule.Condition != nil {
							condition = &caiasset.Expr{
								Expression:  rule.Condition.Expression,
								Title:       rule.Condition.Title,
								Description: rule.Condition.Description,
								Location:    rule.Condition.Location,
							}
						}
						rules = append(rules, &caiasset.PolicyRule{
							Values:    values,
							AllowAll:  rule.AllowAll,
							DenyAll:   rule.DenyAll,
							Enforce:   rule.Enforce,
							Condition: condition,
						})
					}
				}

				fixedTime := time.Date(2021, time.April, 14, 15, 16, 17, 0, time.UTC)
				spec = &caiasset.PolicySpec{
					Etag: o2.PolicySpec.Etag,
					UpdateTime: &caiasset.Timestamp{
						Seconds: int64(fixedTime.Unix()),
						Nanos:   int64(fixedTime.UnixNano()),
					},
					PolicyRules:       rules,
					InheritFromParent: o2.PolicySpec.InheritFromParent,
					Reset:             o2.PolicySpec.Reset,
				}

			}

			v2OrgPolicies = append(v2OrgPolicies, &caiasset.V2OrgPolicies{
				Name:       o2.Name,
				PolicySpec: spec,
			})
		}
	}

	return Asset{
		Name:           cai.Name,
		Type:           cai.Type,
		Resource:       resource,
		IAMPolicy:      policy,
		OrgPolicy:      orgPolicy,
		V2OrgPolicies:  v2OrgPolicies,
		converterAsset: cai,
		Ancestors:      ancestors,
	}, nil
}

func convertWrapper(conv resources.ResourceConverter, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (assets []resources.Asset, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = fmt.Errorf("unknown panic error: %v", v)
			}
			err = fmt.Errorf("%v\n Stack trace: %s", err, string(debug.Stack()))
		}
	}()
	return conv.Convert(d, config)
}
