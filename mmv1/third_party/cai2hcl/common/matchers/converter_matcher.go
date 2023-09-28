package matchers

import (
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"
)

type ConverterMatcher interface {
	GetConverterName() string

	Match(asset *caiasset.Asset) bool
}

func Matcher(regexp string) _AssetNameMatcher {
	return _AssetNameMatcher{Regexp: regexp}
}
