package google

import (
	"encoding/json"
	"fmt"
	"strings"

	iamcredentials "google.golang.org/api/iamcredentials/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleServiceAccountJwt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountJwtRead,
		Schema: map[string]*schema.Schema{
			"payload": {
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A JWT Claims Set that will be included in the signed JWT.`,
			},
			"target_service_account": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp("(" + strings.Join(PossibleServiceAccountNames, "|") + ")"),
			},
			"delegates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRegexp(ServiceAccountLinkRegex),
				},
			},
			"jwt": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountJwtRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	userAgent, err := generateUserAgentString(d, config.userAgent)

	if err != nil {
		return err
	}

	name := fmt.Sprintf("projects/-/serviceAccounts/%s", d.Get("target_service_account").(string))

	payload, err := json.Marshal(d.Get("payload").(map[string]interface{}))

	if err != nil {
		return fmt.Errorf("error marshaling payload JSON: %w", err)
	}

	jwtRequest := &iamcredentials.SignJwtRequest{
		Payload:   string(payload),
		Delegates: convertStringSet(d.Get("delegates").(*schema.Set)),
	}

	service := config.NewIamCredentialsClient(userAgent)

	jwtResponse, err := service.Projects.ServiceAccounts.SignJwt(name, jwtRequest).Do()

	if err != nil {
		return fmt.Errorf("error calling iamcredentials.SignJwt: %w", err)
	}

	d.SetId(name)

	if err := d.Set("jwt", jwtResponse.SignedJwt); err != nil {
		return fmt.Errorf("error setting jwt attribute: %w", err)
	}

	return nil
}
