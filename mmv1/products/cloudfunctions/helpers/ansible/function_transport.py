# Copyright 2019 Google Inc.
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
# Verify that firewall network names match what the API returns

def encode_request(request, module):
    return_vals = {}
    for k, v in request.items():
        if v or v is False:
            return_vals[k] = v

    if module.params['trigger_http'] and not return_vals.get('httpsTrigger'):
        return_vals['httpsTrigger'] = {}

    return return_vals
