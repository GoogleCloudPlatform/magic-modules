package matchers

import (
	"regexp"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"
)

func ByAssetName(regexp string, converterName string) ConverterMatcher {
	return _AssetNameMatcher{Regexp: regexp, ConverterName: converterName}
}

type _AssetNameMatcher struct {
	Regexp        string
	ConverterName string
}

func (inst _AssetNameMatcher) Match(asset *caiasset.Asset) bool {
	matched, _ := regexp.MatchString(inst.Regexp, asset.Name)
	return matched
}

func (inst _AssetNameMatcher) GetConverterName() string {
	return inst.ConverterName
}
