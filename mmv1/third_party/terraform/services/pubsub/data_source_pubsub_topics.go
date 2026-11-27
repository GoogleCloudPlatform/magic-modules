package pubsub

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGooglePubsubTopics returns the schema of the google_pubsub_topics data source.
func DataSourceGooglePubsubTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGooglePubsubTopicsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose topics are listed. If it is not provided, the provider project is used.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The topics of the project.`,
				Elem:        &schema.Resource{Schema: dataSourceGooglePubsubTopicsItemSchema()},
			},
		},
	}
}

// dataSourceGooglePubsubTopicsItemSchema returns the schema of one topic of the list, built from google_pubsub_topic.
func dataSourceGooglePubsubTopicsItemSchema() map[string]*schema.Schema {
	topicSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourcePubsubTopic().Schema)
	// These are arguments of the resource that the API does not return, or that only describe its own configuration.
	tpgresource.DeleteFieldsFromSchema(topicSchema, "project", "deletion_policy", "tags", "terraform_labels", "effective_labels")
	// The description of the resource sends the reader to effective_labels, which is not kept here.
	topicSchema["labels"].Description = `The labels of the topic.`
	// A nested block gets no id from the SDK, so the id of the resource is added by hand, as in google_kms_crypto_keys.
	topicSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: `An identifier for the topic with format projects/{{project}}/topics/{{name}}.`,
	}
	return topicSchema
}

// dataSourceGooglePubsubTopicsRead lists every topic of the project with the list call of the google_pubsub_topic list resource.
func dataSourceGooglePubsubTopicsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	topicSchema := dataSourceGooglePubsubTopicsItemSchema()
	topics := make([]map[string]interface{}, 0)
	// ListPubsubTopics follows the pages and hands over each topic as the data of a google_pubsub_topic resource.
	err = ListPubsubTopics(config, project, func(rd *schema.ResourceData) error {
		topic := make(map[string]interface{}, len(topicSchema))
		for key := range topicSchema {
			switch key {
			case "id":
				topic[key] = rd.Id()
			// In the data of a resource, labels only holds the labels of its configuration and effective_labels holds all of them.
			case "labels":
				topic[key] = rd.Get("effective_labels")
			default:
				topic[key] = rd.Get(key)
			}
		}
		topics = append(topics, topic)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of topics: %w", err)
	}

	if err := d.Set("topics", topics); err != nil {
		return fmt.Errorf("Error setting topics: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/topics", project))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_pubsub_topics",
		ProductName: "pubsub",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGooglePubsubTopics(),
	}.Register()
}
