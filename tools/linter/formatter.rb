# Copyright 2018 Google Inc.
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

require 'rspec'
require 'csv'

# This is an Rspec Formatter that's responsible for outputting the
# linter 'tests' out to a CSV format.
#
# This formatter only works on tests tagged with `property`
#
# Format:
# product | resource | property | api.yaml (y/n)
class CsvFormatterForMM
  RSpec::Core::Formatters.register self, :start, :example_passed, :example_failed, :example_pending

  def initialize(output)
    @output = output
  end

  # Places in the CSV header
  def start(_start_notification)
    @output << ['Product', 'Resource', 'Property', 'api.yaml'].to_csv
  end

  # This property exists in api.yaml
  def example_passed(notification)
    @output << info_to_csv(test_information(notification).merge(api_yaml: true))
  end

  # This property does not exist in api.yaml
  def example_failed(notification)
    @output << info_to_csv(test_information(notification).merge(api_yaml: false))
  end

  # This test isn't being run.
  # Don't do anything.
  def example_pending; end

  private

  def test_information(notification)
    example_group = notification.example.metadata[:example_group]
    {
      product: example_group[:parent_example_group][:parent_example_group][:description],
      resource: example_group[:parent_example_group][:description],
      property: example_group[:description]
    }
  end

  def info_to_csv(test_info)
    [
      test_info[:product], test_info[:resource], test_info[:property], test_info[:api_yaml]
    ].to_csv
  end
end
