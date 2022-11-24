package google

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const nonUniqueWriter = "serviceAccount:cloud-logs@system.gserviceaccount.com"

func resourceLoggingFolderSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingFolderSinkCreate,
		Read:   resourceLoggingFolderSinkRead,
		Delete: resourceLoggingFolderSinkDelete,
		Update: resourceLoggingFolderSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("folder"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["folder"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The folder to be exported to the sink. Note that either [FOLDER_ID] or "folders/[FOLDER_ID]" is accepted.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "folders/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: `Whether or not to include children folders in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided folder are included.`,
	}
	schm.Schema["unique_writer_identity"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: `Whether or not to create a unique identity associated with this sink. If true, then a unique service account is created and used for this sink. If you wish to publish logs across projects, you must set unique_writer_identity to true. Any requests that don't explicitly set 'uniqueWriterIdentity' to true will be rejected.`,
	}

	return schm
}

func resourceLoggingFolderSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	folder := parseFolderId(d.Get("folder"))
	id, sink := expandResourceLoggingSink(d, "folders", folder)
	sink.IncludeChildren = d.Get("include_children").(bool)
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Folders.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingFolderSinkRead(d, meta)
}

// if bigquery_options is set unique_writer_identity must be true
func resourceLoggingFolderSinkCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return resourceLoggingFolderSinkCustomizeDiffFunc(d)
}

func resourceLoggingFolderSinkCustomizeDiffFunc(diff TerraformResourceDiff) error {
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

func resourceLoggingFolderSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	sink, err := config.NewLoggingClient(userAgent).Folders.Sinks.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Logging Sink %s", d.Get("name").(string)))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	if sink.WriterIdentity != nonUniqueWriter {
		if err := d.Set("unique_writer_identity", true); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	} else {
		if err := d.Set("unique_writer_identity", false); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	}

	if err := d.Set("include_children", sink.IncludeChildren); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	return nil
}

func resourceLoggingFolderSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	// It seems the API might actually accept an update for include_children; this is not in the list of updatable
	// properties though and might break in the future. Always include the value to prevent it changing.
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.ForceSendFields = append(sink.ForceSendFields, "IncludeChildren")
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).Folders.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
