# Magic Modules Template SDK

This document describes in detail the resources available to build a new
resource using Magic Module template system.


## Variables

### `@api` (`Api::Product`)

The @api variable contains a reference to the product being generated. The
product is constructed from `api.yaml` ([example][api-yaml-example]) and it is
either a `Api::Product` or a descendant of it.

-   `@api.name`: The full name of the product: `Google Compute Engine`
-   `@api.prefix`: The prefix to uniquely identify the types: `gcompute`
-   `@api.base_url`: The base URL for the service API endpoint:
    `https://www.googleapis.com/compute/v1/`. This should not be set in the
    API definition file, but rather by the generator.
-   `@api.scopes` (list): The list of permission scopes available for the
    service: `https://www.googleapis.com/auth/compute`
-   `@api.versions` (list of `Api::Product::Version`): The API versions that
    can be called for this product.

### `object` (`Api::Resource`)

The `object` variable holds a reference to the type being compiled. From this
variable you have access to everything defined in the object: properties,
parameters, exports, description, etc.

A non-exhaustive list of the most common used properties:

-   `object.base_url` (URI: relative to `@api.base_url` or absolute)
-   `object.description` (string)
-   `object.exports` (list of `string` | `Api::Type::FetchedExternal`)
-   `object.kind` (GCP kind, e.g. `compute#disk`)
-   `object.name` (string)
-   `object.parameters` (list of `Api::Type`)
-   `object.properties` (list of `Api::Type`)
-   `object.references`: Documentation references
    *   `object.references.guides` (hash)
        -   name: The title of the link
        -   value: The URL to navigate on click
    *   `object.references.api` (list of URI)
        -   name: The title of the link
        -   value: The URL to navigate on click

#### Helper Functions

-   `object.all_user_properties`
-   `object.out_name`

For complete list of functions and properties refer to `api/resource.rb`
documentation.

### `@config` (`Provider::Config`)

The `@config` variable holds the provider specific overrides for the object. The
`@config` variable is at the object level (same level as `object` variable).

Templates usually use `@config` to fetch specific one-off product specific
behavior. For example in Cloud DNS updating a Resource Record Set requires
creating a transaction on another object. _This behavior is not REST compliant,
and therefore we cannot use the expected REST create call._ To overcome this
non-standard REST behavior we define an override of the create function
([example][cloud-dns-override-create]) to accommodate it:

    objects: !ruby/object:Api::Resource::HashArray
      ...
      ResourceRecordSet:
        ...
        create: |
          change = create_change nil, updated_record, @resource
          change_id = change['id'].to_i
          debug("created for transaction '#{change_id}' to complete")
          wait_for_change_to_complete change_id, @resource \
            if change['status'] == 'pending'

## Functions: Provider Specific

Magic Modules uses reflection and aligns the ERB stack with the provider object.
That means all methods described in the provider class being executed are
readily available for use within a template file being compiled.

In this document we list the most commonly used functions from `Provider::Core`.
Refer to the provider documentation for a complete list of functions available.

For example the function [`effective_properties()`][tf-effective-properties]
defined in `provider/terraform.rb` can be used directly in the template
[`templates/terraform/resource.erb`][tf-effective-properties-usage]

## Functions: Code Generation

### `compile()`

    def compile(file, caller_frame=1)

Compiles a ERB style template and outputs the result as string. During
compilation the file will have access to all variables and functions described
in this document, as well as any provider specific public functions defined in
the corresponding provider class.

Example:

    <%= compile('templates/async.erb') -%>

It is okay to pair `compile()` with formatting functions such as
[`indent()`](#indent), [`lines()`](#lines) and others.

> Do not change the `caller_frame=1`. It is reserved for future use.

### `compile_if()`

    def compile_if(config, node)

TODO(nelsona): Document

Example:

    <%= lines(indent(compile_if(tests, %w[expectation_helpers]), 2), 1) -%>

### `include()`

    def include(file)

Includes a file _verbatim_ (does not resolve or compile its contents) and
returns as a string.

This method is similar to [`compile()`](#compile) with the difference that no
processing happens on the file. It is used to clearly specify that the file is
not to be changed, e.g. when the file itself is written in Ruby or ERB.

Example:

    <%= indent(include('products/pubsub/test.yaml'), 2) %>

### `emit_method()`

    def emit_method(name, args, code, file_name, opts = {})

Emits a method definition for the function named `name`.

Arguments:

-   `name (string)`: the name of the method
-   `args (list)`: the list of parameters for the function
-   `code (list|string)`: the body of the function
-   `file_name`: the name of the file being processed (used to add style
    guidelines)
-   `opts (hash)`: a name/value pair of options

Example (excerpt from `templates/puppet/resource.erb`):

    <%=
      lines(indent(emit_method('self.resource_to_hash', %w[resource], r2h_code,
                               file_relative), 2), 1)
    -%>

### `emit_link()`

    def emit_link(name, url, emit_self, extra_data=false)

Generates a function that outputs a link to the resource (or its collection).

Arguments:

-   `name`: the name of the function to generate
-   `url`: the URL of the resource
-   `emit_self`: if `true` the function is static

Example (excerpt from `templates/ansible/resource.erb`):

    <%= lines(emit_link('collection', collection_url(object))) -%>

### `emit_requires()`

    def emit_requires(requires)

Generate the list of `require` (`import` or `include` depending on the target
language) for a specific template. This function is useful to because during
compilation of various template fragments each fragment may need a specific set
of libraries to operate. The list will overlap and we want each fragment to be
self-contained, so it can be reused.

Therefore each fragment declares which includes it requires, and the
`emit_requires()` function collects them all, sort them according to the style
guide and output to the template.

Arguments:

-   `requires (list)`: the list of requires collected

Example (excerpt from `templates/puppet/bolt~task.rb.erb`):

    <%= lines(emit_requires(task.requires)) -%>

Example (excerpt from `templates/chef/resource.erb`):

    <%
      requires = generate_requires(object.all_user_properties)
      requires << 'chef/resource'
      requires << 'google/hash_utils'
      requires << emit_google_lib(binding, Compile::Libraries::NETWORK, 'get')
      unless object.readonly
        requires << emit_google_lib(binding, Compile::Libraries::NETWORK, 'delete')
        requires << emit_google_lib(binding, Compile::Libraries::NETWORK, 'post')
        requires << emit_google_lib(binding, Compile::Libraries::NETWORK, 'put')
      end
    -%>
    <%= lines(emit_requires(requires)) -%>

### `self_link_url()`

    def self_link_url(resource)

Generates the uniquely identifying self link for the object.

Arguments:

-   `resource`: the object to generate the self link

Example (excerpt from `templates/terraform/resource.erb`):

    url, err := replaceVars(d, config, "<%= self_link_url(object) %>")

### `collection_url()`

    def collection_url(resource)

Generates the URL that hosts the collection of objects of the same kind as the
resource.

Arguments:

-   `resource`: the object to generate the collection url

Example (excerpt from `templates/terraform/resource.erb`):

    url, err := replaceVars(d, config, "<%= collection_url(object) %>")

### `emit_google_lib()`

    def emit_google_lib(ctx, lib, file)

TODO(nelsona): Document

## Functions: Formatting

### `indent()`

    def indent(text, amount, filler = ' ')

Indents code (or code block) by N characters specified by `filler`. If the input
is empty, it returns empty string. It is very versatile and supports:

-   Single line or multi line (list) input
-   Preserving inline indentation
-   Nesting
-   Pairing with [`lines()`](#lines) or its lookalikes to generate indented
    optionally placed code

Arguments:

-   `text (list | string)`: the text to indent. it can be a string or a list of
    strings.
-   `amount (integer)`: the number of fillers to indent by
-   `filler (string)`: the block to add on each indent (default to space)

Example (pseudo code):

    indent([
      ....
      .... everything else is +2
      ....
      ident([
           ...
           ... everything here would be effectively indented +4
           ...
        ], 2)
      ], 2)

Example (excerpt from `templates/puppet/property/nested_object.rb.erb`):

    [
      "@#{prop.out_name} =",
      indent([
        "#{parser}(",
        indent("args['#{source}']", 2),
        ')'
      ], 2)
    ]

### `indent_list()`

    def indent_list(text, amount, filler = ' ')

Similar to [`indent()`](#indent) but will produce a list of items as output. A
list of items has the following properties:

-   Indented by N (`amount`) instances of `filler`
-   Every line item (except the last) is terminated with a `,`

Example (excerpt from `templates/puppet/resource.erb`):

    f2h_code = [
        '{',
        indent_list(assigns, 2),
        '}.reject { |_, v| v.nil? }',
    ]

Will produce (excerpt from Puppet Compute module
`lib/puppet/provider/gcompute_address/google.rb`):

    {
      address: Google::Compute::Property::String.api_munge(fetch['address']),
      creation_timestamp:
        Google::Compute::Property::Time.api_munge(fetch['creationTimestamp']),
      description:
        Google::Compute::Property::String.api_munge(fetch['description']),
      id: Google::Compute::Property::Integer.api_munge(fetch['id']),
      name: Google::Compute::Property::String.api_munge(fetch['name']),
      users: Google::Compute::Property::StringArray.api_munge(fetch['users'])
    }.reject { |_, v| v.nil? }

### `indent_array()`

    def indent_array(text, spaces)

Similar to [`indent()`](#indent) but it expects an array of strings as input.

### `bullet_line()` + `bullet_lines()`

TODO(nelsona): Document

### `lines()`

    def lines(code, number=0)

Ensures code has a terminating newline and that N extra lines are at the end of
the text.

> This is one of the most useful functions of the SDK (along with
> [`indent()`](#indent). It gives the template writer peace of mind that no
> matter how the functions they are calling are laid out, some specific newline
> guarantees are enforced.

This is useful when:

-   You're emitting code that's optional or conditional to product specific
    settings. This avoids lots of blank lines as artifact of pieces of code
    being removed
-   You need to enforce style guide regarding spaces (e.g. Google Python code
    style requires 2 blank lines between the imports and the class name, or
    between methods)

Example (excerpt from `templates/ansible/async.erb`):

    <%= lines(emit_link('async_op_url', async_operation_url(object), true), 2) -%>

In the example above it will:

-   Add a newline to the output if it does not end with `\n`
-   Add extra 2 newlines after the output to conform with Python code style

When called without the number of lines (default=0) it will ensure that the
output has **at most** 1 newline on the output.

Let's take this function:

    <%
      def my_unpredictable_newliner():
        "my line" + ("\n" * rand(10))
      end
    -%>

By calling with:

    ---
    <%= lines(my_unpredictable_newliner()) -%>
    <%= lines(my_unpredictable_newliner()) -%>
    <%= lines(my_unpredictable_newliner()) -%>
    <%= lines(my_unpredictable_newliner()) -%>
    ---

It will produce:

    ---
    my line
    my line
    my line
    my line
    ---

Whereas if you call with an argument:

    ---
    <%= lines(my_unpredictable_newliner(), 1) -%>
    <%= lines(my_unpredictable_newliner(), 1) -%>
    <%= lines(my_unpredictable_newliner(), 1) -%>
    <%= lines(my_unpredictable_newliner(), 1) -%>
    ---

It will produce:

    ---
    my line

    my line

    my line

    my line

    ---

### `lines_before()`

    def lines_before(code, number=0)

Similar to [`lines()`](#lines) except that it enforces blank lines **before**
the text.

### `format()`

    def format(sources, indent=0, start_indent=0,
               max_columns=DEFAULT_FORMAT_OPTIONS[:max_columns])

When generating artifacts (code, documentation, examples) we rely on templates.
At the same time we want the generated code to look good and follow standards
set by the community.

The main reason formatting issues happen is due to the fact templates will have
variables, loops and other artifacts that will depend on developer input (e.g.
property names will vary in length, some properties will be optional, etc).

The secret to achieve well formatted code is to use the `format()` function.
`format()` will allow you to specify various functionally equivalent,
alternative format outputs, and give it contraints, both in size and
indentation. Magic Modules will determine the best alternative for all the input
once evaluated and emit it.

    <%=
      lines(format([
        ["#{first_name} and #{last_name} are < 80 characters"],
        [
          first_name,
          "is quite long, so let's put it on it's own line #{last_name}"
        ],
        [
          first_name,
          "is quite long, as well as last name, so let's put it next",
          last_name
        ]
      ]))
    -%>

We use `format()` in Puppet to create a property function definition. Whenever
it fits in 80 characters it will be put in a single line:

    newparam(:zone, parent: Google::Container::Property::String) do
      desc 'The zone where the cluster is deployed'
    end

But when it did not fit it is broken into multiple lines according to
alternative formats:

    newproperty(:master_auth,
                parent: Google::Container::Property::ClusterMasterAuth) do
      desc 'The authentication information for accessing the master endpoint.'
    end

To achieve such behavior without various conditionals, we describe the possible
outcomes:

    <%=
      new_property_len = 'newproperty('.length
      format([
        ["newproperty(:#{p.out_name}, parent: #{p.property_type}) do"],
        [
          "newproperty(:#{p.out_name},",
          indent("parent: #{p.property_type}) do", new_property_len)
        ],
      ], 2)
    %>
    <%= format_description(p, 4, 'desc', p.output ? "(output only)" : "") %>
    <%= property_body(p) -%>
      end

### `format_description()`

    def format_description(object, spaces, container, suffix = '')

Formats a long description to fit within a specific `container`. This function
is similar to [`wrap_field()`](#wrap-field) but it also emits the container for
the text, generating **code compatible** output.

Arguments:

-   `object`: the text to format
-   `spaces`: the number of leading spaces to indent by
-   `container`: the container what will wrap the text
-   `suffix`: an optional suffix to append to the text

Example (excerpt from `templates/puppet/type.erb`):

    <%= format_description(p, 4, 'desc', p.output ? "(output only)" : "") %>

It will generate:

    desc 'Creation timestamp in RFC3339 text format. (output only)'

If the text does not fit in the maximum allotted size, it will wrap in a heredoc
block:

    desc <<-DESC
      A gauth_credential name to be used to authenticate with Google Cloud
      Platform.
    DESC

### `wrap_field()`

    def wrap_field(field, spaces)

Fits a long running field text, breaking into word boundaries, to fit in the
`DEFAULT_FORMAT_OPTIONS[:max_columns]` available character limit.

Optionally it will indent the output by N `spaces`.

Example (excerpt from `templates/chef/README.md.erb`):

    <%= lines(wrap_field("Converges the `#{object.out_name}` resource into the final
    state described within the block. If the resource does not exist, Chef will
    attempt to create it.", 0)) -%>

### `comment_block()`

    def comment_block(text, language)

Formats a block for a given language, adding the necessary markup around to make
the multi line text block a valid comment in the language.

Arguments:

-   `text` (list | string): the text to be formatted. It can be a string with
    "\n" or an array. If text is an array, each entry of the array will
    correspond to 1 line.
-   `lang` (symbol): The name of the target language

Currently supported target languages:

-   `:chef`: Chef artifact (cookbook, recipe, control file)
-   `:gemfile`: Gem file
-   `:git`: Git file (.gitignore)
-   `:go`: Go source file
-   `:html`: HTML file
-   `:markdown`: Markdown file
-   `:puppet`: Puppet manifest
-   `:python`: Python source file
-   `:ruby`: Ruby source file
-   `:yaml`: YAML file

## Functions: Language Specific

### `quote_string()` + `unquote_string()` [Ruby, Python]

    def quote_string(string)
    def unquote_string (string)

Adds quotes or double quotes depending on the input. It is common for dynamic
languages -- such as Python and Ruby -- to use either single (') or double (")
quotes to wrap a string. Some code styles will prefer one type over the other:
At the same time it is preferred to quote the string with the other style when
the string contents contain the preferred marker to avoid escaping them:

    'it\'s hard for Matt\'s team to find string end'
    "it will be \"funny but not funny per se\" if this was not \"tragic\""

Versus:

    "it's hard for Matt's team to find string end"
    'it will be "funny but not funny per se" if this was not "tragic"'

This function chooses the most appropriate quotation for the string using these
various criteria.

Example (excerpt from `templates/ansible/resource.erb`):

    auth = GcpSession(module, <%= quote_string(prod_name) -%>)

## Testing

### `Provider::TestMatrix`

We like unit tests. We want to cover all possible cases. You can lay down all
the cases a multi-dimensional matrix. For example to completely test all
possible Puppet expected object transitions, you need to ensure you cover:

-   ensure: present | absent
-   exists: yes | no
-   changed: yes | no
-   has name: yes | no

For each of the combinations above there's an expected action and result
associated with it: a successful and failed test cases. For example:

-   present > exists > changed > has_name ==> success case
-   present > exists > changed > has_name ==> failure case

To make sure we do not miss any of these tests we rely on the
`Provider::TestMatrix` object. It tracks all dimensions that need to be
exercised and fails the compilation if one or more combinations are missing.

To use `TestMatrix` you:

1.  Define the combination matrix of all possible allowed states (nested
    matrixes are also supported)
2.  Wrap the test case that handles a specific set of dimension with `::push()`
    and `::pop()` methods, signaling you covered that combination

Example (test matrix definition) (excerpt from
`templates/puppet/provider_spec.erb`):

    test_matrix = Provider::TestMatrix.new(template, object, self,
      present: {
        exists: {
          changes: { # flush
            no_name: [:pass, :fail],
            has_name: [:pass, :fail]
          },
          no_change: { # no action
            no_name: [:pass, :fail],
            has_name: [:pass, :fail]
          }
        },
        missing: { # create
          # changes == ignore
          no_name: [:pass, :fail],
          has_name: [:pass, :fail]
        }
      },
      absent: {
        exists: { # delete
          # changes == ignore
          no_name: [:pass, :fail],
          has_name: [:pass, :fail]
        },
        missing: { # no action
          # changes == ignore
          no_name: [:pass, :fail],
          has_name: [:pass, :fail]
        }
      }
    )

Also note that the matrix neither need to be symmetric or solid. In the example
above, when testing `absent > exists` we do not care for `changed`, so we
skipped that completely (while we cared for `changed` in the `present` part of
the matrix).

Example (test matrix usage) (excerpt from `templates/puppet/provider_spec.erb`):

    <%= test_matrix.push(:present, :exists, :no_change, :has_name, :pass) -%>
    <%=
      # Ensure no changes happened
      name = true
      compile('templates/puppet/test~present~no_changes.erb')
    -%>
    <%= test_matrix.pop(:present, :exists, :no_change, :has_name, :pass) -%>

If any test cases are missing the compiler will fail:

    The following tests are missing from the matrix: [
      [:present, :missing, :no_name, :fail],
      [:present, :missing, :has_name, :fail]
    ]

## Functions: Ruby specific

### `emit_rubocop()`

TODO(nelsona): Document

## Miscellaneous Helper Objects

-   `google/hash_utils.rb`: Hash helper functions (e.g. navigate)
-   `google/integer_utils.rb`: Integer helper functions
-   `google/logger.rb`: Logging facilities
-   `google/object_store.rb`: In-memory object store (to keep gloal state)
-   `google/string_utils.rb`: String helper functions (e.g. underscore,
    camelize)
-   `google/yaml_validator.rb`: YAML parser and type enforcer (base for all
    `api/` objects)

## Generating License & Copyright notices

Whenever we're dealing with licenses we want to add the correct copyright notice
and year to the files. That gets tricky when the file is a template. If you add
a license year to the template, and use that template to generate new artifacts,
the output artifacts will be tagged with the wrong year: it will feature the
year the template was created instead of the year the artifact was created.

To remediate this Magic Modules deals with copyright notices in a different way:

-   The copyright notices on templates belongs to the template only
-   The effective copyright notice will be added by compiling a copyright
    template

To ensure copyright notices do not bleed from templates Magic Modules strips
from templates any copyright notices (blocks of comments starting at line 1 of
the template being compiled if it contains the word _Copyright_ anywhere in the
first line.

The template owner decides the proper place to add the copyright notice (e.g. in
a bash script it would be after the `#!/bin/bash` that is required to be in
line #1 by compiling the license template:

    <%= compile 'templates/license.erb' -%>

To ensure no doubts exist on the readers' minds, we recommend that you mark the
copyright notice on the template with a ERB block:

    <%# The license inside this block applies to this file
    # Copyright 2017 Google Inc.
    ...
    ...
    ...
    -%>

### Determining the generated artifact "effective year"

In the past the template has a convenient "replace with current year" for the
copyright notice:

    # Copyright <%= Time.now.year -%> Google Inc.

However when you're doing multi year development that can be tricky as updates
to the file would rewrite the year (incorrectly). Instead Magic Modules will
determine the "effective year" range for the file and update the file
accordingly, with the template:

    # Copyright <%= effective_copyright_year(out_file) -%> Google Inc.

## Template Files

-   `templates/CONTRIBUTING.md.erb`: Generates the `CONTRIBUTING.md`
    ([example][contributing-example]) with a complete list of generated
    artifacts tracked by Magic Modules
-   `templates/autogen_notice.erb`: Attached the "automatically generated"
    notice to an artifact
-   `templates/async.yaml.erb`: Skeleton for `async:` block of a product
    `api.yaml` file
-   `templates/transport.erb`: Product independent Ruby implementation of a REST
    transport
-   `templates/CHANGELOG.md.erb`: Generates the `CHANGELOG.md`
    ([example][changelog-example]) with project changes
-   `templates/network_mocks.erb`: Generates the network mocks for unit testing
    hermetic run
-   `templates/license.erb`: Generates (and year tags) project LICENSE
-   `templates/dot~rubocop~root.yml`: Rubocop definitions for generated projects
-   `templates/dot~rubocop~spec.yml`: Rubocop definitions for generated tests

... and many more. Check the [`templates/`][templates]

[mm-styleguide]: https://github.com/GoogleCloudPlatform/magic-modules/blob/master/CONTRIBUTING.md#templating-style-guide
[ruby-style-guide]: https://github.com/bbatsov/ruby-style-guide
[format-example]: https://github.com/GoogleCloudPlatform/magic-modules/blob/a15064abaec6bac31e4ddf37119109b6701b7285/templates/puppet/type.erb#L143
[templates]: https://github.com/GoogleCloudPlatform/magic-modules/blob/master/templates
[changelog-example]: https://github.com/GoogleCloudPlatform/puppet-google-compute/blob/master/CHANGELOG.md
[contributing-example]:  https://github.com/GoogleCloudPlatform/puppet-google-compute/blob/master/CONTRIBUTING.md
[api-yaml-example]: https://github.com/GoogleCloudPlatform/magic-modules/blob/master/product/compute/api.yaml
[cloud-dns-override-create]: https://github.com/GoogleCloudPlatform/magic-modules/blob/a15064abaec6bac31e4ddf37119109b6701b7285/products/dns/puppet.yaml#L40
[tf-effective-properties]: https://github.com/GoogleCloudPlatform/magic-modules/blob/a15064abaec6bac31e4ddf37119109b6701b7285/provider/terraform.rb#L99
[tf-effective-properties-usage]: https://github.com/GoogleCloudPlatform/magic-modules/blob/a15064abaec6bac31e4ddf37119109b6701b7285/templates/terraform/resource.erb#L21
[ruby-enumerable]: https://ruby-doc.org/core-2.5.0/Enumerable.html
