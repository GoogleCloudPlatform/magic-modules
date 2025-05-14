package dataplex

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceDataplexDataQualityRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataplexDataQualityRulesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_scan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The unnested column which this rule is evaluated against. e.g. `,
						},
						"ignore_null": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
							Description: `Rows with null values will automatically fail a rule, unless ignoreNull is true. In that case, such null rows are trivially considered passing. 
											This field is only valid for the following type of rules: RangeExpectation, RegexExpectation, SetExpectation, UniquenessExpectation`,
						},
						"dimension": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dimension a rule belongs to. Supported dimensions are "COMPLETENESS", "ACCURACY", "CONSISTENCY", "VALIDITY", "UNIQUENESS", "FRESHNESS", "VOLUME"`,
						},
						"threshold": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Optional:    true,
							Description: `The minimum ratio of passing_rows / total_rows required to pass this rule, with a range of [0.0, 1.0]. 0 indicates default value (i.e. 1.0). This field is only valid for row-level type rules.`,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							Description: `A mutable name for the rule. 
											The name must contain only letters (a-z, A-Z), numbers (0-9), or hyphens (-).
											The maximum length is 63 characters.
											Must start with a letter.
											Must end with a number or a letter.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `Description of the rule. (The maximum length is 1,024 characters.)`,
						},
						"suspended": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: `Whether the Rule is active or suspended. Default is false.`,
						},
						"range_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The minimum column value allowed for a row to pass this validation.`,
									},
									"max_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The maximum column value allowed for a row to pass this validation.`,
									},
									"strict_min_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										Description: `Whether each value needs to be strictly greater than ('>') the minimum, or if equality is allowed.`,
									},
									"strict_max_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										Description: ` Whether each value needs to be strictly lesser than ('<') the maximum, or if equality is allowed.`,
									},
								},
							},
							Description: `Row-level rule which evaluates whether each column value lies between a specified range.`,
						},
						"non_null_expectation": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Row-level rule which evaluates whether each column value is null.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
							// Elem: &schema.Schema{},
						},
						"set_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Expected values for the column value.`,
									},
								},
							},
							Description: `Row-level rule which evaluates whether each column value is contained by a specified set.`,
						},
						"regex_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"regex": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A regular expression the column value is expected to match.`,
									},
								},
							},

							Description: `Row-level rule which evaluates whether each column value matches a specified regex.`,
						},
						"uniqueness_expectation": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Row-level rule which evaluates whether each column value is unique.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
							// Elem: &schema.Schema{},
						},
						"statistic_range_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"statistic": {
										Type:     schema.TypeString,
										Computed: true,
										Description: `The list of aggregate metrics a rule can be evaluated against. 
																	Possible values: ["STATISTIC_UNDEFINED", "MEAN", "MIN", "MAX"]`,
									},
									"min_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The minimum column value allowed for a row to pass this validation.`,
									},
									"max_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The maximum column value allowed for a row to pass this validation.`,
									},
									"strict_min_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										Description: `Whether each value needs to be strictly greater than ('>') the minimum, or if equality is allowed.`,
									},
									"strict_max_enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Optional:    true,
										Description: ` Whether each value needs to be strictly lesser than ('<') the maximum, or if equality is allowed.`,
									},
								},
							},
							Description: `Aggregate rule which evaluates whether the column aggregate statistic lies between a specified range.`,
						},
						"row_condition_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sql_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The SQL expression.`,
									},
								},
							},
							Description: `Row-level rule which evaluates whether each row in a table passes the specified condition.`,
						},
						"table_condition_expectation": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sql_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The SQL expression.`,
									},
								},
							},
							Description: `Aggregate rule which evaluates whether the provided expression is true for a table.`,
						},
						"sql_assertion": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sql_statement": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The SQL expression.`,
									},
								},
							},
							Description: `Aggregate rule which evaluates the number of rows returned for the provided statement. If any rows are returned, this rule fails.`,
						},
					},
				},
			},
		},
	}
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, ch := range s {
		if unicode.IsUpper(ch) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(ch))
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func parseRulesResponse(res map[string]interface{}) ([]map[string]interface{}, error) {
	rulesToSet := make([]map[string]interface{}, 0)

	// if response doesn't include rule
	if _, ok := res["rule"].([]map[string]interface{}); !ok {
		return rulesToSet, nil
	}

	for _, apiRuleRaw := range res["rule"].([]map[string]interface{}) {
		newRuleMap := make(map[string]interface{})
		for k, v := range apiRuleRaw {
			snakeCaseKey := camelToSnake(k)

			if k == "nonNullExpectation" || k == "uniquenessExpectation" {
				newRuleMap[snakeCaseKey] = []interface{}{}
			} else {
				// For other fields (column, dimension, threshold, etc.), directly assign
				newRuleMap[snakeCaseKey] = v
			}
		}
		rulesToSet = append(rulesToSet, newRuleMap)
	}

	return rulesToSet, nil
}

func dataSourceDataplexDataQualityRulesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	if len(location) == 0 {
		return fmt.Errorf("Cannot determine location: set location in this data source or at provider-level")
	}

	data_scan_id := d.Get("data_scan_id").(string)

	url, err := tpgresource.ReplaceVars(d, config, "{{DataplexBasePath}}projects/{{project}}/locations/{{location}}/dataScans/{{data_scan_id}}:generateDataQualityRules")
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:               config,
		Method:               "POST",
		Project:              project,
		RawURL:               url,
		UserAgent:            userAgent,
		ErrorAbortPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429QuotaError},
	})

	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("DataQualityRules %q", d.Id()), url)
	}

	log.Printf("[jimmyxjc debug] res: %s", res)

	rules, err := parseRulesResponse(res)
	if err != nil {
		return fmt.Errorf("Error parsing rules: %s", err)
	}

	log.Printf("[jimmyxjc debug] rules: %s", rules)

	if err := d.Set("rule", rules); err != nil {
		return fmt.Errorf("Error setting rule: %s", err)
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/locations/%s/dataScans/%s", project, location, data_scan_id)
	d.SetId(id)

	return nil
}
