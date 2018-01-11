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

require 'google/container/network/get'
require 'google/container/property/boolean'
require 'google/container/property/cluster_name'
require 'google/container/property/string'
require 'google/hash_utils'
require 'puppet'
require 'yaml'

Puppet::Type.type(:gcontainer_kube_config).provide(:google) do
  mk_resource_methods

  def self.instances
    debug('instances')
    raise [
      '"puppet resource" is not supported at the moment:',
      'TODO(nelsonjr): https://goto.google.com/graphite-bugs-view?id=167'
    ].join(' ')
  end

  def exists?
    debug("exists? #{@property_hash[:ensure] == :present}")
    @property_hash[:ensure] == :present
  end

  def create
    debug('create')
    @created = true
    IO.write(@resource[:name], kube_config.to_yaml)
    @property_hash[:ensure] = :present
  end

  def destroy
    debug('destroy')
    @deleted = true
    File.delete(@resource[:name]) if File.exist?(@resource[:name])
    @property_hash[:ensure] = :absent
  end

  def flush
    debug('flush')
    # return on !@dirty is for aiding testing (puppet already guarantees that)
    return if @created || @deleted || !@dirty
  end

  private

  def fetch_auth(resource)
    self.class.fetch_auth(resource)
  end

  def self.fetch_auth(resource)
    Puppet::Type.type(:gauth_credential).fetch(resource)
  end

  def debug(message)
    puts("DEBUG: #{message}") if ENV['PUPPET_HTTP_VERBOSE']
    super(message)
  end

  def cluster
    @cluster_ref ||= begin
      id = @resource[:cluster].resource
      Google::ObjectStore.instance[:gcontainer_cluster].each do |entry|
        return entry if entry.title == id
      end
      raise ArgumentError, "gcontainer_cluster[#{id}] not found"
    end
  end

  # rubocop:disable Metrics/MethodLength
  def kube_config
    endpoint = URI.parse("https://#{cluster.exports[:endpoint]}")
    auth = fetch_auth(@resource).authorize(endpoint).last
    context = @resource[:context] || cluster.name

    {
      'apiVersion' => 'v1',
      'clusters' => [
        {
          'name' => context,
          'cluster' => {
            'certificate-authority-data' =>
              cluster.exports[:master_auth]['clusterCaCertificate'],
            'server' => endpoint.to_s
          }
        }
      ],
      'contexts' => [
        {
          'name' => context,
          'context' => {
            'cluster' => context,
            'user' => context
          }
        }
      ],
      'current-context' => context,
      'kind' => 'Config',
      'preferences' => {},
      'users' => [
        {
          'name' => context,
          'user' => {
            'auth-provider' => {
              'config' => {
                'access-token' => auth.token,
                'cmd-args' => 'config config-helper --format=json',
                'cmd-path' => '/usr/lib64/google-cloud-sdk/bin/gcloud',
                'expiry-key' => '{.credential.token_expiry}',
                'token-key' => '{.credential.access_token}'
              },
              'name' => 'gcp'
            },
            'username' => cluster.exports[:master_auth]['username'],
            'password' => cluster.exports[:master_auth]['password']
          }
        }
      ]
    }
    # rubocop:enable Metrics/MethodLength
  end
end
