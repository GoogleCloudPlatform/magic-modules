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

def handle_auto_to_custom_change
  # We allow changing the auto_create_subnetworks from true => false
  # (which will make the network going from Auto to Custom)
  auto_change = @dirty[:auto_create_subnetworks]
  raise 'Cannot convert a network from Custom back to Auto' \
    if auto_change[:from] == false && auto_change[:to] == true
  # TODO(nelsonjr): Enable converting from Auto => Custom via call to
  # special method URL. See tracking work item:
  # https://bugzilla.graphite.cloudnativeapp.com/show_bug.cgi?id=174
  raise [
    'Conversion from Auto to Custom not implemented yet.',
    'See', ['https://bugzilla.graphite.cloudnativeapp.com',
            'show_bug.cgi?id=174'].join('/'),
    'for more details'
  ].join(' ')
end
