package recaptchaenterprise

import "testing"

// Regression test for hashicorp/terraform-provider-google#28502.
func TestCanonicalizeKeyDesiredStatePreservesInitialActionSettings(t *testing.T) {
	initial := &Key{
		Project:     stringPtr("example-project"),
		DisplayName: stringPtr("policy-based-challenge-example"),
		WebSettings: &KeyWebSettings{
			IntegrationType: KeyWebSettingsIntegrationTypeEnumRef("POLICY_BASED_CHALLENGE"),
			ChallengeSettings: &KeyWebSettingsChallengeSettings{
				DefaultSettings: &KeyWebSettingsChallengeSettingsDefaultSettings{
					ScoreThreshold: float64Ptr(0.1),
				},
				ActionSettings: map[string]KeyWebSettingsChallengeSettingsActionSettings{
					"auth_email": {ScoreThreshold: float64Ptr(0.45)},
					"auth_phone": {ScoreThreshold: float64Ptr(0.1)},
				},
			},
		},
	}

	desired := &Key{
		Project:     stringPtr("example-project"),
		DisplayName: stringPtr("policy-based-challenge-example"),
		WebSettings: &KeyWebSettings{
			IntegrationType: KeyWebSettingsIntegrationTypeEnumRef("POLICY_BASED_CHALLENGE"),
			ChallengeSettings: &KeyWebSettingsChallengeSettings{
				DefaultSettings: &KeyWebSettingsChallengeSettingsDefaultSettings{
					ScoreThreshold: float64Ptr(0.1),
				},
				ActionSettings: map[string]KeyWebSettingsChallengeSettingsActionSettings{
					"auth": {ScoreThreshold: float64Ptr(0.45)},
				},
			},
		},
	}

	got, err := canonicalizeKeyDesiredState(desired, initial)
	if err != nil {
		t.Fatalf("canonicalizeKeyDesiredState returned an error: %v", err)
	}

	if got.WebSettings == nil || got.WebSettings.ChallengeSettings == nil || got.WebSettings.ChallengeSettings.ActionSettings == nil {
		t.Fatal("expected canonicalized key to preserve challenge settings action settings")
	}

	actionSettings := got.WebSettings.ChallengeSettings.ActionSettings
	if len(actionSettings) != 1 {
		t.Fatalf("expected only the desired auth action after canonicalization, got %d entries", len(actionSettings))
	}

	if _, ok := actionSettings["auth"]; !ok {
		t.Fatalf("expected the desired auth action to be present")
	}
	if _, ok := actionSettings["auth_email"]; ok {
		t.Fatalf("did not expect auth_email to remain after the update")
	}
	if _, ok := actionSettings["auth_phone"]; ok {
		t.Fatalf("did not expect auth_phone to remain after the update")
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}
