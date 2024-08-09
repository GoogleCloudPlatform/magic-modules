package bigtable

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	btapb "google.golang.org/genproto/googleapis/bigtable/admin/v2"
)

func ResourceBigtableTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableTableCreate,
		Read:   resourceBigtableTableRead,
		Update: resourceBigtableTableUpdate,
		Delete: resourceBigtableTableDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableTableImport,
		},

		// Set a longer timeout for table creation as adding column families can be slow.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
		),
		// ----------------------------------------------------------------------
		// IMPORTANT: Do not add any additional ForceNew fields to this resource.
		// Destroying/recreating tables can lead to data loss for users.
		// ----------------------------------------------------------------------
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the table. Must be 1-50 characters and must only contain hyphens, underscores, periods, letters and numbers.`,
			},

			"column_family": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `A group of columns within a table which share a common configuration. This can be specified multiple times.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the column family.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The type of the column family.`,
						},
					},
				},
			},

			"instance_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The name of the Bigtable instance.`,
			},

			"split_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A list of predefined keys to split the table on. !> Warning: Modifying the split_keys of an existing table will cause Terraform to delete/recreate the entire google_bigtable_table resource.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROTECTED", "UNPROTECTED"}, false),
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `A field to make the table protected against data loss i.e. when set to PROTECTED, deleting the table, the column families in the table, and the instance containing the table would be prohibited. If not provided, currently deletion protection will be set to UNPROTECTED as it is the API default value.`,
			},

			"change_stream_retention": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidateDuration(),
				Description:  `Duration to retain change stream data for the table. Set to 0 to disable. Must be between 1 and 7 days.`,
			},

			"automated_backup_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_period": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidateDuration(),
							Description:  `How long the automated backups should be retained.`,
						},
						"frequency": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidateDuration(),
							Description:  `How frequently automated backups should occur.`,
						},
					},
				},
				Description: `Defines an automated backup policy for a table, specified by Retention Period and Frequency. To disable, set both Retention Period and Frequency to 0.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceBigtableTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	if err := d.Set("instance_name", instanceName); err != nil {
		return fmt.Errorf("Error setting instance_name: %s", err)
	}

	defer c.Close()

	tableId := d.Get("name").(string)
	tblConf := bigtable.TableConf{TableID: tableId}

	// Check if deletion protection is given
	// If not given, currently tblConf.DeletionProtection will be set to false in the API
	deletionProtection := d.Get("deletion_protection")
	if deletionProtection == "PROTECTED" {
		tblConf.DeletionProtection = bigtable.Protected
	} else if deletionProtection == "UNPROTECTED" {
		tblConf.DeletionProtection = bigtable.Unprotected
	}

	if changeStreamRetention, ok := d.GetOk("change_stream_retention"); ok {
		tblConf.ChangeStreamRetention, err = time.ParseDuration(changeStreamRetention.(string))
		if err != nil {
			return fmt.Errorf("Error parsing change stream retention: %s", err)
		}
	}

	if automatedBackupPolicyField, ok := d.GetOk("automated_backup_policy"); ok {
		automatedBackupPolicyElements := automatedBackupPolicyField.(*schema.Set).List()
		if len(automatedBackupPolicyElements) == 0 {
			return fmt.Errorf("Incomplete automated_backup_policy")
		} else {
			automatedBackupPolicy := automatedBackupPolicyElements[0].(map[string]interface{})
			abpRetentionPeriodField, retentionPeriodExists := automatedBackupPolicy["retention_period"]
			if !retentionPeriodExists {
				return fmt.Errorf("Automated backup policy retention period must be specified")
			}
			abpFrequencyField, frequencyExists := automatedBackupPolicy["frequency"]
			if !frequencyExists {
				return fmt.Errorf("Automated backup policy frequency must be specified")
			}
			abpRetentionPeriod, err := ParseDuration(abpRetentionPeriodField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy retention period: %s", err)
			}
			abpFrequency, err := ParseDuration(abpFrequencyField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy frequency: %s", err)
			}
			tblConf.AutomatedBackupConfig = &bigtable.TableAutomatedBackupPolicy{
				RetentionPeriod: abpRetentionPeriod,
				Frequency:       abpFrequency,
			}
		}
	}

	// Set the split keys if given.
	if v, ok := d.GetOk("split_keys"); ok {
		tblConf.SplitKeys = tpgresource.ConvertStringArr(v.([]interface{}))
	}

	// Set the column families if given.
	columnFamilies := make(map[string]bigtable.Family)
	if d.Get("column_family.#").(int) > 0 {
		columns := d.Get("column_family").(*schema.Set).List()

		for _, co := range columns {
			column := co.(map[string]interface{})

			if v, ok := column["family"]; ok {
				valueType, err := getType(column["type"])
				if err != nil {
					return err
				}
				columnFamilies[v.(string)] = bigtable.Family{
					// By default, there is no GC rules.
					GCPolicy:  bigtable.NoGcPolicy(),
					ValueType: valueType,
				}
			}
		}
	}
	tblConf.ColumnFamilies = columnFamilies

	// This method may return before the table's creation is complete - we may need to wait until
	// it exists in the future.
	// Set a longer timeout as creating table and adding column families can be pretty slow.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel() // Always call cancel.
	err = c.CreateTableFromConf(ctxWithTimeout, &tblConf)
	if err != nil {
		return fmt.Errorf("Error creating table. %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	table, err := c.TableInfo(ctx, name)
	if err != nil {
		if tpgresource.IsNotFoundGrpcError(err) {
			log.Printf("[WARN] Removing %s because it's gone", name)
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("column_family", FlattenColumnFamily(table.FamilyInfos)); err != nil {
		return fmt.Errorf("Error setting column_family: %s", err)
	}

	deletionProtection := table.DeletionProtection
	if deletionProtection == bigtable.Protected {
		if err := d.Set("deletion_protection", "PROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else if deletionProtection == bigtable.Unprotected {
		if err := d.Set("deletion_protection", "UNPROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else {
		return fmt.Errorf("Error setting deletion_protection, it should be either PROTECTED or UNPROTECTED")
	}

	changeStreamRetention := table.ChangeStreamRetention
	if changeStreamRetention != nil {
		if err := d.Set("change_stream_retention", changeStreamRetention.(time.Duration).String()); err != nil {
			return fmt.Errorf("Error setting change_stream_retention: %s", err)
		}
	}

	if table.AutomatedBackupConfig != nil {
		switch automatedBackupConfig := table.AutomatedBackupConfig.(type) {
		case *bigtable.TableAutomatedBackupPolicy:
			var tableAbp bigtable.TableAutomatedBackupPolicy = *automatedBackupConfig
			abpRetentionPeriod := tableAbp.RetentionPeriod.(time.Duration).String()
			abpFrequency := tableAbp.Frequency.(time.Duration).String()
			abp := []interface{}{
				map[string]interface{}{
					"retention_period": abpRetentionPeriod,
					"frequency":        abpFrequency,
				},
			}
			if err := d.Set("automated_backup_policy", abp); err != nil {
				return fmt.Errorf("Error setting automated_backup_policy: %s", err)
			}
		default:
			return fmt.Errorf("error: Unknown type of automated backup configuration")
		}
	}

	return nil
}

func resourceBigtableTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	defer c.Close()

	o, n := d.GetChange("column_family")
	oSet := o.(*schema.Set)
	nSet := n.(*schema.Set)
	name := d.Get("name").(string)

	// Add column families that are in new but not in old
	for _, new := range nSet.Difference(oSet).List() {
		column := new.(map[string]interface{})

		if v, ok := column["family"]; ok {
			log.Printf("[DEBUG] adding column family %q", v)
			config := bigtable.Family{
				ValueType: column["type"].(bigtable.Type),
			}
			if err := c.CreateColumnFamilyWithConfig(ctx, name, v.(string), config); err != nil {
				return fmt.Errorf("Error creating column family %q: %s", v, err)
			}
		}
	}

	// Remove column families that are in old but not in new
	for _, old := range oSet.Difference(nSet).List() {
		column := old.(map[string]interface{})

		if v, ok := column["family"]; ok {
			log.Printf("[DEBUG] removing column family %q", v)
			if err := c.DeleteColumnFamily(ctx, name, v.(string)); err != nil {
				return fmt.Errorf("Error deleting column family %q: %s", v, err)
			}
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	if d.HasChange("deletion_protection") {
		deletionProtection := d.Get("deletion_protection")
		if deletionProtection == "PROTECTED" {
			if err := c.UpdateTableWithDeletionProtection(ctxWithTimeout, name, bigtable.Protected); err != nil {
				return fmt.Errorf("Error updating deletion protection in table %v: %s", name, err)
			}
		} else if deletionProtection == "UNPROTECTED" {
			if err := c.UpdateTableWithDeletionProtection(ctxWithTimeout, name, bigtable.Unprotected); err != nil {
				return fmt.Errorf("Error updating deletion protection in table %v: %s", name, err)
			}
		}
	}

	if d.HasChange("change_stream_retention") {
		changeStreamRetention := d.Get("change_stream_retention")
		changeStream, err := time.ParseDuration(changeStreamRetention.(string))
		if err != nil {
			return fmt.Errorf("Error parsing change stream retention: %s", err)
		}
		if changeStream == 0 {
			if err := c.UpdateTableDisableChangeStream(ctxWithTimeout, name); err != nil {
				return fmt.Errorf("Error disabling change stream retention in table %v: %s", name, err)
			}
		} else {
			if err := c.UpdateTableWithChangeStream(ctxWithTimeout, name, changeStream); err != nil {
				return fmt.Errorf("Error updating change stream retention in table %v: %s", name, err)
			}
		}
	}

	if d.HasChange("automated_backup_policy") {
		automatedBackupPolicyField := d.Get("automated_backup_policy").(*schema.Set)
		automatedBackupPolicyElements := automatedBackupPolicyField.List()
		if len(automatedBackupPolicyElements) == 0 {
			return fmt.Errorf("Incomplete automated_backup_policy")
		}
		automatedBackupPolicy := automatedBackupPolicyElements[0].(map[string]interface{})
		abp := bigtable.TableAutomatedBackupPolicy{}

		abpRetentionPeriodField, retentionPeriodExists := automatedBackupPolicy["retention_period"]
		if retentionPeriodExists && abpRetentionPeriodField != "" {
			abpRetentionPeriod, err := ParseDuration(abpRetentionPeriodField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy retention period: %s", err)
			}
			abp.RetentionPeriod = abpRetentionPeriod
		}

		abpFrequencyField, frequencyExists := automatedBackupPolicy["frequency"]
		if frequencyExists && abpFrequencyField != "" {
			abpFrequency, err := ParseDuration(abpFrequencyField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy frequency: %s", err)
			}
			abp.Frequency = abpFrequency
		}

		if abp.RetentionPeriod != nil && abp.RetentionPeriod.(time.Duration) == 0 && abp.Frequency != nil && abp.Frequency.(time.Duration) == 0 {
			// Disable Automated Backups
			if err := c.UpdateTableDisableAutomatedBackupPolicy(ctxWithTimeout, name); err != nil {
				return fmt.Errorf("Error disabling automated backup configuration on table %v: %s", name, err)
			}
		} else {
			// Update Automated Backups config
			if err := c.UpdateTableWithAutomatedBackupPolicy(ctxWithTimeout, name, abp); err != nil {
				return fmt.Errorf("Error updating automated backup configuration on table %v: %s", name, err)
			}
		}
	}

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	err = c.DeleteTable(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting table. %s", err)
	}

	d.SetId("")

	return nil
}

func FlattenColumnFamily(families []bigtable.FamilyInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(families))

	for _, f := range families {
		data := make(map[string]interface{})
		data["family"] = f.Name
		//data["type"] = f.ValueType
		result = append(result, data)
	}

	return result
}

// TODO(rileykarson): Fix the stored import format after rebasing 3.0.0
func resourceBigtableTableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance_name>[^/]+)/tables/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func getType(input interface{}) (bigtable.Type, error) {
	if input == nil || input.(string) == "" {
		return nil, nil
	}
	inputType := input.(string)
	switch inputType {
	case "intsum":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.SumAggregator{},
		}, nil
	case "intmin":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.MinAggregator{},
		}, nil
	case "intmax":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.MaxAggregator{},
		}, nil
	case "inthll":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.HllppUniqueCountAggregator{},
		}, nil
	}
	unm := protojson.UnmarshalOptions{}
	output := &btapb.Type{}
	if err := unm.Unmarshal([]byte(inputType), output); err != nil {
		return nil, err
	}
	return bigtable.ProtoToType(output), nil
}
