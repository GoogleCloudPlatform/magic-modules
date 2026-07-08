package datalineage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/googleapi"
)

var (
	_ = bytes.Clone
	_ = context.WithCancel
	_ = base64.NewDecoder
	_ = json.Marshal
	_ = fmt.Sprintf
	_ = log.Print
	_ = http.Get
	_ = reflect.ValueOf
	_ = regexp.Match
	_ = slices.Min([]int{1})
	_ = sort.IntSlice{}
	_ = strconv.Atoi
	_ = strings.Trim
	_ = time.Now
	_ = errwrap.Wrap
	_ = cty.BoolVal
	_ = diag.Diagnostic{}
	_ = customdiff.All
	_ = id.UniqueId
	_ = logging.LogLevel
	_ = retry.Retry
	_ = schema.Noop
	_ = validation.All
	_ = structure.ExpandJsonFromString
	_ = terraform.State{}
	_ = tpgresource.SetLabels
	_ = transport_tpg.Config{}
	_ = verify.ValidateEnum
	_ = googleapi.Error{}
)

func init() {
	registry.Schema{
		Name:        "google_data_lineage_open_lineage_job",
		ProductName: "datalineage",
		Type:        registry.SchemaTypeResource,
		Schema:      ResourceDataLineageOpenLineageJob(),
	}.Register()
}

func ResourceDataLineageOpenLineageJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataLineageOpenLineageJobCreate,
		ReadContext:   resourceDataLineageOpenLineageJobRead,
		UpdateContext: resourceDataLineageOpenLineageJobUpdate,
		DeleteContext: resourceDataLineageOpenLineageJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		ResourceBehavior: schema.ResourceBehavior{
			MutableIdentity: true,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Name of the OpenLineage job.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Namespace of the OpenLineage job.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description of the OpenLineage job.`,
			},
			"input": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Input datasets consumed by this job.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Name of the dataset.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Namespace of the dataset.`,
						},
						"catalog": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Catalog information for the dataset.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"framework": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog framework.`,
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog entity name.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog entity type.`,
									},
								},
							},
						},
						"symlink": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Symlink targets for the dataset.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Name of the symlink target.`,
									},
									"namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Namespace of the symlink target.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Type of the symlink target.`,
									},
								},
							},
						},
					},
				},
			},
			"output": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Output datasets produced by this job.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Name of the dataset.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Namespace of the dataset.`,
						},
						"catalog": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Catalog information for the dataset.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"framework": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog framework.`,
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog entity name.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Catalog entity type.`,
									},
								},
							},
						},
						"column_lineage": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Column-level lineage information for the output dataset.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset_input": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Input fields affecting whole dataset, e.g. filtering columns.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Source field name.`,
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Name of the source dataset.`,
												},
												"namespace": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Namespace of the source dataset.`,
												},
												"transformation": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: `Transformations applied to fields from this input.`,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: `Transformation type.`,
															},
															"subtype": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: `Transformation subtype.`,
															},
														},
													},
												},
											},
										},
									},
									"field": {
										Type:        schema.TypeList,
										Required:    true,
										Description: `Field-level lineage mappings.`,
										MinItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"input": {
													Type:        schema.TypeList,
													Required:    true,
													Description: `Input fields contributing to this output field.`,
													MinItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"field": {
																Type:        schema.TypeString,
																Required:    true,
																Description: `Source field name.`,
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: `Name of the source dataset.`,
															},
															"namespace": {
																Type:        schema.TypeString,
																Required:    true,
																Description: `Namespace of the source dataset.`,
															},
															"transformation": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: `Transformations applied from source to output field.`,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: `Transformation type.`,
																		},
																		"subtype": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: `Transformation subtype.`,
																		},
																	},
																},
															},
														},
													},
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Output field name.`,
												},
											},
										},
									},
								},
							},
						},
						"symlink": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Symlink targets for the dataset.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Name of the symlink target.`,
									},
									"namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Namespace of the symlink target.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Type of the symlink target.`,
									},
								},
							},
						},
					},
				},
			},
			"owner": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The owner of the OpenLineage job.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Owner name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Owner type.`,
						},
					},
				},
			},
			"knowledge_catalog": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Knowledge Catalog entities generated for this lineage job.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"process": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Knowledge Catalog process identifier.`,
						},
						"run": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Knowledge Catalog run identifier.`,
						},
					},
				},
			},

			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `Whether Terraform will be prevented from destroying the instance. Defaults to "DELETE".
When a 'terraform destroy' or 'terraform apply' would delete the instance,
the command will fail if this field is set to "PREVENT" in Terraform state.
When set to "ABANDON", the command will remove the resource from Terraform
management without updating or deleting the resource in the API.
When set to "DELETE", deleting the resource is allowed.
`,
			},
		},
		UseJSONNumber: true,
	}
}

func flattenDataLineageOpenLineageJobNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOwner(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"name": flattenDataLineageOpenLineageJobOwnerName(original["name"], d, config),
			"type": flattenDataLineageOpenLineageJobOwnerType(original["type"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOwnerName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOwnerType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInput(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace": flattenDataLineageOpenLineageJobInputNamespace(original["namespace"], d, config),
			"name":      flattenDataLineageOpenLineageJobInputName(original["name"], d, config),
			"symlink":   flattenDataLineageOpenLineageJobInputSymlink(original["symlink"], d, config),
			"catalog":   flattenDataLineageOpenLineageJobInputCatalog(original["catalog"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobInputNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputSymlink(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace": flattenDataLineageOpenLineageJobInputSymlinkNamespace(original["namespace"], d, config),
			"name":      flattenDataLineageOpenLineageJobInputSymlinkName(original["name"], d, config),
			"type":      flattenDataLineageOpenLineageJobInputSymlinkType(original["type"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobInputSymlinkNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputSymlinkName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputSymlinkType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputCatalog(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["framework"] =
		flattenDataLineageOpenLineageJobInputCatalogFramework(original["framework"], d, config)
	transformed["type"] =
		flattenDataLineageOpenLineageJobInputCatalogType(original["type"], d, config)
	transformed["name"] =
		flattenDataLineageOpenLineageJobInputCatalogName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenDataLineageOpenLineageJobInputCatalogFramework(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputCatalogType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobInputCatalogName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutput(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace":      flattenDataLineageOpenLineageJobOutputNamespace(original["namespace"], d, config),
			"name":           flattenDataLineageOpenLineageJobOutputName(original["name"], d, config),
			"symlink":        flattenDataLineageOpenLineageJobOutputSymlink(original["symlink"], d, config),
			"catalog":        flattenDataLineageOpenLineageJobOutputCatalog(original["catalog"], d, config),
			"column_lineage": flattenDataLineageOpenLineageJobOutputColumnLineage(original["column_lineage"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputSymlink(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace": flattenDataLineageOpenLineageJobOutputSymlinkNamespace(original["namespace"], d, config),
			"name":      flattenDataLineageOpenLineageJobOutputSymlinkName(original["name"], d, config),
			"type":      flattenDataLineageOpenLineageJobOutputSymlinkType(original["type"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputSymlinkNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputSymlinkName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputSymlinkType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputCatalog(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["framework"] =
		flattenDataLineageOpenLineageJobOutputCatalogFramework(original["framework"], d, config)
	transformed["type"] =
		flattenDataLineageOpenLineageJobOutputCatalogType(original["type"], d, config)
	transformed["name"] =
		flattenDataLineageOpenLineageJobOutputCatalogName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenDataLineageOpenLineageJobOutputCatalogFramework(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputCatalogType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputCatalogName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineage(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["field"] =
		flattenDataLineageOpenLineageJobOutputColumnLineageField(original["field"], d, config)
	transformed["dataset_input"] =
		flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInput(original["dataset_input"], d, config)
	return []interface{}{transformed}
}
func flattenDataLineageOpenLineageJobOutputColumnLineageField(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"name":  flattenDataLineageOpenLineageJobOutputColumnLineageFieldName(original["name"], d, config),
			"input": flattenDataLineageOpenLineageJobOutputColumnLineageFieldInput(original["input"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputColumnLineageFieldName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInput(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace":      flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputNamespace(original["namespace"], d, config),
			"name":           flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputName(original["name"], d, config),
			"field":          flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputField(original["field"], d, config),
			"transformation": flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformation(original["transformation"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputField(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"type":    flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationType(original["type"], d, config),
			"subtype": flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationSubtype(original["subtype"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationSubtype(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInput(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"namespace":      flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputNamespace(original["namespace"], d, config),
			"name":           flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputName(original["name"], d, config),
			"field":          flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputField(original["field"], d, config),
			"transformation": flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformation(original["transformation"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputNamespace(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputField(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformation(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"type":    flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationType(original["type"], d, config),
			"subtype": flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationSubtype(original["subtype"], d, config),
		})
	}
	return transformed
}
func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationSubtype(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobKnowledgeCatalog(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["process"] =
		flattenDataLineageOpenLineageJobKnowledgeCatalogProcess(original["process"], d, config)
	transformed["run"] =
		flattenDataLineageOpenLineageJobKnowledgeCatalogRun(original["run"], d, config)
	return []interface{}{transformed}
}
func flattenDataLineageOpenLineageJobKnowledgeCatalogProcess(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenDataLineageOpenLineageJobKnowledgeCatalogRun(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func expandDataLineageOpenLineageJobNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobDescription(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOwner(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandDataLineageOpenLineageJobOwnerName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedType, err := expandDataLineageOpenLineageJobOwnerType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOwnerName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOwnerType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInput(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobInputNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobInputName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedSymlink, err := expandDataLineageOpenLineageJobInputSymlink(original["symlink"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSymlink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["symlink"] = transformedSymlink
		}

		transformedCatalog, err := expandDataLineageOpenLineageJobInputCatalog(original["catalog"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCatalog); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["catalog"] = transformedCatalog
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobInputNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputSymlink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobInputSymlinkNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobInputSymlinkName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedType, err := expandDataLineageOpenLineageJobInputSymlinkType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobInputSymlinkNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputSymlinkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputSymlinkType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputCatalog(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFramework, err := expandDataLineageOpenLineageJobInputCatalogFramework(original["framework"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFramework); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["framework"] = transformedFramework
	}

	transformedType, err := expandDataLineageOpenLineageJobInputCatalogType(original["type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	transformedName, err := expandDataLineageOpenLineageJobInputCatalogName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandDataLineageOpenLineageJobInputCatalogFramework(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputCatalogType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobInputCatalogName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutput(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobOutputNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobOutputName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedSymlink, err := expandDataLineageOpenLineageJobOutputSymlink(original["symlink"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSymlink); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["symlink"] = transformedSymlink
		}

		transformedCatalog, err := expandDataLineageOpenLineageJobOutputCatalog(original["catalog"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedCatalog); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["catalog"] = transformedCatalog
		}

		transformedColumnLineage, err := expandDataLineageOpenLineageJobOutputColumnLineage(original["column_lineage"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedColumnLineage); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["column_lineage"] = transformedColumnLineage
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputSymlink(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobOutputSymlinkNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobOutputSymlinkName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedType, err := expandDataLineageOpenLineageJobOutputSymlinkType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputSymlinkNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputSymlinkName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputSymlinkType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputCatalog(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedFramework, err := expandDataLineageOpenLineageJobOutputCatalogFramework(original["framework"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedFramework); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["framework"] = transformedFramework
	}

	transformedType, err := expandDataLineageOpenLineageJobOutputCatalogType(original["type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["type"] = transformedType
	}

	transformedName, err := expandDataLineageOpenLineageJobOutputCatalogName(original["name"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["name"] = transformedName
	}

	return transformed, nil
}

func expandDataLineageOpenLineageJobOutputCatalogFramework(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputCatalogType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputCatalogName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineage(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedField, err := expandDataLineageOpenLineageJobOutputColumnLineageField(original["field"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedField); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["field"] = transformedField
	}

	transformedDatasetInput, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInput(original["dataset_input"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDatasetInput); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["dataset_input"] = transformedDatasetInput
	}

	return transformed, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageField(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedName, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedInput, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInput(original["input"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedInput); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["input"] = transformedInput
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInput(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedField, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputField(original["field"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedField); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["field"] = transformedField
		}

		transformedTransformation, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformation(original["transformation"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedTransformation); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["transformation"] = transformedTransformation
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputField(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedType, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		transformedSubtype, err := expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationSubtype(original["subtype"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubtype); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subtype"] = transformedSubtype
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageFieldInputTransformationSubtype(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInput(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNamespace, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputNamespace(original["namespace"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNamespace); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["namespace"] = transformedNamespace
		}

		transformedName, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputName(original["name"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedName); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["name"] = transformedName
		}

		transformedField, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputField(original["field"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedField); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["field"] = transformedField
		}

		transformedTransformation, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformation(original["transformation"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedTransformation); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["transformation"] = transformedTransformation
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputNamespace(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputName(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputField(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformation(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedType, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationType(original["type"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedType); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["type"] = transformedType
		}

		transformedSubtype, err := expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationSubtype(original["subtype"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedSubtype); val.IsValid() && !tpgresource.IsEmptyValue(val) {
			transformed["subtype"] = transformedSubtype
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandDataLineageOpenLineageJobOutputColumnLineageDatasetInputTransformationSubtype(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func ResourceDataLineageOpenLineageJobFlatten(d *schema.ResourceData, meta interface{}, res map[string]interface{}, config *transport_tpg.Config, userAgent string, billingProject string, url string, headers http.Header) error {
	var err error

	if err = d.Set("namespace", flattenDataLineageOpenLineageJobNamespace(res["namespace"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("name", flattenDataLineageOpenLineageJobName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("description", flattenDataLineageOpenLineageJobDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("owner", flattenDataLineageOpenLineageJobOwner(res["owner"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("input", flattenDataLineageOpenLineageJobInput(res["input"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("output", flattenDataLineageOpenLineageJobOutput(res["output"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}
	if err = d.Set("knowledge_catalog", flattenDataLineageOpenLineageJobKnowledgeCatalog(res["knowledge_catalog"], d, config)); err != nil {
		return fmt.Errorf("Error reading OpenLineageJob: %s", err)
	}

	return nil
}
