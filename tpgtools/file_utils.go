package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func mergeYaml(fileA, fileB string) string {
	var objA map[string]interface{}
	bs, err := ioutil.ReadFile(fileA)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, &objA); err != nil {
		panic(err)
	}

	var objB map[string]interface{}
	bs, err = ioutil.ReadFile(fileB)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bs, &objB); err != nil {
		panic(err)
	}

	for k, v := range objB {
		objA[k] = v
	}

	out, err := yaml.Marshal(objA)
	if err != nil {
		panic(err)
	}
	return string(out)
}
