package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

var SharedServicePerimeterProjectPrefix = "tf-bootstrap-sp-"

func BootstrapServicePerimeterProjects(t *testing.T, desiredProjects int) []*cloudresourcemanager.Project {
	config := BootstrapConfig(t)
	if config == nil {
		return nil
	}

	org := GetTestOrgFromEnv(t)

	// The filter endpoint works differently if you provide both the parent id and parent type, and
	// doesn't seem to allow for prefix matching. Don't change this to include the parent type unless
	// that API behavior changes.
	prefixFilter := fmt.Sprintf("id:%s* parent.id:%s", SharedServicePerimeterProjectPrefix, org)
	res, err := config.NewResourceManagerClient(config.UserAgent).Projects.List().Filter(prefixFilter).Do()
	if err != nil {
		t.Fatalf("Error getting shared test projects: %s", err)
	}

	projects := res.Projects
	for len(projects) < desiredProjects {
		pid := SharedServicePerimeterProjectPrefix + RandString(t, 10)
		project := &cloudresourcemanager.Project{
			ProjectId: pid,
			Name:      "TF Service Perimeter Test",
			Parent: &cloudresourcemanager.ResourceId{
				Type: "organization",
				Id:   org,
			},
		}
		op, err := config.NewResourceManagerClient(config.UserAgent).Projects.Create(project).Do()
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		opAsMap, err := google.ConvertToMap(op)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		err = ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 4)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		p, err := config.NewResourceManagerClient(config.UserAgent).Projects.Get(pid).Do()
		if err != nil {
			t.Fatalf("Error getting shared test project: %s", err)
		}
		projects = append(projects, p)
	}

	return projects
}
