package logging

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const nonUniqueWriterAccount = "serviceAccount:cloud-logs@system.gserviceaccount.com"

func ResourceLoggingProjectSink() *schema.Resource {
	schm := &schema.Resource{
		Create:        resourceLoggingProjectSinkCreate,
		Read:          resourceLoggingProjectSinkRead,
		Delete:        resourceLoggingProjectSinkDelete,
		Update:        resourceLoggingProjectSinkUpdate,
		Schema:        resourceLoggingSinkSchema(),
		CustomizeDiff: resourceLoggingProjectSinkCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("project"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["project"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The ID of the project to create the sink in. If omitted, the project associated with the provider is used.`,
	}
	return schm
}

func resourceLoggingProjectSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	id, sink := expandResourceLoggingSink(d, "projects", project)
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())

	return resourceLoggingProjectSinkRead(d, meta)
}

// if bigquery_options is set unique_writer_identity must be true
func resourceLoggingProjectSinkCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return resourceLoggingProjectSinkCustomizeDiffFunc(d)
}

func resourceLoggingProjectSinkCustomizeDiffFunc(diff tpgresource.TerraformResourceDiff) error {
	if !diff.HasChange("bigquery_options.#") {
		return nil
	}

	bigqueryOptions := diff.Get("bigquery_options.#").(int)
	if bigqueryOptions > 0 {
		uwi := diff.Get("unique_writer_identity")
		if !uwi.(bool) {
			return errors.New("unique_writer_identity must be true when bigquery_options is supplied")
		}
	}
	return nil
}

func resourceLoggingProjectSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	sink, err := config.NewLoggingClient(userAgent).Projects.Sinks.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Project Logging Sink %s", d.Get("name").(string)))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	if sink.WriterIdentity != nonUniqueWriterAccount {
		if err := d.Set("unique_writer_identity", true); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	} else {
		if err := d.Set("unique_writer_identity", false); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	}
	return nil
}

func resourceLoggingProjectSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	return resourceLoggingProjectSinkRead(d, meta)
}

func resourceLoggingProjectSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
