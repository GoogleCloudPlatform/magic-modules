package auditmanager

import (
	"context"
	"fmt"

	"google3/third_party/golang/github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"google3/third_party/golang/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google3/third_party/golang/terraform_providers/google/google/tpgresource"
	"google3/third_party/golang/terraform_providers/google/google/transport_tpg"
)

func dataSourceAuditManagerFrameworkAuditScopeReport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuditManagerFrameworkAuditScopeReportRead,
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The scope for the audit scope report.`,
			},
			"report_format": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The format that the scope report is returned in.`,
			},
			"compliance_framework": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The compliance framework that the scope report is generated for.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the audit report.`,
			},
			"scope_report_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The audit scope report content in byte format.`,
			},
		},
	}
}

func dataSourceAuditManagerFrameworkAuditScopeReportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)

	scope := d.Get("scope").(string)
	reportFormat := d.Get("report_format").(string)
	complianceFramework := d.Get("compliance_framework").(string)

	url, err := tpgresource.ReplaceVars(d, config, "{{AuditManagerBasePath}}/"+scope+"/frameworkAuditScopeReports:generateFrameworkAuditScopeReport")
	if err != nil {
		return diag.FromErr(err)
	}

	body := map[string]interface{}{
		"report_format":        reportFormat,
		"compliance_framework": complianceFramework,
	}

	res, err := transport_tpg.SendRequest(transport_tpg.TemplatedRequest{
		Config:    config,
		Method:    "POST",
		RawURL:    url,
		Body:      body,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error generating framework audit scope report: %s", err))
	}

	if err := d.Set("name", res["name"]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err := d.Set("compliance_framework", res["compliance_framework"]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting compliance_framework: %s", err))
	}
	if err := d.Set("scope_report_contents", res["scope_report_contents"]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting scope_report_contents: %s", err))
	}

	d.SetId(res["name"].(string))

	return nil
}
