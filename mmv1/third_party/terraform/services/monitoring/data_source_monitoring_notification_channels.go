package monitoring

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceMonitoringNotificationChannels returns the schema of the google_monitoring_notification_channels data source.
func DataSourceMonitoringNotificationChannels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMonitoringNotificationChannelsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose notification channels are listed. If it is not provided, the provider project is used.`,
			},
			"notification_channels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The notification channels of the project.`,
				Elem:        &schema.Resource{Schema: dataSourceMonitoringNotificationChannelsItemSchema()},
			},
		},
	}
}

// dataSourceMonitoringNotificationChannelsItemSchema returns the schema of one channel of the list, built from google_monitoring_notification_channel.
func dataSourceMonitoringNotificationChannelsItemSchema() map[string]*schema.Schema {
	channelSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMonitoringNotificationChannel().Schema)
	// These fields are left out: project repeats the argument, and the others are settings of the resource that the API never returns.
	tpgresource.DeleteFieldsFromSchema(channelSchema, "project", "deletion_policy", "force_delete", "sensitive_labels")
	// A nested block gets no id from the SDK, so the id of the resource is added by hand, as in google_kms_crypto_keys.
	channelSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: `An identifier for the channel with format projects/{{project}}/notificationChannels/{{channel_id}}, the same value as name.`,
	}
	return channelSchema
}

// dataSourceMonitoringNotificationChannelsRead lists every channel of the project and maps each one the way the google_monitoring_notification_channel resource reads it.
func dataSourceMonitoringNotificationChannelsRead(d *schema.ResourceData, meta interface{}) error {
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

	id := fmt.Sprintf("projects/%s/notificationChannels", project)
	url := transport_tpg.BaseUrl(Product, config) + "v3/" + id

	channelSchema := dataSourceMonitoringNotificationChannelsItemSchema()
	channels := make([]map[string]interface{}, 0)
	// ListPages follows the pages and hands over each channel as the data of a google_monitoring_notification_channel resource.
	err = transport_tpg.ListPages(transport_tpg.ListPagesOptions{
		Config:         config,
		TempData:       ResourceMonitoringNotificationChannel().Data(&terraform.InstanceState{}),
		Resource:       ResourceMonitoringNotificationChannel(),
		ListURL:        url,
		BillingProject: billingProject,
		UserAgent:      userAgent,
		ItemName:       "notificationChannels",
		Flattener: func(res map[string]interface{}, rd *schema.ResourceData, config *transport_tpg.Config) error {
			res, err := resourceMonitoringNotificationChannelDecoder(rd, config, res)
			if err != nil {
				return err
			}
			if res == nil {
				return fmt.Errorf("error decoding NotificationChannel from list response")
			}
			return ResourceMonitoringNotificationChannelFlatten(rd, config, res, config, project, userAgent, billingProject, url, nil)
		},
		Callback: func(rd *schema.ResourceData) error {
			channel := make(map[string]interface{}, len(channelSchema))
			for key := range channelSchema {
				if key == "id" {
					// The id of a google_monitoring_notification_channel is its name.
					channel[key] = rd.Get("name")
				} else {
					channel[key] = rd.Get(key)
				}
			}
			channels = append(channels, channel)
			return nil
		},
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of notification channels: %w", err)
	}

	if err := d.Set("notification_channels", channels); err != nil {
		return fmt.Errorf("Error setting notification_channels: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(id)

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_monitoring_notification_channels",
		ProductName: "monitoring",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceMonitoringNotificationChannels(),
	}.Register()
}
