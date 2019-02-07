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

require 'provider/terraform_example'

module Provider
  # Code generator for a library converting terraform state to gcp objects.
  class TerraformObjectLibrary < Provider::TerraformExample
    def generate_resource(data)
      target_folder = data[:output_folder]

      generate_resource_file data.clone.merge(
        object: data[:object],
        default_template: 'templates/terraform/objectlib/base.go.erb',
        out_file: File.join(target_folder, "google/#{data[:object].name}.go")
      )
    end
  end
end
