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
      # Returns a list of acceptable import id formats for a given resource.
      #
      # For instance, if the resource base url is:
      #   projects/{{project}}/global/networks
      #
      # It returns 3 formats:
      # a) self_link: projects/{{project}}/global/networks/{{name}}
      # b) short id: {{project}}/{{name}}
      # c) short id w/o defaults: {{name}}
      #
      # Fields with default values are `project`, `region` and `zone`.
      def import_id_formats(resource)
        if resource.import_format.empty?
          underscored_base_url = resource.base_url
                                         .gsub(/{{[[:word:]]+}}/) do |field_name|
            Google::StringUtils.underscore(field_name)
          end

          # We assume that all resources have a name field
          id_formats = [underscored_base_url + '/{{name}}']
        else
          id_formats = resource.import_format
        end

        # short id: {{project}}/{{zone}}/{{name}}
        field_markers = id_formats[0].scan(/{{[[:word:]]+}}/)
        short_id_format = field_markers.join('/')

        # short id without fields with provider-level default: {{name}}
        field_markers.delete('{{project}}')
        field_markers.delete('{{region}}')
        field_markers.delete('{{zone}}')
        short_id_default_format = field_markers.join('/')

        id_formats + [short_id_format, short_id_default_format]
      end
    end
  end
end
