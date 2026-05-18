package cmd

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const unusedTmplDesc = "Check whether any template files are not used in product yamls"

var exampleFilePathReg = regexp.MustCompile(".*mmv1/templates/terraform/examples/([a-zA-Z0-9_-]+).tf.tmpl")
var sampleFilePathReg = regexp.MustCompile(".*mmv1/templates/terraform/samples/services/.*\\.tf\\.tmpl")

type unusedTmplOptions struct {
	rootOptions *rootOptions
	stdout      io.Writer
	fileList    []string
}

type tree struct {
	tmplPaths map[string]bool
}

type resourceYaml struct {
	Examples []struct {
		Name string `yaml:"name"`
	} `yaml:"examples,omitempty"`
	Samples []struct {
		Name  string `yaml:"name"`
		Steps []struct {
			Name       string `yaml:"name"`
			ConfigPath string `yaml:"config_path,omitempty"`
		} `yaml:"steps"`
	} `yaml:"samples,omitempty"`
}

func newUnusedTmplCmd(rootOptions *rootOptions) *cobra.Command {
	o := &unusedTmplOptions{
		rootOptions: rootOptions,
		stdout:      os.Stdout,
	}
	command := &cobra.Command{
		Use:   "unused-tmpl",
		Short: unusedTmplDesc,
		Long:  unusedTmplDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	command.Flags().StringSliceVar(&o.fileList, "file-list", []string{}, "file list to check")
	return command
}

func (o *unusedTmplOptions) run() error {
	if len(o.fileList) == 0 {
		return nil
	}
	newCustomTmpls, newExamples, newSamples := processInputFiles(o.fileList)

	found := false
	// get repo dir from tmpl files
	repoPath := strings.Split(o.fileList[0], "/mmv1/")[0]
	dir := filepath.Join(repoPath, "mmv1", "products")

	productFiles, err := yamlFiles(dir)
	if err != nil {
		return err
	}
	if len(newCustomTmpls) > 0 {
		customTempls, err := findTmpls(productFiles)
		if err != nil {
			return err
		}
		for _, file := range newCustomTmpls {
			templatePath := strings.ReplaceAll(file, repoPath+"/mmv1/", "")
			if _, ok := customTempls[templatePath]; !ok {
				found = true
				fmt.Fprintf(os.Stderr, "File %s not used in any product yaml.\n", file)
			}
		}
	}
	if len(newExamples) > 0 {
		ex, err := findExamples(productFiles)
		if err != nil {
			return err
		}
		for _, file := range newExamples {
			baseName := filepath.Base(file)
			newExName := strings.TrimSuffix(baseName, ".tf.tmpl")
			if _, ok := ex[newExName]; !ok {
				found = true
				fmt.Fprintf(os.Stderr, "File %s not used in any product yaml.\n", file)
			}
		}

	}
	if len(newSamples) > 0 {
		samples, err := findSamples(productFiles)
		if err != nil {
			return err
		}
		for _, file := range newSamples {
			templatePath := strings.ReplaceAll(file, repoPath+"/mmv1/", "")
			if _, ok := samples[templatePath]; !ok {
				found = true
				fmt.Fprintf(os.Stderr, "File %s not used in any product yaml.\n", file)
			}
		}
	}
	if found {
		return fmt.Errorf("found templates not used")
	}
	return nil
}

func processInputFiles(fileList []string) (customTmpls []string, examples []string, samples []string) {
	for _, v := range fileList {
		if exampleFilePathReg.MatchString(v) {
			examples = append(examples, v)
		} else if sampleFilePathReg.MatchString(v) {
			samples = append(samples, v)
		} else if strings.Contains(v, "mmv1/templates/terraform") && strings.HasSuffix(v, ".tmpl") {
			customTmpls = append(customTmpls, v)
		} else {
			fmt.Printf("Skipping check for file %s\n", v)
		}
	}
	return
}

func (t *tree) walkTree(tree map[any]any) {
	for _, value := range tree {
		switch v := value.(type) {
		case []any:
			for _, v1 := range v {
				if val, ok := v1.(map[any]any); ok {
					t.walkTree(val)
				}
			}
		case map[any]any:
			t.walkTree(v)
		case string:
			if strings.HasSuffix(v, ".tmpl") {
				t.tmplPaths[v] = true
			}
		default:
		}
	}
}

func yamlFiles(dir string) ([]string, error) {
	var allYamlFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".yaml" {
			allYamlFiles = append(allYamlFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allYamlFiles, nil
}

var (
	underscoreCap = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	underscoreLow = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

func underscore(source string) string {
	tmp := underscoreCap.ReplaceAllString(source, "${1}_${2}")
	tmp = underscoreLow.ReplaceAllString(tmp, "${1}_${2}")
	tmp = strings.ReplaceAll(tmp, "-", "_")
	tmp = strings.ReplaceAll(tmp, ".", "_")
	tmp = strings.ToLower(tmp)
	return tmp
}

// findTmpls parsed yaml files to get values ending with .tmpl.
// It returns a map of tmpl paths where the key is the tmpl path.
func findTmpls(yamlFiles []string) (map[string]bool, error) {
	allTmpls := map[string]bool{}

	productNames := map[string]string{}
	for _, yamlFile := range yamlFiles {
		if filepath.Base(yamlFile) == "product.yaml" {
			b, err := os.ReadFile(yamlFile)
			if err != nil {
				continue
			}
			var m map[any]any
			if err := yaml.Unmarshal(b, &m); err == nil {
				for k, v := range m {
					if keyStr, ok := k.(string); ok && keyStr == "name" {
						if nameStr, ok := v.(string); ok {
							productNames[filepath.Dir(yamlFile)] = nameStr
						}
					}
				}
			}
		}
	}

	for _, yamlFile := range yamlFiles {
		b, err := os.ReadFile(yamlFile)
		if err != nil {
			return nil, err
		}
		var m map[any]any
		if err := yaml.Unmarshal(b, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml file %s: %s", yamlFile, err)
		}

		var resName string
		hasStateUpgraders := false

		for k, v := range m {
			if keyStr, ok := k.(string); ok {
				if keyStr == "name" {
					if nameStr, ok := v.(string); ok {
						resName = nameStr
					}
				}
				if keyStr == "state_upgraders" {
					if su, ok := v.(bool); ok && su {
						hasStateUpgraders = true
					}
				}
			}
		}

		if hasStateUpgraders && resName != "" {
			productName := productNames[filepath.Dir(yamlFile)]
			if productName == "" {
				productName = filepath.Base(filepath.Dir(yamlFile))
			}
			tmplName := fmt.Sprintf("templates/terraform/state_migrations/%s_%s.go.tmpl", underscore(productName), underscore(resName))
			allTmpls[tmplName] = true
		}

		tr := &tree{
			tmplPaths: make(map[string]bool),
		}
		tr.walkTree(m)
		maps.Copy(allTmpls, tr.tmplPaths)
	}
	return allTmpls, nil
}

// findExamples parsed yaml files to get examples.
// It returns a map of examples where the key is the example name.
func findExamples(yamlFiles []string) (map[string]bool, error) {
	allExamples := map[string]bool{}
	for _, yamlFile := range yamlFiles {
		b, err := os.ReadFile(yamlFile)
		if err != nil {
			return nil, err
		}

		var r resourceYaml
		if err := yaml.Unmarshal(b, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml file for examples%s: %s", yamlFile, err)
		}
		for _, v := range r.Examples {
			allExamples[v.Name] = true
		}
	}
	return allExamples, nil
}

// findSamples parsed yaml files to get samples.
// It returns a map of samples where the key is the inferred sample path.
func findSamples(yamlFiles []string) (map[string]bool, error) {
	allSamples := map[string]bool{}
	for _, yamlFile := range yamlFiles {
		b, err := os.ReadFile(yamlFile)
		if err != nil {
			return nil, err
		}

		var r resourceYaml
		if err := yaml.Unmarshal(b, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml file for samples %s: %s", yamlFile, err)
		}
		packageName := filepath.Base(filepath.Dir(yamlFile))
		for _, sample := range r.Samples {
			for _, step := range sample.Steps {
				var tmplPath string
				if step.ConfigPath != "" {
					tmplPath = step.ConfigPath
				} else {
					tmplPath = fmt.Sprintf("templates/terraform/samples/services/%s/%s.tf.tmpl", packageName, step.Name)
				}
				allSamples[tmplPath] = true
			}
		}
	}
	return allSamples, nil
}
