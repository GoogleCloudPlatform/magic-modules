package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/googleapi"
)

func resourceDataflowFlexTemplateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataflowFlexTemplateJobLaunchTemplate,
		Read:   resourceDataflowFlexTemplateJobRead,
		// Update: resourceDataflowFlexTemplateJobUpdateByReplacement,
		Delete: resourceDataflowFlexTemplateJobDelete,
		// CustomizeDiff: customdiff.All(
		// 	resourceDataflowJobTypeCustomizeDiff,
		// ),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_spec_gcs_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceDataflowFlexTemplateJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	request, err := buildLaunchFlexTemplateRequest(d, config)

	job, err := resourceDataflowFlexTemplateJobCreateJob(config, project, region, &request)
	if err != nil {
		return err
	}
	d.SetId(job.Id)

	return resourceDataflowFlexTemplateJobRead(d, meta)
}

func resourceDataflowFlexTemplateJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	id := d.Id()

	job, err := resourceDataflowFlexTemplateJobGetJob(config, project, region, id)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataflow job %s", id))
	}

	d.Set("job_id", job.Id)
	d.Set("state", job.CurrentState)
	d.Set("name", job.Name)
	d.Set("type", job.Type)
	d.Set("project", project)
	d.Set("labels", job.Labels)

	sdkPipelineOptions, err := ConvertToMap(job.Environment.SdkPipelineOptions)
	if err != nil {
		return err
	}
	optionsMap := sdkPipelineOptions["options"].(map[string]interface{})
	d.Set("template_gcs_path", optionsMap["templateLocation"])
	d.Set("temp_gcs_location", optionsMap["tempLocation"])

	if _, ok := dataflowTerminalStatesMap[job.CurrentState]; ok {
		log.Printf("[DEBUG] Removing resource '%s' because it is in state %s.\n", job.Name, job.CurrentState)
		d.SetId("")
		return nil
	}
	d.SetId(job.Id)

	return nil
}

// Stream update method. Batch job changes should have been set to ForceNew via custom diff
func resourceDataflowFlexTemplateJobUpdateByReplacement(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	params := expandStringMap(d, "parameters")

	env, err := resourceDataflowJobSetupEnv(d, config)
	if err != nil {
		return err
	}

	request := dataflow.LaunchTemplateParameters{
		JobName:     d.Get("name").(string),
		Parameters:  params,
		Environment: &env,
		Update:      true,
	}

	var response *dataflow.LaunchTemplateResponse
	err = retryTimeDuration(func() (updateErr error) {
		response, updateErr = resourceDataflowJobLaunchTemplate(config, project, region, d.Get("template_gcs_path").(string), &request)
		return updateErr
	}, time.Minute*time.Duration(5), isDataflowJobUpdateRetryableError)
	if err != nil {
		return err
	}
	d.SetId(response.Job.Id)

	return resourceDataflowJobRead(d, meta)
}

func resourceDataflowFlexTemplateJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	id := d.Id()

	requestedState, err := resourceDataflowJobMapRequestedState(d.Get("on_delete").(string))
	if err != nil {
		return err
	}

	// Retry updating the state while the job is not ready to be canceled/drained.
	err = resource.Retry(time.Minute*time.Duration(15), func() *resource.RetryError {
		// To terminate a dataflow job, we update the job with a requested
		// terminal state.
		job := &dataflow.Job{
			RequestedState: requestedState,
		}

		_, updateErr := resourceDataflowJobUpdateJob(config, project, region, id, job)
		if updateErr != nil {
			gerr, isGoogleErr := updateErr.(*googleapi.Error)
			if !isGoogleErr {
				// If we have an error and it's not a google-specific error, we should go ahead and return.
				return resource.NonRetryableError(updateErr)
			}

			if strings.Contains(gerr.Message, "not yet ready for canceling") {
				// Retry cancelling job if it's not ready.
				// Sleep to avoid hitting update quota with repeated attempts.
				time.Sleep(5 * time.Second)
				return resource.RetryableError(updateErr)
			}

			if strings.Contains(gerr.Message, "Job has terminated") {
				// Job has already been terminated, skip.
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Wait for state to reach terminal state (canceled/drained/done)
	_, ok := dataflowTerminalStatesMap[d.Get("state").(string)]
	for !ok {
		log.Printf("[DEBUG] Waiting for job with job state %q to terminate...", d.Get("state").(string))
		time.Sleep(5 * time.Second)

		err = resourceDataflowJobRead(d, meta)
		if err != nil {
			return fmt.Errorf("Error while reading job to see if it was properly terminated: %v", err)
		}
		_, ok = dataflowTerminalStatesMap[d.Get("state").(string)]
	}

	// Only remove the job from state if it's actually successfully canceled.
	if _, ok := dataflowTerminalStatesMap[d.Get("state").(string)]; ok {
		log.Printf("[DEBUG] Removing dataflow job with final state %q", d.Get("state").(string))
		d.SetId("")
		return nil
	}
	return fmt.Errorf("Unable to cancel the dataflow job '%s' - final state was %q.", d.Id(), d.Get("state").(string))
}

func resourceDataflowFlexTemplateJobMapRequestedState(policy string) (string, error) {
	switch policy {
	case "cancel":
		return "JOB_STATE_CANCELLED", nil
	case "drain":
		return "JOB_STATE_DRAINING", nil
	default:
		return "", fmt.Errorf("Invalid `on_delete` policy: %s", policy)
	}
}

func resourceDataflowFlexTemplateJobLaunchTemplate(config *Config, project string, region string, request *dataflow.LaunchFlexTemplateRequest) (*dataflow.Job, error) {
	response, err := config.clientDataflow.Projects.Locations.FlexTemplates.Launch(project, region, request).Do()
	if err != nil {
		return nil, err
	}

	return response.Job, nil
}

func resourceDataflowFlexTemplateJobGetJob(config *Config, project string, region string, id string) (*dataflow.Job, error) {
	return config.clientDataflow.Projects.Locations.Jobs.Get(project, region, id).View("JOB_VIEW_ALL").Do()
}

// Streaming jobs cannot be updated:
// https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates#limitations
// func resourceDataflowFlexTemplateJobUpdateJob(config *Config, project string, region string, id string, job *dataflow.Job) (*dataflow.Job, error) {
// 	if region == "" {
// 		return config.clientDataflow.Projects.Jobs.Update(project, id, job).Do()
// 	}
// 	return config.clientDataflow.Projects.Locations.Jobs.Update(project, region, id, job).Do()
// }

func buildLaunchFlexTemplateRequest(d *schema.ResourceData, config *Config) dataflow.LaunchFlexTemplateRequest {
	return dataflow.LaunchFlexTemplateRequest{
		LaunchParameter: &dataflow.LaunchFlexTemplateParameter{
			ContainerSpecGcsPath: d.Get("container_spec_gcs_path").(string),
			JobName:              d.Get("name").(string),
			Parameters:           convertStringMap(d, "parameters"),
		},
	}
}
