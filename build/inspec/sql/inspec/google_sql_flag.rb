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

# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
#
# ----------------------------------------------------------------------------
#
#     This file is automatically generated by Magic Modules and manual
#     changes will be clobbered when the file is regenerated.
#
#     Please read more about how to change this file in README.md and
#     CONTRIBUTING.md located at the root of this package.
#
# ----------------------------------------------------------------------------

# Add our google/ lib
$LOAD_PATH.unshift ::File.expand_path('../libraries', ::File.dirname(__FILE__))

require 'google/hash_utils'
require 'inspec/resource'

# A provider to manage Google Cloud SQL resources.
class Flag < Inspec.resource(1)

  name 'google_sql_flag'
  desc 'Flag'
  supports platform: 'gcp2'

  attr_reader :allowed_string_values
  attr_reader :applies_to
  attr_reader :max_value
  attr_reader :min_value
  attr_reader :name
  attr_reader :requires_restart
  attr_reader :type

  def base
    'https://www.googleapis.com/sql/v1beta4/'
  end

  def url
    'flags'
  end

  def initialize(params) 
    @fetched = fetch_wrapped_resource(params, 'sql#flag', 'sql#flagsList', 'items')
    parse unless @fetched.nil?
  end

  def parse
    @allowed_string_values = @fetched['allowedStringValues']
    @applies_to = @fetched['appliesTo']
    @max_value = @fetched['maxValue']
    @min_value = @fetched['minValue']
    @name = @fetched['name']
    @requires_restart = @fetched['requiresRestart']
    @type = @fetched['type']
  end

  def exists?
    !@fetched.nil?
  end

  def fetch_resource(params, kind)
    get_request = inspec.backend.fetch(base, url, params)
    return_if_object get_request.send, kind, true
  end



  def fetch_wrapped_resource(params, kind, wrap_kind, wrap_path)
    result = fetch_resource(params, wrap_kind)
    return if result.nil? || !result.key?(wrap_path)
    result = unwrap_resource(result[wrap_path], params)
    return if result.nil?
    raise "Incorrect result: #{result['kind']} (expected #{kind})" \
      unless result['kind'] == kind
    result
  end


  def unwrap_resource(result, resource)
    query_predicate = unwrap_resource_filter(resource)
    matches = result.select do |entry|
      query_predicate.all? do |k, v|
        entry[k.id2name] == v
      end
    end
    raise "More than 1 result found: #{matches}" if matches.size > 1
    return if matches.empty?
    matches.first
  end

  def self.raise_if_errors(response, err_path, msg_field)
    errors = Google::HashUtils.navigate(response, err_path)
    raise_error(errors, msg_field) unless errors.nil?
  end

  def self.raise_error(errors, msg_field)
    raise IOError, ['Operation failed:',
                    errors.map { |e| e[msg_field] }.join(', ')].join(' ')
  end

  def self.self_link(data)
    URI.join(
      'https://www.googleapis.com/sql/v1beta4/',
      expand_variables(
        'flags',
        data
      )
    )
  end

  def self_link(data)
    self.class.self_link(data)
  end


  # rubocop:disable Metrics/CyclomaticComplexity
  def self.return_if_object(response, kind, allow_not_found = false)
    raise "Bad response: #{response.body}" \
      if response.is_a?(Net::HTTPBadRequest)
    raise "Bad response: #{response}" \
      unless response.is_a?(Net::HTTPResponse)
    return if response.is_a?(Net::HTTPNotFound) && allow_not_found 
    return if response.is_a?(Net::HTTPNoContent)
    result = JSON.parse(response.body)
    raise_if_errors result, %w[error errors], 'message'
    raise "Bad response: #{response}" unless response.is_a?(Net::HTTPOK)
    result
  end
  # rubocop:enable Metrics/CyclomaticComplexity

  def return_if_object(response, kind, allow_not_found = false)
    self.class.return_if_object(response, kind, allow_not_found)
  end


  def unwrap_resource_filter(resource)
    self.class.unwrap_resource_filter(resource)
  end

  def self.unwrap_resource_filter(resource)
    {
      name: resource[:name]
    }
  end

end
