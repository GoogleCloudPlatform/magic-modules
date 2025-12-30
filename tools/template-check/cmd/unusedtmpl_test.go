package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessInput(t *testing.T) {
	fileList := []string{
		"mmv1/templates/terraform/examples/abc.tf.tmpl",
		"mmv1/templates/terraform/examples/abc.go.tmpl",
		"mmv1/templates/terraform/examples/subfolder/abc.tf.tmpl",
		"mmv1/templates/terraform/custom_flatten/abc.go.tmpl",
	}
	tmpl, examples := processInputFiles(fileList)
	wantTmpl, wantExamples := []string{
		"mmv1/templates/terraform/examples/abc.go.tmpl",
		"mmv1/templates/terraform/examples/subfolder/abc.tf.tmpl",
		"mmv1/templates/terraform/custom_flatten/abc.go.tmpl",
	}, []string{
		"mmv1/templates/terraform/examples/abc.tf.tmpl",
	}

	if diff := cmp.Diff(wantTmpl, tmpl); diff != "" {
		t.Errorf("processInputFiles() got diff(-want, got) for template files = %s", diff)
	}
	if diff := cmp.Diff(wantExamples, examples); diff != "" {
		t.Errorf("processInputFiles() got diff(-want, got) for example files = %s", diff)
	}
}

func TestFindTmpls(t *testing.T) {
	yamlFiles := []string{"testdata/resource1.yaml", "testdata/resource2.yaml"}
	got, err := findTmpls(yamlFiles)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{
		"templates/terraform/custom_flatten/bigquery_table_ref_query_destinationtable.go.tmpl": true,
		"templates/terraform/custom_expand/bigquery_table_ref.go.tmpl":                         true,
		"templates/terraform/constants/bigquery_job.go.tmpl":                                   true,
		"templates/terraform/encoders/bigquery_job.go.tmpl":                                    true,
		"templates/terraform/custom_flatten/bigquery_table_ref_extract_sourcetable.go.tmpl":    true,
		"templates/terraform/custom_flatten/bigquery_kms_version.go.tmpl":                      true,
		"templates/terraform/custom_flatten/bigquery_table_ref_copy_destinationtable.go.tmpl":  true,
		"templates/terraform/custom_expand/bigquery_table_ref_array.go.tmpl":                   true,
		"templates/terraform/custom_flatten/bigquery_table_ref_copy_sourcetables.go.tmpl":      true,
		"templates/terraform/custom_flatten/bigquery_table_ref_load_destinationtable.go.tmpl":  true,
		"templates/terraform/custom_expand/bigquery_dataset_ref.go.tmpl":                       true,
		"templates/terraform/custom_flatten/bigquery_dataset_ref.go.tmpl":                      true,
		"templates/terraform/iam/example_config_body/app_engine_service.tf.tmpl":               true,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findTmpls() got unexpected diff(-want, got) = %s", diff)
	}

}

func TestFindExamples(t *testing.T) {
	yamlFiles := []string{"testdata/resource1.yaml", "testdata/resource2.yaml"}
	got, err := findExamples(yamlFiles)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{
		"bigquery_job_query":                   true,
		"bigquery_job_query_continuous":        true,
		"bigquery_job_query_table_reference":   true,
		"bigquery_job_load":                    true,
		"bigquery_job_load_geojson":            true,
		"bigquery_job_load_parquet":            true,
		"bigquery_job_load_table_reference":    true,
		"bigquery_job_copy":                    true,
		"bigquery_job_copy_table_reference":    true,
		"bigquery_job_extract":                 true,
		"bigquery_job_extract_table_reference": true,
		"iap_app_engine_service":               true,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findExamples() got unexpected diff(-want, got) = %s", diff)
	}
}
