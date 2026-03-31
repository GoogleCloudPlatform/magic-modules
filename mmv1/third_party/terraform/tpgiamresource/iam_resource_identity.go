package tpgiamresource

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// IamIdentityParam describes a single parameter of the IAM resource URI.
type IamIdentityParam struct {
	Key         string // raw key for values map and regex captures (e.g. "serviceId")
	IdentityKey string // snake_case key for identity.GetOk (e.g. "service_id")
}

// IamResourceIdentityConfig holds all the per-resource data needed to parse an
// IAM import identity into a canonical resource id.
type IamResourceIdentityConfig struct {
	Params    []IamIdentityParam
	UriFormat string // fmt.Sprintf format producing the canonical resource id
}

var locationDefaultFuncs = map[string]func(*schema.ResourceData, *transport_tpg.Config) (string, error){
	"project":  tpgresource.GetProject,
	"zone":     tpgresource.GetZone,
	"region":   tpgresource.GetRegion,
	"location": tpgresource.GetLocation,
}

// ParseIamResourceIdentity resolves an IAM import identity into a canonical
// resource id string (the same shape as the IAM updater's GetResourceId()).
func ParseIamResourceIdentity(
	d *schema.ResourceData,
	identity *schema.IdentityData,
	config *transport_tpg.Config,
	rc IamResourceIdentityConfig,
) (string, error) {
	resolved := make(map[string]string, len(rc.Params))
	var nonLocParams []IamIdentityParam
	nonLocCount := 0
	for _, p := range rc.Params {
		if fn, isLoc := locationDefaultFuncs[p.Key]; isLoc {
			var val string
			if rv, ok := identity.GetOk(p.IdentityKey); ok {
				if s, ok := rv.(string); ok {
					val = s
				}
			}
			if val == "" {
				res, err := fn(d, config)
				if err != nil {
					return "", err
				}
				val = res
			}
			if val == "" {
				return "", fmt.Errorf("could not determine %q for IAM import identity; set it on the resource or configure the provider", p.IdentityKey)
			}
			resolved[p.Key] = val
		} else {
			nonLocParams = append(nonLocParams, p)
			if v, ok := identity.GetOk(p.IdentityKey); ok {
				if s, ok := v.(string); ok && s != "" {
					nonLocCount++
				}
			}
		}
	}

	if nonLocCount == 0 {
		return formatIamResourceUri(rc, resolved)
	}

	return parseIdentityAttributes(identity, rc, resolved)
}

// parseMultiAttrIdentity handles the case where 2+ non-location params are set:
// each non-location param is extracted individually from the identity. The last
// param in the list gets GetResourceNameFromSelfLink applied.
func parseIdentityAttributes(
	identity *schema.IdentityData,
	rc IamResourceIdentityConfig,
	resolved map[string]string,
) (string, error) {
	for i, p := range rc.Params {
		if _, isLoc := locationDefaultFuncs[p.Key]; isLoc {
			continue
		}
		val, ok := identity.GetOk(p.IdentityKey)
		if !ok {
			return "", fmt.Errorf("import identity is missing attribute %q", p.IdentityKey)
		}
		s, strOk := val.(string)
		if !strOk || s == "" {
			return "", fmt.Errorf("import identity attribute %q must be a non-empty string", p.IdentityKey)
		}
		if i == len(rc.Params)-1 {
			s = tpgresource.GetResourceNameFromSelfLink(s)
		}
		resolved[p.Key] = s
	}
	return formatIamResourceUri(rc, resolved)
}

func formatIamResourceUri(rc IamResourceIdentityConfig, values map[string]string) (string, error) {
	args := make([]any, len(rc.Params))
	for i, p := range rc.Params {
		args[i] = values[p.Key]
	}
	return fmt.Sprintf(rc.UriFormat, args...), nil
}
