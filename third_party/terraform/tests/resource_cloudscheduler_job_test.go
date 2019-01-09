package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudSchedulerJob_pubsub(t *testing.T) {
	t.Parallel()

	jobResourceName := "google_cloud_scheduler_job.job"
	pubSubJobName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSchedulerJob_pubSubConfig(pubSubJobName, project),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudSchedulerJobExists(jobResourceName),
					resource.TestCheckResourceAttr(jobResourceName, "name", pubSubJobName),
					resource.TestCheckResourceAttr(jobResourceName, "description", "test job"),
					resource.TestCheckResourceAttr(jobResourceName, "schedule", "*/2 * * * *"),
					resource.TestCheckResourceAttr(jobResourceName, "time_zone", "Etc/UTC"),
				),
			},
		},
	})
}

func TestAccCloudSchedulerJob_http(t *testing.T) {
	t.Parallel()

	jobResourceName := "google_cloud_scheduler_job.job"
	httpJobName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSchedulerJob_httpConfig(httpJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudSchedulerJobExists(jobResourceName),
					resource.TestCheckResourceAttr(jobResourceName, "name", httpJobName),
					resource.TestCheckResourceAttr(jobResourceName, "description", "test http job"),
					resource.TestCheckResourceAttr(jobResourceName, "schedule", "*/8 * * * *"),
					resource.TestCheckResourceAttr(jobResourceName, "time_zone", "America/New_York"),
				),
			},
		},
	})
}

func TestAccCloudSchedulerJob_appEngine(t *testing.T) {
	t.Parallel()

	jobResourceName := "google_cloud_scheduler_job.job"
	appEngineJobName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudSchedulerJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSchedulerJob_appEngineConfig(appEngineJobName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudSchedulerJobExists(jobResourceName),
					resource.TestCheckResourceAttr(jobResourceName, "name", appEngineJobName),
					resource.TestCheckResourceAttr(jobResourceName, "description", "test app engine job"),
					resource.TestCheckResourceAttr(jobResourceName, "schedule", "*/4 * * * *"),
					resource.TestCheckResourceAttr(jobResourceName, "time_zone", "Europe/London"),
				),
			},
		},
	})
}

func testAccCloudSchedulerJobExists(n string) resource.TestCheckFunc {
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
		region := getTestRegionFromEnv()
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, name)
		_, err := config.clientCloudScheduler.Projects.Locations.Jobs.Get(jobName).Do()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("CloudScheduler Job not present %s %s %s", project, region, name))
		}

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
		region := getTestRegionFromEnv()
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, name)

		_, err := config.clientCloudScheduler.Projects.Locations.Jobs.Get(jobName).Do()
		if err == nil {
			return fmt.Errorf("Function still exists")
		}

	}

	return nil
}

func testAccCloudSchedulerJob_pubSubConfig(name string, project string) string {
	return fmt.Sprintf(`

resource "google_pubsub_topic" "topic" {
	name = "build-triggers"
}

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	description = "test job"
	schedule = "*/2 * * * *"

	pubsub_target = {
		topic_name = "projects/%s/topics/build-triggers"
		data = "${base64encode("test")}"
	}
}

	`, name, project)
}

func testAccCloudSchedulerJob_appEngineConfig(name string) string {
	return fmt.Sprintf(`

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	schedule = "*/4 * * * *"
	description = "test app engine job"
	time_zone = "Europe/London"

	app_engine_http_target = {
		http_method = "POST"
    app_engine_routing = {
      service = "web"
      version = "prod"
      instance = "my-instance-001"
    }
		relative_uri = "/ping"
	}
}

	`, name)
}

func testAccCloudSchedulerJob_httpConfig(name string) string {
	return fmt.Sprintf(`

resource "google_cloud_scheduler_job" "job" {
	name     = "%s"
	description = "test http job"
	schedule = "*/8 * * * *"
	time_zone = "America/New_York"

	http_target = {
		http_method = "POST"
		uri = "https://example.com/ping"
	}
}

	`, name)
}
