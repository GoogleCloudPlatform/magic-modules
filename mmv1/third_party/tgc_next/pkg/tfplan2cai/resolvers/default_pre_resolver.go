package resolvers

import (
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/models"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/tfplan"

	provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ErrDuplicateAsset = errors.New("duplicate asset")

type DefaultPreResolver struct {
	schema *schema.Provider

	// For logging error / status information that doesn't warrant an outright failure
	errorLogger *zap.Logger
}

func NewDefaultPreResolver(errorLogger *zap.Logger) *DefaultPreResolver {
	return &DefaultPreResolver{
		schema:      provider.Provider(),
		errorLogger: errorLogger,
	}
}

func (r *DefaultPreResolver) Resolve(jsonPlan []byte) map[string][]*models.FakeResourceDataWithMeta {
	// ReadResourceChanges
	changes, err := tfplan.ReadResourceChanges(jsonPlan)
	if err != nil {
		return nil
	}

	return r.AddResourceChanges(changes)
}

// AddResourceChange processes the resource changes in two stages:
// 1. Process deletions (fetching canonical resources from GCP as necessary)
// 2. Process creates, updates, and no-ops (fetching canonical resources from GCP as necessary)
// This will give us a deterministic end result even in cases where for example
// an IAM Binding and Member conflict with each other, but one is replacing the
// other.
func (r *DefaultPreResolver) AddResourceChanges(changes []*tfjson.ResourceChange) map[string][]*models.FakeResourceDataWithMeta {
	resourceDataMap := make(map[string][]*models.FakeResourceDataWithMeta, 0)

	for _, rc := range changes {
		// Silently skip non-google resources
		if !strings.HasPrefix(rc.Type, "google_") {
			continue
		}

		// Skip resources not found in the google beta provider's schema
		if _, ok := r.schema.ResourcesMap[rc.Type]; !ok {
			r.errorLogger.Debug(fmt.Sprintf("%s: resource type not found in google beta provider: %s.", rc.Address, rc.Type))
			continue
		}

		var resourceData *models.FakeResourceDataWithMeta
		resource := r.schema.ResourcesMap[rc.Type]
		if tfplan.IsCreate(rc) || tfplan.IsUpdate(rc) || tfplan.IsDeleteCreate(rc) {
			resourceData = models.NewFakeResourceDataWithMeta(
				rc.Type,
				resource.Schema,
				rc.Change.After.(map[string]interface{}),
				false,
				rc.Address,
			)
		} else if tfplan.IsDelete(rc) {
			resourceData = models.NewFakeResourceDataWithMeta(
				rc.Type,
				resource.Schema,
				rc.Change.Before.(map[string]interface{}),
				true,
				rc.Address,
			)
		} else {
			continue
		}

		// TODO: handle the address of iam resources
		if exist := resourceDataMap[rc.Address]; exist == nil {
			resourceDataMap[rc.Address] = make([]*models.FakeResourceDataWithMeta, 0)
		}
		resourceDataMap[rc.Address] = append(resourceDataMap[rc.Address], resourceData)
	}

	return resourceDataMap
}
