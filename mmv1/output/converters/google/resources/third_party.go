package cai2hcl

import "regexp"

func ConvertSelfLinkToV1(link string) string {
	reg := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/")
	return reg.ReplaceAllString(link, "/compute/v1/projects/")
}
