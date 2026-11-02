// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
// ----------------------------------------------------------------------------
//
//	***     HANDWRITTEN CODE     ***
//
// ----------------------------------------------------------------------------
package sql

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

var customerManagedCertContentKeys = []string{
	"server_ca_cert",
	"client_cert",
	"client_key",
}

func ResourceSqlCustomerManagedCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlCustomerManagedCertificateCreate,
		Read:   resourceSqlCustomerManagedCertificateRead,
		Update: resourceSqlCustomerManagedCertificateUpdate,
		Delete: resourceSqlCustomerManagedCertificateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSqlCustomerManagedCertificateImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			tpgresource.DefaultProviderDeletionPolicy("DELETE"),
		),

		Schema: map[string]*schema.Schema{
			"bundle": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`), "bundle must be 1 to 10 alphanumeric characters"),
				Description:  `Customer managed certificate bundle name. Constrained to 1-10 alphanumeric characters.`,
			},

			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the Cloud SQL instance. This does not include the project ID.`,
			},

			"intent": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PG_LOGICAL",
				ValidateFunc: verify.ValidateEnum([]string{"PG_LOGICAL", "NATIVE_REPLICATION", ""}),
				Description:  `Intent for the customer managed certificate bundle configuration. Possible values: ["PG_LOGICAL", "NATIVE_REPLICATION"]. Defaults to "PG_LOGICAL".`,
			},

			"server_ca_cert": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				AtLeastOneOf: customerManagedCertContentKeys,
				Description:  `Input only. The PEM-encoded peer server CA certificate content.`,
			},

			"client_cert": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				AtLeastOneOf: customerManagedCertContentKeys,
				Description:  `Input only. The PEM-encoded peer client certificate content.`,
			},

			"client_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				AtLeastOneOf: customerManagedCertContentKeys,
				Description:  `Input only. The PEM-encoded peer client private key content.`,
			},

			"primary_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `DNS name associated with the primary instance in a replication topology.`,
			},

			"primary_ipv4_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `IPv4 address for the primary DNS name.`,
			},

			"file_paths": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Output only. File paths of the PEM-encoded certificates and key on the instance.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"server_ca_cert_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. File path of the PEM-encoded peer server CA certificate on the instance.`,
			},

			"client_cert_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. File path of the PEM-encoded peer client certificate on the instance.`,
			},

			"client_key_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. File path of the PEM-encoded peer client private key on the instance.`,
			},

			"deletion_policy": tpgresource.DeletionPolicySchemaEntry("DELETE"),

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSqlCustomerManagedCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CustomerManagedCertificate: %s", err)
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	instance := d.Get("instance").(string)
	bundle := d.Get("bundle").(string)

	configObj := expandSqlCustomerManagedCertificateConfig(d)

	body := map[string]interface{}{
		"allowMissing":                     true,
		"customerManagedCertificateConfig": configObj,
	}

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new CustomerManagedCertificate %q in instance %q", bundle, instance)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      body,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating CustomerManagedCertificate %q in instance %q: %s", bundle, instance, err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = SqlAdminOperationWaitTime(
		config, res, project, "Create Customer Managed Certificate", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error waiting to create CustomerManagedCertificate %q in instance %q: %s", bundle, instance, err)
	}

	log.Printf("[DEBUG] Finished creating CustomerManagedCertificate %q: %#v", d.Id(), res)
	return resourceSqlCustomerManagedCertificateRead(d, meta)
}

func resourceSqlCustomerManagedCertificateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CustomerManagedCertificate: %s", err)
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	instance := d.Get("instance").(string)
	bundle := d.Get("bundle").(string)

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
		return transport_tpg.HandleNotFoundError(transformSQLDatabaseReadError(err), d, fmt.Sprintf("SQL CustomerManagedCertificate %q in instance %q", bundle, instance))
	}

	foundBundle := findSqlCustomerManagedCertificateBundle(res, bundle)
	if foundBundle == nil {
		log.Printf("[WARN] Removing SQL CustomerManagedCertificate %q in instance %q because it is gone", bundle, instance)
		d.SetId("")
		return nil
	}

	if err := tpgresource.DeletionPolicyReadDefault(d, config, "DELETE"); err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("instance", instance); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}
	if v, ok := foundBundle["bundle"].(string); ok && v != "" {
		if err := d.Set("bundle", v); err != nil {
			return fmt.Errorf("Error setting bundle: %s", err)
		}
	}
	if err := d.Set("intent", foundBundle["intent"]); err != nil {
		return fmt.Errorf("Error setting intent: %s", err)
	}

	dns, ipv4 := extractPrimaryDnsFields(foundBundle["primaryDnsRecord"])
	if err := d.Set("primary_dns", dns); err != nil {
		return fmt.Errorf("Error setting primary_dns: %s", err)
	}
	if err := d.Set("primary_ipv4_address", ipv4); err != nil {
		return fmt.Errorf("Error setting primary_ipv4_address: %s", err)
	}

	filePaths, serverCaPath, clientCertPath, clientKeyPath := extractCustomerManagedCertFilePaths(foundBundle["customerManagedCertificates"])
	if err := d.Set("file_paths", filePaths); err != nil {
		return fmt.Errorf("Error setting file_paths: %s", err)
	}
	if err := d.Set("server_ca_cert_file_path", serverCaPath); err != nil {
		return fmt.Errorf("Error setting server_ca_cert_file_path: %s", err)
	}
	if err := d.Set("client_cert_file_path", clientCertPath); err != nil {
		return fmt.Errorf("Error setting client_cert_file_path: %s", err)
	}
	if err := d.Set("client_key_file_path", clientKeyPath); err != nil {
		return fmt.Errorf("Error setting client_key_file_path: %s", err)
	}

	return nil
}

func resourceSqlCustomerManagedCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CustomerManagedCertificate: %s", err)
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	instance := d.Get("instance").(string)
	bundle := d.Get("bundle").(string)

	configObj, updateMask := buildSqlCustomerManagedCertificateUpdatePayload(d)
	if len(updateMask) == 0 {
		return resourceSqlCustomerManagedCertificateRead(d, meta)
	}

	body := map[string]interface{}{
		"allowMissing":                     false,
		"updateMask":                       strings.Join(updateMask, ","),
		"customerManagedCertificateConfig": configObj,
	}

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating CustomerManagedCertificate %q in instance %q with mask %v", bundle, instance, updateMask)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      body,
		Timeout:   d.Timeout(schema.TimeoutUpdate),
	})
	if err != nil {
		return fmt.Errorf("Error updating CustomerManagedCertificate %q in instance %q: %s", bundle, instance, err)
	}

	err = SqlAdminOperationWaitTime(
		config, res, project, "Update Customer Managed Certificate", userAgent,
		d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("Error waiting to update CustomerManagedCertificate %q in instance %q: %s", bundle, instance, err)
	}

	log.Printf("[DEBUG] Finished updating CustomerManagedCertificate %q: %#v", d.Id(), res)
	return resourceSqlCustomerManagedCertificateRead(d, meta)
}

func resourceSqlCustomerManagedCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	if ok, err := tpgresource.DeletionPolicyPreDelete(d); err != nil {
		return err
	} else if ok {
		return nil
	}

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for CustomerManagedCertificate: %s", err)
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	instance := d.Get("instance").(string)
	bundle := d.Get("bundle").(string)

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))

	url, err := tpgresource.ReplaceVars(d, config, "{{SQLBasePath}}projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting CustomerManagedCertificate %q in instance %q", bundle, instance)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutDelete),
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SQL CustomerManagedCertificate %q in instance %q", bundle, instance))
	}

	err = SqlAdminOperationWaitTime(
		config, res, project, "Delete Customer Managed Certificate", userAgent,
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error waiting to delete CustomerManagedCertificate %q in instance %q: %s", bundle, instance, err)
	}

	log.Printf("[DEBUG] Finished deleting CustomerManagedCertificate %q: %#v", d.Id(), res)
	return nil
}

func resourceSqlCustomerManagedCertificateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/customerManagedCertificates/(?P<bundle>[^/]+)$",
		"^instances/(?P<instance>[^/]+)/customerManagedCertificates/(?P<bundle>[^/]+)$",
		"^(?P<project>[^/]+)/(?P<instance>[^/]+)/(?P<bundle>[^/]+)$",
		"^(?P<instance>[^/]+)/(?P<bundle>[^/]+)$",
	}, d, config); err != nil {
		return nil, err
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance}}/customerManagedCertificates/{{bundle}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func expandSqlCustomerManagedCertificateConfig(d tpgresource.TerraformResourceData) map[string]interface{} {
	configObj := map[string]interface{}{
		"bundle": d.Get("bundle").(string),
		"intent": d.Get("intent").(string),
	}

	primaryDns, _ := d.Get("primary_dns").(string)
	primaryIpv4, _ := d.Get("primary_ipv4_address").(string)
	if dnsRecord := buildPrimaryDnsRecordPayload(primaryDns, primaryIpv4); dnsRecord != nil {
		configObj["primaryDnsRecord"] = dnsRecord
	}

	serverCa, _ := d.Get("server_ca_cert").(string)
	clientCert, _ := d.Get("client_cert").(string)
	clientKey, _ := d.Get("client_key").(string)
	if certs := buildCustomerManagedCertificatesPayload(serverCa, clientCert, clientKey); len(certs) > 0 {
		configObj["customerManagedCertificates"] = certs
	}

	return configObj
}

func buildSqlCustomerManagedCertificateUpdatePayload(d *schema.ResourceData) (map[string]interface{}, []string) {
	configObj := map[string]interface{}{
		"bundle": d.Get("bundle").(string),
	}
	var updateMask []string

	if d.HasChange("intent") {
		configObj["intent"] = d.Get("intent").(string)
		updateMask = append(updateMask, "intent")
	}

	dnsChanged := d.HasChange("primary_dns")
	ipv4Changed := d.HasChange("primary_ipv4_address")
	if dnsChanged || ipv4Changed {
		primaryDns, _ := d.Get("primary_dns").(string)
		primaryIpv4, _ := d.Get("primary_ipv4_address").(string)
		if dnsRecord := buildPrimaryDnsRecordPayload(primaryDns, primaryIpv4); dnsRecord != nil {
			configObj["primaryDnsRecord"] = dnsRecord
		}
		if dnsChanged && !ipv4Changed {
			updateMask = append(updateMask, "primary_dns_record.dns")
		} else if ipv4Changed && !dnsChanged {
			updateMask = append(updateMask, "primary_dns_record.ipv4_address")
		} else {
			updateMask = append(updateMask, "primary_dns_record")
		}
	}

	if d.HasChange("server_ca_cert") || d.HasChange("client_cert") || d.HasChange("client_key") {
		serverCa, _ := d.Get("server_ca_cert").(string)
		clientCert, _ := d.Get("client_cert").(string)
		clientKey, _ := d.Get("client_key").(string)
		configObj["customerManagedCertificates"] = buildCustomerManagedCertificatesPayload(serverCa, clientCert, clientKey)
		updateMask = append(updateMask, "customer_managed_certificates")
	}

	return configObj, updateMask
}

func buildPrimaryDnsRecordPayload(dns, ipv4 string) map[string]interface{} {
	if dns == "" && ipv4 == "" {
		return nil
	}
	res := make(map[string]interface{})
	if dns != "" {
		res["dns"] = dns
	}
	if ipv4 != "" {
		res["ipv4Address"] = ipv4
	}
	return res
}

func buildCustomerManagedCertificatesPayload(serverCa, clientCert, clientKey string) []map[string]interface{} {
	var certs []map[string]interface{}
	if serverCa != "" {
		certs = append(certs, map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_SERVER_CA",
			"content":         encodeCustomerManagedCertificateContent(serverCa),
		})
	}
	if clientCert != "" {
		certs = append(certs, map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_CERTIFICATE",
			"content":         encodeCustomerManagedCertificateContent(clientCert),
		})
	}
	if clientKey != "" {
		certs = append(certs, map[string]interface{}{
			"certificateType": "CERTIFICATE_TYPE_PEER_CLIENT_KEY",
			"content":         encodeCustomerManagedCertificateContent(clientKey),
		})
	}
	return certs
}

func encodeCustomerManagedCertificateContent(val string) string {
	trimmed := strings.TrimSpace(val)
	if !strings.Contains(trimmed, "-----BEGIN") {
		if decoded, err := base64.StdEncoding.DecodeString(trimmed); err == nil && len(decoded) > 0 {
			return trimmed
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(val))
}

func findSqlCustomerManagedCertificateBundle(res map[string]interface{}, targetBundle string) map[string]interface{} {
	if res == nil {
		return nil
	}
	rawList, ok := res["customerManagedCertificates"].([]interface{})
	if !ok {
		return nil
	}
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if b, ok := m["bundle"].(string); ok && b == targetBundle {
			return m
		}
	}
	return nil
}

func extractPrimaryDnsFields(v interface{}) (string, string) {
	m, ok := v.(map[string]interface{})
	if !ok || m == nil {
		return "", ""
	}
	dns, _ := m["dns"].(string)
	ipv4, _ := m["ipv4Address"].(string)
	return dns, ipv4
}

func extractCustomerManagedCertFilePaths(v interface{}) ([]string, string, string, string) {
	rawList, ok := v.([]interface{})
	if !ok || len(rawList) == 0 {
		return nil, "", "", ""
	}
	filePaths := make([]string, 0, len(rawList))
	var serverCaPath, clientCertPath, clientKeyPath string
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		certType, _ := m["certificateType"].(string)
		filePath, _ := m["filePath"].(string)
		if filePath != "" {
			filePaths = append(filePaths, filePath)
		}
		switch certType {
		case "CERTIFICATE_TYPE_PEER_SERVER_CA":
			serverCaPath = filePath
		case "CERTIFICATE_TYPE_PEER_CLIENT_CERTIFICATE":
			clientCertPath = filePath
		case "CERTIFICATE_TYPE_PEER_CLIENT_KEY":
			clientKeyPath = filePath
		}
	}
	sort.Strings(filePaths)
	return filePaths, serverCaPath, clientCertPath, clientKeyPath
}

func init() {
	registry.Schema{
		Name:        "google_sql_customer_managed_certificate",
		ProductName: "sql",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceSqlCustomerManagedCertificate(),
	}.Register()
}
