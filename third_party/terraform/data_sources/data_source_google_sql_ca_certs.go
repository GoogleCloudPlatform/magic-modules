package google

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func dataSourceGoogleSQLCaCerts() *schema.Resource {
	certSchema := datasourceSchemaFromResourceSchema(resourceSqlSslCert().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleSQLCaCertsRead,

		Schema: map[string]*schema.Schema{
			"active_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certs": {
				Type:     schema.TypeList,
				Elem:     &certSchema,
				Computed: true,
			},
			"instance": {
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
			"self_link": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleSQLCaCertsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var project, instance string
	if v, ok := d.GetOk("instance"); ok {
		p, err := getProject(d, config)
		if err != nil {
			return err
		}
		project = p
		instance = v.(string)
	} else if selfLink, ok := d.GetOk("self_link"); ok {
		fv, err := parseProjectFieldValue("instances", selfLink, "project", d, config, false)
		if err != nil {
			return err
		}
		project = fv.Project
		instance = fv.Name
	} else {
		return fmt.Errorf("one of instance or self_link must be set")
	}

	log.Printf("[DEBUG] Fetching CA certs from instance %s", instance)

	response, err := config.clientSqlAdmin.Service.Instances.ListServerCas(project, instance).Do()
	if err != nil {
		return fmt.Errorf("error retrieving CA certs: %s", err)
	}

	log.Printf("[DEBUG] Fetched CA certs from instance %s", instance)

	d.Set("project", project)
	d.Set("instance", instance)
	d.Set("certs", flattenServerCaCerts(response.Certs))
	d.Set("active_version", response.ActiveVersion)
	d.SetId(fmt.Sprintf("projects/%s/instance/%s", project, instance))

	return nil
}
