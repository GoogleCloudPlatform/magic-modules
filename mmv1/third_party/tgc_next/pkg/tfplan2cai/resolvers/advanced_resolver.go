package resolvers

import (
	"slices"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"

	provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// List of keyword to filter out in order to find the iam resource parent
var filterList = []string{"etag", "policy_data", "id", "role", "members", "condition", "member"}

type AdvancedPreResolver struct {
	schema *schema.Provider

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

func NewAdvancedResolver(errorLogger *zap.Logger) *AdvancedPreResolver {
	return &AdvancedPreResolver{
		schema:      provider.Provider(),
		errorLogger: errorLogger,
	}
}

func (r *AdvancedPreResolver) Resolve(jsonPlan []byte, resourceDataMap map[string][]*models.FakeResourceDataWithMeta) map[string][]*models.FakeResourceDataWithMeta {
	// ReadChanges
	planChanges, err := tfplan.ReadResourceChanges(jsonPlan)
	if err != nil {
		return resourceDataMap
	}

	// Read elements from the resouce config
	resourceConfig, err := tfplan.ReadResourceConfigurations(jsonPlan)
	if err != nil {
		return resourceDataMap
	}

	// Keys are resource IDs, and values are lists of IAM resource addresses.
	idToAddresses := make(map[string][]string)

	for _, rc := range planChanges {
		// Silently skip non-google resources
		if !strings.HasPrefix(rc.Type, "google_") {
			continue
		}
		// Handle iam resources, build an id for each of them and group them together
		if strings.Contains(rc.Type, "iam") {
			var keys []string
			// Take all keys from Change.After and store them in a list
			if afterMap, ok := rc.Change.After.(map[string]interface{}); ok {
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
						resourceId = fmt.Sprintf("%s/%/", key, sVal)
					}
				}

				if _, ok = afterUnknownMap[key]; ok {
					resourceId = resourceId + key + "/"
					unknownValue := ""
					for _, resource := range resourceConfig.RootModule.Resources {
						if resource.Address == rc.Address {
							unknownValue = resource.Expressions[key].References[0]
						}
					}
					resourceId = resourceId + unknownValue + "/"
				}

			}
			if len(idToAddressMap[resourceId]) == 0 {
				idToAddressMap[resourceId] = []string{rc.Address}
			} else {
				idToAddressMap[resourceId] = append(idToAddressMap[resourceId], rc.Address)
			}

		}
	}

	var groupKey [][]string
	// id : [resource1, resource2]
	for _, values := range idToAddressMap {
		tempList := []string{}
		// [resource1, resource2]
		for _, value := range values {
			for key, curResource := range resourceDataMap {
				// i is index, rd is the object we need to combine later on
				if value == curResource[0].Address() {
					tempList = append(tempList, key)
				}
			}
		}
		groupKey = append(groupKey, tempList)
	}

	// Could be something like [key1, key2] [key3, key4]
	for _, addresses := range idToAddresses {
		for i := 1; i < len(row); i++ {
			resourceDataMap[row[0]] = append(resourceDataMap[row[0]], resourceDataMap[row[i]][0])
			delete(resourceDataMap, row[i])
		}
	}

	return resourceDataMap
}
