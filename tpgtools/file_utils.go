package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func mergeYaml(fileA, fileB string) ([]byte, error) {
	var objA map[string]interface{}
	bs, err := ioutil.ReadFile(fileA)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bs, &objA); err != nil {
		return nil, err
	}

	var objB map[string]interface{}
	bs, err = ioutil.ReadFile(fileB)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bs, &objB); err != nil {
		return nil, err
	}

	for k, v := range objB {
		objA[k] = v
	}

	out, err := yaml.Marshal(objA)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func pathExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
