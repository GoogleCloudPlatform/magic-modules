// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"google.golang.org/api/iam/v1"
)

func DataSourceGoogleServiceAccounts() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleServiceAccountsRead,
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidateRegexCompiles(),
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unique_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleServiceAccountsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for service accounts: %s", err)
	}

	prefix := d.Get("prefix").(string)
	regexPattern := d.Get("regex").(string)

	var regex *regexp.Regexp
	if regexPattern != "" {
		regex, err = regexp.Compile(regexPattern)
		if err != nil {
			return fmt.Errorf("Invalid regex pattern: %s", err)
		}
	}

	accounts := make([]map[string]interface{}, 0)

	request := config.NewIamClient(userAgent).Projects.ServiceAccounts.List("projects/" + project)

	err = request.Pages(context.Background(), func(accountList *iam.ListServiceAccountsResponse) error {
		for _, account := range accountList.Accounts {
			accountId := strings.Split(account.Email, "@")[0]

			if prefix != "" && !strings.HasPrefix(accountId, prefix) {
				continue
			}
			if regex != nil && !regex.MatchString(account.Email) {
				continue
			}

			accounts = append(accounts, map[string]interface{}{
				"account_id":   accountId,
				"disabled":     account.Disabled,
				"email":        account.Email,
				"display_name": account.DisplayName,
				"member":       "serviceAccount:" + account.Email,
				"name":         account.Name,
				"unique_id":    account.UniqueId,
			})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error retrieving service accounts: %s", err)
	}

	if err := d.Set("accounts", accounts); err != nil {
		return fmt.Errorf("Error setting service accounts: %s", err)
	}

	idParts := []string{"projects", project}

	if prefix != "" {
		idParts = append(idParts, "prefix/"+prefix)
	}
	if regexPattern != "" {
		idParts = append(idParts, "regex/"+regexPattern)
	}

	// Set the ID dynamically based on the provided attributes
	d.SetId(strings.Join(idParts, "/"))

	return nil
}
