package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func dataSourceSqlDatabase() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceSqlDatabaseRead,

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database.`,
			},
			"charset": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The charset value. See MySQL's
            [Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
            and Postgres' [Character Set Support](https://www.postgresql.org/docs/9.6/static/multibyte.html)
            for more details and supported values. Postgres databases only support
            a value of 'UTF8' at creation time.`,
			},
			"collation": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The collation value. See MySQL's
            [Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
            and Postgres' [Collation Support](https://www.postgresql.org/docs/9.6/static/collation.html)
            for more details and supported values. Postgres databases only support
            a value of 'en_US.UTF8' at creation time.`,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	var database *sqladmin.Database
	err = retryTimeDuration(func() (rerr error) {
		database, rerr = config.NewSqlAdminClient(userAgent).Databases.Get(project, d.Get("instance").(string), d.Get("name").(string)).Do()
		return rerr
	}, d.Timeout(schema.TimeoutRead), isSqlOperationInProgressError)

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Database %q", d.Get("name").(string)))
	}
	if err := d.Set("name", database.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("instance", database.Instance); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}
	if err := d.Set("project", database.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("charset", database.Charset); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("collation", database.Collation); err != nil {
		return fmt.Errorf("Error setting collation: %s", err)
	}
	if err := d.Set("self_link", database.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, d.Get("instance").(string), d.Get("name").(string)))
	return nil
}
