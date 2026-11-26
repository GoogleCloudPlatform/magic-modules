package logging

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/logging/v2"
)

// DataSourceGoogleLoggingProjectSinks returns the schema of the google_logging_project_sinks data source.
func DataSourceGoogleLoggingProjectSinks() *schema.Resource {
	sinkSchema := tpgresource.DatasourceSchemaFromResourceSchema(resourceLoggingSinkSchema())
	// deletion_policy is a setting of the sink resources, not a field of a sink.
	delete(sinkSchema, "deletion_policy")

	return &schema.Resource{
		Read: dataSourceGoogleLoggingProjectSinksRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose sinks are listed. If it is not provided, the provider project is used.`,
			},
			"sinks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The sinks owned by the project, including the built-in _Required and _Default sinks.`,
				Elem:        &schema.Resource{Schema: sinkSchema},
			},
		},
	}
}

// dataSourceGoogleLoggingProjectSinksRead lists every sink owned by the project, following pages.
func dataSourceGoogleLoggingProjectSinksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	parent := fmt.Sprintf("projects/%s", project)

	sinks := make([]map[string]interface{}, 0)
	err = NewClient(config, userAgent).Sinks.List(parent).Pages(context.Background(), func(page *logging.ListSinksResponse) error {
		for _, sink := range page.Sinks {
			sinks = append(sinks, map[string]interface{}{
				"name":             sink.Name,
				"destination":      sink.Destination,
				"filter":           sink.Filter,
				"description":      sink.Description,
				"disabled":         sink.Disabled,
				"writer_identity":  sink.WriterIdentity,
				"exclusions":       flattenLoggingSinkExclusion(sink.Exclusions),
				"bigquery_options": flattenLoggingSinkBigqueryOptions(sink.BigqueryOptions),
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of sinks: %w", err)
	}

	if err := d.Set("sinks", sinks); err != nil {
		return fmt.Errorf("Error setting sinks: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(parent + "/sinks")

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_logging_project_sinks",
		ProductName: "logging",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleLoggingProjectSinks(),
	}.Register()
}
