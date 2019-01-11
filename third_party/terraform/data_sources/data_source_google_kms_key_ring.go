package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleKmsKeyRing() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsKeyRingRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGoogleKmsKeyRingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	keyRingId := &kmsKeyRingId{
		Name:     d.Get("name").(string),
		Location: d.Get("location").(string),
		Project:  project,
	}
	log.Printf("[DEBUG] Executing read for KMS KeyRing %s", keyRingId.keyRingId())

	keyRing, err := config.clientKms.Projects.Locations.KeyRings.Get(keyRingId.keyRingId()).Do()

	if err != nil {
		return fmt.Errorf("Error reading KeyRing: %s", err)
	}

	d.Set("project", project)
	d.Set("self_link", keyRing.Name)

	return nil
}
