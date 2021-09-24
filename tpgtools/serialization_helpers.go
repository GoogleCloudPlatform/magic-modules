package main

import (
	"encoding/json"
	"fmt"
	"strings"

	cloudresourcemanager "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager"
	cloudresourcemanagerBeta "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/cloudresourcemanager/beta"
)

func serializeBetaProjectToHCL(r cloudresourcemanagerBeta.Project) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return serializeProjectToHCL(m)
}

func serializeGAProjectToHCL(r cloudresourcemanager.Project) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return serializeProjectToHCL(m)
}

func serializeProjectToHCL(m map[string]interface{}) (string, error) {
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
	return formatHCL(outputConfig + "}")
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
