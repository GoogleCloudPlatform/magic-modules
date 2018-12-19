package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/container/v1beta1"
)

type ContainerOperationWaiter struct {
	Service  *container.Service
	Op       *container.Operation
	Project  string
	Location string
}

func (w *ContainerOperationWaiter) Conf(timeoutMinutes int) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"DONE"},
		Refresh:    w.RefreshFunc(),
		Timeout:    time.Duration(timeoutMinutes) * time.Minute,
		MinTimeout: 2 * time.Second,
	}
}

func (w *ContainerOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		name := fmt.Sprintf("projects/%s/locations/%s/operations/%s",
			w.Project, w.Location, w.Op.Name)
		resp, err := w.Service.Projects.Locations.Operations.Get(name).Do()

		if err != nil {
			return nil, "", err
		}

		if resp.StatusMessage != "" {
			return resp, resp.Status, fmt.Errorf(resp.StatusMessage)
		}

		log.Printf("[DEBUG] Progress of operation %q: %q", w.Op.Name, resp.Status)

		return resp, resp.Status, err
	}
}

func containerOperationWait(config *Config, op *container.Operation, project, location, activity string, timeoutMinutes int) error {
	if op.Status == "DONE" {
		if op.StatusMessage != "" {
			return fmt.Errorf(op.StatusMessage)
		}
		return nil
	}

	w := &ContainerOperationWaiter{
		Service:  config.clientContainerBeta,
		Op:       op,
		Project:  project,
		Location: location,
	}

	_, err := w.Conf(timeoutMinutes).WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}
	return nil
}
