package resourcemanager

import (
	"fmt"
	"log"
	"time"

	rmClient "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager/client"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const (
	batchKeyTmplServiceUsageEnableServices = "project/%s/services:batchEnable"
	batchKeyTmplServiceUsageListServices   = "project/%s/services"
	batchKeyTmplServiceUsageGetProject     = "project/%s/get"
)

// BatchRequestEnableServices can be used to batch requests to enable services
// across resource nodes, i.e. to batch creation of several
// google_project_service(s) resources.
func BatchRequestEnableService(service string, project string, d *schema.ResourceData, config *transport_tpg.Config) error {
	// Renamed service create calls are relatively likely to fail, so don't try to batch the call.
	if altName, ok := renamedServicesByOldAndNewServiceNames[service]; ok {
		return tryEnableRenamedService(service, altName, project, d, config)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	req := &transport_tpg.BatchRequest{
		ResourceName: project,
		Body:         []string{service},
		CombineF:     combineServiceUsageServicesBatches,
		SendF:        sendBatchFuncEnableServices(config, userAgent, billingProject, d.Timeout(schema.TimeoutCreate)),
		DebugId:      fmt.Sprintf("Enable Project Service %q for project %q", service, project),
	}

	_, err = config.RequestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageEnableServices, project),
		req,
		d.Timeout(schema.TimeoutCreate))
	return err
}

func tryEnableRenamedService(service, altName string, project string, d *schema.ResourceData, config *transport_tpg.Config) error {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] found renamed service %s (with alternate name %s)", service, altName)
	// use a short timeout- failures are likely

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	log.Printf("[DEBUG] attempting enabling service with user-specified name %s", service)
	err = EnableServiceUsageProjectServices([]string{service}, project, billingProject, userAgent, config, 1*time.Minute)
	if err != nil {
		log.Printf("[DEBUG] saw error %s. attempting alternate name %v", err, altName)
		err2 := EnableServiceUsageProjectServices([]string{altName}, project, billingProject, userAgent, config, 1*time.Minute)
		if err2 != nil {
			return fmt.Errorf("Saw 2 subsequent errors attempting to enable a renamed service: %s / %s", err, err2)
		}
	}
	return nil
}

func BatchRequestReadServices(project string, d *schema.ResourceData, config *transport_tpg.Config) (interface{}, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	req := &transport_tpg.BatchRequest{
		ResourceName: project,
		Body:         nil,
		// Use empty CombineF since the request is exactly the same no matter how many services we read.
		CombineF: func(body interface{}, toAdd interface{}) (interface{}, error) { return nil, nil },
		SendF:    sendListServices(config, billingProject, userAgent, d.Timeout(schema.TimeoutRead)),
		DebugId:  fmt.Sprintf("List Project Services %s", project),
	}

	return config.RequestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageListServices, project),
		req,
		d.Timeout(schema.TimeoutRead))
}

func BatchRequestReadProject(project string, d *schema.ResourceData, config *transport_tpg.Config) (*cloudresourcemanager.Project, error) {
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	req := &transport_tpg.BatchRequest{
		ResourceName: project,
		Body:         nil,
		// Use empty CombineF since the request is exactly the same no matter how many callers ask for this project.
		CombineF: func(body interface{}, toAdd interface{}) (interface{}, error) { return nil, nil },
		SendF:    sendGetProject(config, billingProject, userAgent),
		DebugId:  fmt.Sprintf("Get Project %s", project),
	}

	resp, err := config.RequestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageGetProject, project),
		req,
		d.Timeout(schema.TimeoutRead))
	if err != nil {
		return nil, err
	}
	return resp.(*cloudresourcemanager.Project), nil
}

func combineServiceUsageServicesBatches(srvsRaw interface{}, toAddRaw interface{}) (interface{}, error) {
	srvs, ok := srvsRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", srvsRaw)
	}
	toAdd, ok := toAddRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected new request body type to be []string, got %v. This is a provider error.", toAdd)
	}

	return append(srvs, toAdd...), nil
}

func sendBatchFuncEnableServices(config *transport_tpg.Config, userAgent, billingProject string, timeout time.Duration) transport_tpg.BatcherSendFunc {
	return func(project string, toEnableRaw interface{}) (interface{}, error) {
		toEnable, ok := toEnableRaw.([]string)
		if !ok {
			return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", toEnableRaw)
		}
		return nil, EnableServiceUsageProjectServices(toEnable, project, billingProject, userAgent, config, timeout)
	}
}

func sendListServices(config *transport_tpg.Config, billingProject, userAgent string, timeout time.Duration) transport_tpg.BatcherSendFunc {
	return func(project string, _ interface{}) (interface{}, error) {
		return ListCurrentlyEnabledServices(project, billingProject, userAgent, config, timeout)
	}
}

func sendGetProject(config *transport_tpg.Config, billingProject, userAgent string) transport_tpg.BatcherSendFunc {
	return func(project string, _ interface{}) (interface{}, error) {
		projectGetCall := rmClient.NewClient(config, userAgent).Projects.Get(project)
		if config.UserProjectOverride {
			projectGetCall.Header().Add("X-Goog-User-Project", billingProject)
		}
		return projectGetCall.Do()
	}
}
