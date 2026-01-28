package resolvers

import (
	"strings"

	"go.uber.org/zap"

	provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ParentResourceResolver struct {
	schema *schema.Provider

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

func NewParentResourceResolver(errorLogger *zap.Logger) *ParentResourceResolver {
	return &ParentResourceResolver{
		schema:      provider.Provider(),
		errorLogger: errorLogger,
	}
}

func (r *ParentResourceResolver) ResolveParents(jsonPlan []byte) map[string][]string {
	parentToChildMap := make(map[string][]string)

	// Read elements from the resouce config
	resourceConfig, err := tfplan.ReadResourceConfigurations(jsonPlan)
	if err != nil {
		return parentToChildMap
	}

	for _, resource := range resourceConfig.RootModule.Resources {
		for _, expression := range resource.Expressions {
			reference := expression.ExpressionData.References
			if reference != nil {
				if strings.Contains(reference[0], ".id") {
					parentToChildMap[reference[1]] = append(parentToChildMap[reference[1]], resource.Type)
				} else {
					parentToChildMap[reference[0]] = append(parentToChildMap[reference[1]], resource.Type)
				}

			}
		}
	}

	return parentToChildMap
}
