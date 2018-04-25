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

# Add our google/ lib
$LOAD_PATH.unshift ::File.expand_path('../libraries', ::File.dirname(__FILE__))

require 'chef/resource'
require 'google/container/network/get'
require 'google/container/property/boolean'
require 'google/container/property/cluster_name'
require 'google/container/property/string'
require 'google/hash_utils'
require 'yaml'

module Google
  module GCONTAINER
    class KubeConfig < Chef::Resource
      resource_name :gcontainer_kubeconfig

      property :cluster,
               [String, ::Google::Container::Data::ClusterNameRef],
               coerce: ::Google::Container::Property::ClusterNameRef.coerce,
               desired_state: true
      property :zone,
               ::String,
               coerce: ::Google::Container::Property::String.coerce,
               desired_state: true
      property :context,
               ::String,
               coerce: ::Google::Container::Property::String.coerce,
               desired_state: true
      property :credential, String, desired_state: false, required: true
      property :project, String, desired_state: false, required: true

      action :create do
        IO.write(@new_resource.name, kube_config(@new_resource).to_yaml)
      end

      action :delete do
        File.delete(@new_resource.name) if File.exist(@new_resource.name)
      end

      private

      action_class do
        def fetch_auth(resource)
          resource.resources("gauth_credential[#{resource.credential}]")
        end

        def fetch_cluster(resource)
          id = resource.cluster.resource
          Chef.run_context.resource_collection.each do |entry|
            return entry if entry.name == id
          end
          raise ArgumentError, "gcontainer_cluster[#{id}] not found"
        end

        # rubocop:disable Metrics/MethodLength
        def kube_config(resource)
          cluster = fetch_cluster(resource)
          endpoint = URI.parse("https://#{cluster.exports[:endpoint]}")
          auth = fetch_auth(resource).authorization.authorize(endpoint).last
          context = resource.context || cluster.name
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
    end
  end
end
