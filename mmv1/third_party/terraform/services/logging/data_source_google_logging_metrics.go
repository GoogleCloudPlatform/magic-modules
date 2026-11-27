package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGoogleLoggingMetrics returns the schema of the google_logging_metrics data source.
func DataSourceGoogleLoggingMetrics() *schema.Resource {
	metricSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceLoggingMetric().Schema)
	// project repeats the argument, and deletion_policy is a setting of the resource, not a field of a metric.
	tpgresource.DeleteFieldsFromSchema(metricSchema, "project", "deletion_policy")

	return &schema.Resource{
		Read: dataSourceGoogleLoggingMetricsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose logs-based metrics are listed. If it is not provided, the provider project is used.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The logs-based metrics of the project.`,
				Elem:        &schema.Resource{Schema: metricSchema},
			},
		},
	}
}

// dataSourceGoogleLoggingMetricsRead lists every logs-based metric of the project, following pages.
func dataSourceGoogleLoggingMetricsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/metrics")
	if err != nil {
		return err
	}

	params := map[string]string{}
	metrics := make([]map[string]interface{}, 0)
	for {
		url, err := transport_tpg.AddQueryParams(transport_tpg.BaseUrl(Product, config)+id, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving list of metrics: %w", err)
		}

		items, _ := res["metrics"].([]interface{})
		for _, item := range items {
			metric, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Error reading list of metrics: expected an object, got %T", item)
			}
			metrics = append(metrics, flattenLoggingMetricsItem(metric, d, config))
		}

		pageToken, _ := res["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
		params["pageToken"] = pageToken
	}

	if err := d.Set("metrics", metrics); err != nil {
		return fmt.Errorf("Error setting metrics: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(id)

	return nil
}

// flattenLoggingMetricsItem maps one metric from the API response with the flatteners of google_logging_metric.
func flattenLoggingMetricsItem(metric map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) map[string]interface{} {
	return map[string]interface{}{
		"name":              flattenLoggingMetricName(metric["name"], d, config),
		"description":       flattenLoggingMetricDescription(metric["description"], d, config),
		"bucket_name":       flattenLoggingMetricBucketName(metric["bucketName"], d, config),
		"disabled":          flattenLoggingMetricDisabled(metric["disabled"], d, config),
		"filter":            flattenLoggingMetricFilter(metric["filter"], d, config),
		"metric_descriptor": flattenLoggingMetricMetricDescriptor(metric["metricDescriptor"], d, config),
		"label_extractors":  flattenLoggingMetricLabelExtractors(metric["labelExtractors"], d, config),
		"value_extractor":   flattenLoggingMetricValueExtractor(metric["valueExtractor"], d, config),
		"bucket_options":    flattenLoggingMetricBucketOptions(metric["bucketOptions"], d, config),
	}
}

func init() {
	registry.Schema{
		Name:        "google_logging_metrics",
		ProductName: "logging",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleLoggingMetrics(),
	}.Register()
}
