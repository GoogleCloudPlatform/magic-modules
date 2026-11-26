package servicedirectory

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceServiceDirectoryNamespaces returns the schema of the google_service_directory_namespaces data source.
func DataSourceServiceDirectoryNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceDirectoryNamespacesRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// dataSourceServiceDirectoryNamespacesRead lists every namespace in the location, following pages and applying the optional filter.
func dataSourceServiceDirectoryNamespacesRead(d *schema.ResourceData, meta interface{}) error {
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

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{%location}}/namespaces")
	if err != nil {
		return err
	}

	params := map[string]string{}
	if filter := d.Get("filter").(string); filter != "" {
		params["filter"] = filter
	}

	namespaces := make([]map[string]interface{}, 0)
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
			return fmt.Errorf("Error retrieving list of namespaces: %w", err)
		}

		items, _ := res["namespaces"].([]interface{})
		for _, item := range items {
			namespace, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Error reading list of namespaces: expected an object, got %T", item)
			}
			name, _ := namespace["name"].(string)
			namespaces = append(namespaces, map[string]interface{}{
				"name":         name,
				"namespace_id": tpgresource.GetResourceNameFromSelfLink(name),
				"labels":       namespace["labels"],
			})
		}

		pageToken, _ := res["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
		params["pageToken"] = pageToken
	}

	if err := d.Set("namespaces", namespaces); err != nil {
		return fmt.Errorf("Error setting namespaces: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(id)

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_service_directory_namespaces",
		ProductName: "servicedirectory",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceServiceDirectoryNamespaces(),
	}.Register()
}
