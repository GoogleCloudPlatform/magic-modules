package saasruntime

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type saasRuntimeSweeperTiming struct {
	pollInterval       time.Duration
	deleteRetryTimeout time.Duration
	absenceTimeout     time.Duration
	operationTimeout   time.Duration
}

var defaultSaasRuntimeSweeperTiming = saasRuntimeSweeperTiming{
	pollInterval:       2 * time.Second,
	deleteRetryTimeout: 15 * time.Second,
	absenceTimeout:     60 * time.Second,
	operationTimeout:   2 * time.Minute,
}

var saasRuntimeFallbackLocations = []string{
	"global",
	"us-central1",
	"us-east1",
	"europe-west1",
}

var saasRuntimeTerminalUnitOpStates = map[string]bool{
	"UNIT_OPERATION_STATE_SUCCEEDED": true,
	"UNIT_OPERATION_STATE_FAILED":    true,
	"UNIT_OPERATION_STATE_CANCELLED": true,
}

func init() {
	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_unit_operation",
		ListAndAction:  listAndActionSaasRuntimeUnitOperation,
		DeleteFunction: testSweepSaasRuntimeUnitOperation,
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_unit",
		ListAndAction:  listAndActionSaasRuntimeUnit,
		DeleteFunction: testSweepSaasRuntimeUnit,
		Dependencies: []string{
			"google_saas_runtime_unit_operation",
		},
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_rollout_kind",
		ListAndAction:  listAndActionSaasRuntimeRolloutKind,
		DeleteFunction: testSweepSaasRuntimeRolloutKind,
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_release",
		ListAndAction:  listAndActionSaasRuntimeRelease,
		DeleteFunction: testSweepSaasRuntimeRelease,
		Dependencies: []string{
			"google_saas_runtime_unit",
			"google_saas_runtime_unit_operation",
		},
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_unit_kind",
		ListAndAction:  listAndActionSaasRuntimeUnitKind,
		DeleteFunction: testSweepSaasRuntimeUnitKind,
		Dependencies: []string{
			"google_saas_runtime_unit",
			"google_saas_runtime_release",
			"google_saas_runtime_rollout_kind",
		},
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_tenant",
		ListAndAction:  listAndActionSaasRuntimeTenant,
		DeleteFunction: testSweepSaasRuntimeTenant,
		Dependencies: []string{
			"google_saas_runtime_unit",
		},
	})

	sweeper.AddTestSweepers(&sweeper.Sweeper{
		Name:           "google_saas_runtime_saas",
		ListAndAction:  listAndActionSaasRuntimeSaas,
		DeleteFunction: testSweepSaasRuntimeSaas,
		Dependencies: []string{
			"google_saas_runtime_unit_kind",
			"google_saas_runtime_tenant",
		},
	})
}

func newSaasRuntimeSweeperConfig(region string) (*transport_tpg.Config, error) {
	if region == "" {
		region = "us-central1"
	}
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region %q: %s", region, err)
		return nil, err
	}

	config.IsMtls = transport_tpg.IsMtls()
	if os.Getenv("USER_PROJECT_OVERRIDE") == "true" {
		config.UserProjectOverride = true
	}
	if bp := os.Getenv("GOOGLE_BILLING_PROJECT"); bp != "" {
		config.BillingProject = bp
	}
	if customEndpoint := os.Getenv(Product.CustomEndpointEnvVar); customEndpoint != "" {
		if config.CustomEndpoints == nil {
			config.CustomEndpoints = make(map[string]string)
		}
		config.CustomEndpoints[Product.CustomEndpointField] = customEndpoint
	}

	if err := config.LoadAndValidate(context.Background()); err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading config: %s", err)
		return nil, err
	}
	return config, nil
}

// saasRuntimeRelativeResourcePath strips any scheme/version prefix before "projects/".
func saasRuntimeRelativeResourcePath(raw string) string {
	if idx := strings.Index(raw, "projects/"); idx >= 0 {
		return raw[idx:]
	}
	return strings.TrimPrefix(raw, "/")
}

// saasRuntimeCanonicalKey normalizes a resource reference to "locations/{loc}/{collection}/{id}"
// so references using project numbers match references using project IDs.
func saasRuntimeCanonicalKey(raw string) string {
	if idx := strings.Index(raw, "locations/"); idx >= 0 {
		return raw[idx:]
	}
	return saasRuntimeRelativeResourcePath(raw)
}

func saasRuntimeLocationFromName(raw string) string {
	rel := saasRuntimeRelativeResourcePath(raw)
	parts := strings.Split(rel, "/")
	if len(parts) >= 4 && parts[0] == "projects" && parts[2] == "locations" && parts[3] != "" {
		return parts[3]
	}
	return "us-central1"
}

func saasRuntimeProjectFromName(raw, fallback string) string {
	rel := saasRuntimeRelativeResourcePath(raw)
	parts := strings.Split(rel, "/")
	if len(parts) >= 2 && parts[0] == "projects" && parts[1] != "" {
		return parts[1]
	}
	return fallback
}

func dedupSaasRuntimeLocations(extra string, defaults []string) []string {
	seen := make(map[string]bool, len(defaults)+1)
	var out []string
	add := func(loc string) {
		loc = strings.TrimSpace(loc)
		if loc == "" || loc == "-" || seen[loc] {
			return
		}
		seen[loc] = true
		out = append(out, loc)
	}
	add(extra)
	for _, loc := range defaults {
		add(loc)
	}
	return out
}

func listSaasRuntimeAtLocation(config *transport_tpg.Config, location, collection, responseKey string) ([]map[string]interface{}, []string, error) {
	baseUrl := transport_tpg.BaseUrl(Product, config)
	listUrl := fmt.Sprintf("%sprojects/%s/locations/%s/%s", baseUrl, config.Project, location, collection)
	log.Printf("[INFO][SWEEPER_LOG] Listing SaasRuntime %s at %s", collection, listUrl)

	var items []map[string]interface{}
	var unreachable []string
	seenUnreachable := make(map[string]bool)

	pageUrl := listUrl
	for {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   config.Project,
			RawURL:    pageUrl,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return items, unreachable, err
		}

		if rawList, ok := res[responseKey].([]interface{}); ok {
			for _, rawItem := range rawList {
				if m, ok := rawItem.(map[string]interface{}); ok {
					items = append(items, m)
				}
			}
		} else if rawItems, ok := res["items"].([]interface{}); ok {
			for _, rawItem := range rawItems {
				if m, ok := rawItem.(map[string]interface{}); ok {
					items = append(items, m)
				}
			}
		}

		if rawUnreachable, ok := res["unreachable"].([]interface{}); ok {
			for _, u := range rawUnreachable {
				if s, ok := u.(string); ok && s != "" {
					loc := tpgresource.GetResourceNameFromSelfLink(s)
					if loc != "" && loc != "-" && !seenUnreachable[loc] {
						seenUnreachable[loc] = true
						unreachable = append(unreachable, loc)
					}
				}
			}
		}

		pageToken, _ := res["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
		pageUrl, err = transport_tpg.AddQueryParams(listUrl, map[string]string{"pageToken": pageToken})
		if err != nil {
			return items, unreachable, err
		}
	}

	return items, unreachable, nil
}

// listSaasRuntimeCollection lists a collection across all locations via "locations/-",
// falling back to explicit locations if "locations/-" fails or reports unreachable regions.
func listSaasRuntimeCollection(config *transport_tpg.Config, region, collection, responseKey string) ([]map[string]interface{}, error) {
	seen := make(map[string]bool)
	var result []map[string]interface{}
	appendUnique := func(items []map[string]interface{}) {
		for _, item := range items {
			name, _ := item["name"].(string)
			key := saasRuntimeCanonicalKey(name)
			if key != "" {
				if seen[key] {
					continue
				}
				seen[key] = true
			}
			result = append(result, item)
		}
	}

	items, unreachable, err := listSaasRuntimeAtLocation(config, "-", collection, responseKey)
	if err == nil {
		appendUnique(items)
		var lastErr error
		for _, loc := range unreachable {
			locItems, _, locErr := listSaasRuntimeAtLocation(config, loc, collection, responseKey)
			if locErr != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error listing unreachable location %q for %s: %s", loc, collection, locErr)
				lastErr = locErr
				continue
			}
			appendUnique(locItems)
		}
		return result, lastErr
	}

	log.Printf("[INFO][SWEEPER_LOG] Wildcard list failed for %s (%s); falling back to explicit locations", collection, err)
	var lastErr error
	anySuccess := false
	for _, loc := range dedupSaasRuntimeLocations(region, saasRuntimeFallbackLocations) {
		locItems, _, locErr := listSaasRuntimeAtLocation(config, loc, collection, responseKey)
		if locErr != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error listing location %q for %s: %s", loc, collection, locErr)
			lastErr = locErr
			continue
		}
		anySuccess = true
		appendUnique(locItems)
	}
	if anySuccess {
		return result, nil
	}
	return result, lastErr
}

func listAndActionSaasRuntimeCollection(collection, responseKey string, action sweeper.ResourceAction) error {
	config, err := newSaasRuntimeSweeperConfig("us-central1")
	if err != nil {
		return err
	}
	items, err := listSaasRuntimeCollection(config, "us-central1", collection, responseKey)
	var lastErr error
	if err != nil {
		lastErr = err
	}
	for _, obj := range items {
		name, _ := obj["name"].(string)
		loc := saasRuntimeLocationFromName(name)
		mockConfig := &tpgresource.ResourceDataMock{
			FieldsInSchema: map[string]interface{}{
				"project":  saasRuntimeProjectFromName(name, config.Project),
				"region":   loc,
				"location": loc,
				"zone":     loc + "-a",
			},
		}
		if actErr := action(config, mockConfig, obj); actErr != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error in action for %s (%s): %s", collection, name, actErr)
			lastErr = actErr
		}
	}
	return lastErr
}

func isSweepableSaasRuntimeObj(obj map[string]interface{}, knownSweepableParentKeys map[string]bool) bool {
	rawName, _ := obj["name"].(string)
	if rawName == "" {
		return false
	}
	shortName := tpgresource.GetResourceNameFromSelfLink(rawName)
	if sweeper.IsSweepableTestResource(shortName) {
		return true
	}
	// System-imported child resources (e.g. "u-<module>-<hash>" UnitKinds and Releases
	// created when a tf-test-* Saas imports an ApplicationTemplate) do not have a
	// tf-test-* short name of their own, so also check parent resource references.
	for _, refField := range []string{"saas", "unitKind", "tenant", "unit"} {
		refVal, _ := obj[refField].(string)
		if refVal == "" {
			continue
		}
		if sweeper.IsSweepableTestResource(tpgresource.GetResourceNameFromSelfLink(refVal)) {
			return true
		}
		if knownSweepableParentKeys != nil && knownSweepableParentKeys[saasRuntimeCanonicalKey(refVal)] {
			return true
		}
	}
	return false
}

func collectSweepableUnitKinds(config *transport_tpg.Config, region string) ([]map[string]interface{}, map[string]bool, error) {
	allUnitKinds, err := listSaasRuntimeCollection(config, region, "unitKinds", "unitKinds")
	sweepableKeys := make(map[string]bool)
	var sweepable []map[string]interface{}
	for _, uk := range allUnitKinds {
		if isSweepableSaasRuntimeObj(uk, nil) {
			name, _ := uk["name"].(string)
			sweepableKeys[saasRuntimeCanonicalKey(name)] = true
			sweepable = append(sweepable, uk)
		}
	}
	// Propagate through intra-collection UnitKind.dependencies in case an imported
	// UnitKind depends on or is depended on by another sweepable UnitKind.
	changed := true
	for changed {
		changed = false
		for _, uk := range allUnitKinds {
			name, _ := uk["name"].(string)
			key := saasRuntimeCanonicalKey(name)
			if sweepableKeys[key] {
				continue
			}
			for _, depKey := range extractUnitKindDependencyKeys(uk) {
				if sweepableKeys[depKey] {
					sweepableKeys[key] = true
					sweepable = append(sweepable, uk)
					changed = true
					break
				}
			}
		}
	}
	return sweepable, sweepableKeys, err
}

func isRetryableSaasRuntimeDeleteError(err error) bool {
	if err == nil {
		return false
	}
	if !transport_tpg.IsGoogleApiErrorWithCode(err, 400) {
		return false
	}
	msg := err.Error()
	// A provisioned Unit will never become deletable by simply retrying DELETE.
	if strings.Contains(msg, "unit has been provisioned") {
		return false
	}
	for _, transientSubstr := range []string{
		"is referenced as default_release",
		"has dependent",
		"with a dependency on",
		"is in RUNNING state",
		"has Units",
		"is not completed",
		"is executing",
		"still provisioned",
	} {
		if strings.Contains(msg, transientSubstr) {
			return true
		}
	}
	return false
}

func deleteSaasRuntimeResource(config *transport_tpg.Config, resourceName, fullName string, timing saasRuntimeSweeperTiming) (bool, error) {
	relPath := saasRuntimeRelativeResourcePath(fullName)
	if relPath == "" {
		return false, nil
	}
	deleteUrl := transport_tpg.BaseUrl(Product, config) + relPath
	deadline := time.Now().Add(timing.deleteRetryTimeout)

	for {
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "DELETE",
			Project:   saasRuntimeProjectFromName(relPath, config.Project),
			RawURL:    deleteUrl,
			UserAgent: config.UserAgent,
		})
		if err == nil {
			log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s: %s", resourceName, relPath)
			return true, nil
		}
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return false, nil
		}
		if !isRetryableSaasRuntimeDeleteError(err) || time.Now().After(deadline) {
			log.Printf("[INFO][SWEEPER_LOG] Error deleting %s (%s): %s", resourceName, deleteUrl, err)
			return false, err
		}
		time.Sleep(timing.pollInterval)
	}
}

func waitForSaasRuntimeAbsence(config *transport_tpg.Config, resourceName, fullName string, timing saasRuntimeSweeperTiming) error {
	relPath := saasRuntimeRelativeResourcePath(fullName)
	if relPath == "" {
		return nil
	}
	getUrl := transport_tpg.BaseUrl(Product, config) + relPath
	project := saasRuntimeProjectFromName(relPath, config.Project)
	deadline := time.Now().Add(timing.absenceTimeout)

	for {
		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    getUrl,
			UserAgent: config.UserAgent,
		})
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return nil
		}
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error polling %s absence (%s): %s", resourceName, getUrl, err)
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for %s %s to be deleted", resourceName, relPath)
		}
		time.Sleep(timing.pollInterval)
	}
}

func deleteAndWaitSaasRuntimeWave(config *transport_tpg.Config, resourceName string, items []map[string]interface{}, timing saasRuntimeSweeperTiming) error {
	var lastErr error
	var pendingAbsence []string

	for _, obj := range items {
		name, _ := obj["name"].(string)
		deleted, err := deleteSaasRuntimeResource(config, resourceName, name, timing)
		if err != nil {
			lastErr = err
			continue
		}
		if deleted {
			pendingAbsence = append(pendingAbsence, name)
		}
	}

	for _, name := range pendingAbsence {
		if err := waitForSaasRuntimeAbsence(config, resourceName, name, timing); err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error waiting for %s deletion (%s): %s", resourceName, name, err)
			lastErr = err
		}
	}

	return lastErr
}

// saasRuntimeTopologicalWaves groups items into ordered deletion waves such that if
// item A must be deleted before item B (an edge A -> B), A is placed in an earlier wave.
func saasRuntimeTopologicalWaves(items []map[string]interface{}, edgesFn func(items []map[string]interface{}, byKey map[string]int) [][2]int) [][]map[string]interface{} {
	n := len(items)
	if n == 0 {
		return nil
	}

	byKey := make(map[string]int, n)
	for i, item := range items {
		name, _ := item["name"].(string)
		key := saasRuntimeCanonicalKey(name)
		if key != "" {
			byKey[key] = i
		}
	}

	adj := make([][]int, n)
	inDegree := make([]int, n)
	seenEdge := make(map[[2]int]bool)

	for _, e := range edgesFn(items, byKey) {
		from, to := e[0], e[1]
		if from == to || from < 0 || from >= n || to < 0 || to >= n {
			continue
		}
		pair := [2]int{from, to}
		if seenEdge[pair] {
			continue
		}
		seenEdge[pair] = true
		adj[from] = append(adj[from], to)
		inDegree[to]++
	}

	placed := make([]bool, n)
	placedCount := 0
	var waves [][]map[string]interface{}

	for placedCount < n {
		var ready []int
		for i := 0; i < n; i++ {
			if !placed[i] && inDegree[i] == 0 {
				ready = append(ready, i)
			}
		}
		if len(ready) == 0 {
			// Break any unexpected cycle deterministically without dropping items.
			var remaining []map[string]interface{}
			for i := 0; i < n; i++ {
				if !placed[i] {
					placed[i] = true
					placedCount++
					remaining = append(remaining, items[i])
				}
			}
			waves = append(waves, remaining)
			break
		}

		wave := make([]map[string]interface{}, 0, len(ready))
		for _, idx := range ready {
			placed[idx] = true
			placedCount++
			wave = append(wave, items[idx])
		}
		for _, idx := range ready {
			for _, next := range adj[idx] {
				inDegree[next]--
			}
		}
		waves = append(waves, wave)
	}

	return waves
}

func extractUnitKindDependencyKeys(obj map[string]interface{}) []string {
	rawDeps, ok := obj["dependencies"].([]interface{})
	if !ok {
		return nil
	}
	var keys []string
	for _, rawDep := range rawDeps {
		depMap, ok := rawDep.(map[string]interface{})
		if !ok {
			continue
		}
		if uk, _ := depMap["unitKind"].(string); uk != "" {
			keys = append(keys, saasRuntimeCanonicalKey(uk))
		}
	}
	return keys
}

func clearUnitKindDefaultRelease(config *transport_tpg.Config, unitKinds []map[string]interface{}) error {
	baseUrl := transport_tpg.BaseUrl(Product, config)
	var lastErr error
	for _, uk := range unitKinds {
		defaultRelease, _ := uk["defaultRelease"].(string)
		if defaultRelease == "" {
			continue
		}
		name, _ := uk["name"].(string)
		relPath := saasRuntimeRelativeResourcePath(name)
		if relPath == "" {
			continue
		}
		patchUrl, err := transport_tpg.AddQueryParams(baseUrl+relPath, map[string]string{"updateMask": "defaultRelease"})
		if err != nil {
			lastErr = err
			continue
		}
		log.Printf("[INFO][SWEEPER_LOG] Clearing defaultRelease on SaasRuntimeUnitKind %s (was %s)", relPath, defaultRelease)
		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			Project:   saasRuntimeProjectFromName(relPath, config.Project),
			RawURL:    patchUrl,
			UserAgent: config.UserAgent,
			Body: map[string]interface{}{
				"defaultRelease": "",
			},
		})
		if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			log.Printf("[INFO][SWEEPER_LOG] Error clearing defaultRelease on %s: %s", relPath, err)
			lastErr = err
		}
	}
	return lastErr
}

func cancelUnitOperationIfNeeded(config *transport_tpg.Config, obj map[string]interface{}, timing saasRuntimeSweeperTiming) {
	state, _ := obj["state"].(string)
	if state == "" || saasRuntimeTerminalUnitOpStates[state] {
		return
	}
	name, _ := obj["name"].(string)
	unit, _ := obj["unit"].(string)
	relPath := saasRuntimeRelativeResourcePath(name)
	if relPath == "" || unit == "" {
		return
	}
	baseUrl := transport_tpg.BaseUrl(Product, config)
	project := saasRuntimeProjectFromName(relPath, config.Project)
	patchUrl, err := transport_tpg.AddQueryParams(baseUrl+relPath, map[string]string{"updateMask": "cancel"})
	if err != nil {
		return
	}

	log.Printf("[INFO][SWEEPER_LOG] Cancelling non-terminal SaasRuntimeUnitOperation %s (state=%s)", relPath, state)
	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "PATCH",
		Project:   project,
		RawURL:    patchUrl,
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"cancel": true,
			"unit":   unit,
		},
	})
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Unable to cancel SaasRuntimeUnitOperation %s: %s", relPath, err)
		return
	}

	getUrl := baseUrl + relPath
	deadline := time.Now().Add(timing.operationTimeout)
	for !time.Now().After(deadline) {
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    getUrl,
			UserAgent: config.UserAgent,
		})
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			return
		}
		if err != nil {
			return
		}
		curState, _ := res["state"].(string)
		if saasRuntimeTerminalUnitOpStates[curState] {
			return
		}
		time.Sleep(timing.pollInterval)
	}
}

func isSaasRuntimeUnitProvisioned(obj map[string]interface{}) bool {
	if rawConds, ok := obj["conditions"].([]interface{}); ok {
		for _, rawCond := range rawConds {
			cond, ok := rawCond.(map[string]interface{})
			if !ok {
				continue
			}
			cType, _ := cond["type"].(string)
			cStatus, _ := cond["status"].(string)
			if cType == "TYPE_PROVISIONED" {
				return cStatus == "STATUS_TRUE"
			}
		}
	}
	state, _ := obj["state"].(string)
	switch state {
	case "UNIT_STATE_READY", "UNIT_STATE_UPDATING", "UNIT_STATE_PROVISIONING", "UNIT_STATE_DEPROVISIONING":
		return true
	default:
		return false
	}
}

func deprovisionSaasRuntimeUnitIfNeeded(config *transport_tpg.Config, obj map[string]interface{}, timing saasRuntimeSweeperTiming) {
	if !isSaasRuntimeUnitProvisioned(obj) {
		return
	}
	name, _ := obj["name"].(string)
	relPath := saasRuntimeRelativeResourcePath(name)
	if relPath == "" {
		return
	}

	project := saasRuntimeProjectFromName(relPath, config.Project)
	loc := saasRuntimeLocationFromName(relPath)
	baseUrl := transport_tpg.BaseUrl(Product, config)

	h := fnv.New32a()
	_, _ = h.Write([]byte(relPath))
	opID := fmt.Sprintf("tf-test-sweep-deprov-%08x-%d", h.Sum32(), time.Now().Unix()%100000)
	createUrl := fmt.Sprintf("%sprojects/%s/locations/%s/unitOperations?unitOperationId=%s", baseUrl, project, loc, opID)
	opRelPath := fmt.Sprintf("projects/%s/locations/%s/unitOperations/%s", project, loc, opID)

	log.Printf("[INFO][SWEEPER_LOG] Creating deprovision UnitOperation %s for provisioned Unit %s", opRelPath, relPath)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   project,
		RawURL:    createUrl,
		UserAgent: config.UserAgent,
		Body: map[string]interface{}{
			"unit":        relPath,
			"deprovision": map[string]interface{}{},
		},
	})
	if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 409) {
		log.Printf("[INFO][SWEEPER_LOG] Could not create deprovision UnitOperation for %s: %s", relPath, err)
		return
	}

	opGetUrl := baseUrl + opRelPath
	deadline := time.Now().Add(timing.operationTimeout)
	for !time.Now().After(deadline) {
		res, getErr := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    opGetUrl,
			UserAgent: config.UserAgent,
		})
		if transport_tpg.IsGoogleApiErrorWithCode(getErr, 404) {
			break
		}
		if getErr != nil {
			break
		}
		opState, _ := res["state"].(string)
		if saasRuntimeTerminalUnitOpStates[opState] {
			break
		}
		time.Sleep(timing.pollInterval)
	}

	if deleted, delErr := deleteSaasRuntimeResource(config, "SaasRuntimeUnitOperation", opRelPath, timing); delErr == nil && deleted {
		_ = waitForSaasRuntimeAbsence(config, "SaasRuntimeUnitOperation", opRelPath, timing)
	}
}

// --- Per-resource sweepers & ListAndAction functions ---

func testSweepSaasRuntimeUnitOperation(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeUnitOperations(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeUnitOperations(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeUnitOperation")
	_, sweepableUKKeys, _ := collectSweepableUnitKinds(config, region)
	sweepableParentKeys := make(map[string]bool, len(sweepableUKKeys))
	for k, v := range sweepableUKKeys {
		sweepableParentKeys[k] = v
	}
	if allUnits, unitErr := listSaasRuntimeCollection(config, region, "units", "units"); unitErr == nil {
		for _, u := range allUnits {
			if isSweepableSaasRuntimeObj(u, sweepableUKKeys) {
				if uName, _ := u["name"].(string); uName != "" {
					sweepableParentKeys[saasRuntimeCanonicalKey(uName)] = true
				}
			}
		}
	}

	allOps, err := listSaasRuntimeCollection(config, region, "unitOperations", "unitOperations")
	var lastErr error
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allOps {
		if isSweepableSaasRuntimeObj(obj, sweepableParentKeys) {
			cancelUnitOperationIfNeeded(config, obj, timing)
			sweepable = append(sweepable, obj)
		}
	}

	if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeUnitOperation", sweepable, timing); waveErr != nil {
		lastErr = waveErr
	}
	return lastErr
}

func listAndActionSaasRuntimeUnitOperation(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("unitOperations", "unitOperations", action)
}

func testSweepSaasRuntimeUnit(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeUnits(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeUnits(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeUnit")
	_, sweepableUKKeys, _ := collectSweepableUnitKinds(config, region)
	allUnits, err := listSaasRuntimeCollection(config, region, "units", "units")
	var lastErr error
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allUnits {
		if isSweepableSaasRuntimeObj(obj, sweepableUKKeys) {
			sweepable = append(sweepable, obj)
		}
	}

	waves := saasRuntimeTopologicalWaves(sweepable, func(items []map[string]interface{}, byKey map[string]int) [][2]int {
		var edges [][2]int
		for i, item := range items {
			if rawDeps, ok := item["dependencies"].([]interface{}); ok {
				for _, rawDep := range rawDeps {
					if depMap, ok := rawDep.(map[string]interface{}); ok {
						if depUnit, _ := depMap["unit"].(string); depUnit != "" {
							if j, exists := byKey[saasRuntimeCanonicalKey(depUnit)]; exists {
								// Dependent unit i must be deleted before dependency unit j.
								edges = append(edges, [2]int{i, j})
							}
						}
					}
				}
			}
			if rawDependents, ok := item["dependents"].([]interface{}); ok {
				for _, rawDep := range rawDependents {
					if depMap, ok := rawDep.(map[string]interface{}); ok {
						if depUnit, _ := depMap["unit"].(string); depUnit != "" {
							if j, exists := byKey[saasRuntimeCanonicalKey(depUnit)]; exists {
								// Dependent unit j must be deleted before base unit i.
								edges = append(edges, [2]int{j, i})
							}
						}
					}
				}
			}
		}
		return edges
	})

	for _, wave := range waves {
		for _, obj := range wave {
			deprovisionSaasRuntimeUnitIfNeeded(config, obj, timing)
		}
		if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeUnit", wave, timing); waveErr != nil {
			lastErr = waveErr
		}
	}
	return lastErr
}

func listAndActionSaasRuntimeUnit(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("units", "units", action)
}

func testSweepSaasRuntimeRolloutKind(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeRolloutKinds(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeRolloutKinds(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeRolloutKind")
	_, sweepableUKKeys, _ := collectSweepableUnitKinds(config, region)
	allRolloutKinds, err := listSaasRuntimeCollection(config, region, "rolloutKinds", "rolloutKinds")
	var lastErr error
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allRolloutKinds {
		if isSweepableSaasRuntimeObj(obj, sweepableUKKeys) {
			sweepable = append(sweepable, obj)
		}
	}

	if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeRolloutKind", sweepable, timing); waveErr != nil {
		lastErr = waveErr
	}
	return lastErr
}

func listAndActionSaasRuntimeRolloutKind(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("rolloutKinds", "rolloutKinds", action)
}

func testSweepSaasRuntimeRelease(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeReleases(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeReleases(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeRelease")
	var lastErr error

	// Break the UnitKind.defaultRelease <-> Release.unitKind deletion deadlock
	// before attempting to delete Releases.
	sweepableUnitKinds, sweepableUKKeys, ukErr := collectSweepableUnitKinds(config, region)
	if ukErr != nil {
		lastErr = ukErr
	}
	if patchErr := clearUnitKindDefaultRelease(config, sweepableUnitKinds); patchErr != nil {
		lastErr = patchErr
	}

	allReleases, err := listSaasRuntimeCollection(config, region, "releases", "releases")
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allReleases {
		if isSweepableSaasRuntimeObj(obj, sweepableUKKeys) {
			sweepable = append(sweepable, obj)
		}
	}

	if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeRelease", sweepable, timing); waveErr != nil {
		lastErr = waveErr
	}
	return lastErr
}

func listAndActionSaasRuntimeRelease(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("releases", "releases", action)
}

func testSweepSaasRuntimeUnitKind(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeUnitKinds(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeUnitKinds(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeUnitKind")
	sweepable, sweepableUKKeys, err := collectSweepableUnitKinds(config, region)
	var lastErr error
	if err != nil {
		lastErr = err
	}
	if len(sweepable) == 0 {
		return lastErr
	}

	// If invoked directly outside the full DAG, ensure defaultRelease is cleared
	// and any remaining sweepable Releases referencing these UnitKinds are deleted first.
	if patchErr := clearUnitKindDefaultRelease(config, sweepable); patchErr != nil {
		lastErr = patchErr
	}
	if allReleases, relErr := listSaasRuntimeCollection(config, region, "releases", "releases"); relErr == nil {
		var blockingReleases []map[string]interface{}
		for _, r := range allReleases {
			ukRef, _ := r["unitKind"].(string)
			if ukRef != "" && sweepableUKKeys[saasRuntimeCanonicalKey(ukRef)] && isSweepableSaasRuntimeObj(r, sweepableUKKeys) {
				blockingReleases = append(blockingReleases, r)
			}
		}
		if len(blockingReleases) > 0 {
			if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeRelease", blockingReleases, timing); waveErr != nil {
				lastErr = waveErr
			}
		}
	}

	waves := saasRuntimeTopologicalWaves(sweepable, func(items []map[string]interface{}, byKey map[string]int) [][2]int {
		var edges [][2]int
		for i, item := range items {
			for _, depKey := range extractUnitKindDependencyKeys(item) {
				if j, exists := byKey[depKey]; exists {
					// Dependent UnitKind i must be deleted and purged before base UnitKind j.
					edges = append(edges, [2]int{i, j})
				}
			}
		}
		return edges
	})

	for _, wave := range waves {
		if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeUnitKind", wave, timing); waveErr != nil {
			lastErr = waveErr
		}
	}
	return lastErr
}

func listAndActionSaasRuntimeUnitKind(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("unitKinds", "unitKinds", action)
}

func testSweepSaasRuntimeTenant(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeTenants(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeTenants(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeTenant")
	allTenants, err := listSaasRuntimeCollection(config, region, "tenants", "tenants")
	var lastErr error
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allTenants {
		if isSweepableSaasRuntimeObj(obj, nil) {
			sweepable = append(sweepable, obj)
		}
	}

	if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeTenant", sweepable, timing); waveErr != nil {
		lastErr = waveErr
	}
	return lastErr
}

func listAndActionSaasRuntimeTenant(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("tenants", "tenants", action)
}

func testSweepSaasRuntimeSaas(region string) error {
	config, err := newSaasRuntimeSweeperConfig(region)
	if err != nil {
		return err
	}
	return sweepSaasRuntimeSaas(config, region, defaultSaasRuntimeSweeperTiming)
}

func sweepSaasRuntimeSaas(config *transport_tpg.Config, region string, timing saasRuntimeSweeperTiming) error {
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for SaasRuntimeSaas")
	allSaas, err := listSaasRuntimeCollection(config, region, "saas", "saas")
	var lastErr error
	if err != nil {
		lastErr = err
	}

	var sweepable []map[string]interface{}
	for _, obj := range allSaas {
		if isSweepableSaasRuntimeObj(obj, nil) {
			sweepable = append(sweepable, obj)
		}
	}

	if waveErr := deleteAndWaitSaasRuntimeWave(config, "SaasRuntimeSaas", sweepable, timing); waveErr != nil {
		lastErr = waveErr
	}
	return lastErr
}

func listAndActionSaasRuntimeSaas(action sweeper.ResourceAction) error {
	return listAndActionSaasRuntimeCollection("saas", "saas", action)
}
