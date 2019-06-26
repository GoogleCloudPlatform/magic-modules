package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceGoogleProjectService() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleProjectServiceCreate,
		Read:   resourceGoogleProjectServiceRead,
		Delete: resourceGoogleProjectServiceDelete,
		Update: resourceGoogleProjectServiceUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"disable_dependent_services": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGoogleProjectServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	srv := d.Get("service").(string)
	err = enableServiceUsageProjectServices([]string{srv}, project, config)
	if err != nil {
		return err
	}

	id := &projectServiceId{
		project: project,
		service: srv,
	}
	d.SetId(id.terraformId())
	return resourceGoogleProjectServiceRead(d, meta)
}

func resourceGoogleProjectServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, err := parseProjectServiceId(d.Id())
	if err != nil {
		return err
	}

	enabledServices, err := readEnabledServiceUsageProjectServices(id.project, config)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	d.Set("project", id.project)

	srv := id.service
	for _, s := range enabledServices {
		if s == srv {
			d.Set("service", s)
			return nil
		}
	}

	// The service is was not found in enabled services - remove it from state
	log.Printf("[DEBUG] service %s not in enabled services for project %s, removing from state", srv, id.project)
	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if disable := d.Get("disable_on_destroy"); !(disable.(bool)) {
		log.Printf("[WARN] Project service %q disable_on_destroy is false, skip disabling service", d.Id())
		d.SetId("")
		return nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	service := d.Get("service").(string)
	disableDependencies := d.Get("disable_dependent_services").(bool)
	if err = disableServiceUsageProjectService(service, project, config, disableDependencies); err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Service %s", d.Id()))
	}

	d.SetId("")
	return nil
}

func resourceGoogleProjectServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	// This update method is no-op because the only updatable fields
	// are state/config-only, i.e. they aren't sent in requests to the API.
	return nil
}

// Parts that make up the id of a `google_project_service` resource.
// Project is included in order to allow multiple projects to enable the same service within the same Terraform state
type projectServiceId struct {
	project string
	service string
}

func (id projectServiceId) terraformId() string {
	return fmt.Sprintf("%s/%s", id.project, id.service)
}

func parseProjectServiceId(id string) (*projectServiceId, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid google_project_service id format, expecting `{project}/{service}`, found %s", id)
	}

	return &projectServiceId{parts[0], parts[1]}, nil
}
