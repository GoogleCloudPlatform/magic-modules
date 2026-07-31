package datalineage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OpenLineage/openlineage/client/go/pkg/openlineage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func getOpenLineageBillingProject(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return "", fmt.Errorf("error fetching project for OpenLineageJob: %w", err)
	}

	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		project = bp
	}

	return project, nil
}

func getOpenLineageProjectAndLocation(d *schema.ResourceData, config *transport_tpg.Config) (string, string, error) {
	project, err := tpgresource.GetProject(d, config)
	if err != nil || strings.TrimSpace(project) == "" {
		if err != nil {
			return "", "", fmt.Errorf("missing project for processOpenLineageRunEvent: %w", err)
		}
		return "", "", fmt.Errorf("missing project for processOpenLineageRunEvent")
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil || strings.TrimSpace(location) == "" {
		if err != nil {
			return "", "", fmt.Errorf("missing location for processOpenLineageRunEvent: %w", err)
		}
		return "", "", fmt.Errorf("missing location for processOpenLineageRunEvent")
	}

	return project, location, nil
}

func emitEvent(ctx context.Context, d *schema.ResourceData, runEvent *openlineage.RunEvent, config *transport_tpg.Config, userAgent string, timeout time.Duration) (map[string]interface{}, diag.Diagnostics) {
	_ = ctx

	eventJSON, err := json.Marshal(runEvent)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	payload := map[string]any{}
	if err := json.Unmarshal(eventJSON, &payload); err != nil {
		return nil, diag.FromErr(err)
	}

	project, location, err := getOpenLineageProjectAndLocation(d, config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	url := transport_tpg.BaseUrl(Product, config) + fmt.Sprintf("projects/%s/locations/%s:processOpenLineageRunEvent", project, location)

	billingProject, err := getOpenLineageBillingProject(d, config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      payload,
		Timeout:   timeout,
		Headers:   make(http.Header),
	})
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("processOpenLineageRunEvent: %w", err))
	}

	return resp, nil
}

func getLatestRunForProcess(ctx context.Context, d *schema.ResourceData, config *transport_tpg.Config, process string, userAgent string) (string, diag.Diagnostics) {
	_ = ctx

	if _, _, err := getOpenLineageProjectAndLocation(d, config); err != nil {
		return "", diag.FromErr(err)
	}

	billingProject, err := getOpenLineageBillingProject(d, config)
	if err != nil {
		return "", diag.FromErr(err)
	}

	processURL := transport_tpg.BaseUrl(Product, config) + strings.TrimPrefix(process, "/")

	_, pErr := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    processURL,
		UserAgent: userAgent,
		Headers:   make(http.Header),
	})
	if pErr != nil {
		return "", diag.FromErr(transport_tpg.HandleNotFoundError(pErr, d, fmt.Sprintf("DataLineageOpenLineageJob %q", process)))
	}

	runsURL := strings.TrimSuffix(transport_tpg.BaseUrl(Product, config)+strings.TrimPrefix(process, "/"), "/") + "/runs"
	runsURL, err = transport_tpg.AddQueryParams(runsURL, map[string]string{"pageSize": "1"})
	if err != nil {
		return "", diag.FromErr(err)
	}
	runsResponse, rErr := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    runsURL,
		UserAgent: userAgent,
		Headers:   make(http.Header),
	})
	if rErr != nil {
		return "", diag.FromErr(fmt.Errorf("error retrieving latest run for process %s: %w", process, rErr))
	}

	runs, ok := runsResponse["runs"].([]interface{})
	if !ok || len(runs) == 0 {
		return "", diag.Errorf("error retrieving latest run for process %s: no runs found", process)
	}

	firstRun, ok := runs[0].(map[string]interface{})
	if !ok {
		return "", diag.Errorf("error retrieving latest run for process %s: invalid run format", process)
	}

	runName, ok := firstRun["name"].(string)
	if !ok || runName == "" {
		return "", diag.Errorf("error retrieving latest run for process %s: missing run name", process)
	}

	return runName, nil
}

func deleteProcess(ctx context.Context, d *schema.ResourceData, config *transport_tpg.Config, process string, userAgent string) error {
	_ = ctx

	if _, _, err := getOpenLineageProjectAndLocation(d, config); err != nil {
		return err
	}

	billingProject, err := getOpenLineageBillingProject(d, config)
	if err != nil {
		return err
	}

	url := transport_tpg.BaseUrl(Product, config) + strings.TrimPrefix(process, "/")

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "DELETE",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Timeout:   d.Timeout(schema.TimeoutDelete),
		Headers:   make(http.Header),
	})
	if err != nil {
		return err
	}

	return nil
}
