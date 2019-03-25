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
# product | resource | property | api.yaml (y/n) | tf (y/n) | ansible (y/n)
class CsvFormatterForMM
  RSpec::Core::Formatters.register self, :start, :example_passed, :example_failed, :example_pending,
                                   :stop

  def initialize(output)
    @output = output
    @rows = []
  end

  # Places in the CSV header
  def start(_start_notification)
    @output << ['Date',
                'Product', 'Resource', 'Property', 'api.yaml', 'terraform', 'ansible'].to_csv
  end

  # This property exists in api.yaml
  def example_passed(notification)
    add_row(test_information(notification).merge(pass: 1))
  end

  # This property does not exist in api.yaml
  def example_failed(notification)
    add_row(test_information(notification).merge(pass: 0))
  end

  # This test isn't being run.
  # Don't do anything.
  def example_pending; end

  def stop(_stop_notification)
    @rows.map { |r| info_to_csv(r) }
         .each { |r| @output << r }
  end

  private

  def test_information(notification)
    example_group = notification.example.metadata[:example_group]
    {
      date: Time.now.strftime('%Y-%m-%d'),
      product: example_group[:parent_example_group][:parent_example_group][:description],
      resource: example_group[:parent_example_group][:description],
      property: example_group[:description],
      provider: notification.example.metadata[:provider]
    }
  end

  def info_to_csv(test_info)
    [
      test_info[:date],
      test_info[:product], test_info[:resource], test_info[:property], test_info[:api],
      test_info[:terraform], test_info[:ansible]
    ].to_csv
  end

  def add_row(test_info)
    row = @rows.select { |r| %i[product resource property].all? { |v| r[v] == test_info[v] } }.first
    if row
      row[test_info[:provider]] = test_info[:pass]
    else
      test_info[test_info[:provider]] = test_info[:pass]
      %i[pass provider].each { |k| test_info.delete(k) }
      @rows.append(test_info)
    end
  end
end
