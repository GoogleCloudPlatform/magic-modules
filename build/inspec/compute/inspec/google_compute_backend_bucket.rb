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

# A provider to manage Google Compute Engine resources.
class BackendBucket < Inspec.resource(1)

  name 'google_compute_backendbucket'
  desc 'BackendBucket'
  supports platform: 'gcp2'

  attr_reader :bucket_name
  attr_reader :creation_timestamp
  attr_reader :description
  attr_reader :enable_cdn
  attr_reader :id

  def base
    'https://www.googleapis.com/compute/v1/'
  end

  def url
    'projects/{{project}}/global/backendBuckets/{{name}}'
  end

  def initialize(params) 
    @fetched = fetch_resource(params, 'compute#backendBucket')
    parse unless @fetched.nil?
  end  

  def parse
    @bucket_name = @fetched['bucketName']
    @creation_timestamp = @fetched['creationTimestamp']
    @description = @fetched['description']
    @enable_cdn = @fetched['enableCdn']
    @id = @fetched['id']
  end

  def exists?
    !@fetched.nil?
  end

  def fetch_resource(params, kind)
    get_request = inspec.backend.fetch(base, url, params)
    return_if_object get_request.send, kind, true
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
      'https://www.googleapis.com/compute/v1/',
      expand_variables(
        'projects/{{project}}/global/backendBuckets/{{name}}',
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


end
