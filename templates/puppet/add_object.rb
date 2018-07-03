<% if false # the license inside this if block assertains to this file -%>
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
<% end -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :ruby) -%>

require 'puppet'

# Constructs a referenceable object without that object needing to
# be defined by puppet.  Inserts the object into the local state so
# that resource references work without further modification.
class ReferenceableObject
  def initialize(attrs)
    @attrs = attrs
  end

  def method_missing(method_name, *_args, &_blk)
    if @attrs.key? method_name
      @attrs[method_name]
    else
      super
    end
  end

  def respond_to_missing?(method_name, _pvt = false)
    @attrs.key?(method_name) || super
  end

  def exports
    @attrs
  end
end

Puppet::Functions.create_function(:gcompute_external_resource) do
  dispatch :gcompute_external_resource do
    param 'String', :type
    param 'Hash', :attributes
  end

  def gcompute_external_resource(type, attributes)
    attributes['title'] = attributes['name'] unless attributes['name'].nil?
    sym_attrs = {}
    attributes.each { |k, v| sym_attrs[k.to_sym] = v }
    Google::ObjectStore.instance.add(type.to_sym,
                                     ReferenceableObject.new(sym_attrs))
  end
end
