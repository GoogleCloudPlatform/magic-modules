package pubsub

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGooglePubsubSubscriptions returns the schema of the google_pubsub_subscriptions data source.
func DataSourceGooglePubsubSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGooglePubsubSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose subscriptions are listed. If it is not provided, the provider project is used.`,
			},
			"subscriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The subscriptions of the project.`,
				Elem:        &schema.Resource{Schema: dataSourceGooglePubsubSubscriptionsItemSchema()},
			},
		},
	}
}

// dataSourceGooglePubsubSubscriptionsItemSchema returns the schema of one subscription of the list, built from google_pubsub_subscription.
func dataSourceGooglePubsubSubscriptionsItemSchema() map[string]*schema.Schema {
	subscriptionSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourcePubsubSubscription().Schema)
	// These fields are left out: the API does not return them, they repeat the project argument, or labels already holds them.
	tpgresource.DeleteFieldsFromSchema(subscriptionSchema, "project", "deletion_policy", "tags", "terraform_labels", "effective_labels")
	// The description of the resource sends the reader to effective_labels, which is not kept here.
	subscriptionSchema["labels"].Description = `The labels of the subscription.`
	// A nested block gets no id from the SDK, so the id of the resource is added by hand, as in google_kms_crypto_keys.
	subscriptionSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: `An identifier for the subscription with format projects/{{project}}/subscriptions/{{name}}.`,
	}
	return subscriptionSchema
}

// dataSourceGooglePubsubSubscriptionsRead lists every subscription of the project with the list call of the google_pubsub_subscription list resource.
func dataSourceGooglePubsubSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	subscriptionSchema := dataSourceGooglePubsubSubscriptionsItemSchema()
	subscriptions := make([]map[string]interface{}, 0)
	// ListPubsubSubscriptions follows the pages and hands over each subscription as the data of a google_pubsub_subscription resource.
	err = ListPubsubSubscriptions(config, project, func(rd *schema.ResourceData) error {
		subscription := make(map[string]interface{}, len(subscriptionSchema))
		for key := range subscriptionSchema {
			switch key {
			case "id":
				subscription[key] = rd.Id()
			// In the data of a resource, labels only holds the labels of its configuration and effective_labels holds all of them.
			case "labels":
				subscription[key] = rd.Get("effective_labels")
			default:
				subscription[key] = rd.Get(key)
			}
		}
		subscriptions = append(subscriptions, subscription)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of subscriptions: %w", err)
	}

	if err := d.Set("subscriptions", subscriptions); err != nil {
		return fmt.Errorf("Error setting subscriptions: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/subscriptions", project))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_pubsub_subscriptions",
		ProductName: "pubsub",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGooglePubsubSubscriptions(),
	}.Register()
}
