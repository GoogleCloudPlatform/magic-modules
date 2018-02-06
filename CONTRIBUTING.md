# How to become a contributor and submit your own code

## Templating Style Guide

* In-lined Ruby code should follow the [Ruby style guide][ruby-style-guide]
* All ERB template files should have line length < 80 characters whenever
  possible.
* Generated files should have line length < 80 characters (when possible)

* If in-lined Ruby cannot fit in one line, indent the block as followed:
  ```ruby
    <%
      code
    -%>
  ```
* Any reused or complicated in-lined Ruby code should be added to the Ruby
  provider file.
* Use the lines() function to add newline characters.
  Do not manually add the "\n" character in a template or Ruby code.
  This helps avoid adding an accidental second newline.
* When possible, align all "<%" and "<%=" tags to column 1 (no spaces before).
* Align all '<%'/'<%=' to '-%>'/'%>' when possible.
  **Good**
  ```ruby
  <%
    code
  -%>
  ```
  **Bad**
  ```ruby
  <%
    code -%>
  ```
* In ERB, if/else/end tags will probably be separated by many lines.
  Please add a comment by else/end tags to signify which if/loop branch these
  refer to.
  ```ruby
    <% if true -%>
    <% end # if true -%>
  ```
* There may be cases where filling ERB tags will result in > 80 characters
  in some, but not all situations.
  The format() function allows you to choose different formats of
  the same line depending on which fits within 80 characters.
  ```ruby
  <%=
    format([
      ["#{first_name} and #{last_name} are < 80 characters"],
      [
        first_name,
        "is quite long, so let's put it on it's own line #{last_name}"
      ]
    ])
  -%>
  ```
### Separating into Multiple Files
* If multiple providers are sharing code, place that code in a separate file
  in templates/.
* If all products in a single provider share a piece of fully encapsulated and
  complicated code, place it in a separate file under templates/<provider>/
* Extra code for a single product should be placed under
  products/<product>/helpers
* Separate files should be added to the provider template using include()
  and/or compile()

## Contributor License Agreements

We'd love to accept your sample apps and patches! Before we can take them, we
have to jump a couple of legal hurdles.

Please fill out either the individual or corporate Contributor License Agreement
(CLA).

  * If you are an individual writing original source code and you're sure you
    own the intellectual property, then you'll need to sign an [individual CLA]
    (http://code.google.com/legal/individual-cla-v1.0.html).
  * If you work for a company that wants to allow you to contribute your work,
    then you'll need to sign a [corporate CLA]
    (http://code.google.com/legal/corporate-cla-v1.0.html).

Follow either of the two links above to access the appropriate CLA and
instructions for how to sign and return it. Once we receive it, we'll be able to
accept your pull requests.

## Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
2. The repo owner will respond to your issue promptly.
3. If your proposed change is accepted, and you haven't already done so, sign a
   Contributor License Agreement (see details above).
4. Fork the desired repo, develop and test your code changes.
5. Ensure that your code is clear and comprehensible.
6. Ensure that your code has an appropriate set of unit tests which all pass.
7. Run `rake test` to verify that all unit tests and linters pass.
8. Submit a pull request.

[ruby-style-guide][https://github.com/bbatsov/ruby-style-guide]
