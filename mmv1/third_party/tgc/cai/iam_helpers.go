package cai

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// ExpandIamPolicyBindings is used in google_<type>_iam_policy resources.
func ExpandIamPolicyBindings(d tpgresource.TerraformResourceData) ([]IAMBinding, error) {
	ps := d.Get("policy_data").(string)
	var bindings []IAMBinding
	// policy_data is (known after apply) in terraform plan, hence an empty string
	if ps == "" {
		return bindings, nil
	}
	// The policy string is just a marshaled cloudresourcemanager.Policy.
	policy := &cloudresourcemanager.Policy{}
	if err := json.Unmarshal([]byte(ps), policy); err != nil {
		return nil, fmt.Errorf("Could not unmarshal %s: %v", ps, err)
	}

	for _, b := range policy.Bindings {
		bindings = append(bindings, IAMBinding{
			Role:    b.Role,
			Members: b.Members,
		})
	}

	return bindings, nil
}

// ExpandIamRoleBindings is used in google_<type>_iam_binding resources.
func ExpandIamRoleBindings(d tpgresource.TerraformResourceData) ([]IAMBinding, error) {
	var members []string
	for _, m := range d.Get("members").(*schema.Set).List() {
		members = append(members, m.(string))
	}
	return []IAMBinding{
		{
			Role:    d.Get("role").(string),
			Members: members,
		},
	}, nil
}

// ExpandIamMemberBindings is used in google_<type>_iam_member resources.
func ExpandIamMemberBindings(d tpgresource.TerraformResourceData) ([]IAMBinding, error) {
	return []IAMBinding{
		{
			Role:    d.Get("role").(string),
			Members: []string{d.Get("member").(string)},
		},
	}, nil
}

// MergeIamAssets merges an existing asset with the IAM bindings of an incoming
// Asset.
func MergeIamAssets(
	existing, incoming Asset,
	MergeBindings func(existing, incoming []IAMBinding) []IAMBinding,
) Asset {
	if existing.IAMPolicy != nil {
		existing.IAMPolicy.Bindings = MergeBindings(existing.IAMPolicy.Bindings, incoming.IAMPolicy.Bindings)
	} else {
		existing.IAMPolicy = incoming.IAMPolicy
	}
	return existing
}

// incoming is the last known state of an asset prior to deletion
func MergeDeleteIamAssets(
	existing, incoming Asset,
	MergeBindings func(existing, incoming []IAMBinding) []IAMBinding,
) Asset {
	if existing.IAMPolicy != nil {
		existing.IAMPolicy.Bindings = MergeBindings(existing.IAMPolicy.Bindings, incoming.IAMPolicy.Bindings)
	}
	return existing
}

// MergeAdditiveBindings adds members to bindings with the same roles and adds new
// bindings for roles that dont exist.
func MergeAdditiveBindings(existing, incoming []IAMBinding) []IAMBinding {
	existingIdxs := make(map[string]int)
	for i, binding := range existing {
		existingIdxs[binding.Role] = i
	}

	for _, binding := range incoming {
		if ei, ok := existingIdxs[binding.Role]; ok {
			memberExists := make(map[string]bool)
			for _, m := range existing[ei].Members {
				memberExists[m] = true
			}
			for _, m := range binding.Members {
				// Only add members that don't exist.
				if !memberExists[m] {
					existing[ei].Members = append(existing[ei].Members, m)
				}
			}
		} else {
			existing = append(existing, binding)
		}
	}

	// Sort members
	for i := range existing {
		sort.Strings(existing[i].Members)
	}

	return existing
}

// MergeDeleteAdditiveBindings eliminates listed members from roles in the
// existing list. incoming is the last known state of the bindings being deleted.
func MergeDeleteAdditiveBindings(existing, incoming []IAMBinding) []IAMBinding {
	toDelete := make(map[string]struct{})
	for _, binding := range incoming {
		for _, m := range binding.Members {
			key := binding.Role + "-" + m
			toDelete[key] = struct{}{}
		}
	}

	var newExisting []IAMBinding
	for _, binding := range existing {
		var newMembers []string
		for _, m := range binding.Members {
			key := binding.Role + "-" + m
			_, delete := toDelete[key]
			if !delete {
				newMembers = append(newMembers, m)
			}
		}
		if newMembers != nil {
			newExisting = append(newExisting, IAMBinding{
				Role:    binding.Role,
				Members: newMembers,
			})
		}
	}

	return newExisting
}

// MergeAuthoritativeBindings clobbers members to bindings with the same roles
// and adds new bindings for roles that dont exist.
func MergeAuthoritativeBindings(existing, incoming []IAMBinding) []IAMBinding {
	existingIdxs := make(map[string]int)
	for i, binding := range existing {
		existingIdxs[binding.Role] = i
	}

	for _, binding := range incoming {
		if ei, ok := existingIdxs[binding.Role]; ok {
			existing[ei].Members = binding.Members
		} else {
			existing = append(existing, binding)
		}
	}

	// Sort members
	for i := range existing {
		sort.Strings(existing[i].Members)
	}

	return existing
}

// MergeDeleteAuthoritativeBindings eliminates any bindings with matching roles
// in the existing list. incoming is the last known state of the bindings being
// deleted.
func MergeDeleteAuthoritativeBindings(existing, incoming []IAMBinding) []IAMBinding {
	toDelete := make(map[string]struct{})
	for _, binding := range incoming {
		key := binding.Role
		toDelete[key] = struct{}{}
	}

	var newExisting []IAMBinding
	for _, binding := range existing {
		key := binding.Role
		_, delete := toDelete[key]
		if !delete {
			newExisting = append(newExisting, binding)
		}
	}

	return newExisting
}

func FetchIamPolicy(
	newUpdaterFunc tpgiamresource.NewResourceIamUpdaterFunc,
	d tpgresource.TerraformResourceData,
	config *transport_tpg.Config,
	assetNameTmpl string,
	assetType string,
) (Asset, error) {
	updater, err := newUpdaterFunc(d, config)
	if err != nil {
		return Asset{}, err
	}

	iamPolicy, err := updater.GetResourceIamPolicy()
	if transport_tpg.IsGoogleApiErrorWithCode(err, 403) || transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		return Asset{}, ErrResourceInaccessible
	}

	if err != nil {
		return Asset{}, err
	}

	var bindings []IAMBinding
	for _, b := range iamPolicy.Bindings {
		bindings = append(
			bindings,
			IAMBinding{
				Role:    b.Role,
				Members: b.Members,
			},
		)
	}

	name, err := AssetName(d, config, assetNameTmpl)

	return Asset{
		Name: name,
		Type: assetType,
		IAMPolicy: &IAMPolicy{
			Bindings: bindings,
		},
	}, nil
}
