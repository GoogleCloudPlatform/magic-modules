package certificatemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/certificatemanager/v1"
)

func DataSourceGoogleCertificateManagerCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCertificateManagerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Description: `A user-defined name of the certificate. Certificate names must be unique
The name must be 1-64 characters long, and match the regular expression [a-zA-Z][a-zA-Z0-9_-]* which means the first character must be a letter,
and all following characters must be a dash, underscore, letter or digit.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `A human-readable description of the resource.`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: `Set of label tags associated with the Certificate resource.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Certificate Manager location.`,
						},
						"managed": {
							Type:     schema.TypeList,
							Computed: true,
							Description: `Configuration and state of a Managed Certificate.
Certificate Manager provisions and renews Managed Certificates
automatically, for as long as it's authorized to do so.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_authorizations": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Authorizations that will be used for performing domain authorization. Either issuanceConfig or dnsAuthorizations should be specificed, but not both.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"domains": {
										Type:     schema.TypeList,
										Computed: true,
										Description: `The domains for which a managed SSL certificate will be generated.
Wildcard domains are only supported with DNS challenge resolution`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"issuance_config": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `The resource name for a CertificateIssuanceConfig used to configure private PKI certificates in the format projects/*/locations/*/certificateIssuanceConfigs/*.
													If this field is not set, the certificates will instead be publicly signed as documented at https://cloud.google.com/load-balancing/docs/ssl-certificates/google-managed-certs#caa.
													Either issuanceConfig or dnsAuthorizations should be specificed, but not both.`,
									},
									"authorization_attempt_info": {
										Type:     schema.TypeList,
										Computed: true,
										Description: `Detailed state of the latest authorization attempt for each domain
specified for this Managed Certificate.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"details": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Human readable explanation for reaching the state. Provided to help address the configuration issues. Not guaranteed to be stable. For programmatic access use 'failure_reason' field.`,
												},
												"domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Domain name of the authorization attempt.`,
												},
												"failure_reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Reason for failure of the authorization attempt for the domain.`,
												},
												"state": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `State of the domain for managed certificate issuance.`,
												},
											},
										},
									},
									"provisioning_issue": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Information about issues with provisioning this Managed Certificate.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"details": {
													Type:     schema.TypeString,
													Computed: true,
													Description: `Human readable explanation about the issue. Provided to help address
the configuration issues.
Not guaranteed to be stable. For programmatic access use 'reason' field.`,
												},
												"reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `Reason for provisioning failures.`,
												},
											},
										},
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A state of this Managed Certificate.`,
									},
								},
							},
						},
						"san_dnsnames": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of Subject Alternative Names of dnsName type defined in the certificate (see RFC 5280 4.2.1.6).`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The scope of the certificate.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleCertificateManagerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("error fetching project for certificate: %s", err)
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return fmt.Errorf("error fetching region for certificate: %s", err)
	}

	filter := d.Get("filter").(string)

	certificates := make([]map[string]interface{}, 0)
	certificatesList, err := config.NewCertificateManagerClient(userAgent).Projects.Locations.Certificates.List(fmt.Sprintf("projects/%s/locations/%s", project, region)).Filter(filter).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Certificates : %s %s", project, region))
	}

	for _, certificate := range certificatesList.Certificates {
		if certificate != nil {
			certificates = append(certificates, map[string]interface{}{
				"name":         certificate.Name,
				"description":  certificate.Description,
				"labels":       certificate.Labels,
				"location":     region,
				"managed":      flattenCertificateManaged(certificate.Managed),
				"san_dnsnames": certificate.SanDnsnames,
				"scope":        certificate.Scope,
			})
		}
	}

	if err := d.Set("certificates", certificates); err != nil {
		return fmt.Errorf("error setting certificates: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s/locations/%s/certificates",
		project,
		region,
	))

	return nil
}

func flattenCertificateManaged(v *certificatemanager.ManagedCertificate) interface{} {
	if v == nil {
		return nil
	}

	output := make(map[string]interface{})

	output["authorization_attempt_info"] = flattenCertificateManagedAuthorizationAttemptInfo(v.AuthorizationAttemptInfo)
	output["dns_authorizations"] = v.DnsAuthorizations
	output["domains"] = v.Domains
	output["issuance_config"] = v.IssuanceConfig
	output["state"] = v.State
	output["provisioning_issue"] = flattenCertificateManagedProvisioningIssue(v.ProvisioningIssue)

	return []interface{}{output}
}

func flattenCertificateManagedAuthorizationAttemptInfo(v []*certificatemanager.AuthorizationAttemptInfo) interface{} {
	if v == nil {
		return nil
	}

	output := make([]interface{}, 0, len(v))

	for _, authorizationAttemptInfo := range v {
		output = append(output, map[string]interface{}{
			"details":        authorizationAttemptInfo.Details,
			"domain":         authorizationAttemptInfo.Domain,
			"failure_reason": authorizationAttemptInfo.FailureReason,
			"state":          authorizationAttemptInfo.State,
		})
	}

	return output
}

func flattenCertificateManagedProvisioningIssue(v *certificatemanager.ProvisioningIssue) interface{} {
	if v == nil {
		return nil
	}

	output := make(map[string]interface{})

	output["details"] = v.Details
	output["reason"] = v.Reason

	return []interface{}{output}
}
