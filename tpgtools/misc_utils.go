package main

import (
	"reflect"

	directory "github.com/GoogleCloudPlatform/declarative-resource-client-library/services"
)

func isLastIndex(array []string, index int) bool {
	return len(array)-1 == index
}

func getDCLPackageLocation() string {
	var miscPointer *directory.Directory
	return reflect.TypeOf(miscPointer).PkgPath()
}
