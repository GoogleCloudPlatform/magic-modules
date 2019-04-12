// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package azurerm



func resourceArmBatchAccount() *schema.Resource {
    return &schema.Resource{
        Create: resourceArmBatchAccountCreate,
        Read: resourceArmBatchAccountRead,
        Update: resourceArmBatchAccountUpdate,
        Delete: resourceArmBatchAccountDelete,

        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },


        Schema: map[string]*schema.Schema{
            "name": {
                Type: schema.TypeString,
                Required: true,
                ForceNew: true,
                ValidateFunc: validateAzureRMBatchAccountName,
            },

            "location": locationSchema(),

            "resource_group_name": resourceGroupNameSchema(),

            "key_vault_reference": {
                Type: schema.TypeList,
                Optional: true,
                ForceNew: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type: schema.TypeString,
                            Required: true,
                            ForceNew: true,
                        },
                        "url": {
                            Type: schema.TypeString,
                            Required: true,
                            ForceNew: true,
                        },
                    },
                },
            },

            "pool_allocation_mode": {
                Type: schema.TypeString,
                Optional: true,
                ForceNew: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(batch.BatchService),
                    string(batch.UserSubscription),
                }, false),
                Default: string(batch.BatchService),
            },

            "storage_account_id": {
                Type: schema.TypeString,
                Optional: true,
                ValidateFunc: azure.ValidateResourceIDOrEmpty,
            },

            "tags": tagsSchema(),
        },
    }
}

func resourceArmBatchAccountCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*ArmClient).batchAccountClient
    ctx := meta.(*ArmClient).StopContext

    name := d.Get("name").(string)
    resourceGroup := d.Get("resource_group_name").(string)

    if requireResourcesToBeImported {
        resp, err := client.Get(ctx, resourceGroup, name)
        if err != nil {
            if !utils.ResponseWasNotFound(resp.Response) {
                return fmt.Errorf("Error checking for present of existing Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
            }
        }
        if !utils.ResponseWasNotFound(resp.Response) {
            return tf.ImportAsExistsError("azurerm_batch_account", *resp.ID)
        }
    }

    location := azureRMNormalizeLocation(d.Get("location").(string))
    keyVaultReference := d.Get("key_vault_reference").([]interface{})
    poolAllocationMode := d.Get("pool_allocation_mode").(string)
    storageAccountId := d.Get("storage_account_id").(string)
    tags := d.Get("tags").(map[string]interface{})

    parameters := batch.AccountCreateParameters{
        Location: utils.String(location),
        AccountCreateProperties: &batch.AccountCreateProperties{
            AutoStorage: &batch.AutoStorageBaseProperties{
                StorageAccountID: utils.String(storageAccountId),
            },
            KeyVaultReference: expandArmBatchAccountKeyVaultReference(keyVaultReference),
            PoolAllocationMode: batch.PoolAllocationMode(poolAllocationMode),
        },
        Tags: expandTags(tags),
    }


    future, err := client.Create(ctx, resourceGroup, name, parameters)
    if err != nil {
        return fmt.Errorf("Error creating Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }
    if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
        return fmt.Errorf("Error waiting for creation of Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }


    resp, err := client.Get(ctx, resourceGroup, name)
    if err != nil {
        return fmt.Errorf("Error retrieving Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }
    if resp.ID == nil {
        return fmt.Errorf("Cannot read Batch Account %q (Resource Group %q) ID", name, resourceGroup)
    }
    d.SetId(*resp.ID)

    return resourceArmBatchAccountRead(d, meta)
}

func resourceArmBatchAccountRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*ArmClient).batchAccountClient
    ctx := meta.(*ArmClient).StopContext

    id, err := parseAzureResourceID(d.Id())
    if err != nil {
        return err
    }
    resourceGroup := id.ResourceGroup
    name := id.Path["batchAccounts"]

    resp, err := client.Get(ctx, resourceGroup, name)
    if err != nil {
        if utils.ResponseWasNotFound(resp.Response) {
            log.Printf("[INFO] Batch Account %q does not exist - removing from state", d.Id())
            d.SetId("")
            return nil
        }
        return fmt.Errorf("Error reading Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }


    d.Set("name", resp.Name)
    if location := resp.Location; location != nil {
        d.Set("location", azureRMNormalizeLocation(*location))
    }
    if properties := resp.AccountProperties; properties != nil {
        if err := d.Set("key_vault_reference", flattenArmBatchAccountKeyVaultReference(properties.KeyVaultReference)); err != nil {
            return fmt.Errorf("Error setting `key_vault_reference`: %+v", err)
        }
        d.Set("pool_allocation_mode", string(properties.PoolAllocationMode))
        if autoStorage := properties.AutoStorage; autoStorage != nil {
            d.Set("storage_account_id", autoStorage.StorageAccountID)
        }
    }
    flattenAndSetTags(d, resp.Tags)

    return nil
}

func resourceArmBatchAccountUpdate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*ArmClient).batchAccountClient
    ctx := meta.(*ArmClient).StopContext

    name := d.Get("name").(string)
    resourceGroup := d.Get("resource_group_name").(string)
    storageAccountId := d.Get("storage_account_id").(string)
    tags := d.Get("tags").(map[string]interface{})

    parameters := batch.AccountUpdateParameters{
        AccountUpdateProperties: &batch.AccountUpdateProperties{
            AutoStorage: &batch.AutoStorageBaseProperties{
                StorageAccountID: utils.String(storageAccountId),
            },
        },
        Tags: expandTags(tags),
    }

    if _, err := client.Update(ctx, resourceGroup, name, parameters); err != nil {
        return fmt.Errorf("Error updating Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }

    return resourceArmBatchAccountRead(d, meta)
}

func resourceArmBatchAccountDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*ArmClient).batchAccountClient
    ctx := meta.(*ArmClient).StopContext


    id, err := parseAzureResourceID(d.Id())
    if err != nil {
        return err
    }
    resourceGroup := id.ResourceGroup
    name := id.Path["batchAccounts"]

    future, err := client.Delete(ctx, resourceGroup, name)
    if err != nil {
        if response.WasNotFound(future.Response()) {
            return nil
        }
        return fmt.Errorf("Error deleting Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
    }

    if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
        if !response.WasNotFound(future.Response()) {
            return fmt.Errorf("Error waiting for deleting Batch Account %q (Resource Group %q): %+v", name, resourceGroup, err)
        }
    }

    return nil
}

func expandArmBatchAccountKeyVaultReference(input []interface{}) *batch.KeyVaultReference {
    if len(input) == 0 {
        return nil
    }
    v := input[0].(map[string]interface{})

    id := v["id"].(string)
    url := v["url"].(string)

    result := batch.KeyVaultReference{
        ID: utils.String(id),
        URL: utils.String(url),
    }
    return &result
}

func flattenArmBatchAccountKeyVaultReference(input *batch.KeyVaultReference) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    result := make(map[string]interface{})

    if id := input.ID; id != nil {
        result["id"] = *id
    }
    if url := input.URL; url != nil {
        result["url"] = *url
    }

    return []interface{}{result}
}

func validateAzureRMBatchAccountName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)
    if !regexp.MustCompile(`^[a-z0-9]+$`).MatchString(value) {
        errors = append(errors, fmt.Errorf("lowercase letters and numbers only are allowed in %q: %q", k, value))
    }

    if 3 > len(value) {
        errors = append(errors, fmt.Errorf("%q cannot be less than 3 characters: %q", k, value))
    }

    if len(value) > 24 {
        errors = append(errors, fmt.Errorf("%q cannot be longer than 24 characters: %q %d", k, value, len(value)))
    }

    return warnings, errors
}
