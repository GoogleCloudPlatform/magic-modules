package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// newSecurityPolicyRuleForHash returns a fully-populated, schema-decoded rule map,
// matching the shape produced after schema decoding / flattenSecurityPolicyRules. Every
// call returns an independent copy so callers can mutate it freely.
func newSecurityPolicyRuleForHash() map[string]interface{} {
	return map[string]interface{}{
		"priority":    int(1000),
		"action":      "rate_based_ban",
		"preview":     false,
		"description": "a rule",
		"match": []interface{}{
			map[string]interface{}{
				"versioned_expr": "SRC_IPS_V1",
				"config": []interface{}{
					map[string]interface{}{
						"src_ip_ranges": schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8", "192.168.0.0/16"}),
					},
				},
				"expr": []interface{}{
					map[string]interface{}{
						"expression": "origin.region_code == 'US'",
					},
				},
				"expr_options": []interface{}{
					map[string]interface{}{
						"recaptcha_options": []interface{}{
							map[string]interface{}{
								"action_token_site_keys":  []interface{}{"action-key-1"},
								"session_token_site_keys": []interface{}{"session-key-1"},
							},
						},
					},
				},
			},
		},
		"preconfigured_waf_config": []interface{}{
			map[string]interface{}{
				"exclusion": []interface{}{
					map[string]interface{}{
						"request_header": []interface{}{
							map[string]interface{}{"operator": "STARTS_WITH", "value": "x-"},
						},
						"request_cookie":      []interface{}{},
						"request_uri":         []interface{}{},
						"request_query_param": []interface{}{},
						"target_rule_set":     "rce-stable",
						"target_rule_ids":     schema.NewSet(schema.HashString, []interface{}{"owasp-crs-v030001-id941120-xss"}),
					},
				},
			},
		},
		"rate_limit_options": []interface{}{
			map[string]interface{}{
				"rate_limit_threshold": []interface{}{
					map[string]interface{}{"count": int(100), "interval_sec": int(60)},
				},
				"ban_threshold": []interface{}{
					map[string]interface{}{"count": int(1000), "interval_sec": int(300)},
				},
				"conform_action":      "allow",
				"exceed_action":       "deny(403)",
				"enforce_on_key":      "IP",
				"enforce_on_key_name": "",
				"enforce_on_key_configs": []interface{}{
					map[string]interface{}{"enforce_on_key_type": "IP", "enforce_on_key_name": ""},
				},
				"ban_duration_sec": int(600),
				"exceed_redirect_options": []interface{}{
					map[string]interface{}{"type": "EXTERNAL_302", "target": "https://example.com"},
				},
			},
		},
		"redirect_options": []interface{}{
			map[string]interface{}{"type": "EXTERNAL_302", "target": "https://example.com"},
		},
		"header_action": []interface{}{
			map[string]interface{}{
				"request_headers_to_adds": []interface{}{
					map[string]interface{}{"header_name": "X-Goog-Test", "header_value": "true"},
				},
			},
		},
	}
}

// nestedMap descends into a []interface{}-wrapped block (MaxItems: 1) at each key in path
// and returns the innermost map, so tests can mutate a single deeply-nested leaf.
func nestedMap(t *testing.T, rule map[string]interface{}, path ...string) map[string]interface{} {
	t.Helper()
	cur := rule
	for _, key := range path {
		list, ok := cur[key].([]interface{})
		if !ok || len(list) == 0 {
			t.Fatalf("path element %q is not a non-empty block", key)
		}
		cur, ok = list[0].(map[string]interface{})
		if !ok {
			t.Fatalf("path element %q does not contain a map", key)
		}
	}
	return cur
}

// TestUnitComputeSecurityPolicyRuleHash_detectsFieldEdits asserts that editing any
// behavior-defining field changes the rule hash. Before the fix, editing a field other
// than priority/action left the hash unchanged, so schema.Set treated the edited rule as
// identical and d.HasChange("rule") missed the change
// (hashicorp/terraform-provider-google#27936).
func TestUnitComputeSecurityPolicyRuleHash_detectsFieldEdits(t *testing.T) {
	t.Parallel()

	cases := map[string]func(rule map[string]interface{}){
		"action":      func(r map[string]interface{}) { r["action"] = "deny(403)" },
		"preview":     func(r map[string]interface{}) { r["preview"] = true },
		"description": func(r map[string]interface{}) { r["description"] = "changed" },
		"match.versioned_expr": func(r map[string]interface{}) {
			nestedMap(t, r, "match")["versioned_expr"] = ""
		},
		"match.expr.expression": func(r map[string]interface{}) {
			nestedMap(t, r, "match", "expr")["expression"] = "origin.region_code == 'CA'"
		},
		"match.config.src_ip_ranges": func(r map[string]interface{}) {
			nestedMap(t, r, "match", "config")["src_ip_ranges"] = schema.NewSet(schema.HashString, []interface{}{"10.0.0.0/8"})
		},
		"match.expr_options.recaptcha_options.action_token_site_keys": func(r map[string]interface{}) {
			nestedMap(t, r, "match", "expr_options", "recaptcha_options")["action_token_site_keys"] = []interface{}{"action-key-2"}
		},
		"preconfigured_waf_config.exclusion.target_rule_set": func(r map[string]interface{}) {
			nestedMap(t, r, "preconfigured_waf_config", "exclusion")["target_rule_set"] = "sqli-stable"
		},
		"preconfigured_waf_config.exclusion.target_rule_ids": func(r map[string]interface{}) {
			nestedMap(t, r, "preconfigured_waf_config", "exclusion")["target_rule_ids"] = schema.NewSet(schema.HashString, []interface{}{"owasp-crs-v030001-id941130-xss"})
		},
		"preconfigured_waf_config.exclusion.request_header.value": func(r map[string]interface{}) {
			nestedMap(t, r, "preconfigured_waf_config", "exclusion", "request_header")["value"] = "y-"
		},
		// The reviewer's specific concern: the rate limit threshold and its count / interval.
		"rate_limit_options.rate_limit_threshold.count": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options", "rate_limit_threshold")["count"] = int(200)
		},
		"rate_limit_options.rate_limit_threshold.interval_sec": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options", "rate_limit_threshold")["interval_sec"] = int(120)
		},
		"rate_limit_options.ban_threshold.count": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options", "ban_threshold")["count"] = int(2000)
		},
		"rate_limit_options.conform_action": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options")["conform_action"] = "deny"
		},
		"rate_limit_options.exceed_action": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options")["exceed_action"] = "deny(429)"
		},
		"rate_limit_options.enforce_on_key": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options")["enforce_on_key"] = "HTTP_HEADER"
		},
		"rate_limit_options.enforce_on_key_name": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options")["enforce_on_key_name"] = "X-Custom"
		},
		"rate_limit_options.enforce_on_key_configs.enforce_on_key_type": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options", "enforce_on_key_configs")["enforce_on_key_type"] = "HTTP_COOKIE"
		},
		"rate_limit_options.ban_duration_sec": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options")["ban_duration_sec"] = int(1200)
		},
		"rate_limit_options.exceed_redirect_options.target": func(r map[string]interface{}) {
			nestedMap(t, r, "rate_limit_options", "exceed_redirect_options")["target"] = "https://changed.example.com"
		},
		"redirect_options.type": func(r map[string]interface{}) {
			nestedMap(t, r, "redirect_options")["type"] = "GOOGLE_RECAPTCHA"
		},
		"redirect_options.target": func(r map[string]interface{}) {
			nestedMap(t, r, "redirect_options")["target"] = "https://changed.example.com"
		},
		"header_action.request_headers_to_adds.header_name": func(r map[string]interface{}) {
			nestedMap(t, r, "header_action", "request_headers_to_adds")["header_name"] = "X-Goog-Other"
		},
		"header_action.request_headers_to_adds.header_value": func(r map[string]interface{}) {
			nestedMap(t, r, "header_action", "request_headers_to_adds")["header_value"] = "false"
		},
	}

	base := resourceComputeSecurityPolicyRuleHash(newSecurityPolicyRuleForHash())

	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			rule := newSecurityPolicyRuleForHash()
			mutate(rule)
			if got := resourceComputeSecurityPolicyRuleHash(rule); got == base {
				t.Errorf("editing %s did not change the rule hash (%d); the edit would be silently skipped", name, got)
			}
		})
	}
}

// TestUnitComputeSecurityPolicyRuleHash_stable asserts that two independently-built but
// semantically identical rules hash identically. This guards against churn / spurious
// recreation (hashicorp/terraform-provider-google#16882), including set members supplied
// in different orders.
func TestUnitComputeSecurityPolicyRuleHash_stable(t *testing.T) {
	t.Parallel()

	a := newSecurityPolicyRuleForHash()

	b := newSecurityPolicyRuleForHash()
	// Same IP ranges, different insertion order: a set hash must be order-independent.
	nestedMap(t, b, "match", "config")["src_ip_ranges"] = schema.NewSet(schema.HashString, []interface{}{"192.168.0.0/16", "10.0.0.0/8"})

	if resourceComputeSecurityPolicyRuleHash(a) != resourceComputeSecurityPolicyRuleHash(b) {
		t.Errorf("semantically identical rules produced different hashes: %d != %d",
			resourceComputeSecurityPolicyRuleHash(a), resourceComputeSecurityPolicyRuleHash(b))
	}
}
