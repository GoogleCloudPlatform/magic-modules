---
title: "Ruby to Go Migration"
weight: 10
---
# What has changed in the MMv1 Go migration

The Magic Modules code generator has been rewritten from Ruby to Go. For experienced contributors, this reference document lists what the expected changes are to the previous development workflow in Ruby.

## YAML changes

`.yaml` files within `mmv1/products` have had adjustments to the attribute typing. The initial Ruby lines `!ruby/object:Api::Type::<TYPE>` have been removed and replaced with a simpler `type: <TYPE>` line.

Old Ruby YAML
```yaml
- !ruby/object:Api::Type::String
  name: 'apiFieldName'
  description: |
    MULTILINE_FIELD_DESCRIPTION
```

New Go YAML
```yaml
- name: 'apiFieldName'
  type: String
  description: |
    MULTI_LINE_FIELD_DESCRIPTION
```

## Template `.erb` file changes

Template files have all been converted Embedded Ruby (ERB) templates to Go's [text/template](https://pkg.go.dev/text/template) format.
All `.erb` files are replaced with equivalent `.tmpl` files. The MMv1 resource objects are passed to the Go templates for referencing, similar to the previous Ruby templates. For the list of available variables and functions within the templates, please reference:

* [text/template standard library](https://pkg.go.dev/text/template)
* [mmv1/api/resource.go](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource.go)
* [mmv1/google/template_utils.go](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/google/template_utils.go)

### Common templating snippets

#### Version guards

Old Ruby template
```go
<% unless version == 'ga' -%>
// Go code here
<% end -%>
```

New Go template
```go
{{- if ne $.TargetVersionName "ga" }}
// Go code here
{{- end }}
```

#### Example `.tf.erb` variables

Old Ruby template `pubsub_topic_basic.tf.erb`
```hcl
resource "google_pubsub_topic" "<%= ctx[:primary_resource_id] %>" {
  name = "<%= ctx[:vars]['topic_name'] %>"

  labels = {
    foo = "bar"
  }

  message_retention_duration = "86600s"
}
```

New Go template `pubsub_topic_basic.tf.tmpl`
```hcl
resource "google_pubsub_topic" "{{$.PrimaryResourceId}}" {
  name = "{{index $.Vars "topic_name"}}"

  labels = {
    foo = "bar"
  }

  message_retention_duration = "86600s"
}
```

## Advanced: MMv1-specific generator command

Most contributors should use the make commands referenced in [make-commands](https://googlecloudplatform.github.io/magic-modules/reference/make-commands/) reference page to generate the downstream `google` and `google-beta` providers. The input for these commands have not changed, and have already been correctly switched over to use the new Go engine.

Some advanced contributors may be used to running the MMv1 generator commands. These commands have changed from Ruby's `bundle exec` to `go run`.

**These are not generally recommended to use**

Old Ruby MMv1 generator command in mmv1/:
```bash
bundle exec compiler -e terraform -o <output directory> -v <version> -f <MMv1 provider> -p <products/productfolder>
```

New Go MMv1 generator command in mmv1/:
```bash
go run . --output <output directory> --version <version> --provider <MMv1 provider>
```

## Advanced: MMv1 generator code locations

Most previous Ruby compiler code has parallel Go code placed the same file locations.
For example, the Go replacements for `mmv1/compiler.rb` and `mmv1/provider/terraform.rb` are `mmv1/main.go` and `mmv1/provider/terraform.go` respectively.
