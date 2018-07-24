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

# A Puppet resource that specifies the credentials and authentication method to
# use when authorizing Google Cloud Platform API requests.
Puppet::Type.newtype(:gauth_credential) do
  @doc = "Authentication for Google Cloud Platform.

  gauth_credential { 'myuser':
    path     => '/private/my_credentials.json',
    provider => serviceaccount,
  }

  Then you add 'credential' to your type to have access to the object:

  gcompute_instance { 'myvm':
    ensure     => present,
    credential => 'myuser',
    ...
    ...
  }
  "

  feature :serviceaccount, 'Authenticate using a service_account.json file'
  feature :defaultuseraccount, 'Authenticate using user default credential'

  newparam :title, namevar: true
  newparam :path, feature: :serviceaccount
  newparam :scopes, feature: :serviceaccount

  def authorization
    provider.authorization
  end

  # Fetch the authorization for the resource passed as parameter.
  def self.fetch(resource)
    id = get_id(resource)
    Puppet::Type.type(:gauth_credential).instances.each do |entry|
      return entry.authorization if entry.title == id
    end
    raise ArgumentError, \
          "gauth_credential[#{id}] required to authenticate #{resource}"
  end

  private

  def self.get_id(resource)
    id = resource[:credential] if resource.class <= Hash
    id = resource.value(:credential) unless resource.class <= Hash
    raise ArgumentError, "#{resource} lacks 'credential' parameter" if id.nil?
    id
  end
end
