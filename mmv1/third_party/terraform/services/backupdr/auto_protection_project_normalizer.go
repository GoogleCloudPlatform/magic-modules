package backupdr

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
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

func backupDRResourceManagerBillingProjectFromDiff(diff *schema.ResourceDiff, config *transport_tpg.Config) string {
	billingProject := ""
	if project, err := tpgresource.GetProjectFromDiff(diff, config); err == nil && project != "" {
		billingProject = project
	}
	if config != nil && config.BillingProject != "" {
		billingProject = config.BillingProject
	}
	return billingProject
}

func sendBackupDRProjectGetRequest(
	config *transport_tpg.Config,
	product registry.Product,
	billingProject string,
	userAgent string,
	timeout time.Duration,
	projectPathSegment string,
) (map[string]interface{}, error) {
	if config == nil {
		return nil, fmt.Errorf("provider config is nil")
	}
	baseURL := strings.TrimRight(transport_tpg.BaseUrl(product, config), "/")
	url := fmt.Sprintf("%s/projects/%s", baseURL, projectPathSegment)

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   timeout,
	})
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
	res, err := sendBackupDRProjectGetRequest(config, resourcemanagerv3.Product, billingProject, userAgent, d.Timeout(schema.TimeoutRead), projectNumber)
	if err != nil {
		return "", err
	}
	projectID, _ := res["projectId"].(string)
	if projectID == "" {
		return "", fmt.Errorf("Resource Manager v3 response for projects/%s did not include projectId", projectNumber)
	}
	return projectID, nil
}

func resolveBackupDRProjectIDFromNumberForDiff(diff *schema.ResourceDiff, config *transport_tpg.Config, projectNumber string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("provider config is nil")
	}
	billingProject := backupDRResourceManagerBillingProjectFromDiff(diff, config)
	res, err := sendBackupDRProjectGetRequest(config, resourcemanagerv3.Product, billingProject, config.UserAgent, 0, projectNumber)
	if err != nil {
		return "", err
	}
	projectID, _ := res["projectId"].(string)
	if projectID == "" {
		return "", fmt.Errorf("Resource Manager v3 response for projects/%s did not include projectId", projectNumber)
	}
	return projectID, nil
}

func resolveBackupDRProjectNumberFromIDForDiff(diff *schema.ResourceDiff, config *transport_tpg.Config, projectID string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("provider config is nil")
	}
	billingProject := backupDRResourceManagerBillingProjectFromDiff(diff, config)
	res, err := sendBackupDRProjectGetRequest(config, resourcemanager.Product, billingProject, config.UserAgent, 0, projectID)
	if err != nil {
		return "", err
	}
	projectNumber, _ := res["projectNumber"].(string)
	if projectNumber == "" || !isNumericProjectIdentifier(projectNumber) {
		return "", fmt.Errorf("Resource Manager v1 response for projects/%s did not include numeric projectNumber", projectID)
	}
	return projectNumber, nil
}

// backupDRAutoProtectionPolicyBindingScopeForceNewIfChanged enforces
// immutability of AutoProtectionPolicyBinding.scope while allowing a safe
// in-place state-only update when transitioning between `projects/{number}`
// and `projects/{project_id}` for the same project (such as upgrading state
// written by earlier provider versions).
func backupDRAutoProtectionPolicyBindingScopeForceNewIfChanged(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config, _ := meta.(*transport_tpg.Config)
	return backupDRAutoProtectionPolicyBindingScopeForceNewIfChangedWithResolvers(
		diff,
		func(projectNumber string) (string, error) {
			return resolveBackupDRProjectIDFromNumberForDiff(diff, config, projectNumber)
		},
		func(projectID string) (string, error) {
			return resolveBackupDRProjectNumberFromIDForDiff(diff, config, projectID)
		},
	)
}

func backupDRAutoProtectionPolicyBindingScopeForceNewIfChangedWithResolvers(
	diff *schema.ResourceDiff,
	resolveIDFromNumber func(projectNumber string) (string, error),
	resolveNumberFromID func(projectID string) (string, error),
) error {
	if diff.Id() == "" || !diff.HasChange("scope") {
		return nil
	}
	if !diff.NewValueKnown("scope") {
		return diff.ForceNew("scope")
	}

	oldVal, newVal := diff.GetChange("scope")
	oldScope, _ := oldVal.(string)
	newScope, _ := newVal.(string)
	if oldScope == newScope {
		return nil
	}

	oldProj, oldRest, oldOK := parseBackupDRProjectPrefix(oldScope)
	newProj, newRest, newOK := parseBackupDRProjectPrefix(newScope)
	if !oldOK || !newOK || oldRest != "" || newRest != "" {
		return diff.ForceNew("scope")
	}

	oldIsNum := isNumericProjectIdentifier(oldProj)
	newIsNum := isNumericProjectIdentifier(newProj)
	if oldIsNum == newIsNum {
		// Both numeric (different project numbers) or both project IDs
		// (different project IDs): force replacement without a CRM call.
		return diff.ForceNew("scope")
	}

	// Always resolve the existing project in state (oldProj) rather than the
	// incoming config value (newProj), since newProj may reference a different
	// or inaccessible project.
	if oldIsNum {
		if resolveIDFromNumber == nil {
			return fmt.Errorf("cannot resolve project number %q while diffing BackupDR AutoProtectionPolicyBinding scope: resolver is nil", oldProj)
		}
		resolvedID, err := resolveIDFromNumber(oldProj)
		if err != nil {
			return fmt.Errorf(
				"error resolving project number %q while diffing BackupDR AutoProtectionPolicyBinding scope: %w "+
					"(ensure cloudresourcemanager.googleapis.com is enabled and caller has resourcemanager.projects.get, "+
					"or keep scope as \"projects/%s\")",
				oldProj, err, oldProj,
			)
		}
		if resolvedID == newProj {
			return nil
		}
		return diff.ForceNew("scope")
	}

	if resolveNumberFromID == nil {
		return fmt.Errorf("cannot resolve project ID %q while diffing BackupDR AutoProtectionPolicyBinding scope: resolver is nil", oldProj)
	}
	resolvedNumber, err := resolveNumberFromID(oldProj)
	if err != nil {
		return fmt.Errorf(
			"error resolving project ID %q while diffing BackupDR AutoProtectionPolicyBinding scope: %w "+
				"(ensure cloudresourcemanager.googleapis.com is enabled and caller has resourcemanager.projects.get, "+
				"or keep scope as \"projects/%s\")",
			oldProj, err, oldProj,
		)
	}
	if resolvedNumber == newProj {
		return nil
	}
	return diff.ForceNew("scope")
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
