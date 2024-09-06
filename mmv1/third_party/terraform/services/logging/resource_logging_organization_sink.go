package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingOrganizationSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingOrganizationSinkCreate,
		Read:   resourceLoggingOrganizationSinkRead,
		Delete: resourceLoggingOrganizationSinkDelete,
		Update: resourceLoggingOrganizationSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("org_id"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["org_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: `The numeric ID of the organization to be exported to the sink.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "organizations/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: `Whether or not to include children organizations in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided organization are included.`,
	}
	schm.Schema["intercept_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: `Whether or not to intercept logs from child projects. If true, matching logs will not match with sinks in child resources, except _Required sinks. This sink will be visible to child resources when listing sinks.`,
	}

	return schm
}

func resourceLoggingOrganizationSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	org := d.Get("org_id").(string)
	id, sink := expandResourceLoggingSink(d, "organizations", org)
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.InterceptChildren = d.Get("intercept_children").(bool)

	// Must use a unique writer, since all destinations are in projects.
	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Organizations.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingOrganizationSinkRead(d, meta)
}

func resourceLoggingOrganizationSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, err := config.NewLoggingClient(userAgent).Organizations.Sinks.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Organization Logging Sink %s", d.Get("name").(string)))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	if err := d.Set("include_children", sink.IncludeChildren); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	if err := d.Set("intercept_children", sink.InterceptChildren); err != nil {
		return fmt.Errorf("Error setting intercept_children: %s", err)
	}

	return nil
}

func resourceLoggingOrganizationSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Organizations.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingOrganizationSinkRead(d, meta)
}

func resourceLoggingOrganizationSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
