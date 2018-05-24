# Magic Modules Governance (a.k.a. "Golden Rules")

This document specified the basic principles behind the Magic Modules and its
"golden rules" that _shall not be infringed_ without a very strong (and
documented) reason.

In this doc the terms:

  - **must**: means something that _have to be followed at all times_<br/>
    (aside some specific exceptions detailed in the [Exceptions](#exceptions)
    section)
  - **should**: means something that have to be followed _as much as possible_.

Please refer to the top-level [README][readme] and [Philosophy][phylosophy]
before reading this guide.  All code-reviews of Magic Modules (regardless of
changes to core, providers, or products) must follow the guidelines outlined
here.


## Table of Contents

- [Development Principles](#development-principles)
- [Folder Locations & Usage](#folder-locations--usage)
- [Code Style](#code-style)
- [Templates](#templates)
- [Exceptions](#exceptions)
- [Ruby Best Practices & Style Guide](#ruby-best-practices--style-guide)


## Development Principles

  - All changes **must** be tested for:
      * Unit tests
      * Code style compliance
  - Changes to _Provider Independent_ code **must** be tested against **all**
    providers
  - Changes to _Product Indepedent_ code **must** be tested against **all**
    products
  - Changes to core Magic Modules features (which are both _Product_ and
    _Provider_ independent) **must** be tested against *all* products **and**
    **all** providers.
  - Code coverage **must** always stay >80% and **should** be >90%


### Object References & Self Links

Resource URL and other self link constructs **must not** be exposed to
customers. Use resource reference properties instead.

For the cases where a URL is unavoidable, provide a function that constructs the
URL based on the required properties of the resource being referenced. For
example to build an image family source disk a function similar to this is to be
provided: `gcompute_address_self_link(name, region, project)`.


## Folder Locations & Usage

Definitions:

  - Folders:
    * Folders suffixed with `...` means folder and its children.
    * Folders without `...` represents the folder only.
  - _Provider Independent_: Provider specific code, files or especializations
    **must not** be placed in this area
  - _Product Independent_: Product specific code, files or especializations
    **must not** be placed in this area

### Folders

  - `api`: Holds all definitions of objects used for defining products,
    serialization and deserialization
      * product independent
      * provider independent
  - `products/<product>/...`: Holds all definitions, files, helpers, etc
    required to build the product, for all providers
      * `api.yaml`
          - provider independent
      * `examples/<provider>`
      * `helpers/`
      * `files/`
  - `provider`: All core provider features
      * product independent
      * provider independent
  - `provider/<provider>/...`: Holds all provider specific code
      * product independent
  - `templates`: Holds all global templates
      * product independent
      * provider independent
  - `templates/<provider>/...`: Holds all templates for the provider
      * product independent

> Corollary: To completely remove (or add) a product to Magic Modules, removing
> (or adding) the `products/<product>/...` folder is all it takes.


## Code Style

  - All Ruby code **must** strictly abide by Rubocop standards.
    _Specializations to `.rubocop.yml` **should** be avoided at all costs._
  - All rspec code **must** strictly abive by
    [rspec standards][rspec-style-guide]
  - Line Lengths **must** be:
      * Ruby: 80 chars
      * YAML:
          - Magic Modules YAML: 80 chars
          - Product specific YAML: 80 chars
          - Provider specific YAML: up to provider standards
      * Markdown: 80 chars


## Templates

Template embedded code **should** be as simple as possible. Avoid complicated
logic in the templates because:

  1. It is harder to read and maintain
  2. It cannot be easily unit tested

It is best to only use simple iterators and provider specific functions.  If you
have a complicated logic (e.g. 'build a map between URL and properties') it is
best to put this functionality in the `provider/<provider>.rb` class.  All
methods in the provider class are directly available to the template (e.g. the
`indent()` function).

If the function is useful for all providers then consider adding it to the
global (core) provider namespace.


## Exceptions

Any _permanent_ exceptions to these rules **must** be thoroughly documented in
the code. If a longer discussion is required and becomes beyond the code where
it lives an issue **must** be created and referenced for context.

For any _temporary_ exceptions a tracking issue **must** be filed and added as a
"TODO" in the code for future fixing.

Exceptions have to have the LGTM (approval) core team member, or from another
core team member if the exception is introduced by a core team member.


## Ruby Best Practices & Style Guide

To make the code easier to be maintained Magic Modules strives to have a
consistent style guide and general rules outside regular Ruby style guides.
These rules **must** be followed.

### Inlined Ruby code to follow the Ruby style guide

Ruby code **must** observe the rules in the
[Ruby style guide][ruby-style-guide]. They help build better Ruby code.

### Keep ERB template code < 80 characters

We want to keep the inline Ruby following general Ruby guidelines (<80 chars)
and at the same time not dictating how the target files should look like. For
example Go does not have a maximum line limit and it is up to the writer to
break at will (usually following some project wide guidelines set by
themselves). So to keep the best of both worlds:

*   All inline Ruby code **must** fit <80 characters
*   Generated code (and its non-inline Ruby text) is up to the writer's
    discretion

### Keep `<% ... -%>` in the same line if it fits

**Good**

    <%= lines(indent('hello', 10)) -%>

**Bad**

    <%=
      lines(indent('hello', 10))
    -%>

### Do not use `\n` to format and introduce new lines.

Use [`lines()`](#lines) or [`lines_before()`](#lines-before) functions instead.
This both helps avoid tracking source of newlines and accidentally add more/less
spaces than needed

**Good**

    <%= lines("this has 3 new lines afterwards" + some_function(), 3) -%>

**Bad**

    <%= "this has 3 newlines afterwards" + some_function() + "\n\n\n"" -%>

In the example above you don't know how many blank lines `some_function()` would
return.

### Do not `.join("\n")`. Use [`lines()`](#lines) instead.

To keep spotting where `\n` are added it cumbersome and error prone. Use
[`lines`](#lines) function to ensure that your array is properly formatted.

### Prefer `-%>` (does not add `\n`) to `%>`.

If you need to break line use [`lines()`](#lines) function instead.

**Good**

    <%= lines(my_var) -%>

**Bad**

    <%= my_var %>

Although it seems easier to not use the [`lines()`](#lines) function, it is
common to have formatting bugs due to mixing `-%>` and `%>`.

### Align `<%` and `-%>` to column 1

By rooting `<%` (and `-%>` if multiline) to column 1 it will avoid introducing
spurious spaces in the final output.

**Good**

    |----------
    <%=
      lines(["something goes here", "and here"])
    -%>
    |----------

Produces:

    |----------
    something goes here
    and here
    |----------

**Bad**

    |----------
       <%=
         lines(["something goes here", "and here"])
       -%>
    |----------

Produces:

    |----------
       something goes here
    and here
    |----------

If you intend to have spaces in the beginning it is easier to clearly show that,
e.g. with [`indent()`](#indent), so someone does not try to "fix" your code.

### Align `<%` and `-%>` vertically

**Good**

    <%
      my_code
      my_code
    -%>

**Bad**

    <%
      my_code
      my_code -%>

### Add a comment to conditionals to easily identify its structure.

It is common for a template to have various lines between start and end of a
block and tracking them can become hard in complex templates. Adding something
to remind you where they belong improves readability

    <% objects.each do |obj| -%>
    ...
    <%   if obj.virtual -%>
    ...
    ... many lines later
    ...
    <%   elsif !obj.virtual && obj.broken -%>
    ...
    ... many lines later
    ...
    <%   end # if obj.virtual -%>
    ...
    ... many lines later
    ...
    <% end # objects.each -%>

### Align block begin and end vertically.

It will most of time make it easier to see that a block depends on another. In
the example below it is easier to note the if is inside the `each` and there are
2 cascaded `if`s as well

**Good**

    <% objects.each do |obj| -%>
    ...
    <%   if obj.virtual -%>
    ...
    <%     if obj.input -%>
    ...
    ... many lines later
    ...
    <%     end # if obj.input -%>
    <%   end # if obj.virtual -%>
    ...
    ... many lines later
    ...
    <% end # objects.each -%>

**Bad**

    <% objects.each do |obj| -%>
    ...
    <% if obj.virtual -%>
    ...
    <% if obj.input -%>
    ...
    ... many lines later
    ...
    <% end # if obj.input -%>
    <% end # if obj.virtual -%>
    ...
    ... many lines later
    ...
    <% end # objects.each -%>

### Use [`format()`](#format) for input dependent (style violating) output

Rely on [`format()`](#format) function for generated code that depends on
variables or other input specific data that may affect the generated output.

**Good**

    <%=
      format([
        ["this will fit in 1 line. user = #{user.full_name}"],
        [
          "this will fit in 2 lines.",
          user.full_name
        ],
        [
          "ugh, not even in 2 lines :(. let's make it 3 then",
          user.first_name,
          user.last_name
        ]
      ], 10, 40)
    -%>

**Bad**

    <%
      # Try to calculate the effective size of strings
      one_liner = "this will fit in 1 line. user = #{user.full_name}"
    -%>
    <% if one_liner.length + 10 < 40 -%>
    <%= lines(one_liner) -%>
    <% else # one_liner.length did not fit -%>
    <%
      # Try to calculate if fits in 2 lines
      ...
      ...
      ...
    -%>
    ...
    ...
    ...
    <% end # one_liner.length -%>

Please refer to [`format()`](#format) documentation for usage examples.

### Prefer nested [`indent()`](#indent) and [`indent_list()`](#indent-list) over calculating its relative displacements

**Good**

    <%=
      indent([
        'first level',
        indent([
          'second level',
          'also on second level'
        ], 2),
        'also on first level'
      ], 2)
    -%>

**Bad**

    <%=
      indent('first level', 2)
      indent('second level', 4)
      indent('also on second level', 4)
      indent('also on first level', 2)
    -%>

### Prefer functional over procedural for data processing

It makes it easier to read the code if you split the steps into functional
phases.

**Good**

    <%=
      # List all virtual properties that are not nested objects alphabetically
      lines(object.all_properties.select(&:virtual)
                                 .reject { |p| p.is_a?(Api::Type::NestedObject) }
                                 .sort
                                 .map { |p| "- #{p.out_name}" })
    -%>

**Bad**

    <%
      # List all virtual properties that are not nested objects alphabetically
      my_properties = []
      object.all_proerties.each do |p|
        if p.virtual && !p.is_a?(Api::Type::NestedObject)
          my_properties << p
        end # p.virtual
      end # object...each
    %>
    <% my_properties.sort.each do |p| -%>
    - <%= p.out_name -%>
    <% end # my_properties.each -%>

Ruby's [Enumerable][ruby-enumerable] interface contains a list of methods that
can be used in such cases.


[readme]: README.md
[phylosophy]: docs/philosophy.md
[contrib]: CONTRIBUTING.md
[rspec-style-guide]: http://betterspecs.org
[ruby-style-guide]: https://github.com/bbatsov/ruby-style-guide 
[ruby-enumerable]: https://ruby-doc.org/core-2.5.0/Enumerable.html
