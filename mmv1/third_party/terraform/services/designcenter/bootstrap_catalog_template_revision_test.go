package designcenter_test

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const bootstrapDesignCenterAppCatalogTemplateID = "tf-bootstrap-app-catalog-template"

// BootstrapDesignCenterApplicationCatalogTemplate ensures a CatalogTemplate of
// category APPLICATION_TEMPLATE exists in the bootstrapped Catalog and returns
// its catalog_template_id. The API rejects application template revision
// sources in COMPONENT_TEMPLATE catalog templates.
func BootstrapDesignCenterApplicationCatalogTemplate(t *testing.T) string {
	t.Helper()

	project, space := bootstrapDesignCenterSpaceName(t)
	template := space + "/catalogs/" + BootstrapDesignCenterCatalog(t) + "/templates/" + bootstrapDesignCenterAppCatalogTemplateID
	return bootstrapDesignCenterOnce(t, template, bootstrapDesignCenterAppCatalogTemplateID, func(config *transport_tpg.Config) string {
		ensureDesignCenterResource(t, config, project, template, "catalogTemplateId", false, map[string]interface{}{
			"displayName":      "Bootstrap Design Center Application Catalog Template",
			"description":      "Bootstrapped for Terraform acceptance tests",
			"templateCategory": "APPLICATION_TEMPLATE",
		})
		return bootstrapDesignCenterAppCatalogTemplateID
	})
}
