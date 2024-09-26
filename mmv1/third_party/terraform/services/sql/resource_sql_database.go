package sql

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceSQLDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceSQLDatabaseCreate,
		Read:   resourceSQLDatabaseRead,
		Update: resourceSQLDatabaseUpdate,
		Delete: resourceSQLDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSQLDatabaseImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The name of the Cloud SQL instance. This does not include the project
ID.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The name of the database in the Cloud SQL instance.
This does not include the project ID or instance name.`,
			},
			"charset": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: tpgresource.CaseDiffSuppress,
				Description: `The charset value. See MySQL's
[Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
and Postgres' [Character Set Support](https://www.postgresql.org/docs/9.6/static/multibyte.html)
for more details and supported values. Postgres databases only support
a value of 'UTF8' at creation time.`,
			},
			"collation": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: `The collation value. See MySQL's
[Supported Character Sets and Collations](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html)
and Postgres' [Collation Support](https://www.postgresql.org/docs/9.6/static/collation.html)
for more details and supported values. Postgres databases only support
a value of 'en_US.UTF8' at creation time.`,
			},
			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The deletion policy for the database. Setting ABANDON allows the resource
to be abandoned rather than deleted. This is useful for Postgres, where databases cannot be
deleted from the API if there are users other than cloudsqlsuperuser with access. Possible
values are: "ABANDON", "DELETE". Defaults to "DELETE".`,
				Default: "DELETE",
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSQLDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	charsetProp, err := expandSQLDatabaseCharset(d.Get("charset"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("charset"); !tpgresource.IsEmptyValue(reflect.ValueOf(charsetProp)) && (ok || !reflect.DeepEqual(v, charsetProp)) {
		obj["charset"] = charsetProp
	}
	collationProp, err := expandSQLDatabaseCollation(d.Get("collation"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("collation"); !tpgresource.IsEmptyValue(reflect.ValueOf(collationProp)) && (ok || !reflect.DeepEqual(v, collationProp)) {
		obj["collation"] = collationProp
	}
	nameProp, err := expandSQLDatabaseName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	instanceProp, err := expandSQLDatabaseInstance(d.Get("instance"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("instance"); !tpgresource.IsEmptyValue(reflect.ValueOf(instanceProp)) && (ok || !reflect.DeepEqual(v, instanceProp)) {
		obj["instance"] = instanceProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "google-sql-database-instance-{{project}}-{{instance}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/databases")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Database: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error creating Database: %s", err)
	}

	// Store the ID now
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = SqlAdminOperationWaitTime(
		config, res, project, "Creating Database", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Database: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Database %q: %#v", d.Id(), res)

	return resourceSQLDatabaseRead(d, meta)
}

func resourceSQLDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(transformSQLDatabaseReadError(err), d, fmt.Sprintf("SQLDatabase %q", d.Id()))
	}

	// Explicitly set virtual fields to default values if unset
	if _, ok := d.GetOkExists("deletion_policy"); !ok {
		if err := d.Set("deletion_policy", "DELETE"); err != nil {
			return fmt.Errorf("Error setting deletion_policy: %s", err)
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}

	if err := d.Set("charset", flattenSQLDatabaseCharset(res["charset"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("collation", flattenSQLDatabaseCollation(res["collation"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("name", flattenSQLDatabaseName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("instance", flattenSQLDatabaseInstance(res["instance"], d, config)); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(res["selfLink"].(string))); err != nil {
		return fmt.Errorf("Error reading Database: %s", err)
	}

	return nil
}

func resourceSQLDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	charsetProp, err := expandSQLDatabaseCharset(d.Get("charset"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("charset"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, charsetProp)) {
		obj["charset"] = charsetProp
	}
	collationProp, err := expandSQLDatabaseCollation(d.Get("collation"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("collation"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, collationProp)) {
		obj["collation"] = collationProp
	}
	nameProp, err := expandSQLDatabaseName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	instanceProp, err := expandSQLDatabaseInstance(d.Get("instance"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("instance"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, instanceProp)) {
		obj["instance"] = instanceProp
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "google-sql-database-instance-{{project}}-{{instance}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Database %q: %#v", d.Id(), obj)
	headers := make(http.Header)

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PUT",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
		Headers:   headers,
	})

	if err != nil {
		return fmt.Errorf("Error updating Database %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Database %q: %#v", d.Id(), res)
	}

	err = SqlAdminOperationWaitTime(
		config, res, project, "Updating Database", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceSQLDatabaseRead(d, meta)
}

func resourceSQLDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Database: %s", err)
	}
	billingProject = project

	lockName, err := tpgresource.ReplaceVars(d, config, "google-sql-database-instance-{{project}}-{{instance}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	if deletionPolicy := d.Get("deletion_policy"); deletionPolicy == "ABANDON" {
		// Allows for database to be abandoned without deletion to avoid deletion failing
		// for Postgres databases in some circumstances due to existing SQL users
		return nil
	}

	log.Printf("[DEBUG] Deleting Database %q", d.Id())
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   headers,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, "Database")
	}

	err = SqlAdminOperationWaitTime(
		config, res, project, "Deleting Database", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Database %q: %#v", d.Id(), res)
	return nil
}

func resourceSQLDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)$",
		"^instances/(?P<instance>[^/]+)/databases/(?P<name>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<instance>[^/]+)/(?P<name>[^/]+)$",
		"^(?P<name>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/databases/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Explicitly set virtual fields to default values on import
	if err := d.Set("deletion_policy", "DELETE"); err != nil {
		return nil, fmt.Errorf("Error setting deletion_policy: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenSQLDatabaseCharset(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSQLDatabaseCollation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSQLDatabaseName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenSQLDatabaseInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandSQLDatabaseCharset(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSQLDatabaseCollation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSQLDatabaseName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandSQLDatabaseInstance(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}