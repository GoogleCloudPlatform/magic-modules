Hi there, I'm the Modular magician. I've detected the following information about your changes:

## Diff report
{{ $diffsLength := len .diffs }}{{if eq $diffsLength 0 }}
Your PR hasn't generated any diffs, but I'll let you know if a future commit does.
{{else}}
Your PR generated some diffs in downstreams - here they are.

{{range .diffs -}}
{{.repo.Title}}: [Diff](https://github.com/modular-magician/{{.repo.Name}}/compare/auto-pr-{{$.prNumber}}-old..auto-pr-{{$.prNumber}}) ({{.diffStats}})
{{end -}}
{{end -}}

{{- $breakingChangesLength := len .breakingChanges }}
{{- if gt $breakingChangesLength 0}}
## Breaking Change(s) Detected
The following breaking change(s) were detected within your pull request.

If you believe this detection to be incorrect please raise the concern with your reviewer.
If you intend to make this change you will need to wait for a [major release](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-major-number-increments) window.
An `override-breaking-change` label can be added to allow merging.
{{end}}
{{.missingTests}}