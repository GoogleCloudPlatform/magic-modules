package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func getResourcePropertiesFromSQLSelfLinkString(link string) (string, string) {
	parts := strings.Split(link, "/")
	if len(parts) >= 3 {
		return parts[len(parts)-3], parts[len(parts)-1]
	}
	return nil, nil
}

func dataSourceGoogleSQLCaCerts() *schema.Resource {
	certSchema := datasourceSchemaFromResourceSchema(resourceSqlSslCert().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleSQLCaCertsRead,

		Schema: map[string]*schema.Schema{
			"certs": {
				Type:     schema.TypeList,
				Elem:     &certSchema,
				Computed: true,
			},
			"active_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceGoogleSQLCaCertsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	var params []string
	var instance string
	if v, ok := d.GetOk("instance"); ok {
		instance = v.(string)
		params = []string{project, instance}
	} else if v, ok := d.GetOk("self_link"); ok {
		project, instance = getResourcePropertiesFromSQLSelfLinkString(v.(string))
		params = []string{project, instance}
	} else {
		return fmt.Errorf("one of instance or self_link must be set")
	}

	log.Printf("[DEBUG] Fetching CA certs from instance %s", instance)

	response, err := config.clientSqlAdmin.Service.Instances.ListServerCas(project, instance).Do()
	if err != nil {
		return fmt.Errorf("error retrieving CA certs: %s", err)
	}

	d.Set("project", project)
	d.Set("instance", instance)
	d.Set("certs", response.Certs)
	d.Set("active_version", response.ActiveVersion)
	d.SetId(strings.Join(params, "/"))

	return nil
}
