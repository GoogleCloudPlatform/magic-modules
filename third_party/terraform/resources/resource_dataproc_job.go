package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/dataproc/v1"
)

func resourceDataprocJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocJobCreate,
		Update: resourceDataprocJobUpdate,
		Read:   resourceDataprocJobRead,
		Delete: resourceDataprocJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project in which the cluster can be found and jobs subsequently run against. If it is not provided, the provider project is used.`,
			},

			// Ref: https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.jobs#JobReference
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				ForceNew:    true,
				Description: `The Cloud Dataproc region. This essentially determines which clusters are available for this job to be submitted to. If not specified, defaults to global.`,
			},

			// If a job is still running, trying to delete a job will fail. Setting
			// this flag to true however will force the deletion by first cancelling
			// the job and then deleting it
			"force_delete": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: `By default, you can only delete inactive jobs within Dataproc. Setting this to true, and calling destroy, will ensure that the job is first cancelled before issuing the delete.`,
			},

			"reference": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `The reference of the job`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:         schema.TypeString,
							Description:  "The job ID, which must be unique within the project. The job ID is generated by the server upon job submission or provided by the user as a means to perform retries without creating duplicate jobs",
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validateRegexp("^[a-zA-Z0-9_-]{1,100}$"),
						},
					},
				},
			},

			"placement": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `The config of job placement.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:        schema.TypeString,
							Description: "The name of the cluster where the job will be submitted",
							Required:    true,
							ForceNew:    true,
						},
						"cluster_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Output-only. A cluster UUID generated by the Cloud Dataproc service when the job is submitted",
						},
					},
				},
			},

			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The status of the job.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:        schema.TypeString,
							Description: "Output-only. A state message specifying the overall job state",
							Computed:    true,
						},
						"details": {
							Type:        schema.TypeString,
							Description: "Output-only. Optional job state details, such as an error description if the state is ERROR",
							Computed:    true,
						},
						"state_start_time": {
							Type:        schema.TypeString,
							Description: "Output-only. The time when this state was entered",
							Computed:    true,
						},
						"substate": {
							Type:        schema.TypeString,
							Description: "Output-only. Additional state information, which includes status reported by the agent",
							Computed:    true,
						},
					},
				},
			},

			"driver_output_resource_uri": {
				Type:        schema.TypeString,
				Description: "Output-only. A URI pointing to the location of the stdout of the job's driver program",
				Computed:    true,
			},

			"driver_controls_files_uri": {
				Type:        schema.TypeString,
				Description: "Output-only. If present, the location of miscellaneous control files which may be used as part of job setup and handling. If not present, control files may be placed in the same location as driver_output_uri.",
				Computed:    true,
			},

			"labels": {
				Type:        schema.TypeMap,
				Description: "Optional. The labels to associate with this job.",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"scheduling": {
				Type:        schema.TypeList,
				Description: "Optional. Job scheduling configuration.",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_failures_per_hour": {
							Type:         schema.TypeInt,
							Description:  "Maximum number of times per hour a driver may be restarted as a result of driver terminating with non-zero code before job is reported failed.",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtMost(10),
						},
					},
				},
			},

			"pyspark_config":  pySparkSchema,
			"spark_config":    sparkSchema,
			"hadoop_config":   hadoopSchema,
			"hive_config":     hiveSchema,
			"pig_config":      pigSchema,
			"sparksql_config": sparkSqlSchema,
		},
	}
}

func resourceDataprocJobUpdate(d *schema.ResourceData, meta interface{}) error {
	// The only updatable value is currently 'force_delete' which is a local
	// only value therefore we don't need to make any GCP calls to update this.

	return resourceDataprocJobRead(d, meta)
}

func resourceDataprocJobCreate(d *schema.ResourceData, meta interface{}) error {
	var m providerMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return err
	}
	config := meta.(*Config)
	config.clientDataproc.UserAgent = fmt.Sprintf("%s %s", config.clientDataproc.UserAgent, m.ModuleKey)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	clusterName := d.Get("placement.0.cluster_name").(string)
	region := d.Get("region").(string)

	submitReq := &dataproc.SubmitJobRequest{
		Job: &dataproc.Job{
			Placement: &dataproc.JobPlacement{
				ClusterName: clusterName,
			},
			Reference: &dataproc.JobReference{
				ProjectId: project,
			},
		},
	}

	if v, ok := d.GetOk("reference.0.job_id"); ok {
		submitReq.Job.Reference.JobId = v.(string)
	}
	if _, ok := d.GetOk("labels"); ok {
		submitReq.Job.Labels = expandLabels(d)
	}

	if v, ok := d.GetOk("pyspark_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.PysparkJob = expandPySparkJob(config)
	}

	if v, ok := d.GetOk("spark_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.SparkJob = expandSparkJob(config)
	}

	if v, ok := d.GetOk("hadoop_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.HadoopJob = expandHadoopJob(config)
	}

	if v, ok := d.GetOk("hive_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.HiveJob = expandHiveJob(config)
	}

	if v, ok := d.GetOk("pig_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.PigJob = expandPigJob(config)
	}

	if v, ok := d.GetOk("sparksql_config"); ok {
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.SparkSqlJob = expandSparkSqlJob(config)
	}

	// Submit the job
	job, err := config.clientDataproc.Projects.Regions.Jobs.Submit(
		project, region, submitReq).Do()
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/jobs/%s", project, region, job.Reference.JobId))

	waitErr := dataprocJobOperationWait(config, region, project, job.Reference.JobId,
		"Creating Dataproc job", d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc job %s has been submitted", job.Reference.JobId)
	return resourceDataprocJobRead(d, meta)
}

func resourceDataprocJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := d.Get("region").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	parts := strings.Split(d.Id(), "/")
	jobId := parts[len(parts)-1]
	job, err := config.clientDataproc.Projects.Regions.Jobs.Get(
		project, region, jobId).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Job %q", jobId))
	}

	d.Set("force_delete", d.Get("force_delete"))
	d.Set("labels", job.Labels)
	d.Set("driver_output_resource_uri", job.DriverOutputResourceUri)
	d.Set("driver_controls_files_uri", job.DriverControlFilesUri)

	d.Set("placement", flattenJobPlacement(job.Placement))
	d.Set("status", flattenJobStatus(job.Status))
	d.Set("reference", flattenJobReference(job.Reference))
	d.Set("project", project)

	if job.PysparkJob != nil {
		d.Set("pyspark_config", flattenPySparkJob(job.PysparkJob))
	}
	if job.SparkJob != nil {
		d.Set("spark_config", flattenSparkJob(job.SparkJob))
	}
	if job.HadoopJob != nil {
		d.Set("hadoop_config", flattenHadoopJob(job.HadoopJob))
	}
	if job.HiveJob != nil {
		d.Set("hive_config", flattenHiveJob(job.HiveJob))
	}
	if job.PigJob != nil {
		d.Set("pig_config", flattenPigJob(job.PigJob))
	}
	if job.SparkSqlJob != nil {
		d.Set("sparksql_config", flattenSparkSqlJob(job.SparkSqlJob))
	}
	return nil
}

func resourceDataprocJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	forceDelete := d.Get("force_delete").(bool)

	parts := strings.Split(d.Id(), "/")
	jobId := parts[len(parts)-1]
	if forceDelete {
		log.Printf("[DEBUG] Attempting to first cancel Dataproc job %s if it's still running ...", d.Id())

		// ignore error if we get one - job may be finished already and not need to
		// be cancelled. We do however wait for the state to be one that is
		// at least not active
		_, _ = config.clientDataproc.Projects.Regions.Jobs.Cancel(project, region, jobId, &dataproc.CancelJobRequest{}).Do()

		waitErr := dataprocJobOperationWait(config, region, project, jobId,
			"Cancelling Dataproc job", d.Timeout(schema.TimeoutDelete))
		if waitErr != nil {
			return waitErr
		}

	}

	log.Printf("[DEBUG] Deleting Dataproc job %s", d.Id())
	_, err = config.clientDataproc.Projects.Regions.Jobs.Delete(
		project, region, jobId).Do()
	if err != nil {
		return err
	}

	waitErr := dataprocDeleteOperationWait(config, region, project, jobId,
		"Deleting Dataproc job", d.Timeout(schema.TimeoutDelete))
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc job %s has been deleted", d.Id())
	d.SetId("")

	return nil
}

// ---- PySpark Job ----

var loggingConfig = &schema.Schema{
	Type:        schema.TypeList,
	Description: "The runtime logging config of the job",
	Optional:    true,
	Computed:    true,
	MaxItems:    1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver_log_levels": {
				Type:        schema.TypeMap,
				Description: "Optional. The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'.",
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

var pySparkSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of pySpark job.`,
	ExactlyOneOf: []string{"pyspark_config", "spark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"main_python_file_uri": {
				Type:        schema.TypeString,
				Description: "Required. The HCFS URI of the main Python file to use as the driver. Must be a .py file",
				Required:    true,
				ForceNew:    true,
			},

			"args": {
				Type:        schema.TypeList,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as --conf, that can be set as job properties, since a collision may occur that causes an incorrect job submission",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"python_file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of files to be copied to the working directory of Python drivers and distributed tasks. Useful for naively parallel tasks",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Description: "Optional. A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig,
		},
	},
}

func flattenPySparkJob(job *dataproc.PySparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_python_file_uri": job.MainPythonFileUri,
			"args":                 job.Args,
			"python_file_uris":     job.PythonFileUris,
			"jar_file_uris":        job.JarFileUris,
			"file_uris":            job.FileUris,
			"archive_uris":         job.ArchiveUris,
			"properties":           job.Properties,
			"logging_config":       flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandPySparkJob(config map[string]interface{}) *dataproc.PySparkJob {
	job := &dataproc.PySparkJob{}
	if v, ok := config["main_python_file_uri"]; ok {
		job.MainPythonFileUri = v.(string)
	}
	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["python_file_uris"]; ok {
		job.PythonFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.FileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.ArchiveUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Spark Job ----

var sparkSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of the Spark job.`,
	ExactlyOneOf: []string{"pyspark_config", "spark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// main driver: can be only one of the class | jar_file
			"main_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The class containing the main method of the driver. Must be in a provided jar or jar that is already on the classpath. Conflicts with main_jar_file_uri`,
				ExactlyOneOf: []string{"spark_config.0.main_class", "spark_config.0.main_jar_file_uri"},
			},

			"main_jar_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The HCFS URI of jar file containing the driver jar. Conflicts with main_class`,
				ExactlyOneOf: []string{"spark_config.0.main_jar_file_uri", "spark_config.0.main_class"},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The arguments to pass to the driver.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of files to be copied to the working directory of Spark drivers and distributed tasks. Useful for naively parallel tasks.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig,
		},
	},
}

func flattenSparkJob(job *dataproc.SparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_class":        job.MainClass,
			"main_jar_file_uri": job.MainJarFileUri,
			"args":              job.Args,
			"jar_file_uris":     job.JarFileUris,
			"file_uris":         job.FileUris,
			"archive_uris":      job.ArchiveUris,
			"properties":        job.Properties,
			"logging_config":    flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandSparkJob(config map[string]interface{}) *dataproc.SparkJob {
	job := &dataproc.SparkJob{}
	if v, ok := config["main_class"]; ok {
		job.MainClass = v.(string)
	}
	if v, ok := config["main_jar_file_uri"]; ok {
		job.MainJarFileUri = v.(string)
	}

	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.FileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.ArchiveUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Hadoop Job ----

var hadoopSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of Hadoop job`,
	ExactlyOneOf: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// main driver: can be only one of the main_class | main_jar_file_uri
			"main_class": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The class containing the main method of the driver. Must be in a provided jar or jar that is already on the classpath. Conflicts with main_jar_file_uri`,
				ExactlyOneOf: []string{"hadoop_config.0.main_jar_file_uri", "hadoop_config.0.main_class"},
			},

			"main_jar_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The HCFS URI of jar file containing the driver jar. Conflicts with main_class`,
				ExactlyOneOf: []string{"hadoop_config.0.main_jar_file_uri", "hadoop_config.0.main_class"},
			},

			"args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `The arguments to pass to the driver.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of jar files to add to the CLASSPATHs of the Spark driver and tasks.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of files to be copied to the working directory of Spark drivers and distributed tasks. Useful for naively parallel tasks.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `A mapping of property names to values, used to configure Spark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig,
		},
	},
}

func flattenHadoopJob(job *dataproc.HadoopJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_class":        job.MainClass,
			"main_jar_file_uri": job.MainJarFileUri,
			"args":              job.Args,
			"jar_file_uris":     job.JarFileUris,
			"file_uris":         job.FileUris,
			"archive_uris":      job.ArchiveUris,
			"properties":        job.Properties,
			"logging_config":    flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandHadoopJob(config map[string]interface{}) *dataproc.HadoopJob {
	job := &dataproc.HadoopJob{}
	if v, ok := config["main_class"]; ok {
		job.MainClass = v.(string)
	}
	if v, ok := config["main_jar_file_uri"]; ok {
		job.MainJarFileUri = v.(string)
	}

	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.FileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.ArchiveUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Hive Job ----

var hiveSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of hive job`,
	ExactlyOneOf: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Description:  `The list of Hive queries or statements to execute as part of the job. Conflicts with query_file_uri`,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"hive_config.0.query_file_uri", "hive_config.0.query_list"},
			},

			"query_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `HCFS URI of file containing Hive script to execute as the job. Conflicts with query_list`,
				ExactlyOneOf: []string{"hive_config.0.query_file_uri", "hive_config.0.query_list"},
			},

			"continue_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to continue executing queries if a query fails. The default value is false. Setting to true can be useful when executing independent parallel queries. Defaults to false.`,
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `Mapping of query variable names to values (equivalent to the Hive command: SET name="value";).`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `A mapping of property names and values, used to configure Hive. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/hive/conf/hive-site.xml, and classes in user code.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of jar files to add to the CLASSPATH of the Hive server and Hadoop MapReduce (MR) tasks. Can contain Hive SerDes and UDFs.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func flattenHiveJob(job *dataproc.HiveJob) []map[string]interface{} {
	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":          queries,
			"query_file_uri":      job.QueryFileUri,
			"continue_on_failure": job.ContinueOnFailure,
			"script_variables":    job.ScriptVariables,
			"properties":          job.Properties,
			"jar_file_uris":       job.JarFileUris,
		},
	}
}

func expandHiveJob(config map[string]interface{}) *dataproc.HiveJob {
	job := &dataproc.HiveJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["continue_on_failure"]; ok {
		job.ContinueOnFailure = v.(bool)
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job
}

// ---- Pig Job ----

var pigSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of pag job.`,
	ExactlyOneOf: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Description:  `The list of Hive queries or statements to execute as part of the job. Conflicts with query_file_uri`,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"pig_config.0.query_file_uri", "pig_config.0.query_list"},
			},

			"query_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `HCFS URI of file containing Hive script to execute as the job. Conflicts with query_list`,
				ExactlyOneOf: []string{"pig_config.0.query_file_uri", "pig_config.0.query_list"},
			},

			"continue_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to continue executing queries if a query fails. The default value is false. Setting to true can be useful when executing independent parallel queries. Defaults to false.`,
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `Mapping of query variable names to values (equivalent to the Pig command: name=[value]).`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `A mapping of property names to values, used to configure Pig. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/hadoop/conf/*-site.xml, /etc/pig/conf/pig.properties, and classes in user code.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of jar files to add to the CLASSPATH of the Pig Client and Hadoop MapReduce (MR) tasks. Can contain Pig UDFs.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig,
		},
	},
}

func flattenPigJob(job *dataproc.PigJob) []map[string]interface{} {
	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":          queries,
			"query_file_uri":      job.QueryFileUri,
			"continue_on_failure": job.ContinueOnFailure,
			"script_variables":    job.ScriptVariables,
			"properties":          job.Properties,
			"jar_file_uris":       job.JarFileUris,
		},
	}
}

func expandPigJob(config map[string]interface{}) *dataproc.PigJob {
	job := &dataproc.PigJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["continue_on_failure"]; ok {
		job.ContinueOnFailure = v.(bool)
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job

}

// ---- Spark SQL Job ----

var sparkSqlSchema = &schema.Schema{
	Type:         schema.TypeList,
	Optional:     true,
	ForceNew:     true,
	MaxItems:     1,
	Description:  `The config of SparkSql job`,
	ExactlyOneOf: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Description:  `The list of SQL queries or statements to execute as part of the job. Conflicts with query_file_uri`,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"sparksql_config.0.query_file_uri", "sparksql_config.0.query_list"},
			},

			"query_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `The HCFS URI of the script that contains SQL queries. Conflicts with query_list`,
				ExactlyOneOf: []string{"sparksql_config.0.query_file_uri", "sparksql_config.0.query_list"},
			},

			"script_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `Mapping of query variable names to values (equivalent to the Spark SQL command: SET name="value";).`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: `A mapping of property names to values, used to configure Spark SQL's SparkConf. Properties that conflict with values set by the Cloud Dataproc API may be overwritten.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `HCFS URIs of jar files to be added to the Spark CLASSPATH.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig,
		},
	},
}

func flattenSparkSqlJob(job *dataproc.SparkSqlJob) []map[string]interface{} {
	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":       queries,
			"query_file_uri":   job.QueryFileUri,
			"script_variables": job.ScriptVariables,
			"properties":       job.Properties,
			"jar_file_uris":    job.JarFileUris,
		},
	}
}

func expandSparkSqlJob(config map[string]interface{}) *dataproc.SparkSqlJob {
	job := &dataproc.SparkSqlJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job

}

// ---- Other flatten / expand methods ----

func expandLoggingConfig(config map[string]interface{}) *dataproc.LoggingConfig {
	conf := &dataproc.LoggingConfig{}
	if v, ok := config["driver_log_levels"]; ok {
		conf.DriverLogLevels = convertStringMap(v.(map[string]interface{}))
	}
	return conf
}

func flattenLoggingConfig(l *dataproc.LoggingConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"driver_log_levels": l.DriverLogLevels,
		},
	}
}

func flattenJobReference(r *dataproc.JobReference) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"job_id": r.JobId,
		},
	}
}

func flattenJobStatus(s *dataproc.JobStatus) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"state":            s.State,
			"details":          s.Details,
			"state_start_time": s.StateStartTime,
			"substate":         s.Substate,
		},
	}
}

func flattenJobPlacement(jp *dataproc.JobPlacement) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"cluster_name": jp.ClusterName,
			"cluster_uuid": jp.ClusterUuid,
		},
	}
}
