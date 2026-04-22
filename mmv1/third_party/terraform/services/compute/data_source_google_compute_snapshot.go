package compute

import (
	"fmt"
	neturl "net/url"
	"sort"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeSnapshot() *schema.Resource {

	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeSnapshot().Schema)

	dsSchema["filter"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	dsSchema["most_recent"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "name", "filter", "most_recent", "project")

	dsSchema["name"].ExactlyOneOf = []string{"name", "filter"}
	dsSchema["filter"].ExactlyOneOf = []string{"name", "filter"}

	return &schema.Resource{
		Read:   dataSourceGoogleComputeSnapshotRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		return retrieveSnapshot(d, meta, project, v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return err
		}

		billingProject := project
		if config.UserProjectOverride {
			if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
				billingProject = bp
			}
		}

		allSnapshots := make([]map[string]interface{}, 0)
		token := ""
		for paginate := true; paginate; {
			params := neturl.Values{}
			params.Set("filter", v.(string))
			if token != "" {
				params.Set("pageToken", token)
			}
			url := fmt.Sprintf("%sprojects/%s/global/snapshots?%s", config.ComputeBasePath, project, params.Encode())
			resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: userAgent,
			})
			if err != nil {
				return fmt.Errorf("error retrieving list of snapshots: %s", err)
			}
			if items, ok := resp["items"].([]interface{}); ok {
				for _, raw := range items {
					if snap, ok := raw.(map[string]interface{}); ok {
						allSnapshots = append(allSnapshots, snap)
					}
				}
			}
			token, _ = resp["nextPageToken"].(string)
			paginate = token != ""
		}

		mostRecent := d.Get("most_recent").(bool)
		if mostRecent {
			sort.Sort(ByCreationTimestampOfSnapshot(allSnapshots))
		}

		count := len(allSnapshots)
		if count == 1 || count > 1 && mostRecent {
			return retrieveSnapshot(d, meta, project, allSnapshots[0]["name"].(string))
		}

		return fmt.Errorf("your filter has returned %d snapshot(s). Please refine your filter or set most_recent to return exactly one snapshot", len(allSnapshots))

	}

	return fmt.Errorf("one of name or filter must be set")
}

func retrieveSnapshot(d *schema.ResourceData, meta interface{}, project, name string) error {
	d.SetId("projects/" + project + "/global/snapshots/" + name)
	d.Set("name", name)
	if err := resourceComputeSnapshotRead(d, meta); err != nil {
		return err
	}
	return tpgresource.SetDataSourceLabels(d)
}

// ByCreationTimestamp implements sort.Interface for []map[string]interface{} based on
// the creationTimestamp field (RFC 3339 strings are lexicographically comparable).
type ByCreationTimestampOfSnapshot []map[string]interface{}

func (a ByCreationTimestampOfSnapshot) Len() int      { return len(a) }
func (a ByCreationTimestampOfSnapshot) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreationTimestampOfSnapshot) Less(i, j int) bool {
	ti, _ := a[i]["creationTimestamp"].(string)
	tj, _ := a[j]["creationTimestamp"].(string)
	return ti > tj
}

func init() {
	registry.Schema{
		Name:        "google_compute_snapshot",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeSnapshot(),
	}.Register()
}
