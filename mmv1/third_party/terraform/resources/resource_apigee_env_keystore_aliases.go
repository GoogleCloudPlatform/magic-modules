package google

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApigeeEnvKeystoreAliases() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigeeEnvKeystoreAliasesCreate,
		Read:        resourceApigeeEnvKeystoreAliasesRead,
		Update:      resourceApigeeEnvKeystoreAliasesUpdate,
		Delete:      resourceApigeeEnvKeystoreAliasesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceApigeeEnvKeystoreAliasesImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:        schema.TypeString,
				Description: `The HTTP Content-Type header value specifying the content type of the body.`,
			},
			"format": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Required. Format of the data. Valid values include: selfsignedcert, keycertfile, or pkcs12`,
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The Apigee Organization associated with the environment`,
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource ID of the environment.`,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `Alias for the key/certificate pair. Values must match the regular expression [\w\s-.]{1,255}. 
This must be provided for all formats except selfsignedcert; self-signed certs may specify the alias in either 
this parameter or the JSON body.`,
			},
			"cert_validity_in_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Signature algorithm to generate private key. Valid values are SHA512withRSA, SHA384withRSA, and SHA256withRSA`,
			},
			"key_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Key size. Default and maximum value is 2048 bits.`,
			},
			"sig_alg": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `sigAlgName.`,
			},
			"key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Key File.`,
			}
			"cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Cert File.`,
			}
			"file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `File.`,
			}
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive: true,
				Description: `Password for the file.`,
			}
			"subject": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Chain of certificates under this alias.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `X.509 version.`,
						},
						"country_code": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Two-letter country code. Example, IN for India, US for United States of America.`,
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Email address.`,
						},
						"locality": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `City or town name. Maximum length is 128 characters.`,
						},
						"org": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Organization name.`,
						},
						"org_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Organization team name. Maximum length is 64 characters.`,
						},
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `State.`,
						},
					},
				},
			},
			"subject_alternative_dns_names": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `List of alternative host names. Maximum length is 255 characters for each value.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subject_alternative_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `subjectAlternativeName`,
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Optional.Type of Alias`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceApigeeEnvKeystoreAliasesUpdate(d *schema.ResourceData, meta interface{}) error {

// GET https://apigee.googleapis.com/v1/{parent=organizations/*/environments/*/keystores/*}/aliases

// Name of the keystore. Use the following format in your request: organizations/{org}/environments/{env}/keystores/{keystore}.




} 
func resourceApigeeEnvKeystoreAliasesCreate(d *schema.ResourceData, meta interface{}) error {

// keycertfile -- Content-Type: multipart/form-data, keyFile, certFile, password
// pkcs12 -- Content-Type: multipart/form-data, file, password
// selfsignedcert - Content-Type: application/json, CertificateGenerationSpec

// POST https://apigee.googleapis.com/v1/{parent=organizations/*/environments/*/keystores/*}/aliases

	if format == "pkcs12" {

	} else if format == "keycertfile" {

	} else if format == "selfsignedcert" {

	} else {
		fmt.
	}


	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	aliasProp, err := expandApigeeEnvKeystoreAliasesAlias(d.Get("alias"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("alias"); !isEmptyValue(reflect.ValueOf(aliasProp)) && (ok || !reflect.DeepEqual(v, aliasProp)) {
		obj["alias"] = aliasProp
	}
	subjectProp, err := expandApigeeEnvKeystoreAliasesSubject(d.Get("subject"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("subject"); !isEmptyValue(reflect.ValueOf(subjectProp)) && (ok || !reflect.DeepEqual(v, subjectProp)) {
		obj["subject"] = subjectProp
	}
	subjectAlternativeDNSNamesProp, err := expandApigeeEnvKeystoreAliasesSubjectAlternativeDNSNames(d.Get("subject_alternative_dns_names"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("subject_alternative_dns_names"); !isEmptyValue(reflect.ValueOf(subjectAlternativeDNSNamesProp)) && (ok || !reflect.DeepEqual(v, subjectAlternativeDNSNamesProp)) {
		obj["subjectAlternativeDNSNames"] = subjectAlternativeDNSNamesProp
	}
	keySizeProp, err := expandApigeeEnvKeystoreAliasesKeySize(d.Get("key_size"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("key_size"); !isEmptyValue(reflect.ValueOf(keySizeProp)) && (ok || !reflect.DeepEqual(v, keySizeProp)) {
		obj["keySize"] = keySizeProp
	}
	sigAlgProp, err := expandApigeeEnvKeystoreAliasesSigAlg(d.Get("sig_alg"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("sig_alg"); !isEmptyValue(reflect.ValueOf(sigAlgProp)) && (ok || !reflect.DeepEqual(v, sigAlgProp)) {
		obj["sigAlg"] = sigAlgProp
	}
	certValidityInDaysProp, err := expandApigeeEnvKeystoreAliasesCertValidityInDays(d.Get("cert_validity_in_days"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("cert_validity_in_days"); !isEmptyValue(reflect.ValueOf(certValidityInDaysProp)) && (ok || !reflect.DeepEqual(v, certValidityInDaysProp)) {
		obj["certValidityInDays"] = certValidityInDaysProp
	}

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}/organizations/{{org_id}}/environments/{{environment}}/aliases?alias={{alias}}&format={{format}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new EnvKeystoreAliases: %#v", obj)
	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating EnvKeystoreAliases: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/aliases/{{alias}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating EnvKeystoreAliases %q: %#v", d.Id(), res)

	return resourceApigeeEnvKeystoreAliasesRead(d, meta)
}

func resourceApigeeEnvKeystoreAliasesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}/organizations/{{org_id}}/environments/{{environment}}/aliases/{{alias}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("ApigeeEnvKeystoreAliases %q", d.Id()))
	}

	if err := d.Set("type", flattenApigeeEnvKeystoreAliasesType(res["type"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("alias", flattenApigeeEnvKeystoreAliasesAlias(res["alias"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("subject", flattenApigeeEnvKeystoreAliasesSubject(res["subject"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("subject_alternative_dns_names", flattenApigeeEnvKeystoreAliasesSubjectAlternativeDNSNames(res["subjectAlternativeDNSNames"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("key_size", flattenApigeeEnvKeystoreAliasesKeySize(res["keySize"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("sig_alg", flattenApigeeEnvKeystoreAliasesSigAlg(res["sigAlg"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}
	if err := d.Set("cert_validity_in_days", flattenApigeeEnvKeystoreAliasesCertValidityInDays(res["certValidityInDays"], d, config)); err != nil {
		return fmt.Errorf("Error reading EnvKeystoreAliases: %s", err)
	}

	return nil
}

func resourceApigeeEnvKeystoreAliasesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	url, err := replaceVars(d, config, "{{ApigeeBasePath}}/organizations/{{org_id}}/environments/{{environment}}/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting EnvKeystoreAliases %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "EnvKeystoreAliases")
	}

	log.Printf("[DEBUG] Finished deleting EnvKeystoreAliases %q: %#v", d.Id(), res)
	return nil
}

// func resourceApigeeEnvKeystoreAliasesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
// 	config := meta.(*Config)

// 	// current import_formats cannot import fields with forward slashes in their value
// 	if err := parseImportId([]string{
// 		"(?P<env_id>.+)/(?P<environment>.+)/(?P<name>.+)",
// 		"(?P<env_id>.+)/(?P<name>.+)",
// 	}, d, config); err != nil {
// 		return nil, err
// 	}

// 	// Replace import id for the resource id
// 	id, err := replaceVars(d, config, "organizations/{{org_id}}/environments/{{environment}}/{{name}}")
// 	if err != nil {
// 		return nil, fmt.Errorf("Error constructing id: %s", err)
// 	}
// 	d.SetId(id)

// 	return []*schema.ResourceData{d}, nil
// }

func flattenApigeeEnvKeystoreAliasesType(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesAlias(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubject(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["common_name"] =
		flattenApigeeEnvKeystoreAliasesSubjectCommonName(original["commonName"], d, config)
	transformed["org"] =
		flattenApigeeEnvKeystoreAliasesSubjectOrg(original["org"], d, config)
	transformed["org_unit"] =
		flattenApigeeEnvKeystoreAliasesSubjectOrgUnit(original["orgUnit"], d, config)
	transformed["state"] =
		flattenApigeeEnvKeystoreAliasesSubjectState(original["state"], d, config)
	transformed["country_code"] =
		flattenApigeeEnvKeystoreAliasesSubjectCountryCode(original["countryCode"], d, config)
	transformed["locality"] =
		flattenApigeeEnvKeystoreAliasesSubjectLocality(original["locality"], d, config)
	transformed["email"] =
		flattenApigeeEnvKeystoreAliasesSubjectEmail(original["email"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeEnvKeystoreAliasesSubjectCommonName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectOrg(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectOrgUnit(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectCountryCode(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectLocality(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectEmail(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesSubjectAlternativeDNSNames(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["subject_alternative_name"] =
		flattenApigeeEnvKeystoreAliasesSubjectAlternativeDNSNamesSubjectAlternativeName(original["subjectAlternativeName"], d, config)
	return []interface{}{transformed}
}
func flattenApigeeEnvKeystoreAliasesSubjectAlternativeDNSNamesSubjectAlternativeName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesKeySize(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenApigeeEnvKeystoreAliasesSigAlg(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenApigeeEnvKeystoreAliasesCertValidityInDays(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := stringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func expandApigeeEnvKeystoreAliasesAlias(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubject(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCommonName, err := expandApigeeEnvKeystoreAliasesSubjectCommonName(original["common_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCommonName); val.IsValid() && !isEmptyValue(val) {
		transformed["commonName"] = transformedCommonName
	}

	transformedOrg, err := expandApigeeEnvKeystoreAliasesSubjectOrg(original["org"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOrg); val.IsValid() && !isEmptyValue(val) {
		transformed["org"] = transformedOrg
	}

	transformedOrgUnit, err := expandApigeeEnvKeystoreAliasesSubjectOrgUnit(original["org_unit"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOrgUnit); val.IsValid() && !isEmptyValue(val) {
		transformed["orgUnit"] = transformedOrgUnit
	}

	transformedState, err := expandApigeeEnvKeystoreAliasesSubjectState(original["state"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedState); val.IsValid() && !isEmptyValue(val) {
		transformed["state"] = transformedState
	}

	transformedCountryCode, err := expandApigeeEnvKeystoreAliasesSubjectCountryCode(original["country_code"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCountryCode); val.IsValid() && !isEmptyValue(val) {
		transformed["countryCode"] = transformedCountryCode
	}

	transformedLocality, err := expandApigeeEnvKeystoreAliasesSubjectLocality(original["locality"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLocality); val.IsValid() && !isEmptyValue(val) {
		transformed["locality"] = transformedLocality
	}

	transformedEmail, err := expandApigeeEnvKeystoreAliasesSubjectEmail(original["email"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEmail); val.IsValid() && !isEmptyValue(val) {
		transformed["email"] = transformedEmail
	}

	return transformed, nil
}

func expandApigeeEnvKeystoreAliasesSubjectCommonName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectOrg(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectOrgUnit(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectState(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectCountryCode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectLocality(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectEmail(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSubjectAlternativeDNSNames(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSubjectAlternativeName, err := expandApigeeEnvKeystoreAliasesSubjectAlternativeDNSNamesSubjectAlternativeName(original["subject_alternative_name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSubjectAlternativeName); val.IsValid() && !isEmptyValue(val) {
		transformed["subjectAlternativeName"] = transformedSubjectAlternativeName
	}

	return transformed, nil
}

func expandApigeeEnvKeystoreAliasesSubjectAlternativeDNSNamesSubjectAlternativeName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesKeySize(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesSigAlg(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandApigeeEnvKeystoreAliasesCertValidityInDays(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
