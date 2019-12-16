# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'provider/abstract_core'

module Provider
  class Terraform < Provider::AbstractCore
    # Functions to support 'terraform import'.
    module Import
      def import_id_formats_from_resource(resource)
        import_id_formats(resource.import_format, resource.identity, resource.base_url)
      end

      # Returns a list of import id formats for a given resource. If an id
      # contains provider-default values, this fn will return formats both
      # including and omitting the value.
      #
      # If a resource has an explicit import_format value set, that will be the
      # base import url used. Next, the values of `identity` will be used to
      # construct a URL. Finally, `{{name}}` will be used by default.
      #
      # For instance, if the resource base url is:
      #   projects/{{project}}/global/networks
      #
      # It returns 3 formats:
      # a) self_link: projects/{{project}}/global/networks/{{name}}
      # b) short id: {{project}}/{{name}}
      # c) short id w/o defaults: {{name}}
      def import_id_formats(import_format, identity, base_url)
        if import_format.nil? || import_format.empty?
          underscored_base_url = base_url.gsub(
            /{{[[:word:]]+}}/, &:underscore
          )

          if identity.nil? || identity.empty?
            id_formats = [underscored_base_url + '/{{name}}']
          else
            identity_path = identity.map { |v| "{{#{v.name.underscore}}}" }.join('/')
            id_formats = [underscored_base_url + '/' + identity_path]
          end
        else
          id_formats = import_format
        end

        # short id: {{project}}/{{zone}}/{{name}}
        field_markers = id_formats[0].scan(/{{[[:word:]]+}}/)
        short_id_format = field_markers.join('/')

        # short ids without fields with provider-level defaults:

        # without project
        field_markers -= ['{{project}}']
        short_id_default_project_format = field_markers.join('/')

        # without project or location
        field_markers -= ['{{project}}', '{{region}}', '{{zone}}']
        short_id_default_format = field_markers.join('/')

        # Regexes should be unique and ordered from most specific to least specific
        # We sort by number of `/` characters (the standard block separator)
        # followed by number of variables (`{{`) to make `{{name}}` appear last.
        if id_formats[0].include?('%')
          # If the id format can include `/` characters we cannot allow short forms such as:
          # `{{project}}/{{%name}}` as there is no way to differentiate between
          # project-name/resource-name and resource-name/with-slash
          return \
            id_formats.uniq.reject(&:empty?).sort_by { |i| [i.count('/'), i.count('{{')] }.reverse
        end

        (id_formats + [short_id_format, short_id_default_project_format, short_id_default_format])
          .uniq.reject(&:empty?).sort_by { |i| [i.count('/'), i.count('{{')] }.reverse
      end
    end
  end
end
