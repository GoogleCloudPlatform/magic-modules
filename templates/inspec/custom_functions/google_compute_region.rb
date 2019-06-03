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

# helper for returning a list of zone short names rather than fully qualified URLs e.g.
#   https://www.googleapis.com/compute/v1/projects/spaterson-project/zones/asia-east1-a
def zone_names
  return [] if !exists?
  @zones.map { |zone| zone.split('/').last }
end

def up?
  return false if !exists?
  @status == 'UP'
end