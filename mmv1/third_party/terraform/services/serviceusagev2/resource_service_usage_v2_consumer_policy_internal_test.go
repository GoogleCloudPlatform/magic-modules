package serviceusagev2

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"google.golang.org/api/googleapi"
)

func Test_getTransitiveDependents(t *testing.T) {
	type args struct {
		dependents map[string][]string
		service    string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "not existing",
			args: args{
				map[string][]string{"a": {"b"}},
				"b",
			},
			want: nil,
		},
		{
			name: "simple",
			args: args{
				map[string][]string{"a": {"b"}},
				"a",
			},
			want: []string{"b"},
		},
		{
			name: "transitive",
			args: args{
				map[string][]string{"a": {"b", "c"}, "b": {"c", "d"}, "c": {"e", "f"}},
				"a",
			},
			want: []string{"b", "c", "d", "e", "f"},
		},
		{
			name: "transitive circular",
			args: args{
				map[string][]string{"a": {"b", "c"}, "b": {"c", "d"}, "c": {"a", "f"}},
				"a",
			},
			want: []string{"b", "c", "d", "f"},
		},
		{
			name: "circular",
			args: args{
				map[string][]string{"a": {"b"}, "b": {"c", "d"}, "c": {"a"}},
				"a",
			},
			want: []string{"b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTransitiveDependents(tt.args.dependents, tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTransitiveDependents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getServicesFromString(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "getting groups",
			args: args{"The services pubsub.googleapis.com,sql.googleapis.com have usage in the last 30 days or were enabled in the past 3 days. When the policy is removed from 'organizations/949638339900', if another policy does not enable the services then clients in 'organizations/949638339900' will lose API access and any remaining resources from the corresponding services will be deleted. Please specify force if you want to proceed with the destructive policy change.\\nHelp Token: AQAb6BwmwbBxTs_s1vA7H30XlWMWY3hWLrCNDdUjlVSgdoTALMWLCvzrXc-oinqSD-cLL6KFAR584uf7WUQlhGqf7tIhrNjzbFfGc65LPNenS8S_\""},
			want: []string{"pubsub.googleapis.com", "sql.googleapis.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getServicesFromString(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getErrorMessageForServicesInUse(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "build message",
			args: args{"The services pubsub.googleapis.com,sql.googleapis.com have usage in the last 30 days or were enabled in the past 3 days. When the policy is removed from 'organizations/949638339900', if another policy does not enable the services then clients in 'organizations/949638339900' will lose API access and any remaining resources from the corresponding services will be deleted. Please specify force if you want to proceed with the destructive policy change.\\nHelp Token: AQAb6BwmwbBxTs_s1vA7H30XlWMWY3hWLrCNDdUjlVSgdoTALMWLCvzrXc-oinqSD-cLL6KFAR584uf7WUQlhGqf7tIhrNjzbFfGc65LPNenS8S_\""},
			want: "The service{s} pubsub.googleapis.com, sql.googleapis.com has been used in the last 30 days or was enabled in the past 3 days. If you still wish to remove the service{s}, please set the check_usage_on_remove flag to false to proceed.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getErrorMessageForServicesInUse(tt.args.msg); got != tt.want {
				t.Errorf("getErrorMessageForServicesInUse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDependentServicesMap(t *testing.T) {
	type args struct {
		url                 string
		services            []string
		apiDependenciesFunc ApiCallFn
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{name: "get_dependencies",
			args: args{
				url:      "https://serviceusage.googleapis.com/v2alpha/organizations/949638339900/consumerPolicies/default",
				services: []string{"a", "b", "c"},
				apiDependenciesFunc: func(url string) (map[string]interface{}, error) {

					switch {
					case strings.HasSuffix(url, "/a/groups/dependencies/descendantServices"):
						return map[string]interface{}{
							"services": []map[string]interface{}{{"serviceName": "ab"}}}, nil

					case strings.HasSuffix(url, "/b/groups/dependencies/descendantServices"):
						return map[string]interface{}{
							"services":      []map[string]interface{}{{"serviceName": "ba"}, {"serviceName": "bb"}},
							"nextPageToken": "token1"}, nil

					case strings.HasSuffix(url, "/b/groups/dependencies/descendantServices?pageToken=token1"):
						return map[string]interface{}{
							"services": []map[string]interface{}{{"serviceName": "bc"}}}, nil
					case strings.HasSuffix(url, "/c/groups/dependencies/descendantServices"):
						return nil, &googleapi.Error{Details: []interface{}{map[string]interface{}{"reason": ApiErrorSuGroupNotFound}}}
					}
					return nil, nil

				},
			},
			want:    map[string][]string{"a": {"ab"}, "b": {"ba", "bb", "bc"}, "c": nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDependentServicesMap(tt.args.url, tt.args.services, tt.args.apiDependenciesFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDependentServicesMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDependentServicesMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateDependencies(t *testing.T) {
	type args struct {
		oldServices          []string
		newServices          []string
		dependentServicesMap map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "validate_empty",
			args: args{
				oldServices:          nil,
				newServices:          nil,
				dependentServicesMap: nil,
			},
			wantErr: false,
		},
		{
			name: "validate_ok",
			args: args{
				oldServices:          []string{"a"},
				newServices:          []string{"a", "b"},
				dependentServicesMap: map[string][]string{"b": {"a"}, "a": {}},
			},
			wantErr: false,
		},
		{
			name: "validate_missing_dependency",
			args: args{
				oldServices:          []string{"a"},
				newServices:          []string{"b"},
				dependentServicesMap: map[string][]string{"b": {"a"}},
			},
			wantErr: true,
			errMsg:  "There are additional services for which all necessary dependencies haven't been added\\. Please add these missing dependencies:[\\s\\S]*Added service: \\[\"b\"][\\s\\S]*Missing dependencies: \\[\"a\"][\\s\\S]*If you don't want to validate dependencies, set validate_dependencies to false to override\\.",
		},
		{
			name: "validate_removed_dependency",
			args: args{
				oldServices:          []string{"a", "b"},
				newServices:          []string{"a"},
				dependentServicesMap: map[string][]string{"a": {"b"}},
			},
			wantErr: true,
			errMsg:  "There are existing services in configuration which depend on the services to be removed\\. Please remove existing dependent services:[\\s\\S]*Removed service: \\[\"b\"][\\s\\S]*Existing dependents: \\[\"a\"][\\s\\S]*If you don't want to validate dependencies, set validate_dependencies to false to override\\.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDependencies(tt.args.oldServices, tt.args.newServices, tt.args.dependentServicesMap); ((err != nil) != tt.wantErr) || (tt.errMsg != "" && !regexp.MustCompile(tt.errMsg).MatchString(err.Error())) {
				t.Errorf("validateDependencies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
