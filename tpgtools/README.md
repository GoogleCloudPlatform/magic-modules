# tpgtools

`tpgtools` is the generator responsible for creating DCL-based resources in the
Terraform Google Provider (TPG). The DCL provides [OpenAPI schema objects](https://swagger.io/specification/#schema-object)
to describe the available types, and `tpgtools` uses those to construct
Terraform resource schemas.

## Usage

`tpgtools` expects to run targetting a "root service directory", a dir-of-dirs
where the child dirs contain OpenApi specs for resources such as the `api/` path
above. Additionally, overrides are expected in a similar structure (as seen in
the `overrides/` path. For example:

```
go run . --path "api" --overrides "overrides"
```

This will output the file contents of each resource to stdout, for fast
iteration. You can filter by service and resource to make it more useful:

```
go run . --path "api" --overrides "overrides" --service redis --resource instance
```

To persist the output(s) to files you can specify an output path. This is the
most familiar experience for MMv1 developers. For example:

```
go run . --path "api" --overrides "overrides" --output ~/tpg-fork
```

If generation fails, an error should be logged showing what went wrong. The raw
output will be returned, and line numbers (if available in the error) will
correspond to the line numbers in the output.

### Version

You can specify a version such as `beta` using the `--version`:


```
go run . --path "api" --overrides "overrides" --output ~/tpg-fork --version "beta"
```

### Accessory Code

To generate accessory code such as `serializarion`, you can specify the `--mode`:

```
go run . --path "api" --overrides "overrides" --output ~/some/dir --mode "serialization"
```

## Development

`tpgtools` builds resources using Go Templates, with the templates stored under
the `templates/` directory. They're fed the `Resource` type, which contains
resource-level metadata and a list of `Property` types which represent top-level
fields in a Terraform resource. Additionally, `Property`s can contain other
`Property` objects to create a nested structure of fields.

`main.go` parses the OpenAPI schemas you've selected into Go structs and then
transforms them into `Resource`/`Property` before running them through Go
Templates.

### Overrides

Overrides are specified per-resource, with a directory structure parallel to the
OpenAPI specs. Inside each resource file is an unordered array of overrides made
up of an override type (like `CUSTOM_DESCRIPTION` or `VIRTUAL_FIELD`) as well as
a field they affect (if a field is omitted, they affect the resource) and an
optional `details` object that will be parsed for additional metadata.

For example, override entries will look like the following:

```yaml
- type: CUSTOM_DESCRIPTION
  field: lifecycle.rule.condition.age
  details:
    description: Custom description here.
```
