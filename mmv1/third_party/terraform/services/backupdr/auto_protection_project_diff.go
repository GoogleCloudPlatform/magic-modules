package backupdr

import (
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanagerv3"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var (
	backupDRProjectIDByNumberMu sync.RWMutex
	backupDRProjectIDByNumber   = map[string]string{}
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

func recordBackupDRProjectIDAndNumber(projectID, projectNumber string) {
	if projectID == "" || projectNumber == "" {
		return
	}
	if isNumericProjectIdentifier(projectID) || !isNumericProjectIdentifier(projectNumber) {
		return
	}
	backupDRProjectIDByNumberMu.Lock()
	backupDRProjectIDByNumber[projectNumber] = projectID
	backupDRProjectIDByNumberMu.Unlock()
}

func lookupBackupDRProjectIDForNumber(projectNumber string) (string, bool) {
	backupDRProjectIDByNumberMu.RLock()
	projectID, ok := backupDRProjectIDByNumber[projectNumber]
	backupDRProjectIDByNumberMu.RUnlock()
	return projectID, ok
}

// backupDRAutoProtectionProjectDiffSuppress suppresses diffs on BackupDR auto
// protection fields (scope and backup_plan) only when the trailing resource
// path is identical and the numeric project number is verified to correspond to
// the project ID. Unlike tpgresource.ProjectNumberDiffSuppress, changing to a
// different project ID or a different numeric project number is never
// suppressed.
func backupDRAutoProtectionProjectDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}
	oldProj, oldRest, oldOK := parseBackupDRProjectPrefix(old)
	newProj, newRest, newOK := parseBackupDRProjectPrefix(new)
	if !oldOK || !newOK || oldRest != newRest {
		return false
	}

	oldIsNum := isNumericProjectIdentifier(oldProj)
	newIsNum := isNumericProjectIdentifier(newProj)
	if oldIsNum == newIsNum {
		// Both numeric or both project IDs, and already known to differ.
		return false
	}

	if oldIsNum {
		projectID, ok := lookupBackupDRProjectIDForNumber(oldProj)
		return ok && projectID == newProj
	}
	projectID, ok := lookupBackupDRProjectIDForNumber(newProj)
	return ok && projectID == oldProj
}

func resolveBackupDRProjectIDFromNumber(d *schema.ResourceData, config *transport_tpg.Config, projectNumber string) (string, error) {
	if config == nil {
		return "", nil
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return "", err
	}
	client := resourcemanagerv3.NewClient(config, userAgent)
	if client == nil {
		return "", nil
	}
	proj, err := client.Projects.Get("projects/" + projectNumber).Do()
	if err != nil || proj == nil {
		return "", err
	}
	return proj.ProjectId, nil
}

// recordBackupDRProjectNumber resolves projectNumber via Resource Manager v3
// and caches the verified projectNumber -> projectID mapping. If the lookup
// fails and the API response's own "name" echoes the same projectNumber for a
// request addressed by non-numeric resourceProject, that API-verified mapping
// is recorded as a fallback.
func recordBackupDRProjectNumber(
	d *schema.ResourceData,
	config *transport_tpg.Config,
	res map[string]interface{},
	projectNumber string,
	resolve func(projectNumber string) (string, error),
) {
	if !isNumericProjectIdentifier(projectNumber) {
		return
	}
	if _, known := lookupBackupDRProjectIDForNumber(projectNumber); known {
		return
	}

	if resolve != nil {
		projectID, err := resolve(projectNumber)
		if err == nil && projectID != "" {
			recordBackupDRProjectIDAndNumber(projectID, projectNumber)
			return
		}
		if err != nil {
			log.Printf("[WARN] Failed to resolve project number %s via Resource Manager v3: %s", projectNumber, err)
		}
	}

	// Fallback: if the API returned res["name"] under the same numeric project
	// number for a resource addressed using a non-numeric project ID, the
	// mapping is verified by the API response itself.
	if nameVal, ok := res["name"].(string); ok {
		if nameProj, _, nameOK := parseBackupDRProjectPrefix(nameVal); nameOK && nameProj == projectNumber {
			if resourceProj, err := tpgresource.GetProject(d, config); err == nil && resourceProj != "" && !isNumericProjectIdentifier(resourceProj) {
				recordBackupDRProjectIDAndNumber(resourceProj, projectNumber)
			}
		}
	}
}

func recordBackupDRAutoProtectionPolicyProjectNumbers(d *schema.ResourceData, meta interface{}, res map[string]interface{}) {
	config, _ := meta.(*transport_tpg.Config)
	recordBackupDRAutoProtectionPolicyProjectNumbersWithResolver(d, config, res, func(projectNumber string) (string, error) {
		return resolveBackupDRProjectIDFromNumber(d, config, projectNumber)
	})
}

func recordBackupDRAutoProtectionPolicyProjectNumbersWithResolver(
	d *schema.ResourceData,
	config *transport_tpg.Config,
	res map[string]interface{},
	resolve func(projectNumber string) (string, error),
) {
	detailsRaw, ok := res["backupPlanDetails"].([]interface{})
	if !ok {
		return
	}
	for _, item := range detailsRaw {
		detail, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		bp, ok := detail["backupPlan"].(string)
		if !ok {
			continue
		}
		apiProj, rest, parsed := parseBackupDRProjectPrefix(bp)
		if !parsed || rest == "" || !isNumericProjectIdentifier(apiProj) {
			continue
		}
		recordBackupDRProjectNumber(d, config, res, apiProj, resolve)
	}
}

func recordBackupDRAutoProtectionPolicyBindingScopeProjectNumber(d *schema.ResourceData, meta interface{}, res map[string]interface{}) {
	config, _ := meta.(*transport_tpg.Config)
	recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(d, config, res, func(projectNumber string) (string, error) {
		return resolveBackupDRProjectIDFromNumber(d, config, projectNumber)
	})
}

func recordBackupDRAutoProtectionPolicyBindingScopeProjectNumberWithResolver(
	d *schema.ResourceData,
	config *transport_tpg.Config,
	res map[string]interface{},
	resolve func(projectNumber string) (string, error),
) {
	scopeVal, ok := res["scope"].(string)
	if !ok {
		return
	}
	apiProj, rest, parsed := parseBackupDRProjectPrefix(scopeVal)
	if !parsed || rest != "" || !isNumericProjectIdentifier(apiProj) {
		return
	}
	recordBackupDRProjectNumber(d, config, res, apiProj, resolve)
}
