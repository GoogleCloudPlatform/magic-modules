Hi there, I'm the Modular magician. I've detected the following information about your changes:

## Diff report
{{ $diffsLength := len .Diffs }}{{if eq $diffsLength 0 }}
Your PR hasn't generated any diffs, but I'll let you know if a future commit does.
{{else}}
Your PR generated some diffs in downstreams - here they are.

{{range .Diffs -}}
{{.Title}}: [Diff](https://github.com/modular-magician/{{.Repo}}/compare/auto-pr-{{$.PrNumber}}-old..auto-pr-{{$.PrNumber}}) ({{.ShortStat}})
{{end -}}
{{end -}}

{{- $breakingChangesLength := len .BreakingChanges }}
{{- if gt $breakingChangesLength 0}}
## Breaking Change(s) Detected

The following breaking change(s) were detected within your pull request.

{{- range .BreakingChanges}}
- {{.Message}} - [reference]({{.DocumentationReference}}){{end}}

If you believe this detection to be incorrect please raise the concern with your reviewer.
If you intend to make this change you will need to wait for a [major release](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-major-number-increments) window.
An `override-breaking-change` label can be added to allow merging.
{{end}}

{{if gt (len .MissingTests) 0}}
## Missing test report
Your PR includes resource fields which are not covered by any test.
{{ range $resourceName, $missingTestInfo := .MissingTests }}
Resource: `{{ $resourceName }}` ({{ len $missingTestInfo.Tests }} total tests)
Please add an acceptance test which includes these fields. The test should include the following:

```hcl
{{ $missingTestInfo.SuggestedTest }}

```

{{- end }}
{{end}}

{{- $errorsLength := len .Errors}}
{{- if gt $errorsLength 0}}
## Errors
{{range .Errors}}
{{.Title}}:
{{- range .Errors}}
- {{.}}{{end}}
{{end}}
{{- end -}}
