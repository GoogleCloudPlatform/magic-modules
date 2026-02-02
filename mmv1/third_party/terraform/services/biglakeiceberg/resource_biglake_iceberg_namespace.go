// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package biglakeiceberg

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var icebergNamespaceIgnoredProperties = []string{
	"location",
}

func isIgnoredProperty(k string) bool {
	for _, p := range icebergNamespaceIgnoredProperties {
		if k == p {
			return true
		}
	}
	return false
}

func ResourceBiglakeIcebergIcebergNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceBiglakeIcebergIcebergNamespaceCreate,
		Read:   resourceBiglakeIcebergIcebergNamespaceRead,
		Update: resourceBiglakeIcebergIcebergNamespaceUpdate,
		Delete: resourceBiglakeIcebergIcebergNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBiglakeIcebergIcebergNamespaceImport,
		},

		Schema: map[string]*schema.Schema{
			"catalog": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: icebergNamespacePropertiesDiffSuppress,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Hidden field to facilitate import parsing
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func icebergNamespacePropertiesDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// properties.KEY
	parts := strings.Split(k, ".")
	if len(parts) == 2 && isIgnoredProperty(parts[1]) {
		return true
	}
	return false
}

func encodeNamespace(ns []string) string {
	return url.PathEscape(strings.Join(ns, "\x1f"))
}

func decodeNamespace(nsStr string) ([]string, error) {
	decoded, err := url.PathUnescape(nsStr)
	if err != nil {
		return nil, err
	}
	return strings.Split(decoded, "\x1f"), nil
}

func resourceBiglakeIcebergIcebergNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	catalog := d.Get("catalog").(string)

	url, err := tpgresource.ReplaceVars(d, config, "{{BiglakeIcebergBasePath}}iceberg/v1/restcatalog/v1/projects/{{project}}/catalogs/{{catalog}}/namespaces")
	if err != nil {
		return err
	}

	nsRaw := d.Get("namespace").([]interface{})
	ns := make([]string, len(nsRaw))
	for i, v := range nsRaw {
		ns[i] = v.(string)
	}

	body := map[string]interface{}{
		"namespace": ns,
	}
	if v, ok := d.GetOk("properties"); ok {
		body["properties"] = v
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		RawURL:    url,
		UserAgent: config.UserAgent,
		Body:      body,
	})
	if err != nil {
		return fmt.Errorf("Error creating IcebergNamespace: %s", err)
	}

	id := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s", project, catalog, encodeNamespace(ns))
	d.SetId(id)

	return resourceBiglakeIcebergIcebergNamespaceRead(d, meta)
}

func resourceBiglakeIcebergIcebergNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	catalog := d.Get("catalog").(string)
	nsRaw := d.Get("namespace").([]interface{})
	ns := make([]string, len(nsRaw))
	for i, v := range nsRaw {
		ns[i] = v.(string)
	}

	encodedNs := encodeNamespace(ns)
	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{BiglakeIcebergBasePath}}iceberg/v1/restcatalog/v1/projects/{{project}}/catalogs/{{catalog}}/namespaces/%s", encodedNs))
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("IcebergNamespace %q", d.Id()))
	}

	if err := d.Set("namespace", res["namespace"]); err != nil {
		return fmt.Errorf("Error setting namespace: %s", err)
	}
	if err := d.Set("properties", res["properties"]); err != nil {
		return fmt.Errorf("Error setting properties: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("catalog", catalog); err != nil {
		return fmt.Errorf("Error setting catalog: %s", err)
	}
	if err := d.Set("namespace_id", encodedNs); err != nil {
		return fmt.Errorf("Error setting namespace_id: %s", err)
	}

	return nil
}

func resourceBiglakeIcebergIcebergNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	nsRaw := d.Get("namespace").([]interface{})
	ns := make([]string, len(nsRaw))
	for i, v := range nsRaw {
		ns[i] = v.(string)
	}

	if d.HasChange("properties") {
		oldProp, newProp := d.GetChange("properties")
		oldMap := oldProp.(map[string]interface{})
		newMap := newProp.(map[string]interface{})

		removals := []string{}
		for k := range oldMap {
			if isIgnoredProperty(k) {
				continue
			}
			if _, ok := newMap[k]; !ok {
				removals = append(removals, k)
			}
		}

		updates := map[string]string{}
		for k, v := range newMap {
			if isIgnoredProperty(k) {
				continue
			}
			updates[k] = v.(string)
		}

		body := map[string]interface{}{
			"removals": removals,
			"updates":  updates,
		}

		url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{BiglakeIcebergBasePath}}iceberg/v1/restcatalog/v1/projects/{{project}}/catalogs/{{catalog}}/namespaces/%s/properties", encodeNamespace(ns)))
		if err != nil {
			return err
		}

		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      body,
		})
		if err != nil {
			return fmt.Errorf("Error updating IcebergNamespace properties: %s", err)
		}
	}

	return resourceBiglakeIcebergIcebergNamespaceRead(d, meta)
}

func resourceBiglakeIcebergIcebergNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	nsRaw := d.Get("namespace").([]interface{})
	ns := make([]string, len(nsRaw))
	for i, v := range nsRaw {
		ns[i] = v.(string)
	}

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{BiglakeIcebergBasePath}}iceberg/v1/restcatalog/v1/projects/{{project}}/catalogs/{{catalog}}/namespaces/%s", encodeNamespace(ns)))
	if err != nil {
		return err
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		return fmt.Errorf("Error deleting IcebergNamespace: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceBiglakeIcebergIcebergNamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)

	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/catalogs/(?P<catalog>[^/]+)/namespaces/(?P<namespace_id>.+)",
		"(?P<project>[^/]+)/(?P<catalog>[^/]+)/namespaces/(?P<namespace_id>.+)",
		"(?P<catalog>[^/]+)/namespaces/(?P<namespace_id>.+)",
	}, d, config); err != nil {
		return nil, err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}
	d.Set("project", project)

	nsStr := d.Get("namespace_id").(string)
	ns, err := decodeNamespace(nsStr)
	if err != nil {
		return nil, fmt.Errorf("Error decoding namespace from ID: %s", err)
	}

	nsI := make([]interface{}, len(ns))
	for i, v := range ns {
		nsI[i] = v
	}

	if err := d.Set("namespace", nsI); err != nil {
		return nil, fmt.Errorf("Error setting namespace in import: %s", err)
	}

	catalog := d.Get("catalog").(string)
	id := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s", project, catalog, nsStr)
	d.SetId(id)

	if err := resourceBiglakeIcebergIcebergNamespaceRead(d, meta); err != nil {
		return nil, fmt.Errorf("Error calling Read during import: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}