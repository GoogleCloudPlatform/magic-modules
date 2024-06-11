package google

import (
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/tfplan2cai/converters/google/resources/cai"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const AppEngineApplicationAssetType string = "appengine.googleapis.com/Application"

func resourceConverterAppEngineApplication() cai.ResourceConverter {
	return cai.ResourceConverter{
		AssetType: AppEngineApplicationAssetType,
		Convert:   GetAppEngineApplicationCaiObject,
	}
}

func GetAppEngineApplicationCaiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) ([]cai.Asset, error) {
	name, err := cai.AssetName(d, config, "//appengine.googleapis.com/v1/{{name}}")
	if err != nil {
		return []cai.Asset{}, err
	}
	if obj, err := GetAppEngineApplicationApiObject(d, config); err == nil {
		return []cai.Asset{{
			Name: name,
			Type: AppEngineApplicationAssetType,
			Resource: &cai.AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://appengine/docs/admin-api/reference/rest/v1/apps/",
				DiscoveryName:        "AppEngineApplication",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []cai.Asset{}, err
	}
}

func GetAppEngineApplicationApiObject(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	idProp, err := expandAppEngineApplicationId(d.Get("id"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("id"); !tpgresource.IsEmptyValue(reflect.ValueOf(idProp)) && (ok || !reflect.DeepEqual(v, idProp)) {
		obj["id"] = idProp
	}

	dispatchRulesProp, err := expandAppEngineApplicationDispatchRules(d.Get("dispatch_rules"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("dispatch_rules"); !tpgresource.IsEmptyValue(reflect.ValueOf(dispatchRulesProp)) && (ok || !reflect.DeepEqual(v, dispatchRulesProp)) {
		obj["dispatch_rules"] = dispatchRulesProp
	}

	authDomainProp, err := expandAppEngineApplicationAuthDomain(d.Get("authDomain"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("authDomain"); !tpgresource.IsEmptyValue(reflect.ValueOf(authDomainProp)) && (ok || !reflect.DeepEqual(v, authDomainProp)) {
		obj["authDomain"] = authDomainProp
	}

	locationIdProp, err := expandAppEngineApplicationLocationId(d.Get("locationId"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("locationId"); !tpgresource.IsEmptyValue(reflect.ValueOf(locationIdProp)) && (ok || !reflect.DeepEqual(v, locationIdProp)) {
		obj["locationId"] = locationIdProp
	}

	codeBucketProp, err := expandAppEngineApplicationCodeBucket(d.Get("codeBucket"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("codeBucket"); !tpgresource.IsEmptyValue(reflect.ValueOf(codeBucketProp)) && (ok || !reflect.DeepEqual(v, codeBucketProp)) {
		obj["codeBucket"] = codeBucketProp
	}

	defaultCookieExpirationProp, err := expandAppEngineApplicationDefaultCookieExpiration(d.Get("defaultCookieExpiration"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("defaultCookieExpiration"); !tpgresource.IsEmptyValue(reflect.ValueOf(defaultCookieExpirationProp)) && (ok || !reflect.DeepEqual(v, defaultCookieExpirationProp)) {
		obj["defaultCookieExpiration"] = defaultCookieExpirationProp
	}

	servingStatusProp, err := expandAppEngineApplicationServingStatus(d.Get("servingStatus"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("servingStatus"); !tpgresource.IsEmptyValue(reflect.ValueOf(servingStatusProp)) && (ok || !reflect.DeepEqual(v, servingStatusProp)) {
		obj["servingStatus"] = servingStatusProp
	}

	defaultHostnameProp, err := expandAppEngineApplicationDefaultHostname(d.Get("defaultHostname"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("defaultHostname"); !tpgresource.IsEmptyValue(reflect.ValueOf(defaultHostnameProp)) && (ok || !reflect.DeepEqual(v, defaultHostnameProp)) {
		obj["defaultHostname"] = defaultHostnameProp
	}

	defaultBucketProp, err := expandAppEngineApplicationdefaultBucket(d.Get("defaultBucket"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("defaultBucket"); !tpgresource.IsEmptyValue(reflect.ValueOf(defaultBucketProp)) && (ok || !reflect.DeepEqual(v, defaultBucketProp)) {
		obj["defaultBucket"] = defaultBucketProp
	}

	iapProp, err := expandAppEngineApplicationIap(d.Get("iap"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("iap"); !tpgresource.IsEmptyValue(reflect.ValueOf(iapProp)) && (ok || !reflect.DeepEqual(v, iapProp)) {
		obj["iap"] = iapProp
	}

	gcrDomainProp, err := expandAppEngineApplicationGcrDomain(d.Get("gcrDomain"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("gcrDomain"); !tpgresource.IsEmptyValue(reflect.ValueOf(gcrDomainProp)) && (ok || !reflect.DeepEqual(v, gcrDomainProp)) {
		obj["gcrDomain"] = gcrDomainProp
	}

	databaseTypeProp, err := expandAppEngineApplicationDatabaseType(d.Get("databaseType"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("databaseType"); !tpgresource.IsEmptyValue(reflect.ValueOf(databaseTypeProp)) && (ok || !reflect.DeepEqual(v, databaseTypeProp)) {
		obj["databaseType"] = databaseTypeProp
	}

	featureSettingsProp, err := expandAppEngineApplicationFeatureSettings(d.Get("featureSettings"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("featureSettings"); !tpgresource.IsEmptyValue(reflect.ValueOf(featureSettingsProp)) && (ok || !reflect.DeepEqual(v, featureSettingsProp)) {
		obj["featureSettings"] = featureSettingsProp
	}



	return obj, nil
}

func expandAppEngineApplicationId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationDispatchRules(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUrlDispatchRule, err := expandAppEngineApplicationUrlDispatchRule(original["url_dispatch_rule"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUrlDispatchRule); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["url_dispatch_rule"] = transformedUrlDispatchRule
	}

	
	return transformed, nil
}

func expandAppEngineApplicationUrlDispatchRule(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDomain, err := expandAppEngineApplicationDomain(original["domain"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDomain); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["domain"] = transformedDomain
	}

	transformedPath, err := expandAppEngineApplicationPath(original["path"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPath); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["path"] = transformedPath
	}

	transformedService, err := expandAppEngineApplicationService(original["service"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedService); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["service"] = transformedService
	}

	
	return transformed, nil
}

func expandAppEngineApplicationDomain(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationPath(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationService(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationAuthDomain(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationLocationId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationCodeBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationDefaultCookieExpiration(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationServingStatus(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedEnabled, err := expandAppEngineServingStatusEnabled(original["enabled"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedEnabled); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["enabled"] = transformedEnabled
	}

	transformedOauth2ClientId, err := expandAppEngineServingStatusOauth2ClientId(original["oauth2ClientId"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauth2ClientId); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauth2ClientId"] = transformedOauth2ClientId
	}

	transformedOauth2ClientSecret, err := expandAppEngineServingStatusOauth2ClientSecret(original["oauth2ClientSecret"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauth2ClientSecret); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauth2ClientSecret"] = transformedOauth2ClientSecret
	}

	transformedOauth2ClientSecretSha256, err := expandAppEngineServingStatusOauth2ClientSecretSha256(original["oauth2ClientSecretSha256"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauth2ClientSecretSha256); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["oauth2ClientSecretSha256"] = transformedOauth2ClientSecretSha256
	}
	
	
	return transformed, nil
}

func expandAppEngineServingStatusEnabled(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineServingStatusOauth2ClientId(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineServingStatusOauth2ClientSecret(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineServingStatusOauth2ClientSecretSha256(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationDefaultHostname(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationdefaultBucket(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationIap(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationGcrDomain(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationDatabaseType(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineApplicationFeatureSettings(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedSplitHealthChecks, err := expandAppEngineFeatureSettingsSplitHealthChecks(original["splitHealthChecks"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedSplitHealthChecks); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["splitHealthChecks"] = transformedSplitHealthChecks
	}

	transformedUseContainerOptimizedOs, err := expandAppEngineFeatureSettingsUseContainerOptimizedOs(original["useContainerOptimizedOs"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUseContainerOptimizedOs); val.IsValid() && !tpgresource.IsEmptyValue(val) {
		transformed["useContainerOptimizedOs"] = transformedUseContainerOptimizedOs
	}

	return transformed, nil
}

func expandAppEngineFeatureSettingsSplitHealthChecks(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}

func expandAppEngineFeatureSettingsUseContainerOptimizedOs(v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	return v, nil
}










