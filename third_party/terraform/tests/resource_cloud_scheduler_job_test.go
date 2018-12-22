package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	cloudscheduler "google.golang.org/api/cloudscheduler/v1beta1"
)

func TestAccCloudSchedulerJob_basic(t *testing.T) {
	t.Parallel()

	var job cloudscheduler.Job

	jobResourceName := "google_cloud_scheduler_job.job"
	jobName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSchedulerJob_pubSub(jobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudSchedulerJobExists(
						jobResourceName, &job),
					resource.TestCheckResourceAttr(jobResourceName,
						"name", jobName),
					resource.TestCheckResourceAttr(jobResourceName,
						"description", "test job"),
					resource.TestCheckResourceAttr(jobResourceName,
						"schedule", "* * * * *"),
					resource.TestCheckResourceAttr(jobResourceName,
						"time_zone", "Europe/London"),
					resource.TestCheckResourceAttr(jobResourceName,
						"pubsub_target", ""),
				),
			},
			{
				ResourceName:      jobResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudSchedulerJobExists(n string, job *cloudscheduler.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, name)
		found, err := config.clientScheduler.Projects.Locations.Jobs.Get(jobName).Do()
		if err != nil {
			return fmt.Errorf("CloudScheduler Job not present")
		}

		*job = *found

		return nil
	}
}

func testAccCheckCloudSchedulerJobDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloud_scheduler_job" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, name)

		_, err := config.clientScheduler.Projects.Locations.Jobs.Get(jobName).Do()
		if err == nil {
			return fmt.Errorf("Function still exists")
		}

	}

	return nil
}

func testAccCloudSchedulerJob_pubSub(name string) string {
	return fmt.Sprintf(`

resource "google_pubsub_topic" "topic" {
	name = "build-triggers"
}

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	schedule = "*/2 * * * *"

	pubsub_target = {
		topic_name = "${google_pubsub_topic.topic.name}"
	}
}

	`, name)
}

func testAccCloudSchedulerJob_appEngine(name string) string {
	return fmt.Sprintf(`

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	schedule = "*/4 * * * *"

	// TODO defaults to the default service, investigation required
	app_engine_http_target = {
		relative_uri = "/ping"
	}
}

	`, name)
}

func testAccCloudSchedulerJob_http(name string) string {
	return fmt.Sprintf(`

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	schedule = "*/8 * * * *"

	http_target = {
		uri = "https://example.com/ping"
	}
}

	`, name)
}
