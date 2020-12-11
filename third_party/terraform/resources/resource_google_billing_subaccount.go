package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"master_billing_account": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"rename_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"billing_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceBillingSubaccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	displayName := d.Get("display_name").(string)
	masterBillingAccount := d.Get("master_billing_account").(string)

	billingAccount := &cloudbilling.BillingAccount{
		DisplayName:          displayName,
		MasterBillingAccount: canonicalBillingAccountName(masterBillingAccount),
	}

	res, err := config.clientBilling.BillingAccounts.Create(billingAccount).Do()
	if err != nil {
		return fmt.Errorf("Error creating billing subaccount '%s' in master account '%s': %s", displayName, masterBillingAccount, err)
	}

	d.SetId(res.Name)
	d.Set("name", res.Name)
	d.Set("open", res.Open)
	d.Set("billing_account_id", GetResourceNameFromSelfLink(d.Get("name").(string)))

	return nil
}

func resourceBillingSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id := d.Id()

	res, err := config.clientBilling.BillingAccounts.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Billing Subaccount Not Found : %s", id))
	}

	billingAccount := res

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
		billingAccount := &cloudbilling.BillingAccount{
			DisplayName: d.Get("display_name").(string),
		}
		_, err := config.clientBilling.BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Error updating billing account : %s", d.Id()))
		}
	}
	return resourceBillingSubaccountRead(d, meta)
}

func resourceBillingSubaccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	renameOnDestroy := d.Get("rename_on_destroy").(bool)

	if renameOnDestroy {
		t := time.Now()
		billingAccount := &cloudbilling.BillingAccount{
			DisplayName: "Terraform Destroyed " + t.Format("20060102150405"),
		}
		_, err := config.clientBilling.BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Error updating billing account : %s", d.Id()))
		}
	}

	d.SetId("")

	return nil
}
