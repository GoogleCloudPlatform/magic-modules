# changelog-entry

`changelog-entry` is a command that will generate a changelog entry based on the information passed and the information retrieved from the Github repository.

The default changelog entry template is embedded from [`changelog-entry.tmpl`](changelog-entry.tmpl) but a path to a custom template can also can be passed as parameter.

The type parameter can be one of the following:
* bug
* note
* enhancement
* new-resource
* new-datasource
* deprecation
* breaking-change
* feature

## Usage

```sh
$ changelog-entry -type improvement -subcategory monitoring -description "optimize the monitoring endpoint to avoid losing logs when under high load"
```

If parameters are missing the command will prompt to fill them, the pull request number is optional and if not provided the command will try to guess it based on the current branch name and remote if the current directory is in a git repository.

## Output

``````markdown
```release-note:improvement
monitoring: optimize the monitoring endpoint to avoid losing logs when under high load
```
``````

Any failures will be logged to stderr. The entry will be written to a file named `{PR_NUMBER}.txt`, in the current directory unless an output directory is specified.
