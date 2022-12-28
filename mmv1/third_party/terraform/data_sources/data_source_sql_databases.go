package google

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func dataSourceSqlDatabases() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSqlDatabasesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project that contains the instance.`,
			},
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Cloud SQL database instance in which the database belongs.`,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Values for the field.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Name of the field.`,
						},
						"exclude_values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `The returned list would not include databases which match these values`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(resourceSQLDatabase().Schema),
				},
			},
		},
	}
}

func dataSourceSqlDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	var databases *sqladmin.DatabasesListResponse
	err = retryTimeDuration(func() (rerr error) {
		databases, rerr = config.NewSqlAdminClient(userAgent).Databases.List(project, d.Get("instance").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Databases in %q instance", d.Get("instance").(string)))
	}
	var filteredDatabases []*sqladmin.Database
	if v, ok := d.GetOk("filters"); ok {
		filteredDatabases, err = applyFilterOnDatabases(databases.Items, v.([]interface{}))
		if err != nil {
			return err
		}
	} else {
		filteredDatabases = databases.Items
	}

	if err := d.Set("databases", flattenDatabases(filteredDatabases)); err != nil {
		return fmt.Errorf("Error setting databases: %s", err)
	}

	d.SetId(setId(d, project))
	return nil
}

func applyFilterOnDatabases(databases []*sqladmin.Database, databaseFilters []interface{}) ([]*sqladmin.Database, error) {
	filteredDatabases := make([]*sqladmin.Database, 0)
	if len(databases) == 0 {
		return databases, nil
	}
	for _, d := range databases {
		include := true
		for _, f := range databaseFilters {
			if f == nil {
				continue
			}
			if !include {
				break
			}
			filter := f.(map[string]interface{})
			switch filter["name"].(string) {
			case "name":
				i, err := regexMatch(filter, d.Name, include)
				if err != nil {
					return filteredDatabases, err
				}
				include = i
			case "charset":
				i, err := regexMatch(filter, d.Charset, include)
				if err != nil {
					return filteredDatabases, err
				}
				include = i
			case "collation":
				i, err := regexMatch(filter, d.Collation, include)
				if err != nil {
					return filteredDatabases, err
				}
				include = i
			default:
				return filteredDatabases, fmt.Errorf("Invalid filter")
			}
		}
		if include {
			filteredDatabases = append(filteredDatabases, d)
		}
	}

	return filteredDatabases, nil

}

func flattenDatabases(fetchedDatabases []*sqladmin.Database) []map[string]interface{} {
	if fetchedDatabases == nil {
		return make([]map[string]interface{}, 0)
	}

	databases := make([]map[string]interface{}, 0, len(fetchedDatabases))
	for _, rawDatabase := range fetchedDatabases {
		database := make(map[string]interface{})
		database["name"] = rawDatabase.Name
		database["instance"] = rawDatabase.Instance
		database["project"] = rawDatabase.Project
		database["charset"] = rawDatabase.Charset
		database["collation"] = rawDatabase.Collation
		database["self_link"] = rawDatabase.SelfLink

		databases = append(databases, database)
	}
	return databases
}

func regexMatch(filter map[string]interface{}, field string, include bool) (bool, error) {
	b := false
	for _, regex := range filter["values"].([]interface{}) {
		match, err := regexp.MatchString(regex.(string), field)
		if err != nil {
			return include, fmt.Errorf("Invalid regex %s", regex)
		}
		b = b || match
	}
	if _, ok := filter["values"]; ok {
		if len(filter["values"].([]interface{})) > 0 {
			include = include && b
		}
	}
	//exclude has higher priority than include
	for _, regex := range filter["exclude_values"].([]interface{}) {
		match, err := regexp.MatchString(regex.(string), field)
		if err != nil {
			return include, fmt.Errorf("Invalid regex %s", regex)
		}
		if match {
			include = false
		}
	}
	return include, nil
}

func setId(d *schema.ResourceData, project string) string {
	id := fmt.Sprintf("projects/%s/instances/%s", project, d.Get("instance").(string))
	if v, ok := d.GetOk("filters"); ok {
		for i, f := range v.([]interface{}) {
			filter := f.(map[string]interface{})
			id = fmt.Sprintf("%s/%d/%s", id, i, filter["name"])
			if fv, ok := filter["values"]; ok {
				for j, val := range fv.([]interface{}) {
					id = fmt.Sprintf("%s/%d/%s", id, j, val.(string))
				}
			}
			if fev, ok := filter["exclude_values"]; ok {
				for j, val := range fev.([]interface{}) {
					id = fmt.Sprintf("%s/%d/%s", id, j, val.(string))
				}
			}

		}
	}
	return id
}
