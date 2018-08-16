# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Public: Authorizes access to Google API objects.
#
# Examples
#
#   * Uses user credential stored in ~/.config/gcloud
#
#     api = Google::Authorization.new
#         .for!('https://www.googleapis.com/auth/compute.readonly')
#         .from_user_credential!
#         .authorize Google::Apis::ComputeV1::ComputeService.new
#
#   * Uses service account specified by the :file argument (in JSON format)
#
#     api = Google::Authorization.new
#         .for!('https://www.googleapis.com/auth/compute.readonly')
#         .from_service_account_json!(
#             File.join(File.expand_path('~'), "my_account.json"))
#         .authorize Google::Apis::ComputeV1::ComputeService.new
#
# TODO(nelsona): Add support gcloud's beta "app default credential"

begin
  require 'googleauth'
rescue LoadError
  raise 'Please install the googleauth gem'
end
require 'json'
require 'net/http'

# Google authorization handler module.
module Google
  # A helper class to determine if we have Ruby 2
  class Ruby
    def self.two?
      Gem::Version.new(RUBY_VERSION.dup) >= Gem::Version.new('2.0.0')
    end

    def self.ensure_two!
      callee = caller(1..1).first[/`([^']*)'/, 1]
      raise "Ruby ~> 2.0.0 required for '#{callee}'" unless Ruby.two?
    end
  end

  begin
    require 'google/api_client/client_secrets' if Google::Ruby.two?
  rescue LoadError
    raise 'Please install the google-api-client gem'
  end
  # A class to aquire credentials and authorize Google API calls.
  class Authorization
    def initialize
      @authorization = nil
      @scopes = []
    end

    def authorize(obj)
      raise ArgumentError, 'A from_* method needs to be called before' \
        unless @authorization

      if obj.class <= URI::HTTPS || obj.class <= URI::HTTP
        authorize_uri obj
      elsif obj.class < Net::HTTPRequest
        authorize_http obj
      else
        obj.authorization = @authorization
        obj
      end
    end

    def for!(*scopes)
      @scopes = scopes
      self
    end

    def from_user_credential!
      Google::Ruby.ensure_two! # TODO(nelsona): Ensure we can run with Ruby 1.9
      hash = make_secrets_hash(find_credential)
      @authorization = Google::APIClient::ClientSecrets.new(hash)
                                                       .to_authorization
      self
    end

    def from_service_account_json!(service_account_file)
      raise 'Missing argument for scopes' if @scopes.empty?
      @authorization = Google::Auth::ServiceAccountCredentials.make_creds(
        json_key_io: File.open(service_account_file),
        scope: @scopes
      )
      self
    end

    def from_application_default_credentials!
      raise NotImplementedError, ':application_default_credentials'
    end

    private

    def authorize_uri(obj)
      http = Net::HTTP.new(obj.host, obj.port)
      http.use_ssl = obj.instance_of?(URI::HTTPS)
      http.verify_mode = OpenSSL::SSL::VERIFY_PEER
      [http, authorize_http(Net::HTTP::Get.new(obj.request_uri))]
    end

    def authorize_http(req)
      req.extend TokenProperty
      auth = {}
      @authorization.apply!(auth)
      req['Authorization'] = auth[:authorization]
      req.token = auth[:authorization].split(' ')[1]
      req
    end

    def credentials
      JSON.parse(File.read(File.join(ENV['HOME'], '.config', 'gcloud',
                                     'credentials')))
    end

    def find_credential
      credentials['data'].each do |entry|
        return entry['credential'] if entry['credential']['_class'] == 'OAuth2Credentials'
      end

      raise "Credential not found in '#{file}'"
    end

    def make_secrets_hash(cred)
      {
        'installed' => {
          'client_id' => cred['client_id'],
          'client_secret' => cred['client_secret'],
          'refresh_token' => cred['refresh_token']
        }
      }
    end
  end

  # Extension methods to enable retrieving the authentication token.
  module TokenProperty
    attr_reader :token
    attr_writer :token
  end
end
