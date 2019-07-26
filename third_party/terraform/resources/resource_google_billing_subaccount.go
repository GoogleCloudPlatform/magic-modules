package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudbilling/v1"
)

func resourceBillingSubaccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceBillingSubaccountCreate,
		Read:   resourceBillingSubaccountRead,
		Delete: resourceBillingSubaccountDelete,
		Update: resourceBillingSubaccountUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_billing_account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == canonicalBillingAccountName(old) {
						return false
					}
					return true
				},
			},
			"billing_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rename_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceBillingSubaccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	billingAccount := new(cloudbilling.BillingAccount)
	displayName := d.Get("display_name").(string)
	masterBillingAccount := d.Get("master_billing_account").(string)
	billingAccount.DisplayName = displayName
	billingAccount.MasterBillingAccount = canonicalBillingAccountName(masterBillingAccount)

	resp, err := config.clientBilling.BillingAccounts.Create(billingAccount).Do()
	if err != nil {
		return fmt.Errorf("Error creating billing subaccount '%s' in master account '%s': %s", displayName, masterBillingAccount, err)
	}

	d.SetId(resp.Name)
	d.Set("name", resp.Name)
	d.Set("open", resp.Open)
	d.Set("billing_account_id", strings.TrimPrefix(d.Get("name").(string), "billingAccounts/"))

	return nil
}

func resourceBillingSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()

	var billingAccount *cloudbilling.BillingAccount
	resp, err := config.clientBilling.BillingAccounts.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Billing Account Not Found : %s", id))
	}

	billingAccount = resp

	d.Set("name", billingAccount.Name)
	d.Set("display_name", billingAccount.DisplayName)
	d.Set("open", billingAccount.Open)
	d.Set("master_billing_account", billingAccount.MasterBillingAccount)
	d.Set("billing_account_id", strings.TrimPrefix(d.Get("name").(string), "billingAccounts/"))

	return nil
}

func resourceBillingSubaccountUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if ok := d.HasChange("display_name"); ok {
		billingAccount := new(cloudbilling.BillingAccount)
		billingAccount.DisplayName = d.Get("display_name").(string)
		_, err := config.clientBilling.BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return fmt.Errorf("Could not update billing account display name for: %s", d.Id())
		}
	}
	return resourceBillingSubaccountRead(d, meta)
}

func resourceBillingSubaccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	renameOnDestroy := d.Get("rename_on_destroy").(bool)

	if renameOnDestroy {
		t := time.Now()
		billingAccount := new(cloudbilling.BillingAccount)
		billingAccount.DisplayName = "Terraform Destroyed " + t.Format("20060102150405")
		_, err := config.clientBilling.BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return fmt.Errorf("Could not update billing account display name for: %s", d.Id())
		}
	}

	d.SetId("")

	return nil
}
