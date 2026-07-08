package datalineage

import (
	"context"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func resourceDataLineageOpenLineageJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)

	obj := make(map[string]interface{})
	namespaceProp, err := expandDataLineageOpenLineageJobNamespace(d.Get("namespace"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("namespace"); !tpgresource.IsEmptyValue(reflect.ValueOf(namespaceProp)) && (ok || !reflect.DeepEqual(v, namespaceProp)) {
		obj["namespace"] = namespaceProp
	}
	nameProp, err := expandDataLineageOpenLineageJobName(d.Get("name"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandDataLineageOpenLineageJobDescription(d.Get("description"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	ownerProp, err := expandDataLineageOpenLineageJobOwner(d.Get("owner"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("owner"); !tpgresource.IsEmptyValue(reflect.ValueOf(ownerProp)) && (ok || !reflect.DeepEqual(v, ownerProp)) {
		obj["owner"] = ownerProp
	}
	inputProp, err := expandDataLineageOpenLineageJobInput(d.Get("input"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("input"); !tpgresource.IsEmptyValue(reflect.ValueOf(inputProp)) && (ok || !reflect.DeepEqual(v, inputProp)) {
		obj["input"] = inputProp
	}
	outputProp, err := expandDataLineageOpenLineageJobOutput(d.Get("output"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("output"); !tpgresource.IsEmptyValue(reflect.ValueOf(outputProp)) && (ok || !reflect.DeepEqual(v, outputProp)) {
		obj["output"] = outputProp
	}

	event := buildRunEvent(obj)

	log.Printf("[DEBUG] Creating new OpenLineageJob: %#v", obj)

	res, diagnostics := emitEvent(ctx, event, config)
	if diagnostics != nil {
		return diagnostics
	}

	process := res.GetProcess()

	err = d.Set("knowledge_catalog", flattenKnowledgeCatalog(process, res.GetRun()))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(process)

	log.Printf("[DEBUG] Finished creating OpenLineageJob %q: %#v", d.Id(), res)

	return nil
}

func resourceDataLineageOpenLineageJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	run, diagnostics := getLatestRunForProcess(ctx, config, d.Id())
	if diagnostics != nil {
		return diagnostics
	}

	if v, ok := d.GetOk("knowledge_catalog"); ok {
		r := v.([]interface{})[0].(map[string]interface{})["run"].(string)
		if run != r {
			log.Printf("[WARN] Run ID has changed for OpenLineageJob %q: %s -> %s, this suggests external modifications. It will get updated during next apply", d.Id(), r, run)
		}
	}

	return nil
}

func resourceDataLineageOpenLineageJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clientSideFields := map[string]bool{"deletion_policy": true}
	clientSideOnly := true
	for field := range ResourceDataLineageOpenLineageJob().Schema {
		if d.HasChange(field) && !clientSideFields[field] {
			clientSideOnly = false
			break
		}
	}
	if clientSideOnly {
		log.Print("[DEBUG] Only client-side changes detected. Cancelling update operation.")
		return resourceDataLineageOpenLineageJobRead(ctx, d, meta)
	}

	config := meta.(*transport_tpg.Config)

	obj := make(map[string]interface{})
	namespaceProp, err := expandDataLineageOpenLineageJobNamespace(d.Get("namespace"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("namespace"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, namespaceProp)) {
		obj["namespace"] = namespaceProp
	}
	nameProp, err := expandDataLineageOpenLineageJobName(d.Get("name"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	descriptionProp, err := expandDataLineageOpenLineageJobDescription(d.Get("description"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("description"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	ownerProp, err := expandDataLineageOpenLineageJobOwner(d.Get("owner"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("owner"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, ownerProp)) {
		obj["owner"] = ownerProp
	}
	inputProp, err := expandDataLineageOpenLineageJobInput(d.Get("input"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("input"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, inputProp)) {
		obj["input"] = inputProp
	}
	outputProp, err := expandDataLineageOpenLineageJobOutput(d.Get("output"), d, config)
	if err != nil {
		return diag.FromErr(err)
	} else if v, ok := d.GetOkExists("output"); !tpgresource.IsEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, outputProp)) {
		obj["output"] = outputProp
	}

	log.Printf("[DEBUG] Updating OpenLineageJob %q: %#v", d.Id(), obj)

	event := buildRunEvent(obj)

	response, diagnostics := emitEvent(ctx, event, config)
	if diagnostics != nil {
		return diagnostics
	}

	err = d.Set("knowledge_catalog", flattenKnowledgeCatalog(response.GetProcess(), response.GetRun()))

	if err != nil {
		return diag.Errorf("Error updating OpenLineageJob %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating OpenLineageJob %q: %#v", d.Id(), response)
	}

	return nil
}

func resourceDataLineageOpenLineageJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("deletion_policy").(string) == "PREVENT" {
		return diag.Errorf("cannot destroy DataLineageOpenLineageJob without setting deletion_policy=\"DELETE\" and running `terraform apply`")
	}
	if d.Get("deletion_policy").(string) == "ABANDON" {
		log.Printf("[DEBUG] deletion_policy set to \"ABANDON\", removing OpenLineageJob %q from Terraform state without deletion", d.Id())
		return nil
	}
	config := meta.(*transport_tpg.Config)

	err := deleteProcess(ctx, config, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Finished deleting OpenLineageJob %q", d.Id())
	return nil
}
