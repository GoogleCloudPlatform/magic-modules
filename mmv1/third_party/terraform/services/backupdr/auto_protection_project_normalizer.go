package backupdr

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanagerv3"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// parseBackupDRProjectPrefix splits a resource name of the form
// "projects/{project}" or "projects/{project}/{rest}" into the project
// identifier and the trailing suffix (empty string or "/...").
func parseBackupDRProjectPrefix(v string) (project string, rest string, ok bool) {
	if !strings.HasPrefix(v, "projects/") {
		return "", "", false
	}
	after := strings.TrimPrefix(v, "projects/")
	if after == "" {
		return "", "", false
	}
	if idx := strings.IndexByte(after, '/'); idx >= 0 {
		if idx == 0 {
			return "", "", false
		}
		return after[:idx], after[idx:], true
	}
	return after, "", true
}

func isNumericProjectIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func backupDRResourceManagerBillingProject(d *schema.ResourceData, config *transport_tpg.Config) string {
	billingProject := ""
	if project, err := tpgresource.GetProject(d, config); err == nil && project != "" {
		billingProject = project
	}
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil && bp != "" {
		billingProject = bp
	}
	return billingProject
}

// resolveBackupDRProjectIDFromNumber looks up a project number via Cloud
// Resource Manager v3 using transport_tpg.SendRequest so retries, custom
// endpoints, and user_project_override (X-Goog-User-Project) are handled
// consistently with other provider requests.
func resolveBackupDRProjectIDFromNumber(d *schema.ResourceData, config *transport_tpg.Config, projectNumber string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("provider config is nil")
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return "", err
	}

	billingProject := backupDRResourceManagerBillingProject(d, config)
	baseURL := strings.TrimRight(transport_tpg.BaseUrl(resourcemanagerv3.Product, config), "/")
	url := fmt.Sprintf("%s/projects/%s", baseURL, projectNumber)

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutRead),
	})
	if err != nil {
		return "", err
	}
	projectID, _ := res["projectId"].(string)
	if projectID == "" {
		return "", fmt.Errorf("Resource Manager v3 response for projects/%s did not include projectId", projectNumber)
	}
	return projectID, nil
}

func normalizeBackupDRAutoProtectionPolicyBindingScope(d *schema.ResourceData, meta interface{}, res map[string]interface{}) error {
	config, _ := meta.(*transport_tpg.Config)
	return normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(d, res, func(projectNumber string) (string, error) {
		return resolveBackupDRProjectIDFromNumber(d, config, projectNumber)
	})
}

func normalizeBackupDRAutoProtectionPolicyBindingScopeWithResolver(
	d *schema.ResourceData,
	res map[string]interface{},
	resolve func(projectNumber string) (string, error),
) error {
	scopeVal, ok := res["scope"].(string)
	if !ok {
		return nil
	}
	apiProj, apiRest, parsed := parseBackupDRProjectPrefix(scopeVal)
	if !parsed || apiRest != "" || !isNumericProjectIdentifier(apiProj) {
		return nil
	}

	cfgScope, _ := d.Get("scope").(string)
	if cfgScope == "" {
		// Import path (no prior scope in state): resolve to project ID on a
		// best-effort basis so imported state aligns with project-ID configs.
		if resolve == nil {
			return nil
		}
		projectID, err := resolve(apiProj)
		if err != nil {
			log.Printf("[WARN] Failed to resolve project number %s for imported BackupDR AutoProtectionPolicyBinding scope: %s", apiProj, err)
			return nil
		}
		if projectID != "" && !isNumericProjectIdentifier(projectID) {
			res["scope"] = "projects/" + projectID
		}
		return nil
	}

	if cfgScope == scopeVal {
		// Configured/state scope already uses the same numeric project number.
		return nil
	}

	cfgProj, cfgRest, cfgParsed := parseBackupDRProjectPrefix(cfgScope)
	if !cfgParsed || cfgRest != "" || isNumericProjectIdentifier(cfgProj) {
		// Different numeric project number or non-project scope format: keep API value.
		return nil
	}

	if resolve == nil {
		return fmt.Errorf("cannot resolve project number %q for BackupDR AutoProtectionPolicyBinding scope: resolver is nil", apiProj)
	}
	projectID, err := resolve(apiProj)
	if err != nil {
		return fmt.Errorf(
			"error resolving project number %q for BackupDR AutoProtectionPolicyBinding scope: %w "+
				"(ensure cloudresourcemanager.googleapis.com is enabled and caller has resourcemanager.projects.get, "+
				"or configure scope using the numeric project number \"projects/%s\")",
			apiProj, err, apiProj,
		)
	}
	if projectID == cfgProj {
		res["scope"] = cfgScope
	}
	return nil
}

func normalizeBackupDRAutoProtectionPolicyProjectNumbers(d *schema.ResourceData, meta interface{}, res map[string]interface{}) error {
	config, _ := meta.(*transport_tpg.Config)
	return normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, res, func(projectNumber string) (string, error) {
		return resolveBackupDRProjectIDFromNumber(d, config, projectNumber)
	})
}

func normalizeBackupDRAutoProtectionPolicyProjectNumbersWithResolver(
	d *schema.ResourceData,
	res map[string]interface{},
	resolve func(projectNumber string) (string, error),
) error {
	detailsRaw, ok := res["backupPlanDetails"].([]interface{})
	if !ok || len(detailsRaw) == 0 {
		return nil
	}

	var cfgDetails []interface{}
	if raw, ok := d.Get("backup_plan_details").([]interface{}); ok {
		cfgDetails = raw
	}

	cfgProjectsByRest := map[string][]string{}
	for _, item := range cfgDetails {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		bp, _ := m["backup_plan"].(string)
		cfgProj, cfgRest, parsed := parseBackupDRProjectPrefix(bp)
		if !parsed || cfgRest == "" {
			continue
		}
		cfgProjectsByRest[cfgRest] = append(cfgProjectsByRest[cfgRest], cfgProj)
	}

	// Memoize lookups only within this single Read call (never across calls).
	resolvedInRead := map[string]string{}
	resolveOnce := func(projectNumber string) (string, error) {
		if id, ok := resolvedInRead[projectNumber]; ok {
			return id, nil
		}
		if resolve == nil {
			return "", fmt.Errorf("resolver is nil")
		}
		id, err := resolve(projectNumber)
		if err != nil {
			return "", err
		}
		resolvedInRead[projectNumber] = id
		return id, nil
	}

	isImport := len(cfgDetails) == 0

	for _, item := range detailsRaw {
		detail, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		bp, ok := detail["backupPlan"].(string)
		if !ok {
			continue
		}
		apiProj, apiRest, parsed := parseBackupDRProjectPrefix(bp)
		if !parsed || apiRest == "" || !isNumericProjectIdentifier(apiProj) {
			continue
		}

		if isImport {
			projectID, err := resolveOnce(apiProj)
			if err != nil {
				log.Printf("[WARN] Failed to resolve project number %s for imported BackupDR AutoProtectionPolicy backup_plan: %s", apiProj, err)
				continue
			}
			if projectID != "" && !isNumericProjectIdentifier(projectID) {
				detail["backupPlan"] = "projects/" + projectID + apiRest
			}
			continue
		}

		candidates := cfgProjectsByRest[apiRest]
		if len(candidates) == 0 {
			continue
		}

		exactNumericMatch := false
		hasNonNumericCandidate := false
		for _, c := range candidates {
			if c == apiProj {
				exactNumericMatch = true
				break
			}
			if !isNumericProjectIdentifier(c) {
				hasNonNumericCandidate = true
			}
		}
		if exactNumericMatch || !hasNonNumericCandidate {
			continue
		}

		projectID, err := resolveOnce(apiProj)
		if err != nil {
			return fmt.Errorf(
				"error resolving project number %q for BackupDR AutoProtectionPolicy backup_plan: %w "+
					"(ensure cloudresourcemanager.googleapis.com is enabled and caller has resourcemanager.projects.get, "+
					"or configure backup_plan using the numeric project number \"projects/%s%s\")",
				apiProj, err, apiProj, apiRest,
			)
		}
		for _, c := range candidates {
			if c == projectID {
				detail["backupPlan"] = "projects/" + projectID + apiRest
				break
			}
		}
	}

	return nil
}
