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

require 'google/object_store'
require 'puppet'

# A dummy provider that tells the user he needs to specify a provider.
Puppet::Type.type(:gauth_credential).provide(:default) do
  def self.init
    @resource_collector = Google::ObjectStore.instance
  end

  init

  def self.instances
    @resource_collector[:gauth_credential].map(&:provider)
  end

  def self.prefetch(resources)
    resources.each do |title, _|
      raise "gauth_credential[#{title}] does not have a provider"
    end
  end

  def self.default?
    true
  end

  def authorization
    raise 'default credential does not provide authorization'
  end
end
