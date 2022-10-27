package main

import (
	"reflect"
	"sort"
	"testing"

	provider "google/provider/new/google-beta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDetectMissingTest(t *testing.T) {
	for _, test := range []struct {
		resourceName           string
		changedFields          []string
		expectedUntestedFields []string
		errorExpected          bool
	}{
		{
			resourceName: "google_vertex_ai_endpoint",
			changedFields: allFields(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"display_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `Required. The display name of the Endpoint. The name can be up to 128 characters long and can consist of any UTF-8 characters.`,
					},
					"location": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: `The location for the resource`,
					},
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: `The resource name of the Endpoint. The name must be numeric with no leading zeros and can be at most 10 digits.`,
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `The description of the Endpoint.`,
					},
					"encryption_spec": {
						Type:        schema.TypeList,
						Optional:    true,
						ForceNew:    true,
						Description: `Customer-managed encryption key spec for an Endpoint. If set, this Endpoint and all sub-resources of this Endpoint will be secured by this key.`,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kms_key_name": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    true,
									Description: `Required. The Cloud KMS resource identifier of the customer managed encryption key used to protect a resource. Has the form: 'projects/my-project/locations/my-region/keyRings/my-kr/cryptoKeys/my-key'. The key needs to be in the same region as where the compute resource is created.`,
								},
							},
						},
					},
					"labels": {
						Type:        schema.TypeMap,
						Optional:    true,
						Description: `The labels with user-defined metadata to organize your Endpoints. Label keys and values can be no longer than 64 characters (Unicode codepoints), can only contain lowercase letters, numeric characters, underscores and dashes. International characters are allowed. See https://goo.gl/xmQnxf for more information and examples of labels.`,
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					"network": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: `The full name of the Google Compute Engine [network](https://cloud.google.com//compute/docs/networks-and-firewalls#networks) to which the Endpoint should be peered. Private services access must already be configured for the network. If left unspecified, the Endpoint is not peered with any network. Only one of the fields, network or enable_private_service_connect, can be set. [Format](https://cloud.google.com/compute/docs/reference/rest/v1/networks/insert): 'projects/{project}/global/networks/{network}'. Where '{project}' is a project number, as in '12345', and '{network}' is network name.`,
					},
					"create_time": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Output only. Timestamp when this Endpoint was created.`,
					},
					"deployed_models": {
						Type:        schema.TypeList,
						Computed:    true,
						Description: `Output only. The models deployed in this Endpoint. To add or remove DeployedModels use EndpointService.DeployModel and EndpointService.UndeployModel respectively. Models can also be deployed and undeployed using the [Cloud Console](https://console.cloud.google.com/vertex-ai/).`,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"automatic_resources": {
									Type:        schema.TypeList,
									Computed:    true,
									Description: `A description of resources that to large degree are decided by Vertex AI, and require only a modest additional configuration.`,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"max_replica_count": {
												Type:        schema.TypeInt,
												Computed:    true,
												Description: `The maximum number of replicas this DeployedModel may be deployed on when the traffic against it increases. If the requested value is too large, the deployment will error, but if deployment succeeds then the ability to scale the model to that many replicas is guaranteed (barring service outages). If traffic against the DeployedModel increases beyond what its replicas at maximum may handle, a portion of the traffic will be dropped. If this value is not provided, a no upper bound for scaling under heavy traffic will be assume, though Vertex AI may be unable to scale beyond certain replica number.`,
											},
											"min_replica_count": {
												Type:        schema.TypeInt,
												Computed:    true,
												Description: `The minimum number of replicas this DeployedModel will be always deployed on. If traffic against it increases, it may dynamically be deployed onto more replicas up to max_replica_count, and as traffic decreases, some of these extra replicas may be freed. If the requested value is too large, the deployment will error.`,
											},
										},
									},
								},
								"create_time": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `Output only. Timestamp when the DeployedModel was created.`,
								},
								"dedicated_resources": {
									Type:        schema.TypeList,
									Computed:    true,
									Description: `A description of resources that are dedicated to the DeployedModel, and that need a higher degree of manual configuration.`,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"autoscaling_metric_specs": {
												Type:        schema.TypeList,
												Computed:    true,
												Description: `The metric specifications that overrides a resource utilization metric (CPU utilization, accelerator's duty cycle, and so on) target value (default to 60 if not set). At most one entry is allowed per metric. If machine_spec.accelerator_count is above 0, the autoscaling will be based on both CPU utilization and accelerator's duty cycle metrics and scale up when either metrics exceeds its target value while scale down if both metrics are under their target value. The default target value is 60 for both metrics. If machine_spec.accelerator_count is 0, the autoscaling will be based on CPU utilization metric only with default target value 60 if not explicitly set. For example, in the case of Online Prediction, if you want to override target CPU utilization to 80, you should set autoscaling_metric_specs.metric_name to 'aiplatform.googleapis.com/prediction/online/cpu/utilization' and autoscaling_metric_specs.target to '80'.`,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"metric_name": {
															Type:        schema.TypeString,
															Computed:    true,
															Description: `The resource metric name. Supported metrics: * For Online Prediction: * 'aiplatform.googleapis.com/prediction/online/accelerator/duty_cycle' * 'aiplatform.googleapis.com/prediction/online/cpu/utilization'`,
														},
														"target": {
															Type:        schema.TypeInt,
															Computed:    true,
															Description: `The target resource utilization in percentage (1% - 100%) for the given metric; once the real usage deviates from the target by a certain percentage, the machine replicas change. The default value is 60 (representing 60%) if not provided.`,
														},
													},
												},
											},
											"machine_spec": {
												Type:        schema.TypeList,
												Computed:    true,
												Description: `The specification of a single machine used by the prediction.`,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"accelerator_count": {
															Type:        schema.TypeInt,
															Computed:    true,
															Description: `The number of accelerators to attach to the machine.`,
														},
														"accelerator_type": {
															Type:        schema.TypeString,
															Computed:    true,
															Description: `The type of accelerator(s) that may be attached to the machine as per accelerator_count. See possible values [here](https://cloud.google.com/vertex-ai/docs/reference/rest/v1/MachineSpec#AcceleratorType).`,
														},
														"machine_type": {
															Type:        schema.TypeString,
															Computed:    true,
															Description: `The type of the machine. See the [list of machine types supported for prediction](https://cloud.google.com/vertex-ai/docs/predictions/configure-compute#machine-types) See the [list of machine types supported for custom training](https://cloud.google.com/vertex-ai/docs/training/configure-compute#machine-types). For DeployedModel this field is optional, and the default value is 'n1-standard-2'. For BatchPredictionJob or as part of WorkerPoolSpec this field is required. TODO(rsurowka): Try to better unify the required vs optional.`,
														},
													},
												},
											},
											"max_replica_count": {
												Type:        schema.TypeInt,
												Computed:    true,
												Description: `The maximum number of replicas this DeployedModel may be deployed on when the traffic against it increases. If the requested value is too large, the deployment will error, but if deployment succeeds then the ability to scale the model to that many replicas is guaranteed (barring service outages). If traffic against the DeployedModel increases beyond what its replicas at maximum may handle, a portion of the traffic will be dropped. If this value is not provided, will use min_replica_count as the default value. The value of this field impacts the charge against Vertex CPU and GPU quotas. Specifically, you will be charged for max_replica_count * number of cores in the selected machine type) and (max_replica_count * number of GPUs per replica in the selected machine type).`,
											},
											"min_replica_count": {
												Type:        schema.TypeInt,
												Computed:    true,
												Description: `The minimum number of machine replicas this DeployedModel will be always deployed on. This value must be greater than or equal to 1. If traffic against the DeployedModel increases, it may dynamically be deployed onto more replicas, and as traffic decreases, some of these extra replicas may be freed.`,
											},
										},
									},
								},
								"display_name": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `The display name of the DeployedModel. If not provided upon creation, the Model's display_name is used.`,
								},
								"enable_access_logging": {
									Type:        schema.TypeBool,
									Computed:    true,
									Description: `These logs are like standard server access logs, containing information like timestamp and latency for each prediction request. Note that Stackdriver logs may incur a cost, especially if your project receives prediction requests at a high queries per second rate (QPS). Estimate your costs before enabling this option.`,
								},
								"enable_container_logging": {
									Type:        schema.TypeBool,
									Computed:    true,
									Description: `If true, the container of the DeployedModel instances will send 'stderr' and 'stdout' streams to Stackdriver Logging. Only supported for custom-trained Models and AutoML Tabular Models.`,
								},
								"id": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `The ID of the DeployedModel. If not provided upon deployment, Vertex AI will generate a value for this ID. This value should be 1-10 characters, and valid characters are /[0-9]/.`,
								},
								"model": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `The name of the Model that this is the deployment of. Note that the Model may be in a different location than the DeployedModel's Endpoint.`,
								},
								"model_version_id": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `Output only. The version ID of the model that is deployed.`,
								},
								"private_endpoints": {
									Type:        schema.TypeList,
									Computed:    true,
									Description: `Output only. Provide paths for users to send predict/explain/health requests directly to the deployed model services running on Cloud via private services access. This field is populated if network is configured.`,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"explain_http_uri": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `Output only. Http(s) path to send explain requests.`,
											},
											"health_http_uri": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `Output only. Http(s) path to send health check requests.`,
											},
											"predict_http_uri": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `Output only. Http(s) path to send prediction requests.`,
											},
											"service_attachment": {
												Type:        schema.TypeString,
												Computed:    true,
												Description: `Output only. The name of the service attachment resource. Populated if private service connect is enabled.`,
											},
										},
									},
								},
								"service_account": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `The service account that the DeployedModel's container runs as. Specify the email address of the service account. If this service account is not specified, the container runs as a service account that doesn't have access to the resource project. Users deploying the Model must have the 'iam.serviceAccounts.actAs' permission on this service account.`,
								},
								"shared_resources": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: `The resource name of the shared DeploymentResourcePool to deploy on. Format: projects/{project}/locations/{location}/deploymentResourcePools/{deployment_resource_pool}`,
								},
							},
						},
					},
					"etag": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Used to perform consistent read-modify-write updates. If not set, a blind "overwrite" update happens.`,
					},
					"model_deployment_monitoring_job": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Output only. Resource name of the Model Monitoring job associated with this Endpoint if monitoring is enabled by CreateModelDeploymentMonitoringJob. Format: 'projects/{project}/locations/{location}/modelDeploymentMonitoringJobs/{model_deployment_monitoring_job}'`,
					},
					"update_time": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: `Output only. Timestamp when this Endpoint was last updated.`,
					},
					"project": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			}, nil),
		},
		{
			resourceName:  "google_compute_instance",
			changedFields: allFields(provider.ResourceMap()["google_compute_instance"], nil),
			expectedUntestedFields: []string{
				"attached_disk.device_name",
				"attached_disk.kms_key_self_link",
				"boot_disk.auto_delete",
				"boot_disk.device_name",
				"boot_disk.initialize_params.labels",
				"boot_disk.initialize_params.size",
				"description",
				"metadata_startup_script",
			},
			errorExpected: true,
		},
	} {
		if missingTest, err := detectMissingTest(test.resourceName, "testdata", test.changedFields); err != nil && !test.errorExpected {
			t.Errorf("error detecting missing test for resource %s: %v", test.resourceName, err)
		} else if missingTest != nil {
			sort.Strings(missingTest.UntestedFields)
			sort.Strings(test.expectedUntestedFields)
			if !reflect.DeepEqual(missingTest.UntestedFields, test.expectedUntestedFields) {
				t.Errorf("found unexpected untested fields: %v, expected %v", missingTest.UntestedFields, test.expectedUntestedFields)
			}
		} else if len(test.expectedUntestedFields) > 0 {
			t.Errorf("failed to find expected untested fields: %v", test.expectedUntestedFields)
		}
	}
}
