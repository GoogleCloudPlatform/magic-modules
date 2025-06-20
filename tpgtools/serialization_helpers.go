package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	cloudresourcemanager "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager"
	cloudresourcemanagerAlpha "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager/alpha"
	cloudresourcemanagerBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager/beta"
)

func serializeAlphaProjectToHCL(r cloudresourcemanagerAlpha.Project, hasGAEquivalent bool) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return serializeProjectToHCL(m, hasGAEquivalent)
}

func serializeBetaProjectToHCL(r cloudresourcemanagerBeta.Project, hasGAEquivalent bool) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return serializeProjectToHCL(m, hasGAEquivalent)
}

func serializeGAProjectToHCL(r cloudresourcemanager.Project, hasGAEquivalent bool) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return serializeProjectToHCL(m, hasGAEquivalent)
}

func serializeProjectToHCL(m map[string]interface{}, hasGAEquivalent bool) (string, error) {
	outputConfig := "resource \"google_project\" \"output\" {\n"
	if name, ok := m["name"]; ok {
		outputConfig += fmt.Sprintf("\tproject_id = %#v\n", name)
		outputConfig += fmt.Sprintf("\tname = %#v\n", name)
	} else {
		return "", fmt.Errorf("project id was not provided")
	}
	if parentInterface, ok := m["parent"]; ok {
		parent, ok := parentInterface.(string)
		if !ok {
			return "", fmt.Errorf("non-string parent %v", parentInterface)
		}
		if strings.HasPrefix(parent, "folders/") {
			outputConfig += fmt.Sprintf("\tfolder_id = %#v\n", strings.TrimPrefix(parent, "folders/"))
		} else if strings.HasPrefix(parent, "organizations/") {
			outputConfig += fmt.Sprintf("\torg_id = %#v\n", strings.TrimPrefix(parent, "organizations/"))
		} else {
			return "", fmt.Errorf("unknown parent type %v", parent)
		}
	} else {
		return "", fmt.Errorf("parent was not provided")
	}
	formatted, err := formatHCL(outputConfig + "}")
	if err != nil {
		// The formatter will not accept the google-beta symbol because it is injected during testing.
		return "", err
	}
	if !hasGAEquivalent {
		return withProviderLine(formatted), nil
	}
	return formatted, nil
}

// Returns the terraform representation of a three-state boolean value represented by a pointer to bool in DCL.
func serializeEnumBool(v interface{}) string {
	b, ok := v.(*bool)
	if !ok || b == nil {
		return ""
	}
	if *b {
		return "TRUE"
	}
	return "FALSE"
}

// Returns the given formatted hcl with the provider = google-beta line added at the end.
func withProviderLine(hcl string) string {
	// Count the number of characters before the first = to determine how to space the provider line.
	equalsPosition := len(regexp.MustCompile(".*=").FindString(hcl)) - 1
	if equalsPosition < 11 {
		equalsPosition = 11
	}
	return hcl[0:len(hcl)-2] + "  provider" + strings.Repeat(" ", equalsPosition-10) + "= google-beta\n}"
}
