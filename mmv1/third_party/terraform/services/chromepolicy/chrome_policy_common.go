package chromepolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"google.golang.org/api/chromepolicy/v1"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const chromePolicyRetryTimeout = 5 * time.Minute

// Target kind
type chromePolicyTargetKind string

const (
	targetOrgUnit chromePolicyTargetKind = "orgunits"
	targetGroup   chromePolicyTargetKind = "groups"
)

func chromePolicyTargetResource(kind chromePolicyTargetKind, id string) string {
	return string(kind) + "/" + id
}

// PolicyTargetKey builder
func buildPolicyTargetKey(targetResource string, pol map[string]interface{}) *chromepolicy.GoogleChromePolicyVersionsV1PolicyTargetKey {
	key := &chromepolicy.GoogleChromePolicyVersionsV1PolicyTargetKey{
		TargetResource: targetResource,
	}
	if atk, ok := pol["additional_target_keys"]; ok && atk != nil {
		raw := atk.(map[string]interface{})
		if len(raw) > 0 {
			key.AdditionalTargetKeys = make(map[string]string, len(raw))
			for k, v := range raw {
				key.AdditionalTargetKeys[k] = v.(string)
			}
		}
	}
	return key
}

type schemaCache struct {
	service  *chromepolicy.CustomersPolicySchemasService
	customer string
	cache    map[string]*chromepolicy.GoogleChromePolicyVersionsV1PolicySchema
}

func newSchemaCache(service *chromepolicy.CustomersPolicySchemasService, customer string) *schemaCache {
	return &schemaCache{
		service:  service,
		customer: customer,
		cache:    make(map[string]*chromepolicy.GoogleChromePolicyVersionsV1PolicySchema),
	}
}

func (sc *schemaCache) get(ctx context.Context, schemaName string) (*chromepolicy.GoogleChromePolicyVersionsV1PolicySchema, error) {
	if cached, ok := sc.cache[schemaName]; ok {
		return cached, nil
	}

	var schemaDef *chromepolicy.GoogleChromePolicyVersionsV1PolicySchema
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			var retryErr error
			schemaDef, retryErr = sc.service.Get(
				fmt.Sprintf("customers/%s/policySchemas/%s", sc.customer, schemaName),
			).Do()
			return retryErr
		},
		Timeout: chromePolicyRetryTimeout,
	})
	if err != nil {
		return nil, err
	}

	sc.cache[schemaName] = schemaDef
	return schemaDef, nil
}

// Validation
func buildSchemaFieldMap(schemaDef *chromepolicy.GoogleChromePolicyVersionsV1PolicySchema) map[string]*chromepolicy.Proto2FieldDescriptorProto {
	fieldMap := make(map[string]*chromepolicy.Proto2FieldDescriptorProto)
	for _, mt := range schemaDef.Definition.MessageType {
		for i, f := range mt.Field {
			fieldMap[f.Name] = mt.Field[i]
		}
	}
	return fieldMap
}

// validatePolicies validates policy values against their schema definitions.
func validatePolicies(ctx context.Context, policies []interface{}, sc *schemaCache) diag.Diagnostics {
	for _, policy := range policies {
		pol := policy.(map[string]interface{})
		schemaName := pol["schema"].(string)
		schemaValues := pol["value"].(map[string]interface{})

		schemaDef, err := sc.get(ctx, schemaName)
		if err != nil {
			return diag.FromErr(err)
		}

		if schemaDef == nil || schemaDef.Definition == nil || schemaDef.Definition.MessageType == nil {
			return diag.Errorf("schema definition (%s) is empty", schemaName)
		}

		fieldMap := buildSchemaFieldMap(schemaDef)
		for fieldName, jsonVal := range schemaValues {
			field, ok := fieldMap[fieldName]
			if !ok {
				return diag.Errorf("field %q is not found in schema %s", fieldName, schemaName)
			}

			var val interface{}
			if err := json.Unmarshal([]byte(jsonVal.(string)), &val); err != nil {
				return diag.FromErr(err)
			}

			if field.Label == "LABEL_REPEATED" {
				arr, ok := val.([]interface{})
				if !ok {
					return diag.Errorf("value for %s.%s must be an array (got %T)", schemaName, fieldName, val)
				}
				for _, item := range arr {
					if !validatePolicyFieldValueType(field.Type, item) {
						return diag.Errorf("array element in %s.%s has incorrect type (expected %s)", schemaName, fieldName, field.Type)
					}
				}
			} else if !validatePolicyFieldValueType(field.Type, val) {
				return diag.Errorf("value for %s.%s has incorrect type (expected %s)", schemaName, fieldName, field.Type)
			}
		}

		perPolicyATK := identityFromPolicy(pol).AdditionalTargetKeys
		if schemaDef.AdditionalTargetKeyNames != nil && len(perPolicyATK) == 0 {
			return diag.Errorf("schema %s requires additional_target_keys", schemaName)
		}
		if schemaDef.AdditionalTargetKeyNames == nil && len(perPolicyATK) > 0 {
			return diag.Errorf("schema %s does not support additional_target_keys", schemaName)
		}

		if len(perPolicyATK) > 0 {
			allowed := make(map[string]bool)
			for _, tkn := range schemaDef.AdditionalTargetKeyNames {
				allowed[tkn.Key] = true
			}
			for k := range perPolicyATK {
				if !allowed[k] {
					return diag.Errorf("additional_target_key %q is not valid for schema %s", k, schemaName)
				}
			}
		}
	}

	return nil
}

func validatePolicyFieldValueType(fieldType string, fieldValue interface{}) bool {
	switch fieldType {
	case "TYPE_BOOL":
		return reflect.ValueOf(fieldValue).Kind() == reflect.Bool
	case "TYPE_FLOAT", "TYPE_DOUBLE":
		return reflect.ValueOf(fieldValue).Kind() == reflect.Float64
	case "TYPE_INT64", "TYPE_FIXED64", "TYPE_SFIXED64", "TYPE_SINT64", "TYPE_UINT64":
		if reflect.ValueOf(fieldValue).Kind() == reflect.Float64 &&
			fieldValue == float64(int(fieldValue.(float64))) {
			return true
		}
		return false
	case "TYPE_INT32", "TYPE_FIXED32", "TYPE_SFIXED32", "TYPE_SINT32", "TYPE_UINT32":
		if reflect.ValueOf(fieldValue).Kind() == reflect.Float64 &&
			fieldValue == float64(int32(fieldValue.(float64))) {
			return true
		}
		return false
	case "TYPE_MESSAGE":
		return reflect.ValueOf(fieldValue).Kind() == reflect.Map
	default: // TYPE_ENUM, TYPE_STRING, etc.
		return reflect.ValueOf(fieldValue).Kind() == reflect.String
	}
}

func batchModifyPolicies(
	ctx context.Context,
	policiesService *chromepolicy.CustomersPoliciesService,
	customer string,
	targetResource string,
	policies []interface{},
) diag.Diagnostics {
	var requests []*chromepolicy.GoogleChromePolicyVersionsV1ModifyOrgUnitPolicyRequest

	for _, p := range policies {
		pol := p.(map[string]interface{})
		schemaName := pol["schema"].(string)
		schemaValues := pol["value"].(map[string]interface{})

		// Build the JSON directly from the per-field JSON strings to avoid an unnecessary unmarshal/marshal round-trip.
		var updateKeys []string
		jsonParts := make([]string, 0, len(schemaValues))
		for k, v := range schemaValues {
			keyJSON, _ := json.Marshal(k)
			jsonParts = append(jsonParts, string(keyJSON)+":"+v.(string))
			updateKeys = append(updateKeys, k)
		}
		valueJSON := []byte("{" + strings.Join(jsonParts, ",") + "}")

		requests = append(requests, &chromepolicy.GoogleChromePolicyVersionsV1ModifyOrgUnitPolicyRequest{
			PolicyTargetKey: buildPolicyTargetKey(targetResource, pol),
			PolicyValue: &chromepolicy.GoogleChromePolicyVersionsV1PolicyValue{
				PolicySchema: schemaName,
				Value:        valueJSON,
			},
			UpdateMask: strings.Join(updateKeys, ","),
		})
	}

	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			_, retryErr := policiesService.Orgunits.BatchModify(
				fmt.Sprintf("customers/%s", customer),
				&chromepolicy.GoogleChromePolicyVersionsV1BatchModifyOrgUnitPoliciesRequest{Requests: requests},
			).Do()
			return retryErr
		},
		Timeout: chromePolicyRetryTimeout,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func inheritPolicies(
	ctx context.Context,
	policiesService *chromepolicy.CustomersPoliciesService,
	customer string,
	requests []*chromepolicy.GoogleChromePolicyVersionsV1InheritOrgUnitPolicyRequest,
) diag.Diagnostics {
	if len(requests) == 0 {
		return nil
	}
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			_, retryErr := policiesService.Orgunits.BatchInherit(
				fmt.Sprintf("customers/%s", customer),
				&chromepolicy.GoogleChromePolicyVersionsV1BatchInheritOrgUnitPoliciesRequest{Requests: requests},
			).Do()
			return retryErr
		},
		Timeout: chromePolicyRetryTimeout,
	})
	if err != nil {
		if isNonFatalDeleteError(err) {
			log.Printf("[DEBUG] Non-fatal error during policy inherit: %v", err)
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

type resolvedPolicyEntry struct {
	Identity policyIdentity
	Value    *chromepolicy.GoogleChromePolicyVersionsV1PolicyValue
}

// resolveDirectlySetPolicies does a paginated resolve filtered to directly-set policies.
func resolveDirectlySetPolicies(
	ctx context.Context,
	policiesService *chromepolicy.CustomersPoliciesService,
	customer string,
	filter string,
	targetResource string,
) ([]resolvedPolicyEntry, diag.Diagnostics) {
	var result []resolvedPolicyEntry

	req := &chromepolicy.GoogleChromePolicyVersionsV1ResolveRequest{
		PolicySchemaFilter: filter,
		PolicyTargetKey: &chromepolicy.GoogleChromePolicyVersionsV1PolicyTargetKey{
			TargetResource: targetResource,
		},
		PageSize: 1000,
	}

	for {
		var resp *chromepolicy.GoogleChromePolicyVersionsV1ResolveResponse
		err := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				var retryErr error
				resp, retryErr = policiesService.Resolve(
					fmt.Sprintf("customers/%s", customer), req,
				).Do()
				return retryErr
			},
			Timeout: chromePolicyRetryTimeout,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		for _, rp := range resp.ResolvedPolicies {
			if rp.SourceKey == nil || rp.SourceKey.TargetResource != targetResource {
				continue
			}
			result = append(result, resolvedPolicyEntry{
				Identity: identityFromResolved(rp),
				Value:    rp.Value,
			})
		}

		if resp.NextPageToken == "" {
			break
		}
		req.PageToken = resp.NextPageToken
	}

	return result, nil
}

func isNonFatalDeleteError(err error) bool {
	if err == nil {
		return false
	}
	if !transport_tpg.IsGoogleApiErrorWithCode(err, 400) {
		return false
	}
	msg := err.Error()
	nonFatalMessages := []string{
		"apps are not installed",
		"Install Type can only be inherited",
		"BatchInheritOrgUnitPolicies request must contain at least one request",
		"do not exist in Chrome Web Store and do not have a Url specified",
	}
	for _, nfm := range nonFatalMessages {
		if strings.Contains(msg, nfm) {
			return true
		}
	}
	return false
}

type policyIdentity struct {
	SchemaName           string
	AdditionalTargetKeys map[string]string
}

func (p policyIdentity) key() string {
	if len(p.AdditionalTargetKeys) == 0 {
		return p.SchemaName
	}
	sorted := make([]string, 0, len(p.AdditionalTargetKeys))
	for k, v := range p.AdditionalTargetKeys {
		sorted = append(sorted, k+"="+v)
	}
	sort.Strings(sorted)
	return p.SchemaName + "\x00" + strings.Join(sorted, ",")
}

func identityFromPolicy(p map[string]interface{}) policyIdentity {
	id := policyIdentity{
		SchemaName: p["schema"].(string),
	}
	if atk, ok := p["additional_target_keys"]; ok && atk != nil {
		raw := atk.(map[string]interface{})
		if len(raw) > 0 {
			id.AdditionalTargetKeys = make(map[string]string, len(raw))
			for k, v := range raw {
				id.AdditionalTargetKeys[k] = v.(string)
			}
		}
	}
	return id
}

func identityFromResolved(rp *chromepolicy.GoogleChromePolicyVersionsV1ResolvedPolicy) policyIdentity {
	id := policyIdentity{
		SchemaName: rp.Value.PolicySchema,
	}
	if rp.TargetKey != nil && len(rp.TargetKey.AdditionalTargetKeys) > 0 {
		id.AdditionalTargetKeys = rp.TargetKey.AdditionalTargetKeys
	}
	return id
}

func chromePoliciesResourceID(customerID string, kind chromePolicyTargetKind, targetID, filter string) string {
	return customerID + "/" + string(kind) + "/" + targetID + "/" + filter
}

// policyMapByKey converts a list of policy maps into a map keyed by identity.
func policyMapByKey(policies []interface{}) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{}, len(policies))
	for _, p := range policies {
		pol := p.(map[string]interface{})
		result[identityFromPolicy(pol).key()] = pol
	}
	return result
}

// policySetsEqual returns true if two policy maps have the same keys and values.
func policySetsEqual(a, b map[string]map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for key, aPol := range a {
		bPol, exists := b[key]
		if !exists || !policyValuesEqual(aPol, bPol) {
			return false
		}
	}
	return true
}

// policyValuesEqual compares the "value" field of two policy maps.
func policyValuesEqual(a, b map[string]interface{}) bool {
	aVals, _ := a["value"].(map[string]interface{})
	bVals, _ := b["value"].(map[string]interface{})
	if len(aVals) != len(bVals) {
		return false
	}
	for k, av := range aVals {
		bv, ok := bVals[k]
		if !ok || av != bv {
			return false
		}
	}
	return true
}

func schemaNameMatchesFilter(name, filter string) bool {
	if !strings.HasSuffix(filter, ".*") {
		return name == filter
	}
	prefix := strings.TrimSuffix(filter, "*")
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	leaf := name[len(prefix):]
	return len(leaf) > 0 && !strings.Contains(leaf, ".")
}
