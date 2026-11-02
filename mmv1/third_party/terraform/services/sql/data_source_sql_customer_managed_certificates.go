// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
// ----------------------------------------------------------------------------
//
//	***     HANDWRITTEN CODE     ***
//
// ----------------------------------------------------------------------------
package sql

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSqlCustomerManagedCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSqlCustomerManagedCertificatesRead,

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the Cloud SQL instance.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the project in which the instance belongs.`,
			},
			"customer_managed_certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of customer-managed certificate bundles configured on the Cloud SQL instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bundle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Customer managed certificate bundle name.`,
						},
						"intent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Intent for the customer managed certificate bundle configuration.`,
						},
						"primary_dns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `DNS name associated with the primary instance in a replication topology.`,
						},
						"primary_ipv4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `IPv4 address for the primary DNS name.`,
						},
						"file_paths": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `File paths of the PEM-encoded certificates and key on the instance.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"server_ca_cert_file_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `File path of the PEM-encoded peer server CA certificate on the instance.`,
						},
						"client_cert_file_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `File path of the PEM-encoded peer client certificate on the instance.`,
						},
						"client_key_file_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `File path of the PEM-encoded peer client private key on the instance.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSqlCustomerManagedCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CustomerManagedCertificates: %s", err)
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	instance := d.Get("instance").(string)

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/customerManagedCertificates")
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SQL CustomerManagedCertificates in instance %q", instance))
	}

	bundles := flattenDataSourceSqlCustomerManagedCertificates(res)
	if err := d.Set("customer_managed_certificates", bundles); err != nil {
		return fmt.Errorf("Error setting customer_managed_certificates: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/customerManagedCertificates")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return nil
}

func flattenDataSourceSqlCustomerManagedCertificates(res map[string]interface{}) []map[string]interface{} {
	if res == nil {
		return nil
	}
	rawList, ok := res["customerManagedCertificates"].([]interface{})
	if !ok {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rawList))
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		bundle, _ := m["bundle"].(string)
		intent, _ := m["intent"].(string)
		dns, ipv4 := extractPrimaryDnsFields(m["primaryDnsRecord"])
		filePaths, serverCaPath, clientCertPath, clientKeyPath := extractCustomerManagedCertFilePaths(m["customerManagedCertificates"])
		result = append(result, map[string]interface{}{
			"bundle":                   bundle,
			"intent":                   intent,
			"primary_dns":              dns,
			"primary_ipv4_address":     ipv4,
			"file_paths":               filePaths,
			"server_ca_cert_file_path": serverCaPath,
			"client_cert_file_path":    clientCertPath,
			"client_key_file_path":     clientKeyPath,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i]["bundle"].(string) < result[j]["bundle"].(string)
	})

	return result
}

func init() {
	registry.Schema{
		Name:        "google_sql_customer_managed_certificates",
		ProductName: "sql",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceSqlCustomerManagedCertificates(),
	}.Register()
}
