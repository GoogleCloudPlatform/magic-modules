package resolvers

import (
	"slices"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"

	provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// List of keyword to filter out in order to find the iam resource parent
var filterList = []string{"etag", "policy_data", "id", "role", "members", "condition", "member"}

type IamAdvancedPreResolver struct {
	schema *schema.Provider

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

func NewIamAdvancedResolver(errorLogger *zap.Logger) *IamAdvancedPreResolver {
	return &IamAdvancedPreResolver{
		schema:      provider.Provider(),
		errorLogger: errorLogger,
	}
}

func (r *IamAdvancedPreResolver) Resolve(jsonPlan []byte) map[string][]*tfjson.ResourceChange {
	// Keys are resource IDs, and values are resource change objects.
	idToResourceChange := make(map[string][]*tfjson.ResourceChange)
	// ReadChanges
	planChanges, err := tfplan.ReadResourceChanges(jsonPlan)
	if err != nil {
		return idToResourceChange
	}

	// Read elements from the resouce config
	resourceConfig, err := tfplan.ReadResourceConfigurations(jsonPlan)
	if err != nil {
		return idToResourceChange
	}

	// Stores information about resources, address as key, expression as value
	addressToExpressionMap := make(map[string]map[string]*tfjson.Expression)

	for _, resource := range resourceConfig.RootModule.Resources {
		addressToExpressionMap[resource.Address] = resource.Expressions
	}

	for _, rc := range planChanges {
		// Silently skip non-google resources
		if !strings.HasPrefix(rc.Type, "google_") {
			continue
		}
		// Handle iam resources, build an id for each of them and group them together
		if strings.Contains(rc.Type, "iam_member") || strings.Contains(rc.Type, "iam_binding") || strings.Contains(rc.Type, "iam_policy") {
			var keys []string
			// Take all keys from Change.After and store them in a list
			afterMap, ok := rc.Change.After.(map[string]interface{})
			if ok {
				for k := range afterMap {
					if !slices.Contains(filterList, k) {
						keys = append(keys, k)
					}
				}
			}

			// Take all keys from Change.AfterUnknown and store them in a list
			afterUnknownMap, ok := rc.Change.AfterUnknown.(map[string]interface{})
			if ok {
				for k := range afterUnknownMap {
					if !slices.Contains(filterList, k) {
						keys = append(keys, k)
					}
				}
			}

			sort.Strings(keys)
			resourceId := ""
			// Build the id for each iam resource
			for _, key := range keys {
				// variable refers to the parent argument,
				if value, ok := afterMap[key]; ok {
					resourceId = resourceId + key + "/"
					if sVal, ok := value.(string); ok {
						resourceId = resourceId + sVal + "/"
					}
				}

				if _, ok = afterUnknownMap[key]; ok {
					resourceId = resourceId + key + "/"
					unknownValue := addressToExpressionMap[rc.Address][key].References[0]
					resourceId = resourceId + unknownValue + "/"

				}

			}

			if len(idToResourceChange[resourceId]) == 0 {
				idToResourceChange[resourceId] = []*tfjson.ResourceChange{rc}
			} else {
				idToResourceChange[resourceId] = append(idToResourceChange[resourceId], rc)
			}

		}
	}

	return idToResourceChange
}
