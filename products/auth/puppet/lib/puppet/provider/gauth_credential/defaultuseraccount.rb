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

require 'google/authorization'
require 'google/object_store'
require 'puppet'

# A Puppet provider that implements authenticating requests using a Google Cloud
# Platform default user account.
Puppet::Type.type(:gauth_credential).provide(:defaultuseraccount) do
  has_feature :defaultuseraccount

  def self.init
    @resource_collector = Google::ObjectStore.instance
  end

  init

  def self.instances
    @resource_collector[:gauth_credential].select do |resource|
      resource.provider == 'defaultuseraccount'
    end.map(&:provider)
  end

  def self.prefetch(resources)
    resources.each do |title, resource|
      resource.provider = new(title: title,
                              provider: :defaultuseraccount)
      debug "Created resource #{resource}"
      Google::ObjectStore.instance.add(:gauth_credential, resource)
    end
  end

  def authorization
    debug "Acquiring authorization for #{@resource}"
    Google::Authorization.new.from_user_credential!
  end
end
