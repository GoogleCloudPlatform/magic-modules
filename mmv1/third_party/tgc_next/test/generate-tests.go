package test

import (
	"fmt"
	"log"
	"os"
)

func generateTests(TestConfig map[string]TestMetadata, resource, asset string) {
	testsPerResource := make(map[string][]string, 0)
	for test, config := range TestConfig {
		if _, ok := testsPerResource[config.Resource]; !ok {
			testsPerResource[config.Resource] = make([]string, 0)
		}

		testsPerResource[config.Resource] = append(testsPerResource[config.Resource], test)
	}

	for r, tests := range testsPerResource {
		if r != resource {
			continue
		}
		total := fmt.Sprintf("// Total %d tests\n", len(tests))
		for _, test := range tests {
			str := fmt.Sprintf(`
		func %s(t *testing.T) {
			t.Parallel()
		
			test.AssertTestFile(
				t,
				"%s",
				"%s",
				"%s",
				[]string{
			"desired_status",
			"metadata",
			// "boot_disk.initialize_params" is not converted
			"boot_disk.initialize_params",
			"boot_disk.initialize_params.image",
				},
			)
		}`, test, test, r, asset)
			total = fmt.Sprintf("%s\n%s", total, str)
		}

		filePath := fmt.Sprintf("%s.go", resource)
		err := os.WriteFile(filePath, []byte(total), 0644)
		if err != nil {
			log.Fatalf("error writing to file %s: %#v", filePath, err)
		}
	}
}
