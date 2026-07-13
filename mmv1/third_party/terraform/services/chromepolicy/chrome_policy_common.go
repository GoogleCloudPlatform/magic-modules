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

	"github.com/hashicorp/terraform-provider-google/google/registry"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const chromePolicyRequestTimeout = 5 * time.Minute

// Target kind

type chromePolicyTargetKind string

const (
	targetOrgUnit chromePolicyTargetKind = "orgunits"
	targetGroup   chromePolicyTargetKind = "groups"
)

func chromePolicyTargetResource(kind chromePolicyTargetKind, id string) string {
	return string(kind) + "/" + id
}

// API client wrapper

// chromePolicyAPI bundles the dependencies needed to call the Chrome Policy
// REST endpoints via transport_tpg.SendRequest.
type chromePolicyAPI struct {
	config    *transport_tpg.Config
	userAgent string
	customer  string
}

func newChromePolicyAPI(config *transport_tpg.Config, userAgent, customer string) *chromePolicyAPI {
	return &chromePolicyAPI{config: config, userAgent: userAgent, customer: customer}
}

func (a *chromePolicyAPI) url(path string) string {
	return transport_tpg.BaseUrl(registry.GetProduct("chromepolicy"), a.config) + path
}

func (a *chromePolicyAPI) send(method, url string, body map[string]any) (map[string]interface{}, error) {
	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    a.config,
		Method:    method,
		RawURL:    url,
		UserAgent: a.userAgent,
		Body:      body,
		Timeout:   chromePolicyRequestTimeout,
	})
}

// Policy target key

func buildPolicyTargetKey(targetResource string, additionalTargetKeys map[string]string) map[string]any {
	key := map[string]any{"targetResource": targetResource}
	if len(additionalTargetKeys) > 0 {
		key["additionalTargetKeys"] = additionalTargetKeys
	}
	return key
}

func policyAdditionalTargetKeys(pol map[string]interface{}) map[string]string {
	atk, ok := pol["additional_target_keys"]
	if !ok || atk == nil {
		return nil
	}
	raw, ok := atk.(map[string]interface{})
	if !ok || len(raw) == 0 {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = v.(string)
	}
	return out
}

// Schema cache

// policySchema only deserializes the parts of the API response we care about.
type policySchema struct {
	AdditionalTargetKeyNames []policySchemaTargetKeyName `json:"additionalTargetKeyNames"`
	Definition               policySchemaDefinition      `json:"definition"`
}

type policySchemaTargetKeyName struct {
	Key string `json:"key"`
}

type policySchemaDefinition struct {
	MessageType []policySchemaMessageType `json:"messageType"`
}

type policySchemaMessageType struct {
	Field []policySchemaField `json:"field"`
}

type policySchemaField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type schemaCache struct {
	api   *chromePolicyAPI
	cache map[string]*policySchema
}

func newSchemaCache(api *chromePolicyAPI) *schemaCache {
	return &schemaCache{api: api, cache: make(map[string]*policySchema)}
}

func (sc *schemaCache) get(_ context.Context, schemaName string) (*policySchema, error) {
	if cached, ok := sc.cache[schemaName]; ok {
		return cached, nil
	}
	url := sc.api.url(fmt.Sprintf("customers/%s/policySchemas/%s", sc.api.customer, schemaName))
	res, err := sc.api.send("GET", url, nil)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	var ps policySchema
	if err := json.Unmarshal(b, &ps); err != nil {
		return nil, err
	}
	sc.cache[schemaName] = &ps
	return &ps, nil
}

// Validation

func validatePolicies(ctx context.Context, policies []interface{}, sc *schemaCache) error {
	for _, policy := range policies {
		pol := policy.(map[string]interface{})
		schemaName := pol["schema"].(string)
		schemaValues := pol["value"].(map[string]interface{})

		schemaDef, err := sc.get(ctx, schemaName)
		if err != nil {
			return err
		}
		if len(schemaDef.Definition.MessageType) == 0 {
			return fmt.Errorf("schema definition (%s) is empty", schemaName)
		}

		fieldMap := make(map[string]policySchemaField)
		for _, mt := range schemaDef.Definition.MessageType {
			for _, f := range mt.Field {
				fieldMap[f.Name] = f
			}
		}

		for fieldName, jsonVal := range schemaValues {
			field, ok := fieldMap[fieldName]
			if !ok {
				return fmt.Errorf("field %q is not found in schema %s", fieldName, schemaName)
			}
			var val interface{}
			if err := json.Unmarshal([]byte(jsonVal.(string)), &val); err != nil {
				return err
			}
			if field.Label == "LABEL_REPEATED" {
				arr, ok := val.([]interface{})
				if !ok {
					return fmt.Errorf("value for %s.%s must be an array (got %T)", schemaName, fieldName, val)
				}
				for _, item := range arr {
					if !validatePolicyFieldValueType(field.Type, item) {
						return fmt.Errorf("array element in %s.%s has incorrect type (expected %s)", schemaName, fieldName, field.Type)
					}
				}
			} else if !validatePolicyFieldValueType(field.Type, val) {
				return fmt.Errorf("value for %s.%s has incorrect type (expected %s)", schemaName, fieldName, field.Type)
			}
		}

		perPolicyATK := policyAdditionalTargetKeys(pol)
		if len(schemaDef.AdditionalTargetKeyNames) > 0 && len(perPolicyATK) == 0 {
			return fmt.Errorf("schema %s requires additional_target_keys", schemaName)
		}
		if len(schemaDef.AdditionalTargetKeyNames) == 0 && len(perPolicyATK) > 0 {
			return fmt.Errorf("schema %s does not support additional_target_keys", schemaName)
		}
		if len(perPolicyATK) > 0 {
			allowed := make(map[string]bool, len(schemaDef.AdditionalTargetKeyNames))
			for _, tkn := range schemaDef.AdditionalTargetKeyNames {
				allowed[tkn.Key] = true
			}
			for k := range perPolicyATK {
				if !allowed[k] {
					return fmt.Errorf("additional_target_key %q is not valid for schema %s", k, schemaName)
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

// Batch modify

func batchModifyPolicies(_ context.Context, api *chromePolicyAPI, targetResource string, policies []interface{}) error {
	requests := make([]map[string]any, 0, len(policies))
	for _, p := range policies {
		pol := p.(map[string]interface{})
		schemaName := pol["schema"].(string)
		schemaValues := pol["value"].(map[string]interface{})

		// Per-field values are already JSON-encoded strings; embed them raw to avoid an unmarshal/marshal round trip.
		valuesRaw := make(map[string]json.RawMessage, len(schemaValues))
		updateKeys := make([]string, 0, len(schemaValues))
		for k, v := range schemaValues {
			valuesRaw[k] = json.RawMessage(v.(string))
			updateKeys = append(updateKeys, k)
		}
		requests = append(requests, map[string]any{
			"policyTargetKey": buildPolicyTargetKey(targetResource, policyAdditionalTargetKeys(pol)),
			"policyValue": map[string]any{
				"policySchema": schemaName,
				"value":        valuesRaw,
			},
			"updateMask": strings.Join(updateKeys, ","),
		})
	}

	url := api.url(fmt.Sprintf("customers/%s/policies/orgunits:batchModify", api.customer))
	_, err := api.send("POST", url, map[string]any{"requests": requests})
	return err
}

// Inherit (delete-style)

func inheritPolicies(_ context.Context, api *chromePolicyAPI, requests []map[string]any) error {
	if len(requests) == 0 {
		return nil
	}
	url := api.url(fmt.Sprintf("customers/%s/policies/orgunits:batchInherit", api.customer))
	_, err := api.send("POST", url, map[string]any{"requests": requests})
	if err != nil {
		if isNonFatalDeleteError(err) {
			log.Printf("[DEBUG] Non-fatal error during policy inherit: %v", err)
			return nil
		}
		return err
	}
	return nil
}

func buildInheritRequest(targetResource, schemaName string, additionalTargetKeys map[string]string) map[string]any {
	return map[string]any{
		"policyTargetKey": buildPolicyTargetKey(targetResource, additionalTargetKeys),
		"policySchema":    schemaName,
	}
}

// Resolve

type resolvedPolicyEntry struct {
	Identity policyIdentity
	Schema   string
	Value    map[string]interface{}
}

// resolveDirectlySetPolicies pages through the resolve endpoint and keeps only policies
// directly set on the target (i.e. not inherited from a parent).
func resolveDirectlySetPolicies(_ context.Context, api *chromePolicyAPI, filter, targetResource string) ([]resolvedPolicyEntry, error) {
	var result []resolvedPolicyEntry
	url := api.url(fmt.Sprintf("customers/%s/policies:resolve", api.customer))

	body := map[string]any{
		"policySchemaFilter": filter,
		"policyTargetKey":    map[string]any{"targetResource": targetResource},
		"pageSize":           1000,
	}

	for {
		res, err := api.send("POST", url, body)
		if err != nil {
			return nil, err
		}

		policies, _ := res["resolvedPolicies"].([]interface{})
		for _, p := range policies {
			rp, _ := p.(map[string]interface{})
			sourceKey, _ := rp["sourceKey"].(map[string]interface{})
			if sourceKey == nil {
				continue
			}
			if tr, _ := sourceKey["targetResource"].(string); tr != targetResource {
				continue
			}
			value, _ := rp["value"].(map[string]interface{})
			schemaName, _ := value["policySchema"].(string)
			rawValue, _ := value["value"].(map[string]interface{})

			id := policyIdentity{SchemaName: schemaName}
			if tk, _ := rp["targetKey"].(map[string]interface{}); tk != nil {
				if atk, _ := tk["additionalTargetKeys"].(map[string]interface{}); len(atk) > 0 {
					id.AdditionalTargetKeys = make(map[string]string, len(atk))
					for k, v := range atk {
						id.AdditionalTargetKeys[k] = v.(string)
					}
				}
			}

			result = append(result, resolvedPolicyEntry{
				Identity: id,
				Schema:   schemaName,
				Value:    rawValue,
			})
		}

		nextPageToken, _ := res["nextPageToken"].(string)
		if nextPageToken == "" {
			break
		}
		body["pageToken"] = nextPageToken
	}

	return result, nil
}

// Errors

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

// Identity

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
	return policyIdentity{
		SchemaName:           p["schema"].(string),
		AdditionalTargetKeys: policyAdditionalTargetKeys(p),
	}
}

// Resource ID + set diffing

func chromePoliciesResourceID(customerID string, kind chromePolicyTargetKind, targetID, filter string) string {
	return customerID + "/" + string(kind) + "/" + targetID + "/" + filter
}

func policyMapByKey(policies []interface{}) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{}, len(policies))
	for _, p := range policies {
		pol := p.(map[string]interface{})
		result[identityFromPolicy(pol).key()] = pol
	}
	return result
}

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
